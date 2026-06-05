package anim

import (
	"bytes"
	"context"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

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

const (
	eyeGlowPeriod = 30 * time.Second
	eyeGlowSteps  = 180
	eyeGlowFrame  = 33 * time.Millisecond
	eyeGlowRise   = 0.4
	eyeGlowBase   = 0.46
	eyeGlowFlash  = 0.85
	eyeBloomBlur  = 9
)

type glowEye struct{ cx, cy, rx, ry, reveal, bloom, lo, hi, grad float64 }

type glowSpot struct {
	off            int
	pr, pg, pb, pa float64
}

type EyeGlow struct {
	res  fyne.Resource
	tint color.RGBA
	eyes []glowEye

	once  sync.Once
	img   *image.RGBA
	cv    *canvas.Image
	spots []glowSpot
	flick []flickerStep
}

func NewEyeGlow(bg fyne.Resource) *EyeGlow {
	return &EyeGlow{
		res:  bg,
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
		src, _, err := image.Decode(bytes.NewReader(g.res.Content()))
		if err != nil {
			return
		}

		b := src.Bounds()
		g.img = image.NewRGBA(b)
		g.cv = &canvas.Image{Image: g.img, FillMode: canvas.ImageFillContain, ScaleMode: canvas.ImageScaleFastest}

		tintR, tintG, tintB := float64(g.tint.R), float64(g.tint.G), float64(g.tint.B)
		for _, e := range g.eyes {
			margin := 0
			if e.bloom > 0 {
				margin = eyeBloomBlur + 2
			}
			x0 := max(int(e.cx-e.rx)-margin, b.Min.X)
			y0 := max(int(e.cy-e.ry)-margin, b.Min.Y)
			x1 := min(int(e.cx+e.rx)+1+margin, b.Max.X)
			y1 := min(int(e.cy+e.ry)+1+margin, b.Max.Y)
			w, h := x1-x0, y1-y0
			if w <= 0 || h <= 0 {
				continue
			}

			lens := make([]float64, w*h)
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
					lens[(y-y0)*w+(x-x0)] = win * lensMask(float64(lr>>8), float64(lg>>8), float64(lb>>8), e.lo, e.hi)
				}
			}

			var halo []float64
			if e.bloom > 0 {
				halo = blurMask(lens, w, h, eyeBloomBlur)
			}

			for y := y0; y < y1; y++ {
				for x := x0; x < x1; x++ {
					idx := (y-y0)*w + (x - x0)
					aReveal := e.reveal * lens[idx]
					aBloom := 0.0
					if halo != nil {
						aBloom = e.bloom * halo[idx]
					}
					a := aReveal + aBloom
					if a*255 < 1.5 {
						continue
					}
					if a > 1 {
						a = 1
					}
					lr, lg, lb, _ := src.At(x, y).RGBA()
					R, G, B := float64(lr>>8), float64(lg>>8), float64(lb>>8)
					g.spots = append(g.spots, glowSpot{
						off: g.img.PixOffset(x, y),
						pr:  R*aReveal + tintR*aBloom,
						pg:  G*aReveal + tintG*aBloom,
						pb:  B*aReveal + tintB*aBloom,
						pa:  a * 255,
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

func blurMask(m []float64, w, h, radius int) []float64 {
	if radius < 1 {
		return m
	}
	sigma := float64(radius) / 2
	k := make([]float64, 2*radius+1)
	for i := range k {
		d := float64(i - radius)
		k[i] = math.Exp(-d * d / (2 * sigma * sigma))
	}

	tmp := make([]float64, len(m))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var s, ws float64
			for j := -radius; j <= radius; j++ {
				if xx := x + j; xx >= 0 && xx < w {
					s += m[y*w+xx] * k[j+radius]
					ws += k[j+radius]
				}
			}
			tmp[y*w+x] = s / ws
		}
	}

	out := make([]float64, len(m))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var s, ws float64
			for j := -radius; j <= radius; j++ {
				if yy := y + j; yy >= 0 && yy < h {
					s += tmp[yy*w+x] * k[j+radius]
					ws += k[j+radius]
				}
			}
			out[y*w+x] = s / ws
		}
	}
	return out
}

func lensMask(r, g, b, lo, hi float64) float64 {
	l := 0.299*r + 0.587*g + 0.114*b
	return smoothstep((l - lo) / (hi - lo))
}

func (g *EyeGlow) apply(t float64) {
	pix := g.img.Pix
	for i := range g.spots {
		s := &g.spots[i]
		pix[s.off] = clamp8(t * s.pr)
		pix[s.off+1] = clamp8(t * s.pg)
		pix[s.off+2] = clamp8(t * s.pb)
		pix[s.off+3] = clamp8(t * s.pa)
	}
}

func (g *EyeGlow) Overlay() *canvas.Image {
	g.build()
	return g.cv
}

func (g *EyeGlow) Start() {
	if g.cv == nil {
		return
	}
	go func() {
		ticker := time.NewTicker(eyeGlowPeriod)
		defer ticker.Stop()
		for range ticker.C {
			g.pulse()
		}
	}()
}

func (g *EyeGlow) pulse() {
	flicker := g.buildFlicker()
	Frames(eyeGlowSteps, eyeGlowFrame,
		func() { g.apply(0); g.cv.Refresh() },
		func(i int) {
			fr := flicker[i]
			breath := 1 + 0.13*math.Sin(float64(i)*0.07) + 0.04*math.Sin(float64(i)*0.5)
			t := eyeGlowEnvelope(i) * eyeGlowBase * fr.gate * breath
			if s := fr.spark * eyeGlowFlash; s > t {
				t = s
			}
			g.apply(t)
			g.cv.Refresh()
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
