package statuslabel

import (
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"fyne.io/fyne/v2/widget"
)

type Widget struct {
	widget.Label
	tooltip.WidgetExtend
}

func New() *Widget {
	l := &Widget{}
	l.ExtendBaseWidget(l)
	l.ExtendToolTipWidget(l)
	return l
}

func (l *Widget) Set(text, toolTip string) {
	l.SetToolTip(toolTip)
	l.SetText(text)
}
