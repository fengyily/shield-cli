---
title: App Profiles — Save and Manage Remote Applications
description: Manage up to 10 app profiles in Shield CLI's Web UI with AES-256-GCM encrypted local storage. Configure protocol, IP, port, name, and credentials for each app.
head:
  - - meta
    - name: keywords
      content: Shield CLI app profiles, application management, encrypted storage, Web UI configuration
---

# App Profiles

In Web UI mode, Shield CLI supports saving and managing multiple app configurations.

## App Profile Fields

| Field | Description |
|---|---|
| Protocol | SSH / RDP / VNC / HTTP / HTTPS / Telnet |
| IP Address | Target service's internal IP |
| Port | Target service port |
| Name | Custom display name (e.g., "Office PC") |
| Credentials | Username, password (encrypted storage) |

## Limits

- Up to **10** saved app profiles
- Up to **3** concurrent connections

## Storage

App configurations are encrypted with AES-256-GCM and stored locally:

| Platform | Path |
|---|---|
| macOS / Linux | `~/.shield-cli/` |
| Windows | `%LOCALAPPDATA%\ShieldCLI\` |

Configurations are never uploaded to the server — fully local management.

## Management

In the Web UI you can:

- **Add** — Click the add button and fill in app details
- **Edit** — Modify saved app configurations
- **Delete** — Remove apps you no longer need
- **Rename** — Quickly change display names

Each app records its creation time, last modified time, and last connected time.
