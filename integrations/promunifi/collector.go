// Package promunifi provides the bridge between unifi-poller metrics and prometheus.
package promunifi

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/unifi-poller/poller"
	"golift.io/unifi"
)

const (
	// channel buffer, fits at least one batch.
	defaultBuffer     = 50
	defaultHTTPListen = "0.0.0.0:9130"
	// simply fewer letters.
	counter = prometheus.CounterValue
	gauge   = prometheus.GaugeValue
)

type promUnifi struct {
	*Config `json:"prometheus" toml:"prometheus" xml:"prometheus" yaml:"prometheus"`
	Client  *uclient
	Device  *unifiDevice
	UAP     *uap
	USG     *usg
	USW     *usw
	Site    *site
	// This interface is passed to the Collect() method. The Collect method uses
	// this interface to retrieve the latest UniFi measurements and export them.
	Collector poller.Collect
}

// Config is the input (config file) data used to initialize this output plugin.
type Config struct {
	// If non-empty, each of the collected metrics is prefixed by the
	// provided string and an underscore ("_").
	Namespace  string `json:"namespace" toml:"namespace" xml:"namespace" yaml:"namespace"`
	HTTPListen string `json:"http_listen" toml:"http_listen" xml:"http_listen" yaml:"http_listen"`
	// If true, any error encountered during collection is reported as an
	// invalid metric (see NewInvalidMetric). Otherwise, errors are ignored
	// and the collected metrics will be incomplete. Possibly, no metrics
	// will be collected at all.
	ReportErrors bool `json:"report_errors" toml:"report_errors" xml:"report_errors" yaml:"report_errors"`
	Disable      bool `json:"disable" toml:"disable" xml:"disable" yaml:"disable"`
	// Buffer is a channel buffer.
	// Default is probably 50. Seems fast there; try 1 to see if CPU usage goes down?
	Buffer int `json:"buffer" toml:"buffer" xml:"buffer" yaml:"buffer"`
}

type metric struct {
	Desc      *prometheus.Desc
	ValueType prometheus.ValueType
	Value     interface{}
	Labels    []string
}

// Report accumulates counters that are printed to a log line.
type Report struct {
	*Config
	Total   int             // Total count of metrics recorded.
	Errors  int             // Total count of errors recording metrics.
	Zeros   int             // Total count of metrics equal to zero.
	Metrics *poller.Metrics // Metrics collected and recorded.
	Elapsed time.Duration   // Duration elapsed collecting and exporting.
	Fetch   time.Duration   // Duration elapsed making controller requests.
	Start   time.Time       // Time collection began.
	ch      chan []*metric
	wg      sync.WaitGroup
}

// target is used for targeted (sometimes dynamic) metrics scrapes.
type target struct {
	*poller.Filter
	u *promUnifi
}

func init() {
	u := &promUnifi{Config: &Config{}}

	poller.NewOutput(&poller.Output{
		Name:   "prometheus",
		Config: u,
		Method: u.Run,
	})
}

// Run creates the collectors and starts the web server up.
// Should be run in a Go routine. Returns nil if not configured.
func (u *promUnifi) Run(c poller.Collect) error {
	if u.Disable {
		return nil
	}

	u.Namespace = strings.Trim(strings.Replace(u.Namespace, "-", "_", -1), "_")
	if u.Namespace == "" {
		u.Namespace = strings.Replace(poller.AppName, "-", "", -1)
	}

	if u.HTTPListen == "" {
		u.HTTPListen = defaultHTTPListen
	}

	if u.Buffer == 0 {
		u.Buffer = defaultBuffer
	}

	// Later can pass this in from poller by adding a method to the interface.
	u.Collector = c
	u.Client = descClient(u.Namespace + "_client_")
	u.Device = descDevice(u.Namespace + "_device_") // stats for all device types.
	u.UAP = descUAP(u.Namespace + "_device_")
	u.USG = descUSG(u.Namespace + "_device_")
	u.USW = descUSW(u.Namespace + "_device_")
	u.Site = descSite(u.Namespace + "_site_")
	mux := http.NewServeMux()

	prometheus.MustRegister(version.NewCollector(u.Namespace))
	prometheus.MustRegister(u)
	c.Logf("Prometheus exported at https://%s/ - namespace: %s", u.HTTPListen, u.Namespace)
	mux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer,
		promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError},
	))
	mux.HandleFunc("/scrape", u.ScrapeHandler)
	mux.HandleFunc("/", u.DefaultHandler)

	return http.ListenAndServe(u.HTTPListen, mux)
}

