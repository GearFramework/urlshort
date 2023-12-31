package main

import (
	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/server"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	shortener := app.NewShortener(config.GetConfig())
	s := server.NewServer(&server.Config{Addr: shortener.Conf.Addr}, shortener)
	s.InitRoutes()
	return s.Up()
}
