package server

import (
	"fmt"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	Conf   *Config
	HTTP   *http.Server
	Router *gin.Engine
	api    pkg.APIShortener
}

func NewServer(c *Config, api pkg.APIShortener) *Server {
	return &Server{
		Conf:   c,
		Router: gin.New(),
		api:    api,
	}
}

func (s *Server) Up() error {
	s.HTTP = &http.Server{
		Addr:    s.Conf.Addr,
		Handler: s.Router,
	}
	fmt.Printf("Start server at the %s\n", s.Conf.Addr)
	err := s.HTTP.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to Listen and Serve: %v\n", err)
		return err
	}
	return nil
}
