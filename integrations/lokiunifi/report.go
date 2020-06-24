package lokiunifi

import (
	"strconv"
	"strings"
	"time"

	"github.com/unifi-poller/poller"
	"github.com/unifi-poller/unifi"
)

// LogStream contains a stream of logs (like a log file).
// This app uses one stream per log entry because each log may have different labels.
type LogStream struct {
	Labels  map[string]string `json:"stream"` // "the file name"
	Entries [][]string        `json:"values"` // "the log lines"
}

// Logs is the main logs-holding structure.
type Logs struct {
	Streams []LogStream `json:"streams"` // "multiple files"
}

// Report is the temporary data generated and sent to Loki at every interval.
type Report struct {
	Eve    int // Total count of Events.
	IDS    int // Total count of IDS/IPS Events.
	Start  time.Time
	Last   time.Time
	Loki   *Client
	Events *poller.Events
	Logs
	poller.Logger
}

// SendEvents loops the event Logs, matches the interface
// type, calls the appropriate method for the data, and completes the report.
// This runs once per interval, if there was no collection error.
func (r *Report) SendEvents() error {
	for _, e := range r.Events.Logs {
		switch event := e.(type) {
		case *unifi.IDS:
			r.SaveIDS(event)
		case *unifi.Event:
			r.SaveEvent(event)
		default: // unlikely.
			r.LogErrorf("unknown event type: %T", e)
		}
	}

	return r.Loki.Send(r.Logs)
}

// SaveIDS stores a structured IDS Event for batch sending to Loki.
func (r *Report) SaveIDS(ids *unifi.IDS) {
	if ids.Datetime.Before(r.Last) {
		return
	}

	r.IDS++ // increase counter and append new log line.
	r.Streams = append(r.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(ids.Datetime.UnixNano(), 10), ids.Msg}},
		Labels: cleanLabels(map[string]string{
			"application":  "unifi_ids",
			"source":       ids.SourceName,
			"site_name":    ids.SiteName,
			"subsystem":    ids.Subsystem,
			"category":     ids.Catname,
			"event_type":   ids.EventType,
			"key":          ids.Key,
			"app_protocol": ids.AppProto,
			"protocol":     ids.Proto,
			"interface":    ids.InIface,
			"src_country":  ids.SrcIPCountry,
			"usgip":        ids.USGIP,
			"action":       ids.InnerAlertAction,
		}),
	})
}

// SaveEvent stores a structured UniFi Event for batch sending to Loki.
func (r *Report) SaveEvent(event *unifi.Event) {
	if event.Datetime.Before(r.Last) {
		return
	}

	r.Eve++ // increase counter and append new log line.
	r.Streams = append(r.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime.UnixNano(), 10), event.Msg}},
		Labels: cleanLabels(map[string]string{
			"application":  "unifi_event",
			"admin":        event.Admin, // username
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
			"category":     event.Catname,
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

// cleanLabels removes any tag that is empty.
func cleanLabels(labels map[string]string) map[string]string {
	for i := range labels {
		if strings.TrimSpace(labels[i]) == "" {
			delete(labels, i)
		}
	}

	return labels
}
