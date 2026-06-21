package toggle

import (
	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/anim"
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"context"
	"image"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Widget struct {
	widget.DisableableWidget
	tooltip.WidgetExtend

	icon *canvas.Image
	sw   *canvas.Image

	animCancel context.CancelFunc

	textName *canvas.Text

	onChanged func(on bool)
	on        bool
	focused   bool
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

	prewarmOnce.Do(func() { go getSwitchFrames() })

	s.ExtendBaseWidget(s)
	s.showStatic()
	return s
}

func (s *Widget) ExtendBaseWidget(wid fyne.Widget) {
	s.ExtendToolTipWidget(wid)
	s.BaseWidget.ExtendBaseWidget(wid)
}

func (s *Widget) MouseIn(e *desktop.MouseEvent)    { s.WidgetExtend.MouseIn(e) }
func (s *Widget) MouseMoved(e *desktop.MouseEvent) { s.WidgetExtend.MouseMoved(e) }
func (s *Widget) MouseOut()                        { s.WidgetExtend.MouseOut() }

func (s *Widget) SetState(on, notify bool) {
	s.set(on, notify, false)
}

func (s *Widget) Enable() {
	s.DisableableWidget.Enable()
	s.set(s.on, false, false)
}

func (s *Widget) Disable() {
	s.DisableableWidget.Disable()
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

	if s.animCancel != nil {
		s.animCancel()
		s.animCancel = nil
	}
	s.showStatic()
}

func setTranslucency(img *canvas.Image, t float64) {
	if img.Translucency == t {
		return
	}
	img.Translucency = t
	img.Refresh()
}

func (s *Widget) setSwitch(img image.Image) {
	s.sw.Image = img
	s.sw.Refresh()
}

func (s *Widget) showStatic() {
	st := getStaticFrames()
	if s.on {
		s.setSwitch(st.on)
	} else {
		s.setSwitch(st.off)
	}
}

func (s *Widget) animateTo(on bool) {
	frames := getSwitchFrames()
	if len(frames) == 0 {
		s.showStatic()
		return
	}

	if s.animCancel != nil {
		s.animCancel()
	}

	n := len(frames)
	if on {
		s.setSwitch(frames[0])
	} else {
		s.setSwitch(frames[n-1])
	}

	s.animCancel = anim.Frames(n, 16*time.Millisecond, s.showStatic, func(i int) {
		idx := i - 1
		if !on {
			idx = n - i
		}
		s.setSwitch(frames[idx])
	})
}

func (s *Widget) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
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

	if !s.focused {
		if !fyne.CurrentDevice().IsMobile() {
			if c := fyne.CurrentApp().Driver().CanvasForObject(s); c != nil {
				c.Focus(s)
			}
		}
	}

	s.set(!s.on, true, true)
}

func (s *Widget) TappedSecondary(*fyne.PointEvent) {
}

func (s *Widget) CreateRenderer() fyne.WidgetRenderer {
	iconContainer := container.NewPadded(container.NewStack(s.icon))
	swContainer := container.NewStack(s.sw)
	nameContainer := container.NewPadded(s.textName)

	return widget.NewSimpleRenderer(container.NewHBox(
		iconContainer,
		swContainer,
		nameContainer,
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
