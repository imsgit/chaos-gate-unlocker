//go:build js

package anim

import (
	"sync/atomic"
	"syscall/js"
)

var (
	pageHidden atomic.Bool
	onVisible  func()
)

func init() {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return
	}
	read := func() bool { return doc.Get("hidden").Bool() }
	pageHidden.Store(read())

	doc.Call("addEventListener", "visibilitychange", js.FuncOf(func(js.Value, []js.Value) any {
		h := read()
		pageHidden.Store(h)
		if !h && onVisible != nil {
			onVisible()
		}
		return nil
	}))
}

func hidden() bool { return pageHidden.Load() }

func onShown(f func()) { onVisible = f }
