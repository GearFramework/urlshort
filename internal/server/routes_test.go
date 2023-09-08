package server

import (
	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoutes(t *testing.T) {
	a := app.NewShortener(config.ParseFlags())
	s := NewServer(&Config{Addr: a.Conf.Host}, a)
	s.InitRoutes()
	tests := map[string][]struct {
		pathExpected string
		valid        bool
	}{
		"POST": {
			{"/:id", false},
			{"/", true},
			{"/short/:code", false},
		},
		"GET": {
			{"/:id", false},
			{"/", false},
			{"/:code", true},
		},
	}
	routes := s.Router.Routes()
	for method, paths := range tests {
		for _, test := range paths {
			exists := isRouteExists(method, test.pathExpected, routes)
			assert.Equal(t, test.valid, exists)
		}
	}
}

func isRouteExists(method, path string, routes gin.RoutesInfo) bool {
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return true
		}
	}
	return false
}
