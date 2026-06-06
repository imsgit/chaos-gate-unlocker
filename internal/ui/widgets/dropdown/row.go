package dropdown

import (
	"image/color"
	"math"

	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const iconTextGap = 8

type selectRow struct {
	widget.BaseWidget
	tooltip.WidgetExtend

	text     string
	icon     fyne.Resource
	disabled bool
	scale    float32
	onTapped func()

	hovered     bool
	highlighted bool
}

func (r *selectRow) setHighlighted(h bool) {
	if r.highlighted == h {
		return
	}
	r.highlighted = h
	r.Refresh()
}

func newSelectRow(text string, icon fyne.Resource, disabled bool, scale float32, onTapped func()) *selectRow {
	r := &selectRow{text: text, icon: icon, disabled: disabled, scale: scale, onTapped: onTapped}
	r.ExtendBaseWidget(r)
	return r
}

func (r *selectRow) ExtendBaseWidget(wid fyne.Widget) {
	r.ExtendToolTipWidget(wid)
	r.BaseWidget.ExtendBaseWidget(wid)
}

func (r *selectRow) Tapped(*fyne.PointEvent) {
	if r.disabled || r.onTapped == nil {
		return
	}
	r.onTapped()
}

func (r *selectRow) MouseIn(e *desktop.MouseEvent) {
	r.WidgetExtend.MouseIn(e)
	if !r.disabled && !r.hovered {
		r.hovered = true
		r.Refresh()
	}
}

func (r *selectRow) MouseOut() {
	r.WidgetExtend.MouseOut()
	if r.hovered {
		r.hovered = false
		r.Refresh()
	}
}

func (r *selectRow) CreateRenderer() fyne.WidgetRenderer {
	th := r.Theme()

	bg := canvas.NewRectangle(color.Transparent)
	bg.CornerRadius = th.Size(theme.SizeNameSelectionRadius)

	text := canvas.NewText("", color.Transparent)
	text.TextSize = th.Size(theme.SizeNameText)

	rr := &selectRowRenderer{row: r, bg: bg, text: text}
	if img := ui.DecodeMasked(r.icon); img != nil {
		rr.img = ui.NewIconImage(fyne.Size{})
		rr.img.Image = img
	}
	rr.Refresh()
	return rr
}

type selectRowRenderer struct {
	row  *selectRow
	bg   *canvas.Rectangle
	img  *canvas.Image
	text *canvas.Text
}

func (r *selectRowRenderer) pad() float32 {
	return r.row.Theme().Size(theme.SizeNamePadding)
}

func (r *selectRowRenderer) iconSize() float32 {
	return r.text.MinSize().Height * 1.6
}

func (r *selectRowRenderer) Layout(size fyne.Size) {
	pad := r.pad()
	r.bg.Resize(size)

	textX := pad * 2
	if r.img != nil {
		s := r.iconSize()
		r.img.Resize(fyne.NewSize(s, s))
		r.img.Move(fyne.NewPos(pad*2, (size.Height-s)/2))
		textX = pad*2 + s + iconTextGap
	}
	r.text.Move(fyne.NewPos(textX, pad))
	r.text.Resize(fyne.NewSize(size.Width-textX-pad*2, size.Height-pad*2))
}

func (r *selectRowRenderer) MinSize() fyne.Size {
	pad := r.pad()
	ts := r.text.MinSize()
	w := ts.Width + pad*4
	h := ts.Height
	if r.img != nil {
		s := r.iconSize()
		w += s + iconTextGap
		if s > h {
			h = s
		}
	}

	height := h + pad*2
	if r.row.scale > 0 {
		height = float32(math.Ceil(float64(height*r.row.scale))) / r.row.scale
	}
	return fyne.NewSize(w, height)
}

func (r *selectRowRenderer) Refresh() {
	th := r.row.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.text.Text = r.row.text
	if r.row.disabled {
		r.text.Color = th.Color(theme.ColorNameDisabled, v)
	} else {
		r.text.Color = th.Color(theme.ColorNameForeground, v)
	}

	if r.img != nil {
		r.img.Translucency = 0
		if r.row.disabled {
			_, _, _, a := th.Color(theme.ColorNameDisabled, v).RGBA()
			r.img.Translucency = 1 - float64(a)/0xFFFF
		}
		r.img.Refresh()
	}

	if r.row.hovered || r.row.highlighted {
		r.bg.FillColor = th.Color(theme.ColorNameHover, v)
	} else {
		r.bg.FillColor = color.Transparent
	}

	r.text.Refresh()
	r.bg.Refresh()
}

func (r *selectRowRenderer) Objects() []fyne.CanvasObject {
	if r.img != nil {
		return []fyne.CanvasObject{r.bg, r.img, r.text}
	}
	return []fyne.CanvasObject{r.bg, r.text}
}

func (r *selectRowRenderer) Destroy() {}
