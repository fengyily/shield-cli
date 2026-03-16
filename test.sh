#!/bin/bash

echo "=== Shield CLI Test ==="
echo ""

# Build
echo "[1/3] Building shield..."
go build -o shield . || { echo "Build failed!"; exit 1; }
echo "  OK"

# Start mock API server
echo "[2/3] Starting mock API server on :18080..."
go run test/mock_server.go &
MOCK_PID=$!
sleep 1

# Cleanup on exit
cleanup() {
    echo ""
    echo "Cleaning up..."
    kill $MOCK_PID 2>/dev/null
    wait $MOCK_PID 2>/dev/null
    echo "Done."
}
trap cleanup EXIT

# Run shield
echo "[3/3] Starting shield..."
echo ""
go run ./main.go -t ssh -s 172.16.3.137:22 -H https://console.yishield.com/raas -v
