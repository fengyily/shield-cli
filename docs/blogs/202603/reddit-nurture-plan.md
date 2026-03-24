# Reddit 养号计划

**目标**: 养一个有真实社区参与度的 Reddit 账号，2-3 周后自然推广 Shield CLI
**新号注册日**: 2026-03-24
**预计可发推广帖**: 2026-04-14（第 3 周后）

---

## 目标 Subreddit

| Subreddit | 用途 | 备注 |
|-----------|------|------|
| r/homelab | 主推广阵地 | 最核心，必须深度参与 |
| r/selfhosted | 主推广阵地 | 对开源项目最友好 |
| r/linux | 互动攒 karma | 活跃度高，容易获得 upvote |
| r/docker | 互动攒 karma | Shield CLI 有 Docker 部署 |
| r/sysadmin | 互动攒 karma | IT 运维人群 |
| r/opensource | 备用推广 | 专门分享开源项目 |
| r/networking | 偶尔互动 | 网络/隧道相关话题 |

---

## 第 1 周：纯互动期（3/24 - 3/30）

> 目标：karma ≥ 30，建立活跃痕迹，零推广

### Day 1（3/24 周二）

**r/homelab** — 浏览 Hot/New 帖子，找 2-3 个帖子真诚回复

示例场景与回复：

- 有人晒 homelab 硬件 →
  `"Nice setup! What are you running on the Proxmox node? I've been thinking about adding a second NUC for HA."`

- 有人问远程访问方案 →
  `"I've tried Tailscale and WireGuard for this. Tailscale is easier to set up but WireGuard gives you more control. What's your use case — just SSH or do you need GUI access too?"`

**r/linux** — 找 1-2 个新手提问帖回答

- 有人问 Linux 命令 →
  `"You can use 'ss -tulnp' to check which ports are listening. It's the modern replacement for netstat."`

---

### Day 2（3/25 周三）

**r/selfhosted** — 回复 2-3 个帖子

示例场景：

- 有人问自托管方案推荐 →
  `"For dashboards I'd recommend Homarr or Homepage — both are lightweight and easy to set up with Docker Compose."`

- 有人讨论 reverse proxy →
  `"Caddy is great if you want automatic HTTPS with minimal config. Nginx Proxy Manager if you prefer a GUI."`

**r/docker** — 回复 1-2 个帖子

- Docker Compose 问题 →
  `"Try adding 'depends_on' with a healthcheck condition. The default depends_on only waits for the container to start, not for the service to be ready."`

---

### Day 3（3/26 周四）

**r/homelab** — 回复 2 个帖子 + 点赞互动

- 有人问预算方案 →
  `"Used Dell Optiplex micro PCs are amazing for homelabs. Low power, quiet, and you can get them for $50-100. I run Proxmox on one."`

**r/sysadmin** — 回复 1-2 个帖子

- 有人吐槽远程管理 →
  `"RDP over VPN has been my go-to, but the setup overhead is real. Have you looked at browser-based solutions like Apache Guacamole? Zero client install needed."`

---

### Day 4（3/27 周五）

**r/selfhosted** — 回复 2 个帖子

- 有人问 SSH 管理多台机器 →
  `"I use a combination of SSH config file with aliases and tmux. For a web-based approach, there's also Guacamole but it can be heavy to run."`

**r/linux** — 回复 1-2 个帖子

---

### Day 5（3/28 周六）

**r/homelab** — 发第一个非推广帖（提问或分享经验）

标题：`"What's your go-to solution for accessing homelab machines remotely?"`

内容：
```
I've been experimenting with different remote access setups for my homelab and curious what everyone is using.

My current setup:
- WireGuard VPN for SSH access
- RDP through the VPN for Windows VMs
- Guacamole for browser-based access when I'm on the go

The VPN works great from my own devices, but it's annoying when I'm on a work computer or someone else's machine where I can't install a VPN client.

What's your setup? Anyone found a good solution for accessing machines from devices you don't control?
```

> 这个帖子自然引出"浏览器访问远程机器"的痛点，为后续推广做铺垫

---

### Day 6（3/29 周日）

**回复 Day 5 帖子中的评论** — 和每个回复者互动，感谢分享、追问细节

**r/docker** — 回复 1-2 个帖子

---

### Day 7（3/30 周一）

**r/selfhosted** + **r/homelab** — 各回复 1-2 个帖子

---

## 第 2 周：深度参与期（3/31 - 4/6）

> 目标：karma ≥ 80，成为社区认识的面孔

### Day 8（3/31 周二）

**r/selfhosted** — 回复 2-3 个帖子，重点关注远程访问、VPN、Docker 相关话题

### Day 9（4/1 周三）

