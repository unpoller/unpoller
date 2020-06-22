package poller

import (
	"fmt"
	"strconv"
	"strings"
)

// PrintRawMetrics prints raw json from the UniFi Controller. This is currently
// tied into the -j CLI arg, and is probably not very useful outside that context.
func (u *UnifiPoller) PrintRawMetrics() (err error) {
	split := strings.SplitN(u.Flags.DumpJSON, " ", 2)
	filter := &Filter{Kind: split[0]}

	// Allows you to grab a controller other than 0 from config.
	if split2 := strings.Split(filter.Kind, ":"); len(split2) > 1 {
		filter.Kind = split2[0]
		filter.Unit, _ = strconv.Atoi(split2[1])
	}

	// Used with "other"
	if len(split) > 1 {
		filter.Path = split[1]
	}

	// As of now we only have one input plugin, so target that [0].
	m, err := inputs[0].RawMetrics(filter)
	fmt.Println(string(m))

	return err
}
