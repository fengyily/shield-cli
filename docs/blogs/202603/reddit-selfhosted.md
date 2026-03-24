# Reddit r/selfhosted

**发布地址**: https://www.reddit.com/r/selfhosted/submit

---

**Title:**

I built an open-source tool that lets you access RDP, SSH, and VNC through any browser — no VPN or client apps needed

**Body:**

Hey r/selfhosted,

I've been working on **Shield CLI**, an open-source tool that creates encrypted tunnels to your internal services and makes them accessible directly in the browser.

The thing that sets it apart from tools like ngrok or frp: it doesn't just expose a port. It actually **renders RDP desktops, VNC sessions, and SSH terminals as HTML5** in the browser. You don't need an RDP client, a VNC viewer, or even a terminal emulator — everything runs in the browser tab.

**Self-hosting highlights:**

- Written in Go, single binary, runs on macOS/Linux/Windows
- Docker support for easy deployment
- Web UI dashboard to manage your tunnels
- Apache 2.0 license — fully open source
- Encrypted tunnels by default

**Use cases I built it for:**

- Accessing my home server's desktop remotely without setting up a VPN
- Giving temporary SSH access to someone without sharing keys or VPN configs
- Quick remote support for family members (just send them a link)

It's still a relatively young project, so I'd love to hear what features matter most to you.

GitHub: https://github.com/fengyily/shield-cli
