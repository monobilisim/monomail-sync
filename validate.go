package main

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
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		htmx := fmt.Sprintf(`<div class="alert alert-error max-w-sm absolute top-0 right-0 mt-4 mr-4 cursor-pointer" id="error-notification">
        <svg width="48" height="48" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path fill-rule="evenodd" clip-rule="evenodd"
                d="M24 4C12.96 4 4 12.96 4 24C4 35.04 12.96 44 24 44C35.04 44 44 35.04 44 24C44 12.96 35.04 4 24 4ZM24 26C22.9 26 22 25.1 22 24V16C22 14.9 22.9 14 24 14C25.1 14 26 14.9 26 16V24C26 25.1 25.1 26 24 26ZM26 34H22V30H26V34Z"
                fill="#E92C2C" />
        </svg>
        <div class="flex w-full justify-between">
            <div class="flex flex-col">
                <span>Error</span>
                <span class="text-content2">%s</span>
            </div>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" clip-rule="evenodd"
                    d="M18.3007 5.71C17.9107 5.32 17.2807 5.32 16.8907 5.71L12.0007 10.59L7.1107 5.7C6.7207 5.31 6.0907 5.31 5.7007 5.7C5.3107 6.09 5.3107 6.72 5.7007 7.11L10.5907 12L5.7007 16.89C5.3107 17.28 5.3107 17.91 5.7007 18.3C6.0907 18.69 6.7207 18.69 7.1107 18.3L12.0007 13.41L16.8907 18.3C17.2807 18.69 17.9107 18.69 18.3007 18.3C18.6907 17.91 18.6907 17.28 18.3007 16.89L13.4107 12L18.3007 7.11C18.6807 6.73 18.6807 6.09 18.3007 5.71Z"
                    fill="#969696" />
            </svg>
        </div>
    </div>`, err.Error())
		ctx.String(200, htmx)
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
		return fmt.Errorf("Invalid credentials for account: %s", creds.Account)
	}

	return nil
}
