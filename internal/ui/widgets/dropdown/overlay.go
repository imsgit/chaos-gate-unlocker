package dropdown

import (
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type popOverlay struct {
	widget.BaseWidget

	canvas  fyne.Canvas
	content fyne.CanvasObject

	innerPos  fyne.Position
	innerSize fyne.Size
	shown     bool

	onDismiss func()
}

func newPopOverlay(content fyne.CanvasObject, canvas fyne.Canvas, onDismiss func()) *popOverlay {
	p := &popOverlay{canvas: canvas, onDismiss: onDismiss}
	p.ExtendBaseWidget(p)

	if layer := tooltip.AddOverlayToolTipLayer(p, canvas); layer != nil {
		content = container.NewStack(content, layer)
	}
	p.content = content
	return p
}

func (p *popOverlay) showAt(pos fyne.Position, size fyne.Size) {
	p.innerPos = pos
	p.innerSize = size
	if !p.shown {
		p.canvas.Overlays().Add(p)
		p.shown = true
	}
	p.Move(fyne.NewPos(0, 0))
	p.Resize(p.canvas.Size())
	p.Refresh()
}

func (p *popOverlay) Hide() {
	if p.shown {
		p.canvas.Overlays().Remove(p)
		p.shown = false
		tooltip.RemoveOverlayToolTipLayer(p, p.canvas)
	}
	p.BaseWidget.Hide()
}

func (p *popOverlay) isInsideContent(pos fyne.Position) bool {
	return pos.X >= p.innerPos.X && pos.Y >= p.innerPos.Y &&
		pos.X <= p.innerPos.X+p.innerSize.Width &&
		pos.Y <= p.innerPos.Y+p.innerSize.Height
}

func (p *popOverlay) Tapped(e *fyne.PointEvent) {
	if !p.isInsideContent(e.Position) {
		p.onDismiss()
	}
}

func (p *popOverlay) TappedSecondary(e *fyne.PointEvent) {
	if !p.isInsideContent(e.Position) {
		p.onDismiss()
	}
}

func (p *popOverlay) CreateRenderer() fyne.WidgetRenderer {
	return &popOverlayRenderer{pop: p, objects: []fyne.CanvasObject{p.content}}
}

type popOverlayRenderer struct {
	pop     *popOverlay
	objects []fyne.CanvasObject
}

func (r *popOverlayRenderer) Layout(_ fyne.Size) {
	pos := r.pop.innerPos
	size := r.pop.innerSize

	canvasSize := r.pop.canvas.Size()
	if pos.X+size.Width > canvasSize.Width {
		pos.X = fyne.Max(0, canvasSize.Width-size.Width)
	}
	if pos.Y+size.Height > canvasSize.Height {
		pos.Y = fyne.Max(0, canvasSize.Height-size.Height)
	}

	r.pop.content.Resize(size)
	r.pop.content.Move(pos)
}

func (r *popOverlayRenderer) MinSize() fyne.Size { return r.pop.innerSize }

func (r *popOverlayRenderer) Refresh() {
	r.Layout(r.pop.Size())
	r.pop.content.Refresh()
}

func (r *popOverlayRenderer) Objects() []fyne.CanvasObject { return r.objects }

func (r *popOverlayRenderer) Destroy() {}
