package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchPortAnomaly generates a port anomaly datapoint for InfluxDB.
func (u *InfluxUnifi) batchPortAnomaly(r report, a *unifi.PortAnomaly) {
	if a == nil {
		return
	}

	tags := map[string]string{
		"site_name":    a.SiteName,
		"source":       a.SourceName,
		"device_mac":   a.DeviceMAC,
		"port_idx":     a.PortIdx.Txt,
		"anomaly_type": a.AnomalyType,
	}

	fields := map[string]any{
		"count":     a.Count.Val,
		"last_seen": a.LastSeen.Val,
	}

	r.send(&metric{Table: "port_anomaly", Tags: tags, Fields: fields})
}
