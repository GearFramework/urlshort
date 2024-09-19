package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/server"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func stringBuild(b string) string {
	if b == "" {
		return "N/A"
	}
	return b
}

func printGreeting() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		stringBuild(buildVersion),
		stringBuild(buildDate),
		stringBuild(buildCommit),
	)
}

func main() {
	printGreeting()
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	shortener, err := app.NewShortener(config.GetConfig())
	if err != nil {
		return err
	}
	gracefulStop(shortener.StopApp)
	s, err := server.NewServer(shortener.Conf, shortener)
	if err != nil {
		return err
	}
	s.InitRoutes()
	if shortener.Conf.EnableHTTPS {
		return s.UpTLS()
	}
	return s.Up()
}

func gracefulStop(stopCallback func()) {
	gracefulStopChan := make(chan os.Signal, 1)
	signal.Notify(
		gracefulStopChan,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	go func() {
		sig := <-gracefulStopChan
		stopCallback()
		log.Printf("Caught sig: %+v\n", sig)
		log.Println("Application graceful stop!")
		os.Exit(0)
	}()
}
