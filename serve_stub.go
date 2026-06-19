//go:build js || !embedwasm

package main

const browserSupported = false

func openInBrowser() error { return nil }
