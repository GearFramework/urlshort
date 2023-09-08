package app

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestEncodeURL(t *testing.T) {
	codes = make(codesMap, 10)
	urls = make(urlsMap, 10)

	testUrls := []string{
		"http://ya.ru",
		"http://yandex.ru",
	}
	for _, url := range testUrls {
		code := EncodeURL(url)
		assert.NotEmpty(t, code)
		assert.Equal(t, defShortLen, len(code))
		assert.Regexp(t, regexp.MustCompile(`^[a-zA-Z0-9]+$`), code)
	}
}

func TestEncodeURLExists(t *testing.T) {
	codes = codesMap{
		"http://ya.ru":     "dHGfdhj4",
		"http://yandex.ru": "78gsshSd",
	}
	testUrls := []struct {
		url  string
		want string
	}{
		{"http://ya.ru", "dHGfdhj4"},
		{"http://yandex.ru", "78gsshSd"},
	}
	for _, test := range testUrls {
		code := EncodeURL(test.url)
		assert.Equal(t, test.want, code)
	}
}
