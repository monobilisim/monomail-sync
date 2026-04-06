package main

import (
	"flag"
	"imap-sync/api"
	"imap-sync/config"
)

var configPath = flag.String("config", "/etc/monomail-sync.yml", "Path of the configuration file in YAML format")

func main() {
	flag.Parse()
	config.SetConfigPath(*configPath)
	config.ParseConfig()
	api.InitServer()
}
