//go:build js

package anim

import (
	"sync"
	"sync/atomic"
	"syscall/js"
)

var (
	pageHidden atomic.Bool
	visMu      sync.Mutex
	onVisible  []func()
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
		if !h {
			visMu.Lock()
			fns := append([]func(){}, onVisible...)
			visMu.Unlock()
			for _, f := range fns {
				f()
			}
		}
		return nil
	}))
}

func hidden() bool { return pageHidden.Load() }

func onShown(f func()) {
	visMu.Lock()
	onVisible = append(onVisible, f)
	visMu.Unlock()
}
