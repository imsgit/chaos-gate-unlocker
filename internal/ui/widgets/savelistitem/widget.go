package savelistitem

import (
	"image/color"
	"strings"

	"chaos-gate-unlocker/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Widget struct {
	widget.BaseWidget

	hoverBg    *canvas.Rectangle
	textName   *canvas.Text
	textDetail *canvas.Text
}

func New() fyne.CanvasObject {
	i := &Widget{
		hoverBg:    canvas.NewRectangle(color.Transparent),
		textName:   canvas.NewText("", color.White),
		textDetail: canvas.NewText("", ui.MutedForeground),
	}
	i.textName.TextStyle = fyne.TextStyle{Bold: true}
	i.textDetail.TextSize = 12

	i.ExtendBaseWidget(i)
	return i
}

func (i *Widget) MinSize() fyne.Size {
	i.ExtendBaseWidget(i)
	return fyne.NewSize(0, 50)
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
			container.NewPadded(container.NewHBox(leftPad,
				container.NewCenter(container.NewVBox(i.textName, i.textDetail))))))
}

func (i *Widget) Bind(title, detail string) {
	i.textName.Text = strings.ToUpper(title)
	i.textDetail.Text = detail
	i.textName.Refresh()
	i.textDetail.Refresh()
}
