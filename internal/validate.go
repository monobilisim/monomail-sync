package internal

import (
	"errors"

	"github.com/emersion/go-imap/client"
)

func ValidateCredentials(creds Credentials) error {
	if creds.Server == "" {
		return errors.New("server cannot be empty")
	}
	if creds.Account == "" {
		return errors.New("account cannot be empty")
	}
	if creds.Password == "" {
		return errors.New("password cannot be empty")
	}

	c, err := client.DialTLS(creds.Server+":993", nil)
	if err != nil {
		return err
	}
	log.Infof("IMAP Connected")

	defer c.Logout()

	if err := c.Login(creds.Account, creds.Password); err != nil {
		log.Error(err)
		return err
	}
	log.Infof("User %s Logged in to server %s", creds.Account, creds.Server)

	return nil
}
