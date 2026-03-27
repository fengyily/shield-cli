#!/usr/bin/env node

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";

const server = new McpServer({
  name: "shield-cli",
  version: "0.1.0",
});

// ─── Tool: get installation guide ───────────────────────────────────────────

server.tool(
  "shield_install",
  "Get Shield CLI installation instructions for a specific platform",
  {
    platform: z
      .enum(["macos", "windows", "linux", "docker", "all"])
      .default("all")
      .describe("Target platform"),
  },
  async ({ platform }) => {
    const guides: Record<string, string> = {
      macos: `## macOS

**Homebrew (recommended):**
\`\`\`bash
brew tap fengyily/tap && brew install shield-cli
\`\`\`

**One-liner:**
\`\`\`bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
\`\`\`

**China mirror:**
\`\`\`bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
\`\`\``,

      windows: `## Windows

**Scoop (recommended):**
\`\`\`powershell
scoop bucket add shield https://github.com/fengyily/scoop-bucket
scoop install shield-cli
\`\`\`

**PowerShell one-liner:**
\`\`\`powershell
irm https://raw.githubusercontent.com/fengyily/shield-cli/main/install.ps1 | iex
\`\`\``,

      linux: `## Linux

**Auto-detect (apt/yum/dnf):**
\`\`\`bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/scripts/setup-repo.sh | sudo bash
\`\`\`

**One-liner binary:**
\`\`\`bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
\`\`\`

**China mirror:**
\`\`\`bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
\`\`\``,

      docker: `## Docker

\`\`\`bash
# Linux (recommended — uses host network, can access LAN services)
docker run -d --name shield --network host --restart unless-stopped fengyily/shield-cli
\`\`\`

Open http://localhost:8181 to access the Web UI.

> **Note:** \`--network host\` only works on **Linux**. On macOS/Windows, Docker containers cannot directly access the host network, so Shield CLI inside the container cannot reach host or LAN services (e.g., \`127.0.0.1:22\`, \`10.0.0.x\`). For macOS/Windows, install Shield CLI directly via Homebrew or Scoop instead of Docker.`,
    };

    if (platform === "all") {
      return {
        content: [
          {
            type: "text" as const,
            text: `# Shield CLI Installation Guide\n\n${Object.values(guides).join("\n\n---\n\n")}\n\n---\n\nVerify: \`shield --version\`\n\nFull guide: https://docs.yishield.com/en/guide/install`,
          },
        ],
      };
    }

    return {
      content: [
        {
          type: "text" as const,
          text: `# Shield CLI Installation — ${platform}\n\n${guides[platform]}\n\nVerify: \`shield --version\`\n\nFull guide: https://docs.yishield.com/en/guide/install`,
        },
      ],
    };
  }
);

// ─── Tool: usage guide & commands ───────────────────────────────────────────

