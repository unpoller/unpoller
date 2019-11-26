package promunifi

import (
	"fmt"
	"reflect"
	"sync"
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
	Start   time.Time
	ch      chan []*metricExports
	wg      sync.WaitGroup
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

// Collect satisfies the prometheus Collector. This runs the input method to get
// the current metrics (from another package) then exports them for prometheus.
func (u *unifiCollector) Collect(ch chan<- prometheus.Metric) {
	var err error
	r := &Report{Start: time.Now(), ch: make(chan []*metricExports)}
	defer func() {
		r.wg.Wait()
		close(r.ch)
	}()

	if r.Metrics, err = u.Config.CollectFn(); err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewInvalidDesc(fmt.Errorf("metric fetch failed")), err)
		return
	}

	go u.exportMetrics(ch, r)

	r.wg.Add(len(r.Metrics.Clients) + len(r.Metrics.Sites))
	go u.exportClients(r.Metrics.Clients, r.ch)
	go u.exportSites(r.Metrics.Sites, r.ch)

	if r.Metrics.Devices == nil {
		return
	}

	r.wg.Add(len(r.Metrics.Devices.UAPs) + len(r.Metrics.Devices.USGs) + len(r.Metrics.Devices.USWs) + len(r.Metrics.Devices.UDMs))
	go u.exportUAPs(r.Metrics.Devices.UAPs, r.ch)
	go u.exportUSGs(r.Metrics.Devices.USGs, r.ch)
	go u.exportUSWs(r.Metrics.Devices.USWs, r.ch)
	go u.exportUDMs(r.Metrics.Devices.UDMs, r.ch)
}

// This is closely tied to the method above with a sync.WaitGroup.
// This method runs in a go routine and exits when the channel closes.
func (u *unifiCollector) exportMetrics(ch chan<- prometheus.Metric, r *Report) {
	descs := make(map[*prometheus.Desc]bool) // used as a counter
	for newMetrics := range r.ch {
		for _, m := range newMetrics {
			r.Total++
			descs[m.Desc] = true

			switch v := m.Value.(type) {
			case float64:
				ch <- prometheus.MustNewConstMetric(m.Desc, m.ValueType, v, m.Labels...)
			case int64:
				ch <- prometheus.MustNewConstMetric(m.Desc, m.ValueType, float64(v), m.Labels...)
			case int:
				ch <- prometheus.MustNewConstMetric(m.Desc, m.ValueType, float64(v), m.Labels...)
			case unifi.FlexInt:
				ch <- prometheus.MustNewConstMetric(m.Desc, m.ValueType, v.Val, m.Labels...)

			default:
				r.Errors++
				if u.Config.ReportErrors {
					ch <- prometheus.NewInvalidMetric(m.Desc, fmt.Errorf("not a number: %v", m.Value))
				}
			}
		}
		r.wg.Done()
	}

	if u.Config.LoggingFn == nil {
		return
	}
	r.Descs, r.Elapsed = len(descs), time.Since(r.Start)
	u.Config.LoggingFn(r)
}
