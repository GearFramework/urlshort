package handlers

import (
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	neturl "net/url"
)

func EncodeURL(api pkg.APIShortener, ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	url := string(body)
	if len(url) == 0 {
		logger.Log.Errorln("empty url in request")
		ctx.Status(http.StatusBadRequest)
		return
	}
	if _, err = neturl.ParseRequestURI(url); err != nil {
		logger.Log.Errorln("invalid url")
		ctx.Status(http.StatusBadRequest)
		return
	}
	shortURL := api.EncodeURL(url)
	logger.Log.Infof("Request url: %s short url: %s\n", url, shortURL)
	ctx.Data(http.StatusCreated, "text/plain", []byte(shortURL))
}
