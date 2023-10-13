package internal

import (
	"fmt"

	"github.com/emersion/go-imap/client"
)

func ValidateCredentials(creds Credentials) error {
	if creds.Server == "" {
		return fmt.Errorf("server cannot be empty")
	}
	if creds.Account == "" {
		return fmt.Errorf("account cannot be empty")
	}
	if creds.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	c, err := client.DialTLS(creds.Server+":993", nil)
	if err != nil {
		return fmt.Errorf("error connecting to server: %w", err)
	}

	defer c.Logout()

	if err := c.Login(creds.Account, creds.Password); err != nil {
		log.Error(err)
		return fmt.Errorf("error logging in: %w", err)
	}
	log.Infof("User %s Logged in to server %s", creds.Account, creds.Server)

	return nil
}
