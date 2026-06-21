package ui

import "math"

func SnapToPixel(v, scale float32) float32 {
	if scale <= 0 {
		return v
	}
	return float32(math.Round(float64(v*scale))) / scale
}
