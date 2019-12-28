package metrics

import (
	"time"

	"golift.io/unifi"
)

// Metrics is a type shared by the exporting and reporting packages.
type Metrics struct {
	TS time.Time
	unifi.Sites
	unifi.IDSList
	unifi.Clients
	*unifi.Devices
}
