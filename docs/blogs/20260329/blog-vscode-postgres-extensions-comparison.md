---
title: VS Code 里管理 PostgreSQL，有哪些选择？主流扩展横向对比
description: 对比 VS Code 中 5 款主流 PostgreSQL 扩展——SQLTools、PostgreSQL Explorer、Database Client、Microsoft PostgreSQL 和 Shield CLI PostgreSQL，从功能、ER 图、协作、远程分享等维度帮你选出最适合的那一款。
date: 2026-03-29
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: VS Code PostgreSQL 插件, VS Code 数据库管理, SQLTools, Database Client, PostgreSQL Explorer, Shield CLI PostgreSQL, VS Code 扩展对比, ER图, 数据库协作
---

# VS Code 里管理 PostgreSQL，有哪些选择？主流扩展横向对比

> 在 VS Code 里管 PostgreSQL，扩展市场搜一下能出来一整页。到底用哪个？这篇把几个主流选项摆在一起，按功能维度逐个对比，帮你省掉一个个装了试、试了卸的时间。

---

## 参与对比的 5 款扩展

![VS Code 扩展市场搜索 PostgreSQL](./images/shield-cli-postgresSQL-plugin-vscode-list.jpg)

| 扩展 | 作者 | 定位 | 许可 |
|---|---|---|---|
| **SQLTools** | Matheus Teixeira | 通用数据库客户端，需装 PG 驱动插件 | 免费开源 |
| **PostgreSQL Explorer** | Chris Kolkman | PostgreSQL 专用浏览器 | 免费开源 |
| **Database Client** | Weijan Chen | 全能数据库客户端，支持十几种数据库 | 基础免费，高级功能付费 |
| **PostgreSQL** | Microsoft | 微软官方 PG 扩展 | 免费（已停止维护） |
| **Shield CLI PostgreSQL** | Shield CLI | PostgreSQL Web 客户端 + ER 图 + 协作 | 免费开源 |

---

## 核心功能对比

### SQL 查询

所有扩展都支持在 VS Code 里写 SQL、执行查询、查看结果。差异在细节上：

| | SQLTools | PG Explorer | Database Client | Microsoft PG | Shield CLI PG |
|---|---|---|---|---|---|
| 多标签页查询 | ✅ | ❌ | ✅ | ❌ | ✅ |
| 语法高亮 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 结果排序/过滤 | 部分 | ❌ | ✅ | 部分 | ✅ |
| CSV 导出 | ❌ | ❌ | ✅（付费） | ❌ | ✅ |
| 单元格复制 | ✅ | ✅ | ✅ | ✅ | ✅ |

Shield CLI PostgreSQL 的查询界面是内嵌的 Web 客户端，体验接近独立的数据库 IDE，而不是 VS Code 原生的表格渲染。

### Schema 浏览

| | SQLTools | PG Explorer | Database Client | Microsoft PG | Shield CLI PG |
|---|---|---|---|---|---|
| 树形结构 | Schema → Table | Schema → Table → Column | Schema → Table → Column | Database → Table | Schema → Table → Column → Index |
| 搜索过滤 | ❌ | ❌ | ✅ | ❌ | ✅ |

### 表结构管理

这一项差距比较明显：

| | SQLTools | PG Explorer | Database Client | Microsoft PG | Shield CLI PG |
|---|---|---|---|---|---|
| 可视化建表 | ❌ | ❌ | ✅ | ❌ | ✅ |
| 添加/删除字段 | ❌ | ❌ | ✅ | ❌ | ✅ |
| 管理索引 | ❌ | ❌ | ✅ | ❌ | ✅ |
| 创建/删除 Schema | ❌ | ❌ | 部分 | ❌ | ✅ |

SQLTools 和 PG Explorer 定位是"查询工具"，表结构变更得手写 DDL。Database Client 和 Shield CLI PostgreSQL 提供可视化操作。

### 行级编辑

| | SQLTools | PG Explorer | Database Client | Microsoft PG | Shield CLI PG |
|---|---|---|---|---|---|
| 插入行 | ❌ | ❌ | ✅ | ❌ | ✅ |
| 双击编辑 | ❌ | ❌ | ✅ | ❌ | ✅ |
| 删除行 | ❌ | ❌ | ✅ | ❌ | ✅ |

SQLTools 和 PG Explorer 只能看、不能在表格里直接改。需要修改数据得写 INSERT / UPDATE / DELETE。

---

## 差异化功能

上面几项，Database Client 和 Shield CLI PostgreSQL 都做得不错。真正拉开差距的是下面这几个维度。

### ER 图

