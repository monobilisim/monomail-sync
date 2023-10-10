package controller

import (
	"flag"
	"imap-sync/internal"

	"github.com/gin-gonic/gin"
)

var (
	source_server        = flag.String("source_server", "", "Source server")
	source_account       = flag.String("source_account", "", "Source account")
	source_password      = flag.String("source_password", "", "Source password")
	destination_server   = flag.String("destination_server", "", "Destination server")
	destination_account  = flag.String("destination_account", "", "Destination account")
	destination_password = flag.String("destination_password", "", "Destination password")
)

func HandleRoot(ctx *gin.Context) {
	sourceDetails := internal.Credentials{
		Server:   *source_server,
		Account:  *source_account,
		Password: *source_password,
	}

	destinationDetails := internal.Credentials{
		Server:   *destination_server,
		Account:  *destination_account,
		Password: *destination_password,
	}

	data := struct {
		SourceDetails      internal.Credentials
		DestinationDetails internal.Credentials
	}{
		SourceDetails:      sourceDetails,
		DestinationDetails: destinationDetails,
	}
	ctx.HTML(200, "index.html", data)
}
