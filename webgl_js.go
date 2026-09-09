//go:build js

package main

import (
	"strconv"
	"syscall/js"
	"time"
)

const contextLossRetry = time.Minute

var _ = func() bool {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return false
	}
	doc.Call("addEventListener", "webglcontextlost", js.FuncOf(func(js.Value, []js.Value) any {
		onContextLost(doc)
		return nil
	}), true)
	return true
}()

func onContextLost(doc js.Value) {
	now := time.Now()
	if now.Sub(stampContextLoss(now)) < contextLossRetry {
		doc.Get("body").Set("innerHTML",
			`<p style="color:gray;text-align:center;margin-top:40vh">WebGL context lost. Please reload the page.</p>`)
		return
	}

	reload := func() { js.Global().Get("location").Call("reload") }
	if !doc.Get("hidden").Bool() {
		reload()
		return
	}
	var onVisible js.Func
	onVisible = js.FuncOf(func(js.Value, []js.Value) any {
		if !doc.Get("hidden").Bool() {
			doc.Call("removeEventListener", "visibilitychange", onVisible)
			reload()
		}
		return nil
	})
	doc.Call("addEventListener", "visibilitychange", onVisible)
}

func stampContextLoss(now time.Time) (prev time.Time) {
	defer func() { _ = recover() }()
	const key = "cg_ctxlost"
	storage := js.Global().Get("sessionStorage")
	if v := storage.Call("getItem", key); v.Truthy() {
		if ms, err := strconv.ParseInt(v.String(), 10, 64); err == nil {
			prev = time.UnixMilli(ms)
		}
	}
	storage.Call("setItem", key, strconv.FormatInt(now.UnixMilli(), 10))
	return prev
}
