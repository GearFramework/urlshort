package handlers

import (
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func DecodeURL(api pkg.APIShortener, ctx *gin.Context) {
	code := ctx.Param("code")
	url, err := api.DecodeURL(code)
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
