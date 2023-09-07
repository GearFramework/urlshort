package app

import (
	"errors"
)

func (app *ShortApp) DecodeURL(shortURL string) (string, error) {
	url, exists := app.urls[shortURL]
	//TODO: Роберт М. ("Чистый код"): "Утверждения в условиях читаются лучше, чем отрицания"
	if !exists {
		return "", errors.New("invalid short url")
	}
	return url, nil
}
