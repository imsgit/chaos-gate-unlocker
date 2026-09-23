package toggle

import (
	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/anim"
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"image"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Widget struct {
	widget.BaseWidget
	tooltip.WidgetExtend

	icon *canvas.Image
	sw   *canvas.Image

	anim *fyne.Animation

	textName *canvas.Text

	onChanged func(on bool)
	on        bool
	focused   bool
	disabled  bool
}

func New(onChanged func(on bool), icon, name, toolTip string) *Widget {
	s := &Widget{
		icon:      ui.NewIconImage(fyne.NewSize(46, 46)),
		sw:        ui.NewIconImage(fyne.NewSize(46, 46)),
		textName:  canvas.NewText("", color.White),
		onChanged: onChanged,
	}

	s.sw.ScaleMode = canvas.ImageScaleFastest

	s.textName.Text = name
	s.SetToolTip(toolTip)

	s.icon.Image = ui.DecodeMasked(ui.IconByName(icon))

	s.ExtendBaseWidget(s)
	s.showStatic()
	return s
}

func (s *Widget) ExtendBaseWidget(wid fyne.Widget) {
	s.ExtendToolTipWidget(wid)
	s.BaseWidget.ExtendBaseWidget(wid)
}

func (s *Widget) SetState(on, notify bool) {
	s.set(on, notify, false)
}

func (s *Widget) Disabled() bool { return s.disabled }

func (s *Widget) Enable() {
	s.disabled = false
	s.set(s.on, false, false)
}

func (s *Widget) Disable() {
	s.disabled = true
	s.set(s.on, false, false)
}

func (s *Widget) set(on, notify, animate bool) {
	changed := s.on != on
	s.on = on

	if notify && s.onChanged != nil {
		s.onChanged(on)
	}

	c, dim := color.Color(color.White), 0.0
	if s.Disabled() {
		c = ui.MutedForeground
		dim = 0.5
	}
	if s.textName.Color != c {
		s.textName.Color = c
		s.textName.Refresh()
	}
	setTranslucency(s.icon, dim)
	setTranslucency(s.sw, dim)

	if animate && changed && !s.Disabled() {
		s.animateTo(on)
		return
	}

	s.stopAnim()
	s.showStatic()
}

func (s *Widget) stopAnim() {
	if s.anim != nil {
		s.anim.Stop()
		s.anim = nil
	}
}

func setTranslucency(img *canvas.Image, t float64) {
	if img.Translucency == t {
		return
	}
	img.Translucency = t
	img.Refresh()
}

func (s *Widget) setSwitch(img image.Image) {
	if s.sw.Image == img {
		return
	}
	s.sw.Image = img
	s.sw.Refresh()
}

func (s *Widget) showStatic() {
	if s.on {
		s.setSwitch(staticOn)
	} else {
		s.setSwitch(staticOff)
	}
}

func (s *Widget) animateTo(on bool) {
	frames := switchFrames()
	if len(frames) == 0 {
		s.showStatic()
		return
	}

	s.stopAnim()

	n := len(frames)
	if on {
		s.setSwitch(frames[0])
	} else {
		s.setSwitch(frames[n-1])
	}

	s.anim = anim.Steps(n, time.Duration(n)*16*time.Millisecond, func(i int) {
		if !on {
			i = n - 1 - i
		}
		s.setSwitch(frames[i])
	}, s.showStatic)
	s.anim.Start()
}

func (s *Widget) MinSize() fyne.Size {
	return fyne.NewSize(0, 54)
}

func (s *Widget) FocusGained() {
	if s.Disabled() {
		return
	}

	s.focused = true
}

func (s *Widget) FocusLost() {
	s.focused = false
}

func (s *Widget) TypedRune(r rune) {
	if s.Disabled() {
		return
	}

	if r == ' ' {
		s.set(!s.on, true, true)
	}
}

func (s *Widget) TypedKey(*fyne.KeyEvent) {}

func (s *Widget) Tapped(*fyne.PointEvent) {
	if s.Disabled() {
		return
	}

	if !s.focused && !fyne.CurrentDevice().IsMobile() {
		if c := fyne.CurrentApp().Driver().CanvasForObject(s); c != nil {
			c.Focus(s)
		}
	}

	s.set(!s.on, true, true)
}

func (s *Widget) TappedSecondary(*fyne.PointEvent) {
}

func (s *Widget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewHBox(
		container.NewPadded(s.icon),
		s.sw,
		container.NewPadded(s.textName),
	))
}

func Reset(sw *Widget, status func() (bool, bool)) {
	sw.Enable()
	sw.SetState(false, true)
	if available, state := status(); !available {
		sw.Disable()
		sw.SetState(state, false)
	}
}
