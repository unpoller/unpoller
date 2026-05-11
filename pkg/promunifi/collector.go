// Package promunifi provides the bridge between unpoller metrics and prometheus.
package promunifi

import (
	"fmt"
	"net"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promver "github.com/prometheus/common/version"
	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
	"github.com/unpoller/unpoller/pkg/webserver"
	"golang.org/x/sync/singleflight"
	"golift.io/cnfg"
	"golift.io/version"
)

// PluginName is the name of this plugin.
const PluginName = "prometheus"

const (
	// channel buffer, fits at least one batch.
	defaultBuffer     = 50
	defaultHTTPListen = "0.0.0.0:9130"
	// defaultInterval matches the typical UniFi Retry-After window; gives the
	// controller time to recover between background polls without starving
	// scrapes for fresh data.
	defaultInterval = 60 * time.Second
	minimumInterval = 15 * time.Second
	// simply fewer letters.
	counter = prometheus.CounterValue
	gauge   = prometheus.GaugeValue
)

// ErrMetricFetchFailed is reported as an invalid metric description when a
// scrape cannot obtain data from the configured collector.
var ErrMetricFetchFailed = fmt.Errorf("metric fetch failed")

type promUnifi struct {
	*Config           `json:"prometheus" toml:"prometheus" xml:"prometheus" yaml:"prometheus"`
	Client            *uclient
	Device            *unifiDevice
	UAP               *uap
	USG               *usg
	USW               *usw
	PDU               *pdu
	Site              *site
	RogueAP           *rogueap
	SpeedTest         *speedtest
	CountryTraffic    *ucountrytraffic
	DHCPLease         *dhcplease
	WAN               *wan
	Controller        *controller
	FirewallPolicy    *firewallpolicy
	Topology          *topology
	PortAnomaly       *portanomaly
	VPNMesh           *vpnmesh
	IntegrationDevice   *integrationDevice
	WANStatus           *wanStatus
	PortForward         *portForward
	SSLCertificate      *sslCertificate
	UPSDevice           *upsDevice
	WifiBroadcast       *wifiBroadcast
	FirewallZone        *firewallZone
	ACLRule             *aclRule
	VPNServer           *vpnServer
	SiteToSiteTunnel    *siteToSiteTunnel
	LAG                 *lag
	MCLAGDomain         *mclagDomain
	SwitchStack         *switchStack
	DNSPolicy           *dnsPolicy
	RADIUSProfile       *radiusProfile
	TrafficMatchingList *trafficMatchingList
	HotspotVoucher      *hotspotVoucher
	DPIApplication      *dpiApplication
	DPICategory         *dpiCategory
	PendingDevice       *pendingDevice
	Country             *country
	// controllerUp tracks per-controller poll success (1) or failure (0).
	// Reflects the most recent background poll — when /metrics is served from
	// a stale cache, controllerUp lags real-time health; pair with
	// unpoller_prometheus_cache_age_seconds for staleness signals.
	controllerUp *prometheus.GaugeVec
	// refreshFailures counts background refresh failures since process start
	// so operators can alert on failure rate independently of cache staleness.
	refreshFailures prometheus.Counter
	// cache holds the last successful metrics snapshot from the background
	// poller. Run() always initializes it; the nil-guards in cache-using
	// methods exist only for tests that exercise those methods directly
	// without invoking Run().
	cache *metricsCache
	// scrapeFlight coalesces concurrent /scrape requests targeting the same
	// controller URL so a noisy scraper can't multiply upstream load.
	scrapeFlight singleflight.Group
	// This interface is passed to the Collect() method. The Collect method uses
	// this interface to retrieve the latest UniFi measurements and export them.
	Collector poller.Collect
}

