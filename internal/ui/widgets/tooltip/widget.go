package tooltip

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type WidgetExtend struct {
	Obj fyne.CanvasObject

	toolTip string

	tipLayer         *fyne.Container
	absoluteMousePos fyne.Position
	pending          *time.Timer
	pendingActive    bool
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
	t.pendingActive = true
	if t.pending != nil {
		t.pending.Reset(nextDelay())
		return
	}
	t.pending = time.AfterFunc(nextDelay(), func() {
		fyne.Do(func() {
			if !t.pendingActive {
				return
			}
			t.cancel()
			canvas := fyne.CurrentApp().Driver().CanvasForObject(t.Obj)
			t.tipLayer = showAtMousePosition(canvas, t.absoluteMousePos, t.toolTip)
		})
	})
}

func (t *WidgetExtend) cancel() {
	t.pendingActive = false
	if t.pending != nil {
		t.pending.Stop()
	}
	if t.tipLayer != nil {
		hide(t.tipLayer)
		t.tipLayer = nil
	}
}
