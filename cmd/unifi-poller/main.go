package main

import (
	"log"

	unifipoller "github.com/davidnewhall/unifi-poller/pkg/unifi-poller"
)

// Keep it simple.
func main() {
	if err := unifipoller.Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
