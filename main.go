package main

import (
	"flag"
	"imap-sync/api"
)

func main() {
	flag.Parse()
	api.InitServer()
}
