package toggle

import (
	"image"
	"image/color"
	"image/draw"
	"sync"

	"chaos-gate-unlocker/internal/ui"
)

var (
	staticOff    = ui.DecodeMasked(ui.WidgetSwitchOffIcon())
	staticOn     = ui.DecodeMasked(ui.WidgetSwitchOnIcon())
	switchFrames = sync.OnceValue(buildSwitchFrames)
)

func buildSwitchFrames() []image.Image {
	off, on := staticOff, staticOn
	if off == nil || on == nil {
		return nil
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
	return frames
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
	mask := image.NewUniform(color.Alpha{A: uint8(255 * f)})
	draw.DrawMask(dst, dst.Bounds(), src, src.Bounds().Min, mask, image.Point{}, draw.Over)
}

func drawShifted(dst *image.RGBA, src image.Image, dx int) {
	b := src.Bounds()
	target := image.Rect(b.Min.X+dx, b.Min.Y, b.Max.X+dx, b.Max.Y)
	draw.Draw(dst, target, src, b.Min, draw.Over)
}
