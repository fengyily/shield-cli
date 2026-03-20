---
title: 命令参考 — Shield CLI 完整参数指南
description: Shield CLI 所有命令、参数、地址格式和示例的完整参考。涵盖 shield start、ssh、rdp、vnc、http、https、telnet 和 clean 命令。
head:
  - - meta
    - name: keywords
      content: Shield CLI 命令, CLI参考, shield ssh, shield rdp, shield start, 命令参数
---

# 命令参考

## 命令一览

| 命令 | 说明 |
|---|---|
| `shield start [port]` | 启动 Web 管理面板（默认端口 8181） |
| `shield ssh [address]` | 创建 SSH 隧道 |
| `shield rdp [address]` | 创建 RDP 隧道 |
| `shield vnc [address]` | 创建 VNC 隧道 |
| `shield http [address]` | 创建 HTTP 隧道 |
| `shield https [address]` | 创建 HTTPS 隧道 |
| `shield telnet [address]` | 创建 Telnet 隧道 |
| `shield install [--port]` | 安装为系统服务（开机自启） |
| `shield uninstall` | 卸载系统服务 |
| `shield clean` | 清除本地凭证缓存 |

## 地址格式

`[address]` 支持以下格式：

| 格式 | 示例 | 说明 |
|---|---|---|
| 省略 | `shield ssh` | 使用 `127.0.0.1` + 协议默认端口 |
| 仅端口 | `shield ssh 2222` | 使用 `127.0.0.1` + 指定端口 |
| 仅 IP | `shield ssh 10.0.0.5` | 使用指定 IP + 协议默认端口 |
| 完整地址 | `shield ssh 10.0.0.5:2222` | 使用指定 IP + 指定端口 |

### 协议默认端口

| 协议 | 默认端口 |
|---|---|
| SSH | 22 |
| RDP | 3389 |
| VNC | 5900 |
| HTTP | 80 |
| HTTPS | 443 |
| Telnet | 23 |

## 全局参数

| 参数 | 说明 | 示例 |
|---|---|---|
| `--username` | 目标服务用户名 | `--username root` |
| `--auth-pass` | 目标服务密码 | `--auth-pass mypass` |
| `--server` | 自定义服务端地址 | `--server https://my.server/raas` |

## SSH 专用参数

| 参数 | 说明 | 示例 |
|---|---|---|
| `--private-key` | SSH 私钥文件路径 | `--private-key ~/.ssh/id_rsa` |
| `--passphrase` | 私钥密码 | `--passphrase mypass` |
| `--enable-sftp` | 启用 SFTP 文件传输 | `--enable-sftp` |

## 访问模式参数

| 参数 | 说明 | 示例 |
|---|---|---|
| `--visable` | 可见模式（默认） | `--visable` |
| `--visable=<节点>` | 可见模式，指定接入节点 | `--visable=HK` |
| `--invisible` | 隐身模式，需授权码访问 | `--invisible` |

## 服务管理

| 命令 | 说明 |
|---|---|
| `shield install` | 安装为系统服务，使用默认端口 8181 |
| `shield install --port 8182` | 指定端口安装 |
| `shield uninstall` | 卸载系统服务 |

### Install 参数

| 参数 | 默认值 | 说明 |
|---|---|---|
| `--port` | `8181` | Web UI 端口号 |

安装命令会自动检测端口冲突并建议可用的替代端口。详见[系统服务安装](/zh/guide/system-service)了解各平台详情。

## 示例

```bash
# 最简用法
shield ssh

# 完整参数
shield ssh 10.0.0.5:2222 --username root --auth-pass mypass --enable-sftp

# Web UI 模式
shield start
shield start 9090

# 隐身模式连接 RDP
shield rdp 10.0.0.5 --username Administrator --invisible

# 清除缓存
shield clean

# 安装为系统服务
shield install
shield install --port 8182

# 卸载服务
shield uninstall
```
