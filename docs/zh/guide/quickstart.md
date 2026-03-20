---
title: 5 分钟上手 Shield CLI
description: 5 分钟快速上手 Shield CLI，学习通过 Web UI 或命令行创建 SSH、RDP、HTTP 隧道，体验智能地址解析。
head:
  - - meta
    - name: keywords
      content: Shield CLI 入门, 快速上手, SSH隧道, RDP隧道, HTTP隧道, 教程
---

# 5 分钟上手

本教程将带你从安装到成功访问第一个内网服务。

## 方式一：Web UI 模式（推荐）

### 第一步：启动管理面板

```bash
shield start
```

浏览器会自动打开 `http://localhost:8181`。

### 第二步：添加应用

在 Web UI 中点击添加应用，填写：

- **协议**：选择 SSH / RDP / VNC / HTTP 等
- **IP 地址**：目标服务的内网 IP（如 `10.0.0.5`）
- **端口**：目标服务端口（如 SSH 默认 `22`）
- **名称**：给这个应用起个名字（如「办公室电脑」）

### 第三步：连接

点击应用卡片上的 **连接** 按钮，等待状态变为「已连接」后，浏览器会自动打开 Access URL。

你现在可以直接在浏览器中操作远程服务了。

---

## 方式二：命令行模式

### 连接本机 SSH

```bash
shield ssh
```

Shield CLI 会：
1. 自动解析为 `127.0.0.1:22`
2. 提示输入目标服务器的用户名和密码
3. 建立加密隧道
4. 输出 Access URL

将 URL 在浏览器中打开，即可看到 SSH 终端。

### 连接远程 Windows 桌面

```bash
shield rdp 10.0.0.5
```

输入 Windows 登录凭证后，浏览器中即可看到完整的远程桌面。

### 暴露本地 Web 应用

```bash
shield http 3000
```

你在 `localhost:3000` 运行的 Web 应用，现在可以通过公网 URL 访问了。

## 智能地址解析

Shield CLI 支持灵活的地址输入方式：

| 输入 | 解析为 |
|---|---|
| `shield ssh` | `127.0.0.1:22` |
| `shield ssh 2222` | `127.0.0.1:2222` |
| `shield ssh 10.0.0.2` | `10.0.0.2:22` |
| `shield ssh 10.0.0.2:2222` | `10.0.0.2:2222` |
| `shield rdp` | `127.0.0.1:3389` |
| `shield vnc` | `127.0.0.1:5900` |
| `shield http` | `127.0.0.1:80` |
| `shield http 3000` | `127.0.0.1:3000` |

每种协议都有预设的默认端口，你只需要输入不同的部分。

## 下一步

- [Web UI 模式详解](./web-ui.md)
- [命令行模式详解](./cli-mode.md)
- [协议指南](../protocols/ssh.md)
