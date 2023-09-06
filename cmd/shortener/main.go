package main

import (
	"github.com/GearFramework/urlshort/cmd/shortener/server"
	"github.com/GearFramework/urlshort/internal/app"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	app.InitShortener()
	s := server.NewServer(&server.Config{Host: "localhost", Port: 8080})
	s.InitRoutes()
	return s.Up()
}
