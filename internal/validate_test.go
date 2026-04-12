package internal

import (
	"testing"
)

func TestValidateCredentials_EmptyFields(t *testing.T) {
	tests := []struct {
		name  string
		creds Credentials
	}{
		{"empty server", Credentials{Server: "", Account: "user", Password: "pass"}},
		{"empty account", Credentials{Server: "imap.example.com", Account: "", Password: "pass"}},
		{"empty password", Credentials{Server: "imap.example.com", Account: "user", Password: ""}},
		{"all empty", Credentials{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCredentials(tt.creds, false)
			if err == nil {
				t.Error("expected error for missing fields")
			}
		})
	}
}

func TestValidateCredentials_InvalidServer(t *testing.T) {
	creds := Credentials{
		Server:   "nonexistent.invalid.host.example:993",
		Account:  "user@example.com",
		Password: "password123",
	}

	err := ValidateCredentials(creds, true)
	if err == nil {
		t.Error("expected error for invalid server")
	}
}

func TestMailboxStats_InvalidServer(t *testing.T) {
	_, err := GetMailboxStats("nonexistent.invalid.host.example:993", "user", "pass", true)
	if err == nil {
		t.Error("expected error for invalid server")
	}
}

func TestMailboxStats_AutoTLSDetection(t *testing.T) {
	_, err := GetMailboxStats("nonexistent.invalid.host.example:993", "user", "pass", false)
	if err == nil {
		t.Error("expected connection error")
	}
}

func TestCredentials_Struct(t *testing.T) {
	c := Credentials{
		Server:   "imap.gmail.com",
		Account:  "user@gmail.com",
		Password: "secret",
		Source:   true,
	}

	if c.Server != "imap.gmail.com" {
		t.Errorf("Server = %q", c.Server)
	}
	if !c.Source {
		t.Error("Source should be true")
	}
}

func TestMailboxInfo_Struct(t *testing.T) {
	info := MailboxInfo{
		Name:     "INBOX",
		Messages: 100,
		Size:     1048576,
		SizeMB:   1.0,
	}

	if info.Name != "INBOX" {
		t.Errorf("Name = %q", info.Name)
	}
	if info.Messages != 100 {
		t.Errorf("Messages = %d", info.Messages)
	}
	if info.SizeMB != 1.0 {
		t.Errorf("SizeMB = %f", info.SizeMB)
	}
}

func TestMailboxStats_Struct(t *testing.T) {
	stats := MailboxStats{
		Mailboxes: []MailboxInfo{
			{Name: "INBOX", Messages: 50, Size: 524288, SizeMB: 0.5},
			{Name: "Sent", Messages: 30, Size: 262144, SizeMB: 0.25},
		},
		TotalMessages: 80,
		TotalSize:     786432,
		TotalSizeMB:   0.75,
	}

	if len(stats.Mailboxes) != 2 {
		t.Errorf("Mailboxes length = %d, want 2", len(stats.Mailboxes))
	}
	if stats.TotalMessages != 80 {
		t.Errorf("TotalMessages = %d, want 80", stats.TotalMessages)
	}
}
