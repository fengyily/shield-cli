---
title: MySQL Plugin — In-Browser Database Management
description: Shield CLI MySQL plugin provides a browser-based Web SQL Client with database browsing, schema viewing, SQL execution, result sorting, CSV export, and more.
head:
  - - meta
    - name: keywords
      content: Shield CLI MySQL, MySQL Web client, database management, Web SQL, MySQL browser, MariaDB, shield mysql
---

# MySQL Plugin

Access a full MySQL / MariaDB database management interface in your browser through Shield CLI.

## Installation

```bash
shield plugin add mysql
```

Verify installation:

```bash
shield plugin list
# NAME   VERSION  PROTOCOLS       INSTALLED
# mysql  v0.1.0   mysql, mariadb  2026-03-24T10:00:00+08:00
```

## Quick Connect

```bash
# Connect to local MySQL (default port 3306)
shield mysql

# Custom port
shield mysql 3307

# Remote server
shield mysql 10.0.0.5

# Full address
shield mysql 10.0.0.5:3307
```

## Authentication

### Interactive Input (Recommended)

When no credential flags are provided, you'll be prompted:

```bash
shield mysql 127.0.0.1:3306

  🔐 Database credentials (press Enter to skip)

  Username [root]: root
  Password: ****
  Database (optional): mydb

  ✓ Connecting as root
    Database: mydb
```

- **Username** defaults to `root` — press Enter to accept
- **Password** is hidden input
- **Database** is optional — press Enter to skip

### Command-line Flags

```bash
# Database-specific flags
shield mysql 127.0.0.1:3306 --db-user root --db-pass mypassword --db-name mydb

# Generic auth flags also work
shield mysql 127.0.0.1:3306 --username root --auth-pass mypassword
```

| Flag | Alias | Description |
|---|---|---|
| `--db-user` | `--username` | Database username |
| `--db-pass` | `--auth-pass` | Database password |
| `--db-name` | — | Database name (optional) |
| `--readonly` | — | Force read-only mode, block write operations |

## Web Management Interface

After connecting, the browser opens a Web SQL Client automatically:

### Feature Overview

| Feature | Description |
|---|---|
| Database browsing | Sidebar lists all databases, click to switch |
| Table list | Shows all tables after selecting a database, with filter and pagination |
| Schema viewer | Click a table to view columns, types, indexes |
| SQL execution | Supports SELECT, SHOW, DESCRIBE, and other queries |
| Result sorting | Click column headers to sort ascending/descending |
| CSV export | One-click export of query results to CSV |
| Copy results | Copy as tab-separated text, paste directly into Excel |
| Read-only mode | Enabled by default, blocks INSERT/UPDATE/DELETE |

### Read-Only Mode

Read-only / read-write mode is fully controlled by the startup parameters. The Web UI only displays the current state — remote users cannot change it.

**CLI mode**: use the `--readonly` flag:

```bash
# Read-only mode (recommended for sharing)
shield mysql 10.0.0.5:3306 --db-user root --readonly

# Read-write mode (default)
shield mysql 10.0.0.5:3306 --db-user root
```

**Web UI mode**: check the **Read-Only Mode** checkbox when adding or editing an application.

The top-right badge shows the current mode:

- **🔒 Read-Only** (orange badge) — write operations are blocked on both frontend and backend
- **🔓 Read-Write** (green badge) — all operations allowed

In read-only mode, these statements are blocked:

```
INSERT, UPDATE, DELETE, DROP, ALTER, CREATE,
TRUNCATE, RENAME, REPLACE, GRANT, REVOKE
```

::: tip Security Recommendation
When exposing a database management interface to the internet, enable read-only mode and use `--invisible` mode:
```bash
shield mysql 127.0.0.1:3306 --db-user readonly_user --readonly --invisible
```
:::

### Table Filter & Pagination

For databases with many tables:

- **Filter** — search box at the top of the table list supports real-time filtering
- **Pagination** — automatically paginated when there are more than 50 tables

