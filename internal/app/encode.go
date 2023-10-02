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
	app.store.Lock()
	defer app.store.Unlock()
	code, exists := app.store.getCode(url)
	if !exists {
		code = app.getRandomString(defShortLen)
		app.store.add(url, code)
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
