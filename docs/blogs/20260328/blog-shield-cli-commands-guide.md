---
title: Shield CLI 命令全解析：15 个命令覆盖所有远程访问场景
description: 详细介绍 Shield CLI 的全部命令，包括隧道创建、智能地址解析、服务管理、插件系统和访问模式，附带实用示例和参数速查表。
date: 2026-03-28
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: Shield CLI 命令, shield ssh, shield rdp, shield start, 远程访问命令, 内网穿透命令, Shield CLI 教程, 命令行工具
---

# Shield CLI 命令全解析：15 个命令覆盖所有远程访问场景

> Shield CLI 的设计哲学是"一条命令搞定"。15 个命令覆盖了 SSH 终端、Windows 桌面、VNC 共享、HTTP/HTTPS 应用、TCP/UDP 代理、数据库管理、系统服务和插件管理。这篇文章把每个命令的用法、参数和典型场景讲清楚，方便随时查阅。

---

## 一、隧道命令

隧道命令是 Shield CLI 的核心，一条命令创建加密隧道，生成浏览器可访问的公网链接。

### shield ssh — SSH 终端

在浏览器中打开远程服务器的 SSH 终端。

```bash
shield ssh [address]
```

**示例：**

```bash
shield ssh                          # → 127.0.0.1:22
shield ssh 2222                     # → 127.0.0.1:2222
shield ssh 10.0.0.5                 # → 10.0.0.5:22
shield ssh 10.0.0.5:2222            # → 10.0.0.5:2222
```

**SSH 专用参数：**

| 参数 | 说明 |
|------|------|
| `--username` | SSH 用户名 |
| `--auth-pass` | SSH 密码 |
| `--private-key` | 私钥文件路径 |
| `--passphrase` | 私钥密码 |
| `--enable-sftp` | 启用 SFTP 文件传输 |

**完整示例：**

```bash
# 密码登录
shield ssh 10.0.0.5 --username root --auth-pass mypassword

# 密钥登录
shield ssh 10.0.0.5 --username deploy --private-key ~/.ssh/id_rsa

# 密钥 + 密码短语
shield ssh 10.0.0.5 --username deploy --private-key ~/.ssh/id_rsa --passphrase mypass

# 开启 SFTP 文件传输
shield ssh 10.0.0.5 --username root --auth-pass mypass --enable-sftp
```

**典型场景：** 运维人员远程登录内网服务器，或给同事发一个链接让对方在浏览器中操作。

---

### shield rdp — Windows 远程桌面

在浏览器中打开 Windows 远程桌面。

```bash
shield rdp [address]
```

**示例：**

```bash
shield rdp                          # → 127.0.0.1:3389
shield rdp 10.0.0.10                # → 10.0.0.10:3389
shield rdp 10.0.0.10:3390           # → 10.0.0.10:3390
```

**典型场景：** 出差时用手机浏览器操作公司 Windows 电脑，或让客户远程查看演示环境。

---

### shield vnc — VNC 屏幕共享

在浏览器中查看和操作 VNC 远程屏幕。

```bash
shield vnc [address]
```

**示例：**

```bash
shield vnc                          # → 127.0.0.1:5900
shield vnc 10.0.0.15                # → 10.0.0.15:5900
shield vnc 10.0.0.15:5901           # → 10.0.0.15:5901
```

**典型场景：** 远程协助 macOS / Linux 桌面，或在浏览器中查看嵌入式设备屏幕。

---

### shield http / shield https — Web 应用隧道

将本地或内网的 Web 应用暴露到公网。

```bash
shield http [address]
shield https [address]
```

**示例：**

```bash
shield http 3000                    # → 127.0.0.1:3000（本地开发服务器）
shield http 10.0.0.20:8080          # → 内网 Web 应用
shield https 10.0.0.20              # → 10.0.0.20:443
```

**典型场景：** 前端开发时让外网同事预览本地页面，或临时将内网管理后台共享给远程团队。

---

### shield telnet — Telnet 隧道

```bash
shield telnet [address]
```

**示例：**

```bash
shield telnet                       # → 127.0.0.1:23
shield telnet 10.0.0.30             # → 10.0.0.30:23
```

