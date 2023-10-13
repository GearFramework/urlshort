package server

import (
	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var a *app.ShortApp

func TestRoutes(t *testing.T) {
	var err error
	if a == nil {
		a, err = app.NewShortener(config.GetConfig())
	}
	assert.NoError(t, err)
	s, err := NewServer(a.Conf, a)
	assert.NoError(t, err)
	s.InitRoutes()
	tests := map[string][]struct {
		pathExpected string
		valid        bool
	}{
		"POST": {
			{"/:id", false},
			{"/", true},
			{"/short/:code", false},
			{"/api/shorten", true},
		},
		"GET": {
			{"/:id", false},
			{"/", false},
			{"/:code", true},
			{"/api/shorten", false},
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

func TestPing(t *testing.T) {
	var err error
	if a == nil {
		a, err = app.NewShortener(config.GetConfig())
	}
	assert.NoError(t, err)
	s, err := NewServer(a.Conf, a)
	assert.NoError(t, err)
	s.InitRoutes()
	request := httptest.NewRequest(
		http.MethodGet,
		"/ping",
		strings.NewReader(""),
	)
	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, request)
	response := w.Result()
	defer response.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
