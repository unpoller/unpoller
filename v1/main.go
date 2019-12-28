package main

/* The following activates version 1 (instead of v2 or beyond) */

import (
	"log"

	"github.com/davidnewhall/unifi-poller/v1/poller"
)

func main() {
	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
