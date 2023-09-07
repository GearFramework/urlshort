package app

import (
	"github.com/GearFramework/urlshort/internal/config"
)

type mapCodes map[string]string
type mapURLs map[string]string

type ShortApp struct {
	Conf  *config.ServiceConfig
	codes mapCodes
	urls  mapURLs
}

func NewShortener(conf *config.ServiceConfig) *ShortApp {
	shortener := ShortApp{Conf: conf}
	shortener.initApp()
	return &shortener
}

func (app *ShortApp) initApp() {
	app.codes = make(mapCodes, 10)
	app.urls = make(mapURLs, 10)
}
func (app *ShortApp) AddShortly(url, code string) {
	app.codes[url] = code
	app.urls[code] = url
}
