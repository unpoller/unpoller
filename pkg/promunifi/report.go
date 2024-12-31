package promunifi

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unpoller/pkg/poller"
)

// This file contains the report interface.
// This interface can be mocked and overridden for tests.

// report is an internal interface used to "process metrics".
type report interface {
	done()
	send([]*metric)
	metrics() *poller.Metrics
	report(c poller.Logger, descs map[*prometheus.Desc]bool)
	export(m *metric, v float64) prometheus.Metric
	error(ch chan<- prometheus.Metric, d *prometheus.Desc, v any)
	addUDM()
	addUXG()
	addUBB()
	addUCI()
	addUSG()
	addUAP()
	addUSW()
	addPDU()
}

// Satisfy gomnd.
const oneDecimalPoint = 10.0

func (r *Report) done() {
	r.wg.Done()
}

func (r *Report) send(m []*metric) {
	r.wg.Add(1) // notlint: gomnd
	r.ch <- m
}

func (r *Report) metrics() *poller.Metrics {
	return r.Metrics
}

func (r *Report) report(c poller.Logger, descs map[*prometheus.Desc]bool) {
	m := r.Metrics

	c.Logf("UniFi Measurements Exported. Site: %d, Client: %d, "+
		"UAP: %d, USG/UDM: %d, USW: %d, DPI Site/Client: %d/%d, Desc: %d, "+
		"Metric: %d, Err: %d, 0s: %d, Req/Total: %v / %v",
		len(m.Sites), len(m.Clients), r.UAP, r.UDM+r.USG+r.UXG, r.USW, len(m.SitesDPI),
		len(m.ClientsDPI), len(descs), r.Total, r.Errors, r.Zeros,
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

func (r *Report) error(ch chan<- prometheus.Metric, d *prometheus.Desc, v any) {
	r.Errors++

	if r.ReportErrors {
		ch <- prometheus.NewInvalidMetric(d, fmt.Errorf("error: %v", v)) // nolint: goerr113
	}
}

func (r *Report) addUSW() {
	r.USW++
}

func (r *Report) addPDU() {
	r.PDU++
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

func (r *Report) addUXG() {
	r.UXG++
}

func (r *Report) addUBB() {
	r.UCI++
}

func (r *Report) addUCI() {
	r.UCI++
}

// close is not part of the interface.
func (r *Report) close() {
	r.wg.Wait()
	r.Elapsed = time.Since(r.Start)
	close(r.ch)
}