// metricsCache stores the latest background-poller snapshot. Reads and writes
// are serialized by an RWMutex so scrapes observe a consistent (metrics,
// fetchedAt, lastErr) triple while the background ticker refreshes in place.
// Failed refreshes preserve the prior good snapshot so /metrics never blanks
// out under transient 429 backoffs.
type metricsCache struct {
	mu        sync.RWMutex
	metrics   *poller.Metrics
	fetchedAt time.Time
	lastErr   error
}

func (c *metricsCache) get() (*poller.Metrics, time.Time, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.metrics, c.fetchedAt, c.lastErr
}

func (c *metricsCache) set(m *poller.Metrics, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastErr = err
	// Keep the last good snapshot on error so /metrics never blanks out
	// during transient upstream failures (e.g. 429 backoff). Also reject
	// (nil, nil) — a "successful" empty fetch would otherwise leave us with
	// nil metrics but a fresh fetchedAt, fooling the cache-age gauge.
	if err == nil && m != nil {
		c.metrics = m
		c.fetchedAt = time.Now()
	}
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
	// Interval controls how often the background poller refreshes the cached
	// metrics that /metrics scrapes are served from. Decouples Prometheus
	// scrape cadence from upstream UniFi API calls so 429 backoff loops cannot
	// stall scrapes. Defaults to defaultInterval; values below minimumInterval
	// are clamped up. Must be > 0 before use; normalizeInterval applies the
	// default and floor during Run().
	Interval cnfg.Duration `json:"interval" toml:"interval" xml:"interval" yaml:"interval"`
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
	Bytes   int             // Total count of bytes written.
	USG     int             // Total count of USG devices.
	USW     int             // Total count of USW devices.
	PDU     int             // Total count of PDU devices.
	UAP     int             // Total count of UAP devices.
	UDM     int             // Total count of UDM devices.
	UXG     int             // Total count of UXG devices.
	UBB     int             // Total count of UBB devices.
	UCI     int             // Total count of UCI devices.
	UDB     int             // Total count of UDB devices.
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

// DebugOutput validates the Prometheus output configuration: address format
// and (outside of health-check mode) bindability of the listen port.
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

	// Skip network binding check during health checks to avoid "address already in use"
	// errors when the main application is already running and bound to the port.
	if poller.IsHealthCheckMode() {
		return true, nil
	}

	ln, err := net.Listen("tcp", u.HTTPListen)
	if err != nil {
		return false, err
	}

	_ = ln.Close()

	return true, nil
}