// ScrapeHandler allows prometheus to scrape a single source, instead of all sources.
func (u *promUnifi) ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	t := &target{u: u, Filter: &poller.Filter{
		Name: r.URL.Query().Get("input"), // "unifi"
		Path: r.URL.Query().Get("path"),  // url: "https://127.0.0.1:8443"
		Role: r.URL.Query().Get("role"),  // configured role in up.conf.
	}}

	if t.Name == "" {
		u.Collector.LogErrorf("input parameter missing on scrape from %v", r.RemoteAddr)
		http.Error(w, `'input' parameter must be specified (try "unifi")`, 400)

		return
	}

	if t.Role == "" && t.Path == "" {
		u.Collector.LogErrorf("role and path parameters missing on scrape from %v", r.RemoteAddr)
		http.Error(w, "'role' OR 'path' parameter must be specified: configured role OR unconfigured url", 400)

		return
	}

	registry := prometheus.NewRegistry()

	registry.MustRegister(t)
	promhttp.HandlerFor(
		registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError},
	).ServeHTTP(w, r)
}

func (u *promUnifi) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte(poller.AppName + "\n"))
}

// Describe satisfies the prometheus Collector. This returns all of the
// metric descriptions that this packages produces.
func (t *target) Describe(ch chan<- *prometheus.Desc) {
	t.u.Describe(ch)
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

// Collect satisfies the prometheus Collector. This runs for a single controller poll.
func (t *target) Collect(ch chan<- prometheus.Metric) {
	t.u.collect(ch, t.Filter)
}

// Collect satisfies the prometheus Collector. This runs the input method to get
// the current metrics (from another package) then exports them for prometheus.
func (u *promUnifi) Collect(ch chan<- prometheus.Metric) {
	u.collect(ch, nil)
}

func (u *promUnifi) collect(ch chan<- prometheus.Metric, filter *poller.Filter) {
	var err error

	r := &Report{
		Config: u.Config,
		ch:     make(chan []*metric, u.Config.Buffer),
		Start:  time.Now()}
	defer r.close()

	ok := false

	if filter == nil {
		r.Metrics, ok, err = u.Collector.Metrics()
	} else {
		r.Metrics, ok, err = u.Collector.MetricsFrom(filter)
	}

	r.Fetch = time.Since(r.Start)

	if err != nil {
		r.error(ch, prometheus.NewInvalidDesc(err), fmt.Errorf("metric fetch failed"))
		u.Collector.LogErrorf("metric fetch failed: %v", err)

		if !ok {
			return
		}
	}

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
	r.add()
	r.add()
	r.add()
	r.add()
	r.add()
	r.add()
	r.add()

	go func() {
		defer r.done()

		for _, s := range m.Sites {
			u.exportSite(r, s)
		}
	}()

	go func() {
		defer r.done()

		for _, s := range m.SitesDPI {
			u.exportSiteDPI(r, s)
		}
	}()

	go func() {
		defer r.done()

		for _, c := range m.Clients {
			u.exportClient(r, c)
		}
	}()

	go func() {
		defer r.done()

		for _, c := range m.ClientsDPI {
			u.exportClientDPI(r, c)
		}
	}()

	go func() {
		defer r.done()

		for _, d := range m.UAPs {
			u.exportUAP(r, d)
		}
	}()

	go func() {
		defer r.done()

		for _, d := range m.UDMs {
			u.exportUDM(r, d)
		}
	}()

	go func() {
		defer r.done()

		for _, d := range m.USGs {
			u.exportUSG(r, d)
		}
	}()

	go func() {
		defer r.done()

		for _, d := range m.USWs {
			u.exportUSW(r, d)
		}
	}()
}
