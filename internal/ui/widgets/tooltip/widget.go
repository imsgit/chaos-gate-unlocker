package tooltip

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type WidgetExtend struct {
	Obj fyne.CanvasObject

	toolTip string

	tipLayer         *fyne.Container
	absoluteMousePos fyne.Position
	pendingCancel    context.CancelFunc
}

func (t *WidgetExtend) SetToolTip(toolTip string) { t.toolTip = toolTip }

func (t *WidgetExtend) ExtendToolTipWidget(wid fyne.Widget) { t.Obj = wid }

func (t *WidgetExtend) MouseIn(e *desktop.MouseEvent) {
	if t.toolTip == "" {
		return
	}
	t.absoluteMousePos = e.AbsolutePosition
	t.setPending()
}

func (t *WidgetExtend) MouseInUnlessOverlay(e *desktop.MouseEvent) {
	if !OverlayShown(t.Obj) {
		t.MouseIn(e)
	}
}

func (t *WidgetExtend) MouseMoved(e *desktop.MouseEvent) {
	t.absoluteMousePos = e.AbsolutePosition
}

func (t *WidgetExtend) MouseOut() {
	t.cancel()
}

func (t *WidgetExtend) setPending() {
	if t.pendingCancel != nil {
		t.pendingCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.pendingCancel = cancel

	time.AfterFunc(nextDelay(), func() {
		fyne.Do(func() {
			if ctx.Err() != nil {
				return
			}
			t.cancel()
			canvas := fyne.CurrentApp().Driver().CanvasForObject(t.Obj)
			t.tipLayer = showAtMousePosition(canvas, t.absoluteMousePos, t.toolTip)
		})
	})
}

func (t *WidgetExtend) cancel() {
	if t.pendingCancel != nil {
		t.pendingCancel()
		t.pendingCancel = nil
	}
	if t.tipLayer != nil {
		hide(t.tipLayer)
		t.tipLayer = nil
	}
}
