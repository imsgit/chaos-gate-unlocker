#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
source build/lib.sh

ver=$(read_version)

build_linux() {
	local out=fyne-cross/bin/linux-amd64/chaos-gate-unlocker-launcher
	echo "=== launcher v$ver → $out (linux/amd64) ==="
	mkdir -p "$(dirname "$out")"
	CGO_ENABLED=1 go build -trimpath \
		-ldflags "-s -w -X main.version=$ver" \
		-o "$out" ./cmd/launcher
	ls -lh "$out"
}

build_windows() {
	local out=fyne-cross/bin/windows-amd64/ChaosGateUnlockerLauncher.exe
	echo "=== launcher v$ver → $out (windows/amd64) ==="
	mkdir -p "$(dirname "$out")"
	CGO_ENABLED=1 \
		CGO_CFLAGS="-I$PWD/build/winsdk" \
		CGO_CXXFLAGS="-I$PWD/build/winsdk" \
		go build -trimpath \
		-ldflags "-s -w -H windowsgui -X main.version=$ver" \
		-o "$out" ./cmd/launcher
	ls -lh "$out"
}

targets=("$@")
if [ ${#targets[@]} -eq 0 ]; then
	case "$(uname -s)" in
	MINGW* | MSYS* | CYGWIN* | Windows_NT) targets=(windows) ;;
	*) targets=(linux) ;;
	esac
fi

for t in "${targets[@]}"; do
	case $t in
	linux) build_linux ;;
	windows) build_windows ;;
	*)
		echo "[!] unknown target: $t (use linux and/or windows)" >&2
		exit 1
		;;
	esac
done

echo "=== Done ==="
