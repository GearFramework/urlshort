package app

import (
	"encoding/json"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"log"
	"os"
	"sync"
)

const stepFlush = 10

type mapCodes map[string]string
type mapURLs map[string]string

type Storage struct {
	sync.RWMutex
	storageFilePath string
	codeByURL       mapCodes
	urlByCode       mapURLs
	flushCounter    int
}

func NewStorage(filePath string) *Storage {
	return &Storage{
		storageFilePath: filePath,
	}
}

func (s *Storage) initStorage() {
	s.codeByURL = make(mapCodes, stepFlush)
	s.urlByCode = make(mapURLs, stepFlush)
	s.flushCounter = stepFlush
}

func (s *Storage) loadShortlyURLs() error {
	file, err := os.OpenFile(s.storageFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	if err = json.NewDecoder(file).Decode(&s.codeByURL); err != nil {
		return err
	}
	s.urlByCode = make(mapURLs, stepFlush)
	for url, code := range s.codeByURL {
		s.urlByCode[code] = url
	}
	s.flushCounter = s.count() + stepFlush
	return nil
}

func (s *Storage) getCode(url string) (string, bool) {
	code, ok := s.codeByURL[url]
	return code, ok
}

func (s *Storage) getURL(code string) (string, bool) {
	url, ok := s.urlByCode[code]
	return url, ok
}

func (s *Storage) add(url, code string) {
	s.codeByURL[url] = code
	s.urlByCode[code] = url
	logger.Log.Infof("Count urls in storage %d; counter to flush data %d", s.count(), s.flushCounter)
	if s.count() == s.flushCounter {
		if err := s.flush(); err != nil {
			logger.Log.Warn(err.Error())
		}
		s.flushCounter += stepFlush
	}
}

func (s *Storage) count() int {
	return len(s.codeByURL)
}

func (s *Storage) flush() error {
	logger.Log.Infoln("Flush storage to " + s.storageFilePath)
	if s.codeByURL == nil {
		return nil
	}
	data, err := json.Marshal(&s.codeByURL)
	if err != nil {
		return err
	}
	return os.WriteFile(s.storageFilePath, data, 0666)
}

func (s *Storage) reset() {
	s.clear()
	logger.Log.Infoln("Reset storage")
	if err := os.Remove(s.storageFilePath); err != nil {
		log.Println(err.Error())
	}
}

func (s *Storage) clear() {
	logger.Log.Infoln("Clear storage")
	for url, code := range s.codeByURL {
		delete(s.codeByURL, url)
		delete(s.urlByCode, code)
	}
}
