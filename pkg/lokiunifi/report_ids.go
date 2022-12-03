package lokiunifi

import (
	"strconv"

	"github.com/unpoller/unifi"
)

const typeIDS = "IDS"

// event stores a structured event Event for batch sending to Loki.
func (r *Report) IDS(event *unifi.IDS, logs *Logs) {
	if event.Datetime.Before(r.Oldest) {
		return
	}

	r.Counts[typeIDS]++ // increase counter and append new log line.

	logs.Streams = append(logs.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime.UnixNano(), 10), event.Msg}},
		Labels: CleanLabels(map[string]string{
			"application":  "unifi_ids",
			"source":       event.SourceName,
			"site_name":    event.SiteName,
			"subsystem":    event.Subsystem,
			"category":     event.Catname,
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
