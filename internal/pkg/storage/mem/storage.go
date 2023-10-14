package mem

import (
	"sync"
)

type mapCodes map[string]string
type mapURLs map[string]string

type Storage struct {
	sync.RWMutex
	codeByURL mapCodes
	urlByCode mapURLs
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) InitStorage() error {
	s.codeByURL = make(mapCodes, 10)
	s.urlByCode = make(mapURLs, 10)
	return nil
}

func (s *Storage) Close() {
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
	return nil
}

func (s *Storage) Count() int {
	return len(s.codeByURL)
}

func (s *Storage) Truncate() error {
	for url, code := range s.codeByURL {
		delete(s.codeByURL, url)
		delete(s.urlByCode, code)
	}
	return nil
}

func (s *Storage) Ping() error {
	return nil
}