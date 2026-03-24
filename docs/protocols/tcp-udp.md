---
title: TCP/UDP 端口代理 — 转发任意 TCP/UDP 服务
description: 通过 Shield CLI 创建 TCP/UDP 端口代理隧道，转发 MySQL、Redis、PostgreSQL、DNS 等任意服务，无需打开浏览器。
head:
  - - meta
    - name: keywords
      content: TCP代理, UDP代理, 端口转发, MySQL隧道, Redis隧道, DNS隧道, Shield CLI TCP UDP
---

# TCP / UDP 端口代理

通过 Shield CLI 转发任意 TCP/UDP 端口，适用于数据库、缓存、DNS 等不需要浏览器渲染的服务。

## 与其他协议的区别

| | SSH/RDP/VNC/HTTP | TCP/UDP |
|---|---|---|
| 浏览器渲染 | 是（HTML5） | 否 |
| 自动打开浏览器 | 是 | 否 |
| 连接方式 | 浏览器访问 URL | 客户端工具连接域名:端口 |
| 默认端口 | 有 | 无（必须指定） |

## 快速连接

```bash
# TCP 端口代理
shield tcp 3306                      # MySQL (127.0.0.1:3306)
shield tcp 6379                      # Redis (127.0.0.1:6379)
shield tcp 5432                      # PostgreSQL (127.0.0.1:5432)
shield tcp 192.168.1.10:3306         # 远程 MySQL

# UDP 端口代理
shield udp 53                        # DNS (127.0.0.1:53)
shield udp 514                       # Syslog (127.0.0.1:514)
shield udp 192.168.1.1:161           # SNMP
```

::: warning 必须指定端口
TCP/UDP 没有默认端口，省略端口会报错：
```
Error: port is required for tcp protocol

Usage: shield tcp <port> or shield tcp <ip:port>
```
:::

## 连接指南

隧道建立后，Shield CLI 会打印连接信息（不会自动打开浏览器）：

```
  📡 Connection Guide (TCP port proxy):
    your-app.cn01.apps.yishield.com:58379  →  127.0.0.1:3306

    Examples:
      telnet your-app.cn01.apps.yishield.com 58379
      mysql -h your-app.cn01.apps.yishield.com -P 58379 -u root
      redis-cli -h your-app.cn01.apps.yishield.com -p 58379
```

使用打印的**专属域名和端口**，通过对应的客户端工具连接即可。

## 典型场景

- **数据库访问** — 远程连接 MySQL、PostgreSQL、MongoDB、Redis
- **DNS 转发** — 代理内网 DNS 服务器
- **消息队列** — 连接 RabbitMQ、Kafka、NATS
- **自定义服务** — 任意监听在 TCP/UDP 端口的服务

## 地址格式

| 输入 | 解析为 |
|---|---|
| `shield tcp 3306` | `127.0.0.1:3306` |
| `shield tcp 192.168.1.10:3306` | `192.168.1.10:3306` |
| `shield udp 53` | `127.0.0.1:53` |
| `shield udp 192.168.1.1:161` | `192.168.1.1:161` |
