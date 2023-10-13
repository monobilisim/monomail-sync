package controller

import (
	"imap-sync/internal"

	"github.com/gin-gonic/gin"
)

func handleSync(ctx *gin.Context) {
	sourceServer := ctx.PostForm("source_server")
	sourceAccount := ctx.PostForm("source_account")
	sourcePassword := ctx.PostForm("source_password")
	destinationServer := ctx.PostForm("destination_server")
	destinationAccount := ctx.PostForm("destination_account")
	destinationPassword := ctx.PostForm("destination_password")

	sourceDetails := internal.Credentials{
		Server:   sourceServer,
		Account:  sourceAccount,
		Password: sourcePassword,
	}

	destinationDetails := internal.Credentials{
		Server:   destinationServer,
		Account:  destinationAccount,
		Password: destinationPassword,
	}

	creds := struct {
		Source      internal.Credentials
		Destination internal.Credentials
	}{
		Source:      sourceDetails,
		Destination: destinationDetails,
	}

	// Add to queue
	log.Infof("Adding %s to queue", sourceDetails.Account)
	internal.AddTask(sourceDetails, destinationDetails)
	ctx.HTML(200, "sync_success.html", creds)
}
