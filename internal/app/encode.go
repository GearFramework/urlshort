package app

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lenAlpha    = len(alphabet)
	defShortLen = 8
)

func (app *ShortApp) EncodeURL(url string) string {
	code, exists := app.codes[url]
	if !exists {
		code = app.getRandomString(defShortLen)
		app.AddShortly(url, code)
	}
	return fmt.Sprintf("%s/%s", app.Conf.ShortURLHost, code)
}

func (app *ShortApp) getRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = alphabet[rnd.Intn(lenAlpha)]
	}
	return string(b)
}
