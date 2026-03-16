<p align="center">
  <img src="docs/logo.svg" alt="Shield CLI" width="128" height="128">
</p>

<h1 align="center">Shield CLI</h1>

<p align="center">
  <strong>Browser-Based Secure Tunnel for RDP, VNC, SSH & More</strong><br>
  Access internal RDP desktops, VNC sessions, SSH terminals, and web services directly from a browser — no client software required.
</p>

<p align="center">
  <a href="#features">Features</a> &bull;
  <a href="#quick-start">Quick Start</a> &bull;
  <a href="#visibility-modes">Visibility Modes</a> &bull;
  <a href="#installation">Installation</a> &bull;
  <a href="#usage">Usage</a> &bull;
  <a href="#中文文档">中文文档</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-blue?logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-brightgreen" alt="Platform">
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
</p>

---

## Why Shield CLI?

Traditional remote access requires installing dedicated clients (RDP client, VNC viewer, SSH terminal) on every device. Shield CLI eliminates this by tunneling internal services through a secure gateway that renders everything **in the browser**.

- **No client installation** — Open a URL in any browser to access RDP desktops, VNC screens, or SSH terminals
- **Works anywhere** — Access from phones, tablets, locked-down corporate machines, or any device with a browser
- **One binary** — A single `shield` command exposes any internal service to the web

## Features

