package pkg

import (
	"sync"
)

type BatchURLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResultBatchShort struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type APIShortener interface {
	EncodeURL(url string) string
	BatchEncodeURL(batch []BatchURLs) []ResultBatchShort
	DecodeURL(shortURL string) (string, error)
	AddShortly(url, code string)
}

type Storable interface {
	sync.Locker
	InitStorage() error
	GetCode(url string) (string, bool)
	GetCodeBatch(urls []string) map[string]string
	GetURL(code string) (string, bool)
	Insert(url, code string) error
	InsertBatch(batch [][]string) error
	Count() int
	Truncate() error
	Ping() error
	Close()
}
