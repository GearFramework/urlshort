package app

import (
	"fmt"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"math/rand"
	"time"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lenAlpha    = len(alphabet)
	defShortLen = 8
)

func (app *ShortApp) EncodeURL(url string) string {
	app.Store.Lock()
	defer app.Store.Unlock()
	code, exists := app.Store.GetCode(url)
	if !exists {
		code = app.getRandomString(defShortLen)
		if err := app.Store.Insert(url, code); err != nil {
			logger.Log.Error(err.Error())
		}
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
