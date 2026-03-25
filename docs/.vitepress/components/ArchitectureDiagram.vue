<template>
  <div class="arch">
    <svg :viewBox="`0 0 ${W} ${H}`" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="tunnelGrad" x1="0" y1="0" x2="1" y2="0">
          <stop offset="0%" stop-color="var(--c-brand)" stop-opacity="0.15" />
          <stop offset="50%" stop-color="var(--c-brand)" stop-opacity="0.06" />
          <stop offset="100%" stop-color="var(--c-brand)" stop-opacity="0.15" />
        </linearGradient>
        <linearGradient id="shieldGrad" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stop-color="var(--c-brand)" stop-opacity="0.12" />
          <stop offset="100%" stop-color="var(--c-brand)" stop-opacity="0.04" />
        </linearGradient>
        <filter id="shieldGlow" x="-30%" y="-30%" width="160%" height="160%">
          <feGaussianBlur in="SourceGraphic" stdDeviation="6" result="blur" />
          <feFlood flood-color="var(--c-brand)" flood-opacity="0.25" result="color" />
          <feComposite in="color" in2="blur" operator="in" result="glow" />
          <feMerge><feMergeNode in="glow" /><feMergeNode in="SourceGraphic" /></feMerge>
        </filter>

        <!-- ── Icon symbols ── -->
        <!-- SSH terminal -->
        <symbol id="ico-ssh" viewBox="0 0 24 24">
          <rect x="2" y="4" width="20" height="16" rx="2" fill="none" stroke="currentColor" stroke-width="1.5" />
          <text x="6" y="16" font-size="9" font-weight="700" font-family="monospace" fill="currentColor">$_</text>
        </symbol>
        <!-- MySQL database -->
        <symbol id="ico-db" viewBox="0 0 24 24">
          <ellipse cx="12" cy="6" rx="8" ry="3" fill="none" stroke="currentColor" stroke-width="1.5" />
          <path d="M4 6v12c0 1.66 3.58 3 8 3s8-1.34 8-3V6" fill="none" stroke="currentColor" stroke-width="1.5" />
          <path d="M4 12c0 1.66 3.58 3 8 3s8-1.34 8-3" fill="none" stroke="currentColor" stroke-width="1.5" />
        </symbol>
        <!-- RDP/HTTP web page -->
        <symbol id="ico-web" viewBox="0 0 24 24">
          <rect x="2" y="3" width="20" height="16" rx="2" fill="none" stroke="currentColor" stroke-width="1.5" />
          <line x1="2" y1="8" x2="22" y2="8" stroke="currentColor" stroke-width="1.5" />
          <circle cx="5" cy="5.5" r="0.8" fill="currentColor" />
          <circle cx="7.5" cy="5.5" r="0.8" fill="currentColor" />
          <rect x="5" y="10.5" width="6" height="6" rx="1" fill="none" stroke="currentColor" stroke-width="1" />
          <line x1="14" y1="11" x2="19" y2="11" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" />
          <line x1="14" y1="14" x2="18" y2="14" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" />
          <line x1="14" y1="17" x2="17" y2="17" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" />
        </symbol>
        <!-- Shield -->
        <symbol id="ico-shield" viewBox="0 0 24 24">
          <path d="M12 2L4 6v5c0 5.25 3.4 10.15 8 11.25C16.6 21.15 20 16.25 20 11V6l-8-4z"
            fill="none" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round" />
          <polyline points="8.5,12 11,14.5 15.5,9.5" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
        </symbol>
        <!-- Cloud -->
        <symbol id="ico-cloud" viewBox="0 0 24 24">
          <path d="M6.5 19h11A4.5 4.5 0 0 0 22 14.5a4.5 4.5 0 0 0-3.17-4.3A6.5 6.5 0 0 0 6.5 7.5 5.5 5.5 0 0 0 1 13 5.5 5.5 0 0 0 6.5 19z"
            fill="none" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round" />
        </symbol>
        <!-- Browser with globe -->
        <symbol id="ico-browser" viewBox="0 0 24 24">
          <rect x="2" y="3" width="20" height="18" rx="2.5" fill="none" stroke="currentColor" stroke-width="1.5" />
          <line x1="2" y1="8" x2="22" y2="8" stroke="currentColor" stroke-width="1.5" />
          <circle cx="5" cy="5.5" r="0.8" fill="currentColor" />
          <circle cx="7.5" cy="5.5" r="0.8" fill="currentColor" />
          <circle cx="10" cy="5.5" r="0.8" fill="currentColor" />
          <circle cx="12" cy="15" r="4.5" fill="none" stroke="currentColor" stroke-width="1.2" />
          <ellipse cx="12" cy="15" rx="2" ry="4.5" fill="none" stroke="currentColor" stroke-width="1" />
          <line x1="7.5" y1="15" x2="16.5" y2="15" stroke="currentColor" stroke-width="1" />
        </symbol>
        <!-- Lock -->
        <symbol id="ico-lock" viewBox="0 0 24 24">
          <rect x="5" y="11" width="14" height="10" rx="2" fill="none" stroke="currentColor" stroke-width="1.5" />
          <path d="M8 11V7a4 4 0 0 1 8 0v4" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
          <circle cx="12" cy="16" r="1.5" fill="currentColor" />
        </symbol>
        <!-- Device icons for "any device" -->
        <symbol id="ico-laptop" viewBox="0 0 20 20">
          <rect x="3" y="3" width="14" height="10" rx="1.5" fill="none" stroke="currentColor" stroke-width="1.2" />
          <path d="M1 15h18" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" />
        </symbol>
        <symbol id="ico-phone" viewBox="0 0 20 20">
          <rect x="6" y="2" width="8" height="16" rx="1.5" fill="none" stroke="currentColor" stroke-width="1.2" />
          <line x1="9" y1="15" x2="11" y2="15" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" />
        </symbol>
        <symbol id="ico-tablet" viewBox="0 0 20 20">
          <rect x="4" y="2" width="12" height="16" rx="1.5" fill="none" stroke="currentColor" stroke-width="1.2" />
          <line x1="9" y1="15" x2="11" y2="15" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" />
        </symbol>
      </defs>

      <!-- ── Zone backgrounds ── -->
      <rect x="10" y="10" :width="W * 0.40 - 12" :height="H - 20" rx="16" class="zone zone--lan" />
      <rect :x="W * 0.60 + 2" y="10" :width="W * 0.40 - 12" :height="H - 20" rx="16" class="zone zone--wan" />

      <text :x="W * 0.20" y="34" class="zone-label">Internal Services</text>
      <text :x="W * 0.80" y="34" class="zone-label">Browser Access</text>

      <!-- ── App fan nodes (left) ── -->
      <g v-for="(app, i) in apps" :key="'app-'+i"
        class="node app-node" :class="{ active: step >= 2, pulse: step === 2 }">
        <rect :x="app.x - app.w/2" :y="app.y - app.h/2" :width="app.w" :height="app.h" rx="10" class="node-box" />
        <use :href="'#' + app.ico" :x="app.x - 9" :y="app.y - 16" width="18" height="18" class="node-ico" />
        <text :x="app.x" :y="app.y + 19" class="node-name-sm">{{ app.label }}</text>
      </g>

      <!-- Fan-in lines: Apps → Shield CLI -->
      <line v-for="(app, i) in apps" :key="'fan-'+i"
        :x1="app.x + app.w/2 + 1" :y1="app.y"
        :x2="shieldNode.x - shieldNode.w/2 - 2" :y2="shieldNode.y"
        class="link" :class="{ active: step >= 2 }" />

      <!-- ── Shield CLI node (hero) ── -->
      <g class="node node--hero" :class="{ active: step >= 1, pulse: step === 1 }">
        <rect :x="shieldNode.x - shieldNode.w/2" :y="shieldNode.y - shieldNode.h/2"
          :width="shieldNode.w" :height="shieldNode.h" rx="16" class="node-box hero-box" />
        <image href="/logo.svg" :x="shieldNode.x - 18" :y="shieldNode.y - 26" width="36" height="36" class="hero-logo" />
        <text :x="shieldNode.x" :y="shieldNode.y + 26" class="node-name hero-name">Shield CLI</text>
      </g>

      <!-- ── Tunnel (center) — 3 independent channels ── -->
      <g class="tunnel" :class="{ active: step >= 3 }">
        <rect :x="tunnelX1 - 4" :y="CY - tunnelHalf - 10" :width="tunnelX2 - tunnelX1 + 8" :height="tunnelHalf * 2 + 20" rx="12"
          fill="url(#tunnelGrad)" />
        <!-- 3 channel lanes -->
        <line v-for="(ly, i) in laneYs" :key="'lane-'+i"
          :x1="tunnelX1" :y1="ly" :x2="tunnelX2" :y2="ly" class="tunnel-line" />
        <use href="#ico-lock" :x="W * 0.50 - 6" :y="CY - tunnelHalf - 7" width="12" height="12" class="tunnel-ico" />
        <text :x="W * 0.50" :y="CY + tunnelHalf + 14" class="tunnel-label">Tunnel</text>
        <text :x="W * 0.50" :y="CY + tunnelHalf + 23" class="tunnel-sublabel">Encrypted</text>
      </g>

      <!-- Shield CLI → tunnel (3 lines converging into lanes) -->
      <line v-for="(ly, i) in laneYs" :key="'s2t-'+i"
        :x1="shieldNode.x + shieldNode.w/2 + 2" :y1="shieldNode.y"
        :x2="tunnelX1" :y2="ly"
        class="link" :class="{ active: step >= 3 }" />
      <!-- tunnel → Webgate (3 lines converging out) -->
      <line v-for="(ly, i) in laneYs" :key="'t2w-'+i"
        :x1="tunnelX2" :y1="ly"
        :x2="webgateNode.x - webgateNode.w/2 - 2" :y2="webgateNode.y"
        class="link" :class="{ active: step >= 3 }" />

      <!-- ── Webgate node ── -->
      <g class="node" :class="{ active: step >= 3, pulse: step === 3 }">
        <rect :x="webgateNode.x - webgateNode.w/2" :y="webgateNode.y - webgateNode.h/2"
          :width="webgateNode.w" :height="webgateNode.h" rx="14" class="node-box" />
        <use href="#ico-cloud" :x="webgateNode.x - 13" :y="webgateNode.y - 20" width="26" height="26" class="node-ico" />
        <text :x="webgateNode.x" :y="webgateNode.y + 22" class="node-name">Webgate</text>
      </g>

      <!-- URL badge under Webgate -->
      <g class="url-badge" :class="{ show: step >= 3 }">
        <rect :x="webgateNode.x - 56" :y="webgateNode.y + webgateNode.h/2 + 6" width="112" height="20" rx="10" class="url-bg" />
        <text :x="webgateNode.x" :y="webgateNode.y + webgateNode.h/2 + 20" class="url-text">https://xxx.yishield.com</text>
      </g>

      <!-- Webgate → Browser -->
      <line :x1="webgateNode.x + webgateNode.w/2 + 2" :y1="webgateNode.y"
        :x2="browserNode.x - browserNode.w/2 - 2" :y2="browserNode.y"
        class="link" :class="{ active: step >= 4 }" />

      <!-- ── Browser node ── -->
      <g class="node" :class="{ active: step >= 4, pulse: step === 4 }">
        <rect :x="browserNode.x - browserNode.w/2" :y="browserNode.y - browserNode.h/2"
          :width="browserNode.w" :height="browserNode.h" rx="14" class="node-box" />
        <use href="#ico-browser" :x="browserNode.x - 13" :y="browserNode.y - 20" width="26" height="26" class="node-ico" />
        <text :x="browserNode.x" :y="browserNode.y + 22" class="node-name">Browser</text>
      </g>

      <!-- Device silhouettes under Browser -->
      <g class="devices" :class="{ show: step >= 4 }">
        <use href="#ico-laptop" :x="browserNode.x - 24" :y="browserNode.y + browserNode.h/2 + 7" width="14" height="14" class="device-ico" />
        <use href="#ico-phone" :x="browserNode.x - 5" :y="browserNode.y + browserNode.h/2 + 7" width="10" height="14" class="device-ico" />
        <use href="#ico-tablet" :x="browserNode.x + 10" :y="browserNode.y + browserNode.h/2 + 7" width="14" height="14" class="device-ico" />
        <text :x="browserNode.x" :y="browserNode.y + browserNode.h/2 + 32" class="slogan">Any service, any device</text>
      </g>

      <!-- ── Animated packets (step 4) — per-app colors ── -->
      <template v-if="step >= 4">
        <circle v-for="(p, i) in reqPackets" :key="'req-'+i" r="3" class="pkt" :fill="appColors[p.ai]">
          <animateMotion :dur="p.dur + 's'" :begin="p.delay + 's'" repeatCount="indefinite" fill="freeze" :path="p.path" />
          <animate attributeName="opacity" values="0;0.7;0.7;0" keyTimes="0;0.1;0.85;1" :dur="p.dur + 's'" :begin="p.delay + 's'" repeatCount="indefinite" />
        </circle>
        <circle v-for="(p, i) in resPackets" :key="'res-'+i" r="3" class="pkt" :fill="appColors[p.ai]" fill-opacity="0.6">
          <animateMotion :dur="p.dur + 's'" :begin="p.delay + 's'" repeatCount="indefinite" fill="freeze" :path="p.path" />
          <animate attributeName="opacity" values="0;0.7;0.7;0" keyTimes="0;0.1;0.85;1" :dur="p.dur + 's'" :begin="p.delay + 's'" repeatCount="indefinite" />
        </circle>
      </template>
    </svg>

    <!-- ── Step timeline ── -->
    <div class="timeline">
      <div class="timeline-track">
        <div class="timeline-progress" :style="{ width: ((step - 1) / 3 * 100) + '%' }" />
      </div>
      <div class="timeline-steps">
        <button v-for="(s, i) in steps" :key="i"
          class="timeline-step" :class="{ active: step === i + 1, done: step > i + 1 }"
          @click="goStep(i + 1)">
          <span class="step-dot">{{ i + 1 }}</span>
          <span class="step-title">{{ s.title }}</span>
          <span class="step-desc">{{ s.desc }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const W = 760
const H = 240
const step = ref(1)
const playing = ref(false)
let timer = null

const steps = [
  { title: 'Install', desc: 'One binary, any platform' },
  { title: 'Add Service', desc: 'SSH, DB, Desktop, Web app...' },
  { title: 'Connect', desc: 'Secure tunnel, get a URL' },
  { title: 'Access', desc: 'Open in browser, any device' },
]

// ── Layout ──
const CY = 128 // vertical center for main nodes

// Left: 3 stacked app nodes — scene-oriented labels
const apps = [
  { x: W * 0.09, y: CY - 52, w: 86, h: 46, ico: 'ico-ssh', label: 'SSH Terminal' },
  { x: W * 0.09, y: CY,      w: 86, h: 46, ico: 'ico-db',  label: 'DB Admin' },
  { x: W * 0.09, y: CY + 52, w: 86, h: 46, ico: 'ico-web', label: 'Desktop / Web' },
]

const shieldNode  = { x: W * 0.28, y: CY, w: 116, h: 88 }
const webgateNode = { x: W * 0.70, y: CY, w: 100, h: 78 }
const browserNode = { x: W * 0.91, y: CY, w: 90, h: 78 }

// Tunnel geometry — 3 independent lanes
const tunnelX1 = W * 0.43
const tunnelX2 = W * 0.57
const tunnelHalf = 11 // half-height of the lane spread
const laneYs = [CY - tunnelHalf, CY, CY + tunnelHalf] // SSH, MySQL, RDP lane Y positions

// ── Packet paths ──
function lp(x1, y1, x2, y2) { return `M${x1},${y1} L${x2},${y2}` }

// App colors: SSH=blue, MySQL=orange, RDP=green
const appColors = ['#6366f1', '#f59e0b', '#10b981']

// Small stagger between apps so colors don't overlap perfectly
const APP_OFFSET = 0.3

// Each app gets its own full path with color index (ai)
// Request: browser → webgate → tunnel lane → shield → app
const reqPackets = computed(() =>
  apps.map((app, i) => {
    const ly = laneYs[i]
    const d = i * APP_OFFSET // tiny per-app stagger
    const sx = shieldNode.x, sw = shieldNode.w, sy = shieldNode.y
    const wx = webgateNode.x, ww = webgateNode.w, wy = webgateNode.y
    const bx = browserNode.x, bw = browserNode.w, by = browserNode.y
    return [
      { ai: i, path: lp(bx - bw/2 - 2, by, wx + ww/2 + 2, wy), dur: 2.4, delay: d },
      { ai: i, path: lp(wx - ww/2 - 2, wy, tunnelX2, ly), dur: 1.6, delay: d + 1.2 },
      { ai: i, path: lp(tunnelX2, ly, tunnelX1, ly), dur: 1.6, delay: d + 2.0 },
      { ai: i, path: lp(tunnelX1, ly, sx + sw/2 + 2, sy), dur: 1.2, delay: d + 2.8 },
      { ai: i, path: lp(sx - sw/2 - 2, sy, app.x + app.w/2 + 1, app.y), dur: 1.2, delay: d + 3.4 },
    ]
  }).flat()
)

// Response: app → shield → tunnel lane → webgate → browser
const resPackets = computed(() =>
  apps.map((app, i) => {
    const ly = laneYs[i]
    const d = i * APP_OFFSET
    const sx = shieldNode.x, sw = shieldNode.w, sy = shieldNode.y
    const wx = webgateNode.x, ww = webgateNode.w, wy = webgateNode.y
    const bx = browserNode.x, bw = browserNode.w, by = browserNode.y
    return [
      { ai: i, path: lp(app.x + app.w/2 + 1, app.y, sx - sw/2 - 2, sy), dur: 1.2, delay: d + 4.6 },
      { ai: i, path: lp(sx + sw/2 + 2, sy, tunnelX1, ly), dur: 1.2, delay: d + 5.4 },
      { ai: i, path: lp(tunnelX1, ly, tunnelX2, ly), dur: 1.6, delay: d + 6.0 },
      { ai: i, path: lp(tunnelX2, ly, wx - ww/2 - 2, wy), dur: 1.6, delay: d + 7.0 },
      { ai: i, path: lp(wx + ww/2 + 2, wy, bx - bw/2 - 2, by), dur: 2.4, delay: d + 8.0 },
    ]
  }).flat()
)

function goStep(s) {
  step.value = s
  if (playing.value) { playing.value = false; stopAutoPlay() }
}

function startAutoPlay() {
  stopAutoPlay()
  timer = setInterval(() => {
    step.value = step.value >= 4 ? 1 : step.value + 1
  }, 3000)
}
function stopAutoPlay() {
  if (timer) { clearInterval(timer); timer = null }
}

onMounted(() => { playing.value = true; startAutoPlay() })
onUnmounted(() => stopAutoPlay())
</script>

<style scoped>
.arch {
  --c-brand: var(--vp-c-brand-1, #6366f1);
  --c-brand-soft: var(--vp-c-brand-soft, rgba(99, 102, 241, 0.14));
  margin: 28px 0 8px;
  user-select: none;
}
.arch svg { width: 100%; height: auto; display: block; }

/* ── Zones ── */
.zone { transition: opacity 0.4s; }
.zone--lan { fill: #eafaf2; }
.zone--wan { fill: #eef0fc; }
:global(.dark) .zone--lan { fill: rgba(52, 211, 153, 0.06); }
:global(.dark) .zone--wan { fill: rgba(129, 140, 248, 0.06); }

.zone-label {
  font-size: 10px; font-weight: 700; letter-spacing: 1.2px;
  text-transform: uppercase; fill: var(--vp-c-text-3); text-anchor: middle;
}

/* ── Nodes ── */
.node { opacity: 0.18; transition: opacity 0.5s; }
.node.active { opacity: 1; }

.node-box {
  fill: var(--vp-c-bg); stroke: var(--vp-c-divider); stroke-width: 1.5;
  transition: stroke 0.4s, filter 0.4s;
}
.node.active .node-box { stroke: var(--c-brand); stroke-width: 2; }
.node.pulse .node-box { animation: pulse 1.5s ease-in-out; }

@keyframes pulse {
  0%, 100% { filter: drop-shadow(0 0 0 var(--c-brand)); }
  50% { filter: drop-shadow(0 0 12px var(--c-brand)); }
}

.node-ico { color: var(--vp-c-text-3); transition: color 0.4s; }
.node.active .node-ico { color: var(--c-brand); }

.node-name {
  font-size: 11px; font-weight: 700; text-anchor: middle; fill: var(--vp-c-text-1);
}
.node-name-sm {
  font-size: 9.5px; font-weight: 600; text-anchor: middle; fill: var(--vp-c-text-2);
}

/* ── Hero node (Shield CLI) ── */
.hero-box {
  fill: url(#shieldGrad); stroke-width: 2.5;
}
.node--hero.active .hero-box {
  stroke: var(--c-brand); stroke-width: 2.5;
  filter: url(#shieldGlow);
}
.node--hero.active .node-ico { color: var(--c-brand); transform-origin: center; }
.hero-name { font-size: 13px; font-weight: 800; fill: var(--c-brand); }
.hero-logo { border-radius: 6px; }

/* ── Links ── */
.link {
  stroke: var(--vp-c-divider); stroke-width: 1.5; stroke-dasharray: 5 4;
  opacity: 0.2; transition: all 0.5s;
}
.link.active { stroke: var(--c-brand); stroke-dasharray: none; opacity: 0.55; }

/* ── Tunnel ── */
.tunnel { opacity: 0.15; transition: opacity 0.5s; }
.tunnel.active { opacity: 1; }

.tunnel-line {
  stroke: var(--c-brand); stroke-width: 1.5; stroke-dasharray: 6 4; opacity: 0.35;
}
.tunnel.active .tunnel-line { animation: dash 1.5s linear infinite; }
@keyframes dash { to { stroke-dashoffset: -26; } }

.tunnel-ico { color: var(--vp-c-text-3); transition: color 0.4s; }
.tunnel.active .tunnel-ico { color: var(--c-brand); }

.tunnel-label {
  font-size: 9px; font-weight: 700; text-anchor: middle;
  fill: var(--vp-c-text-2); letter-spacing: 0.5px;
}
.tunnel-sublabel {
  font-size: 7px; font-weight: 500; text-anchor: middle;
  fill: var(--vp-c-text-3); letter-spacing: 0.3px;
}

/* ── URL badge ── */
.url-badge { opacity: 0; transition: opacity 0.5s; }
.url-badge.show { opacity: 1; }
.url-bg { fill: var(--c-brand); opacity: 0.88; }
.url-text {
  font-size: 7.5px; font-weight: 600; font-family: ui-monospace, monospace;
  fill: #fff; text-anchor: middle;
}

/* ── Device icons ── */
.devices { opacity: 0; transition: opacity 0.5s; }
.devices.show { opacity: 1; }
.slogan {
  font-size: 8.5px; font-weight: 600;
  fill: var(--vp-c-text-2); text-anchor: middle; letter-spacing: 0.4px;
}
.device-ico { color: var(--vp-c-text-3); opacity: 0.5; }

/* ── Packets ── */
.pkt { opacity: 0; }

/* ── Timeline ── */
.timeline {
  margin-top: 20px; padding: 16px 4px 0;
  border-top: 1px solid var(--vp-c-divider);
}

.timeline-track {
  height: 3px; background: var(--vp-c-divider); border-radius: 2px;
  margin: 0 48px 12px; overflow: hidden;
}
.timeline-progress {
  height: 100%; background: var(--c-brand); border-radius: 2px;
  transition: width 0.5s ease;
}

.timeline-steps { display: flex; justify-content: space-between; }

.timeline-step {
  display: flex; flex-direction: column; align-items: center; gap: 4px;
  background: none; border: none; cursor: pointer; padding: 4px 8px;
  flex: 1; min-width: 0; transition: opacity 0.3s; opacity: 0.45;
}
.timeline-step.active, .timeline-step.done { opacity: 1; }

.step-dot {
  width: 24px; height: 24px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 11px; font-weight: 700;
  border: 2px solid var(--vp-c-divider); color: var(--vp-c-text-3);
  background: var(--vp-c-bg); transition: all 0.3s; flex-shrink: 0;
}
.timeline-step.active .step-dot {
  background: var(--c-brand); border-color: var(--c-brand); color: #fff;
  box-shadow: 0 0 0 4px var(--c-brand-soft);
}
.timeline-step.done .step-dot {
  background: var(--c-brand); border-color: var(--c-brand); color: #fff;
}

.step-title { font-size: 12px; font-weight: 600; color: var(--vp-c-text-1); white-space: nowrap; }
.step-desc { font-size: 10px; color: var(--vp-c-text-3); text-align: center; line-height: 1.3; }

@media (max-width: 640px) {
  .step-desc { display: none; }
  .step-title { font-size: 11px; }
  .timeline-track { margin: 0 24px 10px; }
}
</style>