server.tool(
  "shield_usage",
  "Get Shield CLI usage examples, commands, and smart address resolution rules",
  {
    topic: z
      .enum(["quickstart", "commands", "address", "web-ui", "system-service"])
      .default("quickstart")
      .describe("Usage topic"),
  },
  async ({ topic }) => {
    const topics: Record<string, string> = {
      quickstart: `# Quick Start

## Web UI Mode (recommended)
\`\`\`bash
shield start
\`\`\`
Open http://localhost:8181, add your services, connect with one click.

## CLI Mode
\`\`\`bash
shield ssh              # SSH terminal in browser (127.0.0.1:22)
shield rdp 10.0.0.5     # Windows desktop in browser
shield vnc 10.0.0.10    # VNC screen sharing in browser
shield http 3000         # Expose local web app
shield mysql 10.0.0.20  # Database admin in browser (plugin)
shield postgres 10.0.0.30  # PostgreSQL admin in browser (plugin)
shield tcp 3306          # TCP port proxy
shield udp 53            # UDP port proxy
\`\`\`

Full guide: https://docs.yishield.com/en/guide/quickstart`,

      commands: `# Command Reference

| Command | Description |
|---------|-------------|
| \`shield start\` | Launch Web UI dashboard (port 8181) |
| \`shield ssh [addr]\` | SSH terminal in browser |
| \`shield rdp [addr]\` | Windows Remote Desktop in browser |
| \`shield vnc [addr]\` | VNC screen sharing in browser |
| \`shield http [port/addr]\` | Expose HTTP web app |
| \`shield https [port/addr]\` | Expose HTTPS web app |
| \`shield telnet [addr]\` | Telnet terminal in browser |
| \`shield tcp [port/addr]\` | TCP port proxy |
| \`shield udp [port/addr]\` | UDP port proxy |
| \`shield install\` | Install as system service |
| \`shield uninstall\` | Remove system service |
| \`shield stop\` | Stop the service |
| \`shield clean\` | Clear cached credentials |
| \`shield plugin install <name>\` | Install a plugin |
| \`shield plugin list\` | List installed plugins |
| \`shield plugin uninstall <name>\` | Remove a plugin |

Options: \`--username\`, \`--auth-pass\`, \`--port\`, \`--invisible\`

Full reference: https://docs.yishield.com/en/reference/commands`,

      address: `# Smart Address Resolution

Shield CLI intelligently resolves addresses to minimize typing:

| Input | Resolves To | Rule |
|-------|-------------|------|
| \`shield ssh\` | 127.0.0.1:22 | empty = localhost + default port |
| \`shield ssh 2222\` | 127.0.0.1:2222 | number only = port on localhost |
| \`shield ssh 10.0.0.5\` | 10.0.0.5:22 | IP only = IP + default port |
| \`shield ssh 10.0.0.5:2222\` | 10.0.0.5:2222 | full = as specified |

Default ports: SSH(22), RDP(3389), VNC(5900), HTTP(80), HTTPS(443), Telnet(23), MySQL(3306), Postgres(5432), SQLServer(1433)`,

      "web-ui": `# Web UI Mode

\`\`\`bash
shield start              # Launch at http://localhost:8181
shield start --port 8182  # Custom port
\`\`\`

Features:
- Manage up to 10 saved app profiles
- One-click connect/disconnect
- Up to 3 concurrent connections
- System tray icon on macOS and Windows
- Real-time connection status

Guide: https://docs.yishield.com/en/guide/web-ui`,

      "system-service": `# System Service

\`\`\`bash
shield install              # Install as system service (port 8181)
shield install --port 8182  # Custom port
shield start                # Start service (if stopped)
shield stop                 # Stop service
shield uninstall            # Remove service
\`\`\`

Supports: macOS (launchd), Linux (systemd), Windows.
After install, service auto-starts on boot.

Guide: https://docs.yishield.com/en/guide/system-service`,
    };

    return {
      content: [{ type: "text" as const, text: topics[topic] }],
    };
  }
);

// ─── Tool: protocol information ─────────────────────────────────────────────

server.tool(
  "shield_protocols",
  "Get information about Shield CLI supported protocols and their browser experience",
  {
    protocol: z
      .enum(["ssh", "rdp", "vnc", "http", "telnet", "tcp", "udp", "all"])
      .default("all")
      .describe("Protocol to query"),
  },
  async ({ protocol }) => {
    const protocols: Record<string, string> = {
      ssh: `## SSH
- Default port: 22
- Browser experience: Full xterm.js terminal + optional SFTP file transfer
- Example: \`shield ssh 10.0.0.5 --username root --auth-pass mypass\`
- Guide: https://docs.yishield.com/en/protocols/ssh`,

      rdp: `## RDP (Remote Desktop)
- Default port: 3389
- Browser experience: Complete Windows desktop with mouse/keyboard control
- Example: \`shield rdp 10.0.0.5 --username Administrator --auth-pass mypass\`
- Guide: https://docs.yishield.com/en/protocols/rdp`,

      vnc: `## VNC
- Default port: 5900
- Browser experience: Pixel-perfect remote desktop sharing
- Example: \`shield vnc 10.0.0.10 --auth-pass mypass\`
- Guide: https://docs.yishield.com/en/protocols/vnc`,

      http: `## HTTP / HTTPS
- Default ports: HTTP(80), HTTPS(443)
- Browser experience: Full reverse proxy with header/cookie/WebSocket preservation
- Example: \`shield http 3000\` or \`shield https 8443\`
- Guide: https://docs.yishield.com/en/protocols/http`,

      telnet: `## Telnet
- Default port: 23
- Browser experience: Terminal for network devices and legacy systems
- Example: \`shield telnet 10.0.0.1\`
- Guide: https://docs.yishield.com/en/protocols/telnet`,

      mysql: `## MySQL (Plugin)
- Default port: 3306
- Browser experience: Full database web admin with schema browser, SQL editor, ER diagrams (via plugin)
- Example: \`shield mysql 10.0.0.20 --username root --auth-pass mypass\`
- Guide: https://docs.yishield.com/en/protocols/mysql`,

      postgres: `## PostgreSQL (Plugin)
- Default port: 5432
- Browser experience: Full database web admin with schema browser, SQL editor, ER diagrams (via plugin)
- Example: \`shield postgres 10.0.0.30 --username postgres --auth-pass mypass\`
- Guide: https://docs.yishield.com/en/protocols/postgres`,

      tcp: `## TCP Proxy
- No default port (must specify)
- Raw TCP port forwarding
- Example: \`shield tcp 3306\``,

      udp: `## UDP Proxy
- No default port (must specify)
- Raw UDP port forwarding
- Example: \`shield udp 53\``,
    };

    if (protocol === "all") {
      return {
        content: [
          {
            type: "text" as const,
            text: `# Shield CLI Supported Protocols\n\n${Object.values(protocols).join("\n\n---\n\n")}`,
          },
        ],
      };
    }

    return {
      content: [{ type: "text" as const, text: protocols[protocol] }],
    };
  }
);

