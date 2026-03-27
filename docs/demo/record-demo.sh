#!/bin/bash
# Shield CLI Demo Recorder
#
# 录屏 → 按需输出 GIF / 视频 / 两者都要
#
# Usage:
#   ./record-demo.sh <name>                  # 默认: 录屏 → GIF (删除视频)
#   ./record-demo.sh <name> --video          # 录屏 → 保留 MP4 视频
#   ./record-demo.sh <name> --both           # 录屏 → GIF + MP4 都保留
#   ./record-demo.sh <name> --video --fps 30 # 视频模式 + 自定义帧率
#
# Examples:
#   ./record-demo.sh ssh                     # → demo-ssh.gif
#   ./record-demo.sh ssh --video             # → demo-ssh.mp4
#   ./record-demo.sh ssh --both              # → demo-ssh.gif + demo-ssh.mp4
#   ./record-demo.sh tutorial-install --video --fps 30
#
# 操作步骤：
#   1. 运行本脚本，录屏自动开始
#   2. 在终端左半屏执行 shield 命令
#   3. 在浏览器右半屏打开 Access URL
#   4. 操作完成后回到本终端按回车，自动停止录制

set -eo pipefail

# ── 参数解析 ──

NAME="${1:-demo}"
MODE="gif"          # gif | video | both
FPS=""              # 空=自动 (gif:15, video:30)
GIF_WIDTH="960"
SCREEN_DEVICE="2"   # macOS avfoundation: screen 0=2, screen 1=3

shift 2>/dev/null || true
while [[ $# -gt 0 ]]; do
    case "$1" in
        --video)    MODE="video"; shift ;;
        --both)     MODE="both"; shift ;;
        --gif)      MODE="gif"; shift ;;
        --fps)      FPS="$2"; shift 2 ;;
        --width)    GIF_WIDTH="$2"; shift 2 ;;
        --screen)   SCREEN_DEVICE="$2"; shift 2 ;;
        *)          echo "Unknown option: $1"; exit 1 ;;
    esac
done

# 自动帧率: GIF=15 (文件小), 视频=30 (流畅)
if [ -z "$FPS" ]; then
    case "$MODE" in
        gif)   FPS=15 ;;
        video) FPS=30 ;;
        both)  FPS=30 ;;
    esac
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
RAW_FILE="$SCRIPT_DIR/.recording-${NAME}.mov"
GIF_FILE="$SCRIPT_DIR/demo-${NAME}.gif"
MP4_FILE="$SCRIPT_DIR/demo-${NAME}.mp4"

FFMPEG_PID=""

cleanup() {
    [ -n "$FFMPEG_PID" ] && kill -INT "$FFMPEG_PID" 2>/dev/null || true
    wait 2>/dev/null || true
}
trap cleanup EXIT

# ── 提示信息 ──

echo "============================================"
echo "  Shield CLI Demo Recorder"
echo "============================================"
echo ""
echo "  Name:   ${NAME}"
echo "  Mode:   ${MODE}"
echo "  FPS:    ${FPS}"
case "$MODE" in
    gif)   echo "  Output: demo-${NAME}.gif" ;;
    video) echo "  Output: demo-${NAME}.mp4" ;;
    both)  echo "  Output: demo-${NAME}.gif + demo-${NAME}.mp4" ;;
esac
echo ""
echo "  请提前准备好窗口布局："
echo "    左半屏 → 终端 (运行 shield 命令)"
echo "    右半屏 → 浏览器 (打开 Access URL)"
echo ""
echo "============================================"
echo ""
read -p "按回车开始录制 ... " _

echo ""
echo "[REC] 录制中... 完成后回到这里按回车。"
echo ""

# ── 录屏 ──

ffmpeg -y \
    -f avfoundation \
    -capture_cursor 1 \
    -capture_mouse_clicks 1 \
    -pixel_format uyvy422 \
    -framerate "$FPS" \
    -i "$SCREEN_DEVICE" \
    -c:v libx264 -preset ultrafast -crf 18 \
    -pix_fmt yuv420p \
    "$RAW_FILE" > /dev/null 2>&1 &
FFMPEG_PID=$!
sleep 1

if ! kill -0 "$FFMPEG_PID" 2>/dev/null; then
    echo "[ERR] ffmpeg 启动失败"
    echo "      请在 系统设置 > 隐私与安全 > 屏幕录制 中授权终端。"
    exit 1
fi

read -p "按回车停止录制 ... " _

kill -INT "$FFMPEG_PID" 2>/dev/null || true
sleep 2
wait "$FFMPEG_PID" 2>/dev/null || true
FFMPEG_PID=""

echo ""

if [ ! -s "$RAW_FILE" ]; then
    echo "[ERR] 视频文件未生成"
    exit 1
fi
echo "[OK]  录制完成: $(du -h "$RAW_FILE" | cut -f1)"

# ── 后处理 ──

# GIF 转换
if [ "$MODE" = "gif" ] || [ "$MODE" = "both" ]; then
    echo "[...] 转换 GIF 中..."
    ffmpeg -y -i "$RAW_FILE" \
        -vf "fps=12,scale=${GIF_WIDTH}:-1:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=128[p];[s1][p]paletteuse=dither=sierra2_4a" \
        "$GIF_FILE" 2>/dev/null
    echo "[OK]  GIF: $GIF_FILE ($(du -h "$GIF_FILE" | cut -f1))"
fi

# MP4 转换 (重编码为兼容格式, 适合剪辑导入)
if [ "$MODE" = "video" ] || [ "$MODE" = "both" ]; then
    echo "[...] 转换 MP4 中..."
    ffmpeg -y -i "$RAW_FILE" \
        -c:v libx264 -preset medium -crf 20 \
        -pix_fmt yuv420p \
        -movflags +faststart \
        "$MP4_FILE" 2>/dev/null
    echo "[OK]  MP4: $MP4_FILE ($(du -h "$MP4_FILE" | cut -f1))"
fi

# 清理原始录制文件
rm -f "$RAW_FILE"

# ── 结果 ──

echo ""
echo "============================================"
echo "  Done!"
case "$MODE" in
    gif)   echo "  $GIF_FILE" ;;
    video) echo "  $MP4_FILE" ;;
    both)  echo "  $GIF_FILE"
           echo "  $MP4_FILE" ;;
esac
echo "============================================"
