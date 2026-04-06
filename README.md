[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-shield]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

# monomail-sync

A web interface for imapsync. Sync emails between IMAP servers with a clean UI.

## Quick Start

```bash
export MONOMAIL_SYNC_ENCRYPTION_KEY="<32-byte key or base64-encoded 32 bytes>"
go run .
```

## Encryption Key

Task credentials are encrypted at rest. Set `MONOMAIL_SYNC_ENCRYPTION_KEY` before starting the app.

Generate a key (base64, recommended):

```bash
export MONOMAIL_SYNC_ENCRYPTION_KEY="$(openssl rand -base64 32)"
```

Or use a raw 32-byte key (exactly 32 chars):

```bash
export MONOMAIL_SYNC_ENCRYPTION_KEY="0123456789abcdef0123456789abcdef"
```

Access at http://localhost:8000

Admin login: `admin` / `admin`

## Configuration (Optional)

Create `config/config.yml` to override defaults:

```yaml
language: en
port: 8000
databaseInfo:
  adminName: admin
  adminPass: admin
  databasePath: ./db.db
sourceAndDestination:
  SourceServer: imap.example.com
  SourceMail: "@example.com"
  DestinationServer: imap.example.com
  DestinationMail: "@example.com"
email:
  SMTPHost: smtp.example.com
  SMTPPort: 587
  SMTPFrom: example
  SMTPUser: example@example.com
  SMTPPassword: password
```

Or run with custom config:
```bash
go run . -config config/config.yml
```

## Features

- Single user/bulk email sync
- Real-time progress tracking
- Dashboard with sync statistics
- IMAP account validation
- Health check endpoint (`/health`)
- Log rotation

## Endpoints

- `/` - Main sync interface
- `/admin` - Admin panel with settings
- `/login` - Admin login
- `/health` - Health check for monitoring

## Tech Stack

- Go + Gin
- SQLite
- HTMX + Alpine.js
- Tailwind CSS (daisyUI)

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
