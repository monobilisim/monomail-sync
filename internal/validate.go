package internal

import (
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"strings"

	"github.com/emersion/go-imap"
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
		c, err = client.DialTLS(addr, &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true})
	} else {
		c, err = client.Dial(addr)
		if err == nil && !c.IsTLS() {
			c.StartTLS(&tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true})
		}
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

type MailboxInfo struct {
	Name     string  `json:"name"`
	Messages uint32  `json:"messages"`
	Size     uint64  `json:"size"`
	SizeMB   float64 `json:"size_mb"`
}

type MailboxStats struct {
	Mailboxes     []MailboxInfo `json:"mailboxes"`
	TotalMessages uint32        `json:"total_messages"`
	TotalSize     uint64        `json:"total_size"`
	TotalSizeMB   float64       `json:"total_size_mb"`
}

func GetMailboxStats(server, username, password string, useTLS bool) (*MailboxStats, error) {
	if !useTLS && strings.Contains(server, ":") {
		parts := strings.Split(server, ":")
		port := parts[len(parts)-1]
		if port == "993" {
			useTLS = true
		}
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
		c, err = client.DialTLS(addr, &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true})
	} else {
		c, err = client.Dial(addr)
		if err == nil && !c.IsTLS() {
			c.StartTLS(&tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true})
		}
	}
	if err != nil {
		return nil, fmt.Errorf("Connection failed: %v", err)
	}
	defer c.Logout()

	if err := c.Login(username, password); err != nil {
		return nil, fmt.Errorf("Login failed: %v", err)
	}

	mailboxesCh := make(chan *imap.MailboxInfo, 100)
	listDone := make(chan error, 1)
	go func() {
		listDone <- c.List("", "*", mailboxesCh)
	}()

	type mbEntry struct {
		name     string
		noselect bool
	}
	var entries []mbEntry
	for m := range mailboxesCh {
		noselect := false
		for _, attr := range m.Attributes {
			if strings.EqualFold(attr, "\\Noselect") {
				noselect = true
				break
			}
		}
		entries = append(entries, mbEntry{name: m.Name, noselect: noselect})
	}
	if err := <-listDone; err != nil {
		return nil, fmt.Errorf("LIST failed: %v", err)
	}

	stats := &MailboxStats{}

	for _, entry := range entries {
		if entry.noselect {
			continue
		}

		mbox, err := c.Select(entry.name, true)
		if err != nil {
			continue
		}

		info := MailboxInfo{
			Name:     entry.name,
			Messages: mbox.Messages,
		}

		if mbox.Messages > 0 {
			seqset := new(imap.SeqSet)
			seqset.AddRange(1, mbox.Messages)

			messages := make(chan *imap.Message, 100)
			fetchDone := make(chan error, 1)
			go func() {
				fetchDone <- c.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822Size}, messages)
			}()

			var totalSize uint64
			for msg := range messages {
				totalSize += uint64(msg.Size)
			}
			<-fetchDone

			info.Size = totalSize
			info.SizeMB = math.Round(float64(totalSize)/(1024*1024)*100) / 100
		}

		stats.Mailboxes = append(stats.Mailboxes, info)
		stats.TotalMessages += info.Messages
		stats.TotalSize += info.Size
	}

	stats.TotalSizeMB = math.Round(float64(stats.TotalSize)/(1024*1024)*100) / 100

	return stats, nil
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
