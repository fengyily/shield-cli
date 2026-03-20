---
title: Quick Start — Shield CLI in 5 Minutes
description: Get started with Shield CLI in 5 minutes. Learn to create SSH, RDP, and HTTP tunnels using Web UI mode or CLI mode, with smart address resolution.
head:
  - - meta
    - name: keywords
      content: Shield CLI quickstart, getting started, SSH tunnel, RDP tunnel, HTTP tunnel, tutorial
---

# Quick Start

This tutorial will take you from installation to accessing your first internal service.

## Option 1: Web UI Mode (Recommended)

### Step 1: Launch the Dashboard

```bash
shield start
```

Your browser will automatically open `http://localhost:8181`.

### Step 2: Add an App

In the Web UI, click the add button and fill in:

- **Protocol**: Choose SSH / RDP / VNC / HTTP, etc.
- **IP Address**: The target service's internal IP (e.g., `10.0.0.5`)
- **Port**: The target service port (e.g., `22` for SSH)
- **Name**: A display name (e.g., "Office PC")

### Step 3: Connect

Click the **Connect** button on the app card. Once the status changes to "Connected", your browser will automatically open the Access URL.

You can now operate the remote service directly in your browser.

---

## Option 2: CLI Mode

### Connect to Local SSH

```bash
shield ssh
```

Shield CLI will:
1. Resolve to `127.0.0.1:22`
2. Prompt for username and password
3. Establish an encrypted tunnel
4. Output the Access URL

Open the URL in a browser to see the SSH terminal.

### Connect to a Remote Windows Desktop

```bash
shield rdp 10.0.0.5
```

Enter Windows login credentials, then you'll see the full remote desktop in your browser.

### Expose a Local Web App

```bash
shield http 3000
```

Your web app running on `localhost:3000` is now accessible via a public URL.

## Smart Address Resolution

Shield CLI supports flexible address input:

| Input | Resolves To |
|---|---|
| `shield ssh` | `127.0.0.1:22` |
| `shield ssh 2222` | `127.0.0.1:2222` |
| `shield ssh 10.0.0.2` | `10.0.0.2:22` |
| `shield ssh 10.0.0.2:2222` | `10.0.0.2:2222` |
| `shield rdp` | `127.0.0.1:3389` |
| `shield vnc` | `127.0.0.1:5900` |
| `shield http` | `127.0.0.1:80` |
| `shield http 3000` | `127.0.0.1:3000` |

Each protocol has a preset default port — you only need to specify what's different.

## Next Steps

- [Web UI Mode](./web-ui.md)
- [CLI Mode](./cli-mode.md)
- [Protocol Guide](../protocols/ssh.md)
