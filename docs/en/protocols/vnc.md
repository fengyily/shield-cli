---
title: VNC Tunnel — Remote Desktop Sharing in Browser
description: Share and control remote desktop screens in your browser via Shield CLI VNC tunnel. Pixel-perfect rendering with full mouse and keyboard support.
head:
  - - meta
    - name: keywords
      content: VNC tunnel, VNC browser, remote desktop sharing, screen sharing, Shield CLI VNC
---

# VNC

Share and control remote desktop screens in your browser via Shield CLI.

## Quick Connect

```bash
# Connect to localhost
shield vnc

# Connect to specific IP
shield vnc 10.0.0.10

# Specify port
shield vnc 10.0.0.10:5901
```

## Authentication

```bash
shield vnc 10.0.0.10 --auth-pass vncpassword
```

VNC typically only requires a password, no username.

## Browser Experience

- Pixel-perfect remote desktop rendering
- Full mouse and keyboard mapping
- Works with VNC servers on Linux, macOS, and Windows

## Use Cases

- Remote assistance: share your screen with colleagues or clients
- Manage headless Linux servers with a GUI
- Remotely operate lab or factory equipment

## Default Ports

| Input | Resolves To |
|---|---|
| `shield vnc` | `127.0.0.1:5900` |
| `shield vnc 5901` | `127.0.0.1:5901` |
| `shield vnc 10.0.0.10` | `10.0.0.10:5900` |
