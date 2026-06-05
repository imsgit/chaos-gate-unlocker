package toggle

import (
	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/anim"
	"chaos-gate-unlocker/internal/ui/widgets/snapimage"

	"context"
	"image"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Widget struct {
	widget.DisableableWidget

	icon *snapimage.Widget
	sw   *snapimage.Widget

	animCancel context.CancelFunc

	textName *canvas.Text

	onChanged func(on bool)
	on        bool
	focused   bool
}

func New(onChanged func(on bool), icon, name, toolTip string) *Widget {
	s := &Widget{
		icon:      snapimage.New(fyne.NewSize(46, 46)),
		sw:        snapimage.New(fyne.NewSize(46, 46)),
		textName:  canvas.NewText("", color.White),
		onChanged: onChanged,
	}

	s.textName.Text = name

	s.icon.SetToolTip(toolTip)
	s.sw.SetToolTip(toolTip)

	s.icon.SetResource(ui.GetIconByName(icon))

	prewarmOnce.Do(func() { go getSwitchFrames() })

	s.ExtendBaseWidget(s)
	s.showStatic()
	return s
}

func (s *Widget) SetToolTip(toolTip string) {
	s.icon.SetToolTip(toolTip)
	s.sw.SetToolTip(toolTip)
}

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

	c, style, dim := color.Color(color.White), fyne.TextStyle{Bold: true}, 0.0
	if s.Disabled() {
		c = fyne.CurrentApp().Settings().Theme().Color(theme.ColorNamePrimary, 0)
		style = fyne.TextStyle{}
		dim = 0.5
	}
	if s.textName.Color != c || s.textName.TextStyle != style {
		s.textName.Color = c
		s.textName.TextStyle = style
		s.textName.Refresh()
	}
	s.icon.SetTranslucency(dim)
	s.sw.SetTranslucency(dim)

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

func (s *Widget) showStatic() {
	st := getStaticFrames()

	var img image.Image
	switch {
	case s.on && s.Disabled():
		img = st.onHover
	case s.on:
		img = st.on
	case s.Disabled():
		img = st.offHover
	default:
		img = st.off
	}
	s.sw.SetImage(img)
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
		s.sw.SetImage(frames[0])
	} else {
		s.sw.SetImage(frames[n-1])
	}

	s.animCancel = anim.Frames(n, 16*time.Millisecond, s.showStatic, func(i int) {
		idx := i - 1
		if !on {
			idx = n - i
		}
		s.sw.SetImage(frames[idx])
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
	iconContainer := container.NewPadded(s.icon)
	nameContainer := container.NewPadded(s.textName)

	return widget.NewSimpleRenderer(container.NewHBox(
		iconContainer,
		s.sw,
		nameContainer,
	))
}
