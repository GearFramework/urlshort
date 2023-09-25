package config

import (
	"os"
)

const (
	defaultAddress  = ":8080"
	defaultShortURL = "http://localhost:8080"
	defaultLevel    = "info"
)

type ServiceConfig struct {
	Addr         string
	ShortURLHost string
	LoggerLevel  string
}

func GetConfig() *ServiceConfig {
	conf := ParseFlags(defaultAddress, defaultShortURL, defaultLevel)
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		conf.Addr = envAddr
	}
	if envURLHost := os.Getenv("BASE_URL"); envURLHost != "" {
		conf.ShortURLHost = envURLHost
	}
	if envLoggerLevel := os.Getenv("LOGGER_LEVEL"); envLoggerLevel != "" {
		conf.LoggerLevel = envLoggerLevel
	}
	return conf
}
