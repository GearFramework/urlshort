package storage

import (
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"sync"
)

type Storage interface {
	GetUrl(code string) string
	GetCode(url string) string
	Ping() error
}

type Store struct {
	sync.RWMutex
	connection *StorageConnection
}

func NewStorage(config *StorageConfig) *Store {
	return &Store{
		connection: NewConnection(config),
	}
}

func (store *Store) InitStorage() error {
	return store.connection.Open()
}

func (store *Store) Close() {
	store.connection.Close()
}

func (store *Store) Truncate() error {
	_, err := store.connection.DB.Exec(`
		TRUNCATE urls.shortly RESTART IDENTITY
	`)
	if err != nil {
		return err
	}
	return nil
}

func (store *Store) Count() int {
	var count int
	if err := store.connection.DB.Get(&count, `
        SELECT COUNT(*) AS total_items
          FROM urls.shortly
    `); err != nil {
		logger.Log.Error(err.Error())
		return 0
	}
	return count
}

func (store *Store) Insert(url, code string) error {
	_, err := store.connection.DB.Exec(`
		INSERT INTO urls.shortly (url, code) 
		VALUES ($1, $2)
	`, url, code)
	return err
}

func (store *Store) GetURL(code string) (string, bool) {
	var url string
	err := store.connection.DB.Get(&url, `
		SELECT url 
		  FROM urls.shortly 
		 WHERE code = $1
 	`, code)
	return url, err == nil
}

func (store *Store) GetCode(url string) (string, bool) {
	var code string
	err := store.connection.DB.Get(&code, `
		SELECT code 
		  FROM urls.shortly 
		 WHERE url = $1
 	`, url)
	return code, err == nil
}

func (store *Store) Ping() error {
	return store.connection.Ping()
}
