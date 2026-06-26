//go:build nowebview

package main

import (
	"log"
	"os"
	"os/signal"
)

func openWindow(title, target string) {
	log.Printf("%s would render: %s", title, target)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
