package promunifi

import (
	"fmt"
	"time"

	"github.com/davidnewhall/unifi-poller/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// This file contains the report interface.
// This interface can be mocked and overridden for tests.

// report is an internal interface used to "process metrics"
type report interface {
	add()
	done()
	send([]*metric)
	metrics() *metrics.Metrics
	report(descs map[*prometheus.Desc]bool)
	export(m *metric, v float64) prometheus.Metric
	error(ch chan<- prometheus.Metric, d *prometheus.Desc, v interface{})
}

// satisfy gomnd
const one = 1

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

func (r *Report) metrics() *metrics.Metrics {
	return r.Metrics
}

func (r *Report) report(descs map[*prometheus.Desc]bool) {
	if r.cf.LoggingFn == nil {
		return
	}
	r.Descs = len(descs)
	r.cf.LoggingFn(r)
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
	if r.cf.ReportErrors {
		ch <- prometheus.NewInvalidMetric(d, fmt.Errorf("error: %v", v))
	}
}

// close is not part of the interface.
func (r *Report) close() {
	r.wg.Wait()
	r.Elapsed = time.Since(r.Start)
	close(r.ch)
}
