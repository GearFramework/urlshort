package app

import (
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortener(t *testing.T) {
	err := logger.Initialize("info")
	assert.NoError(t, err)
	app, err := NewShortener(&config.ServiceConfig{
		Addr: ":8080",
	})
	assert.NoError(t, err)
	userID := app.GenerateUserID()
	assert.NotEqual(t, 0, userID)
	newUserID, tk, err := app.CreateToken()
	assert.NoError(t, err)
	assert.NotEqual(t, userID, newUserID)
	assert.NotEmpty(t, tk)
	authUserID, err := app.Auth(tk)
	assert.NoError(t, err)
	assert.Equal(t, newUserID, authUserID)
}
