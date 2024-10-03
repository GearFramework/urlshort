package config

import (
	"flag"
	"fmt"
	"os"
)

var empty = "\x00"

type StringValue string

func (v StringValue) String() string {
	return string(v)
}

func (v StringValue) Set(s string) error {
	v = StringValue(s)
	return nil
}

func (v StringValue) Isset(s string) bool {
	return string(v) != empty
}

// ShortlyFlags struct of command-line flags
type ShortlyFlags struct {
	Addr,
	ShortURLHost,
	LogLevel,
	StorageFilePath,
	DatabaseDSN,
	ConfigFile string
	EnableHTTPS bool
	Test        StringValue
}

// ParseFlags parse command line flags to application config
func ParseFlags() *ShortlyFlags {
	var fl ShortlyFlags
	var t StringValue
	flag.Var(&t, "p", "ssds")
	fmt.Printf("Service started with flags: %v\n", os.Args[1:])
	flag.StringVar(&fl.Addr, "a", empty, "address to run server")
	flag.StringVar(&fl.ShortURLHost, "b", empty, "base address to result short URL")
	flag.StringVar(&fl.LogLevel, "l", empty, "logger level")
	flag.StringVar(&fl.StorageFilePath, "f", empty, "short url storage path")
	flag.StringVar(&fl.DatabaseDSN, "d", empty, "database connection DSN")
	flag.BoolVar(&fl.EnableHTTPS, "s", false, "enable HTTPS support")
	flag.StringVar(&fl.ConfigFile, "c", empty, "use config file")
	flag.StringVar(&fl.ConfigFile, "config", empty, "use config file")
	flag.Parse()
	fmt.Println("dddd p", t)
	fmt.Println("Config from flags: ", fl)
	return &fl
}
