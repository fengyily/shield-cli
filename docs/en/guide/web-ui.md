---
title: Web UI Mode — Browser-Based Management Dashboard
description: Use Shield CLI's Web UI mode to manage up to 10 remote apps with one-click connect/disconnect, real-time status monitoring, and encrypted local storage.
head:
  - - meta
    - name: keywords
      content: Shield CLI Web UI, management dashboard, remote app management, one-click connect
---

# Web UI Mode

Web UI is the recommended way to use Shield CLI, providing a browser-based graphical management dashboard.

## Launch

```bash
shield start
```

Opens at `http://localhost:8181` by default. Specify a different port if needed:

```bash
shield start 9090
```

## Features

### App Management

- **Add**: Configure protocol, IP, port, and display name
- **Edit**: Modify saved app configurations
- **Delete**: Remove apps you no longer need
- **Rename**: Quickly change display names
- **Up to 10 apps**: Encrypted local storage

### Connection Control

- **One-click connect**: Click the connect button on any app card
- **One-click disconnect**: Click disconnect to close the tunnel
- **Up to 3 concurrent connections**: Connect to multiple apps simultaneously
- **Auto-open browser**: Automatically opens the Access URL in a new tab on success

### Status Monitoring

App cards display real-time connection status:

| Status | Description |
|---|---|
| Disconnected | App saved but no tunnel established |
| Connecting | Establishing tunnel and activating service |
| Connected | Tunnel active, accessible via Access URL |
| Failed | Tunnel establishment failed, check error message |

Status refreshes automatically every 2 seconds.

### Theme Toggle

Supports dark and light themes. Toggle via the button in the top right corner.

## Best For

- Day-to-day work managing multiple remote apps
- Non-technical users who prefer a graphical interface
- Frequently connecting/disconnecting different services
- Saving commonly used app configurations

## Next Steps

- [CLI Mode](./cli-mode.md) — for servers and scripting
- [App Profiles](../config/apps.md)
