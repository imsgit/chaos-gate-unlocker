package pixelsnap

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func Scale(obj fyne.CanvasObject) float32 {
	if obj != nil {
		if c := fyne.CurrentApp().Driver().CanvasForObject(obj); c != nil {
			if s := c.Scale(); s > 0 {
				return s
			}
		}
	}
	return 1
}

func snap(v, scale float32) float32 {
	return float32(math.Round(float64(v*scale))) / scale
}

func Fit(box fyne.Size, aspect, scale float32) (fyne.Size, fyne.Position) {
	if aspect <= 0 || box.Width <= 0 || box.Height <= 0 {
		return box, fyne.NewPos(0, 0)
	}

	w, h := box.Width, box.Height
	if box.Width/box.Height > aspect {
		w = box.Height * aspect
	} else {
		h = box.Width / aspect
	}

	w, h = snap(w, scale), snap(h, scale)
	x := snap((box.Width-w)/2, scale)
	y := snap((box.Height-h)/2, scale)
	return fyne.NewSize(w, h), fyne.NewPos(x, y)
}

func Image(img *canvas.Image, box fyne.Size, scaleObj fyne.CanvasObject) {
	sz, pos := Fit(box, img.Aspect(), Scale(scaleObj))
	img.Resize(sz)
	img.Move(pos)
}

func NewLayout() fyne.Layout { return imageLayout{} }

type imageLayout struct{}

func (imageLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var m fyne.Size
	for _, o := range objects {
		m = m.Max(o.MinSize())
	}
	return m
}

func (imageLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		if img, ok := o.(*canvas.Image); ok {
			Image(img, size, img)
			continue
		}
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}
