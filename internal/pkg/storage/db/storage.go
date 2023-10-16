package db

import (
	"context"
	"database/sql"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/jmoiron/sqlx"
	"sync"
)

type Storage struct {
	sync.RWMutex
	connection *StorageConnection
}

func NewStorage(config *StorageConfig) *Storage {
	return &Storage{
		connection: NewConnection(config),
	}
}

func (s *Storage) InitStorage() error {
	if err := s.connection.Open(); err != nil {
		return err
	}
	_, err := s.connection.DB.ExecContext(context.Background(), `
		CREATE SCHEMA IF NOT EXISTS urls
	`)
	if err != nil {
		return err
	}
	_, err = s.connection.DB.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS urls.shortly (
		    code VARCHAR(8),
			url VARCHAR(1024),
		    CONSTRAINT code_url PRIMARY KEY (code, url)
		)
	`)
	return err
}

func (s *Storage) Close() {
	s.connection.Close()
}

func (s *Storage) GetCode(url string) (string, bool) {
	var code string
	err := s.connection.DB.GetContext(context.Background(), &code, `
		SELECT code 
		  FROM urls.shortly 
		 WHERE url = $1
 	`, url)
	return code, err == nil
}

func (s *Storage) GetCodeBatch(batch []string) map[string]string {
	codes := map[string]string{}
	q, args, err := sqlx.In(`
		SELECT code, url 
		  FROM urls.shortly 
		 WHERE url IN (?)
 	`, batch)
	if err != nil {
		logger.Log.Error(err.Error())
		return codes
	}
	q = sqlx.Rebind(sqlx.DOLLAR, q)
	rows, err := s.connection.DB.QueryContext(context.Background(), q, args...)
	if err != nil {
		logger.Log.Error(err.Error())
		return codes
	}
	defer rows.Close()
	for rows.Next() {
		var code, url string
		err := rows.Scan(&code, &url)
		if err != nil {
			break
		}
		codes[url] = code
	}
	if err = rows.Err(); err != nil {
		logger.Log.Warn(err.Error())
	}
	return codes
}

func (s *Storage) GetURL(code string) (string, bool) {
	var url string
	err := s.connection.DB.GetContext(context.Background(), &url, `
		SELECT url 
		  FROM urls.shortly 
		 WHERE code = $1
 	`, code)
	return url, err == nil
}

func (s *Storage) Insert(url, code string) error {
	_, err := s.connection.DB.ExecContext(context.Background(), `
		INSERT INTO urls.shortly (url, code) 
		VALUES ($1, $2)
	`, url, code)
	return err
}

func (s *Storage) InsertBatch(batch [][]string) error {
	ctx := context.Background()
	var err error
	tx, err := s.connection.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO urls.shortly (url, code) 
		VALUES ($1, $2)
	`)
	if err != nil {
		return err
	}
	for _, data := range batch {
		_, err = stmt.ExecContext(ctx, data[0], data[1])
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Log.Error(err.Error())
			}
			return err
		}
	}
	return tx.Commit()
}

func (s *Storage) Count() int {
	var count int
	if err := s.connection.DB.GetContext(context.Background(), &count, `
        SELECT COUNT(*) AS total_items
          FROM urls.shortly
    `); err != nil {
		logger.Log.Error(err.Error())
		return 0
	}
	return count
}

func (s *Storage) Truncate() error {
	_, err := s.connection.DB.ExecContext(context.Background(), `
		TRUNCATE urls.shortly RESTART IDENTITY
	`)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Ping() error {
	return s.connection.Ping()
}
