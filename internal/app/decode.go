package app

import (
	"errors"
)

func (app *ShortApp) DecodeURL(shortURL string) (string, error) {
	url, exists := app.urls[shortURL]
	if !exists {
		return "", errors.New("invalid short url")
	}
	return url, nil
}
