package promunifi

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/unifi-poller/poller"
)

// This file contains the report interface.
// This interface can be mocked and overridden for tests.

// report is an internal interface used to "process metrics"
type report interface {
	add()
	done()
	send([]*metric)
	metrics() *poller.Metrics
	report(c poller.Collect, descs map[*prometheus.Desc]bool)
	export(m *metric, v float64) prometheus.Metric
	error(ch chan<- prometheus.Metric, d *prometheus.Desc, v interface{})
}

// satisfy gomnd
const one = 1
const oneDecimalPoint = 10.0

func (r *Report) add() {
	r.wg.Add(one)
}

func (r *Report) done() {
	r.wg.Add(-one)
}

func (r *Report) send(m []*metric) {
	r.wg.Add(one)
	r.ch <- m
}

func (r *Report) metrics() *poller.Metrics {
	return r.Metrics
}

func (r *Report) report(c poller.Collect, descs map[*prometheus.Desc]bool) {
	m := r.Metrics
	if m == nil {
		return
	}

	c.Logf("UniFi Measurements Exported. Site: %d, Client: %d, "+
		"UAP: %d, USG/UDM: %d, USW: %d, Descs: %d, "+
		"Metrics: %d, Errs: %d, 0s: %d, Reqs/Total: %v / %v",
		len(m.Sites), len(m.Clients), len(m.UAPs), len(m.UDMs)+len(m.USGs), len(m.USWs),
		len(descs), r.Total, r.Errors, r.Zeros,
		r.Fetch.Round(time.Millisecond/oneDecimalPoint),
		r.Elapsed.Round(time.Millisecond/oneDecimalPoint))
}

func (r *Report) export(m *metric, v float64) prometheus.Metric {
	r.Total++

	if v == 0 {
		r.Zeros++
	}

	return prometheus.MustNewConstMetric(m.Desc, m.ValueType, v, m.Labels...)
}

func (r *Report) error(ch chan<- prometheus.Metric, d *prometheus.Desc, v interface{}) {
	r.Errors++

	if r.ReportErrors {
		ch <- prometheus.NewInvalidMetric(d, fmt.Errorf("error: %v", v))
	}
}

// close is not part of the interface.
func (r *Report) close() {
	r.wg.Wait()
	r.Elapsed = time.Since(r.Start)
	close(r.ch)
}
