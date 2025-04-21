package widgets

import (
	"reflect"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	ttwidget "github.com/dweymouth/fyne-tooltip/widget"
)

type Select struct {
	widget.Select
	ttwidget.ToolTipWidgetExtend

	OnBeforeShowPopup func()
}

func NewSelect() *Select {
	s := &Select{}

	s.ExtendBaseWidget(s)
	return s
}

func (s *Select) ExtendBaseWidget(wid fyne.Widget) {
	s.ExtendToolTipWidget(wid)
	s.Select.ExtendBaseWidget(wid)
}

func (s *Select) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return fyne.NewSize(0, 36)
}

func (s *Select) MouseIn(e *desktop.MouseEvent) {
	s.ToolTipWidgetExtend.MouseIn(e)
	s.Select.MouseIn(e)
}

func (s *Select) MouseMoved(e *desktop.MouseEvent) {
	s.ToolTipWidgetExtend.MouseMoved(e)
	s.Select.MouseMoved(e)
}

func (s *Select) MouseOut() {
	s.ToolTipWidgetExtend.MouseOut()
	s.Select.MouseOut()
}

func (s *Select) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeySpace, fyne.KeyUp, fyne.KeyDown:
		if s.OnBeforeShowPopup != nil {
			s.OnBeforeShowPopup()
		}
	}

	s.Select.TypedKey(event)

	s.disableSelectedOnPopup()
}

func (s *Select) Tapped(*fyne.PointEvent) {
	if s.OnBeforeShowPopup != nil {
		s.OnBeforeShowPopup()
	}

	s.Select.Tapped(nil)

	s.disableSelectedOnPopup()
}

func (s *Select) TappedSecondary(*fyne.PointEvent) {
}

func (s *Select) disableSelectedOnPopup() {
	f := reflect.ValueOf(s).Elem().FieldByName("popUp")
	if !f.IsValid() || f.IsNil() {
		return
	}

	ptr := unsafe.Pointer(f.UnsafeAddr())
	fVal := reflect.NewAt(f.Type(), ptr).Elem()
	popup, _ := fVal.Interface().(*widget.PopUpMenu)

	for _, i := range popup.Items {
		f = reflect.ValueOf(i).Elem().FieldByName("Item")
		v, _ := f.Interface().(*fyne.MenuItem)
		if v.Label == s.Selected {
			v.Disabled = true
			i.Refresh()
			break
		}
	}
}
