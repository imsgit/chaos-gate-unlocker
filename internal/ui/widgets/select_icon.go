package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SelectIcon struct {
	widget.BaseWidget

	icon *Icon
	sel  *Select
}

func NewSelectIcon() *SelectIcon {
	s := &SelectIcon{
		icon: NewIcon(),
		sel:  NewSelect(),
	}

	s.ExtendBaseWidget(s)
	return s
}

func (s *SelectIcon) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return fyne.NewSize(0, 54)
}

func (s *SelectIcon) SetResource(resource fyne.Resource) {
	s.icon.SetResource(resource)
}

func (s *SelectIcon) SetToolTip(toolTip string) {
	s.icon.SetToolTip(toolTip)
	s.sel.SetToolTip(toolTip)
}

func (s *SelectIcon) SetPlaceHolder(placeholder string) {
	s.sel.PlaceHolder = placeholder
}

func (s *SelectIcon) SetOptions(options []string) {
	s.sel.SetOptions(options)
}

func (s *SelectIcon) SetSelected(text string) {
	s.sel.SetSelected(text)
}

func (s *SelectIcon) Selected() string {
	return s.sel.Selected
}

func (s *SelectIcon) OnChanged(fn func(newVal string)) {
	s.sel.OnChanged = fn
}

func (s *SelectIcon) OnBeforeShowPopup(fn func()) {
	s.sel.OnBeforeShowPopup = fn
}

func (s *SelectIcon) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewPadded(
			container.NewBorder(nil, nil, s.icon, nil,
				container.NewPadded(
					container.NewVBox(s.sel)))))
}
