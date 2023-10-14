package file

import (
	"encoding/json"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"log"
	"os"
	"sync"
)

type mapCodes map[string]string
type mapURLs map[string]string

type Storage struct {
	sync.RWMutex
	Config       *StorageConfig
	codeByURL    mapCodes
	urlByCode    mapURLs
	flushCounter int
}

func NewStorage(config *StorageConfig) *Storage {
	return &Storage{
		Config: config,
	}
}

func (s *Storage) InitStorage() error {
	s.codeByURL = make(mapCodes, s.Config.FlushPerItems)
	s.urlByCode = make(mapURLs, s.Config.FlushPerItems)
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
	s.urlByCode = make(mapURLs, s.Config.FlushPerItems)
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

func (s *Storage) GetURL(code string) (string, bool) {
	url, ok := s.urlByCode[code]
	return url, ok
}

func (s *Storage) Insert(url, code string) error {
	s.codeByURL[url] = code
	s.urlByCode[code] = url
	if s.Count() == s.flushCounter {
		if err := s.flush(); err != nil {
			logger.Log.Warn(err.Error())
		}
		s.flushCounter += s.Config.FlushPerItems
	}
	return nil
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
		if err != nil {
			return err
		}
		defer fd.Close()
	}
	return nil
}
