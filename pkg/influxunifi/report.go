package influxunifi

import (
	"fmt"
	"sync"
	"time"

	influxV2API "github.com/influxdata/influxdb-client-go/v2/api"
	influxV2Write "github.com/influxdata/influxdb-client-go/v2/api/write"
	influxV1 "github.com/influxdata/influxdb1-client/v2"
	"github.com/unpoller/unpoller/pkg/poller"
)

// Report is returned to the calling procedure after everything is processed.
type Report struct {
	UseV2   bool
	Metrics *poller.Metrics
	Events  *poller.Events
	Errors  []error
	Counts  *Counts
	Start   time.Time
	Elapsed time.Duration
	ch      chan *metric
	wg      sync.WaitGroup
	bp      influxV1.BatchPoints
	writer  influxV2API.WriteAPI
}

// Counts holds counters and has a lock to deal with routines.
type Counts struct {
	Val map[item]int
	sync.RWMutex
}

// report is an internal interface that can be mocked and overridden for tests.
type report interface {
	add()
	done()
	send(m *metric)
	error(err error)
	batchV1(m *metric, pt *influxV1.Point)
	batchV2(m *metric, pt *influxV2Write.Point)
	metrics() *poller.Metrics
	events() *poller.Events
	addCount(item, ...int)
}

func (r *Report) metrics() *poller.Metrics {
	return r.Metrics
}

func (r *Report) events() *poller.Events {
	return r.Events
}

func (r *Report) add() {
	r.wg.Add(1)
}

func (r *Report) done() {
	r.wg.Done()
}

func (r *Report) send(m *metric) {
	r.wg.Add(1)
	r.ch <- m
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

// These constants are used as names for printed/logged counters.
const (
	pointT = item("Point")
	fieldT = item("Fields")
)

func (r *Report) batchV1(m *metric, p *influxV1.Point) {
	r.addCount(pointT)
	r.addCount(fieldT, len(m.Fields))
	r.bp.AddPoint(p)
}

func (r *Report) batchV2(m *metric, p *influxV2Write.Point) {
	r.addCount(pointT)
	r.addCount(fieldT, len(m.Fields))
	r.writer.WritePoint(p)
}

func (r *Report) String() string {
	r.Counts.RLock()
	defer r.Counts.RUnlock()

	m, c := r.Metrics, r.Counts.Val

	return fmt.Sprintf("Site: %d, Client: %d, "+
		"Gateways: %d, %s: %d, %s: %d, %s/%s/%s/%s: %d/%d/%d/%d, "+
		"DPI Site/Client: %d/%d, %s: %d, %s: %d, Err: %d, Dur: %v",
		len(m.Sites), len(m.Clients),
		c[udmT]+c[usgT]+c[uxgT]+c[uciT]+c[ubbT], uapT, c[uapT], uswT, c[uswT],
		idsT, eventT, alarmT, anomalyT, c[idsT], c[eventT], c[alarmT], c[anomalyT],
		len(m.SitesDPI), len(m.ClientsDPI), pointT, c[pointT], fieldT, c[fieldT],
		len(r.Errors), r.Elapsed.Round(time.Millisecond))
}
