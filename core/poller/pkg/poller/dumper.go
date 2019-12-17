package poller

import (
	"fmt"
	"strings"
)

// DumpJSONPayload prints raw json from the UniFi Controller.
// This only works with controller 0 (first one) in the config.
func (u *UnifiPoller) DumpJSONPayload() (err error) {
	u.Config.Quiet = true

	split := strings.SplitN(u.Flags.DumpJSON, " ", 2)
	filter := Filter{Type: split[0]}

	if len(split) > 1 {
		filter.Term = split[1]
	}

	m, err := inputs[0].RawMetrics(filter)
	if err != nil {
		return err
	}

	fmt.Println(string(m))

	return nil
}
