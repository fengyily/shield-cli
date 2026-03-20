---
title: What is Shield CLI — Secure Tunnel Connector
description: Shield CLI is a secure tunnel connector that exposes internal network services (SSH, RDP, VNC, HTTP) to the public internet, accessible through any web browser with a single command. No client installation needed.
head:
  - - meta
    - name: keywords
      content: Shield CLI, what is Shield CLI, secure tunnel, remote access, browser terminal, SSH browser, RDP browser
---

# What is Shield CLI

Shield CLI is a **secure tunnel connector** that exposes internal network services to the public internet, making them accessible through any web browser with a single command.

## The Problem It Solves

You've probably faced these situations:

- Need to remotely access a server on the company network, but setting up a VPN is too complex
- Want to let a colleague temporarily access your local development environment
- Need to operate a Windows desktop remotely from a phone or tablet
- A client needs access to an internally deployed web app for testing

Traditional solutions (VPN, port forwarding, tunneling tools) are either complex to configure or require the other party to install client software.

**The Shield CLI way:**

```bash
shield ssh 10.0.0.5
```

You get a public URL. Open it in any browser to directly operate the SSH terminal. No client installation, no network configuration needed.

## Key Features

| Feature | Description |
|---|---|
| **Browser Access** | RDP, VNC, SSH, web apps — all rendered via HTML5 in the browser |
| **Zero Client** | Visitors only need a browser — phones, tablets, locked-down PCs all work |
| **Encrypted** | WebSocket encrypted tunnel + AES-256-GCM local credential encryption |
| **Smart Defaults** | `shield ssh` resolves to `127.0.0.1:22` automatically |
| **Dual Mode** | Web UI dashboard (recommended) + pure CLI mode |
| **Six Protocols** | SSH, RDP, VNC, HTTP, HTTPS, Telnet |
| **Cross-Platform** | macOS / Linux / Windows, amd64 and arm64 |

## How It Works

```
Your Service ←→ Shield CLI ←→ Public Gateway ←→ Browser
(SSH/RDP/...)    (Encrypted)    (HTML5 Render)   (Any Device)
```

1. Shield CLI runs on your machine and connects to the internal service
2. Establishes an encrypted WebSocket tunnel to the public gateway
3. The gateway assigns a unique Access URL
4. Visitors open the URL in a browser to operate the remote service

## Two Ways to Use

### Web UI Mode (Recommended)

Launch a local management dashboard and manage all apps in the browser:

```bash
shield start
```

Open `http://localhost:8181` to add apps and connect with one click.

### CLI Mode

Create tunnels directly from the terminal — ideal for servers or scripting:

```bash
shield ssh           # Connect to local SSH
shield rdp 10.0.0.5  # Connect to remote Windows desktop
shield http 3000     # Expose local web app
```

## Next Steps

- [Installation](./install.md)
- [Quick Start](./quickstart.md)
