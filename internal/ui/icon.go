package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func NewIcon(min fyne.Size) *canvas.Image {
	img := &canvas.Image{FillMode: canvas.ImageFillContain, ScaleMode: canvas.ImageScaleSmooth}
	if min.Width > 0 || min.Height > 0 {
		img.SetMinSize(min)
	}
	return img
}
