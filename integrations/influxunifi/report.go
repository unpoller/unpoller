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
	Events  *poller.Events
	Errors  []error
	Total   int
	Fields  int
	USG     int // Total count of USG devices.
	USW     int // Total count of USW devices.
	UAP     int // Total count of UAP devices.
	UDM     int // Total count of UDM devices.
	Eve     int // Total count of Events.
	IDS     int // Total count of IDS/IPS Events.
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
	addUDM()
	addUSG()
	addUAP()
	addUSW()
	addEvent()
	addIDS()
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

func (r *Report) addUSW() {
	r.USW++
}
func (r *Report) addUAP() {
	r.UAP++
}
func (r *Report) addUSG() {
	r.USG++
}
func (r *Report) addUDM() {
	r.UDM++
}
func (r *Report) addEvent() {
	r.Eve++
}
func (r *Report) addIDS() {
	r.IDS++
}

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
