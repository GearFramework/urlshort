package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeURL(t *testing.T) {
	urls = urlsMap{
		"dHGfdhj4": "http://ya.ru",
		"78gsshSd": "http://yandex.ru",
	}
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
		url, err := DecodeURL(test.code)
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
