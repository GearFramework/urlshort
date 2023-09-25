package handlers

import (
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DecodeURL(api pkg.APIShortener, ctx *gin.Context) {
	code := ctx.Param("code")
	url, err := api.DecodeURL(code)
	if err != nil {
		logger.Log.Errorf("%s\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}
	logger.Log.Infof("Request short code: %s url: %s", code, url)
	ctx.Header("Location", url)
	ctx.Status(http.StatusTemporaryRedirect)
	ctx.Done()
}
