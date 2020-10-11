package datadogunifi

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/unifi-poller/poller"
)

type Report struct {
	Metrics *poller.Metrics
	Errors  []error
	Total   int
	Fields  int
	Start   time.Time
	End     time.Time
	Elapsed time.Duration

	client statsd.ClientInterface
}

type report interface {
	error(err error)
	metrics() *poller.Metrics
	reportGauge(name string, value float64, tags []string) error
	reportCount(name string, value int64, tags []string) error
	reportDistribution(name string, value float64, tags []string) error
	reportTiming(name string, value time.Duration, tags []string) error
	reportEvent(title string, message string, tags []string) error
	reportServiceCheck(name string, status statsd.ServiceCheckStatus, message string, tags []string) error
}

func (r *Report) metrics() *poller.Metrics {
	return r.Metrics
}

func (r *Report) error(err error) {
	if err != nil {
		r.Errors = append(r.Errors, err)
	}
}

func (r *Report) reportGauge(name string, value float64, tags []string) error {
	return r.client.Gauge(name, value, tags, 1.0)
}

func (r *Report) reportCount(name string, value int64, tags []string) error {
	return r.client.Count(name, value, tags, 1.0)
}

func (r *Report) reportDistribution(name string, value float64, tags []string) error {
	return r.client.Distribution(name, value, tags, 1.0)
}

func (r *Report) reportTiming(name string, value time.Duration, tags []string) error {
	return r.client.Timing(name, value, tags, 1.0)
}

func (r *Report) reportEvent(title string, message string, tags []string) error {
	return r.client.Event(&statsd.Event{
		Title:     title,
		Text:      message,
		Timestamp: time.Now(),
		Tags:      tags,
	})
}

func (r *Report) reportServiceCheck(name string, status statsd.ServiceCheckStatus, message string, tags []string) error {
	return r.client.ServiceCheck(&statsd.ServiceCheck{
		Name:      name,
		Status:    status,
		Timestamp: time.Now(),
		Message:   message,
		Tags:      tags,
	})
}
