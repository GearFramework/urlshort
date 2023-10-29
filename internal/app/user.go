package app

import (
	"context"
	"github.com/GearFramework/urlshort/internal/pkg"
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

func (app *ShortApp) GetUserURLs(ctx context.Context, userID int) []pkg.UserURL {
	urls := app.Store.GetUserURLs(ctx, userID)
	return urls
}
