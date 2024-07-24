package app

import (
	"context"
	"errors"
)

var ErrShortURLIsDeleted = errors.New("short url is deleted")

// DecodeURL application api for return original url by code
func (app *ShortApp) DecodeURL(ctx context.Context, code string) (string, error) {
	app.Store.Lock()
	defer app.Store.Unlock()
	url, exists := app.Store.GetURL(ctx, code)
	if !exists {
		return "", errors.New("invalid short url " + code)
	}
	if url.IsDeleted {
		return "", ErrShortURLIsDeleted
	}
	return url.URL, nil
}
