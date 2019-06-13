package main

import (
	"fmt"
	"log"
	"os"

	unifipoller "github.com/davidnewhall/unifi-poller/pkg/unifi-poller"
)

func main() {
	u := &unifipoller.UnifiPoller{}
	if u.ParseFlags(os.Args[1:]); u.ShowVer {
		fmt.Printf("unifi-poller v%s\n", unifipoller.Version)
		return // don't run anything else w/ version request.
	}
	if err := u.GetConfig(); err != nil {
		u.Flag.Usage()
		log.Fatalf("[ERROR] config file '%v': %v", u.ConfigFile, err)
	}
	if err := u.Run(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}
