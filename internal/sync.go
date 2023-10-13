package internal

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"time"
)

func syncIMAP(details *Task) error {

	details.Status = "In Progress"

	currentTime := time.Now().Format("2006.01.02_15:04:05")

	logname := details.SourceAccount + "_" + details.DestinationAccount + "_" + currentTime + ".log"

	cmd := exec.Command("imapsync",
		"--host1", details.SourceServer,
		"--user1", details.SourceAccount,
		"--password1", details.SourcePassword,
		"--host2", details.DestinationServer,
		"--user2", details.DestinationAccount,
		"--password2", details.DestinationPassword,
		"--logfile", logname)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Run(); err != nil {
		return err
	}

	details.Status = "Done"

	return nil
}
