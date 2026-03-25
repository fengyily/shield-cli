---
title: 家里装了 OpenClaw，在公司也能随时管理——Shield CLI 远程访问方案
description: 家里电脑装了 OpenClaw，出门在外想配置怎么办？Shield CLI 一条命令搞定远程访问，Dashboard 和 Windows 桌面都能随时连。
date: 2026-03-25
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: OpenClaw 远程访问, Shield CLI, OpenClaw 远程管理, RDP 远程桌面, OpenClaw 内网穿透, AI Agent 远程配置, 安全隧道
---

# 家里装了 OpenClaw，在公司也能随时管理

> OpenClaw 火到不用介绍了——GitHub 25 万 Star，一个能真正帮你干活的 AI Agent。很多人装在家里的 Windows 电脑上，配好了 API Key 和各种插件，用着很爽。但一到公司或者出门在外，就只能干等着回家再弄。这篇文章解决一个具体问题：**怎么在任何地方远程管理家里的 OpenClaw**。

---

## 问题

OpenClaw 的 Dashboard 默认监听在 `localhost:18789`，只有本机能访问。你在家配好了一切，但到了公司：

- 想调整 OpenClaw 的配置，改一下 Prompt 模板或者换个模型——打不开 Dashboard
- 想给 OpenClaw 下个任务让它跑着，下班回家看结果——没法远程操作
- OpenClaw 跑的任务出了问题，想上去看看日志、改改配置——只能等回家
- 需要改一下 OpenClaw 所在的系统配置（装个依赖、改个环境变量）——更没办法了

