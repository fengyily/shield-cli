#!/bin/bash
# Shield CLI Demo Recorder (手动操作版)
#
# 脚本只负责：录屏 → 等你按回车停止 → 转 GIF
# 你负责：在终端执行 shield 命令 + 打开浏览器
#
# Usage:
#   ./record-demo.sh ssh
#   ./record-demo.sh rdp
#   ./record-demo.sh http
#   ./record-demo.sh mysql
#
# 操作步骤：
#   1. 运行本脚本，录屏自动开始
#   2. 在终端左半屏执行 shield 命令
#   3. 在浏览器右半屏打开 Access URL
#   4. 操作完成后回到本终端按回车，自动停止录制并转 GIF

set -eo pipefail

PROTOCOL="${1:-demo}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
VIDEO_FILE="$SCRIPT_DIR/demo-${PROTOCOL}.mov"
GIF_FILE="$SCRIPT_DIR/demo-${PROTOCOL}.gif"

# ffmpeg screen device index (Capture screen 0 = main display)
# 修改为你的主屏设备号: screen 0=2, screen 1=3, screen 2=4
SCREEN_DEVICE="2"

FFMPEG_PID=""

cleanup() {
    [ -n "$FFMPEG_PID" ] && kill -INT "$FFMPEG_PID" 2>/dev/null || true
    wait 2>/dev/null || true
}
trap cleanup EXIT

echo "============================================"
echo "  Shield CLI Demo Recorder"
echo "============================================"
echo ""
echo "  录制目标: demo-${PROTOCOL}.gif"
echo "  录屏设备: Capture screen 0 (device $SCREEN_DEVICE)"
echo ""
echo "  请提前准备好窗口布局："
echo "    左半屏 → 终端 (运行 shield 命令)"
echo "    右半屏 → 浏览器 (打开 Access URL)"
echo ""
echo "============================================"
echo ""
read -p "按回车开始录制 ▶ " _

echo ""
echo "🔴 录制中..."
echo "   请在其他窗口操作，完成后回到这里按回车。"
echo ""

# Start screen recording
ffmpeg -y \
    -f avfoundation \
    -capture_cursor 1 \
    -capture_mouse_clicks 1 \
    -pixel_format uyvy422 \
    -framerate 15 \
    -i "$SCREEN_DEVICE" \
    -c:v libx264 -preset ultrafast -crf 18 \
    -pix_fmt yuv420p \
    "$VIDEO_FILE" > /dev/null 2>&1 &
FFMPEG_PID=$!
sleep 1

if ! kill -0 "$FFMPEG_PID" 2>/dev/null; then
    echo "❌ ffmpeg 启动失败！"
    echo "   请在 系统设置 → 隐私与安全 → 屏幕录制 中授权终端。"
    exit 1
fi

read -p "按回车停止录制 ⏹ " _

# Stop ffmpeg gracefully
kill -INT "$FFMPEG_PID" 2>/dev/null || true
sleep 2
wait "$FFMPEG_PID" 2>/dev/null || true
FFMPEG_PID=""

echo ""

if [ ! -s "$VIDEO_FILE" ]; then
    echo "❌ 视频文件未生成: $VIDEO_FILE"
    exit 1
fi
echo "✅ 视频录制完成: $(du -h "$VIDEO_FILE" | cut -f1)"

# Convert to GIF
echo "🔄 转换 GIF 中..."
ffmpeg -y -i "$VIDEO_FILE" \
    -vf "fps=10,scale=960:-1:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=128[p];[s1][p]paletteuse=dither=sierra2_4a" \
    "$GIF_FILE" 2>/dev/null

rm -f "$VIDEO_FILE"

echo ""
echo "============================================"
echo "  ✅ 完成!"
echo "  GIF: $GIF_FILE"
echo "  大小: $(du -h "$GIF_FILE" | cut -f1)"
echo "============================================"
