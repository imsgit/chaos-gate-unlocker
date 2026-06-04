package ui

import (
	"bytes"
	"context"
	"image"
	_ "image/png"
	"math"
	"time"

	"chaos-gate-unlocker/internal/ui/widgets/progress"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

func Frames(n int, interval time.Duration, onDone func(), step func(i int)) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		i := 0
		do := func() { step(i) }
		for i = 1; i <= n; i++ {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fyne.DoAndWait(do)
			}
		}

		if ctx.Err() == nil && onDone != nil {
			fyne.DoAndWait(onDone)
		}
		cancel()
	}()

	return cancel
}

func AquilaFrames(res fyne.Resource, pivotX, fromDeg, toDeg float64, count int) []image.Image {
	src, _, err := image.Decode(bytes.NewReader(res.Content()))
	if err != nil || count < 2 {
		return nil
	}

	b := src.Bounds()
	px := float64(b.Min.X) + pivotX*float64(b.Dx())
	py := float64(b.Min.Y) + float64(b.Dy())/2

	frames := make([]image.Image, count)
	for i := range frames {
		t := float64(i) / float64(count-1)
		sin, cos := math.Sincos((fromDeg + (toDeg-fromDeg)*t) * math.Pi / 180)

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

func AnimateTop(ctx context.Context, im, im2 *canvas.Image, p *progress.ProgressWidget, open bool, leftFrames, rightFrames []image.Image) {
	width := float32(0)
	sOffset := p.Size().Width / 20

	if open {
		im.Translucency = 1
		im2.Translucency = 1
		if len(leftFrames) > 0 {
			im.Image = leftFrames[0]
			im2.Image = rightFrames[0]
		}
	}

	fyne.DoAndWait(p.Reset)

	ticker := time.NewTicker(15 * time.Millisecond)
	defer ticker.Stop()

	i := 0
	frame := func() {
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
	}

	for ; i < 30; i++ {
		select {
		case <-ctx.Done():
			fyne.DoAndWait(p.Reset)
			return
		case <-ticker.C:
			fyne.DoAndWait(frame)
		}
	}
}

func AnimateAbout(ctx context.Context, im *canvas.Image) {
	ticker := time.NewTicker(15 * time.Millisecond)
	defer ticker.Stop()

	fade := func() {
		im.Translucency = clamp(im.Translucency - 0.04)
		im.Refresh()
	}

	for i := 0; i < 30; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if i > 5 {
				fyne.DoAndWait(fade)
			}
		}
	}
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
