package server

import (
	"fmt"
	"github.com/GearFramework/urlshort/internal/pkg/compresser"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func (s *Server) logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		duration := logger.GetDurationInMilliseconds(start)
		logger.Log.Infoln(
			"uri", ctx.Request.RequestURI,
			"method", ctx.Request.Method,
			"status", ctx.Writer.Status(), // получаем перехваченный код статуса ответа
			"duration", fmt.Sprintf("%.4f ms", duration),
			"size", ctx.Writer.Size(), // получаем перехваченный размер ответа
		)
	}
}

func (s *Server) compress() gin.HandlerFunc {
	return compresser.NewCompressor()
}
