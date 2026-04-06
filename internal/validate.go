package internal

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/emersion/go-imap/client"
)

type AccountValidation struct {
	SourceServer        string
	SourceUser          string
	SourcePassword      string
	SourceUseTLS        bool
	DestinationServer   string
	DestinationUser     string
	DestinationPassword string
	DestinationUseTLS   bool
}

type ValidationResult struct {
	SourceValid      bool
	SourceError      string
	DestinationValid bool
	DestinationError string
}

func ValidateAccount(cred *AccountValidation) *ValidationResult {
	result := &ValidationResult{}

	result.SourceValid, result.SourceError = testIMAPConnection(
		cred.SourceServer, cred.SourceUser, cred.SourcePassword, cred.SourceUseTLS)

	result.DestinationValid, result.DestinationError = testIMAPConnection(
		cred.DestinationServer, cred.DestinationUser, cred.DestinationPassword, cred.DestinationUseTLS)

	return result
}

func testIMAPConnection(server, username, password string, useTLS bool) (bool, string) {
	if server == "" {
		return false, "Server address is required"
	}
	if username == "" {
		return false, "Username is required"
	}
	if password == "" {
		return false, "Password is required"
	}

	host, port, err := net.SplitHostPort(server)
	if err != nil {
		host = server
		if useTLS {
			port = "993"
		} else {
			port = "143"
		}
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	var c *client.Client
	if useTLS {
		c, err = client.DialTLS(addr, &tls.Config{MinVersion: tls.VersionTLS12})
	} else {
		c, err = client.Dial(addr)
	}

	if err != nil {
		return false, fmt.Sprintf("Connection failed: %v", err)
	}
	defer c.Logout()

	if err := c.Login(username, password); err != nil {
		return false, fmt.Sprintf("Login failed: %v", err)
	}

	return true, ""
}

func ValidateCredentials(creds Credentials, useTLS bool) error {
	server := creds.Server

	if !useTLS && strings.Contains(server, ":") {
		parts := strings.Split(server, ":")
		port := parts[len(parts)-1]
		if port == "993" {
			useTLS = true
		}
	}

	valid, err := testIMAPConnection(server, creds.Account, creds.Password, useTLS)
	if !valid {
		return fmt.Errorf("%v", err)
	}
	return nil
}

type BulkMigration struct {
	ID                int
	SourceServer      string
	SourceUseTLS      bool
	DestinationServer string
	DestinationUseTLS bool
	Accounts          []BulkAccount
	CreatedAt         int64
	Status            string
}

type BulkAccount struct {
	SourceUser          string
	SourcePassword      string
	DestinationUser     string
	DestinationPassword string
	Status              string
	Error               string
	Progress            int
	TotalMessages       int
	CopiedMessages      int
}

func ParseBulkAccounts(csvContent string) []BulkAccount {
	lines := strings.Split(csvContent, "\n")
	var accounts []BulkAccount

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) >= 4 {
			account := BulkAccount{
				SourceUser:          strings.TrimSpace(parts[0]),
				SourcePassword:      strings.TrimSpace(parts[1]),
				DestinationUser:     strings.TrimSpace(parts[2]),
				DestinationPassword: strings.TrimSpace(parts[3]),
				Status:              "pending",
				Progress:            0,
			}
			accounts = append(accounts, account)
		}
	}

	return accounts
}
