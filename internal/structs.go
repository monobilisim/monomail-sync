package internal

type Task struct {
	ID                  int
	SourceAccount       string
	SourceServer        string
	SourcePassword      string
	DestinationAccount  string
	DestinationServer   string
	DestinationPassword string
	Status              string
	StartedAt           int64
	EndedAt             int64
	LogFile             string
}

type PageData struct {
	Index int
	Tasks []*Task
}

type Settings struct {
	SourceServer             string
	SourceAccountPrefix      string
	SourceUseTLS             bool
	DestinationServer        string
	DestinationAccountPrefix string
	DestinationUseTLS        bool
}

type Provider struct {
	Name   string
	Host   string
	Port   int
	UseTLS bool
}

var Providers = []Provider{
	{Name: "Custom", Host: "", Port: 993, UseTLS: false},
	{Name: "Gmail", Host: "imap.gmail.com", Port: 993, UseTLS: false},
	{Name: "Google Workspace", Host: "imap.gmail.com", Port: 993, UseTLS: false},
	{Name: "Outlook", Host: "outlook.office365.com", Port: 993, UseTLS: false},
	{Name: "Yahoo", Host: "imap.mail.yahoo.com", Port: 993, UseTLS: false},
	{Name: "Yandex", Host: "imap.yandex.com", Port: 993, UseTLS: false},
	{Name: "Zoho Mail", Host: "imap.zoho.com", Port: 993, UseTLS: false},
	{Name: "iCloud", Host: "imap.mail.me.com", Port: 993, UseTLS: false},
	{Name: "GMX", Host: "imap.gmx.com", Port: 993, UseTLS: false},
	{Name: "ProtonMail", Host: "127.0.0.1", Port: 1143, UseTLS: false},
}