| | SQLTools | PG Explorer | Database Client | Microsoft PG | Shield CLI PG |
|---|---|---|---|---|---|
| ER 图 | ❌ | ❌ | ✅（付费） | ❌ | ✅（免费） |
| 可交互操作 | — | — | 仅查看 | — | 拖拽建外键、右键建表改字段 |
| SQL 预览 | — | — | ❌ | — | 每步操作前显示实际 SQL |

![Shield CLI PostgreSQL ER 图——左边 VS Code 内嵌，右边浏览器打开](./images/shield-cli-plugin-postgres-xiezhuo.jpg)

Database Client 的 ER 图是付费功能，且只能查看。Shield CLI PostgreSQL 的 ER 图免费且可交互——拖拽字段就能建外键，右键建表加字段，每一步操作都会先弹出 SQL 让你确认。

### 远程分享 / 浏览器访问

| | SQLTools | PG Explorer | Database Client | Microsoft PG | Shield CLI PG |
|---|---|---|---|---|---|
| 浏览器打开 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 分享链接给他人 | ❌ | ❌ | ❌ | ❌ | ✅ |

这是 Shield CLI PostgreSQL 独有的能力。其他扩展都是"装在谁的 VS Code 里，谁能用"。Shield CLI PostgreSQL 可以点击 **Open in Browser** 在浏览器打开，也可以把链接发给同事，对方不需要装 VS Code，浏览器打开就能查。

典型场景：同事需要查几条数据但他电脑上没装任何数据库工具，发个链接就解决了。

### 实时协作

| | SQLTools | PG Explorer | Database Client | Microsoft PG | Shield CLI PG |
|---|---|---|---|---|---|
| 多人同时使用 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 实时光标 | — | — | — | — | ✅ |
| 操作同步 | — | — | — | — | ✅ |

多人打开同一个数据库连接，每个人都能看到其他人的光标位置、正在拖拽哪张表。讨论数据库设计时，不用共享屏幕，不用截图，直接指给对方看。

### 只读模式

| | SQLTools | PG Explorer | Database Client | Microsoft PG | Shield CLI PG |
|---|---|---|---|---|---|
| 只读模式 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 后端强制拦截 | — | — | — | — | ✅ |

其他扩展要实现只读，只能靠数据库账号权限。Shield CLI PostgreSQL 提供应用层的只读模式：前端禁用写操作按钮，后端拦截 INSERT / UPDATE / DELETE / DROP / ALTER / CREATE 等语句，双重保护。把链接分享给别人时，可以确保对方只能查不能改。

---

## 总结对比

| 维度 | SQLTools | PG Explorer | Database Client | Microsoft PG | Shield CLI PG |
|---|---|---|---|---|---|
| SQL 查询 | ✅ | ✅ | ✅ | ✅ | ✅ |
| Schema 浏览 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 表结构管理 | ❌ | ❌ | ✅ | ❌ | ✅ |
| 行级编辑 | ❌ | ❌ | ✅ | ❌ | ✅ |
| ER 图 | ❌ | ❌ | 💰 付费 | ❌ | ✅ 免费 |
| ER 图可交互 | — | — | ❌ | — | ✅ |
| 浏览器访问 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 远程分享 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 实时协作 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 只读模式 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 多数据库支持 | ✅ | ❌ | ✅ | ❌ | ❌ |
| 维护状态 | 活跃 | 较慢 | 活跃 | 已停维 | 活跃 |
| 费用 | 免费 | 免费 | 部分付费 | 免费 | 免费 |

---

## 怎么选？

**只需要写 SQL 查数据** → SQLTools 或 PG Explorer 就够了，轻量，不折腾。

**需要表结构管理和行级编辑** → Database Client 或 Shield CLI PostgreSQL，看你需不需要后面的能力。

**需要 ER 图** → Database Client（付费，仅查看）或 Shield CLI PostgreSQL（免费，可交互操作）。

**需要把数据库界面分享给别人** → Shield CLI PostgreSQL 是目前唯一支持浏览器访问和链接分享的选项。

**团队协作讨论数据库设计** → Shield CLI PostgreSQL 的实时协作功能是独有的。

**管多种数据库** → SQLTools 和 Database Client 支持 MySQL、SQLite、SQL Server 等多种数据库。Shield CLI PostgreSQL 目前只支持 PostgreSQL（Shield CLI 另有 MySQL 插件）。

---

## 安装 Shield CLI PostgreSQL

在 VS Code 扩展面板搜索 **Shield CLI PostgreSQL**，点 Install 即可。

![安装](./images/shield-cli-postgresSQL-plugin-vscode-list.jpg)

开源地址：

- Shield CLI：https://github.com/fengyily/shield-cli
- 插件源码：https://github.com/fengyily/shield-plugins

有问题或建议欢迎提 [Issue](https://github.com/fengyily/shield-cli/issues)。
