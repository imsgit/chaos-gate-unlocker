#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
source ./lib.sh

build=$(read_build)
echo "=== fyne package wasm | build $build (no auto-bump) ==="

created_vendor=
[ -d vendor ] || { echo "=== go mod vendor (for font stubbing) ==="; go mod vendor; created_vendor=1; }
trap 'restore_swaps; [ -n "$created_vendor" ] && rm -rf vendor' EXIT

stub_fonts

fyne package -os wasm --app-build "$build" --tags no_emoji
write_build "$build"

idx=wasm/index.html

sed -i -E "s#(application-version\">v[0-9.]+)<#\1.$build<#" "$idx"
grep -qE "application-version\">v[0-9.]+\.$build<" "$idx" || { echo "[!] version patch failed"; exit 1; }
sed -i '/application-name/d' "$idx";                              gone "$idx" application-name
sed -i 's/max-width: 130px;/max-width: 180px;/; s/max-height: 130px;/max-height: 180px;/' "$idx"
have "$idx" 'max-width: 180px;'

sed -i 's#<input id="dummyEntry"#<input id="dummyEntry" inputmode="none" readonly#' "$idx"
have "$idx" 'id="dummyEntry" inputmode="none" readonly'

echo "=== Strip wasm 'name' custom section ==="
python3 - wasm/ChaosGateUnlocker.wasm <<'PY'
import sys
p = sys.argv[1]
d = open(p, 'rb').read()
assert d[:4] == b'\x00asm', "not a wasm file"
def rd(b, i):
	r = s = 0
	while True:
		x = b[i]; i += 1; r |= (x & 0x7f) << s
		if not x & 0x80: break
		s += 7
	return r, i
out = bytearray(d[:8]); i = 8; dropped = 0
while i < len(d):
	sid = d[i]
	sz, k = rd(d, i + 1)
	end = k + sz
	if sid == 0:
		nl, m = rd(d, k)
		if d[m:m+nl] == b'name':
			dropped += end - i; i = end; continue
	out += d[i:end]; i = end
open(p, 'wb').write(out)
print(f"  stripped {dropped} bytes")
PY

gzip -9 -f wasm/ChaosGateUnlocker.wasm
sed -i 's#fetch("ChaosGateUnlocker.wasm")#fetch("ChaosGateUnlocker.wasm.gz").then(r=>new Response(r.body.pipeThrough(new DecompressionStream("gzip")),{headers:{"Content-Type":"application/wasm"}}))#' "$idx"
have "$idx" DecompressionStream

sed -i '/webgl-debug\.js/d' "$idx"; gone "$idx" webgl-debug
rm -f wasm/webgl-debug.js

sed -i 's#<meta charset="utf-8">#<meta charset="utf-8"><script>(function(){var dpr=2;try{var p=window.parent;var s=Math.min(p.innerWidth/800,p.innerHeight/600);dpr=Math.min(s*(p.devicePixelRatio||1),3);}catch(e){}Object.defineProperty(window,"devicePixelRatio",{configurable:true,get:function(){return dpr;}});window.__setDPR=function(v){v=Math.min(Math.max(v,0.5),3);if(Math.abs(v-dpr)<0.01)return;dpr=v;window.dispatchEvent(new Event("resize"));};Object.defineProperty(navigator,"userAgent",{configurable:true,get:function(){return "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36";}});})();</script>#' "$idx"
have "$idx" 'window.__setDPR'; have "$idx" 'Chrome/120.0.0.0'

sed -i 's|<style>|<style>html,body{background-color:#151515}@media (prefers-color-scheme: light){html,body{background-color:#fff}}|' "$idx"
have "$idx" 'html,body{background-color:#151515}'
sed -i 's/#141415/#151515/g' wasm/dark.css; have wasm/dark.css '#151515'

mv "$idx" wasm/app.html
cp .github/pages-index.html "$idx"

echo "=== done ==="
ls -la wasm/