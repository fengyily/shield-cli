# Shield CLI MCP Server

An [MCP (Model Context Protocol)](https://modelcontextprotocol.io) server that provides Shield CLI documentation to AI tools like Claude Code, Cursor, Trae, Windsurf, and other MCP-compatible clients.

## What It Provides

| Tool | Description |
|------|-------------|
| `shield_about` | Overview, architecture, and comparison with alternatives |
| `shield_install` | Installation instructions by platform (macOS/Windows/Linux/Docker) |
| `shield_usage` | Quick start, commands, address resolution, Web UI, system service |
| `shield_protocols` | Supported protocols (SSH, RDP, VNC, HTTP, Telnet, TCP, UDP) |
| `shield_plugins` | Plugin info (MySQL, PostgreSQL, SQL Server) |

## Setup

### Claude Code

```bash
claude mcp add shield-cli -- npx -y shield-cli-mcp
```

Or add to `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "shield-cli": {
      "command": "npx",
      "args": ["-y", "shield-cli-mcp"]
    }
  }
}
```

### Cursor

Add to Cursor Settings → MCP Servers:

```json
{
  "mcpServers": {
    "shield-cli": {
      "command": "npx",
      "args": ["-y", "shield-cli-mcp"]
    }
  }
}
```

### Trae

Click the AI sidebar → Settings icon → MCP → Add MCP Server, paste:

```json
{
  "mcpServers": {
    "shield-cli": {
      "command": "npx",
      "args": ["-y", "shield-cli-mcp"]
    }
  }
}
```

### Windsurf / Other MCP Clients

Use the same configuration format — command: `npx`, args: `["-y", "shield-cli-mcp"]`.

## Uninstall

**Claude Code:** `claude mcp remove shield-cli`

**Cursor / Trae / Windsurf:** Remove the `"shield-cli"` entry from your MCP settings JSON.

## Local Development

```bash
cd mcp-server
npm install
npm run build

# Test with Claude Code (local)
claude mcp add shield-cli-dev -- node /path/to/mcp-server/dist/index.js
```

## License

Apache 2.0
