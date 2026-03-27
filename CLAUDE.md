# Shield CLI

Shield CLI is a browser-first internal service gateway. One binary, one command — SSH terminals, remote desktops, database admin, web apps — all accessible through any browser.

## Project Overview

- **Language**: Go 1.25.0 with embedded Vue.js Web UI
- **License**: Apache 2.0
- **Repository**: https://github.com/fengyily/shield-cli
- **Documentation**: https://docs.yishield.com

## Architecture

```
Internal Service ←→ Shield CLI (Chisel WebSocket Tunnel) ←→ Public Gateway (HTML5 Render) ←→ Browser
```

Key packages:
- `cmd/` — CLI commands (cobra), protocol routing, plugin discovery
- `tunnel/` — Chisel-based tunnel management (main + per-app resource tunnels)
- `web/` — HTTP server, connection manager, embedded Web UI
- `plugin/` — Subprocess-based plugin system (stdin/stdout JSON IPC)
- `config/` — AES-256-GCM encrypted credentials & app storage (machine fingerprint key)
- `service/` — Cross-platform system service (launchd/systemd/Windows)
- `plugins/mysql/` — Built-in MySQL web admin plugin

## Installation

```bash
# macOS
brew tap fengyily/tap && brew install shield-cli

# Windows
scoop bucket add shield https://github.com/fengyily/scoop-bucket && scoop install shield-cli

# Linux / macOS (one-liner)
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh

# China mirror
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh

# Docker
docker run -d --name shield --network host --restart unless-stopped fengyily/shield-cli
```

## Quick Start

```bash
shield start                # Launch Web UI at http://localhost:8181
shield ssh 10.0.0.5         # SSH terminal in browser
shield rdp 10.0.0.5         # Windows desktop in browser
shield mysql 10.0.0.20      # Database admin in browser (plugin)
shield http 3000             # Expose local web app
shield vnc 10.0.0.10         # VNC screen sharing
```

## Smart Address Resolution

- `shield ssh` → `127.0.0.1:22`
- `shield ssh 2222` → `127.0.0.1:2222`
- `shield ssh 10.0.0.5` → `10.0.0.5:22`

## Plugin System

Plugins are external binaries communicating via stdin/stdout JSON. Available plugins: `mysql`, `postgres`, `sqlserver`.

```bash
shield plugin install mysql
shield plugin list
shield plugin uninstall mysql
```

## Build & Development

```bash
go build -o shield .
./shield --version
```

## Conventions

- Go code follows standard `gofmt` formatting
- CLI commands defined via `spf13/cobra`
- Credentials encrypted with AES-256-GCM, keyed to machine fingerprint (SHA256 of hostname + MAC + Machine ID)
- Config files stored in `~/.shield-cli/` (macOS/Linux) or `%LOCALAPPDATA%\ShieldCLI\` (Windows)
- Plugin binaries stored in `~/.shield-cli/plugins/`
- Web UI is a single embedded `index.html` in `web/static/`
- Documentation site uses VitePress, source in `docs/`
