package main

import (
	"fmt"
	"log"
	"os"

	unifipoller "github.com/davidnewhall/unifi-poller/pkg/unifi-poller"
)

func main() {
	log.SetFlags(log.LstdFlags)
	if err := run(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}

func run() error {
	unifi := &unifipoller.UnifiPoller{}
	if unifi.ParseFlags(os.Args[1:]); unifi.ShowVer {
		fmt.Printf("unifi-poller v%s\n", unifipoller.Version)
		return nil // don't run anything else w/ version request.
	}
	if err := unifi.GetConfig(); err != nil {
		unifi.Flag.Usage()
		return err
	}
	return unifi.Run()
}
