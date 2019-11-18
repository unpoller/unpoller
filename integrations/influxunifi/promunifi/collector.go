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
}

type unifiCollector struct {
	Config UnifiCollectorCnfg
	Client *client
	UAP    *uap
	USG    *usg
	USW    *usw
	UDM    *udm
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
}

// Collect satisifes the prometheus Collector. This runs the input method to get
// the current metrics (from another package) then exports them for prometheus.
func (u *unifiCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	unifiMetrics, err := u.Config.CollectFn()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewInvalidDesc(fmt.Errorf("metric fetch failed")), err)
		return
	}

	descs := make(map[*prometheus.Desc]bool) // used as a counter
	r := &Report{Metrics: unifiMetrics}
	if u.Config.LoggingFn != nil {
		defer func() {
			r.Elapsed = time.Since(start)
			r.Descs = len(descs)
			u.Config.LoggingFn(r)
		}()
	}

	export := func(metrics []*metricExports) {
		count, errors := u.export(ch, metrics)
		r.Total += count
		r.Errors += errors
		for _, d := range metrics {
			descs[d.Desc] = true
		}
	}

	for _, asset := range r.Metrics.Clients {
		export(u.exportClient(asset))
	}
	for _, asset := range r.Metrics.Sites {
		export(u.exportSite(asset))
	}

	if r.Metrics.Devices == nil {
		return
	}

	for _, asset := range r.Metrics.Devices.UAPs {
		export(u.exportUAP(asset))
	}
	for _, asset := range r.Metrics.Devices.USGs {
		export(u.exportUSG(asset))
	}
	for _, asset := range r.Metrics.Devices.USWs {
		export(u.exportUSW(asset))
	}
	for _, asset := range r.Metrics.Devices.UDMs {
		export(u.exportUDM(asset))
	}
}

func (u *unifiCollector) export(ch chan<- prometheus.Metric, exports []*metricExports) (count, errors int) {
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
		ch <- prometheus.MustNewConstMetric(e.Desc, e.ValueType, val, e.Labels...)
	}

	return
}
