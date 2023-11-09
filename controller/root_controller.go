package controller

import (
	"imap-sync/config"
	"imap-sync/internal"
	"imap-sync/logger"

	"github.com/gin-gonic/gin"
)

var log = logger.Log

var (
	source_server       string
	source_account      string
	destination_server  string
	destination_account string
)

func HandleRoot(ctx *gin.Context) {
	source_server = config.Conf.SourceAndDestination.SourceServer
	source_account = config.Conf.SourceAndDestination.SourceMail
	destination_server = config.Conf.SourceAndDestination.DestinationServer
	destination_account = config.Conf.SourceAndDestination.DestinationMail

	sourceDetails := internal.Credentials{
		Server:  source_server,
		Account: source_account,
	}

	destinationDetails := internal.Credentials{
		Server:  destination_server,
		Account: destination_account,
	}

	data := struct {
		SourceDetails      internal.Credentials
		DestinationDetails internal.Credentials
		Text               map[string]string
		Table              map[string]string
	}{
		SourceDetails:      sourceDetails,
		DestinationDetails: destinationDetails,
		Text:               internal.Data["index"],
		Table:              internal.Data["table"],
	}
	ctx.HTML(200, "index.html", data)
}
