#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

echo "######################## WASM ########################"
build/build_wasm.sh

echo "###################### LAUNCHERS #####################"
build/build_launcher.sh linux windows

echo "######################## DONE ########################"
echo "wasm:      wasm/app.html + wasm/ChaosGateUnlocker.wasm.gz"
echo "launchers: fyne-cross/bin/{linux-amd64,windows-amd64}/"
