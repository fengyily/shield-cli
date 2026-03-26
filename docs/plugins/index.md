---
title: 插件系统 — 按需扩展协议支持
description: Shield CLI 插件系统概览，通过独立插件按需扩展数据库和服务协议支持，包括 MySQL、PostgreSQL、SQL Server 等。
head:
  - - meta
    - name: keywords
      content: Shield CLI 插件, 插件系统, shield plugin, 数据库管理, MySQL插件, PostgreSQL插件, 扩展协议
---

# 插件系统

Shield CLI 通过插件机制扩展协议支持。每个插件是一个独立的二进制文件，按需安装，主程序零膨胀。

## 设计理念

Shield CLI 内置了 SSH、RDP、VNC、HTTP 等通用协议。对于数据库、消息队列等需要专属 Web 管理界面的场景，我们通过插件来支持：

- **按需安装** — 不用的协议不会增加主程序体积
- **独立更新** — 插件版本独立于主程序
- **开放扩展** — 任何人都可以开发和发布 Shield 插件

## 工作原理

```
shield mysql 127.0.0.1:3306 --db-user root
    ↓
Shield 查找已安装的 mysql 插件
    ↓
启动插件进程（独立二进制）
    ↓
插件在本地启动 Web 数据库管理界面
    ↓
Shield 通过加密隧道将 Web 界面暴露到公网
    ↓
用户通过浏览器访问完整的数据库管理平台
```

对 Shield 服务端而言，插件的 Web 界面就是一个普通的 HTTP 应用，无需服务端做任何适配。

## 可用插件

| 插件 | 协议 | 默认端口 | 说明 |
|---|---|---|---|
| [mysql](/plugins/mysql) | `mysql`, `mariadb` | 3306 | MySQL / MariaDB Web 管理客户端 |
| [postgres](/plugins/postgres) | `postgres`, `pg` | 5432 | PostgreSQL Web 管理客户端（含 ER 图） |
| sqlserver | `sqlserver`, `mssql` | 1433 | SQL Server Web 管理客户端（即将发布） |

## 快速开始

```bash
# 安装插件
shield plugin add mysql

# 使用（交互式输入凭证）
shield mysql 127.0.0.1:3306

# 使用（命令行传参）
shield mysql 127.0.0.1:3306 --db-user root --db-pass mypassword --db-name mydb
```

连接成功后，浏览器自动打开 Web 数据库管理界面。

## 管理插件

```bash
# 查看已安装插件
shield plugin list

# 安装插件
shield plugin add <name>

# 从本地二进制安装（开发调试）
shield plugin add <name> --from ./path/to/binary

# 卸载插件
shield plugin remove <name>
```

## 插件通信协议

Shield 主程序与插件之间通过 **stdin/stdout JSON** 通信，协议极其简单：

### 启动请求（Shield → 插件）

```json
{
  "action": "start",
  "config": {
    "host": "127.0.0.1",
    "port": 3306,
    "user": "root",
    "pass": "password",
    "database": "mydb"
  }
}
```

### 就绪响应（插件 → Shield）

```json
{
  "status": "ready",
  "web_port": 19876,
  "name": "MySQL Web Client",
  "version": "0.1.0"
}
```

### 停止请求（Shield → 插件）

```json
{"action": "stop"}
```

## 插件存储位置

| 平台 | 路径 |
|---|---|
| macOS / Linux | `~/.shield-cli/plugins/` |
| Windows | `%LOCALAPPDATA%\ShieldCLI\plugins\` |

```
~/.shield-cli/plugins/
├── registry.json              # 已安装插件索引
├── shield-plugin-mysql        # MySQL 插件二进制
└── shield-plugin-postgres     # PostgreSQL 插件二进制
```

## 开发自定义插件

想为 Shield 开发新插件？请参阅[插件开发指南](/plugins/development)。

## 下一步

- [MySQL 插件详细文档](/plugins/mysql)
- [PostgreSQL 插件详细文档](/plugins/postgres)
- [插件开发指南](/plugins/development)
- [命令参考](/reference/commands)
