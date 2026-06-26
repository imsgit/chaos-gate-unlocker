package dropdown

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type selectPopup struct {
	widget.BaseWidget

	rows      []*selectRow
	scroll    *container.Scroll
	highlight int
	onDismiss func()
}

func newSelectPopup(rows []*selectRow, scroll *container.Scroll, cursor int, onDismiss func()) *selectPopup {
	p := &selectPopup{rows: rows, scroll: scroll, highlight: cursor, onDismiss: onDismiss}
	p.ExtendBaseWidget(p)
	return p
}

func (p *selectPopup) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.scroll)
}

func (p *selectPopup) FocusGained()   {}
func (p *selectPopup) FocusLost()     {}
func (p *selectPopup) TypedRune(rune) {}

func (p *selectPopup) Dragged(e *fyne.DragEvent) {
	if p.scroll != nil {
		p.scroll.ScrollToOffset(fyne.NewPos(p.scroll.Offset.X, p.scroll.Offset.Y-e.Dragged.DY))
	}
}

func (p *selectPopup) DragEnd() {}

func (p *selectPopup) TypedKey(e *fyne.KeyEvent) {
	switch e.Name {
	case fyne.KeyDown:
		p.move(1)
	case fyne.KeyUp:
		p.move(-1)
	case fyne.KeyReturn, fyne.KeyEnter:
		p.activate()
	case fyne.KeyEscape, fyne.KeySpace:
		if p.onDismiss != nil {
			p.onDismiss()
		}
	}
}

func (p *selectPopup) move(dir int) {
	for i := p.highlight + dir; i >= 0 && i < len(p.rows); i += dir {
		if !p.rows[i].selected {
			if p.highlight >= 0 && p.highlight < len(p.rows) {
				p.rows[p.highlight].setHighlighted(false)
			}
			p.highlight = i
			p.rows[i].setHighlighted(true)
			p.scrollTo(i)
			return
		}
	}
}

func (p *selectPopup) activate() {
	if p.highlight < 0 || p.highlight >= len(p.rows) {
		return
	}
	if row := p.rows[p.highlight]; row.onTapped != nil {
		row.onTapped()
	}
}

func (p *selectPopup) scrollTo(i int) {
	if p.scroll == nil || len(p.rows) == 0 {
		return
	}
	rowH := p.rows[i].MinSize().Height
	pad := p.Theme().Size(theme.SizeNamePadding)
	top := float32(i) * (rowH + pad)
	bottom := top + rowH

	offset := p.scroll.Offset.Y
	if viewH := p.scroll.Size().Height; top < offset {
		offset = top
	} else if bottom > offset+viewH {
		offset = bottom - viewH
	}
	p.scroll.ScrollToOffset(fyne.NewPos(0, offset))
}
