package app

import "errors"

func DecodeUrl(shortUrl string) (string, error) {
	url, exists := urls[shortUrl]
	if exists == false {
		return "", errors.New("invalid short url")
	}
	return url, nil
}
