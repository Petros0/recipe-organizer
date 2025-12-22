# justfile for Flutter/Dart project
# Expects the following aliases in your .zshrc:
# alias f="flutter"
# alias d="dart"
#
# All recipes run in interactive zsh, so aliases/functions are available.

set shell := ["zsh", "-i", "-c"]

default:
    @just help

help:
    @just --list
    @echo ""
    @echo "Note: This justfile expects 'f' and 'd' aliases for flutter and dart."

l10n:
    @echo "Generating localization files..."
    f gen-l10n --arb-dir="lib/l10n/arb"
    @echo "Localization files generated successfully!"

format:
    @echo "=== Formatting code ==="
    d fix --apply
    d format . --line-length 120
    @echo "Code formatted successfully!"

dev:
    @echo "Running development flavor..."
    f run --flavor development --target lib/main_development.dart

devices:
    @echo "Listing connected devices..."
    f devices 

pair host port connectPort:
    @echo "Pairing device..."
    adb pair {{ host }}:{{ port }}
    @echo "Pairing successful! Ready to connect."
    @just connect {{ host }} {{ connectPort }}
    
connect host port:
    @echo "Connecting device..."
    adb connect {{ host }}:{{ port }}
    @echo "Device connected successfully!"

deploy-recipe-request:
    @echo "Deploying function..."
    appwrite functions create-deployment \
        --function-id=recipe-request \
        --entrypoint=main.go \
        --code=functions/recipe-request --activate true
    @echo "Function deployed successfully!"

deploy-recipe-request-processor:
    @echo "Deploying function..."
    appwrite functions create-deployment \
        --function-id=recipe-request-processor \
        --entrypoint=main.go \
        --code=functions/recipe-request-processor --activate true
    @echo "Function deployed successfully!"

deploy-all:
    @echo "Deploying all functions..."
    @just deploy-recipe-request
    @just deploy-recipe-request-processor
    @echo "All functions deployed successfully!"

logs:
    @echo "Showing filtered logs (your app logs only, excluding system noise)..."
    @echo "Press Ctrl+C to stop"
    DEVICE=`adb devices | grep -E "device$$" | awk '{print $$1}' | grep -E "^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+" | head -1`; \
    if [ -z "$$DEVICE" ]; then \
        DEVICE=`adb devices | grep -E "device$$" | head -1 | awk '{print $$1}'`; \
    fi; \
    if [ -z "$$DEVICE" ]; then \
        echo "Error: No device found. Please connect a device or start an emulator."; \
        exit 1; \
    fi; \
    echo "Using device: $$DEVICE"; \
    adb -s $$DEVICE logcat -c && adb -s $$DEVICE logcat flutter:* *:S | grep -v -E "(ACCESSIBILITY_EVENT|VRI\[|InsetsController|SurfaceView|BLASTBufferQueue|InputMethodManager|InsetsSourceConsumer|ImeTracker|Choreographer|HWUI|CacheManager|InputTransport|ViewRootImpl|SV\[)"

logs-verbose:
    @echo "Showing all Flutter logs (including framework logs)..."
    @echo "Press Ctrl+C to stop"
    DEVICE=`adb devices | grep -E "device$$" | awk '{print $$1}' | grep -E "^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+" | head -1`; \
    if [ -z "$$DEVICE" ]; then \
        DEVICE=`adb devices | grep -E "device$$" | head -1 | awk '{print $$1}'`; \
    fi; \
    if [ -z "$$DEVICE" ]; then \
        echo "Error: No device found. Please connect a device or start an emulator."; \
        exit 1; \
    fi; \
    echo "Using device: $$DEVICE"; \
    adb -s $$DEVICE logcat -c && adb -s $$DEVICE logcat flutter:* *:S

logs-clean:
    @echo "Showing only your app's logs (developer.log output, minimal filtering)..."
    @echo "Press Ctrl+C to stop"
    DEVICE=`adb devices | grep -E "device$$" | awk '{print $$1}' | grep -E "^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+" | head -1`; \
    if [ -z "$$DEVICE" ]; then \
        DEVICE=`adb devices | grep -E "device$$" | head -1 | awk '{print $$1}'`; \
    fi; \
    if [ -z "$$DEVICE" ]; then \
        echo "Error: No device found. Please connect a device or start an emulator."; \
        exit 1; \
    fi; \
    echo "Using device: $$DEVICE"; \
    adb -s $$DEVICE logcat -c && adb -s $$DEVICE logcat flutter:* | grep -E "(flutter|recipe_organizer)" 