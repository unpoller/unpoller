package main

import (
	"fmt"
	"log"
	"os"

	unifipoller "github.com/davidnewhall/unifi-poller/pkg/unifi-poller"
)

func main() {
	log.SetFlags(log.LstdFlags)
	unifi := &unifipoller.UnifiPoller{}
	if unifi.ParseFlags(os.Args[1:]); unifi.ShowVer {
		fmt.Printf("unifi-poller v%s\n", unifipoller.Version)
		return // don't run anything else w/ version request.
	}
	if err := unifi.GetConfig(); err != nil {
		unifi.Flag.Usage()
		log.Fatalf("[ERROR] config file '%v': %v", unifi.ConfigFile, err)
	}
	if err := unifi.Run(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
