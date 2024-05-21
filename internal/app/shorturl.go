package app

import (
	"fmt"
)

// GetShortURL возвращает сокращённый урл по указанному коду
func (app *ShortApp) GetShortURL(code string) string {
	return fmt.Sprintf("%s/%s", app.Conf.ShortURLHost, code)
}
