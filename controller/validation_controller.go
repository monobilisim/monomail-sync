package controller

import (
	"imap-sync/internal"

	"github.com/gin-gonic/gin"
)

func HandleValidate(ctx *gin.Context) {

	validate := ctx.PostForm("validate")
	submitsync := ctx.PostForm("submit_sync")

	var SServer, SAccount, SPassword string
	var DServer, DAccount, DPassword string

	if validate != "" {
		SServer = ctx.PostForm("source_server")
		SAccount = ctx.PostForm("source_account")
		SPassword = ctx.PostForm("source_password")
		DServer = ctx.PostForm("destination_server")
		DAccount = ctx.PostForm("destination_account")
		DPassword = ctx.PostForm("destination_password")
	}

	if validate == "" && submitsync != "" {
		HandleSync(ctx)
		return
	}

	creds := internal.Credentials{
		Server:   SServer,
		Account:  SAccount,
		Password: SPassword,
		Source:   true,
	}

	log.Infof("Validating credentials for: %s", creds.Account)

	err := internal.ValidateCredentials(creds)
	if err != nil {
		ctx.HTML(200, "error.html", "Couldn't verify for user: "+SAccount)
		return
	}

	creds = internal.Credentials{
		Server:   DServer,
		Account:  DAccount,
		Password: DPassword,
		Source:   false,
	}

	log.Infof("Validating credentials for: %s", creds.Account)

	err = internal.ValidateCredentials(creds)
	if err != nil {
		ctx.HTML(200, "error.html", "Couldn't verify for user: "+DAccount)
		return
	}

	ctx.HTML(200, "success.html", nil)
}
