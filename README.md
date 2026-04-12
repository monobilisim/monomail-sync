[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-shield]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

# monomail-sync

A web interface for [imapsync](https://imapsync.lamiral.info/). Sync emails between IMAP servers with a clean UI.

## Quick Start

### Prerequisites

- Go 1.21+
- [imapsync](https://imapsync.lamiral.info/) installed and available in `$PATH`
- SQLite3

### Run

```bash
export MONOMAIL_SYNC_ENCRYPTION_KEY="$(openssl rand -base64 32)"
go run .
```

Access at http://localhost:8000 — Login: `admin` / `admin`

### Systemd (Production)

```bash
go build -o /usr/local/bin/monomail-sync .

sudo useradd -r -s /usr/sbin/nologin -d /var/lib/monomail-sync monomail-sync
sudo mkdir -p /var/lib/monomail-sync/{templates,static,LOG_imapsync,logs}
sudo cp -r templates/* /var/lib/monomail-sync/templates/
sudo cp -r static/* /var/lib/monomail-sync/static/
sudo cp favicon.ico /var/lib/monomail-sync/
sudo cp config/config.yml /etc/monomail-sync.yml
sudo chown -R monomail-sync:monomail-sync /var/lib/monomail-sync
```

Create `/etc/systemd/system/monomail-sync.service`:

```ini
[Unit]
Description=Monomail Sync - IMAP Email Sync Service
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=monomail-sync
Group=monomail-sync
ExecStart=/usr/local/bin/monomail-sync
WorkingDirectory=/var/lib/monomail-sync
Environment=MONOMAIL_SYNC_ENCRYPTION_KEY=<your-32-byte-key>
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/monomail-sync /var/log/monomail-sync
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now monomail-sync
```

## Features

### Single Account Sync
Validate source/destination IMAP credentials, compare mailbox sizes, and start synchronization — all from a single form.

### Real-time Progress Tracking
Each running sync task displays a live progress bar parsed directly from imapsync output:
- Percentage complete with progress bar
- Current folder being synced (e.g. `[INBOX]`, `[Sent]`)
- Message count (`216/436 msgs`)
- Transfer speed (`11.39 msgs/s`)
- ETA countdown

Progress turns green with a checkmark when the task finishes.

### Mailbox Comparison
Before syncing, a comparison table shows folder-by-folder message counts and sizes for both source and destination. The table automatically refreshes after sync completes to verify the transfer.

### Bulk Migration
Sync multiple accounts at once by providing a list of `source_user:source_pass:dest_user:dest_pass` entries. Includes:
- Per-account status tracking in a collapsible panel
- Live activity log at the bottom of the page with timestamped events
- Toast notifications for completed/failed accounts

### Dashboard
Clickable stat cards (Total, Completed, Failed, In Progress) — clicking a card filters the task table by that status.

### Admin Panel (`/admin`)
- IMAP server settings with provider presets (Gmail, Outlook, Yahoo, etc.)
- Full task queue with search
- Active session management
- Audit log
- System info (uptime, memory, goroutines)
- Filterable stat cards

### Other
- Encrypted credentials at rest (AES-256-GCM)
- Task cancellation and retry
- Email notifications on sync completion/failure (SMTP)
- Sync log viewer
- Dark/light mode
- Multi-language support (English, Turkish)
- Health check endpoint (`/health`)

## Configuration

Default config path: `/etc/monomail-sync.yml`

Override with: `monomail-sync -config /path/to/config.yml`

```yaml
language: en        # en | tr
port: 8000
databaseInfo:
  adminName: admin
  adminPass: admin
  databasePath: ./db.db
sourceAndDestination:
  SourceServer: "imap.example.com:993"
  SourceMail: "@example.com"
  DestinationServer: "imap.example.com:993"
  DestinationMail: "@example.com"
email:
  SMTPHost: smtp.example.com
  SMTPPort: "587"
  SMTPFrom: noreply@example.com
  SMTPUser: noreply@example.com
  SMTPPassword: password
```

## Encryption Key

Task credentials (source/destination passwords) are encrypted at rest using AES-256-GCM. Set the key via environment variable before starting:

```bash
# Base64-encoded (recommended)
export MONOMAIL_SYNC_ENCRYPTION_KEY="$(openssl rand -base64 32)"

# Or raw 32-character string
export MONOMAIL_SYNC_ENCRYPTION_KEY="0123456789abcdef0123456789abcdef"
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Main sync interface |
| `GET` | `/admin` | Admin panel |
| `GET` | `/login` | Login page |
| `GET` | `/health` | Health check (no auth) |
| `POST` | `/api/validate` | Validate IMAP credentials |
| `POST` | `/api/validate/step` | Step-by-step validation |
| `POST` | `/api/validate/stats` | Mailbox statistics |
| `GET` | `/api/sync` | Start/cancel/retry sync |
| `POST` | `/api/bulk` | Start bulk migration |
| `GET` | `/api/bulk/status` | Bulk migration status |
| `GET` | `/api/task/progress` | Real-time task progress |
| `GET` | `/api/stats` | Dashboard statistics |
| `GET` | `/api/queue` | Task queue |
| `POST` | `/api/search` | Search tasks |
| `GET` | `/api/system` | System info |
| `GET` | `/api/audit` | Audit log |
| `GET` | `/api/sessions` | Active sessions |

## Tech Stack

- **Backend:** Go, Gin, SQLite
- **Frontend:** HTMX, Alpine.js, Tailwind CSS
- **Sync engine:** [imapsync](https://imapsync.lamiral.info/)

## Project Structure

```
├── api/             # HTTP router and middleware
├── config/          # Configuration parsing
├── controller/      # Request handlers
├── internal/        # Core logic
│   ├── bulk.go      # Bulk migration orchestration
│   ├── crypto.go    # AES-256-GCM encryption
│   ├── db.go        # SQLite database operations
│   ├── progress.go  # Real-time imapsync output parser
│   ├── queue.go     # Task queue (linked list)
│   ├── sync.go      # imapsync process management
│   └── validate.go  # IMAP credential validation
├── templates/       # Go HTML templates
├── static/          # CSS, JS assets
└── main.go
```

## License

GPL-3.0

[contributors-shield]: https://img.shields.io/github/contributors/monobilisim/monomail-sync?style=for-the-badge
[contributors-url]: https://github.com/monobilisim/monomail-sync/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/monobilisim/monomail-sync?style=for-the-badge
[forks-url]: https://github.com/monobilisim/monomail-sync/network/members
[stars-shield]: https://img.shields.io/github/stars/monobilisim/monomail-sync?style=for-the-badge
[stars-url]: https://github.com/monobilisim/monomail-sync/stargazers
[issues-shield]: https://img.shields.io/github/issues/monobilisim/monomail-sync?style=for-the-badge
[issues-url]: https://github.com/monobilisim/monomail-sync/issues
[license-shield]: https://img.shields.io/github/license/monobilisim/monomail-sync?style=for-the-badge
[license-url]: https://github.com/monobilisim/monomail-sync/blob/master/LICENSE
