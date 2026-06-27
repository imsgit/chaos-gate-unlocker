#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
source build/lib.sh

ver=$(read_version)

created_vendor=
[ -d vendor ] || { echo "=== go mod vendor (for webview.h swap) ==="; go mod vendor; created_vendor=1; }
trap 'restore_swaps; [ -n "$created_vendor" ] && rm -rf vendor' EXIT
hide_webview_window

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

	local env=(CGO_ENABLED=1 CGO_CFLAGS="-I$PWD/build/winsdk" CGO_CXXFLAGS="-I$PWD/build/winsdk")
	case "$(uname -s)" in
	MINGW* | MSYS* | CYGWIN* | Windows_NT) ;;
	*)
		command -v x86_64-w64-mingw32-gcc >/dev/null ||
			{ echo "[!] x86_64-w64-mingw32-gcc not found (needed to cross-build the Windows launcher)" >&2; exit 1; }
		env+=(GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++)
		;;
	esac

	env "${env[@]}" go build -trimpath \
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
