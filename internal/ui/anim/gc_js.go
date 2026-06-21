//go:build js

package anim

import "runtime"

func reclaim() { runtime.GC() }
