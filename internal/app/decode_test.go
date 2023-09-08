package app

import (
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeURL(t *testing.T) {
	if shortener == nil {
		shortener = NewShortener(config.GetConfig())
	}
	shortener.AddShortly("http://ya.ru", "dHGfdhj4")
	shortener.AddShortly("http://yandex.ru", "78gsshSd")

	testCodes := []struct {
		code  string
		want  string
		error bool
	}{
		{"dHGfdhj4", "http://ya.ru", false},
		{"78gsshSd", "http://yandex.ru", false},
		{"dHGfdhj4", "http://ya.ru", false},
		{"7nnDfdds", "", true},
	}
	for _, test := range testCodes {
		url, err := shortener.DecodeURL(test.code)
		if test.error {
			t.Run("has error", func(t *testing.T) {
				assert.Error(t, err)
			})
		} else {
			t.Run("has error", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, test.want, url)
			})
		}
	}
}