// ─── Tool: plugin information ───────────────────────────────────────────────

server.tool(
  "shield_plugins",
  "Get information about Shield CLI plugins (MySQL, PostgreSQL, SQLServer database web admin)",
  {
    plugin: z
      .enum(["mysql", "postgres", "sqlserver", "all"])
      .default("all")
      .describe("Plugin to query"),
  },
  async ({ plugin }) => {
    const plugins: Record<string, string> = {
      mysql: `## MySQL Plugin
- Protocols: mysql, mariadb
- Default port: 3306
- Features: Database browser, table schema, SQL query, ER diagram, collaborative editing
- Install: \`shield plugin install mysql\`
- Usage: \`shield mysql 10.0.0.20 --username root --auth-pass mypass\`
- Also works standalone via Docker: \`docker run -e DB_HOST=10.0.0.20 fengyily/shield-plugin-mysql\``,

      postgres: `## PostgreSQL Plugin
- Protocols: postgres, pg, postgresql
- Default port: 5432
- Features: Database browser, table schema, SQL query, ER diagram
- Install: \`shield plugin install postgres\`
- Usage: \`shield postgres 10.0.0.20 --username postgres --auth-pass mypass\`
- Source: https://github.com/fengyily/shield-plugins`,

      sqlserver: `## SQL Server Plugin
- Protocols: sqlserver, mssql
- Default port: 1433
- Features: Database browser, table schema, SQL query
- Install: \`shield plugin install sqlserver\`
- Usage: \`shield sqlserver 10.0.0.20 --username sa --auth-pass mypass\`
- Source: https://github.com/fengyily/shield-plugins`,
    };

    if (plugin === "all") {
      return {
        content: [
          {
            type: "text" as const,
            text: `# Shield CLI Plugins\n\nPlugins extend Shield CLI with database web admin capabilities.\n\n\`\`\`bash\nshield plugin install <name>   # Install\nshield plugin list             # List installed\nshield plugin uninstall <name> # Remove\n\`\`\`\n\n${Object.values(plugins).join("\n\n---\n\n")}`,
          },
        ],
      };
    }

    return {
      content: [{ type: "text" as const, text: plugins[plugin] }],
    };
  }
);

// ─── Tool: about Shield CLI ─────────────────────────────────────────────────

server.tool(
  "shield_about",
  "Get an overview of what Shield CLI is, how it works, and how it compares to alternatives",
  {},
  async () => {
    return {
      content: [
        {
          type: "text" as const,
          text: `# Shield CLI

**Access any internal service from your browser. No VPN, no client, one command.**

Shield CLI is a browser-first internal service gateway — SSH terminals, remote desktops, database admin, web apps — all accessible through any browser with a single command.

## How It Works

\`\`\`
Internal Service ←→ Shield CLI ←→ Public Gateway ←→ Browser
  (SSH/RDP/...)    (Encrypted      (HTML5 Render)   (Any Device)
                    WebSocket)
\`\`\`

Shield CLI creates encrypted WebSocket tunnels from your machine to a public gateway. The gateway assigns a unique HTTPS Access URL. Visitors open this URL in a browser — the gateway renders the remote service using HTML5.

## Key Differentiator

Traditional tools (ngrok, frp) solve **network reachability** (L4 port forwarding), but users still need protocol-specific clients.
Shield CLI solves **terminal usability** (L7 protocol rendering) — the browser IS the client.

## Comparison

| Feature | Shield CLI | ngrok | frp |
|---------|-----------|-------|-----|
| Browser RDP/VNC/SSH | Yes | No | No |
| Database Web Admin | Yes (plugins) | No | No |
| Zero client install | Yes | No | No |
| Single binary deploy | Yes | Yes | Yes (2 binaries) |
| Plugin extensibility | Yes | No | No |

## Links

- Repository: https://github.com/fengyily/shield-cli
- Documentation: https://docs.yishield.com
- License: Apache 2.0`,
        },
      ],
    };
  }
);

// ─── Start server ───────────────────────────────────────────────────────────

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch((err) => {
  console.error("Fatal:", err);
  process.exit(1);
});
