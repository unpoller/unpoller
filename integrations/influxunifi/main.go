package main

import (
	"log"

	"github.com/davidnewhall/unifi-poller/poller"
)

// Keep it simple.
func main() {
	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
