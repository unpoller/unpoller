package main

import (
	"log"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
	// Enable output plugins!
	_ "github.com/davidnewhall/unifi-poller/pkg/influxunifi"
	_ "github.com/davidnewhall/unifi-poller/pkg/promunifi"
)

// Keep it simple.
func main() {
	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
