package snapimage

import (
	"image"

	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/pixelsnap"
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Widget struct {
	widget.BaseWidget
	tooltip.WidgetExtend

	img *canvas.Image
	min fyne.Size
}

func New(min fyne.Size) *Widget {
	w := &Widget{
		img: &canvas.Image{FillMode: canvas.ImageFillStretch, ScaleMode: canvas.ImageScaleFastest},
		min: min,
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *Widget) ExtendBaseWidget(wid fyne.Widget) {
	w.ExtendToolTipWidget(wid)
	w.BaseWidget.ExtendBaseWidget(wid)
}

func (w *Widget) SetResource(res fyne.Resource) {
	w.SetImage(ui.DecodeIcon(res))
}

func (w *Widget) SetImage(img image.Image) {
	w.img.Resource = nil
	w.img.Image = img
	w.Refresh()
}

func (w *Widget) SetTranslucency(t float64) {
	if w.img.Translucency == t {
		return
	}
	w.img.Translucency = t
	w.img.Refresh()
}

func (w *Widget) MinSize() fyne.Size {
	w.ExtendBaseWidget(w)
	return w.min
}

func (w *Widget) MouseIn(e *desktop.MouseEvent) {
	if tooltip.OverlayShown(w) {
		return
	}
	w.WidgetExtend.MouseIn(e)
}

func (w *Widget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return &renderer{w: w}
}

type renderer struct{ w *Widget }

func (r *renderer) Layout(size fyne.Size)        { pixelsnap.Image(r.w.img, size, r.w) }
func (r *renderer) MinSize() fyne.Size           { return r.w.min }
func (r *renderer) Refresh()                     { r.w.img.Refresh(); r.Layout(r.w.Size()) }
func (r *renderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.w.img} }
func (r *renderer) Destroy()                     {}
