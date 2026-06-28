#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
source build/lib.sh

req=1.6.2
fc=$(fyne-cross version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
[ -n "$fc" ] || { echo "[!] fyne-cross not found in PATH" >&2; exit 1; }
[ "$(printf '%s\n%s\n' "$req" "$fc" | sort -V | head -1)" = "$req" ] || { echo "[!] fyne-cross $fc < $req" >&2; exit 1; }

build=$(read_build)
echo "=== fyne-cross $fc | build $build (no auto-bump) ==="

tags=no_emoji

created_vendor=
[ -d vendor ] || { echo "=== go mod vendor (for font/charset/markdown swaps) ==="; go mod vendor; created_vendor=1; }
trap 'restore_swaps; [ -n "$created_vendor" ] && rm -rf vendor' EXIT
stub_fonts
slim_charset
slim_markdown
round_dialogs

for os in windows linux; do
	echo "=== Build $os/amd64 ==="
	name_args=()
	[ "$os" = windows ] && name_args=(-name ChaosGateUnlocker)
	fyne-cross "$os" -arch=amd64 -app-build "$build" -tags "$tags" "${name_args[@]}"
done
write_build "$build"

strip --strip-all fyne-cross/bin/linux-amd64/chaos-gate-unlocker
rm -rf fyne-cross/dist
sha256sum fyne-cross/bin/windows-amd64/ChaosGateUnlocker.exe \
          fyne-cross/bin/linux-amd64/chaos-gate-unlocker > fyne-cross/bin/SHA256SUMS

echo "=== Done ==="
ls -la fyne-cross/bin/windows-amd64/ fyne-cross/bin/linux-amd64/
