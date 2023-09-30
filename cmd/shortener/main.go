package main

import (
	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	gracefulStop()
	shortener := app.NewShortener(config.GetConfig())
	s, err := server.NewServer(shortener.Conf, shortener)
	if err != nil {
		return err
	}
	s.InitRoutes()
	return s.Up()
}

func gracefulStop() {
	gracefulStopChan := make(chan os.Signal, 1)
	signal.Notify(
		gracefulStopChan,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	go func() {
		_ = <-gracefulStopChan
		os.Exit(0)
	}()
}
