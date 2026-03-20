---
title: 系统服务安装 — 开机自动启动 Shield CLI
description: 将 Shield CLI 安装为系统服务，随操作系统自动启动。支持 macOS (launchd)、Linux (systemd) 和 Windows 服务。
head:
  - - meta
    - name: keywords
      content: Shield CLI 服务, 安装服务, 系统服务, launchd, systemd, Windows 服务, 开机自启, 守护进程
---

# 系统服务安装

Shield CLI 可以安装为系统服务，在电脑启动时自动运行。适用于需要持续访问内网服务的场景。

## 快速开始

```bash
# 使用默认端口 (8181) 安装
shield install

# 指定端口安装
shield install --port 8182
```

安装完成后，Web UI 将在 `http://localhost:8181`（或你指定的端口）上持续运行，重启后自动恢复。

在 **macOS** 和 **Windows** 上，系统托盘会自动显示 Shield 图标 — 点击即可在浏览器中打开 Dashboard。

## 命令

### `shield install`

将 Shield CLI 安装为系统服务。

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--port` | `8181` | Web UI 端口号 |

**执行流程：**
1. 检查是否已安装过服务
2. 检测指定端口是否被占用
3. 向操作系统注册服务
4. 立即启动服务

### `shield uninstall`

卸载 Shield CLI 系统服务。

```bash
shield uninstall
```

**执行流程：**
1. 停止正在运行的服务
2. 从自动启动中移除
3. 配置和凭证**不会被删除**

## 端口配置

默认使用 **8181** 端口。如果端口被占用：

```bash
# Shield 会自动检测冲突并建议可用端口
$ shield install
Error: port 8181 is already in use.
Try an available port: shield install --port 8182
```

你可以指定任意可用端口：

```bash
shield install --port 9090
```

## 各平台说明

### macOS (launchd)

- **类型：** 用户级 Launch Agent（无需 sudo）
- **配置文件：** `~/Library/LaunchAgents/com.yishield.shield-cli.plist`
- **日志：** `~/.shield-cli/logs/shield-cli.log`
- 用户登录时自动启动
- 异常退出后自动重启 (KeepAlive)
- 系统托盘图标，点击快速打开 Dashboard

### Linux (systemd)

- **类型：** 系统服务（需要 sudo）
- **配置文件：** `/etc/systemd/system/shield-cli.service`
- **日志：** `journalctl -u shield-cli`
- 网络就绪后自动启动
- 异常退出后 5 秒重启

```bash
# 手动管理服务
sudo systemctl status shield-cli
sudo systemctl restart shield-cli
sudo journalctl -u shield-cli -f
```

### Windows

- **类型：** Windows 服务（需要管理员权限）
- **服务名称：** `ShieldCLI`
- 开机自动启动
- 系统托盘图标，点击快速打开 Dashboard
- 可通过服务管理器 (`services.msc`) 或 `sc` 命令管理

```powershell
# 手动管理服务
sc query ShieldCLI
sc stop ShieldCLI
sc start ShieldCLI
```

## 更换端口

要更改已安装服务的端口：

```bash
shield uninstall
shield install --port 9090
```

## 故障排查

### 端口被占用

```bash
# 查看端口占用
# macOS / Linux
lsof -i :8181

# Windows
netstat -ano | findstr :8181
```

### 服务无法启动 (Linux)

```bash
# 查看服务日志
sudo journalctl -u shield-cli --no-pager -n 50

# 确认二进制路径
which shield
```

### 服务无法启动 (macOS)

```bash
# 查看系统日志
log show --predicate 'eventMessage contains "shield"' --last 5m

# 验证 plist 文件
plutil ~/Library/LaunchAgents/com.yishield.shield-cli.plist
```
