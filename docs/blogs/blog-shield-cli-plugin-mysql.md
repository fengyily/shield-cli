---
title: Shield CLI v0.3.0：插件系统上线，首发 MySQL Web 管理
description: Shield CLI v0.3.0 新增插件系统，首个插件支持在浏览器中管理 MySQL 数据库。本文介绍如何通过 Web UI 和命令行安装、配置和使用 MySQL 插件。
date: 2026-03-24
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: Shield CLI v0.3.0, MySQL Web管理, 浏览器数据库管理, shield mysql, 插件系统, 远程MySQL管理, Web UI
---

# Shield CLI v0.3.0：插件系统上线，首发 MySQL Web 管理

> Shield CLI v0.3.0 加了一个插件机制，可以通过独立插件扩展协议支持。第一个插件是 MySQL Web 客户端——装完之后，在 Web 管理面板里添加一个 MySQL 应用，点连接就能在浏览器里管理数据库。这篇文章介绍下怎么用，以及插件系统大致是怎么回事。

---

## 先看效果

![Shield MySQL 插件演示](/demo/demo-mysql.gif)

整个流程三步：

1. 安装 MySQL 插件
2. 在 Web 管理面板添加一个 MySQL 应用，填好地址和凭证
3. 点击连接，浏览器自动打开 Web SQL 客户端

对方只需要一个浏览器，不用装 Navicat、DBeaver 或任何 MySQL 客户端。

---

## 安装插件

