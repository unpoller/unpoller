package datadogunifi

import (
	"sync"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/unpoller/unpoller/pkg/poller"
)

// Report is a will report the current collection run data.
type Report struct {
	Metrics *poller.Metrics
	Events  *poller.Events
	Errors  []error
	Counts  *Counts
	Start   time.Time
	End     time.Time
	Elapsed time.Duration

	Collector poller.Collect

	Total  int
	Fields int

	wg sync.WaitGroup

	client statsd.ClientInterface
}

// Counts holds counters and has a lock to deal with routines.
type Counts struct {
	Val map[item]int
	sync.RWMutex
}

type report interface {
	add()
	done()
	error(err error)
	metrics() *poller.Metrics
	events() *poller.Events
	addCount(item, ...int)

	reportGauge(name string, value float64, tags []string) error
	reportCount(name string, value int64, tags []string) error
	reportDistribution(name string, value float64, tags []string) error
	reportTiming(name string, value time.Duration, tags []string) error
	reportEvent(title string, date time.Time, message string, tags []string) error
	reportInfoLog(message string, f ...any)
	reportWarnLog(message string, f ...any)
	reportServiceCheck(name string, status statsd.ServiceCheckStatus, message string, tags []string) error
}

func (r *Report) add() {
	r.wg.Add(1)
}

func (r *Report) done() {
	r.wg.Done()
}

func (r *Report) metrics() *poller.Metrics {
	return r.Metrics
}

func (r *Report) events() *poller.Events {
	return r.Events
}

/* The following methods are not thread safe. */

type item string

func (r *Report) addCount(name item, counts ...int) {
	r.Counts.Lock()
	defer r.Counts.Unlock()

	if len(counts) == 0 {
		r.Counts.Val[name]++
	}

	for _, c := range counts {
		r.Counts.Val[name] += c
	}
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

func (r *Report) reportEvent(title string, date time.Time, message string, tags []string) error {
	if date.IsZero() {
		date = time.Now()
	}

	return r.client.Event(&statsd.Event{
		Title:     title,
		Text:      message,
		Timestamp: date,
		Tags:      tags,
	})
}

func (r *Report) reportInfoLog(message string, f ...any) {
	r.Collector.Logf(message, f)
}

func (r *Report) reportWarnLog(message string, f ...any) {
	r.Collector.Logf(message, f)
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
