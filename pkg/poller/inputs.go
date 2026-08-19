package poller

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	// These are used ot keep track of loaded input plugins.
	inputs    []*InputPlugin // nolint: gochecknoglobals
	inputSync sync.RWMutex   // nolint: gochecknoglobals
)

// Input plugins must implement this interface.
type Input interface {
	Initialize(Logger) error           // Called once on startup to initialize the plugin.
	Metrics(*Filter) (*Metrics, error) // Called every time new metrics are requested.
	Events(*Filter) (*Events, error)   // This is new.
	RawMetrics(*Filter) ([]byte, error)
	DebugInput() (bool, error)
}

// Discoverer is an optional interface for inputs that can discover API endpoints.
type Discoverer interface {
	Discover(outputPath string) error
}

// InputPlugin describes an input plugin's consumable interface.
type InputPlugin struct {
	Name   string
	Config any // Each config is passed into an unmarshaller later.
	Input
}

// Filter is used for metrics filters. Many fields for lots of expansion.
type Filter struct {
	Type string
	Term string
	Name string
	Role string
	Kind string
	Path string
	Text string
	Unit int
	Pass bool
	Skip bool
	Time time.Time
	Dur  time.Duration
}

// NewInput creates a metric input. This should be called by input plugins
// init() functions.
func NewInput(i *InputPlugin) {
	inputSync.Lock()
	defer inputSync.Unlock()

	if i == nil || i.Input == nil {
		panic("nil output or method passed to poller.NewOutput")
	}

	inputs = append(inputs, i)
}

// recoverInitialize runs input.Initialize and converts a panic into an error.
// Each input plugin's Initialize call runs in its own goroutine, so a panic
// there cannot be caught by a recover() in the caller; left unrecovered, it
// crashes the whole process before the app even starts.
// See https://github.com/unpoller/unpoller/issues/1030
func recoverInitialize(input *InputPlugin, l Logger) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("input plugin %s panicked initializing: %v", input.Name, r) //nolint:err113
		}
	}()

	return input.Initialize(l)
}

// InitializeInputs runs the passed-in initializer method for each input plugin.
func (u *UnifiPoller) InitializeInputs() error {
	inputSync.RLock()
	defer inputSync.RUnlock()

	errChan := make(chan error, len(inputs))
	wg := &sync.WaitGroup{}

	// parallelize startup
	u.LogDebugf("initializing %d inputs", len(inputs))

	for _, input := range inputs {
		wg.Add(1)

		go func(input *InputPlugin) {
			defer wg.Done()

			// This must return, or the app locks up here.
			u.LogDebugf("inititalizing input... %s", input.Name)

			if err := recoverInitialize(input, u); err != nil {
				u.LogDebugf("error initializing input ... %s", input.Name)

				errChan <- err

				return
			}

			u.LogDebugf("input successfully initialized ... %s", input.Name)

			errChan <- nil
		}(input)
	}

	wg.Wait()
	close(errChan)

	u.LogDebugf("collecting input errors...")

	// collect errors if any.
	errs := make([]error, 0)

	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf("error initializing inputs")
		for _, e := range errs {
			err = errors.Wrap(err, e.Error())
		}
	}

	u.LogDebugf("returning error: %w", err)

	return err
}

type eventInputResult struct {
	logs []any
	err  error
}

// recoverEvents runs input.Events and converts a panic into an error. See
// recoverInitialize for why this is necessary.
func recoverEvents(input *InputPlugin, filter *Filter) (e *Events, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("input plugin %s panicked collecting events: %v", input.Name, r) //nolint:err113
		}
	}()

	return input.Events(filter)
}

func collectEvents(filter *Filter, inputs []*InputPlugin) (*Events, error) {
	resultChan := make(chan eventInputResult, len(inputs))
	wg := &sync.WaitGroup{}

	for _, input := range inputs {
		wg.Add(1)

		go func(input *InputPlugin) {
			defer wg.Done()

			if filter != nil &&
				filter.Name != "" &&
				!strings.EqualFold(input.Name, filter.Name) {
				resultChan <- eventInputResult{}

				return
			}

			e, err := recoverEvents(input, filter)
			if err != nil {
				resultChan <- eventInputResult{err: err}

				return
			}

			if e == nil {
				resultChan <- eventInputResult{}

				return
			}

			resultChan <- eventInputResult{logs: e.Logs}
		}(input)
	}

	wg.Wait()

	close(resultChan)

	events := Events{}
	errs := make([]error, 0)

	for result := range resultChan {
		if result.err != nil {
			errs = append(errs, result.err)
		} else if result.logs != nil {
			// Logs is the only member to extend at this time.
			events.Logs = append(events.Logs, result.logs...)
		}
	}

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf("error initializing inputs")
		for _, e := range errs {
			err = errors.Wrap(err, e.Error())
		}
	}

	return &events, err
}

// Events aggregates log messages (events) from one or more sources.
func (u *UnifiPoller) Events(filter *Filter) (*Events, error) {
	inputSync.RLock()
	defer inputSync.RUnlock()

	return collectEvents(filter, inputs)
}

type metricInputResult struct {
	metric *Metrics
	err    error
}

