package mem

import (
	"context"
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

func (s *Storage) GetCode(ctx context.Context, url string) (string, bool) {
	code, ok := s.codeByURL[url]
	return code, ok
}

func (s *Storage) GetCodeBatch(ctx context.Context, batch []string) map[string]string {
	codes := map[string]string{}
	for _, url := range batch {
		if _, ok := s.codeByURL[url]; ok {
			codes[url] = s.codeByURL[url]
		}
	}
	return codes
}

func (s *Storage) GetURL(ctx context.Context, code string) (string, bool) {
	url, ok := s.urlByCode[code]
	return url, ok
}

func (s *Storage) Insert(ctx context.Context, url, code string) error {
	s.codeByURL[url] = code
	s.urlByCode[code] = url
	return nil
}

func (s *Storage) InsertBatch(ctx context.Context, batch [][]string) error {
	for _, pack := range batch {
		s.codeByURL[pack[0]] = pack[1]
		s.urlByCode[pack[1]] = pack[0]
	}
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
