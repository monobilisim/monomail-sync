package internal

import (
	"errors"
	"fmt"
	"net/smtp"

	"github.com/gin-gonic/gin"
)

func handleValidate(ctx *gin.Context) {
	Server := ctx.PostForm("server")
	Account := ctx.PostForm("account")
	Password := ctx.PostForm("password")

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

	auth := smtp.PlainAuth("", creds.Account, creds.Password, creds.Server)
	err := smtp.SendMail(creds.Server+":25", auth, creds.Account, []string{creds.Account}, []byte("test"))
	if err != nil {
		return fmt.Errorf("invalid credentials for account: %s", creds.Account)
	}

	return nil
}