// Enabled reports whether this output plugin is configured and active.
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

	u.normalizeInterval()

	u.Client = descClient(u.Namespace + "_client_")
	u.Device = descDevice(u.Namespace + "_device_") // stats for all device types.
	u.UAP = descUAP(u.Namespace + "_device_")
	u.USG = descUSG(u.Namespace + "_device_")
	u.USW = descUSW(u.Namespace + "_device_")
	u.PDU = descPDU(u.Namespace + "_device_")
	u.Site = descSite(u.Namespace + "_site_")
	u.RogueAP = descRogueAP(u.Namespace + "_rogueap_")
	u.SpeedTest = descSpeedTest(u.Namespace + "_speedtest_")
	u.CountryTraffic = descCountryTraffic(u.Namespace + "_countrytraffic_")
	u.DHCPLease = descDHCPLease(u.Namespace + "_")
	u.WAN = descWAN(u.Namespace + "_")
	u.Controller = descController(u.Namespace + "_")
	u.FirewallPolicy = descFirewallPolicy(u.Namespace + "_")
	u.Topology = descTopology(u.Namespace + "_")
	u.PortAnomaly = descPortAnomaly(u.Namespace + "_")
	u.VPNMesh = descVPNMesh(u.Namespace + "_")
	u.IntegrationDevice = descIntegrationDevice(u.Namespace + "_")
	u.WANStatus = descWANStatus(u.Namespace + "_")
	u.PortForward = descPortForward(u.Namespace + "_")
	u.SSLCertificate = descSSLCertificate(u.Namespace + "_")
	u.UPSDevice = descUPSDevice(u.Namespace + "_")
	u.WifiBroadcast = descWifiBroadcast(u.Namespace + "_")
	u.FirewallZone = descFirewallZone(u.Namespace + "_")
	u.ACLRule = descACLRule(u.Namespace + "_")
	u.VPNServer = descVPNServer(u.Namespace + "_")
	u.SiteToSiteTunnel = descSiteToSiteTunnel(u.Namespace + "_")
	u.LAG = descLAG(u.Namespace + "_")
	u.MCLAGDomain = descMCLAGDomain(u.Namespace + "_")
	u.SwitchStack = descSwitchStack(u.Namespace + "_")
	u.DNSPolicy = descDNSPolicy(u.Namespace + "_")
	u.RADIUSProfile = descRADIUSProfile(u.Namespace + "_")
	u.TrafficMatchingList = descTrafficMatchingList(u.Namespace + "_")
	u.HotspotVoucher = descHotspotVoucher(u.Namespace + "_")
	u.DPIApplication = descDPIApplication(u.Namespace + "_")
	u.DPICategory = descDPICategory(u.Namespace + "_")
	u.PendingDevice = descPendingDevice(u.Namespace + "_")
	u.Country = descCountry(u.Namespace + "_")
	u.controllerUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: u.Namespace + "_controller_up",
		Help: "Whether the most recent background poll of the UniFi controller succeeded (1) or failed (0). " +
			"Reflects the last poll attempt, not real-time health; pair with " +
			u.Namespace + "_prometheus_cache_age_seconds for liveness signals when scrapes are served from a stale cache.",
	}, []string{"source"})
	u.refreshFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: u.Namespace + "_prometheus_refresh_failures_total",
		Help: "Total background metrics refresh failures since process start.",
	})

	mux := http.NewServeMux()
	promver.Version = version.Version
	promver.Revision = version.Revision
	promver.Branch = version.Branch

	webserver.UpdateOutput(&webserver.Output{Name: PluginName, Config: u.Config})
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
	prometheus.MustRegister(u.controllerUp)
	prometheus.MustRegister(u.refreshFailures)
	prometheus.MustRegister(u)

	u.cache = &metricsCache{}
	prometheus.MustRegister(u.cacheAgeGauge())
	// safeRefresh (not refreshCache) because a panic in the initial upstream
	// fetch must not kill Run() before the HTTP listener starts.
	u.safeRefresh()

	go u.backgroundPoll()

	u.Logf("Prometheus scrape cache enabled, refresh interval: %v", u.Interval.Duration)

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

// normalizeInterval applies defaults and the minimum-interval floor to the
// configured scrape cache refresh interval. Values <= 0 use the default.
func (u *promUnifi) normalizeInterval() {
	if u.Interval.Duration <= 0 {
		u.Interval.Duration = defaultInterval

		return
	}

	if u.Interval.Duration < minimumInterval {
		u.Logf("Prometheus interval %v is below minimum %v; clamping to minimum",
			u.Interval.Duration, minimumInterval)

		u.Interval.Duration = minimumInterval
	}
}

// backgroundPoll runs forever, refreshing the metrics cache on the configured
// interval. Returns immediately if the cache is not configured. A panic in
// upstream collection is logged and the loop continues so one bad payload
// doesn't silently stop refreshes (operator would only see cache_age climb).
func (u *promUnifi) backgroundPoll() {
	if u.cache == nil {
		return
	}

	ticker := time.NewTicker(u.Interval.Duration)
	defer ticker.Stop()

	for range ticker.C {
		u.safeRefresh()
	}
}

