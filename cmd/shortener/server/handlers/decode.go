package handlers

import (
	"github.com/GearFramework/urlshort/internal/app"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func DecodeURL(ctx *gin.Context) {
	code := ctx.Param("code")
	url, err := app.DecodeURL(code)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}
	log.Printf("Request short code: %s url: %s", code, url)
	ctx.Header("Location", url)
	ctx.Status(http.StatusTemporaryRedirect)
	ctx.Done()
}
