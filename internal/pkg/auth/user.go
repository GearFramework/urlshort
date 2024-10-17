package auth

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
)

// user authorizations settings
const (
	// TokenExpired token time life
	TokenExpired = time.Hour * 24
	// SecretKey salt
	SecretKey = "bu7HBJD&873HVHJdh*Jbhsfdfs8622Dsf"
)

var (
	// ErrTrustedNetworkNotDefined trusted network not defined
	ErrTrustedNetworkNotDefined = errors.New("trusted network not defined")
	// ErrEmptyXRealIP X-Real-IP is undefined
	ErrEmptyXRealIP = errors.New("empty X-Real-IP header")
	// ErrIPNotFromTrustedNetwork user IP not from trusted network
	ErrIPNotFromTrustedNetwork = errors.New("IP not from trusted network")
	// ErrNeedAuthorization error if need authorization
	ErrNeedAuthorization = errors.New("требуется авторизация")
	// ErrInvalidAuthorization error if invalid authorization
	ErrInvalidAuthorization = errors.New("отсутствует ID пользователя")
)

// Claims jwt struct
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// BuildJWT create jwt token
func BuildJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpired)),
		},
		UserID: userID,
	})
	tk, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tk, nil
}

// GetUserIDFromJWT return user ID from token
func GetUserIDFromJWT(tk string) int {
	claims, err := getClaims(tk)
	if err != nil {
		return -1
	}
	logger.Log.Infof("app user ID: %d", claims.UserID)
	return claims.UserID
}

func getClaims(tk string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tk, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil || !token.Valid {
		logger.Log.Error(err.Error())
		return nil, err
	}
	return claims, nil
}

// GetTrustedIP return parsed CIDR
func GetTrustedIP(trustedSubnet string) (net.IP, *net.IPNet, error) {
	if trustedSubnet == "" {
		return nil, nil, ErrTrustedNetworkNotDefined
	}
	return ParseCIDR(trustedSubnet)
}

// ParseIP return result of parsed IP
func ParseIP(IP string) net.IP {
	return net.ParseIP(IP)
}

// ParseCIDR return result of parsed CIDR
func ParseCIDR(IP string) (net.IP, *net.IPNet, error) {
	return net.ParseCIDR(IP)
}
