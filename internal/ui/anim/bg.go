package anim

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	eyeGlowPeriod = 30 * time.Second
	eyeGlowSteps  = 360
	eyeGlowFrame  = 16 * time.Millisecond
	eyeGlowRise   = 0.4
	eyeGlowBase   = 0.46
	eyeGlowFlash  = 0.85
	eyeBloomBlur  = 9
)

func AnimateAbout(cover *canvas.Rectangle, base color.NRGBA) *fyne.Animation {
	an := Steps(30, 30*15*time.Millisecond, func(i int) {
		if i <= 5 {
			return
		}
		c := base
		c.A = uint8(float64(base.A) * (1 - float64(i-5)/24))
		cover.FillColor = c
		cover.Refresh()
	}, nil)
	an.Start()
	return an
}

type glowEye struct{ cx, cy, rx, ry, reveal, bloom, lo, hi, grad float64 }

var glowEyes = []glowEye{
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
}

const eyeGlowShaderBody = `
uniform vec2 frame;
uniform vec4 bounds;
uniform sampler2D glow;
uniform float intensity;
uniform float aspect;

void main() {
	vec2 size = vec2(bounds[2] - bounds[0], bounds[3] - bounds[1]);
	if (size.x <= 0.0 || size.y <= 0.0) {
		discard;
	}
	vec2 p = vec2(gl_FragCoord.x - bounds[0], frame.y - gl_FragCoord.y - bounds[1]);
	vec2 fit = size;
	if (size.x / size.y > aspect) {
		fit.x = size.y * aspect;
	} else {
		fit.y = size.x / aspect;
	}
	vec2 uv = (p - (size - fit) * 0.5) / fit;
	if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0) {
		discard;
	}
	vec4 c = texture2D(glow, uv);
	vec3 q = min(2.0 * c.rgb * intensity, 1.0);
	float a = max(c.a * intensity, max(q.r, max(q.g, q.b)));
	gl_FragColor = vec4(q / max(a, 0.0001), a);
}
`

const (
	eyeGlowShader = "#version 110\n" + eyeGlowShaderBody

	eyeGlowShaderES = "#version 100\n" +
		"#ifdef GL_ES\n" +
		"# ifdef GL_FRAGMENT_PRECISION_HIGH\n" +
		"precision highp float;\n" +
		"# else\n" +
		"precision mediump float;\n" +
		"#endif\n" +
		"precision mediump int;\n" +
		"precision lowp sampler2D;\n" +
		"#endif\n" + eyeGlowShaderBody
)

type EyeGlow struct {
	src  image.Image
	kick chan struct{}

	once sync.Once
	img  *image.RGBA
	sh   *canvas.Shader
	box  *fyne.Container

	animOnce sync.Once
}

func NewEyeGlow(src image.Image) *EyeGlow {
	return &EyeGlow{src: src, kick: make(chan struct{}, 1)}
}

func (g *EyeGlow) Flash() {
	select {
	case g.kick <- struct{}{}:
	default:
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

		const tintR, tintG, tintB = 240.0, 244.0, 255.0

		for _, e := range glowEyes {
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

			lens := make([]float32, w*h)
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
					lens[(y-y0)*w+(x-x0)] = float32(win * lensMask(float64(lr>>8), float64(lg>>8), float64(lb>>8), e.lo, e.hi))
				}
			}

			var halo []float32
			if e.bloom > 0 {
				halo = blurMask(lens, w, h, eyeBloomBlur)
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
					lr, lg, lb, _ := src.At(x, y).RGBA()
					R, G, B := float64(lr>>8), float64(lg>>8), float64(lb>>8)
					off := g.img.PixOffset(x, y)
					g.img.Pix[off] = clamp8((R*aReveal + tintR*aBloom) * 0.5)
					g.img.Pix[off+1] = clamp8((G*aReveal + tintG*aBloom) * 0.5)
					g.img.Pix[off+2] = clamp8((B*aReveal + tintB*aBloom) * 0.5)
					g.img.Pix[off+3] = clamp8(a * 255)
				}
			}
		}

		g.sh = canvas.NewShader("eyeglow", []byte(eyeGlowShader), []byte(eyeGlowShaderES))
		g.sh.Textures = map[string]image.Image{"glow": g.img}
		g.sh.Uniforms = map[string]float32{
			"intensity": 0,
			"aspect":    float32(b.Dx()) / float32(b.Dy()),
		}
		g.box = container.New(fillLayout{}, g.sh)
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

