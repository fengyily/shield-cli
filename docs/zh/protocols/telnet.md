---
title: Telnet 隧道 — 浏览器访问网络设备
description: 通过 Shield CLI 在浏览器中连接 Telnet 服务，适用于路由器、交换机等网络设备和传统系统管理。
head:
  - - meta
    - name: keywords
      content: Telnet隧道, Telnet浏览器, 网络设备管理, 传统系统, Shield CLI Telnet
---

# Telnet

通过 Shield CLI 在浏览器中连接 Telnet 服务，适用于网络设备和传统系统管理。

## 快速连接

```bash
# 连接本机
shield telnet

# 连接指定设备
shield telnet 10.0.0.1

# 指定端口
shield telnet 10.0.0.1:2323
```

## 典型场景

- 管理路由器、交换机等网络设备
- 连接传统工控系统
- 访问不支持 SSH 的遗留服务器

## 默认端口

| 输入 | 解析为 |
|---|---|
| `shield telnet` | `127.0.0.1:23` |
| `shield telnet 2323` | `127.0.0.1:2323` |
| `shield telnet 10.0.0.1` | `10.0.0.1:23` |
