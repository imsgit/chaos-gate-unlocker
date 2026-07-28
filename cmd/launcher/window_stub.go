//go:build !cgo || (!linux && !windows)

package main

import "log"

func openWindow(title, url string) {
	log.Fatalf("%s: launcher requires a cgo webview build (CGO_ENABLED=1, linux or windows); app URL: %s", title, url)
}
