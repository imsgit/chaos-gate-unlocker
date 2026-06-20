#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

build=$(grep -oE 'Build *= *[0-9]+' FyneApp.toml | grep -oE '[0-9]+')
echo "=== fyne package wasm | build $build (no auto-bump) ==="
fyne package -os wasm --app-build "$build"
sed -i -E "s/^(\s*Build *= *).*/\1${build}/" FyneApp.toml

sed -i -E "s#(application-version\">v[0-9.]+)<#\1.$build<#" wasm/index.html
grep -qE "application-version\">v[0-9.]+\.$build<" wasm/index.html || { echo "[!] version patch failed"; exit 1; }

gzip -9 -f wasm/ChaosGateUnlocker.wasm
sed -i 's#fetch("ChaosGateUnlocker.wasm")#fetch("ChaosGateUnlocker.wasm.gz").then(r=>new Response(r.body.pipeThrough(new DecompressionStream("gzip")),{headers:{"Content-Type":"application/wasm"}}))#' wasm/index.html
grep -q DecompressionStream wasm/index.html || { echo "[!] stream patch failed"; exit 1; }

sed -i '/webgl-debug\.js/d' wasm/index.html
! grep -q webgl-debug wasm/index.html || { echo "[!] webgl-debug strip failed"; exit 1; }
rm -f wasm/webgl-debug.js

sed -i 's#<meta charset="utf-8">#<meta charset="utf-8"><script>(function(){var r=window.devicePixelRatio||1;Object.defineProperty(window,"devicePixelRatio",{configurable:true,get:function(){return r>2?2:r;}});})();</script>#' wasm/index.html
grep -qF 'return r>2?2:r' wasm/index.html || { echo "[!] DPR cap inject failed"; exit 1; }

sed -i 's|<style>|<style>html,body{background-color:#151515}@media (prefers-color-scheme: light){html,body{background-color:#fff}}|' wasm/index.html
grep -qF 'html,body{background-color:#151515}' wasm/index.html || { echo "[!] splash bg inject failed"; exit 1; }
sed -i 's/#141415/#151515/g' wasm/dark.css
grep -qF '#151515' wasm/dark.css || { echo "[!] dark.css bg patch failed"; exit 1; }

mv wasm/index.html wasm/app.html
cp .github/pages-index.html wasm/index.html

echo "=== done ==="
ls -la wasm/
