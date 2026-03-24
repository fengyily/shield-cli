# Hacker News — Show HN

**发布地址**: https://news.ycombinator.com/submit

**最佳时间**: 美国东部时间周二至周四上午 9-11 点

---

**Title:**

Show HN: Shield CLI – Access SSH, RDP, VNC through your browser, no VPN needed

**Body:**

Hi HN, I built Shield CLI, an open-source tool (Apache 2.0, written in Go) that creates encrypted tunnels to internal services and renders them directly in the browser via HTML5.

The key difference from ngrok or frp: Shield CLI doesn't just forward ports. It renders full RDP desktops, VNC sessions, and SSH terminals in the browser using HTML5. No client software, no VPN, no RDP viewer — just a URL.

How it works: you run `shield` on the machine with access to internal services, it establishes an encrypted tunnel, and you get a browser-accessible endpoint. The web rendering layer handles the protocol translation (RDP/VNC/SSH → HTML5 canvas/terminal).

Practical uses: access a dev machine's desktop from a coffee shop, give a contractor temporary SSH access without VPN credentials, demo an internal tool to a client.

Supports macOS, Linux, and Windows. Docker images available. Includes a Web UI dashboard for managing tunnels.

Honest caveats: this is still a young project. I'd appreciate feedback on the architecture and UX.

GitHub: https://github.com/fengyily/shield-cli
