package app

import (
	"fmt"
)

func (app *ShortApp) GetShortURL(code string) string {
	return fmt.Sprintf("%s/%s", app.Conf.ShortURLHost, code)
}
