package handlers

import (
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	neturl "net/url"
)

func EncodeURL(api pkg.ApiShortener, ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	url := string(body)
	if len(url) == 0 {
		log.Printf("Error: empty url in request")
		ctx.Status(http.StatusBadRequest)
		return
	}
	if _, err = neturl.ParseRequestURI(url); err != nil {
		log.Printf("Error: invalid url\n")
		ctx.Status(http.StatusBadRequest)
		return
	}
	shortURL := api.EncodeURL(url)
	log.Printf("Request url: %s short url: %s\n", url, shortURL)
	ctx.Data(http.StatusCreated, "text/plain", []byte(shortURL))
}
