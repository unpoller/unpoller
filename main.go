package main

import (
	"log"

	"github.com/davidnewhall/unifi-poller/v2/poller"
	// Load input plugins!
	_ "github.com/davidnewhall/unifi-poller/v2/inputunifi"
	// Load output plugins!
	_ "github.com/davidnewhall/unifi-poller/v2/influxunifi"
	_ "github.com/davidnewhall/unifi-poller/v2/promunifi"
)

// Keep it simple.
func main() {
	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}

/* The following activates version 1 (instead of v2 or beyond) */

/*
import (
	"log"
	"github.com/davidnewhall/unifi-poller/v1/poller"
)

func main() {
	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
*/
