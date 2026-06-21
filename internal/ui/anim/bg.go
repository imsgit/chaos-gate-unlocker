package anim

import (
	"context"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2/canvas"
)

const (
	eyeGlowPeriod = 30 * time.Second
	eyeGlowSteps  = 180
	eyeGlowFrame  = 33 * time.Millisecond
	eyeGlowRise   = 0.4
	eyeGlowBase   = 0.46
	eyeGlowFlash  = 0.85
	eyeBloomBlur  = 9
)

func AnimateAbout(ctx context.Context, im *canvas.Image) {
	runFrames(ctx, 30, 15*time.Millisecond, nil, func(i int) {
		if i > 5 {
			im.Translucency = clamp(im.Translucency - 0.04)
			canvas.Refresh(im)
		}
	})
}

type glowEye struct{ cx, cy, rx, ry, reveal, bloom, lo, hi, grad float64 }

type glowSpot struct {
	off int
	p   [4]float64
}

type EyeGlow struct {
	src  image.Image
	tint color.RGBA
	eyes []glowEye

	once  sync.Once
	img   *image.RGBA
	cv    *canvas.Image
	spots []glowSpot
	flick []flickerStep

	animOnce sync.Once
}

func NewEyeGlow(src image.Image) *EyeGlow {
	return &EyeGlow{
		src:  src,
		tint: color.RGBA{R: 240, G: 244, B: 255},
		eyes: []glowEye{
			{272, 31, 38, 32, 0.22, 0, 100, 150, 0},
			{190, 60, 54, 40, 0.20, 0, 112, 170, 0},
			{356, 60, 54, 40, 0.20, 0, 112, 170, 0},
			{275, 146, 18, 20, 0.28, 0, 100, 150, 0},
			{221, 287, 31, 37, 0.45, 0.28, 100, 150, 0},
			{138, 232, 66, 60, 0.55, 0.25, 102, 155, 0},
			{410, 232, 70, 60, 0.55, 0.25, 102, 155, 0},
			{275, 451, 255, 38, 0.07, 0.03, 122, 205, 0.1},
			{275, 501, 108, 14, 0.10, 0.04, 120, 200, 0},
			{179, 188, 56, 30, 0.90, 0.5, 95, 205, 0},
			{366, 191, 56, 30, 0.90, 0.5, 95, 205, 0},
		},
	}
}

func (g *EyeGlow) build() {
	g.once.Do(func() {
		src := g.src
		if src == nil {
			return
		}

		b := src.Bounds()
		g.img = image.NewRGBA(b)
		g.cv = &canvas.Image{Image: g.img, FillMode: canvas.ImageFillContain, ScaleMode: canvas.ImageScaleFastest}

		tintR, tintG, tintB := float64(g.tint.R), float64(g.tint.G), float64(g.tint.B)

		type eyeBox struct{ x0, y0, x1, y1, w, h int }
		boxes := make([]eyeBox, len(g.eyes))
		maxN := 0
		for i, e := range g.eyes {
			margin := 0
			if e.bloom > 0 {
				margin = eyeBloomBlur + 2
			}
			x0 := max(int(e.cx-e.rx)-margin, b.Min.X)
			y0 := max(int(e.cy-e.ry)-margin, b.Min.Y)
			x1 := min(int(e.cx+e.rx)+1+margin, b.Max.X)
			y1 := min(int(e.cy+e.ry)+1+margin, b.Max.Y)
			w, h := x1-x0, y1-y0
			boxes[i] = eyeBox{x0, y0, x1, y1, w, h}
			if w > 0 && h > 0 && w*h > maxN {
				maxN = w * h
			}
		}
		if maxN == 0 {
			return
		}
		lensBuf := make([]float32, maxN)
		haloBuf := make([]float32, maxN)
		blurTmp := make([]float32, maxN)
		rBuf := make([]float32, maxN)
		gBuf := make([]float32, maxN)
		bBuf := make([]float32, maxN)

		for i, e := range g.eyes {
			bx := boxes[i]
			w, h := bx.w, bx.h
			if w <= 0 || h <= 0 {
				continue
			}
			x0, y0, x1, y1 := bx.x0, bx.y0, bx.x1, bx.y1

			lens := lensBuf[:w*h]
			for j := range lens {
				lens[j] = 0
			}
			for y := y0; y < y1; y++ {
				for x := x0; x < x1; x++ {
					dx := (float64(x) - e.cx) / e.rx
					dy := (float64(y) - e.cy) / e.ry
					win := glowWindow(dx*dx + dy*dy)
					if e.grad > 0 {
						win = gradWindow(dx, dy, e.grad)
					}
					if win <= 0 {
						continue
					}
					lr, lg, lb, _ := src.At(x, y).RGBA()
					R, G, B := float64(lr>>8), float64(lg>>8), float64(lb>>8)
					li := (y-y0)*w + (x - x0)
					rBuf[li], gBuf[li], bBuf[li] = float32(R), float32(G), float32(B)
					lens[li] = float32(win * lensMask(R, G, B, e.lo, e.hi))
				}
			}

			var halo []float32
			if e.bloom > 0 {
				halo = blurMask(lens, haloBuf[:w*h], blurTmp[:w*h], w, h, eyeBloomBlur)
			}

			for y := y0; y < y1; y++ {
				for x := x0; x < x1; x++ {
					idx := (y-y0)*w + (x - x0)
					aReveal := e.reveal * float64(lens[idx])
					aBloom := 0.0
					if halo != nil {
						aBloom = e.bloom * float64(halo[idx])
					}
					a := aReveal + aBloom
					if a*255 < 1.5 {
						continue
					}
					if a > 1 {
						a = 1
					}
					R, G, B := float64(rBuf[idx]), float64(gBuf[idx]), float64(bBuf[idx])
					g.spots = append(g.spots, glowSpot{
						off: g.img.PixOffset(x, y),
						p: [4]float64{
							R*aReveal + tintR*aBloom,
							G*aReveal + tintG*aBloom,
							B*aReveal + tintB*aBloom,
							a * 255,
						},
					})
				}
			}
		}
	})
}

