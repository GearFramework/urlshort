package config

import (
	"os"
)

const (
	defaultAddress     = ":8080"
	defaultShortURL    = "http://localhost:8080"
	defaultLevel       = "info"
	defaultStoragePath = "/tmp/short-url-db.json"
)

type ServiceConfig struct {
	Addr            string
	ShortURLHost    string
	LoggerLevel     string
	StorageFilePath string
}

func GetConfig() *ServiceConfig {
	conf := ParseFlags(defaultAddress, defaultShortURL, defaultLevel, defaultStoragePath)
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		conf.Addr = envAddr
	}
	if envURLHost := os.Getenv("BASE_URL"); envURLHost != "" {
		conf.ShortURLHost = envURLHost
	}
	if envLoggerLevel := os.Getenv("LOGGER_LEVEL"); envLoggerLevel != "" {
		conf.LoggerLevel = envLoggerLevel
	}
	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		conf.StorageFilePath = envStoragePath
	}
	return conf
}
