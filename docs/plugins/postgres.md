---
title: PostgreSQL 插件 — 浏览器内数据库管理
description: Shield CLI PostgreSQL 插件提供浏览器内的 Web SQL Client，支持 Schema 浏览、表结构管理、ER 图、SQL 执行、结果排序、CSV 导出等功能。
head:
  - - meta
    - name: keywords
      content: Shield CLI PostgreSQL, PostgreSQL Web客户端, 数据库管理, Web SQL, PG浏览器, ER图, shield postgres
---

# PostgreSQL 插件

通过 Shield CLI 在浏览器中打开完整的 PostgreSQL 数据库管理界面。

## 安装

```bash
shield plugin add postgres
```

验证安装：

```bash
shield plugin list
# NAME      VERSION  PROTOCOLS     INSTALLED
# postgres  v0.3.0   postgres, pg  2026-03-26T10:00:00+08:00
```

## 快速连接

```bash
# 连接本机 PostgreSQL（默认端口 5432）
shield postgres

# 指定端口
shield postgres 5433

# 连接远程服务器
shield postgres 10.0.0.5

# 完整地址
shield postgres 10.0.0.5:5433
```

## 认证方式

### 交互式输入（推荐）

不传凭证参数时，会提示输入用户名、密码和数据库名：

```bash
shield postgres 127.0.0.1:5432

  🔐 Database credentials (press Enter to skip)

  Username [postgres]: postgres
  Password: ****
  Database (optional): mydb

  ✓ Connecting as postgres
    Database: mydb
```

- **Username** 默认 `postgres`，直接回车使用默认值
- **Password** 隐藏输入，不会回显
- **Database** 可选，回车跳过（默认连接 `postgres` 库）

### 命令行传参

```bash
# 使用数据库专用参数
shield postgres 127.0.0.1:5432 --db-user postgres --db-pass mypassword --db-name mydb

# 也兼容通用认证参数
shield postgres 127.0.0.1:5432 --username postgres --auth-pass mypassword
```

| 参数 | 别名 | 说明 |
|---|---|---|
| `--db-user` | `--username` | 数据库用户名 |
| `--db-pass` | `--auth-pass` | 数据库密码 |
| `--db-name` | — | 数据库名（可选） |
| `--readonly` | — | 强制只读模式，禁止写操作 |

## Web 管理界面

连接成功后，浏览器自动打开 Web SQL Client：

### 功能概览

| 功能 | 说明 |
|---|---|
| Schema 浏览 | 左侧栏列出所有 Schema，展开查看表 → 字段 → 索引 |
| 表结构管理 | 创建/删除表、新增/编辑/删除字段，可视化构建器 |
| ER 图 | 交互式实体关系图，支持拖拽建立/删除外键、表/字段 CRUD |
| SQL 执行 | 多标签页，支持 SELECT、SHOW、EXPLAIN 等查询 |
| 结果排序 | 点击列头升序/降序排列 |
| CSV 导出 | 一键导出查询结果为 CSV 文件 |
| 复制结果 | 复制为 Tab 分隔文本，可直接粘贴到 Excel |
| 只读模式 | 启动参数决定，前后端双重拦截写操作 |

### ER 图

左侧栏每个 Schema 旁的 ◥ 按钮可快捷打开 ER 图。

**操作概要：**

| 操作 | 方式 |
|---|---|
| 创建外键 | 从字段 A 拖动到表 B，自动匹配类型/名称 |
| 删除外键 | 点击关系线选中（变红），按 Delete 键 |
| 编辑表结构 | 点击表头齿轮图标 ⚙ |
| 字段操作 | 鼠标悬停字段行，点击齿轮 → 编辑/删除/新增 |
| 右键菜单 | 表头右键 → 重命名/删除表、新增字段；空白区右键 → 新建表 |
| 布局 | 支持 Grid / Horizontal / Vertical / Center 四种布局模式 |
| 缩放 | Ctrl+滚轮缩放，滚轮平移，Fit 按钮自适应 |
| 导出 | SVG 一键导出 |

所有 DDL 操作均有 SQL 预览弹窗，确认后执行。只读模式下所有写操作被隐藏和拦截。

