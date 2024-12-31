package lokiunifi

import (
	"strconv"

	"github.com/unpoller/unifi/v5"
)

const typeAnomaly = "Anomaly"

// Anomaly stores a structured Anomaly for batch sending to Loki.
func (r *Report) Anomaly(event *unifi.Anomaly, logs *Logs) {
	if event.Datetime.Before(r.Oldest) {
		return
	}

	r.Counts[typeAnomaly]++ // increase counter and append new log line.

	logs.Streams = append(logs.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime.UnixNano(), 10), event.Anomaly}},
		Labels: CleanLabels(map[string]string{
			"application": "unifi_anomaly",
			"source":      event.SourceName,
			"site_name":   event.SiteName,
			"device_mac":  event.DeviceMAC,
		}),
	})
}
