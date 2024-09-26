package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	defaultAddress     = ":8080"
	defaultShortURL    = "http://localhost:8080"
	defaultLevel       = "info"
	defaultStoragePath = ""
	defaultDatabaseDSN = ""
	defaultEnableHTTPS = false
)

// ServiceConfig struct of application config
type ServiceConfig struct {
	Addr            string `json:"server_address"`
	ShortURLHost    string `json:"base_url"`
	LoggerLevel     string
	StorageFilePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	ConfigFile      string
}

// NewConfig constructor of ServiceConfig
func NewConfig() *ServiceConfig {
	return &ServiceConfig{
		Addr:            defaultAddress,
		ShortURLHost:    defaultShortURL,
		LoggerLevel:     defaultLevel,
		StorageFilePath: defaultStoragePath,
		DatabaseDSN:     defaultDatabaseDSN,
		EnableHTTPS:     defaultEnableHTTPS,
	}
}

// GetConfig create and return application config
func GetConfig() *ServiceConfig {
	conf := NewConfig()
	fl := ParseFlags()
	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		fl.ConfigFile = envConfigFile
	}
	if fl.ConfigFile != "" {
		fmt.Printf("Reading config file: %s\n", fl.ConfigFile)
		if err := loadConfigFile(fl.ConfigFile, conf); err != nil {
			fmt.Printf("Error loading config file: %s\n", err)
		}
		fmt.Println("Config from file: ", conf)
	}
	mappingFlagsToConfig(fl, conf)
	fmt.Println("Config after mapping: ", conf)
	mappingEnvToConfig(conf)
	checkEmptyConfigParams(conf)
	fmt.Println("Use config: ", conf)
	return conf
}

func loadConfigFile(filepath string, fl *ServiceConfig) error {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &fl); err != nil {
		return err
	}
	return nil
}

func mappingFlagsToConfig(fl *ShortlyFlags, conf *ServiceConfig) {
	if fl.Addr != "" {
		conf.Addr = fl.Addr
	}
	if fl.ShortURLHost != "" {
		conf.ShortURLHost = fl.ShortURLHost
	}
	if fl.LogLevel != "" {
		conf.LoggerLevel = fl.LogLevel
	}
	if fl.StorageFilePath != "" {
		conf.StorageFilePath = fl.StorageFilePath
	}
	if fl.DatabaseDSN != "" {
		conf.DatabaseDSN = fl.DatabaseDSN
	}
	if fl.EnableHTTPS {
		conf.EnableHTTPS = fl.EnableHTTPS
	}
}

func mappingEnvToConfig(conf *ServiceConfig) {
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
}

func checkEmptyConfigParams(conf *ServiceConfig) {
	if conf.Addr == "" {
		conf.Addr = defaultAddress
	}
	if conf.ShortURLHost == "" {
		conf.ShortURLHost = defaultShortURL
	}
	if conf.LoggerLevel == "" {
		conf.LoggerLevel = defaultLevel
	}
	if conf.StorageFilePath == "" {
		conf.StorageFilePath = defaultStoragePath
	}
	if conf.DatabaseDSN == "" {
		conf.DatabaseDSN = defaultDatabaseDSN
	}
	if !conf.EnableHTTPS {
		conf.EnableHTTPS = defaultEnableHTTPS
	}
}
