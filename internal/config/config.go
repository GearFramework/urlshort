package config

import (
	"os"
)

const (
	defaultAddress  = ":8080"
	defaultShortURL = "http://localhost:8080"
)

type ServiceConfig struct {
	Addr         string
	ShortURLHost string
}

func GetConfig() *ServiceConfig {
	conf := ParseFlags(defaultAddress, defaultShortURL)
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		conf.Addr = envAddr
	}
	if envURLHost := os.Getenv("BASE_URL"); envURLHost != "" {
		conf.ShortURLHost = envURLHost
	}
	return conf
}
