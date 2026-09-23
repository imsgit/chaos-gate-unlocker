//go:build js

package main

import (
	"strconv"
	"syscall/js"
	"time"
)

const (
	contextPollMS   = 30000
	settledAfter    = time.Minute
	attemptsKey     = "cg_ctxlost"
	retryBackoff    = 5 * time.Second
	retryBackoffMax = time.Minute
)

var _ = func() bool {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return false
	}
	doc.Call("addEventListener", "webglcontextlost", js.FuncOf(func(js.Value, []js.Value) any {
		onContextLost(doc)
		return nil
	}), true)

	poll := js.FuncOf(func(js.Value, []js.Value) any {
		if contextIsLost(doc) {
			onContextLost(doc)
		}
		return nil
	})
	doc.Call("addEventListener", "visibilitychange", poll)
	js.Global().Call("addEventListener", "focus", poll)
	js.Global().Call("setInterval", poll, contextPollMS)

	time.AfterFunc(settledAfter, func() { storeAttempts(0) })
	return true
}()

var (
	glContext   js.Value
	lossHandled bool
)

func contextIsLost(doc js.Value) bool {
	if !glContext.Truthy() {
		canvases := doc.Call("getElementsByTagName", "canvas")
		for i := 0; i < canvases.Length() && !glContext.Truthy(); i++ {
			canvas := canvases.Index(i)
			gl := canvas.Call("getContext", "webgl")
			if !gl.Truthy() {
				gl = canvas.Call("getContext", "experimental-webgl")
			}
			if gl.Truthy() {
				glContext = gl
			}
		}
	}
	return glContext.Truthy() && glContext.Call("isContextLost").Bool()
}

func onContextLost(doc js.Value) {
	if lossHandled {
		return
	}
	lossHandled = true

	attempts, remembered := loadAttempts()
	if !remembered {
		attempts = 1
	}
	storeAttempts(attempts + 1)

	reload := func() {
		if attempts == 0 {
			js.Global().Get("location").Call("reload")
			return
		}
		doc.Get("body").Set("innerHTML",
			`<p style="color:gray;text-align:center;margin-top:40vh">WebGL context lost. Reloading&hellip;</p>`)
		time.AfterFunc(backoff(attempts), func() { js.Global().Get("location").Call("reload") })
	}

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

func backoff(attempts int) time.Duration {
	delay := retryBackoff
	for i := 1; i < attempts && delay < retryBackoffMax; i++ {
		delay *= 2
	}
	return min(delay, retryBackoffMax)
}

func loadAttempts() (attempts int, remembered bool) {
	defer func() {
		if recover() != nil {
			attempts, remembered = 0, false
		}
	}()
	if v := js.Global().Get("sessionStorage").Call("getItem", attemptsKey); v.Truthy() {
		attempts, _ = strconv.Atoi(v.String())
	}
	return attempts, true
}

func storeAttempts(attempts int) {
	defer func() { _ = recover() }()
	js.Global().Get("sessionStorage").Call("setItem", attemptsKey, strconv.Itoa(attempts))
}
