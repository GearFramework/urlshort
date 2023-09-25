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
	s, err := server.NewServer(shortener.Conf, shortener)
	if err != nil {
		return err
	}
	s.InitRoutes()
	return s.Up()
}
