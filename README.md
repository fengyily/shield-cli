<p align="center">
  <img src="docs/logo.svg" alt="Shield CLI" width="128" height="128">
</p>

<h1 align="center">Shield CLI</h1>

<p align="center">
  <strong>Access any internal service from your browser. No VPN, no client, one command.</strong><br>
  Shield CLI is a browser-first internal service gateway — SSH terminals, remote desktops, database admin, web apps — all accessible through any browser with a single command.
</p>

<p align="center">
  <a href="https://docs.yishield.com/en/guide/what-is-shield">Documentation</a> &bull;
  <a href="https://docs.yishield.com/en/guide/install">Installation</a> &bull;
  <a href="https://docs.yishield.com/en/guide/quickstart">Quick Start</a> &bull;
  <a href="README_CN.md">中文文档</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-blue?logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-brightgreen" alt="Platform">
  <img src="https://img.shields.io/badge/license-Apache%202.0-green" alt="License">
</p>

---

## How It Works

<p align="center">
  <img src="docs/demo/architecture.gif" alt="Shield CLI Architecture" width="800">
</p>

---

## Demo

### RDP — Browser Remote Desktop

<p align="center">
  <img src="docs/demo/demo-rdp.gif" alt="Shield CLI RDP Demo" width="960">
</p>

### SSH — Browser Terminal

<p align="center">
  <img src="docs/demo/demo-ssh.gif" alt="Shield CLI SSH Demo" width="960">
</p>

---

## Why Shield CLI?

Traditional tools solve **network reachability** (ngrok, frp) or **access control** (Teleport, Boundary) — but they still require protocol-specific clients or complex setup.

Shield CLI is a **unified browser gateway** for all your internal services. One binary, one command — SSH terminals, remote desktops, database admin, web apps — all rendered in the browser via HTML5.

| Capability | Shield CLI | ngrok/frp | Teleport/Boundary |
|-----------|-----------|-----------|-------------------|
| Browser RDP/VNC/SSH | Yes | No | Partial |
| Database Web Admin | Yes (plugins) | No | No |
| Zero client install | Yes | No | No |
| Single binary deploy | Yes | Yes | No |
| Plugin extensibility | Yes | No | No |

## Installation

```bash
# macOS
brew tap fengyily/tap && brew install shield-cli

# Windows
scoop bucket add shield https://github.com/fengyily/scoop-bucket && scoop install shield-cli

# Linux (apt) — Debian / Ubuntu
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/scripts/setup-repo.sh | sudo bash

# Linux (yum) — RHEL / CentOS / Fedora
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/scripts/setup-repo.sh | sudo bash

# Linux / macOS (one-liner binary)
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh

# China mirror (jsDelivr CDN)
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

### Docker

```bash
# Use the prebuilt image (recommended)
docker run -d --name shield \
  --network host \
  --restart unless-stopped \
  fengyily/shield-cli

# Or build from source
docker build -t shield-cli .
docker run -d --name shield --network host --restart unless-stopped shield-cli
```

> **Note:** `--network host` shares the host's network stack, allowing Shield CLI to reach local and LAN services (e.g., `10.0.0.x`, `192.168.x.x`). Open `http://localhost:8181` to access the Web UI.
>
> **Caveat:** `--network host` only works on **Linux**. On macOS/Windows Docker Desktop, use port mapping instead:
>
> ```bash
> docker run -d --name shield -p 8181:8181 --restart unless-stopped fengyily/shield-cli
> ```

More installation methods (apt, yum, deb, rpm, PowerShell, source build): [Installation Guide](https://docs.yishield.com/en/guide/install)

## Quick Start

### Web UI (Recommended)

```bash
shield start
```

Open `http://localhost:8181`, add your services, and connect with one click. On macOS and Windows, a system tray icon provides quick access to the Dashboard.

![Web UI Dashboard](docs/images/shieldcli-webui-001.jpg)

![RDP via Web UI](docs/images/shieldcli-rdp-web-001.jpg)

### System Service (Auto-Start on Boot)

```bash
shield install              # Install as system service (port 8181)
shield install --port 8182  # Use custom port if 8181 is occupied
shield start                # Start the service (if stopped)
shield stop                 # Stop the service
shield uninstall            # Remove the service
```

After `shield install`, the service starts automatically and will restart on boot. If the service is stopped, use `shield start` to restart it — no need to reinstall.

Supports macOS (launchd), Linux (systemd), and Windows. See [System Service Guide](https://docs.yishield.com/en/guide/system-service) for details.

### Command Line

```bash
shield ssh              # SSH terminal in browser (127.0.0.1:22)
shield rdp 10.0.0.5     # Windows desktop in browser
shield mysql 10.0.0.20  # Database admin in browser (plugin)
shield http 3000        # Expose local web app
shield vnc 10.0.0.10    # VNC screen sharing in browser
shield tcp 3306         # TCP port proxy
shield udp 53           # UDP port proxy
```

![Shield CLI Terminal](docs/images/shieldcli-ssh-001.jpg)

![Browser SSH Terminal](docs/images/shieldcli-ssh-web-002.jpg)

### Smart Defaults

| Command | Resolves To |
|---------|-------------|
| `shield ssh` | `127.0.0.1:22` |
| `shield ssh 2222` | `127.0.0.1:2222` |
| `shield ssh 10.0.0.2` | `10.0.0.2:22` |
| `shield rdp` | `127.0.0.1:3389` |
| `shield http 3000` | `127.0.0.1:3000` |
| `shield tcp 3306` | `127.0.0.1:3306` |
| `shield udp 53` | `127.0.0.1:53` |

Protocols: `ssh`, `rdp`, `vnc`, `http`, `https`, `telnet`, `tcp`, `udp` — [Full Commands Reference](https://docs.yishield.com/en/reference/commands)

## Security

- **AES-256-GCM encryption** — credentials encrypted with machine fingerprint-derived keys
- **Password masking** — all passwords hidden in logs
- **WebSocket transport** — authenticated encrypted tunnels
- **0600 permissions** — credential files readable only by owner

Details: [Credentials](https://docs.yishield.com/en/security/credentials) | [Access Modes](https://docs.yishield.com/en/security/access-modes)

## Documentation

Full documentation is available at **[docs.yishield.com](https://docs.yishield.com)**:

- [What is Shield CLI](https://docs.yishield.com/en/guide/what-is-shield) — overview and key features
- [Installation](https://docs.yishield.com/en/guide/install) — all installation methods
- [Quick Start](https://docs.yishield.com/en/guide/quickstart) — 5-minute tutorial
- [Protocol Guides](https://docs.yishield.com/en/protocols/ssh) — SSH, RDP, VNC, HTTP, Telnet
- [Plugin System](https://docs.yishield.com/en/plugins/) — MySQL and more
- [Commands Reference](https://docs.yishield.com/en/reference/commands) — full parameter guide
- [FAQ](https://docs.yishield.com/en/reference/faq) — frequently asked questions
- [Troubleshooting](https://docs.yishield.com/en/troubleshooting/errors) — common errors and fixes

## License

Apache 2.0
