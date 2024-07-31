package handlers

import (
	"net/http"

	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/gin-gonic/gin"
)

// Ping check connection to storage
func Ping(ctx *gin.Context, api pkg.APIShortener) {
	if err := api.(*app.ShortApp).Store.Ping(); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}
