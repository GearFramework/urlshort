package config

import (
	"flag"
)

func ParseFlags(defaultAddr, defaultShortURL string, defaultLevel string) *ServiceConfig {
	var conf ServiceConfig
	flag.StringVar(&conf.Addr, "a", defaultAddr, "address to run server")
	flag.StringVar(&conf.ShortURLHost, "b", defaultShortURL, "base address to result short URL")
	flag.StringVar(&conf.LoggerLevel, "l", defaultLevel, "logger level")
	flag.Parse()
	return &conf
}
