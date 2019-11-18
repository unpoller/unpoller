package promunifi

import (
	"fmt"
	"reflect"
	"time"

	"github.com/davidnewhall/unifi-poller/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

// UnifiCollectorCnfg defines the data needed to collect and report UniFi Metrics.
type UnifiCollectorCnfg struct {
	// If non-empty, each of the collected metrics is prefixed by the
	// provided string and an underscore ("_").
	Namespace string
	// If true, any error encountered during collection is reported as an
	// invalid metric (see NewInvalidMetric). Otherwise, errors are ignored
	// and the collected metrics will be incomplete. (Possibly, no metrics
	// will be collected at all.)
	ReportErrors bool
	// This function is passed to the Collect() method. The Collect method runs This
	// function to retreive the latest UniFi
	CollectFn func() (*metrics.Metrics, error)
	// provide a logger function if you want to run a routine *after* prometheus checks in.
	LoggingFn func(*Report)
	// Setting this to true will enable IDS exports.
	CollectIDS bool
}

type unifiCollector struct {
	Config UnifiCollectorCnfg
	Client *client
	UAP    *uap
	USG    *usg
	USW    *usw
	UDM    *udm
	IDS    *ids
	Site   *site
}

type metricExports struct {
	Desc      *prometheus.Desc
	ValueType prometheus.ValueType
	Value     interface{}
	Labels    []string
}

// Report is passed into LoggingFn to log the export metrics to stdout (outside this package).
type Report struct {
	Total   int
	Errors  int
	Descs   int
	Metrics *metrics.Metrics
	Elapsed time.Duration
}

// NewUnifiCollector returns a prometheus collector that will export any available
// UniFi metrics. You must provide a collection function in the opts.
func NewUnifiCollector(opts UnifiCollectorCnfg) prometheus.Collector {
	if opts.CollectFn == nil {
		panic("nil collector function")
	}

	return &unifiCollector{
		Config: opts,
		Client: descClient(opts.Namespace),
		UAP:    descUAP(opts.Namespace),
		USG:    descUSG(opts.Namespace),
		USW:    descUSW(opts.Namespace),
		UDM:    descUDM(opts.Namespace),
		Site:   descSite(opts.Namespace),
		IDS:    descIDS(opts.Namespace),
	}
}

// Describe satisfies the prometheus Collector. This returns all of the
// metric descriptions that this packages produces.
func (u *unifiCollector) Describe(ch chan<- *prometheus.Desc) {
	describe := func(from interface{}) {
		v := reflect.Indirect(reflect.ValueOf(from))

		// Loop each struct member and send it to the provided channel.
		for i := 0; i < v.NumField(); i++ {
			desc, ok := v.Field(i).Interface().(*prometheus.Desc)
			if ok && desc != nil {
				ch <- desc
			}
		}
	}

	describe(u.Client)
	describe(u.UAP)
	describe(u.USG)
	describe(u.USW)
	describe(u.UDM)
	describe(u.Site)
	if u.Config.CollectIDS {
		describe(u.IDS)
	}
}

// Collect satisifes the prometheus Collector. This runs the input method to get
// the current metrics (from another package) then exports them for prometheus.
func (u *unifiCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	m, err := u.Config.CollectFn()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewInvalidDesc(fmt.Errorf("metric fetch failed")), err)
		return
	}

	descs := make(map[*prometheus.Desc]bool) // used as a counter
	r := &Report{Metrics: m}
	if u.Config.LoggingFn != nil {
		defer func() {
			r.Elapsed = time.Since(start)
			r.Descs = len(descs)
			u.Config.LoggingFn(r)
		}()
	}

	process := func(m []*metricExports, ts time.Time) {
		count, errors := u.export(ch, m, ts)
		r.Total += count
		r.Errors += errors
		for _, d := range m {
			descs[d.Desc] = true
		}
	}

	for _, asset := range m.Clients {
		process(u.exportClient(asset), m.TS)
	}
	for _, asset := range m.Sites {
		process(u.exportSite(asset), m.TS)
	}
	if u.Config.CollectIDS {
		for _, asset := range m.IDSList {
			process(u.exportIDS(asset), m.TS)
		}
	}

	if m.Devices == nil {
		return
	}

	for _, asset := range m.Devices.UAPs {
		process(u.exportUAP(asset), m.TS)
	}
	for _, asset := range m.Devices.USGs {
		process(u.exportUSG(asset), m.TS)
	}
	for _, asset := range m.Devices.USWs {
		process(u.exportUSW(asset), m.TS)
	}
	for _, asset := range m.Devices.UDMs {
		process(u.exportUDM(asset), m.TS)
	}
}

func (u *unifiCollector) export(ch chan<- prometheus.Metric, exports []*metricExports, ts time.Time) (count, errors int) {
	for _, e := range exports {
		var val float64
		switch v := e.Value.(type) {
		case float64:
			val = v
		case int64:
			val = float64(v)
		case unifi.FlexInt:
			val = v.Val
		default:
			errors++
			if u.Config.ReportErrors {
				ch <- prometheus.NewInvalidMetric(e.Desc, fmt.Errorf("not a number: %v", e.Value))
			}
			continue
		}
		count++
		ch <- prometheus.NewMetricWithTimestamp(ts, prometheus.MustNewConstMetric(e.Desc, e.ValueType, val, e.Labels...))
	}
	return
}
