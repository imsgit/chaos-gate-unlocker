package progress

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ProgressWidget struct {
	widget.BaseWidget

	bar  *canvas.Rectangle
	edge *canvas.Rectangle

	width  float32
	active bool
}

func NewProgressWidget() *ProgressWidget {
	p := &ProgressWidget{
		bar:  canvas.NewRectangle(color.Transparent),
		edge: canvas.NewRectangle(color.White),
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *ProgressWidget) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return fyne.NewSize(0, 4)
}

func (p *ProgressWidget) Grow(width float32) {
	p.width = width
	p.active = true
	p.bar.FillColor = p.Theme().Color(theme.ColorNameShadow, 0)
	p.Refresh()
}

func (p *ProgressWidget) Complete() {
	p.active = false
	p.Refresh()
}

func (p *ProgressWidget) Reset() {
	p.width = 0
	p.active = false
	p.bar.FillColor = color.Transparent
	p.Refresh()
}

func (p *ProgressWidget) CreateRenderer() fyne.WidgetRenderer {
	return &progressRenderer{progress: p}
}

type progressRenderer struct {
	progress *ProgressWidget
}

func (r *progressRenderer) Layout(size fyne.Size) {
	p := r.progress
	p.bar.Resize(fyne.NewSize(p.width, size.Height))
	p.bar.Move(fyne.NewPos(0, 0))

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
	r.progress.bar.Refresh()
	r.progress.edge.Refresh()
}

func (r *progressRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.progress.bar, r.progress.edge}
}

func (r *progressRenderer) Destroy() {}
