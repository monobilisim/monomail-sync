package controller

import (
	"imap-sync/internal"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

func HandleValidate(ctx *gin.Context) {

	validate := ctx.PostForm("validate")
	submitsync := ctx.PostForm("submit_sync")

	var SServer, SAccount, SPassword string
	var DServer, DAccount, DPassword string
	var SUseTLS, DUseTLS bool

	if validate != "" {
		SServer = ctx.PostForm("source_server")
		SAccount = ctx.PostForm("source_account")
		SPassword = ctx.PostForm("source_password")
		SUseTLS = ctx.PostForm("source_use_tls") == "on"
		DServer = ctx.PostForm("destination_server")
		DAccount = ctx.PostForm("destination_account")
		DPassword = ctx.PostForm("destination_password")
		DUseTLS = ctx.PostForm("destination_use_tls") == "on"
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

	log.Infof("Validating source credentials for: %s @ %s (TLS: %v)", creds.Account, creds.Server, SUseTLS)

	err := internal.ValidateCredentials(creds, SUseTLS)
	if err != nil {
		log.Errorf("Source validation failed for %s: %v", SAccount, err)
		ctx.HTML(200, "error.html", "Couldn't verify for user: "+SAccount)
		return
	}

	creds = internal.Credentials{
		Server:   DServer,
		Account:  DAccount,
		Password: DPassword,
		Source:   false,
	}

	log.Infof("Validating destination credentials for: %s @ %s (TLS: %v)", creds.Account, creds.Server, DUseTLS)

	err = internal.ValidateCredentials(creds, DUseTLS)
	if err != nil {
		log.Errorf("Destination validation failed for %s: %v", DAccount, err)
		ctx.HTML(200, "error.html", "Couldn't verify for user: "+DAccount)
		return
	}

	ctx.HTML(200, "success.html", nil)
}

func HandleValidateStep(ctx *gin.Context) {
	side := ctx.Query("side")

	var server, account, password string
	var useTLS bool

	switch side {
	case "source":
		server = ctx.PostForm("source_server")
		account = ctx.PostForm("source_account")
		password = ctx.PostForm("source_password")
		useTLS = ctx.PostForm("source_use_tls") == "on"
	case "destination":
		server = ctx.PostForm("destination_server")
		account = ctx.PostForm("destination_account")
		password = ctx.PostForm("destination_password")
		useTLS = ctx.PostForm("destination_use_tls") == "on"
	default:
		ctx.JSON(400, gin.H{"success": false, "error_type": "invalid", "message": "Invalid side parameter"})
		return
	}

	creds := internal.Credentials{
		Server:   server,
		Account:  account,
		Password: password,
	}

	log.Infof("Validating %s credentials for: %s @ %s (TLS: %v)", side, account, server, useTLS)

	err := internal.ValidateCredentials(creds, useTLS)
	if err != nil {
		errMsg := err.Error()
		errorType := "unknown"
		if strings.Contains(errMsg, "Connection failed") {
			errorType = "connection"
		} else if strings.Contains(errMsg, "Login failed") {
			errorType = "auth"
		} else if strings.Contains(errMsg, "required") {
			errorType = "missing"
		}

		log.Errorf("%s validation failed for %s: %v", side, account, err)
		ctx.JSON(200, gin.H{
			"success":    false,
			"error_type": errorType,
			"message":    errMsg,
		})
		return
	}

	log.Infof("%s validation successful for: %s", side, account)
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "",
	})
}

func HandleMailboxStats(ctx *gin.Context) {
	sourceServer := ctx.PostForm("source_server")
	sourceAccount := ctx.PostForm("source_account")
	sourcePassword := ctx.PostForm("source_password")
	sourceUseTLS := ctx.PostForm("source_use_tls") == "on"

	destServer := ctx.PostForm("destination_server")
	destAccount := ctx.PostForm("destination_account")
	destPassword := ctx.PostForm("destination_password")
	destUseTLS := ctx.PostForm("destination_use_tls") == "on"

	var sourceStats, destStats *internal.MailboxStats
	var sourceErr, destErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		sourceStats, sourceErr = internal.GetMailboxStats(sourceServer, sourceAccount, sourcePassword, sourceUseTLS)
	}()

	go func() {
		defer wg.Done()
		destStats, destErr = internal.GetMailboxStats(destServer, destAccount, destPassword, destUseTLS)
	}()

	wg.Wait()

	response := gin.H{}

	if sourceErr != nil {
		response["source_error"] = sourceErr.Error()
	} else {
		response["source"] = sourceStats
	}

	if destErr != nil {
		response["dest_error"] = destErr.Error()
	} else {
		response["destination"] = destStats
	}

	ctx.JSON(200, response)
}
