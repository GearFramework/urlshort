package config

import (
	"encoding/json"
	"fmt"
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
	defaultConfigFile  = ""
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
	if fl.Addr != empty {
		conf.Addr = fl.Addr
	}
	if fl.ShortURLHost != empty {
		conf.ShortURLHost = string(fl.ShortURLHost)
	}
	if fl.LogLevel != empty {
		conf.LoggerLevel = fl.LogLevel
	}
	if fl.StorageFilePath != empty {
		conf.StorageFilePath = fl.StorageFilePath
	}
	if fl.DatabaseDSN != empty {
		conf.DatabaseDSN = fl.DatabaseDSN
	}
	if fl.EnableHTTPS {
		conf.EnableHTTPS = fl.EnableHTTPS
	}
}

func mappingEnvToConfig(conf *ServiceConfig) {
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		conf.Addr = os.Getenv("SERVER_ADDRESS")
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
