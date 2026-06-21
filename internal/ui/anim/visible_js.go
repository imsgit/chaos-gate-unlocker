//go:build js

package anim

import "syscall/js"

func hidden() bool {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return false
	}
	return doc.Get("hidden").Bool()
}
