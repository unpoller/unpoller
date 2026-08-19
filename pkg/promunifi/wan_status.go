package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v6"
)

type wanStatus struct {
	InterfaceState *prometheus.Desc
}

func descWANStatus(ns string) *wanStatus {
	labels := []string{"site_name", "wan_interface", "wan_networkgroup", "state"}

	return &wanStatus{
		InterfaceState: prometheus.NewDesc(ns+"wan_interface_state",
			"WAN interface state: 1=ACTIVE, 0=other. The state label carries the raw "+
				"controller value (ACTIVE, BACKUP, or DISCONNECTED) so BACKUP and DISCONNECTED "+
				"can be distinguished.",
			labels, nil),
	}
}

func wanStateValue(state string) float64 {
	if state == "ACTIVE" {
		return 1
	}

	return 0
}

func (u *promUnifi) exportWANStatus(r report, ws *unifi.WANStatus) {
	if ws == nil {
		return
	}

	for _, iface := range ws.WANInterfaces {
		labels := []string{ws.SiteName, iface.Name, iface.WANNetworkgroup, iface.State}

		r.send([]*metric{
			{u.WANStatus.InterfaceState, gauge, wanStateValue(iface.State), labels},
		})
	}
}
