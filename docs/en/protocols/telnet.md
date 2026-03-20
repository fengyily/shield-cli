---
title: Telnet Tunnel — Network Device Access in Browser
description: Connect to Telnet services in your browser via Shield CLI. Ideal for managing routers, switches, network devices, and legacy systems.
head:
  - - meta
    - name: keywords
      content: Telnet tunnel, Telnet browser, network device management, legacy system access, Shield CLI Telnet
---

# Telnet

Connect to Telnet services in your browser via Shield CLI — ideal for network devices and legacy systems.

## Quick Connect

```bash
# Connect to localhost
shield telnet

# Connect to specific device
shield telnet 10.0.0.1

# Specify port
shield telnet 10.0.0.1:2323
```

## Use Cases

- Manage routers, switches, and other network devices
- Connect to legacy industrial control systems
- Access servers that don't support SSH

## Default Ports

| Input | Resolves To |
|---|---|
| `shield telnet` | `127.0.0.1:23` |
| `shield telnet 2323` | `127.0.0.1:2323` |
| `shield telnet 10.0.0.1` | `10.0.0.1:23` |
