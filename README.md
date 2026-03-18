<p align="center">
  <img src="docs/logo.svg" alt="Shield CLI" width="128" height="128">
</p>

<h1 align="center">Shield CLI</h1>

<p align="center">
  <strong>One command. One URL. Access anything from a browser.</strong><br>
  Shield CLI creates encrypted tunnels to your internal services — RDP desktops, VNC sessions, SSH terminals, web apps — and makes them accessible through any browser. No VPN. No client software. No port forwarding.
</p>

<p align="center">
  <a href="#installation">Installation</a> &bull;
  <a href="#quick-start">Quick Start</a> &bull;
  <a href="#how-it-works">How It Works</a> &bull;
  <a href="#usage">Usage</a> &bull;
  <a href="README_CN.md">中文文档</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-blue?logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-brightgreen" alt="Platform">
  <img src="https://img.shields.io/badge/license-Apache%202.0-green" alt="License">
</p>

---

## The Problem

Remote access today is painful. Want to reach an internal RDP desktop? Install a client. SSH into a server? Find a terminal app. Show a colleague your local web app? Set up port forwarding. Every protocol needs its own tool, every device needs its own setup.

## The Solution

Shield CLI replaces all of that with a single command:

```bash
shield ssh 10.0.0.2
```

That's it. You get a URL. Open it in any browser — phone, tablet, laptop, locked-down corporate machine — and you're connected. A full SSH terminal, RDP desktop, or VNC session, rendered in HTML5, right in the browser.

- **Zero client install** — if it has a browser, it works
- **Zero config** — protocol and address is all you need, everything else is automatic
- **One binary, all protocols** — SSH, RDP, VNC, HTTP, HTTPS, Telnet
- **Encrypted by default** — WebSocket tunnels with AES-256-GCM credential encryption

## Installation

### macOS (Homebrew)

```bash
brew tap fengyily/tap
brew install shield-cli
```

### Windows (Scoop)

```powershell
scoop bucket add shield https://github.com/fengyily/scoop-bucket
scoop install shield-cli
```

### Windows (PowerShell one-liner)

```powershell
irm https://raw.githubusercontent.com/fengyily/shield-cli/main/install.ps1 | iex
```

### Linux / macOS (curl one-liner)

```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
```

### Debian / Ubuntu (.deb)

```bash
# Download from GitHub Releases
sudo dpkg -i shield-cli_<version>_linux_amd64.deb
```

### RHEL / CentOS (.rpm)

```bash
sudo rpm -i shield-cli_<version>_linux_amd64.rpm
```

### From Source

```bash
git clone https://github.com/fengyily/shield-cli.git
cd shield-cli
go build -o shield .
```

## Quick Start

```bash
# Expose local SSH — access a terminal in your browser
shield ssh

# Expose a remote RDP desktop
shield rdp 10.0.0.5

# Expose a local dev server on port 3000
shield http 3000

# Expose a remote VNC server on a custom port
shield vnc 10.0.0.10:5901
```

Once the tunnel is established, open the **Access URL** in any browser — that's it.

**Step 1:** Run `shield ssh` to create the tunnel

![Shield CLI Terminal](docs/images/shieldcli-ssh-001.jpg)

**Step 2:** Open the Access URL — automatic authorization redirect

![Browser Authorization](docs/images/shieldcli-ssh-web-001.jpg)

**Step 3:** SSH terminal in your browser — no client needed

![Browser SSH Terminal](docs/images/shieldcli-ssh-web-002.jpg)

### Smart Defaults

| Command | Resolves To |
|---|---|
| `shield ssh` | `127.0.0.1:22` |
| `shield ssh 2222` | `127.0.0.1:2222` |
| `shield ssh 10.0.0.2` | `10.0.0.2:22` |
| `shield ssh 10.0.0.2:2222` | `10.0.0.2:2222` |
| `shield http` | `127.0.0.1:80` |
| `shield http 3000` | `127.0.0.1:3000` |
| `shield rdp` | `127.0.0.1:3389` |
| `shield vnc` | `127.0.0.1:5900` |
| `shield https` | `127.0.0.1:443` |
| `shield telnet` | `127.0.0.1:23` |

Supported protocols: `ssh`, `http`, `https`, `rdp`, `vnc`, `telnet`

## Visible Mode (default)

By default, the tunnel is in **visible mode** — anyone with the Access URL can connect directly. The Access URL is printed to the terminal after the tunnel is established.

```bash
shield ssh 10.0.0.2
shield rdp 10.0.0.5
```

You can filter a specific AC node by name:

```bash
shield --visable=HK ssh 10.0.0.2
```

**Use cases:** Development servers, demos, staging environments, team collaboration.

## Usage

