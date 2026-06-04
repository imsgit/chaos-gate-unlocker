package toggle

import (
	"chaos-gate-unlocker/internal/ui"
	iconw "chaos-gate-unlocker/internal/ui/widgets/icon"

	"context"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ToggleWidget struct {
	widget.DisableableWidget

	icon         *iconw.IconWidget
	onIcon       *iconw.IconWidget
	offIcon      *iconw.IconWidget
	onHoverIcon  *iconw.IconWidget
	offHoverIcon *iconw.IconWidget

	anim       *canvas.Image
	animCancel context.CancelFunc

	textName *canvas.Text

	onChanged func(on bool)
	on        bool
	focused   bool
}

func NewToggleWidget(onChanged func(on bool), icon, name, toolTip string) *ToggleWidget {
	s := &ToggleWidget{
		icon:         iconw.NewIconWidget(),
		onIcon:       iconw.NewIconWidget(),
		offIcon:      iconw.NewIconWidget(),
		onHoverIcon:  iconw.NewIconWidget(),
		offHoverIcon: iconw.NewIconWidget(),
		textName:     canvas.NewText("", color.White),
		onChanged:    onChanged,
	}

	s.textName.Text = name

	s.icon.SetToolTip(toolTip)
	s.onIcon.SetToolTip(toolTip)
	s.offIcon.SetToolTip(toolTip)

	s.icon.SetResource(ui.GetIconByName(icon))
	s.onIcon.SetResource(ui.GetWidgetSwitchOnIcon())
	s.offIcon.SetResource(ui.GetWidgetSwitchOffIcon())
	s.onHoverIcon.SetResource(ui.GetWidgetSwitchOnHoverIcon())
	s.offHoverIcon.SetResource(ui.GetWidgetSwitchOffHoverIcon())

	s.anim = &canvas.Image{FillMode: canvas.ImageFillContain}
	s.anim.SetMinSize(fyne.NewSize(46, 46))
	s.anim.Hide()

	prewarmOnce.Do(func() { go getSwitchFrames() })

	s.ExtendBaseWidget(s)
	return s
}

func (s *ToggleWidget) SetToolTip(toolTip string) {
	s.icon.SetToolTip(toolTip)
	s.onIcon.SetToolTip(toolTip)
	s.offIcon.SetToolTip(toolTip)
}

func (s *ToggleWidget) SetState(on, notify bool) {
	s.set(on, notify, false)
}

func (s *ToggleWidget) set(on, notify, animate bool) {
	changed := s.on != on
	s.on = on

	if notify && s.onChanged != nil {
		s.onChanged(on)
	}

	c, style := color.Color(color.White), fyne.TextStyle{Bold: true}
	if s.Disabled() {
		c = fyne.CurrentApp().Settings().Theme().Color(theme.ColorNamePrimary, 0)
		style = fyne.TextStyle{}
	}
	if s.textName.Color != c || s.textName.TextStyle != style {
		s.textName.Color = c
		s.textName.TextStyle = style
		s.textName.Refresh()
	}

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

func (s *ToggleWidget) showStatic() {
	s.anim.Hide()
	s.onIcon.Hide()
	s.onHoverIcon.Hide()
	s.offIcon.Hide()
	s.offHoverIcon.Hide()

	switch {
	case s.on && s.Disabled():
		s.onHoverIcon.Show()
	case s.on:
		s.onIcon.Show()
	case s.Disabled():
		s.offHoverIcon.Show()
	default:
		s.offIcon.Show()
	}
}

func (s *ToggleWidget) animateTo(on bool) {
	frames := getSwitchFrames()
	if len(frames) == 0 {
		s.showStatic()
		return
	}

	if s.animCancel != nil {
		s.animCancel()
	}

	s.onIcon.Hide()
	s.onHoverIcon.Hide()
	s.offIcon.Hide()
	s.offHoverIcon.Hide()

	n := len(frames)
	if on {
		s.anim.Image = frames[0]
	} else {
		s.anim.Image = frames[n-1]
	}
	s.anim.Show()

	s.animCancel = ui.Frames(n, 16*time.Millisecond, s.showStatic, func(i int) {
		idx := i - 1
		if !on {
			idx = n - i
		}
		s.anim.Image = frames[idx]
		s.anim.Refresh()
	})
}

func (s *ToggleWidget) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return fyne.NewSize(0, 54)
}

func (s *ToggleWidget) FocusGained() {
	if s.Disabled() {
		return
	}

	s.focused = true
}

func (s *ToggleWidget) FocusLost() {
	s.focused = false
}

func (s *ToggleWidget) TypedRune(r rune) {
	if s.Disabled() {
		return
	}

	if r == ' ' {
		s.set(!s.on, true, true)
	}
}

func (s *ToggleWidget) TypedKey(*fyne.KeyEvent) {}

func (s *ToggleWidget) Tapped(*fyne.PointEvent) {
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

func (s *ToggleWidget) TappedSecondary(*fyne.PointEvent) {
}

func (s *ToggleWidget) CreateRenderer() fyne.WidgetRenderer {
	iconContainer := container.NewPadded(s.icon)

	switchContainer := container.NewStack(
		s.offIcon,
		s.onIcon,
		s.offHoverIcon,
		s.onHoverIcon,
		s.anim,
	)

	nameContainer := container.NewPadded(s.textName)

	return widget.NewSimpleRenderer(container.NewHBox(
		iconContainer,
		switchContainer,
		nameContainer,
	))
}
