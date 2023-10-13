package app

import (
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/GearFramework/urlshort/internal/pkg/storage"
	"io"
)

type ShortApp struct {
	Conf         *config.ServiceConfig
	store        *Storage
	DB           *storage.Store
	flushCounter int
}

func NewShortener(conf *config.ServiceConfig) (*ShortApp, error) {
	shortener := ShortApp{Conf: conf}
	err := shortener.initApp()
	return &shortener, err
}

func (app *ShortApp) initApp() error {
	app.store = NewStorage(app.Conf.StorageFilePath)
	if err := app.store.loadShortlyURLs(); err != nil {
		if err != io.EOF {
			return err
		}
		app.store.initStorage()
	}
	return nil
}

func (app *ShortApp) AddShortly(url, code string) {
	app.store.add(url, code)
}

func (app *ShortApp) ClearShortly(hard bool) {
	app.store.clear()
	if hard {
		app.store.reset()
	}
}

func (app *ShortApp) StopApp() {
	if err := app.store.flush(); err != nil {
		logger.Log.Warn(err.Error())
	}
}
