package controller

import (
	"imap-sync/internal"

	"github.com/gin-gonic/gin"
)

func HandleValidate(ctx *gin.Context) {

	sourceCreds := ctx.PostForm("source_creds")
	destCreds := ctx.PostForm("destination_creds")
	submitsync := ctx.PostForm("submit_sync")

	var Server, Account, Password string

	if sourceCreds != "" {
		Server = ctx.PostForm("source_server")
		Account = ctx.PostForm("source_account")
		Password = ctx.PostForm("source_password")
	}

	if destCreds != "" {
		Server = ctx.PostForm("destination_server")
		Account = ctx.PostForm("destination_account")
		Password = ctx.PostForm("destination_password")
	}

	if destCreds == "" && sourceCreds == "" && submitsync != "" {
		handleSync(ctx)
		return
	}

	creds := internal.Credentials{
		Server:   Server,
		Account:  Account,
		Password: Password,
	}

	log.Infof("Validating credentials for: %s", creds.Account)

	err := internal.ValidateCredentials(creds)
	if err != nil {
		ctx.HTML(200, "error.html", err.Error())
		return
	}
	ctx.HTML(200, "success.html", creds)
}
