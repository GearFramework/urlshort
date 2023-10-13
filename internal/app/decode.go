package app

import (
	"errors"
)

func (app *ShortApp) DecodeURL(code string) (string, error) {
	app.Store.Lock()
	defer app.Store.Unlock()
	url, exists := app.Store.GetURL(code)
	if !exists {
		return "", errors.New("invalid short url " + code)
	}
	return url, nil
}
