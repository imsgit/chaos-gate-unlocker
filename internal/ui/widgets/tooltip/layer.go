package tooltip

import (
	"errors"
	"time"

	"chaos-gate-unlocker/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

const (
	initialDelay        = 750 * time.Millisecond
	subsequentDelay     = 300 * time.Millisecond
	subsequentValidTime = 1500 * time.Millisecond
	maxWidth            = 600
	belowMouseDistance  = 16
	aboveMouseDistance  = 8
)

type layer struct {
	Container fyne.Container
	overlays  map[fyne.CanvasObject]*fyne.Container
}

var (
	layers    = make(map[fyne.Canvas]*layer)
	lastShown time.Time
)

func AddWindowToolTipLayer(content fyne.CanvasObject, canvas fyne.Canvas) fyne.CanvasObject {
	l := &layer{}
	layers[canvas] = l
	return container.NewStack(content, &l.Container)
}

func AddOverlayToolTipLayer(overlay fyne.CanvasObject, canvas fyne.Canvas) *fyne.Container {
	parent := layers[canvas]
	if parent == nil {
		fyne.LogError("", errors.New("no tooltip layer for parent canvas"))
		return nil
	}
	if parent.overlays == nil {
		parent.overlays = make(map[fyne.CanvasObject]*fyne.Container)
	}
	c := &fyne.Container{}
	parent.overlays[overlay] = c
	return c
}

func RemoveOverlayToolTipLayer(overlay fyne.CanvasObject, canvas fyne.Canvas) {
	if parent := layers[canvas]; parent != nil {
		delete(parent.overlays, overlay)
	}
}

func OverlayShown(obj fyne.CanvasObject) bool {
	c := fyne.CurrentApp().Driver().CanvasForObject(obj)
	if c == nil {
		return false
	}
	return c.Overlays().Top() != nil
}

func nextDelay() time.Duration {
	if time.Since(lastShown) < subsequentValidTime {
		return subsequentDelay
	}
	return initialDelay
}

func showAtMousePosition(canvas fyne.Canvas, pos fyne.Position, text string) *fyne.Container {
	if canvas == nil {
		return nil
	}

	lastShown = time.Now()
	l := layers[canvas]
	if l == nil {
		return nil
	}
	c := &l.Container
	if overlay := canvas.Overlays().Top(); overlay != nil {
		c = l.overlays[overlay]
		if c == nil {
			return nil
		}
	}

	t := newTip(text)
	c.Objects = []fyne.CanvasObject{t}

	zeroPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(c)
	sizeAndPosition(zeroPos, pos, t, canvas)
	c.Refresh()
	return c
}

func hide(c *fyne.Container) {
	if c != nil {
		c.Objects = nil
		c.Refresh()
	}
}

func sizeAndPosition(zeroPos, pos fyne.Position, t *tip, canvas fyne.Canvas) {
	canvasSize := canvas.Size()
	pad := theme.Padding()

	w := fyne.Min(t.textWidth(), fyne.Min(canvasSize.Width-pad*2, maxWidth))
	t.Resize(fyne.NewSize(w, 1))
	t.Resize(fyne.NewSize(w, t.textMinSize().Height))

	if rightEdge := pos.X + w; rightEdge > canvasSize.Width-pad {
		pos.X -= rightEdge - canvasSize.Width + pad
	}
	if bottomEdge := pos.Y + t.Size().Height + belowMouseDistance; bottomEdge > canvasSize.Height-pad {
		pos.Y -= t.Size().Height + aboveMouseDistance
	} else {
		pos.Y += belowMouseDistance
	}

	scale := canvas.Scale()
	t.Move(fyne.NewPos(ui.SnapToPixel(pos.X, scale)-zeroPos.X, ui.SnapToPixel(pos.Y, scale)-zeroPos.Y))
}
