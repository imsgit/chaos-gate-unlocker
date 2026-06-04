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

echo "=== Windows build (amd64) ==="
fyne-cross windows -arch=amd64 -app-build "$build"

echo "=== Linux build (amd64) ==="
fyne-cross linux -arch=amd64 -app-build "$build"

sed -i -E "s/^(\s*Build *= *).*/\1${build}/" FyneApp.toml

echo "=== Strip Linux binary ==="
strip --strip-all fyne-cross/bin/linux-amd64/chaos-gate-unlocker

echo "=== Repackage Linux tar.xz with stripped binary ==="
tarball="$(pwd)/fyne-cross/dist/linux-amd64/ChaosGateUnlocker.tar.xz"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
tar -xJf "$tarball" -C "$work"
find "$work" -type f -exec sh -c 'file -b "$1" | grep -q ELF && strip --strip-all "$1"' _ {} \;
tar -cJf "$tarball" -C "$work" .

echo "=== Checksums ==="
win_zip="fyne-cross/dist/windows-amd64/ChaosGateUnlocker.exe.zip"
lin_txz="fyne-cross/dist/linux-amd64/ChaosGateUnlocker.tar.xz"
sums="fyne-cross/dist/SHA256SUMS"
sha256sum "$win_zip" "$lin_txz" > "$sums"
echo "[i] wrote $sums"

echo "=== Done ==="
file fyne-cross/bin/windows-amd64/* fyne-cross/bin/linux-amd64/*
ls -la fyne-cross/bin/windows-amd64/ fyne-cross/bin/linux-amd64/
ls -la fyne-cross/dist/windows-amd64/ fyne-cross/dist/linux-amd64/
