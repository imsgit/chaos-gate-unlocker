package tooltip

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Button struct {
	widget.Button
	WidgetExtend
}

func NewButton(text string, onTapped func()) *Button {
	b := &Button{Button: widget.Button{Text: text, OnTapped: onTapped}}
	b.ExtendBaseWidget(b)
	return b
}

func (b *Button) ExtendBaseWidget(wid fyne.Widget) {
	b.ExtendToolTipWidget(wid)
	b.Button.ExtendBaseWidget(wid)
}

func (b *Button) Tapped(e *fyne.PointEvent) {
	b.WidgetExtend.cancel()
	b.Button.Tapped(e)
}

func (b *Button) MouseIn(e *desktop.MouseEvent) {
	b.WidgetExtend.MouseIn(e)
	b.Button.MouseIn(e)
}

func (b *Button) MouseMoved(e *desktop.MouseEvent) {
	b.WidgetExtend.MouseMoved(e)
	b.Button.MouseMoved(e)
}

func (b *Button) MouseOut() {
	b.WidgetExtend.MouseOut()
	b.Button.MouseOut()
}