// safeRefresh wraps refreshCache with a recover so a panic in an input
// plugin doesn't kill the background poller. The panic message and stack
// trace are logged separately so log aggregators (Sentry, Datadog, Loki)
// group panics by their headline rather than treating each stack line as
// an independent event.
func (u *promUnifi) safeRefresh() {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		u.LogErrorf("background metrics refresh panicked; continuing: %v", r)
		u.LogDebugf("panic stack:\n%s", debug.Stack())

		if u.cache != nil {
			u.cache.set(nil, fmt.Errorf("refresh panicked: %v", r))
		}

		if u.refreshFailures != nil {
			u.refreshFailures.Inc()
		}
	}()

	u.refreshCache()
}

// refreshCache polls upstream once and updates the cache. On error the last
// successful snapshot is preserved so scrapes keep returning data during
// transient upstream failures (e.g. 429 backoff loops). Failures are also
// counted in refreshFailures so operators can alert on failure rate
// independently of cache staleness.
//
// The cache.set call happens before logging and counter increment so the
// cache update is the most-protected statement — if anything in the
// post-set bookkeeping ever panics, safeRefresh's recover still sees a
// correctly updated cache.
func (u *promUnifi) refreshCache() {
	if u.cache == nil {
		return
	}

	m, err := u.Collector.Metrics(nil)
	u.cache.set(m, err)

	if err != nil {
		u.LogErrorf("background metrics refresh failed (serving last good snapshot): %v", err)

		if u.refreshFailures != nil {
			u.refreshFailures.Inc()
		}
	}
}

// cacheAgeGauge returns a GaugeFunc reporting seconds since the last
// successful background refresh. Returns -1 if no refresh has succeeded yet.
func (u *promUnifi) cacheAgeGauge() prometheus.Collector {
	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: u.Namespace + "_prometheus_cache_age_seconds",
			Help: "Seconds since the last successful background metrics refresh. -1 means no refresh has succeeded yet.",
		},
		func() float64 {
			if u.cache == nil {
				return -1
			}

			_, fetchedAt, _ := u.cache.get()
			if fetchedAt.IsZero() {
				return -1
			}

			return time.Since(fetchedAt).Seconds()
		},
	)
}

