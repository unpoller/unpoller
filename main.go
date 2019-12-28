package main

import (
	"log"

	"github.com/davidnewhall/unifi-poller/v2/pkg/poller"
	// Load input plugins!
	_ "github.com/davidnewhall/unifi-poller/v2/pkg/inputunifi"
	// Load output plugins!
	_ "github.com/davidnewhall/unifi-poller/v2/pkg/influxunifi"
	_ "github.com/davidnewhall/unifi-poller/v2/pkg/promunifi"
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
	"github.com/davidnewhall/unifi-poller/pkg/poller"
)

func main() {
	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
*/
