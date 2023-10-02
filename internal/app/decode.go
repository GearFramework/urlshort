package app

import (
	"errors"
)

func (app *ShortApp) DecodeURL(code string) (string, error) {
	app.store.Lock()
	defer app.store.Unlock()
	url, exists := app.store.getURL(code)
	if !exists {
		return "", errors.New("invalid short url " + code)
	}
	return url, nil
}
