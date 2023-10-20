package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// func syncIMAP( ,details *Task) error {

// 	updateTaskStatus(details, "In Progress")

// 	currentTime := time.Now().Format("2006.01.02_15:04:05")

// 	logname := details.SourceAccount + "_" + details.DestinationAccount + "_" + currentTime + ".log"

// 	cmd := exec.Command("imapsync",
// 		"--host1", details.SourceServer,
// 		"--user1", details.SourceAccount,
// 		"--password1", details.SourcePassword,
// 		"--host2", details.DestinationServer,
// 		"--user2", details.DestinationAccount,
// 		"--password2", details.DestinationPassword,
// 		"--logfile", logname)

// 	updateTaskLogFile(details, logname)
// 	var stdBuffer bytes.Buffer
// 	mw := io.MultiWriter(os.Stdout, &stdBuffer)

// 	cmd.Stdout = mw
// 	cmd.Stderr = mw

// 	if err := cmd.Run(); err != nil {
// 		updateTaskStatus(details, "Error")
// 		return fmt.Errorf("error running imapsync: %w", err)
// 	}

// 	updateTaskStatus(details, "Done")

// 	return nil
// }

func syncIMAP(ctx context.Context, details *Task) error {
	updateTaskStatus(details, "In Progress")

	currentTime := time.Now().Format("2006.01.02_15:04:05")

	logname := details.SourceAccount + "_" + details.DestinationAccount + "_" + currentTime + ".log"

	cmd := exec.CommandContext(ctx, "imapsync",
		"--host1", details.SourceServer,
		"--user1", details.SourceAccount,
		"--password1", details.SourcePassword,
		"--host2", details.DestinationServer,
		"--user2", details.DestinationAccount,
		"--password2", details.DestinationPassword,
		"--logfile", logname)

	updateTaskLogFile(details, logname)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Start(); err != nil {
		updateTaskStatus(details, "Error")
		return fmt.Errorf("error starting imapsync: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("failed to terminate imapsync process: %w", err)
		}
		updateTaskStatus(details, "Cancelled")
		return ctx.Err()
	case err := <-done:
		if err != nil {
			updateTaskStatus(details, "Error")
			return fmt.Errorf("error running imapsync: %w", err)
		}
		updateTaskStatus(details, "Done")
		return nil
	}
}
