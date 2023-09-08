package app

import (
	"errors"
)

func DecodeURL(shortURL string) (string, error) {
	url, exists := urls[shortURL]
	if !exists {
		return "", errors.New("invalid short url")
	}
	return url, nil
}