OpenClaw 官方提供了 [Remote over SSH](https://docs.openclaw.ai/zh-CN/gateway/remote) 模式——在远端机器上开 SSH 隧道转发 18789 端口，然后客户端连 `ws://127.0.0.1:18789`。能用，但有几个前提：

- 家里电脑要开 SSH 服务（Windows 默认没开，需要手动配置 OpenSSH Server）
- 需要路由器做端口映射，把 SSH 22 端口暴露到公网（或者配 DDNS）
- 每次连接要先手动开 SSH 隧道：`ssh -L 18789:127.0.0.1:18789 user@home-ip`
- 只解决了 Dashboard 访问，如果还想远程桌面操作系统，得再开一条 RDP 隧道

Shield CLI 的做法更简单：**不开端口、不碰路由器、不配 SSH，两条命令覆盖所有场景**。

---

## 思路

Shield CLI 做两件事就够了：

1. **代理 OpenClaw Dashboard**（HTTP 隧道）— 在浏览器里管理 OpenClaw 的配置、下任务、看结果
2. **代理 Windows 远程桌面**（RDP 隧道）— 需要动系统层面的东西时，直接远程桌面进去操作

两条命令跑起来，记下专属 URL，在公司、在咖啡厅、在手机上，随时随地都能管理家里的 OpenClaw。

---

## 准备工作

在家里的 Windows 电脑上：

### 1. 确认 OpenClaw 已跑起来

Dashboard 默认在 `localhost:18789`，浏览器打开能看到就行。

### 2. 开启 Windows 远程桌面

设置 → 系统 → 远程桌面 → 打开。RDP 默认端口 `3389`。

### 3. 安装 Shield CLI

```powershell
# PowerShell
irm https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.ps1 | iex
```

---

## 方式一：通过 Web 管理面板（推荐）

Shield CLI 自带 Web 管理面板，可以把 OpenClaw Dashboard 和 Windows 远程桌面统一管理，最直观。

```bash
shield start
```

浏览器打开 `http://localhost:8181`，点 **Add App**，分别添加两个应用：

| 应用 | Protocol | Target IP | Port |
| --- | --- | --- | --- |
| OpenClaw Dashboard | `http` | `127.0.0.1` | `18789` |
| Windows 远程桌面 | `rdp` | `127.0.0.1` | `3389` |

点击连接，Shield 自动建立加密隧道并生成专属 URL。两个 URL 记在手机备忘录里，到了公司随时用：

- **轻操作**（改配置、下任务、看结果）→ 打开 OpenClaw Dashboard URL
- **重操作**（装依赖、改系统配置）→ 打开 RDP URL

Web 面板的好处是所有应用集中管理，随时连接/断开，一目了然。

---

## 方式二：命令行直接代理

如果偏好命令行，也可以直接用命令：

### 代理 OpenClaw Dashboard

```bash
shield http 127.0.0.1:18789
```

输出：

```text
  ✓ Tunnel established

  Public URL:  https://abc123.yishield.com
  Local:       http://127.0.0.1:18789

  Share this URL. Press Ctrl+C to stop.
```

**记下这个 URL**。到了公司，浏览器打开就是家里的 OpenClaw Dashboard。能做的事：

- 调整 Prompt 模板、切换模型（Claude / DeepSeek / GPT）
- 给 OpenClaw 下任务，让它跑着，下班回家看结果
- 查看任务执行日志和历史
- 管理插件、配置 MCP Server
- 通过 OpenClaw 内置的聊天界面和 AI 对话

适合改配置、下任务、看结果这类**轻操作**，手机上也能用。

### 代理 Windows 远程桌面

有时候光靠 Dashboard 不够——需要装个 Python 依赖、改个环境变量、调试一个 OpenClaw 插件、看看系统日志。这时候需要完整的桌面操作。

```bash
shield rdp 127.0.0.1
```

输出：

```text
  ✓ Tunnel established

  Remote Desktop:  https://xyz789.yishield.com
  Local:           127.0.0.1:3389

  Open the URL in browser to connect. Press Ctrl+C to stop.
```

**记下这个 URL**。到了公司，浏览器打开，输入 Windows 用户名和密码，就能看到家里电脑的桌面。适合安装依赖、改系统配置、查日志等**重操作**。

---

## 让 Shield 开机自启

每次开机手动跑命令太麻烦。用 Shield 的服务模式，开机自动启动：

```bash
shield start --install
```

这会把 Shield 注册为系统服务，开机自动运行。之后家里电脑开机就自动建立隧道，URL 不变，到了公司直接用。

---

## 安全考量

**OpenClaw 自带认证。** Dashboard 需要 token 才能访问（`?token=xxx`），远程用户没有 token 打不开界面。

**API Key 不暴露。** Key 存在家里电脑的 `~/.openclaw/` 里，远程访问的是 Web 界面，看不到也拿不到。

**加密传输。** Shield 隧道全程 HTTPS/WSS 加密，公司网络管理员也看不到传输内容。

**RDP 有密码保护。** Windows 远程桌面本身需要用户名密码认证，建议设置强密码。

**随时断开。** 不用的时候 `Ctrl+C` 或在 Web 面板断开，链接立即失效。

---

## 不只是 OpenClaw

同样的方式适用于家里电脑上跑的任何服务：

```bash
# 家里的 Jupyter Notebook
shield http 127.0.0.1:8888

# 家里的 Home Assistant 智能家居
shield http 127.0.0.1:8123

# 家里 NAS 的管理页面
shield http 192.168.1.100:5000

# 家里电脑上的 MySQL
shield mysql 127.0.0.1:3306 --db-user root --readonly
```

一个 Shield CLI，把家里的所有服务都带出门。

---

## 试一下

```bash
# 安装 Shield CLI（Windows PowerShell）
irm https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.ps1 | iex

# 代理 OpenClaw Dashboard（记下 URL）
shield http 127.0.0.1:18789

# 代理 Windows 远程桌面（记下 URL）
shield rdp 127.0.0.1

# 或者用 Web 面板统一管理
shield start
```

两个 URL 记在手机里，上班路上就能开始用。

开源地址：https://github.com/fengyily/shield-cli

有问题欢迎提 [Issue](https://github.com/fengyily/shield-cli/issues)。
