package handlers

import (
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserURLs(ctx *gin.Context, api pkg.APIShortener) {
	userID, ok := ctx.Get(pkg.UserIDParamName)
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	userURLs := api.GetUserURLs(ctx, userID.(int))
	if len(userURLs) == 0 {
		ctx.JSON(http.StatusNoContent, userURLs)
		return
	}
	ctx.JSON(http.StatusOK, userURLs)
}
