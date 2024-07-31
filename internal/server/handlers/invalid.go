package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InvalidMethod handler if request is bad
func InvalidMethod(ctx *gin.Context) {
	ctx.Status(http.StatusBadRequest)
}
