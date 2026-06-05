package tooltip

import (
	"errors"
	"math"
	"time"

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

var (
	layers        = make(map[fyne.Canvas]*layer)
	lastShownUnix int64
)

type handle struct {
	canvas  fyne.Canvas
	overlay fyne.CanvasObject
}

type layer struct {
	Container fyne.Container
	overlays  map[fyne.CanvasObject]*layer
}

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
	l := &layer{}
	if parent.overlays == nil {
		parent.overlays = make(map[fyne.CanvasObject]*layer)
	}
	parent.overlays[overlay] = l
	return &l.Container
}

func RemoveOverlayToolTipLayer(overlay fyne.CanvasObject, canvas fyne.Canvas) {
	if parent := layers[canvas]; parent != nil {
		delete(parent.overlays, overlay)
	}
}

func nextDelay() time.Duration {
	if time.Now().UnixMilli()-lastShownUnix < subsequentValidTime.Milliseconds() {
		return subsequentDelay
	}
	return initialDelay
}

func showAtMousePosition(canvas fyne.Canvas, pos fyne.Position, text string) *handle {
	if canvas == nil {
		return nil
	}

	lastShownUnix = time.Now().UnixMilli()
	overlay := canvas.Overlays().Top()
	h := &handle{canvas: canvas, overlay: overlay}
	l := findLayer(h)
	if l == nil {
		return nil
	}

	t := newTip(text)
	l.Container.Objects = []fyne.CanvasObject{t}

	zeroPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(&l.Container)

	sizeAndPosition(zeroPos, pos.Subtract(zeroPos), t, canvas)
	l.Container.Refresh()
	return h
}

func hide(h *handle) {
	if h == nil {
		return
	}
	if l := findLayer(h); l != nil {
		l.Container.Objects = nil
		l.Container.Refresh()
	}
}

func findLayer(h *handle) *layer {
	l := layers[h.canvas]
	if l == nil {
		return nil
	}
	if h.overlay != nil {
		return l.overlays[h.overlay]
	}
	return l
}

func sizeAndPosition(zeroPos, relPos fyne.Position, t *tip, canvas fyne.Canvas) {
	canvasSize := canvas.Size()
	pad := theme.Padding()

	w := fyne.Min(t.textWidth(), fyne.Min(canvasSize.Width-pad*2, maxWidth))
	t.Resize(fyne.NewSize(w, 1))
	t.Resize(fyne.NewSize(w, t.textMinSize().Height))

	if rightEdge := relPos.X + zeroPos.X + w; rightEdge > canvasSize.Width-pad {
		relPos.X -= rightEdge - canvasSize.Width + pad
	}
	if bottomEdge := relPos.Y + zeroPos.Y + t.Size().Height + belowMouseDistance; bottomEdge > canvasSize.Height-pad {
		relPos.Y -= t.Size().Height + aboveMouseDistance
	} else {
		relPos.Y += belowMouseDistance
	}

	if scale := canvas.Scale(); scale > 0 {
		relPos.X = float32(math.Round(float64((relPos.X+zeroPos.X)*scale)))/scale - zeroPos.X
		relPos.Y = float32(math.Round(float64((relPos.Y+zeroPos.Y)*scale)))/scale - zeroPos.Y
	}

	t.Move(relPos)
}
