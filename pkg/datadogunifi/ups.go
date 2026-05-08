package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchUPSDevice generates UPS device selector datapoints for Datadog.
func (u *DatadogUnifi) batchUPSDevice(r report, ups *unifi.UPSDeviceSelector) {
	if ups == nil {
		return
	}

	metricName := metricNamespace("ups_device")

	tags := cleanTags(map[string]string{
		"site_name": ups.SiteName,
		"source":    ups.SourceName,
		"id":        ups.ID,
		"mac":       ups.MAC,
		"label":     ups.Label,
	})

	// Emit a presence gauge (1.0 = device exists in the adoption list).
	_ = r.reportGauge(metricName("present"), 1.0, tagMapToTags(tags))
}

// batchWANStatus generates WAN status datapoints for Datadog.
func (u *DatadogUnifi) batchWANStatus(r report, ws *unifi.WANStatus) {
	if ws == nil {
		return
	}

	metricName := metricNamespace("wan_status")

	for i := range ws.WANInterfaces {
		iface := &ws.WANInterfaces[i]

		active := 0.0
		if iface.State == "ACTIVE" {
			active = 1.0
		}

		tags := cleanTags(map[string]string{
			"site_name":        ws.SiteName,
			"wan_name":         iface.Name,
			"wan_networkgroup": iface.WANNetworkgroup,
			"state":            iface.State,
		})

		_ = r.reportGauge(metricName("active"), active, tagMapToTags(tags))
	}
}
