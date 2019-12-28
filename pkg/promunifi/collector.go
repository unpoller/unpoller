// Package promunifi provides the bridge between unifi metrics and prometheus.
package promunifi

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/davidnewhall/unifi-poller/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

// channel buffer, fits at least one batch.
const buffer = 50

// simply fewer letters.
const counter = prometheus.CounterValue
const gauge = prometheus.GaugeValue

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

type promUnifi struct {
	Config UnifiCollectorCnfg
	Client *uclient
	Device *unifiDevice
	UAP    *uap
	USG    *usg
	USW    *usw
	Site   *site
}

type metric struct {
	Desc      *prometheus.Desc
	ValueType prometheus.ValueType
	Value     interface{}
	Labels    []string
}

// Report is passed into LoggingFn to log the export metrics to stdout (outside this package).
type Report struct {
	Total   int              // Total count of metrics recorded.
	Errors  int              // Total count of errors recording metrics.
	Zeros   int              // Total count of metrics equal to zero.
	Descs   int              // Total count of unique metrics descriptions.
	Metrics *metrics.Metrics // Metrics collected and recorded.
	Elapsed time.Duration    // Duration elapsed collecting and exporting.
	Fetch   time.Duration    // Duration elapsed making controller requests.
	Start   time.Time        // Time collection began.
	ch      chan []*metric
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
	return &promUnifi{
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
func (u *promUnifi) Describe(ch chan<- *prometheus.Desc) {
	for _, f := range []interface{}{u.Client, u.Device, u.UAP, u.USG, u.USW, u.Site} {
		v := reflect.Indirect(reflect.ValueOf(f))
		// Loop each struct member and send it to the provided channel.
		for i := 0; i < v.NumField(); i++ {
			desc, ok := v.Field(i).Interface().(*prometheus.Desc)
			if ok && desc != nil {
				ch <- desc
			}
		}
	}
}

// Collect satisfies the prometheus Collector. This runs the input method to get
// the current metrics (from another package) then exports them for prometheus.
func (u *promUnifi) Collect(ch chan<- prometheus.Metric) {
	var err error
	r := &Report{cf: u.Config, ch: make(chan []*metric, buffer), Start: time.Now()}
	defer r.close()

	if r.Metrics, err = r.cf.CollectFn(); err != nil {
		r.error(ch, prometheus.NewInvalidDesc(fmt.Errorf("metric fetch failed")), err)
		return
	}
	r.Fetch = time.Since(r.Start)
	if r.Metrics.Devices == nil {
		r.Metrics.Devices = &unifi.Devices{}
	}

	// Pass Report interface into our collecting and reporting methods.
	go u.exportMetrics(r, ch, r.ch)
	u.loopExports(r)
}

// This is closely tied to the method above with a sync.WaitGroup.
// This method runs in a go routine and exits when the channel closes.
// This is where our channels connects to the prometheus channel.
func (u *promUnifi) exportMetrics(r report, ch chan<- prometheus.Metric, ourChan chan []*metric) {
	descs := make(map[*prometheus.Desc]bool) // used as a counter
	defer r.report(descs)
	for newMetrics := range ourChan {
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
				r.error(ch, m.Desc, fmt.Sprintf("not a number: %v", m.Value))
			}
		}
		r.done()
	}
}

func (u *promUnifi) loopExports(r report) {
	m := r.metrics()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.Sites {
			u.exportSite(r, s)
		}
	}()

	r.add()
	go func() {
		defer r.done()
		for _, d := range m.UAPs {
			u.exportUAP(r, d)
		}
	}()

	r.add()
	go func() {
		defer r.done()
		for _, d := range m.UDMs {
			u.exportUDM(r, d)
		}
	}()

	r.add()
	go func() {
		defer r.done()
		for _, d := range m.USGs {
			u.exportUSG(r, d)
		}
	}()

	r.add()
	go func() {
		defer r.done()
		for _, d := range m.USWs {
			u.exportUSW(r, d)
		}
	}()

	r.add()
	go func() {
		defer r.done()
		for _, c := range m.Clients {
			u.exportClient(r, c)
		}
	}()
}
