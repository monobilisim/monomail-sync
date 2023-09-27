package internal

import "flag"

var (
	port                 = flag.String("port", "8080", "Port to listen on")
	source_server        = flag.String("source_server", "", "Source server")
	source_account       = flag.String("source_account", "", "Source account")
	source_password      = flag.String("source_password", "", "Source password")
	destination_server   = flag.String("destination_server", "", "Destination server")
	destination_account  = flag.String("destination_account", "", "Destination account")
	destination_password = flag.String("destination_password", "", "Destination password")
)

type Credentials struct {
	Server   string
	Account  string
	Password string
}
