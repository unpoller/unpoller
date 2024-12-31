package lokiunifi

import (
	"strconv"

	"github.com/unpoller/unifi/v5"
)

const typeAlarm = "Alarm"

// Alarm stores a structured Alarm for batch sending to Loki.
func (r *Report) Alarm(event *unifi.Alarm, logs *Logs) {
	if event.Datetime.Before(r.Oldest) {
		return
	}

	r.Counts[typeAlarm]++ // increase counter and append new log line.

	logs.Streams = append(logs.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime.UnixNano(), 10), event.Msg}},
		Labels: CleanLabels(map[string]string{
			"application":  "unifi_alarm",
			"host":         event.Host,
			"source":       event.SourceName,
			"site_name":    event.SiteName,
			"subsystem":    event.Subsystem,
			"category":     event.Catname.String(),
			"event_type":   event.EventType,
			"key":          event.Key,
			"app_protocol": event.AppProto,
			"protocol":     event.Proto,
			"interface":    event.InIface,
			"src_country":  event.SrcIPCountry,
			"usgip":        event.USGIP,
			"action":       event.InnerAlertAction,
		}),
	})
}
