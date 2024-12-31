package lokiunifi

import (
	"strconv"

	"github.com/unpoller/unifi/v5"
)

const typeIDs = "IDs"

// event stores a structured event Event for batch sending to Loki.
func (r *Report) IDs(event *unifi.IDS, logs *Logs) {
	if event.Datetime.Before(r.Oldest) {
		return
	}

	r.Counts[typeIDs]++ // increase counter and append new log line.

	logs.Streams = append(logs.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime.UnixNano(), 10), event.Msg}},
		Labels: CleanLabels(map[string]string{
			"application":      "unifi_ids",
			"source":           event.SourceName,
			"host":             event.Host,
			"site_name":        event.SiteName,
			"subsystem":        event.Subsystem,
			"category":         event.Catname.String(),
			"event_type":       event.EventType,
			"key":              event.Key,
			"app_protocol":     event.AppProto,
			"protocol":         event.Proto,
			"interface":        event.InIface,
			"src_country":      event.SrcIPCountry,
			"src_city":         event.SourceIPGeo.City,
			"src_continent":    event.SourceIPGeo.ContinentCode,
			"src_country_code": event.SourceIPGeo.CountryCode,
			"usgip":            event.USGIP,
			"action":           event.InnerAlertAction,
		}),
	})
}