// recoverMetrics runs input.Metrics and converts a panic into an error (e.g.
// from a malformed controller response, such as an unexpected Site Speed
// Test aggregated-dashboard payload). See recoverInitialize for why this is
// necessary.
func recoverMetrics(input *InputPlugin, filter *Filter) (m *Metrics, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("input plugin %s panicked collecting metrics: %v", input.Name, r) //nolint:err113
		}
	}()

	return input.Metrics(filter)
}

func collectMetrics(filter *Filter, inputs []*InputPlugin) (*Metrics, error) {
	resultChan := make(chan metricInputResult, len(inputs))
	wg := &sync.WaitGroup{}

	for _, input := range inputs {
		wg.Add(1)

		go func(input *InputPlugin) {
			defer wg.Done()

			if filter != nil &&
				filter.Name != "" &&
				!strings.EqualFold(input.Name, filter.Name) {
				resultChan <- metricInputResult{}

				return
			}

			m, err := recoverMetrics(input, filter)
			resultChan <- metricInputResult{metric: m, err: err}
		}(input)
	}

	wg.Wait()

	close(resultChan)

	errs := make([]error, 0)
	metrics := &Metrics{}

	for result := range resultChan {
		if result.err != nil {
			errs = append(errs, result.err)
		} else if result.metric != nil {
			metrics = AppendMetrics(metrics, result.metric)
		}
	}

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf("error initializing inputs")
		for _, e := range errs {
			err = errors.Wrap(err, e.Error())
		}
	}

	return metrics, err
}

// Metrics aggregates all the measurements from filtered inputs and returns them.
// Passing a null filter returns everything!
func (u *UnifiPoller) Metrics(filter *Filter) (*Metrics, error) {
	inputSync.RLock()
	defer inputSync.RUnlock()

	return collectMetrics(filter, inputs)
}

// AppendMetrics combines the metrics from two sources.
func AppendMetrics(existing *Metrics, m *Metrics) *Metrics {
	if existing == nil {
		return m
	}

	if m == nil {
		return existing
	}

	existing.SitesDPI = append(existing.SitesDPI, m.SitesDPI...)
	existing.Sites = append(existing.Sites, m.Sites...)
	existing.ClientsDPI = append(existing.ClientsDPI, m.ClientsDPI...)
	existing.RogueAPs = append(existing.RogueAPs, m.RogueAPs...)
	existing.Clients = append(existing.Clients, m.Clients...)
	existing.Devices = append(existing.Devices, m.Devices...)
	existing.CountryTraffic = append(existing.CountryTraffic, m.CountryTraffic...)
	existing.DHCPLeases = append(existing.DHCPLeases, m.DHCPLeases...)
	existing.WANConfigs = append(existing.WANConfigs, m.WANConfigs...)
	existing.Sysinfos = append(existing.Sysinfos, m.Sysinfos...)
	existing.FirewallPolicies = append(existing.FirewallPolicies, m.FirewallPolicies...)
	existing.Topologies = append(existing.Topologies, m.Topologies...)
	existing.PortAnomalies = append(existing.PortAnomalies, m.PortAnomalies...)
	existing.VPNMeshes = append(existing.VPNMeshes, m.VPNMeshes...)
	existing.ControllerStatuses = append(existing.ControllerStatuses, m.ControllerStatuses...)
	existing.WANStatuses = append(existing.WANStatuses, m.WANStatuses...)
	existing.PortForwards = append(existing.PortForwards, m.PortForwards...)
	existing.SSLCertificates = append(existing.SSLCertificates, m.SSLCertificates...)
	existing.UPSDevices = append(existing.UPSDevices, m.UPSDevices...)
	existing.IntegrationDevStats = append(existing.IntegrationDevStats, m.IntegrationDevStats...)
	existing.WifiBroadcasts = append(existing.WifiBroadcasts, m.WifiBroadcasts...)
	existing.FirewallZones = append(existing.FirewallZones, m.FirewallZones...)
	existing.ACLRules = append(existing.ACLRules, m.ACLRules...)
	existing.VPNServers = append(existing.VPNServers, m.VPNServers...)
	existing.SiteToSiteTunnels = append(existing.SiteToSiteTunnels, m.SiteToSiteTunnels...)
	existing.LAGs = append(existing.LAGs, m.LAGs...)
	existing.MCLAGDomains = append(existing.MCLAGDomains, m.MCLAGDomains...)
	existing.SwitchStacks = append(existing.SwitchStacks, m.SwitchStacks...)
	existing.DNSPolicies = append(existing.DNSPolicies, m.DNSPolicies...)
	existing.RADIUSProfiles = append(existing.RADIUSProfiles, m.RADIUSProfiles...)
	existing.TrafficMatchingLists = append(existing.TrafficMatchingLists, m.TrafficMatchingLists...)
	existing.HotspotVouchers = append(existing.HotspotVouchers, m.HotspotVouchers...)
	existing.DPIApplications = append(existing.DPIApplications, m.DPIApplications...)
	existing.DPICategories = append(existing.DPICategories, m.DPICategories...)
	existing.PendingDevices = append(existing.PendingDevices, m.PendingDevices...)
	existing.Countries = append(existing.Countries, m.Countries...)
	existing.UNASDevices = append(existing.UNASDevices, m.UNASDevices...)

	return existing
}

// Inputs allows output plugins to see the list of loaded input plugins.
func (u *UnifiPoller) Inputs() (names []string) {
	inputSync.RLock()
	defer inputSync.RUnlock()

	for i := range inputs {
		names = append(names, inputs[i].Name)
	}

	return names
}
