package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	ttwidget "github.com/dweymouth/fyne-tooltip/widget"
)

type Icon struct {
	widget.Icon
	ttwidget.ToolTipWidgetExtend
}

func NewIcon() *Icon {
	i := &Icon{}

	i.ExtendBaseWidget(i)
	return i
}

func (i *Icon) ExtendBaseWidget(wid fyne.Widget) {
	i.ExtendToolTipWidget(wid)
	i.Icon.ExtendBaseWidget(wid)
}

func (i *Icon) MinSize() fyne.Size {
	i.ExtendBaseWidget(i)
	return fyne.NewSize(46, 46)
}

func (i *Icon) CreateRenderer() fyne.WidgetRenderer {
	return i.Icon.CreateRenderer()
}
