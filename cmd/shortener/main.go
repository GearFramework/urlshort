package main

import (
	"errors"
	"fmt"
	"github.com/GearFramework/urlshort/internal/app"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	urlPattern = "http://localhost:8080/%s"
)

func init() {
	app.InitShortener()
}
func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	log.Println("Start server :8080")
	gracefulStop()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handleService)
	return http.ListenAndServe(`:8080`, mux)
}

func gracefulStop() {
	gracefulStopChan := make(chan os.Signal)
	signal.Notify(
		gracefulStopChan,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	go func() {
		sig := <-gracefulStopChan
		log.Printf("Caught signal: %+v\nStop server\n", sig)
		os.Exit(0)
	}()
}

func handleService(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleEncodeUrl(w, r)
		return
	}
	if r.Method == http.MethodGet {
		handleDecodeUrl(w, r)
		return
	}
	log.Println("Error: invalid request method")
	w.WriteHeader(http.StatusBadRequest)
}

func handleEncodeUrl(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responseError(w, err)
		return
	}
	url := string(body)
	if len(url) == 0 {
		responseError(w, errors.New("empty url in request"))
		return
	}
	if _, err = r.URL.Parse(url); err != nil {
		responseError(w, errors.New("invalid url in request"))
		return
	}
	code := app.EncodeUrl(url)
	log.Printf("Request url: %s short code: %s\n", url, code)
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write([]byte(fmt.Sprintf(urlPattern, code))); err != nil {
		log.Fatal(fmt.Sprintf("Error: %s\n", err.Error()))
	}
}

func handleDecodeUrl(w http.ResponseWriter, r *http.Request) {
	code, _ := strings.CutPrefix(r.URL.Path, "/")
	url, err := app.DecodeUrl(code)
	log.Printf("Request short code: %s url: %s", code, url)
	if err != nil {
		responseError(w, err)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func responseError(w http.ResponseWriter, err error) {
	log.Printf("Error: %s\n", err.Error())
	w.WriteHeader(http.StatusBadRequest)
}
