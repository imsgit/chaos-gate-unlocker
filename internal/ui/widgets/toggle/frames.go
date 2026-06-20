package toggle

import (
	"image"
	"image/color"
	"image/draw"
	"sync"

	"chaos-gate-unlocker/internal/ui"
)

var (
	switchFramesOnce sync.Once
	switchFrames     []image.Image

	staticFramesOnce sync.Once
	staticImgs       switchStatics

	prewarmOnce sync.Once
)

type switchStatics struct {
	off, on image.Image
}

func getSwitchFrames() []image.Image {
	switchFramesOnce.Do(buildSwitchFrames)
	return switchFrames
}

func getStaticFrames() switchStatics {
	staticFramesOnce.Do(func() {
		off, on := switchBase()
		staticImgs = switchStatics{off: off, on: on}
	})
	return staticImgs
}

func switchBase() (off, on image.Image) {
	return ui.DecodeMasked(ui.WidgetSwitchOffIcon()),
		ui.DecodeMasked(ui.WidgetSwitchOnIcon())
}

func buildSwitchFrames() {
	off, on := switchBase()
	if off == nil || on == nil {
		return
	}

	b := off.Bounds()
	r := b.Dy() / 2
	cy := b.Min.Y + r
	cxOff := b.Min.X + r
	cxOn := b.Max.X - r
	travel := cxOn - cxOff

	cog := maskCircle(off, cxOff, cy, r, true)
	blockOff := maskCircle(off, cxOff, cy, r, false)
	blockOn := maskCircle(on, cxOn, cy, r, false)

	const n = 8
	frames := make([]image.Image, n)
	frames[0] = off
	frames[n-1] = on
	for i := 1; i < n-1; i++ {
		p := float64(i) / float64(n-1)
		dst := image.NewRGBA(b)
		drawAlpha(dst, blockOff, 1-p)
		drawAlpha(dst, blockOn, p)
		drawShifted(dst, cog, int(float64(travel)*p+0.5))
		frames[i] = dst
	}
	switchFrames = frames
}

func maskCircle(src image.Image, cx, cy, r int, inside bool) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	r2 := r * r
	for y := b.Min.Y; y < b.Max.Y; y++ {
		dy := y - cy
		for x := b.Min.X; x < b.Max.X; x++ {
			dx := x - cx
			if (dx*dx+dy*dy <= r2) == inside {
				dst.Set(x, y, src.At(x, y))
			}
		}
	}
	return dst
}

func drawAlpha(dst *image.RGBA, src image.Image, f float64) {
	if f <= 0 {
		return
	}
	if f > 1 {
		f = 1
	}
	mask := image.NewUniform(color.Alpha{A: uint8(255 * f)})
	draw.DrawMask(dst, dst.Bounds(), src, src.Bounds().Min, mask, image.Point{}, draw.Over)
}

func drawShifted(dst *image.RGBA, src image.Image, dx int) {
	b := src.Bounds()
	target := image.Rect(b.Min.X+dx, b.Min.Y, b.Max.X+dx, b.Max.Y)
	draw.Draw(dst, target, src, b.Min, draw.Over)
}
