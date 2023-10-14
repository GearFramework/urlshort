package pkg

import "sync"

type APIShortener interface {
	EncodeURL(url string) string
	DecodeURL(shortURL string) (string, error)
	AddShortly(url, code string)
}

type Storable interface {
	sync.Locker
	InitStorage() error
	GetCode(url string) (string, bool)
	GetURL(code string) (string, bool)
	Insert(url, code string) error
	Count() int
	Truncate() error
	Ping() error
	Close()
}
