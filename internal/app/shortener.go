package app

type codesMap map[string]string
type urlsMap map[string]string

var codes codesMap
var urls urlsMap

func InitShortener() {
	codes = make(codesMap, 10)
	urls = make(urlsMap, 10)
}