详细交互逻辑参见 [ER 图技术文档](https://github.com/fengyily/shield-plugins/blob/main/shield-plugin-postgres/docs/er-foreign-key.md)。

### 只读模式

只读/读写模式完全由启动参数决定，页面上仅展示当前状态，远程用户无法更改。

**CLI 模式**：通过 `--readonly` 参数控制：

```bash
# 只读模式（推荐用于分享场景）
shield postgres 10.0.0.5:5432 --db-user postgres --readonly

# 读写模式（默认）
shield postgres 10.0.0.5:5432 --db-user postgres
```

**Web UI 模式**：在添加/编辑应用时，勾选 **Read-Only Mode** 复选框。

页面右上角显示当前模式：

- **🔒 Read-Only**（橙色标签）— 写操作被前后端双重拦截
- **🔓 Read-Write**（绿色标签）— 允许所有操作

只读模式下以下语句会被阻止：

```
INSERT, UPDATE, DELETE, DROP, ALTER, CREATE,
TRUNCATE, GRANT, REVOKE
```

::: tip 安全建议
在公网暴露数据库管理界面时，建议开启只读模式，并配合 `--invisible` 隐身模式使用：
```bash
shield postgres 127.0.0.1:5432 --db-user readonly_user --readonly --invisible
```
:::

### 快捷键

| 快捷键 | 功能 |
|---|---|
| `Ctrl+Enter` / `Cmd+Enter` | 执行当前 SQL |
| `Tab` | 插入两个空格（不会切换焦点） |
| `Delete` / `Backspace` | 删除选中的 ER 关系线 |
| `Escape` | 关闭 ER 图 |
| `Ctrl+滚轮` | ER 图缩放 |

## Docker

无需 Shield CLI，直接以 Docker 容器方式运行独立 Web UI：

```bash
docker run -d --name shield-postgres \
  -e DB_HOST=10.0.0.20 \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASS=mypass \
  -e DB_NAME=mydb \
  -p 8080:8080 \
  fengyily/shield-postgres
```

打开 http://localhost:8080 — pgAdmin 的轻量替代（~9 MB vs ~400 MB）。

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `DB_HOST` | `127.0.0.1` | 数据库地址 |
| `DB_PORT` | `5432` | 数据库端口 |
| `DB_USER` | `postgres` | 数据库用户 |
| `DB_PASS` | — | 数据库密码 |
| `DB_NAME` | `postgres` | 默认数据库 |
| `DB_READONLY` | `false` | 只读模式 |
| `WEB_PORT` | `8080` | Web UI 端口 |

## 默认端口

| 输入 | 解析为 |
|---|---|
| `shield postgres` | `127.0.0.1:5432` |
| `shield postgres 5433` | `127.0.0.1:5433` |
| `shield postgres 10.0.0.5` | `10.0.0.5:5432` |
| `shield postgres 10.0.0.5:5433` | `10.0.0.5:5433` |

`pg` 是 `postgres` 的别名，行为完全一致：

```bash
shield pg 127.0.0.1:5432 --db-user postgres
```

## 使用场景

### 场景一：远程调试生产数据库

```bash
shield postgres 10.0.0.5:5432 --db-user readonly --db-pass xxx --readonly --invisible
```

分享 Auth URL 给需要查看数据的同事，无需安装任何客户端。通过 ER 图快速了解表关系。

### 场景二：Docker 中的 PostgreSQL

```bash
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:16

shield postgres 127.0.0.1:5432 --db-user postgres --db-pass postgres
```

### 场景三：Schema 可视化与数据建模

以读写模式连接开发库，通过 ER 图拖拽创建外键、右键新建表和字段，实时预览 DDL：

```bash
shield postgres 127.0.0.1:5432 --db-user dev --db-pass xxx --db-name dev_db
```

### 场景四：内网数据库审计

```bash
shield postgres 192.168.1.100:5432 --db-user auditor --db-pass xxx --readonly --invisible
```

- 只读模式，无法执行写操作
- 隐身模式需要授权码才能访问
- 查询结果可导出 CSV 用于报告

## API 接口

PostgreSQL 插件在本地启动一个 HTTP 服务，提供以下 API：

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/info` | 服务器信息（版本、主机、用户、只读状态） |
| GET | `/api/schemas` | Schema 列表 |
| GET | `/api/tables?schema=public` | 表列表 |
| GET | `/api/columns?schema=public&table=users` | 表字段信息 |
| GET | `/api/indexes?schema=public&table=users` | 表索引信息 |
| GET | `/api/er?schema=public` | ER 图数据（表结构 + 外键关系） |
| POST | `/api/query` | 执行 SQL 查询 |

## 故障排除

### 连接被拒绝

```
plugin error: cannot connect to PostgreSQL at 127.0.0.1:5432: connection refused
```

检查 PostgreSQL 服务是否正在运行：

```bash
# macOS
brew services list | grep postgresql

# Linux
systemctl status postgresql

# Docker
docker ps | grep postgres
```

### 认证失败

```
plugin error: password authentication failed for user "postgres"
```

- 确认用户名和密码正确
- 检查 `pg_hba.conf` 是否允许连接来源
- Docker 环境下注意 `listen_addresses` 配置

### 插件未安装

```
unsupported protocol "postgres"
```

先安装插件：

```bash
shield plugin add postgres
```

## 下一步

- [插件系统概览](/plugins/)
- [MySQL 插件](/plugins/mysql)
- [插件开发指南](/plugins/development)
