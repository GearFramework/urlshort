package handlers

import (
	"errors"
	"net"
	"net/http"

	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ResponseStats struct {
	pkg.Stats
}

var (
	errTrustedNetworkNotDefined = errors.New("trusted network not defined")
	errEmptyXRealIP             = errors.New("empty X-Real-IP header")
)

// GetInternalStats return internal statistics about short urls and uers
func GetInternalStats(ctx *gin.Context, api pkg.APIShortener, conf *config.ServiceConfig) {
	if err := validateUserIP(ctx, conf); err != nil {
		logger.Log.Errorf("unauthorized access: %s\n", err)
		ctx.Status(http.StatusForbidden)
	}
	stats := ResponseStats{
		*api.GetStats(ctx),
	}
	ctx.JSON(http.StatusOK, stats)
}

func validateUserIP(ctx *gin.Context, conf *config.ServiceConfig) error {
	_, trustNet, err := getTrustedIP(conf)
	if err != nil {
		return err
	}
	userIP, err := getXRealIP(ctx)
	if err != nil {
		return err
	}
	if !trustNet.Contains(userIP) {
		return err
	}
	return nil
}

func getTrustedIP(conf *config.ServiceConfig) (net.IP, *net.IPNet, error) {
	if conf.TrustedSubnet == "" {
		return nil, nil, errTrustedNetworkNotDefined
	}
	return parseCIDR(conf.TrustedSubnet)
}

func getXRealIP(ctx *gin.Context) (net.IP, error) {
	IP := ctx.Request.Header.Get("X-Real-IP")
	if IP == "" {
		return nil, errEmptyXRealIP
	}
	return parseIP(IP), nil
}

func parseIP(IP string) net.IP {
	return net.ParseIP(IP)
}

func parseCIDR(IP string) (net.IP, *net.IPNet, error) {
	return net.ParseCIDR(IP)
}
