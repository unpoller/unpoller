package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type portanomaly struct {
	AnomalyCount    *prometheus.Desc
	AnomalyLastSeen *prometheus.Desc
}

func descPortAnomaly(ns string) *portanomaly {
	labels := []string{"site_name", "source", "device_mac", "port_idx", "anomaly_type"}

	nd := prometheus.NewDesc

	return &portanomaly{
		AnomalyCount:    nd(ns+"port_anomaly_count", "Number of anomaly events on this port", labels, nil),
		AnomalyLastSeen: nd(ns+"port_anomaly_last_seen", "Unix timestamp of the last anomaly event on this port", labels, nil),
	}
}

func (u *promUnifi) exportPortAnomalies(r report, anomalies []*unifi.PortAnomaly) {
	for _, a := range anomalies {
		labels := []string{
			a.SiteName,
			a.SourceName,
			a.DeviceMAC,
			a.PortIdx.Txt,
			a.AnomalyType,
		}

		r.send([]*metric{
			{u.PortAnomaly.AnomalyCount, gauge, a.Count.Val, labels},
			{u.PortAnomaly.AnomalyLastSeen, gauge, a.LastSeen.Val, labels},
		})
	}
}
