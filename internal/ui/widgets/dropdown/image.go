package dropdown

import (
	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type iconImage struct {
	widget.BaseWidget
	tooltip.WidgetExtend

	img *canvas.Image
}

func newIconImage(min fyne.Size) *iconImage {
	w := &iconImage{img: ui.NewIcon(min)}
	w.ExtendBaseWidget(w)
	return w
}

func (w *iconImage) ExtendBaseWidget(wid fyne.Widget) {
	w.ExtendToolTipWidget(wid)
	w.BaseWidget.ExtendBaseWidget(wid)
}

func (w *iconImage) SetResource(res fyne.Resource) {
	w.img.Resource = nil
	w.img.Image = ui.DecodeIcon(res)
	w.Refresh()
}

func (w *iconImage) MouseIn(e *desktop.MouseEvent) {
	if tooltip.OverlayShown(w) {
		return
	}
	w.WidgetExtend.MouseIn(e)
}

func (w *iconImage) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return widget.NewSimpleRenderer(w.img)
}
