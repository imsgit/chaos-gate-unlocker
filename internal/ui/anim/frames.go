package anim

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
)

func runFrames(ctx context.Context, n int, interval time.Duration, onCancel func(), step func(i int)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			if onCancel != nil {
				fyne.DoAndWait(onCancel)
			}
			return
		case <-ticker.C:
			frame := i
			fyne.DoAndWait(func() { step(frame) })
		}
	}
}

func Frames(n int, interval time.Duration, onDone func(), step func(i int)) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		i := 0
		do := func() { step(i) }
		for i = 1; i <= n; i++ {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fyne.DoAndWait(do)
			}
		}

		if ctx.Err() == nil && onDone != nil {
			fyne.DoAndWait(onDone)
		}
		cancel()
	}()

	return cancel
}
