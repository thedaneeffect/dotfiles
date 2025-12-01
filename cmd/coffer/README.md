# coffer - Encrypted File Management

Medieval-themed CLI for managing encrypted dotfiles.

**Key Achievement:** OpenSSL-compatible encryption - can decrypt bash-encrypted files! ✅

## Features

- **Local registry** - Track files in `~/.coffer` (or `~/.coffer.<group>`)
- **Encrypted sync** - Push/pull to Cloudflare Workers
- **Groups** - Organize files by purpose (github, work, personal, etc.)
- **OpenSSL compatible** - AES-256-CBC with PBKDF2
- **Comprehensive tests** - Unit + integration coverage

## Installation

```bash
# Install from source (during setup.sh or manually)
cd ~/projects/dotfiles/cmd/coffer
go install

# Verify
coffer --version
coffer --help
```

The setup script automatically runs `go install` if Go is available.

## Quick Start

```bash
# Add files to track
coffer add ~/.ssh/id_rsa ~/.env

# Configure worker (optional, for sync)
export COFFER_URL="https://secrets.your-subdomain.workers.dev"
export COFFER_PASSPHRASE="your-passphrase"

# Push to worker
coffer push

# On another machine: pull files
coffer pull
```

## Commands

### Local Operations

```bash
coffer add <path> [<path>...] [-g GROUP]     # Add files/directories
coffer remove <path> [<path>...] [-g GROUP]  # Remove from registry
coffer list [GROUP]                           # List tracked files
coffer status [GROUP]                         # Show file status
```

### Worker Operations

```bash
coffer push [GROUP]      # Upload to worker
coffer pull [GROUP]      # Download from worker
coffer delete [GROUP]    # Delete from worker
coffer groups            # List all groups
```

### Aliases

- `rm` → `remove`
- `ls` → `list`
- `st` → `status`
- `del` → `delete`

## Groups

Organize files by purpose:

```bash
# Default group
coffer add ~/.env
coffer push

# Named groups
coffer add ~/.ssh/github_rsa -g github
coffer add ~/.aws/credentials -g work

coffer push github
coffer pull work
```

Each group has its own registry file:
- `~/.coffer` (default)
- `~/.coffer.github`
- `~/.coffer.work`

## Configuration

Set environment variables:

```bash
export COFFER_URL="https://secrets.your-subdomain.workers.dev"
export COFFER_PASSPHRASE="your-secret-passphrase"
```

Add to `~/.zshrc` to persist.

## Worker Setup

Deploy the Cloudflare Worker from `/worker` directory in dotfiles repo.

See main dotfiles README for worker deployment instructions.

## Testing

```bash
# All tests (local only)
go test ./... -v

# Specific test suites
go test ./internal/crypto -v              # Crypto + OpenSSL compatibility
go test ./cmd -v -run Local               # Local command integration

# Worker integration tests (requires COFFER_URL and COFFER_PASSPHRASE)
export COFFER_URL="https://your-worker.workers.dev"
export COFFER_PASSPHRASE="test-passphrase"
go test ./cmd -v -run Integration
```

Integration tests automatically skip if worker credentials aren't configured.

## Architecture

```
cmd/coffer/
├── main.go              # Entry point
├── cmd/                 # Cobra commands
│   ├── root.go
│   ├── add.go
│   ├── remove.go
│   ├── list.go
│   ├── status.go
│   ├── push.go
│   ├── pull.go
│   ├── delete.go
│   └── groups.go
├── internal/
│   ├── api/             # Worker HTTP client
│   ├── config/          # Environment config
│   ├── crypto/          # Encryption (OpenSSL compatible)
│   ├── registry/        # Local file tracking
│   └── ui/              # Terminal output
└── *.md                 # Documentation
```

## Examples

```bash
# Track SSH keys for github
coffer add ~/.ssh/github_rsa ~/.ssh/github_rsa.pub -g github
coffer push github

# Track work credentials
coffer add ~/.aws/credentials ~/.kube/config -g work
coffer push work

# On new machine
coffer pull github
coffer pull work
coffer status

# Clean up
coffer delete test-group
```

## Limits

- 25MB per group (Cloudflare KV limit)
- Free tier: 1GB storage, 100k reads/day, 1k writes/day

## License

Part of thedaneeffect/dotfiles repository.
