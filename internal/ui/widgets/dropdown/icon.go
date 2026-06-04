package dropdown

import (
	iconw "chaos-gate-unlocker/internal/ui/widgets/icon"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DropdownIconWidget struct {
	widget.BaseWidget

	icon *iconw.IconWidget
	sel  *DropdownWidget
}

func NewDropdownIconWidget() *DropdownIconWidget {
	s := &DropdownIconWidget{
		icon: iconw.NewIconWidget(),
		sel:  NewDropdownWidget(),
	}

	container.NewThemeOverride(s.sel, selectPadTheme{fyne.CurrentApp().Settings().Theme()})

	s.ExtendBaseWidget(s)
	return s
}

func (s *DropdownIconWidget) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return fyne.NewSize(0, 54)
}

func (s *DropdownIconWidget) SetResource(resource fyne.Resource) {
	s.icon.SetResource(resource)
}

func (s *DropdownIconWidget) SetToolTip(toolTip string) {
	s.icon.SetToolTip(toolTip)
	s.sel.SetToolTip(toolTip)
}

func (s *DropdownIconWidget) SetPlaceHolder(placeholder string) {
	s.sel.PlaceHolder = placeholder
}

func (s *DropdownIconWidget) SetOptions(options []string) {
	s.sel.SetOptions(options)
}

func (s *DropdownIconWidget) SetSelected(text string) {
	s.sel.SetSelected(text)
}

func (s *DropdownIconWidget) Selected() string {
	return s.sel.Selected
}

func (s *DropdownIconWidget) OnChanged(fn func(newVal string)) {
	s.sel.OnChanged = fn
}

func (s *DropdownIconWidget) OnBeforeShowPopup(fn func()) {
	s.sel.OnBeforeShowPopup = fn
}

func (s *DropdownIconWidget) SetOptionToolTip(fn func(option string) string) {
	s.sel.ToolTipForOption = fn
}

func (s *DropdownIconWidget) SetOptionIcon(fn func(option string) fyne.Resource) {
	s.sel.IconForOption = fn
}

func (s *DropdownIconWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewPadded(
			container.NewBorder(nil, nil, s.icon, nil,
				container.NewPadded(
					container.NewVBox(s.sel)))))
}
