read_build() { grep -oE 'Build *= *[0-9]+' FyneApp.toml | grep -oE '[0-9]+'; }
write_build() { sed -i -E "s/^(\s*Build *= *).*/\1${1}/" FyneApp.toml; }
read_version() { grep -oE 'Version *= *"[^"]+"' FyneApp.toml | grep -oE '[0-9]+(\.[0-9]+)*'; }

declare -A _bak=()
swap() { _bak["$1"]="$(mktemp)"; cp "$1" "${_bak[$1]}"; }
write_swap() { swap "$1"; cat > "$1"; }
restore_swaps() { for s in "${!_bak[@]}"; do mv -f "${_bak[$s]}" "$s"; done; }

have() { grep -qF "$2" "$1" || { echo "[!] expected '$2' in $1" >&2; exit 1; }; }
gone() { ! grep -qF "$2" "$1" || { echo "[!] '$2' still present in $1" >&2; exit 1; }; }

sub() { sed -i "$2" "$1"; have "$1" "$3"; }
del() { sed -i "$2" "$1"; gone "$1" "$3"; }

replace_block() { OLD="$2" NEW="$3" perl -0777 -i -pe '
	my $i = index($_, $ENV{OLD}); die "block not found in '"$1"'\n" if $i < 0;
	substr($_, $i, length($ENV{OLD})) = $ENV{NEW};' "$1"; }

FONT_SUBSET_RANGES="U+0000-00FF,U+0100-017F,U+0400-04FF,U+2010-2027,U+2030-205E,U+20A0-20BF,U+2116,U+2122,U+2026"

round_dialogs() {
	echo "=== Round dialog/popup corners (match listitem selection radius) ==="
	local base=vendor/fyne.io/fyne/v2/dialog/base.go
	swap "$base"
	sub "$base" \
		's|rect := canvas.NewRectangle(theme.Color(theme.ColorNameOverlayBackground))|&\n\trect.CornerRadius = theme.Size(theme.SizeNameInputRadius)|' \
		"rect.CornerRadius"

	local popup=vendor/fyne.io/fyne/v2/widget/popup.go
	swap "$popup"
	sub "$popup" \
		's|background := canvas.NewRectangle(th.Color(theme.ColorNameOverlayBackground, v))|&\n\tbackground.CornerRadius = th.Size(theme.SizeNameInputRadius)|' \
		"background.CornerRadius"
}

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
				--unicodes="$FONT_SUBSET_RANGES" --no-hinting --desubroutinize \
				--drop-tables+=TTFA
		done
	else
		echo "[!] pyftsubset not found — skipping subset (install: pip install fonttools)"
	fi

	echo "=== Disable system-font scan ==="
	local fontprod=vendor/fyne.io/fyne/v2/internal/painter/font_prod.go
	write_swap "$fontprod" <<'GOEOF'
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

hide_webview_window() {
	local wv=vendor/github.com/webview/webview_go/libs/webview/include/webview.h
	echo "=== Keep owned webview window hidden during New() (kills Windows white flash) ==="
	swap "$wv"
	replace_block "$wv" \
		'    if (m_owns_window) {
      ShowWindow(m_window, SW_SHOW);
      UpdateWindow(m_window);
      SetFocus(m_window);
    }' \
		'    if (m_owns_window) {
      SetFocus(m_window);
    }'
	gone "$wv" "ShowWindow(m_window, SW_SHOW)"

	echo "=== Close() the WebView2 controller on teardown (stop msedgewebview2 lingering) ==="
	replace_block "$wv" \
		'    if (m_controller) {
      m_controller->Release();
      m_controller = nullptr;
    }' \
		'    if (m_controller) {
      m_controller->Close();
      m_controller->Release();
      m_controller = nullptr;
    }'
	have "$wv" "m_controller->Close()"
}

link_webkit() {
	local f=vendor/github.com/webview/webview_go/webview.go
	echo "=== Link launcher against webkit2gtk-4.1 (libsoup3) instead of 4.0 ==="
	swap "$f"
	sed -i 's/webkit2gtk-4\.0/webkit2gtk-4.1/' "$f"
	gone "$f" "webkit2gtk-4.0"
}

slim_charset() {
	local cf=vendor/golang.org/x/net/html/charset/charset.go
	echo "=== Slim SVG charset reader (UTF-8 only; drops x/text CJK tables) ==="
	write_swap "$cf" <<'GOEOF'
package charset

import "io"

func NewReaderLabel(label string, input io.Reader) (io.Reader, error) {
	return input, nil
}
GOEOF
}

slim_markdown() {
	local mf=vendor/fyne.io/fyne/v2/widget/markdown.go
	echo "=== Stub markdown parser (drops goldmark ~2.9MB + html5entities ~525KB) ==="
	write_swap "$mf" <<'GOEOF'
package widget

func NewRichTextFromMarkdown(content string) *RichText {
	return NewRichText(parseMarkdown(content)...)
}

func (t *RichText) ParseMarkdown(content string) {
	t.Segments = parseMarkdown(content)
	t.Refresh()
}

func (t *RichText) AppendMarkdown(content string) {
	t.Segments = append(t.Segments, parseMarkdown(content)...)
	t.Refresh()
}

func parseMarkdown(content string) []RichTextSegment {
	if content == "" {
		return nil
	}
	return []RichTextSegment{&TextSegment{Style: RichTextStyleParagraph, Text: content}}
}
GOEOF
}
