package lokiunifi

import (
	"strconv"

	"github.com/unpoller/unifi/v5"
)

const typeEvent = "Event"

// Event stores a structured UniFi Event for batch sending to Loki.
func (r *Report) Event(event *unifi.Event, logs *Logs) {
	if event.Datetime.Before(r.Oldest) {
		return
	}

	r.Counts[typeEvent]++ // increase counter and append new log line.

	logs.Streams = append(logs.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime.UnixNano(), 10), event.Msg}},
		Labels: CleanLabels(map[string]string{
			"application":  "unifi_event",
			"admin":        event.Admin, // username
			"host":         event.Host,
			"hostname":     event.Hostname,
			"site_name":    event.SiteName,
			"source":       event.SourceName,
			"subsystem":    event.Subsystem,
			"ap_from":      event.ApFrom,
			"ap_to":        event.ApTo,
			"ap":           event.Ap,
			"ap_name":      event.ApName,
			"gw":           event.Gw,
			"gw_name":      event.GwName,
			"sw":           event.Sw,
			"sw_name":      event.SwName,
			"category":     event.Catname.String(),
			"radio":        event.Radio,
			"radio_from":   event.RadioFrom,
			"radio_to":     event.RadioTo,
			"key":          event.Key,
			"interface":    event.InIface,
			"event_type":   event.EventType,
			"ssid":         event.SSID,
			"channel":      event.Channel.Txt,
			"channel_from": event.ChannelFrom.Txt,
			"channel_to":   event.ChannelTo.Txt,
			"usgip":        event.USGIP,
			"network":      event.Network,
			"app_protocol": event.AppProto,
			"protocol":     event.Proto,
			"action":       event.InnerAlertAction,
			"src_country":  event.SrcIPCountry,
		}),
	})
}
