---
title: Plugin System — Extend Protocols On Demand
description: Shield CLI plugin system overview. Extend database and service protocol support with independent plugins including MySQL, PostgreSQL, SQL Server, and more.
head:
  - - meta
    - name: keywords
      content: Shield CLI plugins, plugin system, shield plugin, database management, MySQL plugin, PostgreSQL plugin, extend protocols
---

# Plugin System

Shield CLI extends protocol support through plugins. Each plugin is an independent binary, installed on demand — zero bloat to the main program.

## Design Philosophy

Shield CLI has built-in support for SSH, RDP, VNC, HTTP, and other common protocols. For databases, message queues, and other services that need a dedicated Web management interface, we use plugins:

- **Install on demand** — unused protocols don't increase the main binary size
- **Independent updates** — plugin versions are independent of the main program
- **Open extension** — anyone can develop and publish Shield plugins

## How It Works

```
shield mysql 127.0.0.1:3306 --db-user root
    ↓
Shield looks up the installed mysql plugin
    ↓
Starts the plugin process (independent binary)
    ↓
Plugin launches a local Web database management UI
    ↓
Shield exposes the Web UI to the internet via encrypted tunnel
    ↓
User accesses the full database management platform in the browser
```

From the Shield server's perspective, the plugin's Web UI is just a regular HTTP application — no server-side adaptation needed.

## Available Plugins

| Plugin | Protocols | Default Port | Description |
|---|---|---|---|
| [mysql](/en/plugins/mysql) | `mysql`, `mariadb` | 3306 | MySQL / MariaDB Web management client |
| postgres | `postgres`, `pg` | 5432 | PostgreSQL Web management client (coming soon) |
| sqlserver | `sqlserver`, `mssql` | 1433 | SQL Server Web management client (coming soon) |

## Quick Start

```bash
# Install plugin
shield plugin add mysql

# Use (interactive credential input)
shield mysql 127.0.0.1:3306

# Use (pass credentials via flags)
shield mysql 127.0.0.1:3306 --db-user root --db-pass mypassword --db-name mydb
```

After connecting, the browser automatically opens the Web database management interface.

## Managing Plugins

```bash
# List installed plugins
shield plugin list

# Install a plugin
shield plugin add <name>

# Install from a local binary (for development)
shield plugin add <name> --from ./path/to/binary

# Remove a plugin
shield plugin remove <name>
```

## Plugin Communication Protocol

Shield communicates with plugins via **stdin/stdout JSON** — the protocol is extremely simple:

### Start Request (Shield to Plugin)

```json
{
  "action": "start",
  "config": {
    "host": "127.0.0.1",
    "port": 3306,
    "user": "root",
    "pass": "password",
    "database": "mydb"
  }
}
```

### Ready Response (Plugin to Shield)

```json
{
  "status": "ready",
  "web_port": 19876,
  "name": "MySQL Web Client",
  "version": "0.1.0"
}
```

### Stop Request (Shield to Plugin)

```json
{"action": "stop"}
```

## Plugin Storage Location

| Platform | Path |
|---|---|
| macOS / Linux | `~/.shield-cli/plugins/` |
| Windows | `%LOCALAPPDATA%\ShieldCLI\plugins\` |

```
~/.shield-cli/plugins/
├── registry.json              # Installed plugins index
├── shield-plugin-mysql        # MySQL plugin binary
└── shield-plugin-postgres     # PostgreSQL plugin binary
```

## Developing Custom Plugins

Want to develop a new plugin for Shield? See the [Plugin Development Guide](/en/plugins/development).

## Next Steps

- [MySQL Plugin Documentation](/en/plugins/mysql)
- [Plugin Development Guide](/en/plugins/development)
- [Command Reference](/en/reference/commands)
