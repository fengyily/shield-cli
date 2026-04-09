# Shield CLI

Shield CLI is a browser-first internal service gateway. One binary, one command — SSH terminals, remote desktops, database admin, web apps — all accessible through any browser.

## Tech Stack

- **Language**: Go 1.25.0
- **Web UI**: Embedded Vue.js single-page app (`web/static/index.html`)
- **CLI Framework**: spf13/cobra
- **Tunnel**: Chisel WebSocket tunnel

## Architecture

```
Internal Service ←→ Shield CLI (Chisel WebSocket Tunnel) ←→ Public Gateway (HTML5 Render) ←→ Browser
```

## Key Packages

- `cmd/` — CLI commands (cobra), protocol routing, plugin discovery
- `tunnel/` — Chisel-based tunnel management (main + per-app resource tunnels)
- `web/` — HTTP server, connection manager, embedded Web UI
- `plugin/` — Subprocess-based plugin system (stdin/stdout JSON IPC)
- `config/` — AES-256-GCM encrypted credentials & app storage (machine fingerprint key)
- `service/` — Cross-platform system service (launchd/systemd/Windows)
- `plugins/mysql/` — Built-in MySQL web admin plugin

## Conventions

- Go code follows standard `gofmt` formatting
- CLI commands defined via spf13/cobra
- Config files stored in `~/.shield-cli/` (macOS/Linux) or `%LOCALAPPDATA%\ShieldCLI\` (Windows)
- Plugin binaries stored in `~/.shield-cli/plugins/`
- Web UI is a single embedded `index.html` in `web/static/`
- Documentation site uses VitePress, source in `docs/`

## Build

```bash
go build -o shield .
```
