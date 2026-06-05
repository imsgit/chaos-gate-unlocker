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

func DecodeIconByName(name string) image.Image {
	return DecodeIcon(GetIconByName(name))
}

const blackKnee = 160

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
