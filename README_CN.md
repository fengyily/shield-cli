<p align="center">
  <img src="docs/logo.svg" alt="Shield CLI" width="128" height="128">
</p>

<h1 align="center">Shield CLI</h1>

<p align="center">
  <strong>一条命令，一个 URL，浏览器即可访问内网资源</strong><br>
  通过浏览器直接访问内网的 RDP 远程桌面、VNC 屏幕、SSH 终端和 Web 服务 — 无需 VPN，无需安装任何客户端。
</p>

<p align="center">
  <a href="https://docs.yishield.com/guide/what-is-shield">文档中心</a> &bull;
  <a href="https://docs.yishield.com/guide/install">安装</a> &bull;
  <a href="https://docs.yishield.com/guide/quickstart">快速开始</a> &bull;
  <a href="README.md">English</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-blue?logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-brightgreen" alt="Platform">
  <img src="https://img.shields.io/badge/license-Apache%202.0-green" alt="License">
</p>

---

## 演示

### RDP — 浏览器远程桌面

<p align="center">
  <img src="docs/demo/demo-rdp.gif" alt="Shield CLI RDP 演示" width="960">
</p>

### SSH — 浏览器终端

<p align="center">
  <img src="docs/demo/demo-ssh.gif" alt="Shield CLI SSH 演示" width="960">
</p>

---

## 为什么选择 Shield CLI？

传统隧道工具（ngrok、frp）解决的是**网络可达** — 端口映射到公网，但对方仍需安装 RDP 客户端、SSH 终端、VNC Viewer。

Shield CLI 解决的是**终端可用** — 在浏览器中直接渲染 RDP 桌面、VNC 屏幕和 SSH 终端（HTML5），对方只需一个浏览器。

| 特性 | Shield CLI | ngrok | frp |
|------|-----------|-------|-----|
| 浏览器内 RDP/VNC | 支持 | 不支持 | 不支持 |
| 浏览器内 SSH 终端 | 支持 | 不支持 | 不支持 |
| TCP 隧道免费 | 免费 | 需付费 | 免费（需自建） |
| 零客户端安装 | 是 | 否 | 否 |
| 国内安装友好 | jsDelivr 镜像 | 需翻墙 | 可直连 |

## 安装

```bash
# macOS
brew tap fengyily/tap && brew install shield-cli

# Windows
scoop bucket add shield https://github.com/fengyily/scoop-bucket && scoop install shield-cli

# Linux (apt) — Debian / Ubuntu
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/scripts/setup-repo.sh | sudo bash

# Linux (yum) — RHEL / CentOS / Fedora
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/scripts/setup-repo.sh | sudo bash

# Linux / macOS 一键安装（直接下载二进制）
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh

# 国内镜像（推荐）
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

### Docker

```bash
# 使用预构建镜像（推荐）
docker run -d --name shield \
  --network host \
  --restart unless-stopped \
  fengyily/shield-cli

# 或从源码构建
docker build -t shield-cli .
docker run -d --name shield --network host --restart unless-stopped shield-cli
```

> **说明：** `--network host` 让容器直接使用宿主机网络栈，Shield CLI 可以访问宿主机本机服务以及宿主机所在的内网资源（如 `10.0.0.x`、`192.168.x.x`）。启动后访问 `http://localhost:8181` 即可使用 Web UI。
>
> **注意：** `--network host` 仅在 **Linux** 上生效。macOS 和 Windows 的 Docker Desktop 不支持 host 网络模式，可改用端口映射：
>
> ```bash
> docker run -d --name shield -p 8181:8181 --restart unless-stopped fengyily/shield-cli
> ```

更多安装方式（apt、yum、deb、rpm、PowerShell、源码编译）：[安装指南](https://docs.yishield.com/guide/install)

## 快速开始

### Web UI 模式（推荐）

```bash
shield start
```

打开 `http://localhost:8181`，添加服务，一键连接。macOS 和 Windows 上会在系统托盘显示图标，点击即可快速打开 Dashboard。

![Web 管理面板](docs/images/shieldcli-webui-001.jpg)

![通过 Web UI 访问 RDP](docs/images/shieldcli-rdp-web-001.jpg)

### 系统服务安装（开机自启）

```bash
shield install              # 安装为系统服务（默认端口 8181）
shield install --port 8182  # 如果 8181 被占用，指定其他端口
shield uninstall            # 卸载服务
```

支持 macOS (launchd)、Linux (systemd) 和 Windows。详见[系统服务安装指南](https://docs.yishield.com/guide/system-service)。

### 命令行模式

```bash
shield ssh              # 浏览器内 SSH 终端 (127.0.0.1:22)
shield rdp 10.0.0.5     # 浏览器内 Windows 桌面
shield http 3000        # 暴露本地 Web 应用
shield vnc 10.0.0.10    # 浏览器内 VNC 屏幕共享
shield tcp 3306         # TCP 端口代理（MySQL）
shield udp 53           # UDP 端口代理（DNS）
```

![Shield CLI 终端](docs/images/shieldcli-ssh-001.jpg)

![浏览器 SSH 终端](docs/images/shieldcli-ssh-web-002.jpg)

### 智能默认值

| 命令 | 解析为 |
|------|--------|
| `shield ssh` | `127.0.0.1:22` |
| `shield ssh 2222` | `127.0.0.1:2222` |
| `shield ssh 10.0.0.2` | `10.0.0.2:22` |
| `shield rdp` | `127.0.0.1:3389` |
| `shield http 3000` | `127.0.0.1:3000` |
| `shield tcp 3306` | `127.0.0.1:3306` |
| `shield udp 53` | `127.0.0.1:53` |

支持协议：`ssh`、`rdp`、`vnc`、`http`、`https`、`telnet`、`tcp`、`udp` — [完整命令参考](https://docs.yishield.com/reference/commands)

## 工作原理

```
                      浏览器 (通过 HTML5 访问 RDP/VNC/SSH)
                                │
                                ▼
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│  内网服务      │ ◄──► │  Shield CLI   │ ◄══► │  公网网关      │
│  10.0.0.5:    │ 本地  │  (加密隧道)    │chisel│  + HTML5     │
│  3389/5900/22 │      └──────────────┘wss://│  协议渲染     │
└──────────────┘                             └──────────────┘
```

了解更多：[连接流程](https://docs.yishield.com/security/connection-flow) | [安全模型](https://docs.yishield.com/security/credentials)

## 安全性

- **AES-256-GCM 加密** — 凭证使用机器指纹派生密钥加密
- **密码脱敏** — 日志中所有密码自动隐藏
- **WebSocket 传输** — 带认证的加密隧道
- **0600 权限** — 凭证文件仅当前用户可读

详情：[凭证管理](https://docs.yishield.com/security/credentials) | [访问模式](https://docs.yishield.com/security/access-modes)

## 文档中心

完整文档请访问 **[docs.yishield.com](https://docs.yishield.com)**：

- [Shield CLI 是什么](https://docs.yishield.com/guide/what-is-shield) — 概述和核心特性
- [安装指南](https://docs.yishield.com/guide/install) — 所有安装方式
- [5 分钟上手](https://docs.yishield.com/guide/quickstart) — 快速入门教程
- [协议指南](https://docs.yishield.com/protocols/ssh) — SSH、RDP、VNC、HTTP、Telnet
- [命令参考](https://docs.yishield.com/reference/commands) — 完整参数列表
- [常见问题](https://docs.yishield.com/reference/faq) — FAQ
- [故障排查](https://docs.yishield.com/troubleshooting/errors) — 常见错误和解决方案

## 许可证

Apache 2.0
