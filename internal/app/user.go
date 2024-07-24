package app

import (
	"context"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"sync"
)

type UserGenID struct {
	sync.RWMutex
	lastID int
}

func (id *UserGenID) GetID() int {
	id.Lock()
	defer id.Unlock()
	id.lastID++
	return id.lastID
}

// GetUserURLs application api for get total stored urls by user
func (app *ShortApp) GetUserURLs(ctx context.Context, userID int) []pkg.UserURL {
	urls := app.Store.GetUserURLs(ctx, userID)
	for idx, userURL := range urls {
		urls[idx].ShortURL = app.GetShortURL(userURL.Code)
	}
	return urls
}

// DeleteUserURLs application api delete user urls by slice of short codes
func (app *ShortApp) DeleteUserURLs(ctx context.Context, userID int, codes []string) {
	go func(codeShortURL []string) {
		app.Store.DeleteBatch(ctx, userID, codeShortURL)
	}(codes)
	logger.Log.Infof("mark as delete short urls %v", codes)
}
