package server

import (
	"github.com/GearFramework/urlshort/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func (s *Server) InitRoutes() {
	s.Router.POST("/", func(ctx *gin.Context) { handlers.EncodeURL(s.api, ctx) })
	s.Router.GET("/:code", func(ctx *gin.Context) { handlers.DecodeURL(s.api, ctx) })
	s.Router.POST("/api/shorten", func(ctx *gin.Context) { handlers.EncodeURLFromJSON(s.api, ctx) })
	s.Router.POST("/ping", func(ctx *gin.Context) { handlers.Ping(s.api, ctx) })
	s.Router.NoRoute(handlers.InvalidMethod)
}
