//go:build nowebview

package main

import (
	"log"
	"os"
	"os/signal"
)

func openWindow(title, html string) {
	log.Printf("%s would render: %s", title, html)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
