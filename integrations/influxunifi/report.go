package influxunifi

import (
	"sync"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/unifi-poller/poller"
)

// Report is returned to the calling procedure after everything is processed.
type Report struct {
	Metrics *poller.Metrics
	Errors  []error
	Total   int
	Fields  int
	Start   time.Time
	Elapsed time.Duration
	ch      chan *metric
	wg      sync.WaitGroup
	bp      influx.BatchPoints
}

// report is an internal interface that can be mocked and overrridden for tests.
type report interface {
	add()
	done()
	send(m *metric)
	error(err error)
	batch(m *metric, pt *influx.Point)
	metrics() *poller.Metrics
}

func (r *Report) metrics() *poller.Metrics {
	return r.Metrics
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

func (r *Report) error(err error) {
	if err != nil {
		r.Errors = append(r.Errors, err)
	}
}

func (r *Report) batch(m *metric, p *influx.Point) {
	r.Total++
	r.Fields += len(m.Fields)
	r.bp.AddPoint(p)
}
