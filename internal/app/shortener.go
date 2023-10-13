package app

import (
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/GearFramework/urlshort/internal/pkg/storage"
)

type ShortApp struct {
	Conf         *config.ServiceConfig
	Store        *storage.Store
	flushCounter int
}

func NewShortener(conf *config.ServiceConfig) (*ShortApp, error) {
	shortener := ShortApp{Conf: conf}
	err := shortener.initApp()
	return &shortener, err
}

func (app *ShortApp) initApp() error {
	app.Store = storage.NewStorage(&storage.StorageConfig{
		ConnectionDSN:   app.Conf.DatabaseDSN,
		ConnectMaxOpens: 10,
	})
	if err := app.Store.InitStorage(); err != nil {
		return err
	}
	return nil
}

func (app *ShortApp) AddShortly(url, code string) {
	app.Store.Add(url, code)
}

func (app *ShortApp) ClearShortly() {
	if err := app.Store.Truncate(); err != nil {
		logger.Log.Error(err.Error())
	}
}

func (app *ShortApp) StopApp() {
	app.Store.Close()
}
