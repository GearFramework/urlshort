package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	neturl "net/url"
	"strings"

	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

// EncodeURL return short url for requested url
func EncodeURL(ctx *gin.Context, api pkg.APIShortener) {
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
	userID, _ := ctx.Get("userID")
	shortURL, conflict := api.EncodeURL(ctx, userID.(int), url)
	status := http.StatusCreated
	logger.Log.Infof("Request url: %s short url: %s\n", url, shortURL)
	if conflict {
		status = http.StatusConflict
	}
	ctx.Data(status, "text/plain", []byte(shortURL))
}

// RequestJSON request struct for handler EncodeURLFromJSON
type RequestJSON struct {
	URL string `json:"url"`
}

// ResponseJSON response on request RequestJSON
type ResponseJSON struct {
	Result string `json:"result"`
}

// RequestBatchJSON request struct for handler BatchEncodeURLs
type RequestBatchJSON pkg.BatchURLs

// ResponseBatchJSON response on request RequestBatchJSON
type ResponseBatchJSON pkg.ResultBatchShort

// EncodeURLFromJSON return short url for url in json content type request
func EncodeURLFromJSON(ctx *gin.Context, api pkg.APIShortener) {
	if !strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/json") {
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
	userID, _ := ctx.Get(pkg.UserIDParamName)
	res, conflict := api.EncodeURL(ctx, userID.(int), req.URL)
	resp := ResponseJSON{res}
	status := http.StatusCreated
	logger.Log.Infof("Request url: %s short url: %s\n", req.URL, resp.Result)
	if conflict {
		status = http.StatusConflict
	}
	ctx.JSON(status, resp)
}

// BatchEncodeURLs return short urls for urls in batch json content type request
func BatchEncodeURLs(ctx *gin.Context, api pkg.APIShortener) {
	if !strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/json") {
		logger.Log.Errorf(
			"invalid request header: Content-Type %s\n",
			ctx.Request.Header.Get("Content-Type"),
		)
		ctx.Status(http.StatusBadRequest)
		return
	}
	defer ctx.Request.Body.Close()
	dec := json.NewDecoder(ctx.Request.Body)
	var req []pkg.BatchURLs
	if err := dec.Decode(&req); err != nil {
		logger.Log.Errorln("invalid url")
		ctx.Status(http.StatusBadRequest)
		return
	}
	userID, ok := ctx.Get(pkg.UserIDParamName)
	if !ok {
		ctx.Status(http.StatusNoContent)
		return
	}
	resp := api.BatchEncodeURL(ctx, userID.(int), req)
	//logger.Log.Infof("Request url: %s short url: %s\n", req.URL, resp.Result)
	ctx.JSON(http.StatusCreated, resp)
}
