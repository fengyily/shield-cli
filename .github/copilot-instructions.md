# Shield CLI — Copilot Instructions

## About

Shield CLI is a browser-first internal service gateway written in Go. It creates encrypted WebSocket tunnels to expose internal services (SSH, RDP, VNC, HTTP, databases) through any browser with a single command. Unlike traditional tunnel tools, Shield CLI renders protocols directly in the browser via HTML5.

- **Repository**: https://github.com/fengyily/shield-cli
- **Docs**: https://docs.yishield.com
- **License**: Apache 2.0

## Installation

```bash
# macOS
brew tap fengyily/tap && brew install shield-cli

# Windows
scoop bucket add shield https://github.com/fengyily/scoop-bucket && scoop install shield-cli

# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh

# China mirror
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh

# Docker
docker run -d --name shield --network host --restart unless-stopped fengyily/shield-cli
```

## Usage

```bash
shield start              # Launch Web UI at http://localhost:8181
shield ssh 10.0.0.5       # SSH terminal in browser
shield rdp 10.0.0.5       # Windows desktop in browser
shield mysql 10.0.0.20    # Database admin in browser (plugin)
shield http 3000           # Expose local web app
shield vnc 10.0.0.10       # VNC screen sharing in browser
```

## Project Structure

- `cmd/` — CLI commands (spf13/cobra)
- `tunnel/` — Chisel tunnel management
- `web/` — HTTP server + embedded Web UI
- `plugin/` — Plugin system (subprocess JSON IPC)
- `config/` — Encrypted credentials & app storage
- `service/` — System service (launchd/systemd/Windows)
- `plugins/mysql/` — MySQL web admin plugin
- `docs/` — VitePress documentation

## Code Conventions

- Go 1.25.0, `gofmt` formatting
- AES-256-GCM encryption keyed to machine fingerprint
- Config: `~/.shield-cli/` (macOS/Linux), `%LOCALAPPDATA%\ShieldCLI\` (Windows)
- Plugins: external binaries, JSON IPC over stdin/stdout
