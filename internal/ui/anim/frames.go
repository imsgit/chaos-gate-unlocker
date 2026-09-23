package anim

import (
	"time"

	"fyne.io/fyne/v2"
)

func Steps(n int, d time.Duration, step func(i int), done func()) *fyne.Animation {
	return &fyne.Animation{
		Duration: d,
		Curve:    fyne.AnimationLinear,
		Tick: func(p float32) {
			step(min(int(p*float32(n)), n-1))
			if p == 1 && done != nil {
				done()
			}
		},
	}
}
