// Package influx provides the methods to turn UniFi measurements into influx
// data-points with appropriate tags and fields.
package influxunifi

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/davidnewhall/unifi-poller/pkg/metrics"
	influx "github.com/influxdata/influxdb1-client/v2"
)

// Config defines the data needed to store metrics in InfluxDB
type Config struct {
	Database string
	URL      string
	User     string
	Pass     string
	BadSSL   bool
}

// InfluxUnifi is returned by New() after you provide a Config.
type InfluxUnifi struct {
	cf     *Config
	influx influx.Client
}

type metric struct {
	Table  string
	Tags   map[string]string
	Fields map[string]interface{}
}

// New returns an InfluxDB interface.
func New(c *Config) (*InfluxUnifi, error) {
	i, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:      c.URL,
		Username:  c.User,
		Password:  c.Pass,
		TLSConfig: &tls.Config{InsecureSkipVerify: c.BadSSL},
	})
	return &InfluxUnifi{cf: c, influx: i}, err
}

// ReportMetrics batches all device and client data into influxdb data points.
// Call this after you've collected all the data you care about.
// Returns an error if influxdb calls fail, otherwise returns a report.
func (u *InfluxUnifi) ReportMetrics(m *metrics.Metrics) (*Report, error) {
	r := &Report{Metrics: m, ch: make(chan *metric), Start: time.Now()}
	defer close(r.ch)
	// Make a new Influx Points Batcher.
	var err error
	r.bp, err = influx.NewBatchPoints(influx.BatchPointsConfig{Database: u.cf.Database})
	if err != nil {
		return nil, fmt.Errorf("influx.NewBatchPoints: %v", err)
	}

	go u.collect(r, r.ch)
	// Batch all the points.
	u.loopPoints(r)
	r.wg.Wait() // wait for all points to finish batching!

	// Send all the points.
	if err = u.influx.Write(r.bp); err != nil {
		return nil, fmt.Errorf("influxdb.Write(points): %v", err)
	}
	r.Elapsed = time.Since(r.Start)
	return r, nil
}

// collect runs in a go routine and batches all the points.
func (u *InfluxUnifi) collect(r report, ch chan *metric) {
	for m := range ch {
		pt, err := influx.NewPoint(m.Table, m.Tags, m.Fields, r.metrics().TS)
		if err != nil {
			r.error(err)
		} else {
			r.batch(m, pt)
		}
		r.done()
	}
}

// loopPoints kicks off 3 or 7 go routines to process metrics and send them
// to the collect routine through the metric channel.
func (u *InfluxUnifi) loopPoints(r report) {
	m := r.metrics()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.Sites {
			u.batchSite(r, s)
		}
	}()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.Clients {
			u.batchClient(r, s)
		}
	}()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.IDSList {
			u.batchIDS(r, s)
		}
	}()
	if m.Devices == nil {
		return
	}

	r.add()
	go func() {
		defer r.done()
		for _, s := range m.UAPs {
			u.batchUAP(r, s)
		}
	}()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.USGs {
			u.batchUSG(r, s)
		}
	}()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.USWs {
			u.batchUSW(r, s)
		}
	}()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.UDMs {
			u.batchUDM(r, s)
		}
	}()
}
