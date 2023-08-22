package main

//nolint:gci
import (
	"log"
	"os"
	"time"

	"github.com/kepath/unpollerkg/poller"
	// Load input plugins!
	_ "github.com/kepath/unpoller/pkg/inputunifi"
	// Load output plugins!
	_ "github.com/kepath/unpoller/pkg/datadogunifi"
	_ "github.com/kepath/unpoller/pkg/influxunifi"
	_ "github.com/kepath/unpoller/pkg/lokiunifi"
	_ "github.com/kepath/unpoller/pkg/promunifi"
)

// Keep it simple.
func main() {
	// Set time zone based on TZ env variable.
	setTimeZone(os.Getenv("TZ"))

	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}

func setTimeZone(timezone string) {
	if timezone == "" {
		return
	}

	var err error

	if time.Local, err = time.LoadLocation(timezone); err != nil {
		log.Printf("[ERROR] Loading TZ Location '%s': %v\n", timezone, err)
	}
}