- **Browser-Based Access** — RDP, VNC, SSH all rendered in the browser via HTML5, no plugins or client software needed
- **Zero Config** — Just specify protocol and target, everything else is automatic
- **Visible & Invisible Modes** — Control whether unauthorized users can access the service via the Site URL alone
- **Cross-Platform** — Linux, macOS, Windows native support
- **Encrypted Tunnels** — Built on [chisel](https://github.com/jpillora/chisel) with WebSocket transport
- **Auto Credentials** — Machine fingerprint-based identity, encrypted local storage
- **Dynamic Tunnels** — Local API for runtime tunnel management
- **Reconnection** — Automatic retry with exponential backoff

## Quick Start

```bash
# Expose a Windows Remote Desktop — access via browser, no RDP client needed
shield -t rdp -s 10.0.0.5:3389

# Expose a VNC server — view the desktop in your browser
shield -t vnc -s 10.0.0.10:5900

# Expose an SSH server — get a terminal in your browser
shield -t ssh -s 10.0.0.2:22

# Expose an HTTP service
shield -t http -s 192.168.1.100:8080
```

Once the tunnel is established, open the **Access URL** in any browser — that's it.

## Visibility Modes

Shield CLI supports two access modes controlled by the `--visable` flag.

When a tunnel is established, two URLs are generated:

- **Site URL** — The application address (e.g., `https://xxxx.hk01.apps.yishield.com`). This URL alone is **not accessible** without authorization.
- **Access URL** — Contains an embedded authorization key. Anyone with this URL can access the service directly.

### Visible Mode

The service is **open to unauthorized users** — anyone who knows the Site URL can access it without an authorization key. Suitable for services that are intended to be publicly reachable.

```bash
# Visible: the Site URL is accessible without authorization
shield --visable -t rdp -s 10.0.0.5:3389

# Filter a specific AC node by name
shield --visable=HK -t ssh -s 10.0.0.2:22
```

**Use cases:**

- Public demo environments
- Shared development servers for the whole team
- Staging environments for QA testing

### Invisible Mode (default)

The service **requires authorization** — the Site URL alone will not grant access. Users must use the Access URL (which contains the authorization key) to connect. This is the secure default for all services.

```bash
# Invisible (default): only the Access URL (with key) grants access
shield -t rdp -s 10.0.0.5:3389

# Secure SSH tunnel — share the Access URL with specific people
shield -t ssh -s 10.0.0.2:22

# Secure VNC access to a lab machine
shield -t vnc -s 192.168.1.50:5900
```

**Use cases:**

- Production servers that must not be publicly accessible
- Temporary access during incident response — share the Access URL, revoke when done
- Sensitive machines where access is restricted to authorized users only

> Both modes print the Site URL and Access URL to the terminal. The difference is whether the Site URL alone is sufficient to access the service, or whether the authorization key in the Access URL is required.

## Installation

### From Source

```bash
git clone https://github.com/user/shield-cli.git
cd shield-cli
go build -o shield .
```

### Pre-built Binaries

Download from the [Releases](https://github.com/user/shield-cli/releases) page. Available for Linux, macOS, and Windows on amd64, arm64, and more.

## Usage

```
shield [flags]

Flags:
  -t, --type string          Protocol type (ssh, http, https, tcp, rdp, vnc) [required]
  -s, --source string        Target address in ip:port format     [required]
  -H, --server string        API server URL (default: https://console.yishield.com/raas)
  -p, --tunnel-port int      Chisel tunnel server port (default: 62888)
      --visable [filter]     Enable visible mode (optional: AC node name filter)
  -v, --verbose              Enable verbose log output
  -h, --help                 Help for shield
```

### Example Output

```
   _____ __    _      __    __   ________    ____
  / ___// /_  (_)__  / /___/ /  / ____/ /   /  _/
  \__ \/ __ \/ // _ \/ // __  / / /   / /    / /
 ___/ / / / / //  __/ // /_/ / / /___/ /____/ /
/____/_/ /_/_/ \___/_/ \__,_/  \____/_____/___/
  Shield CLI - Secure Tunnel Connector

  ⚡ Tunnel Mapping
    App Tunnel:   remote:58845  ←→  172.16.3.137:22
    Server:       121.43.154.105:62888

  ✓ Tunnel established successfully!

  Site URL:
    https://xxxx.hk01.apps.yishield.com

  Access URL:
    https://hk.svc.yishield.com/plugins/auth?resid=...

  Press Ctrl+C to stop
```

Open the **Access URL** in your browser to start using the service. No client software to install.

### Local API

Once running, Shield CLI exposes a local API on `127.0.0.1:<port>`:

| Endpoint | Method | Description |
|---|---|---|
| `/health` | GET | Health check |
| `/connectors` | GET | List all active tunnels |
| `/connector?rport=&lip=&lport=` | GET | Create a dynamic tunnel |
| `/connector?rport=` | DELETE | Close a tunnel |

## Architecture

```
                          Browser (RDP/VNC/SSH via HTML5)
                                    │
                                    ▼
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│  Internal     │ ◄─────► │  Shield CLI   │ ◄═════► │  Public      │
│  Service      │  local   │  (tunnel)     │  chisel  │  Gateway     │
│  10.0.0.5:    │         │  127.0.0.1    │  wss://  │  + Web UI    │
│  3389/5900/22 │         └──────────────┘         └──────────────┘
└──────────────┘
```

## Security

- Credentials are encrypted with AES-256-GCM using a machine-specific fingerprint
- Passwords are masked in all log output
- Tunnel connections use authenticated WebSocket transport
- Credential files are stored with `0600` permissions

## License

MIT

---

<a id="中文文档"></a>

## 中文文档

### 简介

Shield CLI 是一个安全内网穿透工具，支持通过浏览器直接访问内网的 RDP 远程桌面、VNC 屏幕、SSH 终端等服务，**无需安装任何客户端软件**。

与传统方案不同，Shield CLI 不需要在访问端安装 RDP 客户端、VNC Viewer 或 SSH 终端。只需一个浏览器，打开链接即可操作远程桌面或终端。

### 功能特性

- **浏览器直接访问** — RDP、VNC、SSH 均通过 HTML5 在浏览器中渲染，无需安装客户端
- **零配置** — 只需指定协议和目标地址，其余自动完成
- **可见/隐身模式** — 控制未授权用户是否可以通过 Site URL 直接访问服务
- **跨平台** — 原生支持 Linux、macOS、Windows
- **加密隧道** — 基于 [chisel](https://github.com/jpillora/chisel) 的 WebSocket 传输
- **自动凭证** — 基于机器指纹的身份标识，本地加密存储
- **动态隧道** — 运行时通过本地 API 管理隧道
- **断线重连** — 指数退避自动重试

### 快速开始

```bash
# 暴露 Windows 远程桌面 — 用浏览器访问，无需 RDP 客户端
shield -t rdp -s 10.0.0.5:3389

# 暴露 VNC 服务 — 在浏览器中查看桌面
shield -t vnc -s 10.0.0.10:5900

# 暴露 SSH 服务 — 在浏览器中获得终端
shield -t ssh -s 10.0.0.2:22

# 暴露 HTTP 服务
shield -t http -s 192.168.1.100:8080
```

隧道建立后，在任意浏览器中打开 **Access URL** 即可访问。

### 可见与隐身模式

隧道建立后会生成两个 URL：

- **Site URL** — 应用地址（如 `https://xxxx.hk01.apps.yishield.com`）。单独使用此 URL **无法访问**，需要授权。
- **Access URL** — 包含内嵌授权密钥的链接。拥有此 URL 的人可以直接访问服务。

通过 `--visable` 参数控制服务的访问权限：

#### 可见模式

服务**对未授权用户开放** — 任何知道 Site URL 的人都可以直接访问，无需授权密钥。适合需要公开访问的服务。

```bash
# 可见模式：Site URL 无需授权即可访问
shield --visable -t rdp -s 10.0.0.5:3389

# 指定特定 AC 节点
shield --visable=HK -t ssh -s 10.0.0.2:22
```

**适用场景：**

- 公开演示环境
- 团队共享的开发服务器
- QA 测试的预发布环境

#### 隐身模式（默认）

服务**需要授权** — 仅凭 Site URL 无法访问。用户必须使用 Access URL（包含授权密钥）才能连接。这是所有服务的安全默认值。

```bash
# 隐身模式（默认）：只有 Access URL（带密钥）才能访问
shield -t rdp -s 10.0.0.5:3389

# 安全的 SSH 隧道 — 将 Access URL 分享给指定人员
shield -t ssh -s 10.0.0.2:22

# 安全的 VNC 访问实验室机器
shield -t vnc -s 192.168.1.50:5900
```

**适用场景：**

- 不能公开访问的生产服务器
- 故障处理期间的临时访问 — 分享 Access URL，处理完毕即撤销
- 敏感机器，访问仅限授权用户

> 两种模式下终端都会打印 Site URL 和 Access URL。区别在于：可见模式下 Site URL 即可访问；隐身模式下必须使用包含授权密钥的 Access URL。

### 安装

#### 从源码编译

```bash
git clone https://github.com/user/shield-cli.git
cd shield-cli
go build -o shield .
```

#### 下载预编译包

前往 [Releases](https://github.com/user/shield-cli/releases) 页面下载对应平台的二进制文件。支持 Linux、macOS、Windows 的 amd64、arm64 等架构。

### 命令参数

```
shield [flags]

参数:
  -t, --type string          协议类型 (ssh, http, https, tcp, rdp, vnc)  [必填]
  -s, --source string        目标地址，格式 ip:port                       [必填]
  -H, --server string        API 服务器地址 (默认: https://console.yishield.com/raas)
  -p, --tunnel-port int      隧道服务器端口 (默认: 62888)
      --visable [过滤词]     启用可见模式 (可选: AC 节点名称过滤)
  -v, --verbose              启用详细日志输出
  -h, --help                 显示帮助信息
```

### 本地 API

Shield CLI 运行后会在 `127.0.0.1:<port>` 上提供本地管理接口：

| 接口 | 方法 | 说明 |
|---|---|---|
| `/health` | GET | 健康检查 |
| `/connectors` | GET | 列出所有活跃隧道 |
| `/connector?rport=&lip=&lport=` | GET | 创建动态隧道 |
| `/connector?rport=` | DELETE | 关闭隧道 |

### 工作原理

```
                      浏览器 (通过 HTML5 访问 RDP/VNC/SSH)
                                    │
                                    ▼
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│  内网服务      │ ◄─────► │  Shield CLI   │ ◄═════► │  公网网关      │
│  10.0.0.5:    │  本地    │  (隧道)       │  chisel  │  + Web UI    │
│  3389/5900/22 │         └──────────────┘  wss://  └──────────────┘
└──────────────┘
```

### 安全性

- 凭证使用 AES-256-GCM 加密，密钥基于机器指纹派生
- 所有日志输出中的密码均已脱敏
- 隧道连接使用带认证的 WebSocket 传输
- 凭证文件权限为 `0600`
