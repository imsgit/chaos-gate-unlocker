package ui

import (
	"bytes"
	"image"
	_ "image/png"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	xdraw "golang.org/x/image/draw"
)

const iconBaseW = 96

var (
	decodedMu sync.Mutex
	decoded   = map[fyne.Resource]image.Image{}
)

func NewIconImage(min fyne.Size) *canvas.Image {
	img := &canvas.Image{FillMode: canvas.ImageFillContain, ScaleMode: canvas.ImageScaleSmooth}
	if min.Width > 0 || min.Height > 0 {
		img.SetMinSize(min)
	}
	return img
}

func DecodeMasked(res fyne.Resource) image.Image {
	if res == nil {
		return nil
	}
	decodedMu.Lock()
	defer decodedMu.Unlock()
	if img, ok := decoded[res]; ok {
		return img
	}
	src, _, err := image.Decode(bytes.NewReader(res.Content()))
	if err != nil {
		return nil
	}
	img := ScaleDown(src, iconBaseW)
	decoded[res] = img
	return img
}

func ScaleDown(src image.Image, w int) image.Image {
	if src == nil {
		return nil
	}
	b := src.Bounds()
	if b.Dx() <= w {
		return src
	}
	h := b.Dy() * w / b.Dx()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, b, xdraw.Over, nil)
	return dst
}

func Decode(res fyne.Resource) image.Image {
	if res == nil {
		return nil
	}
	img, _, err := image.Decode(bytes.NewReader(res.Content()))
	if err != nil {
		return nil
	}
	return img
}
