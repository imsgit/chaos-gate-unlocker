package widgets

import (
	"chaos-gate-unlocker/internal/ui"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Switch struct {
	widget.DisableableWidget

	icon         *Icon
	onIcon       *Icon
	offIcon      *Icon
	onHoverIcon  *Icon
	offHoverIcon *Icon

	textName *canvas.Text

	onChanged func(on bool)
	on        bool
	focused   bool
}

func NewSwitch(onChanged func(on bool), icon, name, toolTip string) *Switch {
	s := &Switch{
		icon:         NewIcon(),
		onIcon:       NewIcon(),
		offIcon:      NewIcon(),
		onHoverIcon:  NewIcon(),
		offHoverIcon: NewIcon(),
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

	s.ExtendBaseWidget(s)
	return s
}

func (s *Switch) SetState(on, notify bool) {
	s.on = on

	if notify && s.onChanged != nil {
		s.onChanged(on)
	}

	s.textName.Color = color.White
	s.textName.TextStyle = fyne.TextStyle{Bold: true}

	s.onIcon.Hide()
	s.onHoverIcon.Hide()
	s.offIcon.Hide()
	s.offHoverIcon.Hide()

	if s.on {
		if s.Disabled() {
			s.textName.Color = fyne.CurrentApp().Settings().Theme().Color(theme.ColorNamePrimary, 0)
			s.textName.TextStyle = fyne.TextStyle{Bold: false}
			s.onHoverIcon.Show()
		} else {
			s.onIcon.Show()
		}
	} else {
		if s.Disabled() {
			s.textName.Color = fyne.CurrentApp().Settings().Theme().Color(theme.ColorNamePrimary, 0)
			s.textName.TextStyle = fyne.TextStyle{Bold: false}
			s.offHoverIcon.Show()
		} else {
			s.offIcon.Show()
		}
	}

	s.textName.Refresh()
}

func (s *Switch) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return fyne.NewSize(0, 54)
}

func (s *Switch) FocusGained() {
	if s.Disabled() {
		return
	}
	s.focused = true
}

func (s *Switch) FocusLost() {
	s.focused = false
}

func (s *Switch) TypedRune(r rune) {
	if s.Disabled() {
		return
	}

	if r == ' ' {
		s.SetState(!s.on, true)
	}
}

func (s *Switch) TypedKey(*fyne.KeyEvent) {}

func (s *Switch) Tapped(*fyne.PointEvent) {
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

	s.SetState(!s.on, true)
}

func (s *Switch) TappedSecondary(*fyne.PointEvent) {
}

func (s *Switch) CreateRenderer() fyne.WidgetRenderer {
	iconContainer := container.NewPadded(s.icon)

	switchContainer := container.NewStack(
		s.offIcon,
		s.onIcon,
		s.offHoverIcon,
		s.onHoverIcon,
	)

	nameContainer := container.NewPadded(s.textName)

	return widget.NewSimpleRenderer(container.NewHBox(
		iconContainer,
		switchContainer,
		nameContainer,
	))
}
