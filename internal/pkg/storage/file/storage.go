package file

import (
	"encoding/json"
	"errors"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"io"
	"log"
	"os"
	"sync"
)

type Storage struct {
	sync.RWMutex
	Config       *StorageConfig
	codeByURL    map[string]string
	urlByCode    map[string]string
	flushCounter int
}

func NewStorage(config *StorageConfig) *Storage {
	return &Storage{
		Config: config,
	}
}

func (s *Storage) InitStorage() error {
	s.codeByURL = make(map[string]string, s.Config.FlushPerItems)
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
	defer file.Close()
	if err = json.NewDecoder(file).Decode(&s.codeByURL); err != nil {
		return err
	}
	s.urlByCode = make(map[string]string, s.Config.FlushPerItems)
	for url, code := range s.codeByURL {
		s.urlByCode[code] = url
	}
	s.flushCounter = s.Count() + s.Config.FlushPerItems
	return nil
}

func (s *Storage) Close() {
	if err := s.flush(); err != nil {
		logger.Log.Error(err.Error())
	}
}

func (s *Storage) GetCode(url string) (string, bool) {
	code, ok := s.codeByURL[url]
	return code, ok
}

func (s *Storage) GetCodeBatch(batch []string) map[string]string {
	codes := map[string]string{}
	for _, url := range batch {
		if _, ok := s.codeByURL[url]; ok {
			codes[url] = s.codeByURL[url]
		}
	}
	return codes
}

func (s *Storage) GetURL(code string) (string, bool) {
	url, ok := s.urlByCode[code]
	return url, ok
}

func (s *Storage) Insert(url, code string) error {
	s.codeByURL[url] = code
	s.urlByCode[code] = url
	var err error
	if s.mustFlush() {
		err = s.flush()
		s.flushCounter += s.Config.FlushPerItems
	}
	return err
}

func (s *Storage) InsertBatch(batch [][]string) error {
	for _, pack := range batch {
		s.codeByURL[pack[0]] = pack[1]
		s.urlByCode[pack[1]] = pack[0]
	}
	if s.mustFlush() {
		if err := s.flush(); err != nil {
			logger.Log.Warn(err.Error())
		}
		s.flushCounter += s.Config.FlushPerItems
	}
	return nil
}

func (s *Storage) mustFlush() bool {
	return s.Count() == s.flushCounter
}

func (s *Storage) Count() int {
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

func (s *Storage) Truncate() error {
	for url, code := range s.codeByURL {
		delete(s.codeByURL, url)
		delete(s.urlByCode, code)
	}
	return s.flush()
}

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
