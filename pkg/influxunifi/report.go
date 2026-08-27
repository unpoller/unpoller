package influxunifi

import (
	"fmt"
	"sync"
	"time"

	influxdb3 "github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	influxV2API "github.com/influxdata/influxdb-client-go/v2/api"
	influxV2Write "github.com/influxdata/influxdb-client-go/v2/api/write"
	influxV1 "github.com/influxdata/influxdb1-client/v2"
	"github.com/unpoller/unpoller/pkg/poller"
)

// InfluxVersion selects the InfluxDB client and write API.
type InfluxVersion uint8

const (
	InfluxV1 InfluxVersion = 1
	InfluxV2 InfluxVersion = 2
	InfluxV3 InfluxVersion = 3
)

// Report is returned to the calling procedure after everything is processed.
type Report struct {
	Version InfluxVersion
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
	v3      []*influxdb3.Point
	v3Mu    sync.Mutex
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
	batchV3(m *metric, pt *influxdb3.Point)
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
	bytesT = item("Bytes")
)

// calculateMetricBytes estimates the size of a metric in InfluxDB line protocol format.
// Format: measurement,tag1=value1,tag2=value2 field1=value1,field2=value2 timestamp
func calculateMetricBytes(m *metric) int {
	bytes := len(m.Table) // measurement name

	// Add tags
	for k, v := range m.Tags {
		bytes += len(k) + len(v) + 2 // tag key + tag value + '=' and ','
	}

	// Add fields
	for k, v := range m.Fields {
		bytes += len(k) + len(fmt.Sprint(v)) + 2 // field key + field value + '=' and ','
	}

	bytes += 20 // approximate size for timestamp and separators

	return bytes
}

func (r *Report) batchV1(m *metric, p *influxV1.Point) {
	r.addCount(pointT)
	r.addCount(fieldT, len(m.Fields))
	r.addCount(bytesT, calculateMetricBytes(m))
	r.bp.AddPoint(p)
}

func (r *Report) batchV2(m *metric, p *influxV2Write.Point) {
	r.addCount(pointT)
	r.addCount(fieldT, len(m.Fields))
	r.addCount(bytesT, calculateMetricBytes(m))
	r.writer.WritePoint(p)
}

func (r *Report) batchV3(m *metric, p *influxdb3.Point) {
	r.addCount(pointT)
	r.addCount(fieldT, len(m.Fields))
	r.addCount(bytesT, calculateMetricBytes(m))
	r.v3Mu.Lock()
	r.v3 = append(r.v3, p)
	r.v3Mu.Unlock()
}

func (r *Report) String() string {
	r.Counts.RLock()
	defer r.Counts.RUnlock()

	m, c := r.Metrics, r.Counts.Val

	return fmt.Sprintf("Site: %d, Client: %d, "+
		"Gateways: %d, %s: %d, %s: %d, %s/%s/%s/%s: %d/%d/%d/%d, "+
		"DPI Site/Client: %d/%d, %s: %d, %s: %d, %s: %d, Err: %d, Dur: %v",
		len(m.Sites), len(m.Clients),
		c[udmT]+c[usgT]+c[uxgT]+c[uciT]+c[ubbT], uapT, c[uapT], uswT, c[uswT],
		idsT, eventT, alarmT, anomalyT, c[idsT], c[eventT], c[alarmT], c[anomalyT],
		len(m.SitesDPI), len(m.ClientsDPI), pointT, c[pointT], fieldT, c[fieldT],
		bytesT, c[bytesT], len(r.Errors), r.Elapsed.Round(time.Millisecond))
}
