---
title: Shield CLI 是什么 — 安全隧道连接器
description: Shield CLI 是一个安全隧道连接器，通过一条命令将内网服务（SSH、RDP、VNC、HTTP）暴露到公网，任何人只需浏览器即可访问，无需安装客户端。
head:
  - - meta
    - name: keywords
      content: Shield CLI, 安全隧道, 内网穿透, 远程访问, 浏览器终端, SSH浏览器, RDP浏览器, ngrok替代
---

# Shield CLI 是什么

Shield CLI 是一个**安全隧道连接器**，让你通过一条命令将内网服务暴露到公网，任何人只需一个浏览器即可访问。

## 它解决什么问题

在日常工作中，你可能经常遇到这些场景：

- 需要远程访问公司内网的服务器，但配置 VPN 太麻烦
- 想让同事临时访问你本地的开发环境
- 需要从手机或平板远程操作一台 Windows 桌面
- 客户需要访问内部部署的 Web 应用做验收

传统方案（VPN、端口转发、内网穿透）要么配置复杂，要么需要对方安装客户端。

**Shield CLI 的方式是：**

```bash
shield ssh 10.0.0.5
```

执行后你会得到一个公网 URL，对方在浏览器中打开即可直接操作 SSH 终端。无需安装任何客户端、无需配置网络。

## 核心特性

| 特性 | 说明 |
|---|---|
| **浏览器直连** | RDP、VNC、SSH、Web 应用全部通过 HTML5 在浏览器中运行 |
| **零客户端** | 访问者只需一个浏览器，手机、平板、受限电脑均可 |
| **加密传输** | WebSocket 加密隧道 + AES-256-GCM 本地凭证加密 |
| **智能默认** | `shield ssh` 自动解析为 `127.0.0.1:22`，减少输入 |
| **双模式** | Web UI 管理面板（推荐）+ 纯命令行模式 |
| **六大协议** | SSH、RDP、VNC、HTTP、HTTPS、Telnet |
| **全平台** | macOS / Linux / Windows，支持 amd64 和 arm64 |

## 工作原理

```
你的内网服务 ←→ Shield CLI ←→ 公网网关 ←→ 浏览器
   (SSH/RDP/...)    (加密隧道)    (HTML5 渲染)    (任意设备)
```

1. Shield CLI 在你的机器上运行，连接到内网服务
2. 通过加密 WebSocket 隧道与公网网关建立连接
3. 网关分配一个唯一的 Access URL
4. 访问者在浏览器中打开 URL，即可操作远程服务

## 两种使用方式

### Web UI 模式（推荐）

启动本地管理面板，在浏览器中管理所有应用：

```bash
shield start
```

打开 `http://localhost:8181`，通过图形界面添加应用、一键连接。

### 命令行模式

直接在终端创建隧道，适合服务器环境或脚本自动化：

```bash
shield ssh           # 连接本机 SSH
shield rdp 10.0.0.5  # 连接远程 Windows 桌面
shield http 3000     # 暴露本地 Web 应用
```

## 下一步

- [安装 Shield CLI](./install.md)
- [5 分钟上手教程](./quickstart.md)
