package main

import (
	"flag"
)

type Credentials struct {
	Server   string
	Account  string
	Password string
}

var (
	port                 = flag.String("port", "8080", "Port to listen on")
	source_server        = flag.String("source_server", "https://example.com", "Source server")
	source_account       = flag.String("source_account", "default source account", "Source account")
	source_password      = flag.String("source_password", "default source password", "Source password")
	destination_server   = flag.String("destination_server", "https://dest.example.com", "Destination server")
	destination_account  = flag.String("destination_account", "default destination account", "Destination account")
	destination_password = flag.String("destination_password", "default destination password", "Destination password")
)

func main() {
	flag.Parse()

	r := initServer()

	r.Run(":" + *port)
}
