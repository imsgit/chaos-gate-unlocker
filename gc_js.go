//go:build js

package main

import (
	"runtime"
	"time"
)

func init() {
	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for range t.C {
			runtime.GC()
		}
	}()
}
