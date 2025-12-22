#!/bin/bash

# Filter Android logs to show only Flutter/app logs, excluding system noise
# Usage: ./scripts/filter_logs.sh

# Prefer wireless device (IP address), fall back to first device if no wireless
DEVICE=$(adb devices | grep -E "device$" | awk '{print $1}' | grep -E "^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+" | head -1)

if [ -z "$DEVICE" ]; then
    # Fall back to first device if no wireless device found
    DEVICE=$(adb devices | grep -E "device$" | head -1 | awk '{print $1}')
fi

if [ -z "$DEVICE" ]; then
    echo "Error: No device found. Please connect a device or start an emulator."
    exit 1
fi

echo "Using device: $DEVICE"
echo "Clearing logcat buffer..."
adb -s "$DEVICE" logcat -c

echo "Showing filtered logs (your app logs only)..."
echo "Press Ctrl+C to stop"
echo ""

# Show only Flutter logs, then filter out common Android system noise
adb -s "$DEVICE" logcat flutter:* *:S | grep -v -E "(ACCESSIBILITY_EVENT|VRI\[|InsetsController|SurfaceView|BLASTBufferQueue|InputMethodManager|InsetsSourceConsumer|ImeTracker|Choreographer|HWUI|CacheManager|InputTransport|ViewRootImpl|SV\[|ThreadedRenderer|dRenderer)"