func blurMask(m []float32, w, h, radius int) []float32 {
	sigma := float64(radius) / 2
	k := make([]float64, 2*radius+1)
	for i := range k {
		d := float64(i - radius)
		k[i] = math.Exp(-d * d / (2 * sigma * sigma))
	}
	tmp := make([]float32, w*h)
	out := make([]float32, w*h)
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

func (g *EyeGlow) set(t float64) {
	v := float32(t)
	if g.sh.Uniforms["intensity"] == v {
		return
	}
	g.sh.Uniforms["intensity"] = v
	if !hidden() {
		g.show()
	}
}

func (g *EyeGlow) show() {
	if g.sh.Uniforms["intensity"] <= 0 {
		g.sh.Hide()
		return
	}
	if g.sh.Visible() {
		g.sh.Refresh()
		return
	}
	g.sh.Show()
	g.box.Refresh()
}

func (g *EyeGlow) Overlay() fyne.CanvasObject {
	g.build()
	return g.box
}

type fillLayout struct{}

func (fillLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Move(fyne.Position{})
		o.Resize(size)
	}
}

func (fillLayout) MinSize([]fyne.CanvasObject) fyne.Size { return fyne.Size{} }

func (g *EyeGlow) Animate() {
	if g.sh == nil {
		return
	}
	g.animOnce.Do(func() {
		onShown(func() {
			fyne.Do(func() {
				if g.sh != nil {
					g.show()
				}
			})
		})
		go func() {
			ticker := time.NewTicker(eyeGlowPeriod)
			defer ticker.Stop()
			var cur *fyne.Animation
			for {
				select {
				case <-ticker.C:
				case <-g.kick:
					ticker.Reset(eyeGlowPeriod)
				}
				fyne.Do(func() {
					if cur != nil {
						cur.Stop()
						cur = nil
					}
					if !hidden() {
						cur = g.pulse()
						cur.Start()
					}
				})
			}
		}()
	})
}

func (g *EyeGlow) pulse() *fyne.Animation {
	flicker := buildFlicker()
	last := -1
	return Steps(eyeGlowSteps, eyeGlowSteps*eyeGlowFrame,
		func(i int) {
			fr := flicker[i]
			if i > last {
				for _, f := range flicker[last+1 : i] {
					fr.spark = max(fr.spark, f.spark)
					fr.gate = min(fr.gate, f.gate)
				}
				last = i
			}
			breath := 1 + 0.13*math.Sin(float64(i)*0.035) + 0.04*math.Sin(float64(i)*0.25)
			t := eyeGlowEnvelope(i) * eyeGlowBase * fr.gate * breath
			if s := fr.spark * eyeGlowFlash; s > t {
				t = s
			}
			g.set(t)
		},
		func() { g.set(0) })
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

func buildFlicker() []flickerStep {
	fl := make([]flickerStep, eyeGlowSteps+1)
	for i := range fl {
		fl[i].gate = 0.97 + 0.06*rand.Float64()
	}

	lo := eyeGlowSteps * 4 / 10
	hi := eyeGlowSteps * 8 / 10
	bursts := rand.Intn(3)
	for b := 0; b < bursts; b++ {
		i := lo + rand.Intn(hi-lo)
		strikes := 1 + rand.Intn(3)
		amp := 0.6 + rand.Float64()*0.4
		for k := 0; k < strikes && i < hi; k++ {
			on := 1 + rand.Intn(2)
			for s := 0; s < on && i < hi; s++ {
				if amp > fl[i].spark {
					fl[i].spark = amp
				}
				i++
			}
			for s := 0; i < hi; s++ {
				v := amp * math.Exp(-float64(s+1)*0.45)
				if v < 0.04 {
					break
				}
				if v > fl[i].spark {
					fl[i].spark = v
				}
				i++
			}
			amp *= 0.65 + rand.Float64()*0.15
			i += 2 + rand.Intn(4)
		}
		if rand.Intn(2) == 0 {
			drop := 0.05 + rand.Float64()*0.15
			for s := 0; s < 1+rand.Intn(2) && i < hi; s++ {
				fl[i].gate = drop
				i++
			}
		}
	}
	return fl
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
