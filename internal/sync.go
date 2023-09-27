package internal

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func handleSync(ctx *gin.Context) {
	sourceServer := ctx.PostForm("sourceServer")
	sourceAccount := ctx.PostForm("sourceAccount")
	sourcePassword := ctx.PostForm("sourcePassword")
	destinationServer := ctx.PostForm("destinationServer")
	destinationAccount := ctx.PostForm("destinationAccount")
	destinationPassword := ctx.PostForm("destinationPassword")

	sourceDetails := Credentials{
		Server:   sourceServer,
		Account:  sourceAccount,
		Password: sourcePassword,
	}

	destinationDetails := Credentials{
		Server:   destinationServer,
		Account:  destinationAccount,
		Password: destinationPassword,
	}

	log.Infof("Syncing %s to %s", sourceDetails.Account, destinationDetails.Account)

	err := syncIMAP(sourceDetails, destinationDetails)
	if err != nil {
		ctx.HTML(200, "error.html", err.Error())
		return
	}
	ctx.HTML(200, "success.html", "Synced "+sourceDetails.Account+" to "+destinationDetails.Account)
}

func syncIMAP(sourceDetails Credentials, destinationDetails Credentials) error {
	cmd := exec.Command("imapsync", "--host1", sourceDetails.Server, "--user1", sourceDetails.Account, "--password1", sourceDetails.Password, "--host2", destinationDetails.Server, "--user2", destinationDetails.Account, "--password2", destinationDetails.Password)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Run(); err != nil {
		return err
	}

	// Command output realtime
	// log.Println(stdBuffer.String())

	return nil
}
