package main

import (
	"log"

	"github.com/davidnewhall/unifi-poller/pkg/unifipoller"
)

// Keep it simple.
func main() {
	if err := unifipoller.Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
