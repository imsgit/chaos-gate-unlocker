package dragscroll

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Wrapper struct {
	widget.BaseWidget
	content  fyne.CanvasObject
	scrollBy func(dy float32)
}

func newWrapper(content fyne.CanvasObject, scrollBy func(dy float32)) *Wrapper {
	w := &Wrapper{content: content, scrollBy: scrollBy}
	w.ExtendBaseWidget(w)
	return w
}

func Scroll(s *container.Scroll) *Wrapper {
	return newWrapper(s, func(dy float32) {
		s.ScrollToOffset(fyne.NewPos(s.Offset.X, s.Offset.Y-dy))
	})
}

func List(l *widget.List) *Wrapper {
	return newWrapper(l, func(dy float32) {
		l.ScrollToOffset(l.GetScrollOffset() - dy)
	})
}

func (w *Wrapper) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.content)
}

func (w *Wrapper) Dragged(e *fyne.DragEvent) {
	w.scrollBy(e.Dragged.DY)
}

func (w *Wrapper) DragEnd() {}
