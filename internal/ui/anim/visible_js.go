//go:build js

package anim

import "syscall/js"

func hidden() bool {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return false
	}
	if doc.Get("hidden").Bool() {
		return true
	}
	return !doc.Call("hasFocus").Bool()
}
