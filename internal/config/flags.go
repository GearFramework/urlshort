package config

import (
	"flag"
	"fmt"
	"os"
)

var empty = "\x00"

// ShortlyFlags struct of command-line flags
type ShortlyFlags struct {
	Addr,
	ShortURLHost,
	LogLevel,
	StorageFilePath,
	DatabaseDSN,
	TrustedSubnet,
	ConfigFile string
	EnableHTTPS bool
}

// ParseFlags parse command line flags to application config
func ParseFlags() *ShortlyFlags {
	var fl ShortlyFlags
	fmt.Printf("Service started with flags: %v\n", os.Args[1:])
	flag.StringVar(&fl.Addr, "a", empty, "address to run server")
	flag.StringVar(&fl.ShortURLHost, "b", empty, "base address to result short URL")
	flag.StringVar(&fl.LogLevel, "l", empty, "logger level")
	flag.StringVar(&fl.StorageFilePath, "f", empty, "short url storage path")
	flag.StringVar(&fl.DatabaseDSN, "d", empty, "database connection DSN")
	flag.BoolVar(&fl.EnableHTTPS, "s", false, "enable HTTPS support")
	flag.StringVar(&fl.ConfigFile, "c", empty, "use config file")
	flag.StringVar(&fl.TrustedSubnet, "t", empty, "trusted subnet")
	flag.StringVar(&fl.ConfigFile, "config", empty, "use config file")
	flag.Parse()
	fmt.Println("Config from flags: ", fl)
	return &fl
}