**r/homelab** — 回复 2 个帖子

**r/linux** — 分享一个实用技巧帖

标题：`"TIL: You can use 'ssh -D' as a quick SOCKS proxy for browser access to your home network"`

内容：
```
Just a quick tip that saved me today. If you have SSH access to a machine on your home network, you can use:

    ssh -D 1080 user@your-server

Then set your browser's SOCKS proxy to localhost:1080. Now your browser traffic routes through your home network — great for accessing internal web UIs without a full VPN setup.

Not a replacement for a proper VPN, but handy in a pinch.
```

---

### Day 10（4/2 周四）

**r/homelab** — 回复 2-3 个帖子
**r/sysadmin** — 回复 1 个帖子

### Day 11（4/3 周五）

**r/selfhosted** — 发一个讨论帖

标题：`"Browser-based remote desktop vs traditional clients — anyone made the full switch?"`

内容：
```
I've been exploring browser-based remote access tools (Guacamole, meshcentral, etc.) as an alternative to traditional RDP/VNC clients.

Pros I've found:
- Zero client install — works from any device with a browser
- Easy to share access with others (just send a URL)
- Works through firewalls without port forwarding

Cons:
- Performance isn't quite as good as native RDP
- Some tools are complex to set up (looking at you, Guacamole)

Has anyone fully replaced their traditional remote desktop clients with a browser-based solution? What's been your experience?
```

> 继续强化"浏览器远程访问"话题的存在感

### Day 12（4/4 周六）

**回复 Day 11 帖子中的评论** — 深度互动

**r/docker** — 回复 1-2 个帖子

### Day 13（4/5 周日）

**r/homelab** — 回复 2 个帖子
**r/linux** — 回复 1 个帖子

### Day 14（4/6 周一）

**r/selfhosted** + **r/sysadmin** — 各回复 1-2 个帖子

---

## 第 3 周：软推广期（4/7 - 4/13）

> 目标：自然提及 Shield CLI，不硬推

### Day 15（4/7 周二）

在相关讨论中**首次自然提及**（在评论里，不是发帖）

当有人问远程访问方案时：
```
"I actually built a tool for this — Shield CLI. It creates encrypted tunnels
and renders RDP/SSH/VNC in the browser. Still early stage but it works well
for my homelab. Happy to share the GitHub link if you're interested."
```

> 关键：不主动贴链接，等人问再给。显得自然，不像 spam

### Day 16（4/8 周三）

**r/homelab** + **r/selfhosted** — 正常回帖互动 2-3 个

### Day 17（4/9 周四）

如果 Day 15 有人问链接 → 回复 GitHub 链接
如果没人问 → 继续正常互动，找下一个合适的时机

### Day 18（4/10 周五）

**r/selfhosted** — 正常互动，在合适的帖子评论中再次自然提及

### Day 19（4/11 周六）

**r/homelab** + **r/docker** — 正常互动

### Day 20（4/12 周日）

准备正式推广帖的内容，根据前两周的社区反馈调整措辞

### Day 21（4/13 周一）

复盘 karma 和互动情况，确认是否达到发帖条件

---

## 第 4 周：正式推广（4/14 起）

> 前提：karma ≥ 100，帖子/评论历史看起来是真实用户

### Day 22（4/14 周二）— r/selfhosted 首发

标题：`"I built an open-source tool for browser-based RDP/SSH/VNC access — no VPN needed"`

内容：使用 [reddit-selfhosted.md](reddit-selfhosted.md) 中准备的内容，根据社区反馈微调

### Day 24（4/16 周四）— r/homelab 发帖

标题：`"Access your homelab machines (RDP, SSH, VNC) from any browser — open-source, no VPN required"`

内容：使用 [reddit-homelab.md](reddit-homelab.md) 中准备的内容

> 两个帖子间隔 2 天，避免被认为是批量 spam

---

## 每日核心规则

1. **每天至少 2-3 条真诚回复**，不要一次性灌水
2. **回复要有实质内容**，不要只说 "nice" "thanks" "cool"
3. **不同时间段发帖/回复**，模拟真实用户行为
4. **前 2 周绝对不提 Shield CLI**
5. **永远不要在帖子标题放 GitHub 链接**
6. **被人质疑 self-promotion 时态度诚恳**：承认是自己的项目，强调开源免费
7. **帖子中的讨论一定要回复**，不要发完就跑

## Karma 获取技巧

- **回复新帖比热帖更容易拿 karma** — 排在前面的评论会持续吃 upvote
- **帮人解决具体问题** 比发表观点更容易拿 upvote
- **r/linux 的新手问题帖** 是最容易攒 karma 的地方
- **带格式的回复**（代码块、列表）比纯文字更容易拿 upvote
