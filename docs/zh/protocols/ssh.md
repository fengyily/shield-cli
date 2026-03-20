---
title: SSH 隧道 — 浏览器内 SSH 终端
description: 通过 Shield CLI 在浏览器中打开完整的 SSH 终端，支持密码认证、私钥认证和 SFTP 文件传输，基于 xterm.js。
head:
  - - meta
    - name: keywords
      content: SSH隧道, SSH浏览器, Web终端, xterm.js, SFTP, Shield CLI SSH, 远程SSH
---

# SSH

通过 Shield CLI 在浏览器中打开完整的 SSH 终端，支持密码认证和私钥认证。

## 快速连接

```bash
# 连接本机
shield ssh

# 连接指定服务器
shield ssh 10.0.0.5

# 指定端口
shield ssh 10.0.0.5:2222
```

## 认证方式

### 密码认证

```bash
shield ssh 10.0.0.5 --username root --auth-pass mypassword
```

不传参数时会交互式提示输入。

### 私钥认证

```bash
shield ssh 10.0.0.5 --username root --private-key ~/.ssh/id_rsa
```

### 加密私钥

```bash
shield ssh 10.0.0.5 --username root --private-key ~/.ssh/id_rsa --passphrase mypass
```

## SFTP 文件传输

启用 SFTP 支持，可以在浏览器中进行文件上传和下载：

```bash
shield ssh 10.0.0.5 --enable-sftp
```

## 浏览器终端

连接成功后，浏览器中会呈现基于 xterm.js 的完整终端：

- 支持完整的终端交互（vim、top 等）
- 支持复制粘贴
- 自适应窗口大小

## 默认端口

| 输入 | 解析为 |
|---|---|
| `shield ssh` | `127.0.0.1:22` |
| `shield ssh 2222` | `127.0.0.1:2222` |
| `shield ssh 10.0.0.5` | `10.0.0.5:22` |
