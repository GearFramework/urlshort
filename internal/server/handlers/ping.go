package handlers

import (
	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(api pkg.APIShortener, ctx *gin.Context) {
	if err := api.(*app.ShortApp).DB.Ping(); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}
