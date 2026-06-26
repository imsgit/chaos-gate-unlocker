#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
source build/lib.sh

ver=$(read_version)
tags=${TAGS:-}

build_one() {
	local goos=$1 out=$2 cc=$3 ldextra=$4
	echo "=== launcher v$ver → $out ($goos/amd64) ==="
	mkdir -p "$(dirname "$out")"
	GOOS="$goos" GOARCH=amd64 CGO_ENABLED=1 ${cc:+CC="$cc"} \
		go build -trimpath -ldflags "-s -w -X main.version=$ver $ldextra" \
		${tags:+-tags "$tags"} -o "$out" ./cmd/launcher
	ls -lh "$out"
}

targets=("$@")
[ ${#targets[@]} -eq 0 ] && targets=(linux windows)

for t in "${targets[@]}"; do
	case $t in
	linux)
		build_one linux fyne-cross/bin/linux-amd64/chaos-gate-unlocker-launcher "${LINUX_CC:-}" ""
		;;
	windows)
		build_one windows fyne-cross/bin/windows-amd64/ChaosGateUnlockerLauncher.exe "${WIN_CC:-x86_64-w64-mingw32-gcc}" "-H windowsgui"
		;;
	*)
		echo "[!] unknown target: $t (use linux and/or windows)" >&2
		exit 1
		;;
	esac
done

echo "=== Done ==="
