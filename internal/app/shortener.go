// Package app for generating short urls
package app

import (
	"context"
	"errors"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/auth"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/GearFramework/urlshort/internal/pkg/storage/db"
	"github.com/GearFramework/urlshort/internal/pkg/storage/file"
	"github.com/GearFramework/urlshort/internal/pkg/storage/mem"
)

// ShortApp struct of application
type ShortApp struct {
	Conf  *config.ServiceConfig
	Store pkg.Storable
	GenID *UserGenID
}

// NewShortener make and return short urls application
func NewShortener(conf *config.ServiceConfig) (*ShortApp, error) {
	shortener := ShortApp{
		Conf: conf,
	}
	err := shortener.initApp()
	return &shortener, err
}

func (app *ShortApp) initApp() error {
	var err error
	app.Store, err = app.factoryStorage()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	maxUserID, err := app.Store.GetMaxUserID(ctx)
	if err != nil {
		return err
	}
	app.GenID = &UserGenID{lastID: maxUserID}
	return nil
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
		} else {
			log.Println(err.Error())
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

// Auth process authorization of user, return user ID, error if not authorised
func (app *ShortApp) Auth(token string) (int, error) {
	userID := auth.GetUserIDFromJWT(token)
	if userID == -1 {
		return userID, auth.ErrNeedAuthorization
	}
	if userID == 0 {
		return userID, auth.ErrInvalidAuthorization
	}
	return userID, nil
}

// CreateToken make access token for authorized user
func (app *ShortApp) CreateToken() (int, string, error) {
	userID := app.GenerateUserID()
	token, err := auth.BuildJWT(userID)
	return userID, token, err
}

// GenerateUserID make unique user ID
func (app *ShortApp) GenerateUserID() int {
	return app.GenID.GetID()
}

// AddShortly save url and short code in storage
func (app *ShortApp) AddShortly(ctx context.Context, userID int, url, code string) {
	if err := app.Store.Insert(ctx, userID, url, code); err != nil {
		logger.Log.Error(err.Error())
	}
}

// ClearShortly delete all short urls in storage
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

// StopApp running when application shut down
func (app *ShortApp) StopApp() {
	app.Store.Close()
}
