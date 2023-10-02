package app

import (
	"encoding/json"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"os"
	"sync"
)

const stepFlush = 2

type mapCodes map[string]string
type mapURLs map[string]string

type Storage struct {
	sync.RWMutex
	storageFilePath string
	codes           mapCodes
	urls            mapURLs
	flushCounter    int
}

func NewStorage(filePath string) *Storage {
	return &Storage{
		storageFilePath: filePath,
	}
}

func (s *Storage) initStorage() {
	s.codes = make(mapCodes, 10)
	s.urls = make(mapURLs, 10)
}

func (s *Storage) loadShortlyURLs() error {
	file, err := os.OpenFile(s.storageFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	if err = json.NewDecoder(file).Decode(&s.codes); err != nil {
		return err
	}
	s.urls = make(mapURLs, 10)
	for url, code := range s.codes {
		s.urls[code] = url
	}
	s.flushCounter = s.Count() + stepFlush
	return nil
}

func (s *Storage) GetCode(url string) (string, bool) {
	code, ok := s.codes[url]
	return code, ok
}

func (s *Storage) GetURL(code string) (string, bool) {
	url, ok := s.urls[code]
	return url, ok
}

func (s *Storage) Add(url, code string) {
	s.codes[url] = code
	s.urls[code] = url
	if s.Count() == s.flushCounter {
		if err := s.Flush(); err != nil {
			logger.Log.Warn(err.Error())
		}
		s.flushCounter += stepFlush
	}
}

func (s *Storage) Count() int {
	return len(s.codes)
}

func (s *Storage) Flush() error {
	if s.codes == nil {
		return nil
	}
	data, err := json.Marshal(&s.codes)
	if err != nil {
		return err
	}
	return os.WriteFile(s.storageFilePath, data, 0666)
}
