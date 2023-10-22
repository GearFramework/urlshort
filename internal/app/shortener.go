package app

import (
	"context"
	"errors"
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/GearFramework/urlshort/internal/pkg/storage/db"
	"github.com/GearFramework/urlshort/internal/pkg/storage/file"
	"github.com/GearFramework/urlshort/internal/pkg/storage/mem"
	"io"
	"log"
	"math/rand"
	"time"
)

type ShortApp struct {
	Conf         *config.ServiceConfig
	Store        pkg.Storable
	flushCounter int
}

func NewShortener(conf *config.ServiceConfig) (*ShortApp, error) {
	shortener := ShortApp{Conf: conf}
	err := shortener.initApp()
	return &shortener, err
}

func (app *ShortApp) initApp() error {
	var err error
	app.Store, err = app.factoryStorage()
	return err
}

func (app *ShortApp) factoryStorage() (pkg.Storable, error) {
	if app.Conf.DatabaseDSN != "" {
		store := db.NewStorage(&db.StorageConfig{
			ConnectionDSN:   app.Conf.DatabaseDSN,
			ConnectMaxOpens: 10,
		})
		if err := app.isValidStorage(store); err == nil {
			log.Println("Use database urls storage")
			return store, nil
		}
	} else if app.Conf.StorageFilePath != "" {
		store := file.NewStorage(&file.StorageConfig{
			StorageFilePath: app.Conf.StorageFilePath,
		})
		err := app.isValidStorage(store)
		if err == nil || errors.Is(err, io.EOF) {
			log.Println("Use file urls storage")
			return store, nil
		}
	}
	log.Println("Use in memory urls storage")
	store := mem.NewStorage()
	if err := store.InitStorage(); err != nil {
		return nil, err
	}
	return store, nil
}

func (app *ShortApp) isValidStorage(store pkg.Storable) error {
	if err := store.InitStorage(); err != nil {
		return err
	}
	return store.Ping()
}

func (app *ShortApp) AddShortly(ctx context.Context, url, code string) {
	if err := app.Store.Insert(ctx, url, code); err != nil {
		logger.Log.Error(err.Error())
	}
}

func (app *ShortApp) ClearShortly() {
	if err := app.Store.Truncate(); err != nil {
		logger.Log.Error(err.Error())
	}
}

func (app *ShortApp) getRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = alphabet[rnd.Intn(lenAlpha)]
	}
	return string(b)
}

func (app *ShortApp) StopApp() {
	app.Store.Close()
}
