package config

import (
	"os"
)

const (
	defaultAddress  = ":8080"
	defaultShortURL = "http://localhost:8080"
	defaultLevel    = "info"
	//defaultStoragePath = "/tmp/short-url-db.json"
	defaultStoragePath = ""
	//defaultDatabaseDSN = "postgres://pgadmin:159753@localhost:5432/urlshortly"
	defaultDatabaseDSN = ""
	defaultEnableHTTPS = false
)

// ServiceConfig struct of application config
type ServiceConfig struct {
	Addr            string
	ShortURLHost    string
	LoggerLevel     string
	StorageFilePath string
	DatabaseDSN     string
	EnableHTTPS     bool
}

// GetConfig create and return application config
func GetConfig() *ServiceConfig {
	conf := ParseFlags()
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
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		conf.DatabaseDSN = envDatabaseDSN
	}
	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		conf.EnableHTTPS = true
	}
	return conf
}
