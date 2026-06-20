read_build() { grep -oE 'Build *= *[0-9]+' FyneApp.toml | grep -oE '[0-9]+'; }
write_build() { sed -i -E "s/^(\s*Build *= *).*/\1${1}/" FyneApp.toml; }

declare -A _bak=()
swap() { _bak["$1"]="$(mktemp)"; cp "$1" "${_bak[$1]}"; }
restore_swaps() { for s in "${!_bak[@]}"; do mv -f "${_bak[$s]}" "$s"; done; }

have() { grep -qF "$2" "$1" || { echo "[!] expected '$2' in $1" >&2; exit 1; }; }
gone() { ! grep -qF "$2" "$1" || { echo "[!] '$2' still present in $1" >&2; exit 1; }; }

stub_fonts() {
	local fontdir=vendor/fyne.io/fyne/v2/theme/font f
	echo "=== Stub unused fonts (italic/bolditalic/mono) ==="
	for f in NotoSans-Italic.ttf NotoSans-BoldItalic.ttf DejaVuSansMono-Powerline.ttf; do
		swap "$fontdir/$f"
		cp "$fontdir/InterSymbols-Regular.ttf" "$fontdir/$f"
	done

	echo "=== Disable system-font scan ==="
	local fontprod=vendor/fyne.io/fyne/v2/internal/painter/font_prod.go
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
}
