//go:build js

package main

import (
	"runtime"
	"time"
)

func init() {
	go func() {
		var m runtime.MemStats
		for {
			time.Sleep(3 * time.Second)
			runtime.ReadMemStats(&m)
			println("MEM heapAlloc(MB)", m.HeapAlloc>>20,
				"heapSys(MB)", m.HeapSys>>20,
				"sys(MB)", m.Sys>>20,
				"numGC", m.NumGC)
		}
	}()
}
