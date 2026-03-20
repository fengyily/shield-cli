---
title: 命令行模式 — 终端创建隧道
description: 使用 Shield CLI 命令行模式在终端快速创建 SSH、RDP、VNC、HTTP 隧道，适合服务器环境和脚本自动化。
head:
  - - meta
    - name: keywords
      content: Shield CLI 命令行, CLI模式, 终端, 脚本自动化, SSH隧道命令
---

# 命令行模式

命令行模式适合服务器环境、脚本自动化，或者你只是想快速建立一个隧道。

## 基本用法

```bash
shield <协议> [地址]
```

### 示例

```bash
# 连接本机 SSH (127.0.0.1:22)
shield ssh

# 连接指定 IP 的 RDP
shield rdp 10.0.0.5

# 暴露本地 3000 端口的 Web 应用
shield http 3000

# 连接指定 IP 和端口的 VNC
shield vnc 10.0.0.10:5901
```

## 认证参数

### 交互式输入（默认）

不带认证参数时，Shield CLI 会交互式提示输入：

```bash
shield ssh 10.0.0.5
# → 提示输入用户名
# → 提示输入密码（隐藏输入）
```

### 命令行传参

```bash
# 用户名 + 密码
shield ssh 10.0.0.5 --username root --auth-pass mypassword

# SSH 私钥认证
shield ssh 10.0.0.5 --username root --private-key ~/.ssh/id_rsa

# 带密码的私钥
shield ssh 10.0.0.5 --username root --private-key ~/.ssh/id_rsa --passphrase mypass
```

## 访问模式

```bash
# 可见模式（默认）— 公开 Access URL
shield ssh 10.0.0.5

# 指定接入节点
shield ssh 10.0.0.5 --visable=HK

# 隐身模式 — 需要授权码访问
shield ssh 10.0.0.5 --invisible
```

详见 [访问模式](../security/access-modes.md)。

## SFTP 支持

SSH 协议可以同时启用 SFTP 文件传输：

```bash
shield ssh 10.0.0.5 --enable-sftp
```

## 连接输出

成功建立隧道后，终端会显示：

```
Shield CLI v1.x.x

Protocol: SSH
Target:   10.0.0.5:22
Status:   Connected

Access URL: https://xxxxx.yishield.com
```

将 Access URL 分享给需要访问的人，在浏览器中打开即可。

## 退出

按 `Ctrl+C` 断开隧道并退出。

## 适用场景

命令行模式最适合：

- 无桌面的 Linux 服务器
- Shell 脚本和自动化流程
- 只需快速建立单个隧道
- CI/CD 环境中临时暴露服务

## 下一步

- [Web UI 模式](./web-ui.md) — 图形化管理多个应用
- [命令参考](../reference/commands.md) — 查看完整参数列表
