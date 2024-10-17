package app

import (
	"context"
	"testing"
	"time"

	"github.com/GearFramework/urlshort/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestDecodeURL(t *testing.T) {
	var err error
	if shortener == nil {
		shortener, err = NewShortener(config.GetConfig())
		assert.NoError(t, err)
	}
	shortener.ClearShortly()
	c, err := shortener.Store.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, c)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shortener.AddShortly(ctx, 1, "http://ya.ru", "dHGfdhj4")
	shortener.AddShortly(ctx, 1, "http://yandex.ru", "78gsshSd")
	c, err = shortener.Store.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 2, c)
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
		url, err := shortener.DecodeURL(ctx, test.code)
		if test.error {
			t.Run("has error", func(t *testing.T) {
				assert.Error(t, err)
			})
		} else {
			t.Run("has no error", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, test.want, url)
			})
		}
	}
}
