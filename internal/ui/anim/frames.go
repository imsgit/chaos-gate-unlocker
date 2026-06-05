package anim

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
)

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
