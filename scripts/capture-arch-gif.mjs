/**
 * Shield CLI Architecture Diagram Capture
 *
 * Usage:
 *   node scripts/capture-arch-gif.mjs            # → architecture.gif
 *   node scripts/capture-arch-gif.mjs --video     # → architecture.mp4
 *   node scripts/capture-arch-gif.mjs --both      # → .gif + .mp4
 *
 * Requires: puppeteer, ffmpeg
 * Requires: VitePress dev server running on localhost:5173
 */

import puppeteer from 'puppeteer'
import { execSync } from 'child_process'
import { mkdirSync, rmSync, statSync } from 'fs'
import path from 'path'

// ── Config ──
const PAGE_URL = 'http://localhost:5173/guide/what-is-shield.html'
const FRAMES_DIR = '/tmp/arch-gif-frames'
const OUTPUT_DIR = path.resolve('docs/demo')
const DURATION_S = 12

// ── Parse args ──
const args = process.argv.slice(2)
const mode = args.includes('--video') ? 'video'
           : args.includes('--both')  ? 'both'
           : 'gif'

// GIF: 12fps (small file), video: 24fps (smooth for editing)
const TARGET_FPS = mode === 'gif' ? 12 : 24
const FRAME_INTERVAL = 1000 / TARGET_FPS

async function main() {
  rmSync(FRAMES_DIR, { recursive: true, force: true })
  mkdirSync(FRAMES_DIR, { recursive: true })

  console.log(`Mode: ${mode} | FPS: ${TARGET_FPS} | Duration: ${DURATION_S}s`)

  // ── Launch browser & navigate ──
  const browser = await puppeteer.launch({ headless: true })
  const page = await browser.newPage()
  await page.setViewport({ width: 900, height: 600, deviceScaleFactor: 2 })
  await page.goto(PAGE_URL, { waitUntil: 'networkidle0' })

  // Activate step 4 (full animation state)
  await page.evaluate(() => {
    const btns = document.querySelectorAll('.timeline-step')
    if (btns[3]) btns[3].click()
  })
  await new Promise(r => setTimeout(r, 1000))

  const archHandle = await page.$('.arch')
  if (!archHandle) {
    console.error('Could not find .arch element')
    await browser.close()
    process.exit(1)
  }

  await archHandle.scrollIntoView()
  await new Promise(r => setTimeout(r, 500))

  // ── Capture frames (wall-clock timing) ──
  const totalFrames = TARGET_FPS * DURATION_S
  console.log(`Capturing ${totalFrames} frames...`)

  const captureStart = Date.now()

  for (let i = 0; i < totalFrames; i++) {
    const frameStart = Date.now()
    const framePath = path.join(FRAMES_DIR, `frame-${String(i).padStart(4, '0')}.png`)
    await archHandle.screenshot({ path: framePath })

    const remaining = FRAME_INTERVAL - (Date.now() - frameStart)
    if (remaining > 0) await new Promise(r => setTimeout(r, remaining))

    if ((i + 1) % (TARGET_FPS * 2) === 0) console.log(`  ${i + 1}/${totalFrames}`)
  }

  const actualDuration = (Date.now() - captureStart) / 1000
  const actualFps = Math.round(totalFrames / actualDuration)
  console.log(`Captured in ${actualDuration.toFixed(1)}s (effective ${actualFps}fps)`)

  await browser.close()

  // ── Generate outputs ──
  const outputs = []

  if (mode === 'gif' || mode === 'both') {
    const gifFile = path.join(OUTPUT_DIR, 'architecture.gif')
    console.log('Generating GIF...')
    const palettePath = path.join(FRAMES_DIR, 'palette.png')
    execSync(
      `ffmpeg -y -framerate ${actualFps} -i "${FRAMES_DIR}/frame-%04d.png" -vf "scale=800:-1:flags=lanczos,palettegen=max_colors=128:stats_mode=diff" -frames:v 1 -update 1 "${palettePath}"`,
      { stdio: 'pipe' }
    )
    execSync(
      `ffmpeg -y -framerate ${actualFps} -i "${FRAMES_DIR}/frame-%04d.png" -i "${palettePath}" -lavfi "scale=800:-1:flags=lanczos [x]; [x][1:v] paletteuse=dither=bayer:bayer_scale=3" -loop 0 "${gifFile}"`,
      { stdio: 'pipe' }
    )
    outputs.push({ file: gifFile, size: (statSync(gifFile).size / 1024).toFixed(0) + ' KB' })
  }

  if (mode === 'video' || mode === 'both') {
    const mp4File = path.join(OUTPUT_DIR, 'architecture.mp4')
    console.log('Generating MP4...')
    execSync(
      `ffmpeg -y -framerate ${actualFps} -i "${FRAMES_DIR}/frame-%04d.png" -c:v libx264 -preset medium -crf 18 -pix_fmt yuv420p -movflags +faststart "${mp4File}"`,
      { stdio: 'pipe' }
    )
    outputs.push({ file: mp4File, size: (statSync(mp4File).size / 1024).toFixed(0) + ' KB' })
  }

  // ── Cleanup & report ──
  rmSync(FRAMES_DIR, { recursive: true, force: true })

  console.log('\nDone!')
  for (const o of outputs) {
    console.log(`  ${o.file} (${o.size})`)
  }
}

main().catch(e => { console.error(e); process.exit(1) })
