#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

required_fyne_cross="1.6.2"
fc_version=$(fyne-cross version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
if [ -z "$fc_version" ]; then
	echo "[!] fyne-cross not found in PATH" >&2
	exit 1
fi
if [ "$(printf '%s\n%s\n' "$required_fyne_cross" "$fc_version" | sort -V | head -1)" != "$required_fyne_cross" ]; then
	echo "[!] fyne-cross $fc_version is older than required $required_fyne_cross" >&2
	exit 1
fi
echo "=== fyne-cross $fc_version (>= $required_fyne_cross) ==="

build=$(grep -oE 'Build *= *[0-9]+' FyneApp.toml | grep -oE '[0-9]+')
echo "=== Build number pinned to $build (no auto-bump) ==="

fontdir="vendor/fyne.io/fyne/v2/theme/font"
declare -A fontbak=()
cleanup() {
	for src in "${!fontbak[@]}"; do mv -f "${fontbak[$src]}" "$src"; done
}
trap cleanup EXIT

echo "=== Stub unused vendored fonts (italic/bolditalic/mono) ==="
for f in NotoSans-Italic.ttf NotoSans-BoldItalic.ttf DejaVuSansMono-Powerline.ttf; do
	fontbak["$fontdir/$f"]="$(mktemp)"
	cp "$fontdir/$f" "${fontbak[$fontdir/$f]}"
	cp "$fontdir/InterSymbols-Regular.ttf" "$fontdir/$f"
done

echo "=== Disable system-font fallback scan ==="
fontprod="vendor/fyne.io/fyne/v2/internal/painter/font_prod.go"
fontbak["$fontprod"]="$(mktemp)"
cp "$fontprod" "${fontbak[$fontprod]}"
cat > "$fontprod" <<'GOEOF'
//go:build !test

package painter

import (
	"errors"

	"github.com/go-text/typesetting/fontscan"
)

func loadSystemFonts(_ *fontscan.FontMap) error {
	return errors.New("system fonts disabled")
}
GOEOF

echo "=== Windows build (amd64) ==="
fyne-cross windows -arch=amd64 -app-build "$build" -tags no_emoji

echo "=== Linux build (amd64) ==="
fyne-cross linux -arch=amd64 -app-build "$build" -tags no_emoji

sed -i -E "s/^(\s*Build *= *).*/\1${build}/" FyneApp.toml

echo "=== Strip Linux binary ==="
strip --strip-all fyne-cross/bin/linux-amd64/chaos-gate-unlocker

echo "=== Checksums ==="
win_bin="fyne-cross/bin/windows-amd64/ChaosGateUnlocker.exe"
lin_bin="fyne-cross/bin/linux-amd64/chaos-gate-unlocker"
sums="fyne-cross/bin/SHA256SUMS"
sha256sum "$win_bin" "$lin_bin" > "$sums"
echo "[i] wrote $sums"

echo "=== Done ==="
file fyne-cross/bin/windows-amd64/* fyne-cross/bin/linux-amd64/*
ls -la fyne-cross/bin/windows-amd64/ fyne-cross/bin/linux-amd64/
