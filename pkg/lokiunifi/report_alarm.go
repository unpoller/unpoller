package lokiunifi

import (
	"encoding/json"
	"strconv"

	"github.com/unpoller/unifi/v5"
)

const typeAlarm = "Alarm"

// Alarm stores a structured Alarm for batch sending to Loki.
// Logs the raw JSON for parsing with Loki's `| json` pipeline.
func (r *Report) Alarm(event *unifi.Alarm, logs *Logs) {
	if event.Datetime.Before(r.Oldest) {
		return
	}

	r.Counts[typeAlarm]++ // increase counter and append new log line.

	// Marshal event to JSON for the log line
	msg, err := json.Marshal(event)
	if err != nil {
		msg = []byte(event.Msg)
	}

	logs.Streams = append(logs.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime.UnixNano(), 10), string(msg)}},
		Labels: CleanLabels(map[string]string{
			"application": "unifi_alarm",
			"source":      event.SourceName,
			"site_name":   event.SiteName,
		}),
	})
}
