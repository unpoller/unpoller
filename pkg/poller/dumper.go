package poller

import (
	"fmt"
	"strconv"
	"strings"
)

// DumpJSONPayload prints raw json from the UniFi Controller. This is currently
// tied into the -j CLI arg, and is probably not very useful outside that context.
func (u *UnifiPoller) DumpJSONPayload() (err error) {
	u.Config.Quiet = true
	split := strings.SplitN(u.Flags.DumpJSON, " ", 2)
	filter := &Filter{Kind: split[0]}

	if split2 := strings.Split(filter.Kind, ":"); len(split2) > 1 {
		filter.Kind = split2[0]
		filter.Unit, _ = strconv.Atoi(split2[1])
	}

	if len(split) > 1 {
		filter.Path = split[1]
	}

	m, err := inputs[0].RawMetrics(filter)
	if err != nil {
		return err
	}

	fmt.Println(string(m))

	return nil
}
