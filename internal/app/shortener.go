package app

import (
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/GearFramework/urlshort/internal/pkg/storage/db"
	"github.com/GearFramework/urlshort/internal/pkg/storage/file"
	"github.com/GearFramework/urlshort/internal/pkg/storage/mem"
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
	app.Store = app.factoryStorage()
	if err := app.Store.InitStorage(); err != nil {
		return err
	}
	return nil
}

func (app *ShortApp) factoryStorage() pkg.Storable {
	if app.Conf.DatabaseDSN != "" {
		return db.NewStorage(&db.StorageConfig{
			ConnectionDSN:   app.Conf.DatabaseDSN,
			ConnectMaxOpens: 10,
		})
	} else if app.Conf.StorageFilePath != "" {
		return file.NewStorage(&file.StorageConfig{
			StorageFilePath: app.Conf.StorageFilePath,
		})
	}
	return mem.NewStorage()
}

func (app *ShortApp) AddShortly(url, code string) {
	if err := app.Store.Insert(url, code); err != nil {
		logger.Log.Error(err.Error())
	}
}

func (app *ShortApp) ClearShortly() {
	if err := app.Store.Truncate(); err != nil {
		logger.Log.Error(err.Error())
	}
}

func (app *ShortApp) StopApp() {
	app.Store.Close()
}
