package savelistitem

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Widget struct {
	widget.BaseWidget

	hoverBg  *canvas.Rectangle
	textName *canvas.Text
}

func New() fyne.CanvasObject {
	i := &Widget{
		hoverBg:  canvas.NewRectangle(color.Transparent),
		textName: canvas.NewText("", color.White),
	}
	i.textName.TextStyle = fyne.TextStyle{Bold: true}

	i.ExtendBaseWidget(i)
	return i
}

func (i *Widget) MinSize() fyne.Size {
	i.ExtendBaseWidget(i)
	return fyne.NewSize(0, 44)
}

func (i *Widget) MouseIn(*desktop.MouseEvent) {
	i.hoverBg.FillColor = i.Theme().Color(theme.ColorNameHover, fyne.CurrentApp().Settings().ThemeVariant())
	i.hoverBg.Refresh()
}

func (i *Widget) MouseMoved(*desktop.MouseEvent) {}

func (i *Widget) MouseOut() {
	i.hoverBg.FillColor = color.Transparent
	i.hoverBg.Refresh()
}

func (i *Widget) CreateRenderer() fyne.WidgetRenderer {
	i.hoverBg.CornerRadius = i.Theme().Size(theme.SizeNameSelectionRadius)

	leftPad := canvas.NewRectangle(color.Transparent)
	leftPad.SetMinSize(fyne.NewSize(i.Theme().Size(theme.SizeNamePadding)*2, 0))

	return widget.NewSimpleRenderer(
		container.NewStack(i.hoverBg,
			container.NewPadded(container.NewHBox(leftPad, container.NewCenter(i.textName)))))
}

func (i *Widget) Bind(name string) {
	i.textName.Text = strings.TrimSuffix(name, ".gksave")
	i.textName.Refresh()
}
