// Package promunifi provides the bridge between unpoller metrics and prometheus.
package promunifi

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promver "github.com/prometheus/common/version"
	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
	"github.com/unpoller/unpoller/pkg/webserver"
	"golift.io/version"
)

// PluginName is the name of this plugin.
const PluginName = "prometheus"

const (
	// channel buffer, fits at least one batch.
	defaultBuffer     = 50
	defaultHTTPListen = "0.0.0.0:9130"
	// simply fewer letters.
	counter = prometheus.CounterValue
	gauge   = prometheus.GaugeValue
)

var ErrMetricFetchFailed = fmt.Errorf("metric fetch failed")

type promUnifi struct {
	*Config `json:"prometheus" toml:"prometheus" xml:"prometheus" yaml:"prometheus"`
	Client  *uclient
	Device  *unifiDevice
	UAP     *uap
	USG     *usg
	USW     *usw
	PDU     *pdu
	Site    *site
	RogueAP *rogueap
	// This interface is passed to the Collect() method. The Collect method uses
	// this interface to retrieve the latest UniFi measurements and export them.
	Collector poller.Collect
}

var _ poller.OutputPlugin = &promUnifi{}

// Config is the input (config file) data used to initialize this output plugin.
type Config struct {
	// If non-empty, each of the collected metrics is prefixed by the
	// provided string and an underscore ("_").
	Namespace  string `json:"namespace"   toml:"namespace"   xml:"namespace"   yaml:"namespace"`
	HTTPListen string `json:"http_listen" toml:"http_listen" xml:"http_listen" yaml:"http_listen"`
	// If these are provided, the app will attempt to listen with an SSL connection.
	SSLCrtPath string `json:"ssl_cert_path" toml:"ssl_cert_path" xml:"ssl_cert_path" yaml:"ssl_cert_path"`
	SSLKeyPath string `json:"ssl_key_path"  toml:"ssl_key_path"  xml:"ssl_key_path"  yaml:"ssl_key_path"`
	// Buffer is a channel buffer.
	// Default is probably 50. Seems fast there; try 1 to see if CPU usage goes down?
	Buffer int `json:"buffer" toml:"buffer" xml:"buffer" yaml:"buffer"`
	// If true, any error encountered during collection is reported as an
	// invalid metric (see NewInvalidMetric). Otherwise, errors are ignored
	// and the collected metrics will be incomplete. Possibly, no metrics
	// will be collected at all.
	ReportErrors bool `json:"report_errors" toml:"report_errors" xml:"report_errors" yaml:"report_errors"`
	Disable      bool `json:"disable"       toml:"disable"       xml:"disable"       yaml:"disable"`
	// Save data for dead ports? ie. ports that are down or disabled.
	DeadPorts bool `json:"dead_ports" toml:"dead_ports" xml:"dead_ports" yaml:"dead_ports"`
}

type metric struct {
	Desc      *prometheus.Desc
	ValueType prometheus.ValueType
	Value     any
	Labels    []string
}

// Report accumulates counters that are printed to a log line.
type Report struct {
	*Config
	Total   int             // Total count of metrics recorded.
	Errors  int             // Total count of errors recording metrics.
	Zeros   int             // Total count of metrics equal to zero.
	USG     int             // Total count of USG devices.
	USW     int             // Total count of USW devices.
	PDU     int             // Total count of PDU devices.
	UAP     int             // Total count of UAP devices.
	UDM     int             // Total count of UDM devices.
	UXG     int             // Total count of UXG devices.
	UBB     int             // Total count of UBB devices.
	UCI     int             // Total count of UCI devices.
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

// init is how this modular code is initialized by the main app.
// This module adds itself as an output module to the poller core.
func init() { // nolint: gochecknoinits
	u := &promUnifi{Config: &Config{}}

	poller.NewOutput(&poller.Output{
		Name:         PluginName,
		Config:       u,
		OutputPlugin: u,
	})
}

func (u *promUnifi) DebugOutput() (bool, error) {
	if u == nil {
		return true, nil
	}

	if !u.Enabled() {
		return true, nil
	}

	if u.HTTPListen == "" {
		return false, fmt.Errorf("invalid listen string")
	}

	// check the port
	parts := strings.Split(u.HTTPListen, ":")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid listen address: %s (must be of the form \"IP:Port\"", u.HTTPListen)
	}

	ln, err := net.Listen("tcp", u.HTTPListen)
	if err != nil {
		return false, err
	}

	_ = ln.Close()

	return true, nil
}

func (u *promUnifi) Enabled() bool {
	if u == nil {
		return false
	}

	if u.Config == nil {
		return false
	}

	return !u.Disable
}

