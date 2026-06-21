read_build() { grep -oE 'Build *= *[0-9]+' FyneApp.toml | grep -oE '[0-9]+'; }
write_build() { sed -i -E "s/^(\s*Build *= *).*/\1${1}/" FyneApp.toml; }

declare -A _bak=()
swap() { _bak["$1"]="$(mktemp)"; cp "$1" "${_bak[$1]}"; }
restore_swaps() { for s in "${!_bak[@]}"; do mv -f "${_bak[$s]}" "$s"; done; }

have() { grep -qF "$2" "$1" || { echo "[!] expected '$2' in $1" >&2; exit 1; }; }
gone() { ! grep -qF "$2" "$1" || { echo "[!] '$2' still present in $1" >&2; exit 1; }; }

sub() { sed -i "$2" "$1"; have "$1" "$3"; }
del() { sed -i "$2" "$1"; gone "$1" "$3"; }

FONT_SUBSET_RANGES="U+0000-00FF,U+0100-017F,U+0400-04FF,U+2010-2027,U+2030-205E,U+20A0-20BF,U+2116,U+2122,U+2026"

stub_fonts() {
	local fontdir=vendor/fyne.io/fyne/v2/theme/font f
	echo "=== Stub unused fonts (italic/bolditalic/mono) ==="
	for f in NotoSans-Italic.ttf NotoSans-BoldItalic.ttf DejaVuSansMono-Powerline.ttf; do
		swap "$fontdir/$f"
		cp "$fontdir/InterSymbols-Regular.ttf" "$fontdir/$f"
	done

	echo "=== Subset primary fonts to Latin+Cyrillic (~910KB -> ~130KB) ==="
	if command -v pyftsubset >/dev/null; then
		for f in NotoSans-Regular.ttf NotoSans-Bold.ttf; do
			swap "$fontdir/$f"
			pyftsubset "${_bak[$fontdir/$f]}" --output-file="$fontdir/$f" \
				--unicodes="$FONT_SUBSET_RANGES" --no-hinting --desubroutinize
		done
	else
		echo "[!] pyftsubset not found — skipping subset (install: pip install fonttools)"
	fi

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

slim_charset() {
	local cf=vendor/golang.org/x/net/html/charset/charset.go
	echo "=== Slim SVG charset reader (UTF-8 only; drops x/text CJK tables) ==="
	swap "$cf"
	cat > "$cf" <<'GOEOF'
package charset // import "golang.org/x/net/html/charset"

import "io"

func NewReaderLabel(label string, input io.Reader) (io.Reader, error) {
	return input, nil
}
GOEOF
}
