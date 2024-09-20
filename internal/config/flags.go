package config

import (
	"flag"
	"fmt"
	"os"
)

type ShortlyFlags struct {
	Addr,
	ShortURLHost,
	LogLevel,
	StorageFilePath,
	DatabaseDSN,
	ConfigFile string
	EnableHTTPS bool
}

// ParseFlags parse command line flags to application config
func ParseFlags() *ShortlyFlags {
	var fl ShortlyFlags
	fmt.Printf("Service started with flags: %v\n", os.Args[1:])
	flag.StringVar(&fl.Addr, "a", "", "address to run server")
	flag.StringVar(&fl.ShortURLHost, "b", "", "base address to result short URL")
	flag.StringVar(&fl.LogLevel, "l", "", "logger level")
	flag.StringVar(&fl.StorageFilePath, "f", "", "short url storage path")
	flag.StringVar(&fl.DatabaseDSN, "d", "", "database connection DSN")
	flag.BoolVar(&fl.EnableHTTPS, "s", false, "enable HTTPS support")
	flag.StringVar(&fl.ConfigFile, "c", "", "use config file")
	flag.StringVar(&fl.ConfigFile, "config", fl.ConfigFile, "use config file")
	flag.Parse()
	fmt.Println("Config from flags: ", fl)
	return &fl
}

func GetConfigFile() string {
	var confFile string
	flag.StringVar(&confFile, "c", defaultConfigFile, "use config file")
	flag.Parse()
	fmt.Println("Config file: ", confFile)
	if confFile != "" {
		return confFile
	}
	flag.StringVar(&confFile, "config", defaultConfigFile, "use config file")
	return confFile
}
