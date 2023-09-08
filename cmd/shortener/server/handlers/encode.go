package handlers

import (
	"fmt"
	"github.com/GearFramework/urlshort/internal/app"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	neturl "net/url"
)

const (
	urlPattern = "http://localhost:8080/%s"
)

func EncodeURL(ctx *gin.Context) {
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
	code := app.EncodeURL(url)
	log.Printf("Request url: %s short code: %s\n", url, code)
	ctx.Data(http.StatusCreated, "text/plain", []byte(fmt.Sprintf(urlPattern, code)))
}
