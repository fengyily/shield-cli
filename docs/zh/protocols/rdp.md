---
title: RDP 隧道 — 浏览器内 Windows 远程桌面
description: 通过 Shield CLI 在浏览器中访问 Windows 远程桌面（RDP），完整鼠标键盘控制，无需安装 RDP 客户端。
head:
  - - meta
    - name: keywords
      content: RDP隧道, RDP浏览器, Windows远程桌面, 远程桌面浏览器, Shield CLI RDP
---

# RDP

通过 Shield CLI 在浏览器中访问 Windows 远程桌面，无需安装 RDP 客户端。

## 快速连接

```bash
# 连接本机 Windows
shield rdp

# 连接指定 IP
shield rdp 10.0.0.5

# 指定端口
shield rdp 10.0.0.5:3390
```

## 认证

```bash
shield rdp 10.0.0.5 --username Administrator --auth-pass mypassword
```

不传参数时会交互式提示输入 Windows 登录凭证。

## 浏览器体验

连接成功后，浏览器中会呈现完整的 Windows 桌面：

- 完整的鼠标和键盘控制
- 屏幕自适应浏览器窗口
- 支持从手机、平板等任意设备访问

## 典型场景

- 远程操作办公室的 Windows 电脑
- 访问内网部署的 Windows 服务器
- 让不在现场的同事临时使用某台 Windows 机器

## 默认端口

| 输入 | 解析为 |
|---|---|
| `shield rdp` | `127.0.0.1:3389` |
| `shield rdp 3390` | `127.0.0.1:3390` |
| `shield rdp 10.0.0.5` | `10.0.0.5:3389` |
