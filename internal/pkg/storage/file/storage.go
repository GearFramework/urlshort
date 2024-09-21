package file

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"

	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
)

// Codes storage struct of urls
type Codes struct {
	Code      string `json:"code"`
	UserID    int    `json:"user_id"`
	IsDeleted bool   `json:"is_deleted"`
}

// Storage in-file storage
type Storage struct {
	sync.RWMutex
	Config       *StorageConfig
	codeByURL    map[string]Codes
	urlByCode    map[string]string
	flushCounter int
}

var lastUserID int = 0

// NewStorage return new in-file storage
func NewStorage(config *StorageConfig) *Storage {
	return &Storage{
		Config: config,
	}
}

// InitStorage initialize in-file storage
func (s *Storage) InitStorage() error {
	s.codeByURL = make(map[string]Codes, s.Config.FlushPerItems)
	s.urlByCode = make(map[string]string, s.Config.FlushPerItems)
	s.flushCounter = s.Config.FlushPerItems
	err := s.loadShortlyURLs()
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func (s *Storage) loadShortlyURLs() error {
	file, err := os.OpenFile(s.Config.StorageFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			log.Println(err.Error())
		}
	}()
	if err = json.NewDecoder(file).Decode(&s.codeByURL); err != nil {
		return err
	}
	s.urlByCode = make(map[string]string, s.Config.FlushPerItems)
	for url, data := range s.codeByURL {
		s.urlByCode[data.Code] = url
		if lastUserID < data.UserID {
			lastUserID = data.UserID
		}
	}
	s.flushCounter = s.Count(context.Background()) + s.Config.FlushPerItems
	return nil
}

// Close storage
func (s *Storage) Close() {
	if err := s.flush(); err != nil {
		logger.Log.Error(err.Error())
	}
}

// GetCode return code from storage by url
func (s *Storage) GetCode(ctx context.Context, url string) (string, bool) {
	data, ok := s.codeByURL[url]
	if ok {
		return data.Code, ok
	}
	return "", ok
}

// GetCodeBatch return batch of codes by slice urls
func (s *Storage) GetCodeBatch(ctx context.Context, batch []string) map[string]string {
	codes := map[string]string{}
	for _, url := range batch {
		if _, ok := s.codeByURL[url]; ok {
			codes[url] = s.codeByURL[url].Code
		}
	}
	return codes
}

// GetURL return url by code
func (s *Storage) GetURL(ctx context.Context, code string) (pkg.ShortURL, bool) {
	url, ok := s.urlByCode[code]
	short := pkg.ShortURL{}
	if ok {
		short.URL = url
		short.IsDeleted = s.codeByURL[url].IsDeleted
	}
	return short, ok
}

// GetMaxUserID get last user ID
func (s *Storage) GetMaxUserID(ctx context.Context) (int, error) {
	return lastUserID, nil
}

// GetUserURLs return slice urls of user
func (s *Storage) GetUserURLs(ctx context.Context, userID int) []pkg.UserURL {
	userURLs := []pkg.UserURL{}
	for url, userShortURL := range s.codeByURL {
		if userShortURL.UserID == userID {
			userURLs = append(userURLs, pkg.UserURL{Code: userShortURL.Code, URL: url})
		}
	}
	return userURLs
}

// Insert url into storage
func (s *Storage) Insert(ctx context.Context, userID int, url, code string) error {
	s.codeByURL[url] = Codes{UserID: userID, Code: code}
	s.urlByCode[code] = url
	var err error
	if s.mustFlush(ctx) {
		err = s.flush()
		s.flushCounter += s.Config.FlushPerItems
	}
	if lastUserID < userID {
		lastUserID = userID
	}
	return err
}

// InsertBatch batch url into storage
func (s *Storage) InsertBatch(ctx context.Context, userID int, batch [][]string) error {
	for _, pack := range batch {
		s.codeByURL[pack[0]] = Codes{UserID: userID, Code: pack[1]}
		s.urlByCode[pack[1]] = pack[0]
	}
	if s.mustFlush(ctx) {
		if err := s.flush(); err != nil {
			logger.Log.Warn(err.Error())
		}
		s.flushCounter += s.Config.FlushPerItems
	}
	if lastUserID < userID {
		lastUserID = userID
	}
	return nil
}

// DeleteBatch mark urls as deleted
func (s *Storage) DeleteBatch(ctx context.Context, userID int, batch []string) {
	s.Lock()
	defer s.Unlock()
	for _, code := range batch {
		url, ok := s.urlByCode[code]
		if ok && s.codeByURL[url].UserID == userID {
			short := s.codeByURL[url]
			short.IsDeleted = true
			s.codeByURL[url] = short
		}
	}
}

func (s *Storage) mustFlush(ctx context.Context) bool {
	return s.Count(ctx) == s.flushCounter
}

// GetUniqueUsers return slice of unique user ID
func (s *Storage) GetUniqueUsers(ctx context.Context) []int {
	s.Lock()
	defer s.Unlock()
	users := []int{}
	hash := map[int]int{}
	for _, code := range s.codeByURL {
		if _, ok := hash[code.UserID]; !ok {
			hash[code.UserID] = 1
			users = append(users, code.UserID)
		}
	}
	return users
}

// Count return count url in storage
func (s *Storage) Count(ctx context.Context) int {
	return len(s.codeByURL)
}

func (s *Storage) flush() error {
	if s.codeByURL == nil {
		return nil
	}
	data, err := json.Marshal(&s.codeByURL)
	if err != nil {
		return err
	}
	return os.WriteFile(s.Config.StorageFilePath, data, 0666)
}

// Truncate clear storage
func (s *Storage) Truncate() error {
	for url, code := range s.codeByURL {
		delete(s.codeByURL, url)
		delete(s.urlByCode, code.Code)
	}
	return s.flush()
}

// Ping storage
func (s *Storage) Ping() error {
	_, err := os.Stat(s.Config.StorageFilePath)
	if os.IsNotExist(err) {
		fd, err := os.OpenFile(s.Config.StorageFilePath, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		return fd.Close()
	}
	return nil
}
