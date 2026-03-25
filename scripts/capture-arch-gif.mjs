import puppeteer from 'puppeteer'
import { execSync } from 'child_process'
import { mkdirSync, rmSync, statSync } from 'fs'
import path from 'path'

const URL = 'http://localhost:5173/guide/what-is-shield.html'
const FRAMES_DIR = '/tmp/arch-gif-frames'
const OUTPUT = path.resolve('docs/demo/architecture.gif')
const TARGET_FPS = 12
const DURATION_S = 12
const FRAME_INTERVAL = 1000 / TARGET_FPS // ms between frames in real time

async function main() {
  rmSync(FRAMES_DIR, { recursive: true, force: true })
  mkdirSync(FRAMES_DIR, { recursive: true })

  const browser = await puppeteer.launch({ headless: true })
  const page = await browser.newPage()
  await page.setViewport({ width: 900, height: 600, deviceScaleFactor: 2 })
  await page.goto(URL, { waitUntil: 'networkidle0' })

  // Click step 4 to activate animations
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

  // Capture frames with wall-clock timing to match real animation speed
  const totalFrames = TARGET_FPS * DURATION_S
  console.log(`Capturing ${totalFrames} frames (${DURATION_S}s at ${TARGET_FPS}fps)...`)

  const captureStart = Date.now()

  for (let i = 0; i < totalFrames; i++) {
    const frameStart = Date.now()
    const framePath = path.join(FRAMES_DIR, `frame-${String(i).padStart(4, '0')}.png`)
    await archHandle.screenshot({ path: framePath })

    // Wait remaining time to hit the target interval
    const elapsed = Date.now() - frameStart
    const remaining = FRAME_INTERVAL - elapsed
    if (remaining > 0) {
      await new Promise(r => setTimeout(r, remaining))
    }

    if ((i + 1) % 24 === 0) console.log(`  ${i + 1}/${totalFrames}`)
  }

  const actualDuration = (Date.now() - captureStart) / 1000
  const actualFps = totalFrames / actualDuration
  console.log(`Captured in ${actualDuration.toFixed(1)}s (effective ${actualFps.toFixed(1)}fps)`)

  await browser.close()

  // Use actual FPS for GIF so playback matches real-time
  const gifFps = Math.round(actualFps)
  console.log(`Generating GIF at ${gifFps}fps...`)

  const palettePath = path.join(FRAMES_DIR, 'palette.png')
  execSync(
    `ffmpeg -y -framerate ${gifFps} -i "${FRAMES_DIR}/frame-%04d.png" -vf "scale=800:-1:flags=lanczos,palettegen=max_colors=128:stats_mode=diff" -frames:v 1 -update 1 "${palettePath}"`,
    { stdio: 'inherit' }
  )
  execSync(
    `ffmpeg -y -framerate ${gifFps} -i "${FRAMES_DIR}/frame-%04d.png" -i "${palettePath}" -lavfi "scale=800:-1:flags=lanczos [x]; [x][1:v] paletteuse=dither=bayer:bayer_scale=3" -loop 0 "${OUTPUT}"`,
    { stdio: 'inherit' }
  )

  rmSync(FRAMES_DIR, { recursive: true, force: true })

  const size = (statSync(OUTPUT).size / 1024).toFixed(0)
  console.log(`Done! ${OUTPUT} (${size} KB)`)
}

main().catch(e => { console.error(e); process.exit(1) })
