package influxunifi

import (
	"fmt"
	"sync"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/unifi-poller/poller"
)

// Report is returned to the calling procedure after everything is processed.
type Report struct {
	Metrics *poller.Metrics
	Events  *poller.Events
	Errors  []error
	Counts  map[item]int
	Start   time.Time
	Elapsed time.Duration
	ch      chan *metric
	wg      sync.WaitGroup
	bp      influx.BatchPoints
}

// report is an internal interface that can be mocked and overridden for tests.
type report interface {
	add()
	done()
	send(m *metric)
	error(err error)
	batch(m *metric, pt *influx.Point)
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
	if len(counts) == 0 {
		r.Counts[name]++
	}

	for _, c := range counts {
		r.Counts[name] += c
	}
}

func (r *Report) error(err error) {
	if err != nil {
		r.Errors = append(r.Errors, err)
	}
}

const (
	Ttotal = item("Point")
	Tfield = item("Fields")
)

func (r *Report) batch(m *metric, p *influx.Point) {
	r.addCount(Ttotal)
	r.addCount(Tfield, len(m.Fields))
	r.bp.AddPoint(p)
}

func (r *Report) String() string {
	return fmt.Sprintf("Site: %d, Client: %d, "+
		"%s: %d, %s/%s: %d, %s: %d, %s/%s/%s/%s: %d/%d/%d/%d, "+
		"DPI Site/Client: %d/%d, %s: %d, %s: %d, Err: %d, Dur: %v",
		len(r.Metrics.Sites), len(r.Metrics.Clients),
		TUAP, r.Counts[TUAP], TUDM, TUSG, r.Counts[TUDM]+r.Counts[TUSG], TUSW, r.Counts[TUSW],
		TIDS, Tevent, Talarm, Tanomaly, r.Counts[TIDS], r.Counts[Tevent], r.Counts[Talarm], r.Counts[Tanomaly],
		len(r.Metrics.SitesDPI), len(r.Metrics.ClientsDPI), Ttotal, r.Counts[Ttotal],
		Tfield, r.Counts[Tfield], len(r.Errors), r.Elapsed.Round(time.Millisecond))
}
