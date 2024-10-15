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
	errChan := make(chan error)
	gracefulStop(shortener.StopApp)
	go func() {
		errChan <- runHTTPServer(shortener)
	}()
	go func() {
		errChan <- runRPCServer(shortener)
	}()
	if err = <-errChan; err != nil {
		return err
	}
	return nil
}

func runHTTPServer(a *app.ShortApp) error {
	s, err := server.NewServer(a.Conf, a)
	if err != nil {
		return err
	}
	s.InitRoutes()
	if a.Conf.EnableHTTPS {
		return s.UpTLS()
	}
	return s.Up()
}

func runRPCServer(a *app.ShortApp) error {
	r, err := server.NewRPCServer(a.Conf, a)
	if err != nil {
		return err
	}
	return r.Up()
}

func gracefulStop(stopCallback func()) {
	gracefulStopChan := make(chan os.Signal, 1)
	signal.Notify(
		gracefulStopChan,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	go func() {
		sig := <-gracefulStopChan
		stopCallback()
		log.Printf("Caught sig: %+v\n", sig)
		log.Println("Application graceful stop!")
		os.Exit(0)
	}()
}
