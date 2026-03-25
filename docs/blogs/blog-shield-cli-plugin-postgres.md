---
title: Shield CLI 新增 PostgreSQL 插件：浏览器里管理 PG 数据库
description: Shield CLI 发布 PostgreSQL Web 客户端插件，支持 Schema 浏览、SQL 查询、表结构管理、行级编辑，一条命令安装，浏览器直接用。
date: 2026-03-25
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: Shield CLI, PostgreSQL Web客户端, 浏览器管理PostgreSQL, shield postgres, PG数据库管理, Web SQL客户端, 远程数据库
---

# Shield CLI 新增 PostgreSQL 插件：浏览器里管理 PG 数据库

> 上周发了 MySQL 插件，这周 PostgreSQL 跟上了。同样的思路——一条命令装好，浏览器打开就能用，不需要在对方机器上装任何客户端。

---

## 一条命令安装

```bash
shield plugin add postgres
```

验证：

```bash
shield plugin list
# NAME      VERSION  PROTOCOLS              INSTALLED
# postgres  v0.1.0   postgres, pg, postgresql  2026-03-25T10:00:00+08:00
```

前提是主程序升级到 v0.3.3+：

```bash
brew update && brew upgrade shield-cli
# 或
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh
```

---

## 使用方式

### Web UI（推荐）

```bash
shield start
```

在 `http://localhost:8181` 添加应用，Protocol 选 `postgres`，填好地址、用户名、密码，点连接。

### 命令行

```bash
# 连接远程 PostgreSQL
shield postgres 10.0.0.5:5432 --db-user postgres --db-pass mypass --database mydb

# 只读模式
shield postgres 10.0.0.5:5432 --db-user postgres --readonly
```

不传凭证会交互式提示输入。`pg` 和 `postgresql` 是别名，效果一样：

```bash
shield pg 10.0.0.5:5432
```

---

## Web 界面功能

连接成功后浏览器自动打开 Web SQL 客户端，界面和 MySQL 插件类似，但针对 PostgreSQL 做了适配：

### Schema 树形浏览

PostgreSQL 按 Schema 组织，不像 MySQL 按 Database 切换。左侧栏展示 Schema → Table → Column/Index 三级树结构，支持搜索过滤。

### SQL 编辑器

- 多标签页，`Ctrl+Enter` 执行
- 结果排序、CSV 导出、单元格复制
- 双击单元格查看完整内容（长文本、JSONB 友好）

### 表结构管理

- 查看字段类型、默认值、约束
- 可视化创建表（支持 PG 类型：`SERIAL`、`BIGSERIAL`、`TIMESTAMPTZ`、`JSONB`、`UUID`、`INET` 等）
- 创建/删除索引
- 添加/删除字段
- 创建/删除 Schema

### 行级操作

- 插入记录
- 编辑单元格（双击）
- 删除行（带确认）

### 只读模式

和 MySQL 插件一样，只读/读写由启动方控制。只读模式下 INSERT、UPDATE、DELETE、DROP、ALTER、CREATE 等写操作被前后端双重拦截，远程用户无法绕过。

---

## 和 MySQL 插件的区别

两个插件功能结构一致，但在 SQL 层面做了完整的 PostgreSQL 适配：

| | MySQL 插件 | PostgreSQL 插件 |
|---|---|---|
| 组织方式 | Database → Table | Schema → Table |
| 标识符引用 | 反引号 `` ` `` | 双引号 `"` |
| 自增类型 | `AUTO_INCREMENT` | `SERIAL` / `BIGSERIAL` |
| 元数据查询 | `SHOW` 命令 | `information_schema` |
| 特有类型 | — | `JSONB`、`UUID`、`INET`、`TIMESTAMPTZ` |

---

## 插件源码独立维护

PostgreSQL 插件的代码不在 Shield CLI 主仓库里，而是放在独立的插件 monorepo：

[github.com/fengyily/shield-plugins](https://github.com/fengyily/shield-plugins)

后续 Redis、SQL Server 等插件也会放在这里。每个插件是独立的 Go module，互不依赖，CI 自动检测哪个插件有改动就构建哪个。

---

## 实际场景

### 临时给同事查 PG 数据

同事需要查几条数据，但他那台机器上没有 pgAdmin 也没有 DBeaver：

```bash
shield postgres 192.168.1.100:5432 --db-user readonly --readonly
```

把链接发给他，浏览器打开就能查。用完断开。

### 排查线上 Schema 结构

生产环境的 PG 不方便直连，通过 Shield 建立加密隧道，只读模式浏览 Schema、表结构和索引：

```bash
shield postgres prod-db.internal:5432 --db-user ops --readonly --invisible
```

`--invisible` 隐身模式需要授权码才能访问，防止链接泄露。

### Docker 里的 PG

```bash
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:16
shield postgres 127.0.0.1:5432 --db-user postgres --db-pass postgres
```

---

## 当前插件状态

| 插件 | 状态 |
|---|---|
| mysql | ✅ 已发布 |
| postgres | ✅ 已发布 |
| redis | 计划中 |
| sqlserver | 计划中 |

---

## 试一下

```bash
# 安装/升级 Shield CLI
brew install fengyily/tap/shield-cli
# 或
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh

# 安装 PostgreSQL 插件
shield plugin add postgres

# 连接
shield postgres 127.0.0.1:5432 --db-user postgres
```

开源地址：https://github.com/fengyily/shield-cli

插件仓库：https://github.com/fengyily/shield-plugins

有问题或建议欢迎提 [Issue](https://github.com/fengyily/shield-cli/issues)。