前提是已经装好了 Shield CLI（[安装指南](https://docs.yishield.com/guide/quickstart)），然后装 MySQL 插件：

```bash
shield plugin add mysql
```

验证一下：

```bash
shield plugin list
# NAME   VERSION  PROTOCOLS       INSTALLED
# mysql  v0.1.0   mysql, mariadb  2026-03-24T10:00:00+08:00
```

插件是一个独立的二进制文件，安装在 `~/.shield-cli/plugins/` 下，不影响主程序。

---

## 通过 Web UI 使用（推荐）

### 1. 启动管理面板

```bash
shield start
```

浏览器打开 `http://localhost:8181`。

### 2. 添加 MySQL 应用

点击 **Add App**，填写配置：

| 字段 | 说明 | 示例 |
|---|---|---|
| Protocol | 选择 `mysql` | mysql |
| Target IP | 数据库地址 | 10.0.0.5 |
| Port | 端口 | 3306 |
| DB User | 数据库用户名 | root |
| DB Password | 数据库密码 | mypassword |
| DB Name | 数据库名（可选） | mydb |
| Read-Only Mode | 是否只读 | ✅ 勾选 |

**Read-Only Mode** 默认勾选。勾选后，远程访问的人无法执行 INSERT、UPDATE、DELETE、DROP 等写操作——控制权在你手上，页面端只能看到状态，不能更改。

如果需要开放写权限，取消勾选即可。

### 3. 连接

点击应用卡片上的连接按钮，等待隧道建立。连接成功后浏览器自动打开 Web SQL 客户端。

### 4. 编辑和管理

点击应用卡片的编辑按钮可以随时修改配置，包括切换只读/读写模式。修改后下次连接生效。

---

## 通过命令行使用

如果偏好命令行，也可以直接用 `shield mysql` 命令：

```bash
# 连接本机 MySQL（默认 127.0.0.1:3306）
shield mysql

# 连接远程 MySQL
shield mysql 10.0.0.5:3306 --db-user root

# 只读模式
shield mysql 10.0.0.5:3306 --db-user root --readonly
```

不传凭证参数时会交互式提示输入：

```bash
shield mysql 10.0.0.5:3306

  🔐 Database credentials (press Enter to skip)

  Username [root]: readonly
  Password: ****
  Database (optional): orders

  ✓ Connecting as readonly
    Database: orders
```

### 命令行参数

| 参数 | 说明 |
|---|---|
| `--db-user` | 数据库用户名 |
| `--db-pass` | 数据库密码 |
| `--db-name` | 数据库名（可选，连接后也能切换） |
| `--readonly` | 强制只读模式 |

MariaDB 用法完全一样，`mariadb` 是 `mysql` 的别名：

```bash
shield mariadb 10.0.0.5:3306 --db-user root
```

---

## Web 界面能做什么

连接成功后浏览器自动打开，界面功能不算多，但日常查数据够用：

- **数据库浏览** — 左侧栏列出所有数据库，点击切换
- **表列表** — 支持搜索过滤，超过 50 个表自动分页
- **表结构查看** — 单击表名查看字段、类型、索引
- **SQL 执行** — 支持 SELECT、SHOW、DESCRIBE 等查询，`Ctrl+Enter` 快捷执行
- **结果操作** — 点击列头排序、导出 CSV、一键复制（可直接粘贴到 Excel）
- **快速查询** — 双击表名自动执行 `SELECT * FROM ... LIMIT 100`

### 关于只读模式

只读/读写状态完全由启动方决定：

- **Web UI**：添加或编辑应用时勾选 **Read-Only Mode**
- **CLI**：加 `--readonly` 参数

页面右上角显示当前模式标识，远程用户只能看到，不能切换：

- **🔒 Read-Only**（橙色）— 写操作被前后端双重拦截
- **🔓 Read-Write**（绿色）— 允许所有操作

只读模式下以下语句会被阻止：

```
INSERT, UPDATE, DELETE, DROP, ALTER, CREATE,
TRUNCATE, RENAME, REPLACE, GRANT, REVOKE
```

对于分享给外部人员的场景，建议只读模式 + 只读数据库账户双重保险，再配合 `--invisible` 隐身模式：

```bash
shield mysql 10.0.0.5:3306 --db-user readonly_user --readonly --invisible
```

---

## 几个实际场景

### 临时给同事查数据

内网有台 MySQL，同事需要查几条数据但没有数据库客户端。在 Web 管理面板添加一个 MySQL 应用，勾选只读模式，连接后把链接发给他。用完断开就行。

或者命令行：

```bash
shield mysql 192.168.1.100:3306 --db-user readonly --db-pass xxx --readonly
```

### Docker 里的 MySQL

本地 Docker 跑了个 MySQL 实例，想临时分享给远程的人看：

```bash
docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:8
```

在 Web 管理面板添加应用：Protocol 选 `mysql`，IP 填 `127.0.0.1`，端口 `3306`，DB User 填 `root`，DB Password 填 `root`，勾选 Read-Only Mode，点连接。

### 内网数据库审计

以只读模式连接，让审计人员可以通过浏览器查看数据库。在 Web 管理面板添加应用时勾选 Read-Only Mode 和 Invisible Mode，分享 Auth URL 给审计人员。

- 只读模式，无法执行写操作
- 隐身模式需要授权码才能访问
- 查询结果可导出 CSV 用于报告

---

## 插件系统简介

这次做 MySQL 支持没有直接写进主程序，而是通过插件实现。主要是不想让主程序越来越臃肿——不是所有人都需要数据库管理功能，没必要让只用 SSH 的用户也下载这部分代码。

简单说下原理：插件是一个独立的可执行文件，Shield 启动它之后通过 stdin/stdout 交换 JSON 消息。插件在本地启动一个 Web 服务，告诉 Shield 端口号，Shield 再把这个端口通过隧道暴露出去。对 Shield 服务端来说，插件的 Web 界面就是一个普通的 HTTP 应用。

```
Web UI 点击连接 / shield mysql 10.0.0.5:3306
    ↓
Shield 启动 mysql 插件进程，传入连接信息和只读配置
    ↓
插件连接数据库，启动 Web 界面（本地随机端口）
    ↓
Shield 通过加密隧道将 Web 界面暴露到公网
    ↓
浏览器打开，直接使用
```

插件管理命令：

```bash
# 查看已安装插件
shield plugin list

# 安装
shield plugin add <name>

# 从本地二进制安装（开发调试用）
shield plugin add <name> --from ./path/to/binary

# 卸载
shield plugin remove <name>
```

如果有兴趣开发自己的插件，可以参考[插件开发指南](https://docs.yishield.com/plugins/development)，协议很简单，本质上就是写一个 HTTP 服务。

---

## 后续计划

MySQL 是第一个插件，后面打算陆续支持 PostgreSQL 和 SQL Server。不过精力有限，进度不好说，先把 MySQL 这个做稳。

| 插件 | 状态 |
|---|---|
| mysql | 已发布 |
| postgres | 开发中 |
| sqlserver | 计划中 |

如果你有想支持的数据库或服务类型，欢迎提 issue 聊聊。

---

## 不足

坦诚说几个目前的问题：

1. **Web 界面功能有限** — 和 Navicat、DBeaver 这类专业工具没法比，目前只支持基础的查询和浏览，没有可视化建表、数据编辑等高级功能
2. **只读模式不能替代数据库权限** — 只读模式在前后端双重拦截写操作，远程用户无法绕过，但安全敏感场景仍建议配合只读数据库账户
3. **插件生态刚起步** — 目前只有 MySQL 一个插件，还谈不上"生态"
4. **服务端不开源** — 数据经过官方网关，介意的话可以用 `shield tcp 3306` 走纯端口转发，自己选客户端连

---

## 试一下

```bash
# 装 Shield CLI
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh

# 装 MySQL 插件
shield plugin add mysql

# Web UI 方式（推荐）
shield start
# 然后在 http://localhost:8181 添加 MySQL 应用

# 命令行方式
shield mysql 127.0.0.1:3306
```

文档：https://docs.yishield.com/plugins/mysql

开源地址：https://github.com/fengyily/shield-cli
