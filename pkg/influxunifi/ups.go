package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchUPSDevice generates InfluxDB points for a UPS device selector entry.
// UPSDeviceSelector is a lightweight inventory record; the only meaningful
// time-series value is a presence/count gauge (1 per device per poll cycle).
func (u *InfluxUnifi) batchUPSDevice(r report, d *unifi.UPSDeviceSelector) {
	if d == nil {
		return
	}

	tags := map[string]string{
		"site_name": d.SiteName,
		"source":    d.SourceName,
		"device_id": d.ID,
		"mac":       d.MAC,
		"label":     d.Label,
	}

	r.send(&metric{
		Table:  "ups_device",
		Tags:   tags,
		Fields: map[string]any{"present": 1},
	})
}
