package icon

import (
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type IconWidget struct {
	widget.Icon
	tooltip.ToolTipWidgetExtend
}

func NewIconWidget() *IconWidget {
	i := &IconWidget{}

	i.ExtendBaseWidget(i)
	return i
}

func (i *IconWidget) ExtendBaseWidget(wid fyne.Widget) {
	i.ExtendToolTipWidget(wid)
	i.Icon.ExtendBaseWidget(wid)
}

func (i *IconWidget) MinSize() fyne.Size {
	i.ExtendBaseWidget(i)
	return fyne.NewSize(46, 46)
}

func (i *IconWidget) MouseIn(e *desktop.MouseEvent) {
	if tooltip.OverlayShown(i) {
		return
	}
	i.ToolTipWidgetExtend.MouseIn(e)
}
