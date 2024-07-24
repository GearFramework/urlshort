package config

import (
	"flag"
	"fmt"
	"os"
)

// ParseFlags parse command line flags to application config
func ParseFlags() *ServiceConfig {
	var conf ServiceConfig
	fmt.Printf("Service started with flags: %v\n", os.Args[1:])
	flag.StringVar(&conf.Addr, "a", defaultAddress, "address to run server")
	flag.StringVar(&conf.ShortURLHost, "b", defaultShortURL, "base address to result short URL")
	flag.StringVar(&conf.LoggerLevel, "l", defaultLevel, "logger level")
	flag.StringVar(&conf.StorageFilePath, "f", defaultStoragePath, "short url storage path")
	flag.StringVar(&conf.DatabaseDSN, "d", defaultDatabaseDSN, "database connection DSN")
	flag.Parse()
	return &conf
}
