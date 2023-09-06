package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	Conf   *Config
	HTTP   *http.Server
	Router *gin.Engine
}

type Config struct {
	Host string
	Port int
}

func NewServer(c *Config) *Server {
	return &Server{
		Conf:   c,
		Router: gin.New(),
	}
}

func (s *Server) Up() error {
	s.HTTP = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Conf.Host, s.Conf.Port),
		Handler: s.Router,
	}
	fmt.Printf("Start server at the port %d\n", s.Conf.Port)
	err := s.HTTP.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to Listen and Serve: %v\n", err)
		return err
	}
	return nil
}