**典型场景：** 远程管理网络设备（交换机、路由器）。

---

### shield tcp / shield udp — 通用端口代理

代理任意 TCP 或 UDP 端口，不限协议。

```bash
shield tcp <port|address>
shield udp <port|address>
```

> 注意：TCP/UDP 没有默认端口，必须指定。

**示例：**

```bash
shield tcp 3306                     # → 127.0.0.1:3306（MySQL）
shield tcp 10.0.0.20:6379           # → 内网 Redis
shield udp 10.0.0.50:53             # → DNS 服务器
```

**典型场景：** 代理数据库端口、Redis、消息队列等任何基于 TCP/UDP 的服务。

---

## 二、智能地址解析

所有隧道命令都支持四种地址格式，Shield CLI 会根据协议自动补全默认端口：

| 输入格式 | 解析结果 | 示例 |
|----------|---------|------|
| 省略 | `127.0.0.1:默认端口` | `shield ssh` → `127.0.0.1:22` |
| 仅端口 | `127.0.0.1:指定端口` | `shield ssh 2222` → `127.0.0.1:2222` |
| 仅 IP | `IP:默认端口` | `shield ssh 10.0.0.5` → `10.0.0.5:22` |
| 完整地址 | `IP:指定端口` | `shield ssh 10.0.0.5:2222` → `10.0.0.5:2222` |

**各协议默认端口：**

| 协议 | 默认端口 |
|------|---------|
| SSH | 22 |
| RDP | 3389 |
| VNC | 5900 |
| HTTP | 80 |
| HTTPS | 443 |
| Telnet | 23 |
| TCP/UDP | 无（必须指定） |

---

## 三、访问模式

Shield CLI 支持两种访问模式，控制隧道链接的可见性。

### 可见模式（默认）

```bash
shield ssh 10.0.0.5 --visable
```

默认行为，任何人拿到链接即可访问。还可以指定接入节点：

```bash
shield ssh 10.0.0.5 --visable=HK    # 指定香港节点接入
```

### 隐身模式

```bash
shield ssh 10.0.0.5 --invisible
```

隐身模式下，访问链接需要授权码才能打开，适合对安全性要求更高的场景。

---

## 四、全局参数

以下参数适用于所有隧道命令：

| 参数 | 说明 | 示例 |
|------|------|------|
| `--username` | 目标服务用户名 | `--username root` |
| `--auth-pass` | 目标服务密码 | `--auth-pass mypass` |
| `--server` | 自定义服务端地址 | `--server gateway.example.com` |
| `--visable` | 可见模式（默认） | `--visable=HK` |
| `--invisible` | 隐身模式，需授权码 | `--invisible` |

---

## 五、Web 管理面板

### shield start — 启动管理面板

```bash
shield start [port]
```

**示例：**

```bash
shield start                        # → http://localhost:8181
shield start 9090                   # → http://localhost:9090
```

启动后在浏览器中打开，可以通过 Web UI 添加、管理和连接所有应用，适合不想记命令的用户。

---

## 六、系统服务管理

将 Shield CLI 注册为系统服务，实现开机自启。

| 命令 | 功能 |
|------|------|
| `shield install` | 安装为系统服务（macOS: launchd / Linux: systemd / Windows: 服务） |
| `shield install --port 9090` | 安装时指定 Web UI 端口 |
| `shield uninstall` | 卸载系统服务 |
| `shield stop` | 停止服务 |
| `shield clean` | 清除本地凭证缓存 |

**示例：**

```bash
# 安装为系统服务，开机自启
shield install

# 自定义端口
shield install --port 9090

# 卸载服务
shield uninstall

# 清除本地保存的凭证
shield clean
```

> `shield install` 会自动检测端口冲突。如果 8181 已被占用，会提示你指定其他端口。

---

## 七、插件管理

Shield CLI 通过插件系统扩展协议支持。插件是独立的二进制文件，通过 stdin/stdout JSON 与主程序通信。

| 命令 | 功能 |
|------|------|
| `shield plugin add <name>` | 安装插件 |
| `shield plugin list` | 查看已安装插件 |
| `shield plugin remove <name>` | 卸载插件 |

**示例：**

