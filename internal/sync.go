package internal

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func syncIMAP(details Task) error {

	cmd := exec.Command("imapsync", "--host1", details.SourceServer, "--user1", details.SourceAccount, "--password1", details.SourcePassword, "--host2", details.DestinationServer, "--user2", details.DestinationAccount, "--password2", details.DestinationPassword)
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
