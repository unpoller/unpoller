package main

import (
	"log"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
	// Load input plugins!
	_ "github.com/davidnewhall/unifi-poller/pkg/inputunifi"
	// Load output plugins!
	_ "github.com/davidnewhall/unifi-poller/pkg/influxunifi"
	_ "github.com/davidnewhall/unifi-poller/pkg/promunifi"
)

// Keep it simple.
func main() {
	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
