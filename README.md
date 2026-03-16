<p align="center">
  <img src="docs/logo.svg" alt="Shield CLI" width="128" height="128">
</p>

<h1 align="center">Shield CLI</h1>

<p align="center">
  <strong>Secure Tunnel Connector</strong><br>
  Expose internal network resources to the public through encrypted tunnels.
</p>

<p align="center">
  <a href="#features">Features</a> &bull;
  <a href="#quick-start">Quick Start</a> &bull;
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

## Features

- **Zero Config** — Just specify protocol and target, everything else is automatic
- **Cross-Platform** — Linux, macOS, Windows native support
- **Encrypted Tunnels** — Built on [chisel](https://github.com/jpillora/chisel) with WebSocket transport
- **Auto Credentials** — Machine fingerprint-based identity, encrypted local storage
- **Dynamic Tunnels** — Local API for runtime tunnel management
- **Reconnection** — Automatic retry with exponential backoff

## Quick Start

```bash
# Expose an SSH service
shield -t ssh -s 10.0.0.2:22

# Expose an HTTP service
shield -t http -s 192.168.1.100:8080

# Expose a Windows Remote Desktop (RDP)
shield -t rdp -s 10.0.0.5:3389

# Expose a VNC service
shield -t vnc -s 10.0.0.10:5900

# Custom API server
shield -t http -s 192.168.1.100:8080 -H https://your-server.com/raas

# Verbose mode for debugging
shield -v -t ssh -s 10.0.0.2:22
```

## Installation

### From Source

```bash
git clone https://github.com/user/shield-cli.git
cd shield-cli
go build -o shield .
```

### Pre-built Binaries

Download from the [Releases](https://github.com/user/shield-cli/releases) page.

## Usage

```
shield [flags]

Flags:
  -t, --type string          Protocol type (ssh, http, https, tcp, rdp, vnc) [required]
  -s, --source string        Target address in ip:port format     [required]
  -H, --server string        API server URL (default: https://console.yishield.com/raas)
  -p, --tunnel-port int      Chisel tunnel server port (default: 62888)
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
    API Tunnel:   remote:63203  ←→  local:4000
    App Tunnel:   remote:58845  ←→  172.16.3.137:22
    Server:       121.43.154.105:62888

  ✓ Tunnel established successfully!

  Site URL:
    https://xxxx.hk01.apps.yishield.com

  Access URL:
    https://hk.svc.yishield.com/plugins/auth?resid=...

  Press Ctrl+C to stop
```

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
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│  Internal     │ ◄─────► │  Shield CLI   │ ◄═════► │  Public      │
│  Service      │  local   │  (tunnel)     │  chisel  │  Server      │
│  10.0.0.2:22  │         │  127.0.0.1    │  wss://  │  1.2.3.4     │
└──────────────┘         └──────────────┘         └──────────────┘
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

Shield CLI 是一个安全隧道连接工具，类似于 ngrok，可以将内网服务安全地暴露到公网。

### 功能特性

- **零配置** — 只需指定协议和目标地址，其余自动完成
- **跨平台** — 原生支持 Linux、macOS、Windows
- **加密隧道** — 基于 [chisel](https://github.com/jpillora/chisel) 的 WebSocket 传输
- **自动凭证** — 基于机器指纹的身份标识，本地加密存储
- **动态隧道** — 运行时通过本地 API 管理隧道
- **断线重连** — 指数退避自动重试

### 快速开始

```bash
# 暴露 SSH 服务
shield -t ssh -s 10.0.0.2:22

# 暴露 HTTP 服务
shield -t http -s 192.168.1.100:8080

# 暴露 Windows 远程桌面 (RDP)
shield -t rdp -s 10.0.0.5:3389

# 暴露 VNC 服务
shield -t vnc -s 10.0.0.10:5900

# 指定 API 服务器
shield -t http -s 192.168.1.100:8080 -H https://your-server.com/raas

# 调试模式
shield -v -t ssh -s 10.0.0.2:22
```

### 安装

#### 从源码编译

```bash
git clone https://github.com/user/shield-cli.git
cd shield-cli
go build -o shield .
```

#### 下载预编译包

前往 [Releases](https://github.com/user/shield-cli/releases) 页面下载对应平台的二进制文件。

### 命令参数

```
shield [flags]

参数:
  -t, --type string          协议类型 (ssh, http, https, tcp, rdp, vnc) [必填]
  -s, --source string        目标地址，格式 ip:port              [必填]
  -H, --server string        API 服务器地址 (默认: https://console.yishield.com/raas)
  -p, --tunnel-port int      隧道服务器端口 (默认: 62888)
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
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│  内网服务      │ ◄─────► │  Shield CLI   │ ◄═════► │  公网服务器    │
│  10.0.0.2:22  │  本地    │  (隧道)       │  chisel  │  1.2.3.4     │
└──────────────┘         └──────────────┘  wss://  └──────────────┘
```

### 安全性

- 凭证使用 AES-256-GCM 加密，密钥基于机器指纹派生
- 所有日志输出中的密码均已脱敏
- 隧道连接使用带认证的 WebSocket 传输
- 凭证文件权限为 `0600`