// fetchMetrics returns the metrics for a scrape, using the cache for global
// /metrics scrapes and singleflight-coalesced live calls for per-target
// /scrape requests.
func (u *promUnifi) fetchMetrics(filter *poller.Filter) (*poller.Metrics, error) {
	if filter == nil {
		if u.cache == nil {
			return u.Collector.Metrics(nil)
		}

		m, _, err := u.cache.get()
		if m != nil {
			// Serve cached data even if the most recent refresh errored.
			return m, nil
		}

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("metrics cache not yet populated")
	}

	// /scrape path: coalesce concurrent scrapes for the same target so a
	// noisy scraper can't multiply upstream API load.
	key := filter.Path
	if key == "" {
		key = filter.Name
	}

	result, err, _ := u.scrapeFlight.Do(key, func() (any, error) {
		return u.Collector.Metrics(filter)
	})
	if err != nil {
		return nil, err
	}

	// Strict assertion: silently dropping a wrong-type result would turn an
	// upstream regression into an empty 200 OK scrape with no log entry.
	m, ok := result.(*poller.Metrics)
	if !ok {
		return nil, fmt.Errorf("singleflight returned unexpected type %T", result)
	}

	return m, nil
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

// DefaultHandler serves the HTTP root with a simple liveness response naming
// the application.
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
	for _, f := range []any{
		u.Client, u.Device, u.UAP, u.USG, u.USW, u.PDU, u.Site, u.SpeedTest,
		u.DHCPLease, u.WAN, u.FirewallPolicy, u.Topology, u.PortAnomaly, u.VPNMesh,
		u.IntegrationDevice, u.WANStatus,
		u.PortForward, u.SSLCertificate, u.UPSDevice, u.WifiBroadcast,
		u.FirewallZone, u.ACLRule, u.VPNServer, u.SiteToSiteTunnel,
		u.LAG, u.MCLAGDomain, u.SwitchStack, u.DNSPolicy, u.RADIUSProfile,
		u.TrafficMatchingList, u.HotspotVoucher,
		u.DPIApplication, u.DPICategory, u.PendingDevice, u.Country,
	} {
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
		ch:     make(chan []*metric, u.Buffer),
		Start:  time.Now(),
	}
	defer r.close()

	r.Metrics, err = u.fetchMetrics(filter)
	r.Fetch = time.Since(r.Start)

	if err != nil {
		r.error(ch, prometheus.NewInvalidDesc(err), ErrMetricFetchFailed)
		u.LogErrorf("metric fetch failed: %v", err)

		return
	}

	// Export per-controller up/down gauge values.
	if u.controllerUp != nil && r.Metrics != nil {
		for _, cs := range r.Metrics.ControllerStatuses {
			val := 0.0
			if cs.Up {
				val = 1.0
			}

			u.controllerUp.WithLabelValues(cs.Source).Set(val)
		}
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

	for _, st := range m.SpeedTests {
		u.switchExport(r, st)
	}

	appTotal := make(totalsDPImap)
	catTotal := make(totalsDPImap)

	for _, c := range m.ClientsDPI {
		u.exportClientDPI(r, c, appTotal, catTotal)
	}

	for _, ct := range m.CountryTraffic {
		u.exportCountryTraffic(r, ct)
	}

	// Export network-level pool metrics first (once per network)
	dhcpLeases := make([]*unifi.DHCPLease, 0, len(m.DHCPLeases))
	for _, lease := range m.DHCPLeases {
		if l, ok := lease.(*unifi.DHCPLease); ok {
			dhcpLeases = append(dhcpLeases, l)
		}
	}

	if len(dhcpLeases) > 0 {
		u.exportDHCPNetworkPool(r, dhcpLeases)
	}

	// Export per-lease metrics
	for _, lease := range m.DHCPLeases {
		if l, ok := lease.(*unifi.DHCPLease); ok {
			u.exportDHCPLease(r, l)
		}
	}

	// Export WAN metrics
	for _, wanConfig := range m.WANConfigs {
		if w, ok := wanConfig.(*unifi.WANEnrichedConfiguration); ok {
			u.exportWAN(r, w)
		}
	}

	// Export controller sysinfo metrics
	for _, s := range m.Sysinfos {
		if sysinfo, ok := s.(*unifi.Sysinfo); ok {
			u.exportSysinfo(r, sysinfo)
		}
	}

	// Export firewall policy metrics
	firewallPolicies := make([]*unifi.FirewallPolicy, 0, len(m.FirewallPolicies))
	for _, p := range m.FirewallPolicies {
		if policy, ok := p.(*unifi.FirewallPolicy); ok {
			firewallPolicies = append(firewallPolicies, policy)
		}
	}

	u.exportFirewallPolicies(r, firewallPolicies)

	for _, t := range m.Topologies {
		if topo, ok := t.(*unifi.Topology); ok {
			u.exportTopology(r, topo)
		}
	}

	portAnomalies := make([]*unifi.PortAnomaly, 0, len(m.PortAnomalies))
	for _, a := range m.PortAnomalies {
		if anomaly, ok := a.(*unifi.PortAnomaly); ok {
			portAnomalies = append(portAnomalies, anomaly)
		}
	}

	u.exportPortAnomalies(r, portAnomalies)

	for _, v := range m.VPNMeshes {
		if mesh, ok := v.(*unifi.MagicSiteToSiteVPN); ok {
			u.exportVPNMesh(r, mesh)
		}
	}

	// v5.26.0 additions.
	for _, ds := range m.IntegrationDevStats {
		if d, ok := ds.(*unifi.IntegrationDeviceStats); ok {
			u.exportIntegrationDeviceStats(r, d)
		}
	}

	for _, ws := range m.WANStatuses {
		if w, ok := ws.(*unifi.WANStatus); ok {
			u.exportWANStatus(r, w)
		}
	}

	for _, pf := range m.PortForwards {
		if v, ok := pf.(*unifi.PortForward); ok {
			u.exportPortForward(r, v)
		}
	}

	for _, sc := range m.SSLCertificates {
		if v, ok := sc.(*unifi.SSLCertificate); ok {
			u.exportSSLCertificate(r, v)
		}
	}

	for _, ud := range m.UPSDevices {
		if v, ok := ud.(*unifi.UPSDeviceSelector); ok {
			u.exportUPSDevice(r, v)
		}
	}

	for _, wb := range m.WifiBroadcasts {
		if v, ok := wb.(*unifi.WifiBroadcast); ok {
			u.exportWifiBroadcast(r, v)
		}
	}

	for _, fz := range m.FirewallZones {
		if v, ok := fz.(*unifi.FirewallZone); ok {
			u.exportFirewallZone(r, v)
		}
	}

	for _, ar := range m.ACLRules {
		if v, ok := ar.(*unifi.ACLRule); ok {
			u.exportACLRule(r, v)
		}
	}

	for _, vs := range m.VPNServers {
		if v, ok := vs.(*unifi.VPNServer); ok {
			u.exportVPNServer(r, v)
		}
	}

	for _, st := range m.SiteToSiteTunnels {
		if v, ok := st.(*unifi.SiteToSiteTunnel); ok {
			u.exportSiteToSiteTunnel(r, v)
		}
	}

	for _, l := range m.LAGs {
		if v, ok := l.(*unifi.LAG); ok {
			u.exportLAG(r, v)
		}
	}

	for _, md := range m.MCLAGDomains {
		if v, ok := md.(*unifi.MCLAGDomain); ok {
			u.exportMCLAGDomain(r, v)
		}
	}

	for _, ss := range m.SwitchStacks {
		if v, ok := ss.(*unifi.SwitchStack); ok {
			u.exportSwitchStack(r, v)
		}
	}

	for _, dp := range m.DNSPolicies {
		if v, ok := dp.(*unifi.DNSPolicy); ok {
			u.exportDNSPolicy(r, v)
		}
	}

	for _, rp := range m.RADIUSProfiles {
		if v, ok := rp.(*unifi.RADIUSProfile); ok {
			u.exportRADIUSProfile(r, v)
		}
	}

	for _, tml := range m.TrafficMatchingLists {
		if v, ok := tml.(*unifi.TrafficMatchingList); ok {
			u.exportTrafficMatchingList(r, v)
		}
	}

	for _, hv := range m.HotspotVouchers {
		if v, ok := hv.(*unifi.HotspotVoucher); ok {
			u.exportHotspotVoucher(r, v)
		}
	}

	for _, app := range m.DPIApplications {
		if v, ok := app.(*unifi.DPIApplication); ok {
			u.exportDPIApplication(r, v)
		}
	}

	for _, cat := range m.DPICategories {
		if v, ok := cat.(*unifi.DPICategory); ok {
			u.exportDPICategory(r, v)
		}
	}

	for _, pd := range m.PendingDevices {
		if v, ok := pd.(*unifi.PendingDevice); ok {
			u.exportPendingDevice(r, v)
		}
	}

	for _, c := range m.Countries {
		if v, ok := c.(*unifi.Country); ok {
			u.exportCountry(r, v)
		}
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
	case *unifi.UDB:
		r.addUDB()
		u.exportUDB(r, v)
	case *unifi.UDM:
		r.addUDM()
		u.exportUDM(r, v)
	case *unifi.Site:
		u.exportSite(r, v)
	case *unifi.Client:
		u.exportClient(r, v)
	case *unifi.SpeedTestResult:
		u.exportSpeedTest(r, v)
	case *unifi.UsageByCountry:
		u.exportCountryTraffic(r, v)
	default:
		if u.Collector.Poller().LogUnknownTypes {
			u.LogDebugf("unknown type: %T", v)
		}
	}
}
