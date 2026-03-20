---
title: RDP Tunnel — Windows Remote Desktop in Browser
description: Access Windows Remote Desktop (RDP) directly in your browser via Shield CLI. Full mouse and keyboard control, no RDP client needed. Works from any device.
head:
  - - meta
    - name: keywords
      content: RDP tunnel, RDP browser, Windows remote desktop, remote desktop browser, Shield CLI RDP
---

# RDP

Access Windows Remote Desktop in your browser via Shield CLI — no RDP client needed.

## Quick Connect

```bash
# Connect to local Windows
shield rdp

# Connect to specific IP
shield rdp 10.0.0.5

# Specify port
shield rdp 10.0.0.5:3390
```

## Authentication

```bash
shield rdp 10.0.0.5 --username Administrator --auth-pass mypassword
```

Omit flags for interactive prompts.

## Browser Experience

Once connected, you get a full Windows desktop in the browser:

- Full mouse and keyboard control
- Screen adapts to browser window size
- Accessible from any device — phones, tablets, laptops

## Use Cases

- Remotely operate an office Windows PC
- Access internally deployed Windows servers
- Let off-site colleagues temporarily use a Windows machine

## Default Ports

| Input | Resolves To |
|---|---|
| `shield rdp` | `127.0.0.1:3389` |
| `shield rdp 3390` | `127.0.0.1:3390` |
| `shield rdp 10.0.0.5` | `10.0.0.5:3389` |
