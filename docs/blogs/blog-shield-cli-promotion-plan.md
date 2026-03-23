# Shield CLI 开源推广计划

> 目标：让更多开发者和运维人员了解 Shield CLI，提升 GitHub star 数和社区活跃度。

---

## 第一阶段：完善项目基础（第 1 周）

### 1. ✅ 添加 GitHub Topics（已完成 2026-03-23）

已通过 API 添加以下 14 个 Topics：

```
browser, cli, devops, golang, networking, rdp, remote-access,
remote-desktop, reverse-proxy, selfhosted, shield, ssh, tunnel, vnc
```

### 2. 添加社交预览图（Social Preview）

- 在 GitHub 仓库 Settings → Social preview 上传预览图
- 分享链接时会显示项目 Logo 和简介，提升点击率

### 3. 录制 Demo GIF / 视频

- 使用 [vhs](https://github.com/charmbracelet/vhs) 或 [asciinema](https://asciinema.org/) 录制 30 秒终端操作 GIF
- 重点展示 `shield rdp` / `shield vnc` 一键浏览器访问远程桌面的效果
- 放在 README 顶部，直观展示核心卖点

---

## 第二阶段：提交到开源目录（第 2-3 周）

### 4. 提交 PR 到 github/explore

- 仓库地址：https://github.com/github/explore
- 将 shield-cli 加入 `topics/` 下相关 topic 的 collection
- 有助于出现在 GitHub Explore 页面

### 5. 提交到 Awesome Lists

| Awesome List | 地址 | 说明 | 准入条件 |
|-------------|------|------|----------|
| awesome-tunneling | https://github.com/anderspitman/awesome-tunneling | 隧道工具专项列表，最精准 | ⚠️ 需 100+ stars（2026-02-16 新政策） |
| awesome-selfhosted | https://github.com/awesome-selfhosted/awesome-selfhosted-data | 自托管工具合集，流量大 | ⚠️ 需首次发布超过 4 个月（最早 2026-07-13） |
| awesome-go | https://github.com/avelino/awesome-go | Go 项目合集 | 待确认 |

**awesome-tunneling 预备 PR 条目（待 star 达标后提交）：**

```markdown
* [Shield CLI](https://github.com/fengyily/shield-cli) [![shield-cli github stars badge](https://img.shields.io/github/stars/fengyily/shield-cli?style=flat)](https://github.com/fengyily/shield-cli/stargazers) - Exposes internal services (SSH, RDP, VNC, HTTP) accessible from any browser via HTML5 rendering. No VPN, no client install needed. Written in Go.
```

**awesome-selfhosted 预备 YML 文件（待时间达标后提交到 awesome-selfhosted-data）：**

```yaml
name: "Shield CLI"
website_url: "https://docs.yishield.com/"
source_code_url: "https://github.com/fengyily/shield-cli"
description: "Expose internal services (SSH, RDP, VNC, HTTP) accessible from any browser via HTML5 rendering, no VPN or client install needed."
licenses:
  - Apache-2.0
platforms:
  - Go
  - Docker
tags:
  - Remote Access - VPN
```

提交时需遵守各列表的 Contributing 指南，确保项目描述简洁准确。

### 6. 注册到工具索引网站

- [Product Hunt](https://www.producthunt.com) — 发布 launch，适合有 UI 的工具
- [AlternativeTo](https://alternativeto.net) — 注册为 ngrok / frp 的替代方案
- [LibHunt](https://www.libhunt.com) — 开源项目发现平台

---

## 第三阶段：社区推广（第 3-5 周）

### 7. Hacker News — Show HN

- 标题格式：`Show HN: Shield CLI – Access RDP/VNC/SSH from any browser, no VPN needed`
- 最佳发布时间：美国东部时间周二至周四上午 9-11 点
- 准备好回复评论，展示技术深度

### 8. Reddit 发帖

| 子版块 | 受众 |
|--------|------|
| r/selfhosted | 自托管爱好者，最精准的目标受众 |
| r/homelab | 家庭服务器用户 |
| r/devops | DevOps 工具用户 |
| r/golang | Go 语言社区 |
| r/sysadmin | 系统管理员 |

每个子版块发帖前先阅读规则，避免被当作广告删除。建议以"分享项目 + 征求反馈"的口吻撰写。

### 9. 中文社区

| 平台 | 说明 |
|------|------|
| [V2EX](https://www.v2ex.com) | `/go/programmer` 或 `/go/share` 节点 |
| [掘金](https://juejin.cn) | 发布技术文章，配合教程 |
| [SegmentFault](https://segmentfault.com) | 技术问答 + 文章 |
| [知乎](https://www.zhihu.com) | 回答相关问题，附带项目链接 |
| [少数派](https://sspai.com) | 如有桌面客户端或 GUI |

---

## 第四阶段：持续运营（长期）

### 10. 技术博客输出

持续产出高质量文章，以下为推荐选题：

- Shield CLI vs ngrok vs frp 详细对比（已有）
- Shield CLI Docker 部署指南（已有）
- 如何用 Shield CLI 实现零客户端远程桌面
- Shield CLI 在 CI/CD 中的应用
- 从浏览器直接 SSH：Shield CLI 实践

发布平台：dev.to、Medium、掘金、知乎专栏。

### 11. 参与社区讨论

- 在 ngrok / frp 的 GitHub Issues 中，当用户需要浏览器 RDP/VNC 功能时，礼貌推荐
- 回答 StackOverflow 上 remote access、tunnel 相关问题
- 参与 Reddit / V2EX 相关话题讨论

### 12. 保持发布节奏

- 每 2-4 周发布新版本，持续出现在 GitHub Release feed 中
- 每次发布附带 Changelog，展示项目活跃度
- 通过 GitHub Discussions 与社区保持互动

---

## 优先级总览

| 优先级 | 动作 | 预期效果 | 难度 | 状态 |
|--------|------|----------|------|------|
| **P0** | 添加 Topics | 被搜索发现的前提 | 低 | ✅ 已完成 |
| **P0** | 提交 awesome-tunneling | 精准受众，审核快 | 低 | ⏳ 需 100+ stars |
| **P1** | 录制 Demo GIF | 提升 README 转化率 | 中 | 待执行 |
| **P1** | Hacker News Show HN | 短期大量曝光 | 中 | 待执行 |
| **P1** | 提交 awesome-selfhosted | 长期稳定流量 | 低 | ⏳ 需发布满 4 个月 |
| **P2** | Product Hunt Launch | 一次性大曝光 | 中 | 待执行 |
| **P2** | 中文社区推广 | 国内用户增长 | 低 | 待执行 |
| **P2** | Reddit 多版块发帖 | 持续引流 | 低 | 待执行 |
| **P3** | 技术博客持续输出 | SEO 长尾流量 | 高 | 待执行 |
| **P3** | 社区讨论参与 | 口碑积累 | 中 | 待执行 |

---

## 关键指标跟踪

- GitHub Stars 增长趋势
- README 页面访问量（GitHub Traffic）
- Discussions 参与度
- 外部引荐来源（GitHub Traffic → Referring sites）
