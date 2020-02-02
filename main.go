package main

import (
	"log"

	"github.com/unifi-poller/poller"

	// Load input plugins!
	_ "github.com/unifi-poller/inputunifi"
	_ "github.com/unifi-poller/unifi"

	// Load output plugins!
	_ "github.com/unifi-poller/influxunifi"
	_ "github.com/unifi-poller/promunifi"
)

// Keep it simple.
func main() {
	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
