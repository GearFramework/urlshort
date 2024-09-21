package pkg

import (
	"context"
	"sync"
)

// User params names
const (
	// UserIDParamName param name for access in cookie
	UserIDParamName = "userID"
)

// BatchURLs struct of batch urls
type BatchURLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ResultBatchShort struct of result batch urls
type ResultBatchShort struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserURL struct of user url
type UserURL struct {
	Code     string `json:"-"`
	ShortURL string `json:"short_url"`
	URL      string `json:"original_url"`
}

// ShortURL struct of state url
type ShortURL struct {
	URL       string `db:"url"`
	IsDeleted bool   `db:"is_deleted"`
}

// Stats struct of counts url and users in storage
type Stats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// APIShortener interface of applications shortener urls
type APIShortener interface {
	Auth(token string) (int, error)
	GenerateUserID() int
	CreateToken() (int, string, error)
	EncodeURL(ctx context.Context, userID int, url string) (string, bool)
	BatchEncodeURL(ctx context.Context, userID int, batch []BatchURLs) []ResultBatchShort
	DecodeURL(ctx context.Context, shortURL string) (string, error)
	AddShortly(ctx context.Context, UserID int, url, code string)
	GetUserURLs(ctx context.Context, userID int) []UserURL
	DeleteUserURLs(ctx context.Context, userID int, codes []string)
	GetShortURL(code string) string
	GetStats(ctx context.Context) *Stats
}

// Storable interface of storage urls
type Storable interface {
	sync.Locker
	InitStorage() error
	GetCode(ctx context.Context, url string) (string, bool)
	GetCodeBatch(ctx context.Context, urls []string) map[string]string
	GetURL(ctx context.Context, code string) (ShortURL, bool)
	GetMaxUserID(ctx context.Context) (int, error)
	GetUserURLs(ctx context.Context, userID int) []UserURL
	Insert(ctx context.Context, userID int, url, code string) error
	InsertBatch(ctx context.Context, userID int, batch [][]string) error
	DeleteBatch(ctx context.Context, userID int, batch []string)
	Count(ctx context.Context) int
	GetUniqueUsers(ctx context.Context) []int
	Truncate() error
	Ping() error
	Close()
}

// GeneratorID interface of generators user ID
type GeneratorID interface {
	GetID() int
}
