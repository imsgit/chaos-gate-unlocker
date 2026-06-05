package dropdown

import (
	"chaos-gate-unlocker/internal/ui/widgets/snapimage"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type IconWidget struct {
	widget.BaseWidget

	icon *snapimage.Widget
	sel  *Widget
}

func NewIconWidget() *IconWidget {
	s := &IconWidget{
		icon: snapimage.New(fyne.NewSize(46, 46)),
		sel:  New(),
	}

	container.NewThemeOverride(s.sel, selectPadTheme{fyne.CurrentApp().Settings().Theme()})

	s.ExtendBaseWidget(s)
	return s
}

func (s *IconWidget) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return fyne.NewSize(0, 54)
}

func (s *IconWidget) SetResource(resource fyne.Resource) {
	s.icon.SetResource(resource)
}

func (s *IconWidget) SetToolTip(toolTip string) {
	s.icon.SetToolTip(toolTip)
	s.sel.SetToolTip(toolTip)
}

func (s *IconWidget) SetPlaceHolder(placeholder string) {
	s.sel.PlaceHolder = placeholder
}

func (s *IconWidget) SetOptions(options []string) {
	s.sel.SetOptions(options)
}

func (s *IconWidget) SetSelected(text string) {
	s.sel.SetSelected(text)
}

func (s *IconWidget) Selected() string {
	return s.sel.Selected
}

func (s *IconWidget) OnChanged(fn func(newVal string)) {
	s.sel.OnChanged = fn
}

func (s *IconWidget) OnBeforeShowPopup(fn func()) {
	s.sel.OnBeforeShowPopup = fn
}

func (s *IconWidget) SetOptionToolTip(fn func(option string) string) {
	s.sel.ToolTipForOption = fn
}

func (s *IconWidget) SetOptionIcon(fn func(option string) fyne.Resource) {
	s.sel.IconForOption = fn
}

func (s *IconWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewPadded(
			container.NewBorder(nil, nil, s.icon, nil,
				container.NewPadded(
					container.NewVBox(s.sel)))))
}
