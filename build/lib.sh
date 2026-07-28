_toml() { grep -oE "$1" FyneApp.toml | grep -oE "$2"; }
read_build() { _toml 'Build *= *[0-9]+' '[0-9]+'; }
read_version() { _toml 'Version *= *"[^"]+"' '[0-9]+(\.[0-9]+)*'; }
write_build() { sed -i -E "s/^(\s*Build *= *).*/\1${1}/" FyneApp.toml; }

declare -A _bak=()
declare -a _added=()
swap() { [ -n "${_bak[$1]:-}" ] || { _bak["$1"]="$(mktemp)"; cp "$1" "${_bak[$1]}"; }; }
write_swap() { swap "$1"; cat > "$1"; }
add_swap() { _added+=("$1"); cat > "$1"; }
restore_swaps() {
	for s in "${!_bak[@]}"; do mv -f "${_bak[$s]}" "$s"; done
	for f in "${_added[@]}"; do rm -f "$f"; done
}

ensure_vendor() {
	created_vendor=
	[ -d vendor ] || { echo "=== go mod vendor ($1) ==="; go mod vendor; created_vendor=1; }
	trap 'restore_swaps; [ -n "$created_vendor" ] && rm -rf vendor' EXIT
}

have() { grep -qF "$2" "$1" || { echo "[!] expected '$2' in $1" >&2; exit 1; }; }
gone() { ! grep -qF "$2" "$1" || { echo "[!] '$2' still present in $1" >&2; exit 1; }; }

sub() { sed -i "$2" "$1"; have "$1" "$3"; }
del() { sed -i "$2" "$1"; gone "$1" "$3"; }

replace_block() { OLD="$2" NEW="$3" perl -0777 -i -pe '
	my $i = index($_, $ENV{OLD}); die "block not found in '"$1"'\n" if $i < 0;
	substr($_, $i, length($ENV{OLD})) = $ENV{NEW};' "$1"; }

vsub() { swap "$1"; sub "$1" "$2" "$3"; }
vblock() { swap "$1"; replace_block "$1" "$2" "$3"; }

FONT_SUBSET_RANGES="U+0000-00FF,U+0100-017F,U+0180-024F,U+0370-03FF,U+0400-04FF,U+1E00-1EFF,U+2010-2027,U+2030-205E,U+20A0-20BF,U+2116,U+2122,U+2026"

slim_common() {
	stub_fonts
	slim_charset
	slim_markdown
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
//go:build !ci && !test

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
	vblock "$wv" \
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
	vblock "$wv" \
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
	vsub "$f" 's/webkit2gtk-4\.0/webkit2gtk-4.1/' 'webkit2gtk-4.1'
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

enable_touch_scroll() {
	local ww=vendor/fyne.io/fyne/v2/internal/driver/glfw/window_wasm.go
	echo "=== Flush pending move before click ==="
	vblock "$ww" \
		'	runOnMain(func() {
		button, modifiers := convertMouseButton(btn, mods)' \
		'	runOnMain(func() {
		if !w.mousePosUpdateProcessed {
			w.processMouseMoved(w.newMousePosX, w.newMousePosY)
			w.mousePosUpdateProcessed = true
		}
		button, modifiers := convertMouseButton(btn, mods)'
	have "$ww" 'if !w.mousePosUpdateProcessed {'

	local wc=vendor/fyne.io/fyne/v2/internal/driver/glfw/window.go
	echo "=== Widen drag slop 2->12 for touch (finger jitter must not turn a tap into a scroll-drag) ==="
	vsub "$wc" 's/dragMoveThreshold = 2 /dragMoveThreshold = 12 /' 'dragMoveThreshold = 12'

	local bw=vendor/github.com/fyne-io/glfw-js/browser_wasm.go
	echo "=== Enable touch->mouse emulation in glfw-js (touch scroll on mobile/Deck browsers) ==="
	have "$bw" 'addDocumentEventListener.Invoke("beforeUnload"'
	vblock "$bw" \
		'	addDocumentEventListener.Invoke("beforeUnload",' \
		'	touchPos := func(te js.Value) (float64, float64, bool) {
		touches := te.Get("touches")
		if touches.Length() == 0 {
			return 0, 0, false
		}
		t := touches.Index(0)
		return t.Get("clientX").Float() * w.devicePixelRatio, t.Get("clientY").Float() * w.devicePixelRatio, true
	}
	touchStart := newJsFuncFrom(func(this js.Value, args []js.Value) any {
		te := args[0]
		w.goFullscreenIfRequested()
		if x, y, ok := touchPos(te); ok {
			w.cursorPos[0], w.cursorPos[1] = x, y
			if w.cursorPosCallback != nil {
				w.cursorPosCallback(w, x, y)
			}
		}
		w.mouseButton[0] = Press
		if w.mouseButtonCallback != nil {
			go w.mouseButtonCallback(w, MouseButton1, Press, 0)
		}
		te.Call("preventDefault")
		return nil
	})
	touchMove := newJsFuncFrom(func(this js.Value, args []js.Value) any {
		te := args[0]
		if x, y, ok := touchPos(te); ok {
			mvX, mvY := x-w.cursorPos[0], y-w.cursorPos[1]
			w.cursorPos[0], w.cursorPos[1] = x, y
			if w.cursorPosCallback != nil {
				w.cursorPosCallback(w, x, y)
			}
			if w.mouseMovementCallback != nil {
				go w.mouseMovementCallback(w, x, y, mvX, mvY)
			}
		}
		te.Call("preventDefault")
		return nil
	})
	touchEnd := newJsFuncFrom(func(this js.Value, args []js.Value) any {
		te := args[0]
		w.mouseButton[0] = Release
		if w.mouseButtonCallback != nil {
			go w.mouseButtonCallback(w, MouseButton1, Release, 0)
		}
		te.Call("preventDefault")
		return nil
	})
	touchOpts := map[string]any{"passive": false}
	addDocumentEventListener.Invoke("touchstart", touchStart, touchOpts)
	addDocumentEventListener.Invoke("touchmove", touchMove, touchOpts)
	addDocumentEventListener.Invoke("touchend", touchEnd, touchOpts)
	addDocumentEventListener.Invoke("touchcancel", touchEnd, touchOpts)

	addDocumentEventListener.Invoke("beforeUnload",'
	have "$bw" 'Invoke("touchstart"'
}

drag_scroll_widget() {
	local f=vendor/fyne.io/fyne/v2/internal/widget/scroller_web.go
	echo "=== Make Scroll draggable on wasm (finger drag scrolls list + dropdown content) ==="
	gone vendor/fyne.io/fyne/v2/internal/widget/scroller.go 'func (s *Scroll) Dragged'
	add_swap "$f" <<'GOEOF'
//go:build js

package widget

import "fyne.io/fyne/v2"

func (s *Scroll) DragEnd() {
}

func (s *Scroll) Dragged(e *fyne.DragEvent) {
	if s.updateOffset(e.Dragged.DX, e.Dragged.DY) {
		s.refreshWithoutOffsetUpdate()
	}
}
GOEOF
}
