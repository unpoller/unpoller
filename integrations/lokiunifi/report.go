package lokiunifi

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/unifi-poller/poller"
	"github.com/unifi-poller/unifi"
)

// LogStream contains a stream of logs (like a log file).
// This app uses one stream per log entry because each log may have different labels.
type LogStream struct {
	Labels  map[string]string `json:"stream"` // "the file name"
	Entries [][]string        `json:"values"` // "the log lines"
}

// LogStreams is the main logs-holding structure.
type LogStreams struct {
	Streams []LogStream `json:"streams"` // "multiple files"
}

// Report is the temporary data generated and sent to Loki at every interval.
type Report struct {
	Counts map[string]int
	Start  time.Time
	Last   *time.Time
	Client *Client
	Events *poller.Events
	LogStreams
	poller.Logger
}

// ReportEvents should be easy to test.
// Reports events to Loki, updates last check time, and prints a log message.
func (r *Report) Execute(skipDur time.Duration) error {
	// Sometimes it gets stuck on old messages. This gets it past that.
	if time.Since(*r.Last) > skipDur {
		*r.Last = time.Now().Add(-skipDur)
	}

	r.ProcessEventLogs() // Compile report.

	// Send report to Loki.
	if err := r.Client.Post(r.LogStreams); err != nil {
		return errors.Wrap(err, "sending to Loki failed")
	}

	*r.Last = r.Start
	r.Logf("Events sent to Loki. Events: %d, IDS: %d, Alarm: %d, Anomalies: %d, Dur: %v",
		r.Counts[typeEvent], r.Counts[typeIDS], r.Counts[typeAlarm], r.Counts[typeAnomaly],
		time.Since(r.Start).Round(time.Millisecond))

	return nil
}

// ProcessEventLogs loops the event Logs, matches the interface
// type, calls the appropriate method for the data, and compiles the report.
// This runs once per interval, if there was no collection error.
func (r *Report) ProcessEventLogs() {
	for _, e := range r.Events.Logs {
		switch event := e.(type) {
		case *unifi.IDS:
			r.IDS(event)
		case *unifi.Event:
			r.Event(event)
		case *unifi.Alarm:
			r.Alarm(event)
		case *unifi.Anomaly:
			r.Anomaly(event)
		default: // unlikely.
			r.LogErrorf("unknown event type: %T", e)
		}
	}
}

// CleanLabels removes any tag that is empty.
func CleanLabels(labels map[string]string) map[string]string {
	for i := range labels {
		if strings.TrimSpace(labels[i]) == "" {
			delete(labels, i)
		}
	}

	return labels
}
