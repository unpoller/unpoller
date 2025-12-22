package lokiunifi

import (
	"encoding/json"
	"strconv"

	"github.com/unpoller/unifi/v5"
)

const typeProtectLog = "ProtectLog"
const typeProtectThumbnail = "ProtectThumbnail"

// ProtectLogEvent stores a structured UniFi Protect Log Entry for batch sending to Loki.
// Logs the raw JSON for parsing with Loki's `| json` pipeline.
// If the event has a thumbnail, it's sent as a separate log line.
func (r *Report) ProtectLogEvent(event *unifi.ProtectLogEntry, logs *Logs) {
	if event.Datetime().Before(r.Oldest) {
		return
	}

	r.Counts[typeProtectLog]++ // increase counter and append new log line.

	// Store thumbnail separately before marshaling (it's excluded from JSON by default now)
	thumbnailBase64 := event.ThumbnailBase64

	// Marshal event to JSON for the log line (without thumbnail to keep it small)
	event.ThumbnailBase64 = "" // Temporarily clear for marshaling
	msg, err := json.Marshal(event)
	if err != nil {
		msg = []byte(event.Msg())
	}
	event.ThumbnailBase64 = thumbnailBase64 // Restore

	// Add event log line
	logs.Streams = append(logs.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime().UnixNano(), 10), string(msg)}},
		Labels: CleanLabels(map[string]string{
			"application": "unifi_protect_log",
			"source":      event.SourceName,
			"event_type":  event.GetEventType(),
			"category":    event.GetCategory(),
			"severity":    event.GetSeverity(),
			"camera":      event.Camera,
		}),
	})

	// Add thumbnail as separate log line if present
	if thumbnailBase64 != "" {
		r.Counts[typeProtectThumbnail]++

		thumbnailJSON, _ := json.Marshal(map[string]string{
			"event_id":         event.ID,
			"thumbnail_base64": thumbnailBase64,
			"mime_type":        "image/jpeg",
		})

		// Use timestamp + 1 nanosecond to ensure ordering (thumbnail after event)
		logs.Streams = append(logs.Streams, LogStream{
			Entries: [][]string{{strconv.FormatInt(event.Datetime().UnixNano()+1, 10), string(thumbnailJSON)}},
			Labels: CleanLabels(map[string]string{
				"application": "unifi_protect_thumbnail",
				"source":      event.SourceName,
				"event_id":    event.ID,
				"camera":      event.Camera,
			}),
		})
	}
}

