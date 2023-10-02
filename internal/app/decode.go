package app

import (
	"errors"
)

func (app *ShortApp) DecodeURL(shortURL string) (string, error) {
	app.store.Lock()
	defer app.store.Unlock()
	url, exists := app.store.GetCode(shortURL)
	if !exists {
		return "", errors.New("invalid short url")
	}
	return url, nil
}