// Run creates the collectors and starts the web server up.
// Should be run in a Go routine. Returns nil if not configured.
func (u *promUnifi) Run(c poller.Collect) error {
	u.Collector = c
	if u.Config == nil || !u.Enabled() {
		u.LogDebugf("Prometheus config missing (or disabled), Prometheus HTTP listener disabled!")

		return nil
	}

	u.Logf("Prometheus is enabled")

	u.Namespace = strings.Trim(strings.ReplaceAll(u.Namespace, "-", "_"), "_")
	if u.Namespace == "" {
		u.Namespace = strings.ReplaceAll(poller.AppName, "-", "")
	}

	if u.HTTPListen == "" {
		u.HTTPListen = defaultHTTPListen
	}

	if u.Buffer == 0 {
		u.Buffer = defaultBuffer
	}

	u.Client = descClient(u.Namespace + "_client_")
	u.Device = descDevice(u.Namespace + "_device_") // stats for all device types.
	u.UAP = descUAP(u.Namespace + "_device_")
	u.USG = descUSG(u.Namespace + "_device_")
	u.USW = descUSW(u.Namespace + "_device_")
	u.PDU = descPDU(u.Namespace + "_device_")
	u.Site = descSite(u.Namespace + "_site_")
	u.RogueAP = descRogueAP(u.Namespace + "_rogueap_")

	mux := http.NewServeMux()
	promver.Version = version.Version
	promver.Revision = version.Revision
	promver.Branch = version.Branch

	webserver.UpdateOutput(&webserver.Output{Name: PluginName, Config: u.Config})
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
	prometheus.MustRegister(u)
	mux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer,
		promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError},
	))
	mux.HandleFunc("/scrape", u.ScrapeHandler)
	mux.HandleFunc("/", u.DefaultHandler)

	switch u.SSLKeyPath == "" && u.SSLCrtPath == "" {
	case true:
		u.Logf("Prometheus exported at http://%s/ - namespace: %s", u.HTTPListen, u.Namespace)

		return http.ListenAndServe(u.HTTPListen, mux)
	default:
		u.Logf("Prometheus exported at https://%s/ - namespace: %s", u.HTTPListen, u.Namespace)

		return http.ListenAndServeTLS(u.HTTPListen, u.SSLCrtPath, u.SSLKeyPath, mux)
	}
}

// ScrapeHandler allows prometheus to scrape a single source, instead of all sources.
func (u *promUnifi) ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	t := &target{u: u, Filter: &poller.Filter{
		Name: r.URL.Query().Get("input"),  // "unifi"
		Path: r.URL.Query().Get("target"), // url: "https://127.0.0.1:8443"
	}}

	if t.Name == "" {
		t.Name = "unifi" // the default
	}

	if pathOld := r.URL.Query().Get("path"); pathOld != "" {
		u.LogErrorf("deprecated 'path' parameter used; update your config to use 'target'")

		if t.Path == "" {
			t.Path = pathOld
		}
	}

	if roleOld := r.URL.Query().Get("role"); roleOld != "" {
		u.LogErrorf("deprecated 'role' parameter used; update your config to use 'target'")

		if t.Path == "" {
			t.Path = roleOld
		}
	}

	if t.Path == "" {
		u.LogErrorf("'target' parameter missing on scrape from %v", r.RemoteAddr)
		http.Error(w, "'target' parameter must be specified: configured OR unconfigured url", 400)

		return
	}

	registry := prometheus.NewRegistry()

	registry.MustRegister(t)
	promhttp.HandlerFor(
		registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError},
	).ServeHTTP(w, r)
}

func (u *promUnifi) DefaultHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
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
	for _, f := range []any{u.Client, u.Device, u.UAP, u.USG, u.USW, u.PDU, u.Site} {
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
		Start:  time.Now(),
	}
	defer r.close()

	r.Metrics, err = u.Collector.Metrics(filter)
	r.Fetch = time.Since(r.Start)

	if err != nil {
		r.error(ch, prometheus.NewInvalidDesc(err), ErrMetricFetchFailed)
		u.LogErrorf("metric fetch failed: %v", err)

		return
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
	defer r.report(u, descs)

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
			case bool:
				if v {
					ch <- r.export(m, 1)
				} else {
					ch <- r.export(m, 0)
				}
			default:
				r.error(ch, m.Desc, fmt.Sprintf("not a number: %v", m.Value))
			}
		}

		r.done()
	}
}

func (u *promUnifi) loopExports(r report) {
	m := r.metrics()

	for _, s := range m.RogueAPs {
		u.switchExport(r, s)
	}

	for _, s := range m.Sites {
		u.switchExport(r, s)
	}

	for _, s := range m.SitesDPI {
		u.exportSiteDPI(r, s)
	}

	for _, c := range m.Clients {
		u.switchExport(r, c)
	}

	for _, d := range m.Devices {
		u.switchExport(r, d)
	}

	appTotal := make(totalsDPImap)
	catTotal := make(totalsDPImap)

	for _, c := range m.ClientsDPI {
		u.exportClientDPI(r, c, appTotal, catTotal)
	}

	u.exportClientDPItotals(r, appTotal, catTotal)
}

func (u *promUnifi) switchExport(r report, v any) {
	switch v := v.(type) {
	case *unifi.RogueAP:
		// r.addRogueAP()
		u.exportRogueAP(r, v)
	case *unifi.UAP:
		r.addUAP()
		u.exportUAP(r, v)
	case *unifi.USW:
		r.addUSW()
		u.exportUSW(r, v)
	case *unifi.PDU:
		r.addPDU()
		u.exportPDU(r, v)
	case *unifi.USG:
		r.addUSG()
		u.exportUSG(r, v)
	case *unifi.UXG:
		r.addUXG()
		u.exportUXG(r, v)
	case *unifi.UBB:
		r.addUBB()
		u.exportUBB(r, v)
	case *unifi.UCI:
		r.addUCI()
		u.exportUCI(r, v)
	case *unifi.UDM:
		r.addUDM()
		u.exportUDM(r, v)
	case *unifi.Site:
		u.exportSite(r, v)
	case *unifi.Client:
		u.exportClient(r, v)
	default:
		u.LogErrorf("invalid type: %T", v)
	}
}
