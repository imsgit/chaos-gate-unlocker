package dropdown

import (
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Widget struct {
	widget.Select
	tooltip.WidgetExtend

	OnBeforeShowPopup func()

	ToolTipForOption func(option string) string
	IconForOption    func(option string) fyne.Resource

	popup *popOverlay
}

func New() *Widget {
	s := &Widget{}

	s.ExtendBaseWidget(s)
	return s
}

func (s *Widget) ExtendBaseWidget(wid fyne.Widget) {
	s.ExtendToolTipWidget(wid)
	s.Select.ExtendBaseWidget(wid)
}

func (s *Widget) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return fyne.NewSize(0, 36)
}

func (s *Widget) MouseIn(e *desktop.MouseEvent) {
	if !tooltip.OverlayShown(s) {
		s.WidgetExtend.MouseIn(e)
	}
	s.Select.MouseIn(e)
}

func (s *Widget) MouseMoved(e *desktop.MouseEvent) {
	s.WidgetExtend.MouseMoved(e)
	s.Select.MouseMoved(e)
}

func (s *Widget) MouseOut() {
	s.WidgetExtend.MouseOut()
	s.Select.MouseOut()
}

func (s *Widget) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeySpace, fyne.KeyUp, fyne.KeyDown:
		s.showPopup()
	default:
		s.Select.TypedKey(event)
	}
}

func (s *Widget) Tapped(*fyne.PointEvent) {
	if s.Disabled() {
		return
	}
	if c := fyne.CurrentApp().Driver().CanvasForObject(s); c != nil {
		c.Focus(s)
	}
	s.showPopup()
}

func (s *Widget) TappedSecondary(*fyne.PointEvent) {
}

func (s *Widget) Hide() {
	s.hidePopup()
	s.Select.Hide()
}

func (s *Widget) showPopup() {
	s.WidgetExtend.MouseOut()

	if s.OnBeforeShowPopup != nil {
		s.OnBeforeShowPopup()
	}

	cv := fyne.CurrentApp().Driver().CanvasForObject(s)
	if cv == nil {
		return
	}

	s.hidePopup()

	rows := make([]*selectRow, len(s.Options))
	objects := make([]fyne.CanvasObject, len(s.Options))
	selectedIdx := -1
	for i := range s.Options {
		text := s.Options[i]
		disabled := text == s.Selected
		if disabled {
			selectedIdx = i
		}
		var icon fyne.Resource
		if s.IconForOption != nil {
			icon = s.IconForOption(text)
		}
		row := newSelectRow(text, icon, disabled, cv.Scale(), func() {
			s.hidePopup()
			s.SetSelected(text)
		})
		if s.ToolTipForOption != nil {
			row.SetToolTip(s.ToolTipForOption(text))
		}
		rows[i] = row
		objects[i] = row
	}

	box := container.NewVBox(objects...)
	scroll := container.NewVScroll(box)
	boxMin := box.MinSize()

	th := s.Theme()
	tv := fyne.CurrentApp().Settings().ThemeVariant()
	pad := th.Size(theme.SizeNamePadding)
	border := th.Size(theme.SizeNameInputBorder)
	const contentPad float32 = 2

	var rowH float32
	if len(rows) > 0 {
		rowH = rows[0].MinSize().Height
	}
	step := rowH + pad

	abs := fyne.CurrentApp().Driver().AbsolutePositionForObject(s)
	belowY := abs.Y + s.Size().Height - border
	spaceBelow := cv.Size().Height - belowY - pad - 2*contentPad
	spaceAbove := abs.Y + border - pad - 2*contentPad

	openAbove := boxMin.Height > spaceBelow && spaceAbove > spaceBelow
	avail := spaceBelow
	if openAbove {
		avail = spaceAbove
	}

	height := boxMin.Height
	if avail > 0 && height > avail {
		if step > 0 {
			n := int((avail + pad) / step)
			if n < 1 {
				n = 1
			}
			height = float32(n)*step - pad
		} else {
			height = avail
		}
	}

	panel := canvas.NewRectangle(th.Color(theme.ColorNameOverlayBackground, tv))
	panel.StrokeColor = th.Color(theme.ColorNameInputBorder, tv)
	panel.StrokeWidth = th.Size(theme.SizeNameInputBorder)
	panel.CornerRadius = th.Size(theme.SizeNameSelectionRadius)

	popHeight := height + 2*contentPad
	popWidth := fyne.Max(s.Size().Width, boxMin.Width+2*contentPad)

	pos := fyne.NewPos(abs.X+(s.Size().Width-popWidth)/2, belowY)
	if openAbove {
		pos.Y = abs.Y + border - popHeight
	}

	if scale := cv.Scale(); scale > 0 {
		snap := func(v float32) float32 { return float32(math.Round(float64(v*scale))) / scale }
		pos.X, pos.Y = snap(pos.X), snap(pos.Y)
		popWidth, popHeight = snap(popWidth), snap(popHeight)
	}

	content := newSelectPopup(rows, scroll, selectedIdx, s.hidePopup)
	padded := container.New(layout.NewCustomPaddedLayout(contentPad, contentPad, contentPad, contentPad), content)
	pop := newPopOverlay(container.NewStack(panel, padded), cv, s.hidePopup)

	pop.showAt(pos, fyne.NewSize(popWidth, popHeight))
	cv.Focus(content)

	if selectedIdx >= 0 && step > 0 {
		target := float32(selectedIdx)*step - (height-rowH)/2
		if target > 0 {
			rowsAbove := int(target/step + 0.5)
			scroll.ScrollToOffset(fyne.NewPos(0, float32(rowsAbove)*step))
		}
	}

	s.popup = pop
}

func (s *Widget) hidePopup() {
	if s.popup == nil {
		return
	}
	s.popup.Hide()
	s.popup = nil
}
