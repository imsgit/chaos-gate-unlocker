package tooltip

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var tipTextStyle = widget.RichTextStyle{SizeName: theme.SizeNameCaptionText}

type tip struct {
	widget.BaseWidget

	Text     string
	richtext *widget.RichText
}

func newTip(text string) *tip {
	t := &tip{Text: text}
	t.ExtendBaseWidget(t)
	return t
}

func (t *tip) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (t *tip) Resize(size fyne.Size) {
	t.updateRichText()
	t.richtext.Resize(size)
	t.BaseWidget.Resize(size)
}

func (t *tip) pad() float32 {
	return t.Theme().Size(theme.SizeNameInnerPadding) * 1.25
}

func (t *tip) textMinSize() fyne.Size {
	t.updateRichText()
	innerPad := t.Theme().Size(theme.SizeNameInnerPadding)
	contentH := t.richtext.MinSize().Height - 2*innerPad
	return fyne.NewSize(t.textWidth(), contentH+2*t.pad())
}

func (t *tip) textWidth() float32 {
	th := t.Theme()
	size := th.Size(tipTextStyle.SizeName)

	var widest float32
	for _, line := range strings.Split(t.Text, "\n") {
		if w := fyne.MeasureText(line, size, tipTextStyle.TextStyle).Width; w > widest {
			widest = w
		}
	}
	return widest + 2*t.pad()
}

func (t *tip) updateRichText() {
	if t.richtext == nil {
		t.richtext = widget.NewRichTextWithText(t.Text)
		t.richtext.Wrapping = fyne.TextWrapWord
	}
	seg := t.richtext.Segments[0].(*widget.TextSegment)
	seg.Text = t.Text
	seg.Style = tipTextStyle
}

func (t *tip) CreateRenderer() fyne.WidgetRenderer {
	t.updateRichText()
	bg := canvas.NewRectangle(color.Transparent)
	bg.CornerRadius = t.Theme().Size(theme.SizeNameSelectionRadius)
	return &tipRenderer{tip: t, bg: bg, objects: []fyne.CanvasObject{bg, t.richtext}}
}

type tipRenderer struct {
	tip     *tip
	bg      *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *tipRenderer) Layout(s fyne.Size) {
	r.bg.Resize(s)
	r.bg.Move(fyne.NewPos(0, 0))

	innerPad := r.tip.Theme().Size(theme.SizeNameInnerPadding)
	off := r.tip.pad() - innerPad
	r.tip.richtext.Resize(s)
	r.tip.richtext.Move(fyne.NewPos(off, off))
}

func (r *tipRenderer) MinSize() fyne.Size {
	return r.tip.textMinSize()
}

func (r *tipRenderer) Refresh() {
	th := r.tip.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.bg.FillColor = th.Color(theme.ColorNameOverlayBackground, v)
	r.bg.StrokeColor = th.Color(theme.ColorNameInputBorder, v)
	r.bg.StrokeWidth = th.Size(theme.SizeNameInputBorder)
	r.bg.Refresh()

	r.tip.updateRichText()
	r.tip.richtext.Refresh()
	canvas.Refresh(r.tip)
}

func (r *tipRenderer) Objects() []fyne.CanvasObject { return r.objects }

func (r *tipRenderer) Destroy() {}
