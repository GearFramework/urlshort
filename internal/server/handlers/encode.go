package handlers

import (
	"encoding/json"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
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

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

func EncodeURLFromJSON(api pkg.APIShortener, ctx *gin.Context) {
	ct := strings.Split(ctx.Request.Header.Get("Content-Type"), ";")
	if ct[0] != "application/json" {
		logger.Log.Errorf(
			"invalid request header: Content-Type %s\n",
			ctx.Request.Header.Get("Content-Type"),
		)
		ctx.Status(http.StatusBadRequest)
		return
	}
	defer ctx.Request.Body.Close()
	dec := json.NewDecoder(ctx.Request.Body)
	var req RequestJSON
	if err := dec.Decode(&req); err != nil {
		logger.Log.Errorln("invalid url")
		ctx.Status(http.StatusBadRequest)
		return
	}
	resp := ResponseJSON{api.EncodeURL(req.URL)}
	logger.Log.Infof("Request url: %s short url: %s\n", req.URL, resp.Result)
	ctx.JSON(http.StatusCreated, resp)
}
