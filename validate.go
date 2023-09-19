package main

import (
	"errors"
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
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.String(200, "<div id=\"error-notification\" class=\"bg-red-500 text-white p-2 rounded absolute top-0 right-0 mt-4 mr-4 cursor-pointer\">Error validating credentials</div>")
		return
	}

}

func validateCredentials(creds Credentials) error {
	if creds.Server == "" {
		return errors.New("Server cannot be empty")
	}
	if creds.Account == "" {
		return errors.New("Account cannot be empty")
	}
	if creds.Password == "" {
		return errors.New("Password cannot be empty")
	}

	auth := smtp.PlainAuth("", creds.Account, creds.Password, creds.Server)
	err := smtp.SendMail(creds.Server+":25", auth, creds.Account, []string{creds.Account}, []byte("test"))
	if err != nil {
		return errors.New("Invalid credentials")
	}

	return nil
}
