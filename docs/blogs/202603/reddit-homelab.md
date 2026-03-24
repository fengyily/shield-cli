# Reddit r/homelab

**发布地址**: https://www.reddit.com/r/homelab/submit

---

**Title:**

Access your homelab machines (RDP, SSH, VNC) from any browser — open-source, no VPN required

**Body:**

Fellow homelabbers,

I built **Shield CLI** to solve a problem I kept running into: I want to access my homelab machines from anywhere, but I don't always want to deal with VPN configs, and I don't always have my usual laptop with the right clients installed.

Shield CLI creates encrypted tunnels to your services and renders them **directly in the browser**. Full RDP desktop sessions, VNC, and SSH terminals — all in an HTML5 browser tab. No client software needed on the accessing side.

**Homelab scenarios:**

- RDP into your Windows VM from your phone's browser
- SSH into your Proxmox host from a work computer where you can't install anything
- Give a friend VNC access to help debug something — just send a link
- Manage multiple machines from the built-in Web UI dashboard

**Tech details:**

- Single Go binary, runs on macOS/Linux/Windows
- Docker support (`docker run` and you're set)
- Apache 2.0 open source
- Encrypted tunnels, no port forwarding needed on your router

Still actively developing — feedback and feature requests very welcome.

GitHub: https://github.com/fengyily/shield-cli
