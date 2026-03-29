---
title: Shield CLI 视频教程：一条命令，浏览器访问一切内部服务
description: 5 分钟视频教程，演示如何使用 Shield CLI 在浏览器中远程访问 SSH 终端、Windows 桌面和 MySQL 数据库，无需 VPN，无需安装客户端。
date: 2026-03-25
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: Shield CLI 教程, 远程访问教程, 浏览器SSH, 浏览器RDP, MySQL Web管理, 内网穿透, 视频教程
---

# Shield CLI 视频教程：一条命令，浏览器访问一切内部服务

> 远程访问内网服务器，VPN 配置太复杂；让同事临时看一下你的开发环境，对方还得装一堆客户端；出差时想用手机操作 Windows 桌面，没有合适的方案。Shield CLI 用一条命令解决这些问题——生成一个链接，对方在浏览器里直接操作。这篇文章配合 5 分钟视频教程，带你快速上手。

---

## 视频教程

<!-- 将下方 src 替换为实际视频地址（Bilibili / YouTube） -->

<div style="position: relative; padding-bottom: 56.25%; height: 0; overflow: hidden;">
  <iframe
    src="//player.bilibili.com/player.html?bvid=YOUR_BVID&page=1&high_quality=1"
    style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; border: 0;"
    allowfullscreen="true"
    scrolling="no"
  ></iframe>
</div>

> 如果视频加载缓慢，也可以直接访问 [Bilibili 观看](https://www.bilibili.com/video/YOUR_BVID)。

---

## 视频内容概览

这个视频用 5 分钟演示了 Shield CLI 的核心使用场景，从安装到实际操作，全程屏幕录制。

### 1. 安装（0:30）

三个平台，各一条命令：

```bash
# macOS
brew tap fengyily/tap && brew install shield-cli

# Linux
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh

# Windows
scoop bucket add shield https://github.com/fengyily/scoop-bucket && scoop install shield-cli
```

安装完成后执行 `shield --version` 验证。Shield CLI 是单二进制文件，没有任何运行时依赖。

### 2. 启动 Web 管理面板（1:00）

```bash
shield start
```

打开 `http://localhost:8181`，在 Web UI 中添加、管理和连接所有内网服务。填入协议、IP、端口和凭证，点击连接即可生成公网访问链接。

### 3. SSH 终端 — 浏览器里操作远程服务器（2:10）

```bash
shield ssh 10.0.0.5
```

一条命令建立加密隧道，浏览器自动打开 SSH 终端。支持颜色输出、Tab 补全、快捷键，体验与原生 SSH 客户端一致。

![SSH 终端演示](/demo/demo-ssh.gif)

### 4. Windows 远程桌面 — 浏览器里操作 Windows（2:50）

```bash
shield rdp 10.0.0.10
```

整个 Windows 桌面在浏览器中实时渲染，鼠标、键盘、剪贴板全部支持。把链接发给同事，对方不需要安装 RDP 客户端。视频中还演示了在手机浏览器上打开同一个链接，直接操作 Windows 桌面。

![RDP 桌面演示](/demo/demo-rdp.gif)

### 5. MySQL 数据库管理 — 浏览器里查数据（3:40）

```bash
shield plugin add mysql
shield mysql 10.0.0.20:3306 --db-user root --db-pass ****
```

通过插件系统扩展协议支持。安装 MySQL 插件后，浏览器中可以浏览表结构、翻页查看数据、执行 SQL 查询、一键导出 CSV。默认只读模式，防止误操作。

![MySQL 管理演示](/demo/demo-mysql.gif)

---

## 工作原理

```
内网服务 ←→ Shield CLI（Chisel 加密隧道）←→ 公网网关（HTML5 渲染）←→ 浏览器
```

Shield CLI 在本地和公网网关之间建立一条 WebSocket 加密隧道。SSH 通过 xterm.js 在浏览器中渲染终端，RDP 通过 Guacamole 协议渲染桌面，数据库通过内置 Web 客户端管理。所有协议统一通过浏览器访问，对方不需要安装任何软件。

---

## 适用场景

- **远程运维**：在浏览器中 SSH 到内网服务器，不用配 VPN
- **临时协作**：发个链接，同事直接在浏览器操作你的开发环境
- **出差办公**：手机浏览器打开 Windows 远程桌面
- **数据库查询**：浏览器中管理 MySQL，不用装 Navicat / DBeaver
- **演示环境**：给客户展示内网系统，一个链接搞定

---

## 开始使用

```bash
# 安装
brew tap fengyily/tap && brew install shield-cli

# 启动
shield start

# 试试 SSH
shield ssh 127.0.0.1
```

- GitHub：[github.com/fengyily/shield-cli](https://github.com/fengyily/shield-cli)
- 文档：[docs.yishield.com](https://docs.yishield.com)
- 安装指南：[docs.yishield.com/guide/quickstart](https://docs.yishield.com/guide/quickstart)

如果觉得有用，欢迎在 GitHub 上 Star 支持一下。有问题可以提 Issue，也欢迎在视频下方留言讨论。
