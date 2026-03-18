<p align="center">
  <img src="docs/logo.svg" alt="Shield CLI" width="128" height="128">
</p>

<h1 align="center">Shield CLI</h1>

<p align="center">
  <strong>Browser-Based Secure Tunnel for RDP, VNC, SSH & More</strong><br>
  Access internal RDP desktops, VNC sessions, SSH terminals, and web services directly from a browser — no client software required.
</p>

<p align="center">
  <a href="#installation">Installation</a> &bull;
  <a href="#quick-start">Quick Start</a> &bull;
  <a href="#visibility-modes">Visibility Modes</a> &bull;
  <a href="#usage">Usage</a> &bull;
  <a href="README_CN.md">中文文档</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-blue?logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-brightgreen" alt="Platform">
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
</p>

---

## Why Shield CLI?

Traditional remote access requires installing dedicated clients (RDP client, VNC viewer, SSH terminal) on every device. Shield CLI eliminates this by tunneling internal services through a secure gateway that renders everything **in the browser**.

- **No client installation** — Open a URL in any browser to access RDP desktops, VNC screens, or SSH terminals
- **Works anywhere** — Access from phones, tablets, locked-down corporate machines, or any device with a browser
- **One binary** — A single `shield` command exposes any internal service to the web

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

## Visibility Modes

Shield CLI supports two access modes. When a tunnel is established, two URLs are generated:

- **Site URL** — The application address (e.g., `https://xxxx.hk01.apps.yishield.com`). This URL alone is **not accessible** without authorization.
- **Access URL** — Contains an embedded authorization key. Anyone with this URL can access the service directly.

### Invisible Mode (default)

The service **requires authorization** — the Site URL alone will not grant access. Users must use the Access URL.

```bash
shield rdp 10.0.0.5
shield ssh 10.0.0.2
```

**Use cases:** Production servers, incident response, sensitive machines.

### Visible Mode

The service is **open to unauthorized users** — anyone who knows the Site URL can access it.

```bash
shield --visable ssh 10.0.0.2
shield --visable=HK rdp 10.0.0.5
```

**Use cases:** Public demos, shared dev servers, QA staging environments.

## Usage

```
shield <protocol> [ip:port] [flags]

Flags:
  -H, --server string         API server URL (default: https://console.yishield.com/raas)
  -p, --tunnel-port int       Chisel tunnel server port (default: 62888)
      --visable [filter]      Enable visible mode (optional: AC node name filter)
      --invisible             Invisible mode with authorization key
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

### Example Output

<p align="center">
  <img src="docs/images/shieldcli-ssh-001.jpg" alt="Shield CLI SSH" width="600">
</p>

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

## License

MIT
