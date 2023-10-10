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

	// Add to queue
	log.Infof("Adding %s to queue", sourceDetails.Account)
	internal.AddTask(sourceDetails, destinationDetails)

	// log.Infof("Syncing %s to %s", sourceDetails.Account, destinationDetails.Account)

	// err := syncIMAP(sourceDetails, destinationDetails)
	// if err != nil {
	// 	ctx.HTML(200, "error.html", err.Error())
	// 	return
	// }
	// ctx.HTML(200, "success.html", "Synced "+sourceDetails.Account+" to "+destinationDetails.Account)
}
