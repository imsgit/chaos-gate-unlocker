#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

req=1.6.2
fc=$(fyne-cross version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
[ -n "$fc" ] || { echo "[!] fyne-cross not found in PATH" >&2; exit 1; }
[ "$(printf '%s\n%s\n' "$req" "$fc" | sort -V | head -1)" = "$req" ] || { echo "[!] fyne-cross $fc < $req" >&2; exit 1; }

build=$(grep -oE 'Build *= *[0-9]+' FyneApp.toml | grep -oE '[0-9]+')
echo "=== fyne-cross $fc | build $build (no auto-bump) ==="

declare -A bak=()
trap 'for s in "${!bak[@]}"; do mv -f "${bak[$s]}" "$s"; done' EXIT
swap() { bak["$1"]="$(mktemp)"; cp "$1" "${bak[$1]}"; }

fontdir=vendor/fyne.io/fyne/v2/theme/font
echo "=== Stub unused fonts (italic/bolditalic/mono) ==="
for f in NotoSans-Italic.ttf NotoSans-BoldItalic.ttf DejaVuSansMono-Powerline.ttf; do
	swap "$fontdir/$f"
	cp "$fontdir/InterSymbols-Regular.ttf" "$fontdir/$f"
done

echo "=== Disable system-font scan ==="
fontprod=vendor/fyne.io/fyne/v2/internal/painter/font_prod.go
swap "$fontprod"
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

echo "=== Slim file-open dialog (favorites + top-right buttons) ==="
filego=vendor/fyne.io/fyne/v2/dialog/file.go
swap "$filego"
perl -0777 -i -pe '
my $a = "\theader := container.NewBorder(\n\t\tnil, nil, nil, optionsbuttons,\n\t\tf.title,\n\t)";
my $b = "\t_ = optionsbuttons\n\theader := container.NewBorder(\n\t\tnil, nil, nil, nil,\n\t\tf.title,\n\t)";
my $i = index($_, $a); die "file.go header block not found\n" if $i < 0; substr($_, $i, length($a)) = $b;' "$filego"
perl -0777 -i -pe '
my $a = "\tbody := container.NewHSplit(\n\t\tf.favoritesList,\n\t\tcontainer.NewBorder(\n\t\t\tf.breadcrumbScroll, nil, nil, nil,\n\t\t\tf.filesScroll,\n\t\t),\n\t)\n\tbody.SetOffset(0) // Set the minimum offset so that the favoritesList takes only its minimal width";
my $b = "\tbody := container.NewBorder(\n\t\tf.breadcrumbScroll, nil, nil, nil,\n\t\tf.filesScroll,\n\t)";
my $i = index($_, $a); die "file.go body block not found\n" if $i < 0; substr($_, $i, length($a)) = $b;' "$filego"

# Embedded browser build — disabled for now; browser version is hosted on Pages.
# Re-enable by uncommenting these and adding embedwasm to the -tags below.
# echo "=== Generate embedded wasm bundle ==="
# fyne package -os wasm
# gzip -9 -f wasm/ChaosGateUnlocker.wasm

for os in windows linux; do
	echo "=== Build $os/amd64 ==="
	fyne-cross "$os" -arch=amd64 -app-build "$build" -tags no_emoji
done
sed -i -E "s/^(\s*Build *= *).*/\1${build}/" FyneApp.toml

strip --strip-all fyne-cross/bin/linux-amd64/chaos-gate-unlocker
rm -rf fyne-cross/dist
sha256sum fyne-cross/bin/windows-amd64/ChaosGateUnlocker.exe \
          fyne-cross/bin/linux-amd64/chaos-gate-unlocker > fyne-cross/bin/SHA256SUMS

echo "=== Done ==="
ls -la fyne-cross/bin/windows-amd64/ fyne-cross/bin/linux-amd64/
