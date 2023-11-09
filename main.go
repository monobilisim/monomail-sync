package main

import (
	"flag"
	"imap-sync/api"
	"imap-sync/config"
)

func main() {
	flag.Parse()
	config.ParseConfig()
	api.InitServer()
}
