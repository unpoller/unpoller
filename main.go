package main

import (
	"log"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
)

// Keep it simple.
func main() {
	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
