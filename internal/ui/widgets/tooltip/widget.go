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
	tipObj           fyne.CanvasObject
	absoluteMousePos fyne.Position
	pending          *time.Timer
	pendingActive    bool
	popUp            *time.Timer
	popUpActive      bool
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

func (t *WidgetExtend) PopUpToolTip() {
	if t.toolTip == "" || t.Obj == nil {
		return
	}
	t.cancel()
	t.popUpActive = true
	t.popUp = time.AfterFunc(popUpDelay, func() {
		fyne.Do(func() {
			if !t.popUpActive || t.toolTip == "" {
				return
			}
			driver := fyne.CurrentApp().Driver()
			pos := driver.AbsolutePositionForObject(t.Obj)
			pos.X += t.Obj.Size().Width / 2
			t.show(driver.CanvasForObject(t.Obj), pos)
			t.popUp = time.AfterFunc(popUpDuration, func() { fyne.Do(t.DismissToolTip) })
		})
	})
}

func (t *WidgetExtend) DismissToolTip() { t.cancel() }

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
			t.show(fyne.CurrentApp().Driver().CanvasForObject(t.Obj), t.absoluteMousePos)
		})
	})
}

func (t *WidgetExtend) show(canvas fyne.Canvas, pos fyne.Position) {
	t.tipLayer = showAtMousePosition(canvas, pos, t.toolTip)
	if t.tipLayer != nil && len(t.tipLayer.Objects) > 0 {
		t.tipObj = t.tipLayer.Objects[0]
	}
}

func (t *WidgetExtend) cancel() {
	t.pendingActive = false
	if t.pending != nil {
		t.pending.Stop()
	}
	t.popUpActive = false
	if t.popUp != nil {
		t.popUp.Stop()
		t.popUp = nil
	}
	if t.tipLayer != nil {
		if len(t.tipLayer.Objects) > 0 && t.tipLayer.Objects[0] == t.tipObj {
			hide(t.tipLayer)
		}
		t.tipLayer, t.tipObj = nil, nil
	}
}
