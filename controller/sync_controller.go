package controller

import (
	"imap-sync/internal"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleSync(ctx *gin.Context) {

	if ctx.Request.FormValue("retry") != "" {
		handleRetry(ctx)
		return
	} else if ctx.Request.FormValue("cancel") != "" {
		handleCancel(ctx)
		return
	}

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

func handleCancel(ctx *gin.Context) {
	id := ctx.Request.FormValue("cancel")

	val, err := strconv.Atoi(id)
	if err != nil {
		log.Errorf("Error converting %s to int", id)
	}

	task := internal.GetTaskFromID(val)

	internal.CancelTask(task)

	log.Infof("%#v", task)
}

func handleRetry(ctx *gin.Context) {
	id := ctx.Request.FormValue("retry")

	val, err := strconv.Atoi(id)
	if err != nil {
		log.Errorf("Error converting %s to int", id)
	}

	task := internal.GetTaskFromID(val)
	internal.RetryTask(task)

	sourceServer := task.SourceServer
	sourceAccount := task.SourceAccount
	sourcePassword := task.SourcePassword
	destinationServer := task.DestinationServer
	destinationAccount := task.DestinationAccount
	destinationPassword := task.DestinationPassword

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

	ctx.HTML(200, "sync_success.html", creds)
}