```bash
# 安装 MySQL 插件
shield plugin add mysql

# 查看已安装插件
shield plugin list
# NAME      VERSION  PROTOCOLS         INSTALLED
# mysql     v0.1.0   mysql, mariadb    2026-03-24T10:00:00+08:00

# 卸载插件
shield plugin remove mysql
```

**数据库插件专用参数：**

| 参数 | 说明 |
|------|------|
| `--db-user` | 数据库用户名 |
| `--db-pass` | 数据库密码 |
| `--db-name` | 数据库名 |
| `--readonly` | 强制只读模式，禁止写操作 |

**数据库插件示例：**

```bash
# MySQL
shield plugin add mysql
shield mysql 10.0.0.20:3306 --db-user root --db-pass mypass

# PostgreSQL
shield plugin add postgres
shield postgres 10.0.0.20:5432 --db-user postgres --db-pass mypass --db-name mydb

# 只读模式
shield mysql 10.0.0.20:3306 --db-user root --db-pass mypass --readonly
```

目前已有的插件：`mysql`、`postgres`、`sqlserver`。

---

## 八、命令速查表

一张表看完所有命令：

| 命令 | 用途 | 示例 |
|------|------|------|
| `shield ssh [addr]` | SSH 终端 | `shield ssh 10.0.0.5` |
| `shield rdp [addr]` | Windows 桌面 | `shield rdp 10.0.0.10` |
| `shield vnc [addr]` | VNC 屏幕共享 | `shield vnc 10.0.0.15` |
| `shield http [addr]` | HTTP 应用 | `shield http 3000` |
| `shield https [addr]` | HTTPS 应用 | `shield https 10.0.0.20` |
| `shield telnet [addr]` | Telnet | `shield telnet 10.0.0.30` |
| `shield tcp <addr>` | TCP 代理 | `shield tcp 3306` |
| `shield udp <addr>` | UDP 代理 | `shield udp 10.0.0.50:53` |
| `shield start [port]` | 启动 Web UI | `shield start 9090` |
| `shield install` | 安装系统服务 | `shield install --port 9090` |
| `shield uninstall` | 卸载系统服务 | `shield uninstall` |
| `shield stop` | 停止服务 | `shield stop` |
| `shield clean` | 清除凭证缓存 | `shield clean` |
| `shield plugin add` | 安装插件 | `shield plugin add mysql` |
| `shield plugin list` | 查看插件 | `shield plugin list` |
| `shield plugin remove` | 卸载插件 | `shield plugin remove mysql` |
| `shield mysql [addr]` | MySQL 管理 | `shield mysql 10.0.0.20:3306` |
| `shield postgres [addr]` | PostgreSQL 管理 | `shield postgres 10.0.0.20:5432` |

---

## 九、安装 Shield CLI

三个平台，各一条命令：

### macOS

```bash
brew tap fengyily/tap && brew install shield-cli
```

### Linux / macOS（一键脚本）

```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
```

国内加速：

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

### Windows

```powershell
scoop bucket add shield https://github.com/fengyily/scoop-bucket
scoop install shield-cli
```

### Docker

```bash
docker run -d --name shield --network host --restart unless-stopped fengyily/shield-cli
```

### 验证安装

```bash
shield --version
```

Shield CLI 是单二进制文件，没有运行时依赖。安装完成后直接使用，不需要配置环境变量或安装额外组件。

---

## 十、了解更多

Shield CLI 完全开源，Apache 2.0 协议。

- **GitHub**：[github.com/fengyily/shield-cli](https://github.com/fengyily/shield-cli) — 欢迎 Star 支持
- **完整文档**：[docs.yishield.com](https://docs.yishield.com)
- **安装指南**：[docs.yishield.com/guide/quickstart](https://docs.yishield.com/guide/quickstart)
- **命令参考**：[docs.yishield.com/reference/commands](https://docs.yishield.com/reference/commands)
- **视频教程**：[5 分钟上手 Shield CLI](https://docs.yishield.com/blogs/20260325/blog-shield-cli-video-tutorial)

有问题可以在 GitHub 提 Issue，也欢迎加入社区讨论。如果 Shield CLI 对你有帮助，给个 Star 就是最大的支持。
