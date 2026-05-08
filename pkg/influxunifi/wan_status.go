package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchWANStatus generates InfluxDB points for WAN interface state.
func (u *InfluxUnifi) batchWANStatus(r report, ws *unifi.WANStatus) {
	if ws == nil {
		return
	}

	for _, iface := range ws.WANInterfaces {
		tags := map[string]string{
			"site_name":        ws.SiteName,
			"wan_interface":    iface.Name,
			"wan_networkgroup": iface.WANNetworkgroup,
		}

		r.send(&metric{
			Table:  "wan_status",
			Tags:   tags,
			Fields: map[string]any{"state": iface.State},
		})
	}
}
