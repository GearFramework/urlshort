package pkg

import (
	"context"
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
	EncodeURL(ctx context.Context, url string) (string, bool)
	BatchEncodeURL(ctx context.Context, batch []BatchURLs) []ResultBatchShort
	DecodeURL(ctx context.Context, shortURL string) (string, error)
	AddShortly(ctx context.Context, url, code string)
}

type Storable interface {
	sync.Locker
	InitStorage() error
	GetCode(ctx context.Context, url string) (string, bool)
	GetCodeBatch(ctx context.Context, urls []string) map[string]string
	GetURL(ctx context.Context, code string) (string, bool)
	Insert(ctx context.Context, url, code string) error
	InsertBatch(ctx context.Context, batch [][]string) error
	Count() int
	Truncate() error
	Ping() error
	Close()
}
