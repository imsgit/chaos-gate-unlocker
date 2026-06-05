package dropdown

import (
	"chaos-gate-unlocker/internal/ui/anim"
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type flatPopupTheme struct{ fyne.Theme }

func (t flatPopupTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameShadow, theme.ColorNameOverlayBackground:
		return color.Transparent
	}
	return t.Theme.Color(n, v)
}

func (t flatPopupTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameInnerPadding {
		return 0
	}
	return t.Theme.Size(n)
}

const selectTextNudge float32 = 3

type selectPadTheme struct{ fyne.Theme }

func (t selectPadTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameInnerPadding {
		return t.Theme.Size(n) + selectTextNudge
	}
	return t.Theme.Size(n)
}

func (t selectPadTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	c := t.Theme.Color(n, v)
	switch n {
	case theme.ColorNameButton, theme.ColorNameInputBackground, theme.ColorNameHover:
		nc := color.NRGBAModel.Convert(c).(color.NRGBA)
		nc.A = uint8(float64(nc.A) * 0.7)
		return nc
	}
	return c
}

type Widget struct {
	widget.Select
	tooltip.WidgetExtend

	OnBeforeShowPopup func()

	ToolTipForOption func(option string) string
	IconForOption    func(option string) fyne.Resource

	popup *widget.PopUp
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

	pad := s.Theme().Size(theme.SizeNamePadding)
	border := s.Theme().Size(theme.SizeNameInputBorder)
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

	th := fyne.CurrentApp().Settings().Theme()
	tv := fyne.CurrentApp().Settings().ThemeVariant()
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
	pop := widget.NewPopUp(container.NewStack(panel, padded), cv)
	tooltip.AddPopUpToolTipLayer(pop)
	container.NewThemeOverride(pop, flatPopupTheme{th})
	container.NewThemeOverride(scroll, th)

	pop.Resize(fyne.NewSize(popWidth, 1))
	if openAbove {
		pop.ShowAtPosition(fyne.NewPos(pos.X, pos.Y+popHeight-1))
	} else {
		pop.ShowAtPosition(pos)
	}
	cv.Focus(content)

	if selectedIdx >= 0 && step > 0 {
		target := float32(selectedIdx)*step - (height-rowH)/2
		if target > 0 {
			rowsAbove := int(target/step + 0.5)
			scroll.ScrollToOffset(fyne.NewPos(0, float32(rowsAbove)*step))
		}
	}

	s.popup = pop

	const growSteps = 12
	anim.Frames(growSteps, 16*time.Millisecond, nil, func(i int) {
		if s.popup != pop {
			return
		}
		t := float32(i) / growSteps
		h := popHeight * (1 - (1-t)*(1-t))
		if openAbove {
			pop.Move(fyne.NewPos(pos.X, pos.Y+popHeight-h))
		}
		pop.Resize(fyne.NewSize(popWidth, h))
	})
}

func (s *Widget) hidePopup() {
	if s.popup == nil {
		return
	}
	tooltip.DestroyPopUpToolTipLayer(s.popup)
	s.popup.Hide()
	s.popup = nil
}
