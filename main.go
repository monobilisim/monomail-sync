package main

import (
	"flag"
	"imap-sync/internal"
)

func main() {
	flag.Parse()
	internal.InitServer()
}