### Keyboard Shortcuts

| Shortcut | Action |
|---|---|
| `Ctrl+Enter` / `Cmd+Enter` | Execute current SQL |
| `Tab` | Insert two spaces |

### Database Navigation

1. Sidebar shows all databases
2. Click a database name to switch to its table list
3. Click **← Databases** to go back to the database list
4. Single-click a table → view schema
5. Double-click a table → auto-fill `SELECT * FROM ... LIMIT 100` and execute

## Default Ports

| Input | Resolves To |
|---|---|
| `shield mysql` | `127.0.0.1:3306` |
| `shield mysql 3307` | `127.0.0.1:3307` |
| `shield mysql 10.0.0.5` | `10.0.0.5:3306` |
| `shield mysql 10.0.0.5:3307` | `10.0.0.5:3307` |

`mariadb` is an alias for `mysql` with identical behavior:

```bash
shield mariadb 127.0.0.1:3306 --db-user root
```

## Use Cases

### Remote Database Debugging

Connect to a remote MySQL locally and expose a Web management interface to team members:

```bash
shield mysql 10.0.0.5:3306 --db-user readonly --db-pass xxx --invisible
```

Share the Auth URL with colleagues who need to view data — no client installation required.

### MySQL in Docker

Connect directly to MySQL running in a Docker container:

```bash
# MySQL container listening on 3306
docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:8

# Expose via Shield
shield mysql 127.0.0.1:3306 --db-user root --db-pass root
```

### Database Audit

Connect in read-only mode for auditors to review data via browser:

```bash
shield mysql 192.168.1.100:3306 --db-user auditor --db-pass xxx --invisible
```

- Read-only mode by default
- Invisible mode requires authorization code
- Export results to CSV for reports

### Temporary Sharing with External Partners

```bash
shield mysql 127.0.0.1:3306 --db-user report_user --db-pass xxx
```

Partners only need a browser — no MySQL client, VPN, or firewall rules needed.

## API Endpoints

The MySQL plugin runs a local HTTP service with these APIs:

| Method | Path | Description |
|---|---|---|
| GET | `/api/info` | Server info (version, host, user) |
| GET | `/api/databases` | List databases |
| GET | `/api/tables?db=mydb` | List tables |
| GET | `/api/schema?db=mydb&table=users` | Table schema (columns, types, indexes) |
| POST | `/api/query` | Execute SQL query |

### Query API Example

```bash
curl -X POST http://localhost:19876/api/query \
  -H 'Content-Type: application/json' \
  -d '{"sql": "SELECT * FROM users LIMIT 10", "db": "mydb"}'
```

Response:

```json
{
  "code": 200,
  "data": {
    "columns": ["id", "name", "email"],
    "rows": [
      {"id": 1, "name": "Alice", "email": "alice@example.com"}
    ],
    "count": 1,
    "duration": "2.3ms"
  }
}
```

## Troubleshooting

### Connection Refused

```
plugin error: cannot connect to MySQL at 127.0.0.1:3306: connection refused
```

Check if MySQL is running:

```bash
# macOS
brew services list | grep mysql

# Linux
systemctl status mysql

# Docker
docker ps | grep mysql
```

### Authentication Failed

```
plugin error: Access denied for user 'root'@'172.17.0.1'
```

- Verify username and password
- In Docker, check MySQL `bind-address` and user permissions
- Try creating a user with remote access:
  ```sql
  CREATE USER 'shield'@'%' IDENTIFIED BY 'password';
  GRANT SELECT ON *.* TO 'shield'@'%';
  ```

### Plugin Not Installed

```
unsupported protocol "mysql"
```

Install the plugin first:

```bash
shield plugin add mysql
```

## Next Steps

- [Plugin System Overview](/en/plugins/)
- [Plugin Development Guide](/en/plugins/development)
- [TCP Port Proxy](/en/protocols/tcp-udp) (alternative: proxy port 3306 directly)
