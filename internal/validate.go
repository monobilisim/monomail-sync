package internal

import (
	"errors"

	"github.com/emersion/go-imap/client"
	"github.com/gin-gonic/gin"
)

func handleValidate(ctx *gin.Context) {

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

	creds := Credentials{
		Server:   Server,
		Account:  Account,
		Password: Password,
	}

	log.Infof("Validating credentials for: %s", creds.Account)

	err := validateCredentials(creds)
	if err != nil {
		ctx.HTML(200, "error.html", err.Error())
		return
	}
	ctx.HTML(200, "success.html", creds)
}

func validateCredentials(creds Credentials) error {
	if creds.Server == "" {
		return errors.New("server cannot be empty")
	}
	if creds.Account == "" {
		return errors.New("account cannot be empty")
	}
	if creds.Password == "" {
		return errors.New("password cannot be empty")
	}

	c, err := client.DialTLS(creds.Server+":993", nil)
	if err != nil {
		return err
	}
	log.Infof("IMAP Connected")

	defer c.Logout()

	if err := c.Login(creds.Account, creds.Password); err != nil {
		log.Error(err)
		return err
	}
	log.Infof("User %s Logged in to server %s", creds.Account, creds.Server)

	return nil
}
