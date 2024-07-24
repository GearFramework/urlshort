package auth

import (
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser(t *testing.T) {
	err := logger.Initialize("info")
	assert.NoError(t, err)
	userID := 1
	tk, err := BuildJWT(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tk)
	userIDFromJWT := GetUserIDFromJWT(tk)
	assert.Equal(t, userID, userIDFromJWT)
}
