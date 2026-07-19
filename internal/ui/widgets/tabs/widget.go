package tabs

import (
	"image/color"

	"chaos-gate-unlocker/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const iconSize = 40

type Item struct {
	Title   string
	Icon    fyne.Resource
	Content fyne.CanvasObject
}

type Tabs struct {
	widget.BaseWidget

	OnSelected func(*Item)

	items    []*Item
	buttons  []*tabButton
	contents *fyne.Container
	selected int
	renderer *tabsRenderer
}

func New(items ...*Item) *Tabs {
	t := &Tabs{items: items, selected: -1}

	contents := make([]fyne.CanvasObject, len(items))
	for i, it := range items {
		t.buttons = append(t.buttons, newTabButton(it, func() { t.SelectIndex(i) }))
		it.Content.Hide()
		contents[i] = it.Content
	}
	t.contents = container.NewStack(contents...)

	t.ExtendBaseWidget(t)
	t.applySelect(0, false)
	return t
}

func (t *Tabs) SelectIndex(i int) { t.applySelect(i, true) }

func (t *Tabs) applySelect(i int, notify bool) {
	if i < 0 || i >= len(t.items) || i == t.selected {
		return
	}
	if t.selected >= 0 {
		t.items[t.selected].Content.Hide()
		t.buttons[t.selected].setSelected(false)
	}
	t.selected = i
	t.items[i].Content.Show()
	t.buttons[i].setSelected(true)
	t.contents.Refresh()
	if t.renderer != nil {
		t.renderer.moveSelection(t.Visible())
	}
	if notify && t.OnSelected != nil {
		t.OnSelected(t.items[i])
	}
}

func (t *Tabs) CreateRenderer() fyne.WidgetRenderer {
	objs := make([]fyne.CanvasObject, len(t.buttons))
	for i, b := range t.buttons {
		objs[i] = b
	}
	v := fyne.CurrentApp().Settings().ThemeVariant()
	sel := canvas.NewRectangle(t.Theme().Color(theme.ColorNameSelection, v))
	sel.CornerRadius = t.Theme().Size(theme.SizeNameSelectionRadius)
	sel.Hide()

	t.renderer = &tabsRenderer{
		tabs:      t,
		bar:       container.NewVBox(objs...),
		selection: sel,
	}
	return t.renderer
}

type tabsRenderer struct {
	tabs       *Tabs
	bar        *fyne.Container
	selection  *canvas.Rectangle
	anim       *fyne.Animation
	animTarget fyne.Position
}

func (r *tabsRenderer) Layout(size fyne.Size) {
	barW := r.bar.MinSize().Width
	r.bar.Resize(fyne.NewSize(barW, size.Height))
	r.bar.Move(fyne.NewPos(size.Width-barW, 0))
	r.tabs.contents.Resize(fyne.NewSize(size.Width-barW, size.Height))
	r.tabs.contents.Move(fyne.NewPos(0, 0))

	if target, ok := r.placeSelection(); ok && (r.anim == nil || target != r.animTarget) {
		r.stopAnim()
		r.selection.Move(target)
	}
}

func (r *tabsRenderer) placeSelection() (fyne.Position, bool) {
	i := r.tabs.selected
	if i < 0 || i >= len(r.tabs.buttons) {
		return fyne.Position{}, false
	}
	b := r.tabs.buttons[i]
	r.selection.Resize(b.Size())
	r.selection.Show()
	return r.bar.Position().Add(b.Position()), true
}

func (r *tabsRenderer) stopAnim() {
	if r.anim != nil {
		r.anim.Stop()
		r.anim = nil
	}
}

func (r *tabsRenderer) moveSelection(animate bool) {
	target, ok := r.placeSelection()
	if !ok {
		return
	}
	r.stopAnim()
	if !animate {
		r.selection.Move(target)
		return
	}
	r.animTarget = target
	r.anim = canvas.NewPositionAnimation(r.selection.Position(), target,
		canvas.DurationShort, func(p fyne.Position) { r.selection.Move(p) })
	r.anim.Curve = fyne.AnimationEaseInOut
	r.anim.Start()
}

func (r *tabsRenderer) MinSize() fyne.Size {
	bar := r.bar.MinSize()
	content := r.tabs.contents.MinSize()
	return fyne.NewSize(content.Width+bar.Width, fyne.Max(content.Height, bar.Height))
}

func (r *tabsRenderer) Refresh() {
	v := fyne.CurrentApp().Settings().ThemeVariant()
	r.selection.FillColor = r.tabs.Theme().Color(theme.ColorNameSelection, v)
	r.selection.CornerRadius = r.tabs.Theme().Size(theme.SizeNameSelectionRadius)
	r.selection.Refresh()
	canvas.Refresh(r.tabs)
}

func (r *tabsRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.tabs.contents, r.selection, r.bar}
}

func (r *tabsRenderer) Destroy() {
	r.stopAnim()
}

type tabButton struct {
	widget.BaseWidget

	icon     *canvas.Image
	label    *canvas.Text
	bg       *canvas.Rectangle
	onTapped func()
	selected bool
	hovered  bool
}

func newTabButton(item *Item, onTapped func()) *tabButton {
	b := &tabButton{
		icon:     ui.NewIconImage(fyne.NewSize(iconSize, iconSize)),
		label:    canvas.NewText(item.Title, color.White),
		bg:       canvas.NewRectangle(color.Transparent),
		onTapped: onTapped,
	}
	b.icon.Image = ui.DecodeMasked(item.Icon)
	b.label.Alignment = fyne.TextAlignCenter
	b.label.TextStyle = fyne.TextStyle{Bold: true}
	b.ExtendBaseWidget(b)
	b.refresh()
	return b
}

func (b *tabButton) CreateRenderer() fyne.WidgetRenderer {
	b.bg.CornerRadius = b.Theme().Size(theme.SizeNameSelectionRadius)
	pad := b.Theme().Size(theme.SizeNamePadding) * 2
	body := container.New(layout.NewCustomPaddedLayout(pad, pad, pad, pad),
		container.NewVBox(
			container.NewCenter(b.icon),
			b.label,
		))
	return widget.NewSimpleRenderer(container.NewStack(b.bg, body))
}

func (b *tabButton) setSelected(sel bool) {
	b.selected = sel
	b.refresh()
}

func (b *tabButton) refresh() {
	v := fyne.CurrentApp().Settings().ThemeVariant()
	th := b.Theme()

	if b.selected {
		b.label.Color = th.Color(theme.ColorNameForeground, v)
	} else {
		b.label.Color = th.Color(theme.ColorNamePlaceHolder, v)
	}
	if b.hovered && !b.selected {
		b.bg.FillColor = th.Color(theme.ColorNameHover, v)
	} else {
		b.bg.FillColor = color.Transparent
	}

	b.label.TextSize = th.Size(theme.SizeNameText)
	b.label.Refresh()
	b.bg.Refresh()
}

func (b *tabButton) Tapped(*fyne.PointEvent) {
	if b.onTapped != nil {
		b.onTapped()
	}
}

func (b *tabButton) MouseIn(*desktop.MouseEvent)    { b.hovered = true; b.refresh() }
func (b *tabButton) MouseMoved(*desktop.MouseEvent) {}
func (b *tabButton) MouseOut()                      { b.hovered = false; b.refresh() }
