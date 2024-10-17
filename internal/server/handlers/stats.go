package handlers

import (
	"github.com/GearFramework/urlshort/internal/pkg/auth"
	"net"
	"net/http"

	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ResponseStats struct of statistic
type ResponseStats struct {
	pkg.Stats
}

// GetInternalStats return internal statistics about short urls and uers
func GetInternalStats(ctx *gin.Context, api pkg.APIShortener, conf *config.ServiceConfig) {
	if err := validateUserIP(ctx, conf.TrustedSubnet); err != nil {
		logger.Log.Errorf("unauthorized access: %s\n", err)
		ctx.Status(http.StatusForbidden)
		return
	}
	stats, err := api.GetStats(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, ResponseStats{
		*stats,
	})
}

func validateUserIP(ctx *gin.Context, trustedSubnet string) error {
	_, trustNet, err := auth.GetTrustedIP(trustedSubnet)
	if err != nil {
		return err
	}
	userIP, err := getXRealIP(ctx)
	if err != nil {
		return err
	}
	if !trustNet.Contains(userIP) {
		return auth.ErrIPNotFromTrustedNetwork
	}
	return nil
}

func getXRealIP(ctx *gin.Context) (net.IP, error) {
	IP := ctx.Request.Header.Get("X-Real-IP")
	if IP == "" {
		return nil, auth.ErrEmptyXRealIP
	}
	return auth.ParseIP(IP), nil
}
