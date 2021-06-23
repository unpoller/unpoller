package main

//nolint:gci
import (
	"log"
	"os"
	"time"

	"github.com/unpoller/poller"
	// Load input plugins!
	_ "github.com/unpoller/inputunifi"
	// Load output plugins!
	_ "github.com/unpoller/influxunifi"
	_ "github.com/unpoller/lokiunifi"
	_ "github.com/unpoller/promunifi"
)

// Keep it simple.
func main() {
	// Set time zone based on TZ env variable.
	setTimeZone(os.Getenv("TZ"))

	if err := poller.New().Start(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}

func setTimeZone(tz string) {
	if tz == "" {
		return
	}

	var err error

	if time.Local, err = time.LoadLocation(tz); err != nil {
		log.Printf("[ERROR] Loading TZ Location '%s': %v\n", tz, err)
	}
}
