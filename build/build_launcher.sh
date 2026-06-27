#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
source build/lib.sh

ver=$(read_version)
tags=${TAGS:-}

WIN_IMAGE=${WIN_IMAGE:-fyneio/fyne-cross-images:windows}

build_linux() {
	local out=fyne-cross/bin/linux-amd64/chaos-gate-unlocker-launcher
	echo "=== launcher v$ver → $out (linux/amd64) ==="
	mkdir -p "$(dirname "$out")"
	CGO_ENABLED=1 ${LINUX_CC:+CC="$LINUX_CC"} \
		go build -trimpath -ldflags "-s -w -X main.version=$ver" \
		${tags:+-tags "$tags"} -o "$out" ./cmd/launcher
	ls -lh "$out"
}

build_windows() {
	local out=fyne-cross/bin/windows-amd64/ChaosGateUnlockerLauncher.exe
	echo "=== launcher v$ver → $out (windows/amd64 via $WIN_IMAGE) ==="
	mkdir -p "$(dirname "$out")"
	[ -d vendor ] || go mod vendor
	docker run --rm --user "$(id -u):$(id -g)" \
		-v "$PWD":/src -w /src \
		-e HOME=/tmp -e GOOS=windows -e GOARCH=amd64 -e CGO_ENABLED=1 \
		-e CC="zig cc -target x86_64-windows-gnu" \
		-e CXX="zig c++ -target x86_64-windows-gnu" \
		-e CGO_CFLAGS="-I/src/build/winsdk" -e CGO_CXXFLAGS="-I/src/build/winsdk" \
		-e GOFLAGS=-mod=vendor -e GOCACHE=/tmp/.gocache \
		-e VER="$ver" -e TAGS="$tags" -e OUT="$out" \
		"$WIN_IMAGE" sh -euc '
			export PATH=$PATH:/usr/local/zig
			go build -buildvcs=false -trimpath \
				-ldflags "-s -w -H windowsgui -X main.version=$VER" \
				${TAGS:+-tags "$TAGS"} -o "$OUT" ./cmd/launcher
		'
	ls -lh "$out"
}

targets=("$@")
[ ${#targets[@]} -eq 0 ] && targets=(linux windows)

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
