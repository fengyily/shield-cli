---
title: 常见错误 — Shield CLI 故障排查
description: Shield CLI 常见错误的解决方案，包括连接超时、认证失败、端口冲突和安装问题。
head:
  - - meta
    - name: keywords
      content: Shield CLI 错误, 故障排查, 连接超时, 认证失败, 端口占用, 安装错误
---

# 常见错误

## 连接相关

### 连接超时

**症状：** 长时间停留在「连接中」状态

**可能原因：**
- 目标服务未启动或端口不正确
- 本机网络无法连接到公网网关
- 防火墙阻止了 WebSocket 连接

**解决方法：**
1. 确认目标服务正在运行：`telnet <ip> <port>`
2. 检查是否能访问 `console.yishield.com`
3. 检查防火墙是否允许出站 62888 端口

### 认证失败 (401)

**症状：** 提示认证失败或 401 错误

**可能原因：**
- 本地凭证已过期或损坏

**解决方法：**
```bash
shield clean
```
清除凭证后重新连接，会自动生成新的凭证。

### 目标服务认证失败

**症状：** 隧道建立成功，但浏览器中显示认证错误

**可能原因：**
- 用户名或密码错误
- SSH 私钥不匹配

**解决方法：**
- 确认用户名和密码正确
- 如果使用私钥，确认私钥文件路径正确且权限为 `600`

### 端口被占用

**症状：** `shield start` 启动失败

**可能原因：**
- 8181 端口已被其他程序占用

**解决方法：**
```bash
shield start 9090  # 使用其他端口
```

## 安装相关

### Homebrew 安装失败

**解决方法：**
```bash
brew update
brew tap fengyily/tap
brew install shield-cli
```

如果仍然失败，尝试一键安装脚本：
```bash
curl -fsSL https://raw.githubusercontent.com/fengyily/shield-cli/main/install.sh | sh
```

### 权限不足

**症状：** Linux 下运行报 Permission denied

**解决方法：**
```bash
chmod +x ./shield
```
