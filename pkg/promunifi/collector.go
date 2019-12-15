// Package promunifi provides the bridge between unifi-poller metrics and prometheus.
package promunifi

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
	"golift.io/unifi"
)

const (
	// channel buffer, fits at least one batch.
	buffer            = 50
	defaultHTTPListen = "0.0.0.0:9130"
	// simply fewer letters.
	counter = prometheus.CounterValue
	gauge   = prometheus.GaugeValue
)

type promUnifi struct {
	*Prometheus
	Client *uclient
	Device *unifiDevice
	UAP    *uap
	USG    *usg
	USW    *usw
	Site   *site
	// This interface is passed to the Collect() method. The Collect method uses
	// this interface to retrieve the latest UniFi measurements and export them.
	Collector poller.Collect
}

// Prometheus allows the data to be nested in the config file.
type Prometheus struct {
	Config Config `json:"prometheus" toml:"prometheus" xml:"prometheus" yaml:"prometheus"`
}

// Config is the input (config file) data used to initialize this output plugin.
type Config struct {
	Disable bool `json:"disable" toml:"disable" xml:"disable" yaml:"disable"`
	// If non-empty, each of the collected metrics is prefixed by the
	// provided string and an underscore ("_").
	Namespace string `json:"namespace" toml:"namespace" xml:"namespace" yaml:"namespace"`
	// If true, any error encountered during collection is reported as an
	// invalid metric (see NewInvalidMetric). Otherwise, errors are ignored
	// and the collected metrics will be incomplete. Possibly, no metrics
	// will be collected at all.
	ReportErrors bool   `json:"report_errors" toml:"report_errors" xml:"report_errors" yaml:"report_errors"`
	HTTPListen   string `json:"http_listen" toml:"http_listen" xml:"http_listen" yaml:"http_listen"`
}

type metric struct {
	Desc      *prometheus.Desc
	ValueType prometheus.ValueType
	Value     interface{}
	Labels    []string
}

// Report accumulates counters that are printed to a log line.
type Report struct {
	Total   int             // Total count of metrics recorded.
	Errors  int             // Total count of errors recording metrics.
	Zeros   int             // Total count of metrics equal to zero.
	Metrics *poller.Metrics // Metrics collected and recorded.
	Elapsed time.Duration   // Duration elapsed collecting and exporting.
	Fetch   time.Duration   // Duration elapsed making controller requests.
	Start   time.Time       // Time collection began.
	ch      chan []*metric
	wg      sync.WaitGroup
	Config
}

func init() {
	u := &promUnifi{Prometheus: &Prometheus{}}
	poller.NewOutput(&poller.Output{
		Name:   "prometheus",
		Config: u.Prometheus,
		Method: u.Run,
	})
}

// Run creates the collectors and starts the web server up.
// Should be run in a Go routine. Returns nil if not configured.
func (u *promUnifi) Run(c poller.Collect) error {
	if u.Config.Disable {
		return nil
	}

	if u.Config.Namespace == "" {
		u.Config.Namespace = strings.Replace(poller.AppName, "-", "", -1)
	}

	if u.Config.HTTPListen == "" {
		u.Config.HTTPListen = defaultHTTPListen
	}

	name := strings.Replace(u.Config.Namespace, "-", "_", -1)

	ns := name
	if ns = strings.Trim(ns, "_") + "_"; ns == "_" {
		ns = ""
	}

	prometheus.MustRegister(version.NewCollector(name))
	prometheus.MustRegister(&promUnifi{
		Collector: c,
		Client:    descClient(ns + "client_"),
		Device:    descDevice(ns + "device_"), // stats for all device types.
		UAP:       descUAP(ns + "device_"),
		USG:       descUSG(ns + "device_"),
		USW:       descUSW(ns + "device_"),
		Site:      descSite(ns + "site_"),
	})
	c.Logf("Exporting Measurements for Prometheus at https://%s/metrics, namespace: %s", u.Config.HTTPListen, u.Config.Namespace)

	return http.ListenAndServe(u.Config.HTTPListen, nil)
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

	r := &Report{Config: u.Config, ch: make(chan []*metric, buffer), Start: time.Now()}
	defer r.close()

	if r.Metrics, err = u.Collector.Metrics(); err != nil {
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
	defer r.report(u.Collector, descs)

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
