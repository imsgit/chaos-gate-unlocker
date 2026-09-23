package progress

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Widget struct {
	widget.BaseWidget

	bg   *canvas.Rectangle
	edge *canvas.Rectangle

	width  float32
	active bool
}

func New() *Widget {
	p := &Widget{
		bg:   canvas.NewRectangle(color.Transparent),
		edge: canvas.NewRectangle(color.White),
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *Widget) Grow(width float32) { p.set(width, true) }

func (p *Widget) Complete() { p.set(p.width, false) }

func (p *Widget) Reset() { p.set(0, false) }

func (p *Widget) set(width float32, active bool) {
	if p.width == width && p.active == active {
		return
	}
	p.width, p.active = width, active
	p.Refresh()
}

func (p *Widget) CreateRenderer() fyne.WidgetRenderer {
	p.bg.FillColor = p.Theme().Color(theme.ColorNameDisabledButton, 0)
	return &progressRenderer{progress: p}
}

type progressRenderer struct {
	progress *Widget
}

func (r *progressRenderer) Layout(size fyne.Size) {
	p := r.progress
	p.bg.Resize(size)

	if p.active {
		p.edge.Resize(fyne.NewSize(1, size.Height))
		p.edge.Move(fyne.NewPos(p.width-1, 0))
		p.edge.Show()
	} else {
		p.edge.Hide()
	}
}

func (r *progressRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 4)
}

func (r *progressRenderer) Refresh() {
	r.Layout(r.progress.Size())
	r.progress.edge.Refresh()
}

func (r *progressRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.progress.bg, r.progress.edge}
}

func (r *progressRenderer) Destroy() {}
