package ui

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"sync"

	"fyne.io/fyne/v2"
)

var (
	decodedMu sync.Mutex
	decoded   = map[fyne.Resource]image.Image{}
)

func DecodeIcon(res fyne.Resource) image.Image {
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
	img := dropBlack(src)
	decoded[res] = img
	return img
}

// DecodeRaw decodes a resource without the dropBlack pass. Use it for images
// assigned to a canvas.Image so refreshes don't re-decode the resource.
func DecodeRaw(res fyne.Resource) image.Image {
	if res == nil {
		return nil
	}
	img, _, err := image.Decode(bytes.NewReader(res.Content()))
	if err != nil {
		return nil
	}
	return img
}

const blackKnee = 180

func dropBlack(src image.Image) image.Image {
	b := src.Bounds()
	out := image.NewNRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)
			lum := (uint32(c.R)*299 + uint32(c.G)*587 + uint32(c.B)*114) / 1000
			if lum < blackKnee {
				c.A = uint8(uint32(c.A) * lum / blackKnee)
			}
			out.SetNRGBA(x, y, c)
		}
	}
	return out
}
