package anim

import (
	"bytes"
	"context"
	"image"
	_ "image/png"
	"math"
	"sync"
	"time"

	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/widgets/progress"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

type Aquila struct {
	leftRes, rightRes fyne.Resource

	once        sync.Once
	left, right []image.Image
}

func NewAquila(left, right fyne.Resource) *Aquila {
	return &Aquila{leftRes: left, rightRes: right}
}

func (a *Aquila) Prewarm() { go a.frames() }

func (a *Aquila) frames() (left, right []image.Image) {
	a.once.Do(func() {
		a.left = aquilaFrames(a.leftRes, 1.0, -30)
		a.right = aquilaFrames(a.rightRes, 0.0, 30)
	})
	return a.left, a.right
}

const (
	aquilaBaseW      = 200
	aquilaFrameCount = 16
)

func aquilaFrames(res fyne.Resource, pivotX, fromDeg float64) []image.Image {
	src, _, err := image.Decode(bytes.NewReader(res.Content()))
	if err != nil {
		return nil
	}
	src = ui.ScaleDown(src, aquilaBaseW)

	b := src.Bounds()
	px := float64(b.Min.X) + pivotX*float64(b.Dx())
	py := float64(b.Min.Y) + float64(b.Dy())/2

	frames := make([]image.Image, aquilaFrameCount)
	for i := range frames {
		t := float64(i) / float64(aquilaFrameCount-1)
		sin, cos := math.Sincos(fromDeg * (1 - t) * math.Pi / 180)

		m := f64.Aff3{
			cos, -sin, px - cos*px + sin*py,
			sin, cos, py - sin*px - cos*py,
		}

		dst := image.NewRGBA(b)
		xdraw.CatmullRom.Transform(dst, m, src, b, xdraw.Over, nil)
		frames[i] = dst
	}
	return frames
}

func (a *Aquila) Animate(ctx context.Context, im, im2 *canvas.Image, p *progress.Widget, open bool) {
	leftFrames, rightFrames := a.frames()

	width := float32(0)
	sOffset := p.Size().Width / 20

	if open {
		im.Translucency = 1
		im2.Translucency = 1
		if len(leftFrames) > 0 {
			im.Resource = nil
			im2.Resource = nil
			im.Image = leftFrames[0]
			im2.Image = rightFrames[0]
		}
	}

	fyne.DoAndWait(p.Reset)

	runFrames(ctx, 30, 15*time.Millisecond, p.Reset, func(i int) {
		if i < 20 {
			width += sOffset
			p.Grow(width)
		} else if open {
			p.Complete()
		} else {
			p.Reset()
		}

		if i > 4 {
			t := float64(i-5) / float64(29-5)
			if open {
				im.Translucency = 1 - t*t
				im2.Translucency = im.Translucency
				if len(leftFrames) > 0 {
					idx := int(t*float64(len(leftFrames)-1) + 0.5)
					im.Image = leftFrames[idx]
					im2.Image = rightFrames[idx]
				}
			} else {
				im.Translucency = t
				im2.Translucency = t
			}
			im.Refresh()
			im2.Refresh()
		}
	})
}
