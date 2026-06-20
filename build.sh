#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
source ./lib.sh

req=1.6.2
fc=$(fyne-cross version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
[ -n "$fc" ] || { echo "[!] fyne-cross not found in PATH" >&2; exit 1; }
[ "$(printf '%s\n%s\n' "$req" "$fc" | sort -V | head -1)" = "$req" ] || { echo "[!] fyne-cross $fc < $req" >&2; exit 1; }

build=$(read_build)
echo "=== fyne-cross $fc | build $build (no auto-bump) ==="

tags=no_emoji
if [ -n "${EMBED:-}" ]; then
	echo "=== EMBED=1 → build wasm bundle for offline 'Try it online' ==="
	bash build-wasm.sh
	tags=no_emoji,embedwasm
fi

trap restore_swaps EXIT
stub_fonts

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

for os in windows linux; do
	echo "=== Build $os/amd64 ==="
	fyne-cross "$os" -arch=amd64 -app-build "$build" -tags "$tags"
done
write_build "$build"

strip --strip-all fyne-cross/bin/linux-amd64/chaos-gate-unlocker
rm -rf fyne-cross/dist
sha256sum fyne-cross/bin/windows-amd64/ChaosGateUnlocker.exe \
          fyne-cross/bin/linux-amd64/chaos-gate-unlocker > fyne-cross/bin/SHA256SUMS

echo "=== Done ==="
ls -la fyne-cross/bin/windows-amd64/ fyne-cross/bin/linux-amd64/
