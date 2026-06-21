#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
source build/lib.sh

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
del "$idx" '/application-name/d' application-name
sub "$idx" 's/max-width: 130px;/max-width: 160px;/; s/max-height: 130px;/max-height: 160px;/' 'max-width: 160px;'
sub "$idx" 's#<input id="dummyEntry"#<input id="dummyEntry" inputmode="none" readonly#' 'id="dummyEntry" inputmode="none" readonly'

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
sub "$idx" 's#fetch("ChaosGateUnlocker.wasm")#fetch("ChaosGateUnlocker.wasm.gz").then(r=>new Response(r.body.pipeThrough(new DecompressionStream("gzip")),{headers:{"Content-Type":"application/wasm"}}))#' DecompressionStream

del "$idx" '/webgl-debug\.js/d' webgl-debug
rm -f wasm/webgl-debug.js

sub "$idx" 's#<meta charset="utf-8">#<meta charset="utf-8"><script>(function(){var dpr=2;try{var p=window.parent;var s=Math.min(p.innerWidth/800,p.innerHeight/600);dpr=Math.min(s*(p.devicePixelRatio||1),3);}catch(e){}Object.defineProperty(window,"devicePixelRatio",{configurable:true,get:function(){return dpr;}});window.__setDPR=function(v){v=Math.min(Math.max(v,0.5),3);if(Math.abs(v-dpr)<0.01)return;dpr=v;window.dispatchEvent(new Event("resize"));};Object.defineProperty(navigator,"userAgent",{configurable:true,get:function(){return "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36";}});})();</script>#' 'window.__setDPR'
have "$idx" 'Chrome/140.0.0.0'

sub "$idx" 's|<style>|<style>html,body{background-color:#151515}@media (prefers-color-scheme: light){html,body{background-color:#fff}}|' 'html,body{background-color:#151515}'
sub wasm/dark.css 's/#141415/#151515/g' '#151515'

mv "$idx" wasm/app.html
cp .github/pages-index.html "$idx"

echo "=== done ==="
ls -la wasm/