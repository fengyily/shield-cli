---
title: PostgreSQL 插件 v0.4.0：用 ER 图设计表关系，还能多人协作
description: Shield PostgreSQL 插件 v0.4.0 发布，新增交互式 ER 图——不只是看表关系，还能直接拖拽建外键、右键建表改字段，支持多人实时协作，远程团队同步讨论数据库设计。
date: 2026-03-26
author: Shield CLI Team
head:
  - - meta
    - name: keywords
      content: ER图, 数据库设计, PostgreSQL ER diagram, 表关系设计, 外键设计, 数据库协作, 实时协作, shield-postgres, Shield CLI, 远程数据库管理
---

# PostgreSQL 插件 v0.4.0：用 ER 图设计表关系，还能多人协作

> 今天早些时候聊了[用 Docker 一行命令查数据](/blogs/blog-shield-cli-docker-db)的事——解决的是"手边没工具"的问题。但查数据只是日常的一半，另一半是理解和设计表之间的关系。这次 PostgreSQL 插件 v0.4.0 的重点就是这个：**ER 图，而且是能直接操作的 ER 图**。

---

## 不只是看，是能动手的 ER 图

大多数 ER 图工具是"展示型"的——把表关系画出来给你看。这个不一样，它是直接连着真实数据库的，你在图上做的操作会真正执行。

具体来说：

### 拖拽建外键

从 A 表的某个字段拖到 B 表的字段上，外键就建好了。

- 类型匹配的字段，直接建立关联
- 类型不匹配，自动在目标表创建一个对应字段再建关联
- 拖的过程中有视觉提示——绿色表示类型匹配，蓝色表示需要新建字段
- 不想要的外键？点击关系线选中，按 Delete 删除

不用写 `ALTER TABLE ... ADD CONSTRAINT ... FOREIGN KEY ...`，拖一下就行。

### 右键建表、改字段

- 空白处右键 → 新建表，可视化添加字段，支持 PG 类型（`SERIAL`、`JSONB`、`UUID`、`TIMESTAMPTZ` 等）
- 表头右键 → 重命名、删除表、添加字段
- 字段右键 → 编辑类型、默认值、NOT NULL，或删除

### SQL 预览

每一步操作，执行前都会弹出实际要执行的 SQL。你清楚地知道它要做什么，确认了再执行。不是黑盒。

### 布局和导航

- 四种布局：网格 / 水平 / 垂直 / 放射状
- 拖拽移动表的位置，Ctrl+滚轮缩放
- 位置自动保存到 localStorage，下次打开还在原来的位置

---

## 多人实时协作

ER 图好用，但如果只能一个人看，很多场景就缺了一块。v0.4.0 加了实时协作。

### 怎么用

不需要额外配置。多个人打开同一个数据库的 ER 图，自动进入协作状态。每个人会被分配一个随机身份（熊猫、狐狸、老鹰之类的）和一个颜色。

### 能看到什么

- **在线列表**：顶部显示当前有哪些人在看这个 Schema 的 ER 图
- **实时光标**：别人的鼠标位置会实时显示在你的画布上，带名字标签和颜色
- **拖拽提示**：有人在拖动某张表时，你会看到那张表周围出现彩色虚线框和操作者名字
- **结构同步**：任何人修改了表结构（建表、加字段、建外键），所有人的 ER 图自动刷新

### 典型场景

**团队讨论方案时**

不用画板、不用截图、不用 "你看第三张表的第五个字段"。打开 ER 图，所有人看到同一个画面，谁在看哪里一目了然。讨论到哪张表，鼠标指过去，大家都能看到。

**远程工作时**

配合 Shield CLI 的远程访问能力，异地的同事也能接入同一个 ER 图：

```bash
shield postgres 10.0.0.5:5432 --db-user designer --db-pass xxx
```

把链接发给同事，浏览器打开就能一起看、一起讨论。比共享屏幕好用——每个人可以自己缩放、移动、查看不同区域，同时还能看到别人在看哪里。

**新人了解项目**

让新人打开 ER 图，你在旁边（或远程）用鼠标指着讲："这张 users 表和 orders 表是一对多关系，通过这个外键关联……" 对方能实时看到你的光标在哪里。

---

## 还是那么轻

ER 图和协作功能没有引入额外依赖。整个插件依然是一个 Go 静态二进制，Docker 镜像依然是个位数 MB，启动依然不到 1 秒。

ER 图的渲染是纯 SVG + 原生 JavaScript，没有 D3、没有 React、没有 Canvas。协作基于 WebSocket，服务端 300 行 Go 代码。

---

## 升级方式

```bash
# 升级插件
shield plugin update postgres

# 验证版本
shield plugin list
# NAME      VERSION  PROTOCOLS                 INSTALLED
# postgres  v0.4.0   postgres, pg, postgresql  2026-03-26T...
```

Docker 用户：

```bash
docker pull fengyily/shield-postgres
```

---

## 串一下

今天发了两篇，可以这么理解这两个东西的关系：

- **Docker 镜像**（[上一篇](/blogs/blog-shield-cli-docker-db)）：解决"我只是想查个数据"——一行命令，浏览器查数据，用完删掉
- **Shield CLI + 插件**（这一篇）：解决"我需要理解和设计表关系，还要和团队一起讨论"——ER 图 + 协作 + 远程访问

简单场景用 Docker 镜像，够了。需要 ER 设计、团队协作、远程分享的场景，用 Shield CLI。

---

## 接下来

ER 图和协作是 PostgreSQL 插件先上的。MySQL 插件的 ER 支持在做了，逻辑类似，预计很快跟上。

另外，Redis 和 SQL Server 插件也在计划中。

---

## 试一下

```bash
# 安装 Shield CLI
brew install fengyily/tap/shield-cli
# 或
curl -fsSL https://cdn.jsdelivr.net/gh/fengyily/shield-cli@main/install.sh | sh

# 安装/升级 PostgreSQL 插件
shield plugin add postgres

# 连接数据库，打开 ER 图
shield postgres 127.0.0.1:5432 --db-user postgres
```

开源地址：https://github.com/fengyily/shield-cli

插件仓库：https://github.com/fengyily/shield-plugins

有问题或建议欢迎提 [Issue](https://github.com/fengyily/shield-cli/issues)。