func smoothstep(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	return t * t * (3 - 2*t)
}

func glowWindow(d float64) float64 {
	return 1 - smoothstep((d-0.7)/0.3)
}

func gradWindow(dx, dy, grad float64) float64 {
	if dy < -1 || dy > 1 {
		return 0
	}
	h := 1 - smoothstep((math.Abs(dx)-0.75)/0.25)
	fy := (dy + 1) / 2
	return h * (grad + (1-grad)*math.Pow(1-fy, 1.6))
}

func blurMask(m, out, tmp []float32, w, h, radius int) []float32 {
	if radius < 1 {
		return m
	}
	sigma := float64(radius) / 2
	k := make([]float64, 2*radius+1)
	for i := range k {
		d := float64(i - radius)
		k[i] = math.Exp(-d * d / (2 * sigma * sigma))
	}
	blurAxis(m, tmp, w, h, radius, k, false)
	blurAxis(tmp, out, w, h, radius, k, true)
	return out
}

func blurAxis(src, dst []float32, w, h, radius int, k []float64, vertical bool) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var s, ws float64
			for j := -radius; j <= radius; j++ {
				xx, yy := x, y
				if vertical {
					yy = y + j
				} else {
					xx = x + j
				}
				if xx >= 0 && xx < w && yy >= 0 && yy < h {
					s += float64(src[yy*w+xx]) * k[j+radius]
					ws += k[j+radius]
				}
			}
			dst[y*w+x] = float32(s / ws)
		}
	}
}

func lensMask(r, g, b, lo, hi float64) float64 {
	l := 0.299*r + 0.587*g + 0.114*b
	return smoothstep((l - lo) / (hi - lo))
}

func (g *EyeGlow) apply(t float64) bool {
	pix := g.img.Pix
	changed := false
	for i := range g.spots {
		s := &g.spots[i]
		for j := 0; j < 4; j++ {
			if v := clamp8(t * s.p[j]); pix[s.off+j] != v {
				pix[s.off+j] = v
				changed = true
			}
		}
	}
	return changed
}

func (g *EyeGlow) Overlay() *canvas.Image {
	g.build()
	return g.cv
}

func (g *EyeGlow) Animate() {
	if g.cv == nil {
		return
	}
	g.animOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(eyeGlowPeriod)
			defer ticker.Stop()
			var cancel context.CancelFunc
			for range ticker.C {
				if cancel != nil {
					cancel()
				}
				cancel = g.pulse()
			}
		}()
	})
}

func (g *EyeGlow) pulse() context.CancelFunc {
	flicker := g.buildFlicker()
	return Frames(eyeGlowSteps, eyeGlowFrame,
		func() {
			if g.apply(0) {
				g.cv.Refresh()
			}
		},
		func(i int) {
			fr := flicker[i]
			breath := 1 + 0.13*math.Sin(float64(i)*0.07) + 0.04*math.Sin(float64(i)*0.5)
			t := eyeGlowEnvelope(i) * eyeGlowBase * fr.gate * breath
			if s := fr.spark * eyeGlowFlash; s > t {
				t = s
			}
			if g.apply(t) {
				g.cv.Refresh()
			}
		})
}

func eyeGlowEnvelope(i int) float64 {
	p := float64(i) / float64(eyeGlowSteps)
	if p < eyeGlowRise {
		return 0.5 - 0.5*math.Cos(math.Pi*p/eyeGlowRise)
	}
	x := (p - eyeGlowRise) / (1 - eyeGlowRise)
	e := math.Exp(-3 * x)
	if p > 0.9 {
		e *= 1 - (p-0.9)/0.1
	}
	return e
}

type flickerStep struct {
	spark float64
	gate  float64
}

func (g *EyeGlow) buildFlicker() []flickerStep {
	if g.flick == nil {
		g.flick = make([]flickerStep, eyeGlowSteps+1)
	}
	fl := g.flick
	for i := range fl {
		fl[i] = flickerStep{gate: 1}
	}

	lo := eyeGlowSteps * 4 / 10
	hi := eyeGlowSteps * 8 / 10
	bursts := rand.Intn(4)
	for b := 0; b < bursts; b++ {
		i := lo + rand.Intn(hi-lo)
		events := 4 + rand.Intn(6)
		for k := 0; k < events && i < hi; k++ {
			dur := 1 + rand.Intn(2)
			if k%2 == 0 {
				amp := 0.6 + rand.Float64()*0.4
				for s := 0; s < dur && i < hi; s++ {
					if amp > fl[i].spark {
						fl[i].spark = amp
					}
					i++
				}
				tail := 2 + rand.Intn(3)
				for s := 0; s < tail && i < hi; s++ {
					if v := amp * math.Exp(-float64(s+1)*0.8); v > fl[i].spark {
						fl[i].spark = v
					}
					i++
				}
			} else {
				g := 0.05 + rand.Float64()*0.2
				for s := 0; s < dur && i < hi; s++ {
					fl[i].gate = g
					i++
				}
			}
			i += rand.Intn(3)
		}
	}
	return fl
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

func clamp8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
