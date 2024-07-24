package server

import (
	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServer(t *testing.T) {
	err := logger.Initialize("info")
	assert.NoError(t, err)
	conf := &config.ServiceConfig{
		Addr: ":8080",
	}
	sh, err := app.NewShortener(conf)
	assert.NoError(t, err)
	_, err = NewServer(conf, sh)
	assert.NoError(t, err)
}
