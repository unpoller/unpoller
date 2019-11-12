package main

import (
	"log"

	"github.com/davidnewhall/unifi-poller/pollerunifi"
)

// Keep it simple.
func main() {
	if err := pollerunifi.Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
