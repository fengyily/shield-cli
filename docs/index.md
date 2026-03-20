---
layout: home
title: Shield CLI — 安全隧道连接器
description: Shield CLI 是一个安全隧道连接器，一条命令将内网 SSH、RDP、VNC、HTTP 服务暴露到浏览器，无需安装客户端。
head:
  - - meta
    - name: keywords
      content: Shield CLI, 安全隧道, 内网穿透, SSH浏览器, RDP浏览器, VNC浏览器, 远程访问, ngrok替代, frp替代
hero:
  name: Shield CLI
  text: 安全隧道连接器
  tagline: 一条命令，一个 URL，浏览器即可访问内网资源
  image:
    src: /logo.svg
    alt: Shield CLI
  actions:
    - theme: brand
      text: 快速入门
      link: /zh/guide/what-is-shield
    - theme: alt
      text: English
      link: /en/guide/what-is-shield
    - theme: alt
      text: GitHub
      link: https://github.com/fengyily/shield-cli
features:
  - icon: 🌐
    title: 浏览器即终端
    details: SSH、RDP、VNC、Web 应用，全部通过浏览器 HTML5 直接访问，无需安装客户端
  - icon: 🔒
    title: 端到端加密
    details: 基于 WebSocket 的加密隧道，AES-256-GCM 凭证加密，机器指纹身份绑定
  - icon: ⚡
    title: 一条命令连接
    details: 智能地址解析，shield ssh 即可连接，自动分配公网访问地址
  - icon: 🖥️
    title: Web 管理面板
    details: 本地 Web UI 管理多达 10 个应用，一键连接/断开，实时状态监控
  - icon: 📦
    title: 全平台支持
    details: macOS / Linux / Windows，amd64 / arm64，Homebrew / Scoop / deb / rpm
  - icon: 🔌
    title: 六大协议
    details: SSH、RDP、VNC、HTTP、HTTPS、Telnet，覆盖主流远程访问场景
---
