package internal

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

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
