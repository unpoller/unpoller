package lokiunifi

import (
	"encoding/json"
	"strconv"

	"github.com/unpoller/unifi/v5"
)

const typeEvent = "Event"
const typeSystemLog = "SystemLog"

// Event stores a structured UniFi Event for batch sending to Loki.
// Logs the raw JSON for parsing with Loki's `| json` pipeline.
func (r *Report) Event(event *unifi.Event, logs *Logs) {
	if event.Datetime.Before(r.Oldest) {
		return
	}

	r.Counts[typeEvent]++ // increase counter and append new log line.

	// Marshal event to JSON for the log line
	msg, err := json.Marshal(event)
	if err != nil {
		msg = []byte(event.Msg)
	}

	logs.Streams = append(logs.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime.UnixNano(), 10), string(msg)}},
		Labels: CleanLabels(map[string]string{
			"application": "unifi_event",
			"site_name":   event.SiteName,
			"source":      event.SourceName,
		}),
	})
}

// SystemLogEvent stores a structured UniFi v2 System Log Entry for batch sending to Loki.
// Logs the raw JSON for parsing with Loki's `| json` pipeline.
func (r *Report) SystemLogEvent(event *unifi.SystemLogEntry, logs *Logs) {
	if event.Datetime().Before(r.Oldest) {
		return
	}

	r.Counts[typeSystemLog]++ // increase counter and append new log line.

	// Marshal event to JSON for the log line
	msg, err := json.Marshal(event)
	if err != nil {
		msg = []byte(event.TitleRaw)
	}

	logs.Streams = append(logs.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime().UnixNano(), 10), string(msg)}},
		Labels: CleanLabels(map[string]string{
			"application": "unifi_system_log",
			"site_name":   event.SiteName,
			"source":      event.SourceName,
			"category":    event.Category,
			"severity":    event.Severity,
		}),
	})
}
