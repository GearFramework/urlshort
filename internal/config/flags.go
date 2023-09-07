package config

import "flag"

func ParseFlags() *ServiceConfig {
	conf := ServiceConfig{}
	flag.StringVar(&conf.Host, "a", ":8080", "address to run server")
	flag.StringVar(&conf.ShortURLHost, "b", "http://localhost:8080", "base address to result short URL")
	flag.Parse()
	return &conf
}
