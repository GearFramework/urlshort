package server

import (
	"github.com/GearFramework/urlshort/cmd/shortener/server/handlers"
)

func (s *Server) InitRoutes() {
	s.Router.POST("/", handlers.EncodeURL)
	s.Router.GET("/:code", handlers.DecodeURL)
	s.Router.NoRoute(handlers.InvalidMethod)
}
