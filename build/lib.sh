read_build() { grep -oE 'Build *= *[0-9]+' FyneApp.toml | grep -oE '[0-9]+'; }
write_build() { sed -i -E "s/^(\s*Build *= *).*/\1${1}/" FyneApp.toml; }

declare -A _bak=()
swap() { _bak["$1"]="$(mktemp)"; cp "$1" "${_bak[$1]}"; }
write_swap() { swap "$1"; cat > "$1"; }
declare -a _added=()
add_file() { _added+=("$1"); }
restore_swaps() {
	for s in "${!_bak[@]}"; do mv -f "${_bak[$s]}" "$s"; done
	local f; for f in "${_added[@]:-}"; do [ -n "$f" ] && rm -f "$f"; done
}

have() { grep -qF "$2" "$1" || { echo "[!] expected '$2' in $1" >&2; exit 1; }; }
gone() { ! grep -qF "$2" "$1" || { echo "[!] '$2' still present in $1" >&2; exit 1; }; }

sub() { sed -i "$2" "$1"; have "$1" "$3"; }
del() { sed -i "$2" "$1"; gone "$1" "$3"; }

replace_block() { OLD="$2" NEW="$3" perl -0777 -i -pe '
	my $i = index($_, $ENV{OLD}); die "block not found in '"$1"'\n" if $i < 0;
	substr($_, $i, length($ENV{OLD})) = $ENV{NEW};' "$1"; }

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

enable_drag_scroll() {
	local wdir=vendor/fyne.io/fyne/v2/internal/widget
	local f="$wdir/scroller_desktop_drag.go"
	echo "=== Enable drag-to-scroll on desktop (finger/touchscreen drag, e.g. Steam Deck) ==="
	[ -f "$wdir/scroller_mobile.go" ] || { echo "[!] scroller_mobile.go missing — Fyne internals moved" >&2; exit 1; }
	[ -f "$f" ] && { echo "[!] $f already exists" >&2; exit 1; }
	add_file "$f"
	cat > "$f" <<'GOEOF'
//go:build !ci && !no_glfw && !android && !ios && !mobile

package widget

import "fyne.io/fyne/v2"

func (s *Scroll) Dragged(e *fyne.DragEvent) {
	if s.updateOffset(e.Dragged.DX, e.Dragged.DY) {
		s.refreshWithoutOffsetUpdate()
	}
}

func (s *Scroll) DragEnd() {}
GOEOF
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

slim_test_asserts() {
	echo "=== Stub fyne test-assert helpers (drops testify+go-spew+difflib+yaml; ~349KB stripped exe) ==="
	local f
	for f in \
		vendor/fyne.io/fyne/v2/test/test_helper.go \
		vendor/fyne.io/fyne/v2/test/notification_helper.go \
		vendor/fyne.io/fyne/v2/test/theme_helper.go \
		vendor/fyne.io/fyne/v2/internal/test/util_helper.go; do
		have "$f" "stretchr/testify"
		write_swap "$f" <<<'package test'
		gone "$f" "stretchr/testify"
	done
}

slim_filedialog() {
	local filego=vendor/fyne.io/fyne/v2/dialog/file.go
	echo "=== Slim file-open dialog (drop favorites, keep gear + add Copy-path button) ==="
	swap "$filego"

	replace_block "$filego" \
		$'\toptionsbuttons := container.NewHBox(' \
		$'\tcopyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {\n\t\tif f.dir != nil {\n\t\t\tfyne.CurrentApp().Clipboard().SetContent(f.dir.Path())\n\t\t}\n\t})\n\n\toptionsbuttons := container.NewHBox('
	have "$filego" "copyButton := widget.NewButtonWithIcon"

	replace_block "$filego" \
		$'\theader := container.NewBorder(\n\t\tnil, nil, nil, optionsbuttons,\n\t\tf.title,\n\t)' \
		$'\t_ = optionsbuttons\n\theader := container.NewBorder(\n\t\tnil, nil, nil, container.NewHBox(copyButton, optionsButton),\n\t\tf.title,\n\t)'
	replace_block "$filego" \
		$'\tbody := container.NewHSplit(\n\t\tf.favoritesList,\n\t\tcontainer.NewBorder(\n\t\t\tf.breadcrumbScroll, nil, nil, nil,\n\t\t\tf.filesScroll,\n\t\t),\n\t)\n\tbody.SetOffset(0)' \
		$'\tbody := container.NewBorder(\n\t\tf.breadcrumbScroll, nil, nil, nil,\n\t\tf.filesScroll,\n\t)'
}
