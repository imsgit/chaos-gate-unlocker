#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

build=$(grep -oE 'Build *= *[0-9]+' FyneApp.toml | grep -oE '[0-9]+')
echo "=== fyne package wasm | build $build (no auto-bump) ==="
fyne package -os wasm --app-build "$build" --tags no_emoji
sed -i -E "s/^(\s*Build *= *).*/\1${build}/" FyneApp.toml

sed -i -E "s#(application-version\">v[0-9.]+)<#\1.$build<#" wasm/index.html
grep -qE "application-version\">v[0-9.]+\.$build<" wasm/index.html || { echo "[!] version patch failed"; exit 1; }

sed -i '/application-name/d' wasm/index.html
! grep -q application-name wasm/index.html || { echo "[!] app-name strip failed"; exit 1; }

sed -i 's/max-width: 130px;/max-width: 180px;/; s/max-height: 130px;/max-height: 180px;/' wasm/index.html
grep -qF 'max-width: 180px;' wasm/index.html || { echo "[!] splash logo resize failed"; exit 1; }

sed -i 's#<input id="dummyEntry"#<input id="dummyEntry" inputmode="none" readonly#' wasm/index.html
grep -qF 'id="dummyEntry" inputmode="none" readonly' wasm/index.html || { echo "[!] dummyEntry keyboard suppress failed"; exit 1; }

gzip -9 -f wasm/ChaosGateUnlocker.wasm
sed -i 's#fetch("ChaosGateUnlocker.wasm")#fetch("ChaosGateUnlocker.wasm.gz").then(r=>new Response(r.body.pipeThrough(new DecompressionStream("gzip")),{headers:{"Content-Type":"application/wasm"}}))#' wasm/index.html
grep -q DecompressionStream wasm/index.html || { echo "[!] stream patch failed"; exit 1; }

sed -i '/webgl-debug\.js/d' wasm/index.html
! grep -q webgl-debug wasm/index.html || { echo "[!] webgl-debug strip failed"; exit 1; }
rm -f wasm/webgl-debug.js

sed -i 's#<meta charset="utf-8">#<meta charset="utf-8"><script>(function(){var dpr=2;try{var p=window.parent;var s=Math.min(p.innerWidth/800,p.innerHeight/600);dpr=Math.min(s*(p.devicePixelRatio||1),3);}catch(e){}Object.defineProperty(window,"devicePixelRatio",{configurable:true,get:function(){return dpr;}});window.__setDPR=function(v){v=Math.min(Math.max(v,0.5),3);if(Math.abs(v-dpr)<0.01)return;dpr=v;window.dispatchEvent(new Event("resize"));};Object.defineProperty(navigator,"userAgent",{configurable:true,get:function(){return "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36";}});})();</script>#' wasm/index.html
grep -qF 'window.__setDPR' wasm/index.html || { echo "[!] dynamic DPR inject failed"; exit 1; }
grep -qF 'Chrome/120.0.0.0' wasm/index.html || { echo "[!] desktop UA inject failed"; exit 1; }

sed -i 's|<style>|<style>html,body{background-color:#151515}@media (prefers-color-scheme: light){html,body{background-color:#fff}}|' wasm/index.html
grep -qF 'html,body{background-color:#151515}' wasm/index.html || { echo "[!] splash bg inject failed"; exit 1; }
sed -i 's/#141415/#151515/g' wasm/dark.css
grep -qF '#151515' wasm/dark.css || { echo "[!] dark.css bg patch failed"; exit 1; }

mv wasm/index.html wasm/app.html
cp .github/pages-index.html wasm/index.html

echo "=== done ==="
ls -la wasm/
