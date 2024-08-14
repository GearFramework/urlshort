package server

import (
	"net/http"

	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// Server http-server
type Server struct {
	Conf   *config.ServiceConfig
	HTTP   *http.Server
	Router *gin.Engine
	api    pkg.APIShortener
}

// NewServer return new http server
func NewServer(c *config.ServiceConfig, api pkg.APIShortener) (*Server, error) {
	if err := logger.Initialize(c.LoggerLevel); err != nil {
		return nil, err
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	s := Server{
		Conf:   c,
		Router: router,
		api:    api,
	}
	pprof.Register(router)
	router.Use(s.logger())
	router.Use(s.compress())
	router.Use(s.auth())
	return &s, nil
}

// Up run server
func (s *Server) Up() error {
	s.HTTP = &http.Server{
		Addr:    s.Conf.Addr,
		Handler: s.Router,
	}
	logger.Log.Infof("Start server at the %s\n", s.Conf.Addr)
	err := s.Router.Run(s.Conf.Addr)
	if err != nil {
		logger.Log.Infof("Failed to Listen and Serve: %v\n", err)
		return err
	}
	return nil
}
