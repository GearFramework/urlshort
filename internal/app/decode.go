package app

import (
	"context"
	"errors"
)

func (app *ShortApp) DecodeURL(ctx context.Context, code string) (string, error) {
	app.Store.Lock()
	defer app.Store.Unlock()
	url, exists := app.Store.GetURL(ctx, code)
	if !exists {
		return "", errors.New("invalid short url " + code)
	}
	return url, nil
}
