package mem

import (
	"context"
	"sync"

	"github.com/GearFramework/urlshort/internal/pkg"
)

type mapCodes map[string]Codes

// Codes storage struct of urls
type Codes struct {
	Code      string
	UserID    int
	IsDeleted bool
}

type mapURLs map[string]string

// Storage in-memory storage
type Storage struct {
	sync.RWMutex
	codeByURL mapCodes
	urlByCode mapURLs
	users     map[int]int
}

var lastUserID int = 0

// NewStorage return new in-memory storage
func NewStorage() *Storage {
	return &Storage{}
}

// InitStorage initialize in-memory storage
func (s *Storage) InitStorage() error {
	s.codeByURL = make(mapCodes, 10)
	s.urlByCode = make(mapURLs, 10)
	return nil
}

// Close storage
func (s *Storage) Close() {
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
	if lastUserID < userID {
		lastUserID = userID
	}
	s.incUserStat(userID, 1)
	return nil
}

// InsertBatch batch url into storage
func (s *Storage) InsertBatch(ctx context.Context, userID int, batch [][]string) error {
	for _, pack := range batch {
		s.codeByURL[pack[0]] = Codes{UserID: userID, Code: pack[1]}
		s.urlByCode[pack[1]] = pack[0]
	}
	if lastUserID < userID {
		lastUserID = userID
	}
	s.incUserStat(userID, len(batch))
	return nil
}

func (s *Storage) incUserStat(userID, added int) {
	if v, ok := s.users[userID]; ok {
		s.users[userID] = v + added
	} else {
		s.users[userID] = added
	}
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

// GetUniqueUsers return slice of unique user ID
func (s *Storage) GetUniqueUsers(ctx context.Context) (int, error) {
	return len(s.users), nil
}

// Count return count url in storage
func (s *Storage) Count(ctx context.Context) (int, error) {
	return len(s.codeByURL), nil
}

// Truncate clear storage
func (s *Storage) Truncate() error {
	for url, code := range s.codeByURL {
		delete(s.codeByURL, url)
		delete(s.urlByCode, code.Code)
	}
	return nil
}

// Ping storage
func (s *Storage) Ping() error {
	return nil
}
