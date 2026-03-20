---
title: HTTP/HTTPS 隧道 — 暴露本地 Web 应用
description: 通过 Shield CLI 将本地或内网 Web 应用暴露到公网，完整代理 HTTP 请求，保留 Headers、Cookies 和 WebSocket。
head:
  - - meta
    - name: keywords
      content: HTTP隧道, HTTPS隧道, 暴露本地Web应用, 反向代理, localhost隧道, Shield CLI HTTP
---

# HTTP / HTTPS

通过 Shield CLI 将本地或内网的 Web 应用暴露到公网。

## 快速连接

```bash
# HTTP - 本地 80 端口
shield http

# 暴露本地开发服务器
shield http 3000

# 暴露内网 Web 应用
shield http 10.0.0.5:8080

# HTTPS
shield https
shield https 10.0.0.5:443
```

## 使用场景

### 本地开发预览

将本地开发中的 Web 应用分享给同事或客户预览：

```bash
# React / Vue / Next.js 开发服务器
shield http 3000

# Python Flask / Django
shield http 5000

# 任意本地端口
shield http 8080
```

### 内网应用访问

让外部人员临时访问内网部署的管理后台、监控面板等：

```bash
shield http 10.0.0.5:8080
```

### HTTPS 服务

如果目标服务使用 HTTPS：

```bash
shield https 10.0.0.5:443
```

## 特性

- 完整代理 HTTP 请求，保留原始 Headers 和 Cookies
- 支持 WebSocket
- 自动分配公网 HTTPS 访问地址

## 默认端口

| 输入 | 解析为 |
|---|---|
| `shield http` | `127.0.0.1:80` |
| `shield http 3000` | `127.0.0.1:3000` |
| `shield https` | `127.0.0.1:443` |