```
shield <protocol> [ip:port] [flags]

Flags:
  -H, --server string         API server URL (default: https://console.yishield.com/raas)
  -p, --tunnel-port int       Chisel tunnel server port (default: 62888)
      --visable [filter]      AC node name filter (default: visible mode)
      --display-name string   Connector display name
      --site-name string      Application site name
      --username string       Target service username (SSH/RDP/VNC)
      --auth-pass string      Target service password (SSH/RDP/VNC)
      --private-key string    SSH private key
      --passphrase string     SSH private key passphrase
      --enable-sftp           Enable SFTP (SSH only)
  -v, --verbose               Enable verbose log output
  -h, --help                  Help for shield

Commands:
  clean                       Clear cached credentials
```

### Local API

Once running, Shield CLI exposes a local API on `127.0.0.1:<port>`:

| Endpoint | Method | Description |
|---|---|---|
| `/health` | GET | Health check |
| `/connectors` | GET | List all active tunnels |
| `/connector?rport=&lip=&lport=` | GET | Create a dynamic tunnel |
| `/connector?rport=` | DELETE | Close a tunnel |

## Architecture

```
                      Browser (RDP/VNC/SSH via HTML5)
                                │
                                ▼
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│  Internal     │ ◄──► │  Shield CLI   │ ◄══► │  Public      │
│  Service      │ local│  (tunnel)     │chisel│  Gateway     │
│  10.0.0.5:    │      │  127.0.0.1   │wss://│  + Web UI    │
│  3389/5900/22 │      └──────────────┘      └──────────────┘
└──────────────┘
```

## Security

- Credentials are encrypted with AES-256-GCM using a machine-specific fingerprint
- Passwords are masked in all log output
- Tunnel connections use authenticated WebSocket transport
- Credential files are stored with `0600` permissions

## Roadmap

### Core

- [x] Encrypted tunnel — secure WebSocket transport based on chisel
- [x] Multi-protocol support — SSH, HTTP, HTTPS, RDP, VNC, Telnet
- [x] Smart defaults — auto-detect IP and port per protocol, minimal input required
- [x] Cross-platform — native binaries for Linux, macOS, Windows (amd64/arm64)
- [x] Package manager distribution — Homebrew, Scoop, deb, rpm, curl/PowerShell one-liner
- [x] Auto credentials — machine fingerprint-based identity with AES-256-GCM encryption
- [x] Visible mode — control access authorization per tunnel
- [ ] Invisible mode — require Access URL with authorization key for secure access
- [x] Dynamic tunnels — runtime tunnel management via local REST API
- [x] Auto-reconnect — exponential backoff retry on connection failure
- [ ] Open-source server — self-hosted deployment, full control over data and infrastructure
- [ ] Persistent configuration — save tunnel profiles, reconnect with `shield up`
- [ ] Local Web UI — browser-based dashboard at `localhost` for managing tunnels, logs, and status
- [ ] Multi-tunnel mode — run multiple tunnels in a single process (`shield up` from config file)

### User Experience

- [x] Zero-config quick start — `shield ssh` just works, no flags required
- [x] Smart address parsing — supports `shield ssh`, `shield ssh 2222`, `shield ssh 10.0.0.2`, `shield ssh 10.0.0.2:2222`
- [x] Clean terminal UI — banner, tunnel mapping, and Access URL displayed clearly
- [x] Post-install usage hints — Homebrew caveats show examples after install
- [ ] Interactive setup wizard — guided first-run with `shield init`
- [ ] QR code output — scan to open Access URL on mobile devices
- [ ] Connection health monitor — real-time latency, bandwidth, and uptime stats in terminal
- [ ] Auto-reconnect with session resume — seamless recovery without generating new URLs
- [ ] Notification hooks — webhook / Slack / email alerts on tunnel connect/disconnect events

### Team & Enterprise

- [ ] Team workspace — shared tunnel dashboard, invite members, role-based access control
- [ ] Audit logs — who accessed what, when, with full session recording for compliance
- [ ] SSO integration — SAML / OIDC login for Access URL authorization
- [ ] Custom domains — use your own domain instead of `*.apps.yishield.com`
- [ ] IP allowlist — restrict Access URL to specific source IPs or CIDR ranges
- [ ] Tunnel expiration policies — auto-shutdown after N hours, scheduled access windows

### Value-Added Services

- [ ] File transfer — browser-based upload/download via tunneled services (SFTP, SCP)
- [ ] Session recording & playback — record RDP/VNC/SSH sessions for training and audit
- [ ] Multi-region relay — choose exit nodes in different regions for lower latency
- [ ] API gateway mode — expose REST/gRPC APIs with rate limiting, auth, and monitoring
- [ ] Mobile app — iOS/Android companion app for tunnel management and quick access

## License

Apache 2.0
