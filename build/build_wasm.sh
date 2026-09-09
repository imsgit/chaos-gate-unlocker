#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
source build/lib.sh

version=$(read_version)
build=$(read_build)
echo "=== go build wasm | v$version build $build (no auto-bump) ==="

ensure_vendor "font stubbing"

slim_common
enable_touch_scroll

mkdir -p wasm
GOOS=js GOARCH=wasm go build -tags no_emoji -trimpath \
	-ldflags "-s -w -X main.appVersion=$version -X main.appBuild=$build" \
	-o wasm/ChaosGateUnlocker.wasm .

gzip -9 -f wasm/ChaosGateUnlocker.wasm

cp web/index.html web/app.html wasm/
sed -i "s/__VERSION__/$version.$build/" wasm/app.html
have wasm/app.html "v$version.$build"
cp Icon.png wasm/icon.png
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" wasm/

echo "=== done ==="
ls -la wasm/
