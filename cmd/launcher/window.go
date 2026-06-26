//go:build !nowebview

package main

import webview "github.com/webview/webview_go"

func openWindow(title, target string) {
	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle(title)
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate(target)
	w.Run()
}
