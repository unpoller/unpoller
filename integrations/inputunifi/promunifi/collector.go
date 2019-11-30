// Package promunifi provides the bridge between unifi metrics and prometheus.
package promunifi

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/davidnewhall/unifi-poller/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

// channel buffer, fits at least one batch.
const buffer = 50

// UnifiCollectorCnfg defines the data needed to collect and report UniFi Metrics.
type UnifiCollectorCnfg struct {
	// If non-empty, each of the collected metrics is prefixed by the
	// provided string and an underscore ("_").
	Namespace string
	// If true, any error encountered during collection is reported as an
	// invalid metric (see NewInvalidMetric). Otherwise, errors are ignored
	// and the collected metrics will be incomplete. Possibly, no metrics
	// will be collected at all.
	ReportErrors bool
	// This function is passed to the Collect() method. The Collect method runs
	// this function to retrieve the latest UniFi measurements and export them.
	CollectFn func() (*metrics.Metrics, error)
	// Provide a logger function if you want to run a routine *after* prometheus checks in.
	LoggingFn func(*Report)
}

type unifiCollector struct {
	Config UnifiCollectorCnfg
	Client *uclient
	Device *unifiDevice
	UAP    *uap
	USG    *usg
	USW    *usw
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
	Zeros   int
	Descs   int
	Metrics *metrics.Metrics
	Elapsed time.Duration
	Start   time.Time
	ch      chan []*metricExports
	wg      sync.WaitGroup
	cf      UnifiCollectorCnfg
}

// NewUnifiCollector returns a prometheus collector that will export any available
// UniFi metrics. You must provide a collection function in the opts.
func NewUnifiCollector(opts UnifiCollectorCnfg) prometheus.Collector {
	if opts.CollectFn == nil {
		panic("nil collector function")
	}
	if opts.Namespace = strings.Trim(opts.Namespace, "_") + "_"; opts.Namespace == "_" {
		opts.Namespace = ""
	}
	return &unifiCollector{
		Config: opts,
		Client: descClient(opts.Namespace + "client_"),
		Device: descDevice(opts.Namespace + "device_"), // stats for all device types.
		UAP:    descUAP(opts.Namespace + "device_"),
		USG:    descUSG(opts.Namespace + "device_"),
		USW:    descUSW(opts.Namespace + "device_"),
		Site:   descSite(opts.Namespace + "site_"),
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
	describe(u.Device)
	describe(u.UAP)
	describe(u.USG)
	describe(u.USW)
	describe(u.Site)
}

// Collect satisfies the prometheus Collector. This runs the input method to get
// the current metrics (from another package) then exports them for prometheus.
func (u *unifiCollector) Collect(ch chan<- prometheus.Metric) {
	r := &Report{
		cf:    u.Config,
		Start: time.Now(),
		ch:    make(chan []*metricExports, buffer),
	}
	defer func() {
		r.wg.Wait()
		close(r.ch)
	}()

	var err error
	if r.Metrics, err = u.Config.CollectFn(); err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewInvalidDesc(fmt.Errorf("metric fetch failed")), err)
		return
	}

	go u.exportMetrics(r, ch)
	// in loops.go.
	u.loopClients(r)
	u.loopSites(r)
	u.loopUAPs(r)
	u.loopUSWs(r)
	u.loopUSGs(r)
	u.loopUDMs(r)
}

// This is closely tied to the method above with a sync.WaitGroup.
// This method runs in a go routine and exits when the channel closes.
func (u *unifiCollector) exportMetrics(r report, ch chan<- prometheus.Metric) {
	descs := make(map[*prometheus.Desc]bool) // used as a counter
	defer r.report(descs)
	for newMetrics := range r.channel() {
		for _, m := range newMetrics {
			descs[m.Desc] = true
			switch v := m.Value.(type) {
			case unifi.FlexInt:
				ch <- r.export(m, v.Val)
			case float64:
				ch <- r.export(m, v)
			case int64:
				ch <- r.export(m, float64(v))
			case int:
				ch <- r.export(m, float64(v))
			default:
				r.error(ch, m.Desc, m.Value)
			}
		}
		r.done()
	}
}
