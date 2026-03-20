---
title: VNC 隧道 — 浏览器内远程桌面共享
description: 通过 Shield CLI 在浏览器中共享和控制远程桌面屏幕，像素级渲染，完整鼠标键盘映射。
head:
  - - meta
    - name: keywords
      content: VNC隧道, VNC浏览器, 远程桌面共享, 屏幕共享, Shield CLI VNC
---

# VNC

通过 Shield CLI 在浏览器中共享和控制远程桌面屏幕。

## 快速连接

```bash
# 连接本机
shield vnc

# 连接指定 IP
shield vnc 10.0.0.10

# 指定端口
shield vnc 10.0.0.10:5901
```

## 认证

```bash
shield vnc 10.0.0.10 --auth-pass vncpassword
```

VNC 通常只需要密码，不需要用户名。

## 浏览器体验

- 像素级远程桌面渲染
- 鼠标和键盘完整映射
- 适用于 Linux、macOS 和 Windows 上的 VNC 服务器

## 典型场景

- 远程协助：共享屏幕给同事或客户
- 管理无头 Linux 服务器的图形界面
- 远程操作实验室或工厂设备

## 默认端口

| 输入 | 解析为 |
|---|---|
| `shield vnc` | `127.0.0.1:5900` |
| `shield vnc 5901` | `127.0.0.1:5901` |
| `shield vnc 10.0.0.10` | `10.0.0.10:5900` |
