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
	t := &tip{Text: text, richtext: widget.NewRichTextWithText(text)}
	t.richtext.Wrapping = fyne.TextWrapWord
	t.richtext.Segments[0].(*widget.TextSegment).Style = tipTextStyle
	t.ExtendBaseWidget(t)
	return t
}

func (t *tip) MinSize() fyne.Size {
	return fyne.Size{}
}

func (t *tip) pad() float32 {
	return t.Theme().Size(theme.SizeNameInnerPadding) * 1.25
}

func (t *tip) textMinSize() fyne.Size {
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

func (t *tip) CreateRenderer() fyne.WidgetRenderer {
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

	r.tip.richtext.Refresh()
}

func (r *tipRenderer) Objects() []fyne.CanvasObject { return r.objects }

func (r *tipRenderer) Destroy() {}
