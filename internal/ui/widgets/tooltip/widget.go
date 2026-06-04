package tooltip

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type ToolTipWidgetExtend struct {
	Obj fyne.CanvasObject

	toolTip string

	handle           *handle
	absoluteMousePos fyne.Position
	pendingCancel    context.CancelFunc
}

func (t *ToolTipWidgetExtend) SetToolTip(toolTip string) { t.toolTip = toolTip }

func (t *ToolTipWidgetExtend) ToolTip() string { return t.toolTip }

func (t *ToolTipWidgetExtend) ExtendToolTipWidget(wid fyne.Widget) { t.Obj = wid }

func (t *ToolTipWidgetExtend) MouseIn(e *desktop.MouseEvent) {
	if t.toolTip == "" {
		return
	}
	t.absoluteMousePos = e.AbsolutePosition
	t.setPending()
}

func (t *ToolTipWidgetExtend) MouseMoved(e *desktop.MouseEvent) {
	t.absoluteMousePos = e.AbsolutePosition
}

func (t *ToolTipWidgetExtend) MouseOut() {
	t.cancel()
}

func (t *ToolTipWidgetExtend) setPending() {
	ctx, cancel := context.WithCancel(context.Background())
	t.pendingCancel = cancel

	delay := nextDelay()
	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(delay):
			fyne.Do(func() {
				if ctx.Err() != nil {
					return
				}
				t.cancel()
				canvas := fyne.CurrentApp().Driver().CanvasForObject(t.Obj)
				t.handle = showAtMousePosition(canvas, t.absoluteMousePos, t.toolTip)
			})
		}
	}()
}

func (t *ToolTipWidgetExtend) cancel() {
	if t.pendingCancel != nil {
		t.pendingCancel()
		t.pendingCancel = nil
	}
	if t.handle != nil {
		hide(t.handle)
		t.handle = nil
	}
}
