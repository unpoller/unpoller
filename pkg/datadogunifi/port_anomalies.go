package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchPortAnomaly generates port anomaly datapoints for Datadog.
func (u *DatadogUnifi) batchPortAnomaly(r report, a *unifi.PortAnomaly) {
	if a == nil {
		return
	}

	metricName := metricNamespace("port_anomaly")

	tags := []string{
		tag("site_name", a.SiteName),
		tag("source", a.SourceName),
		tag("device_mac", a.DeviceMAC),
		tag("port_idx", a.PortIdx.Txt),
		tag("anomaly_type", a.AnomalyType),
	}

	data := map[string]float64{
		"count":     a.Count.Val,
		"last_seen": a.LastSeen.Val,
	}

	for name, value := range data {
		_ = r.reportGauge(metricName(name), value, tags)
	}
}
