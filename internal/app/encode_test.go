package app

import (
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/stretchr/testify/assert"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestEncodeURL(t *testing.T) {
	if shortener == nil {
		shortener = NewShortener(config.ParseFlags())
	}
	testURLs := []string{
		"http://ya.ru",
		"http://yandex.ru",
	}
	for _, testURI := range testURLs {
		shortURI := shortener.EncodeURL(testURI)
		assert.NotEmpty(t, shortURI)
		parsedURI, _ := url.ParseRequestURI(shortURI)
		assert.Equal(t, defShortLen, len(strings.TrimLeft(parsedURI.Path, "/")))
		assert.Regexp(t, regexp.MustCompile(`^/[a-zA-Z0-9]+$`), parsedURI.Path)
	}
}

func TestEncodeURLExists(t *testing.T) {
	if shortener == nil {
		shortener = NewShortener(config.ParseFlags())
	}
	shortener.AddShortly("http://ya.ru", "dHGfdhj4")
	shortener.AddShortly("http://yandex.ru", "78gsshSd")
	testURLs := []struct {
		url  string
		want string
	}{
		{"http://ya.ru", shortener.Conf.ShortURLHost + "/dHGfdhj4"},
		{"http://yandex.ru", shortener.Conf.ShortURLHost + "/78gsshSd"},
	}
	for _, test := range testURLs {
		shortURL := shortener.EncodeURL(test.url)
		assert.Equal(t, test.want, shortURL)
	}
}
