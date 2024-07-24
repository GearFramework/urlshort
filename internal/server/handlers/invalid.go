package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// InvalidMethod handler if request is bad
func InvalidMethod(ctx *gin.Context) {
	ctx.Status(http.StatusBadRequest)
}
