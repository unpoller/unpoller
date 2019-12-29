package poller

import (
	"fmt"
	"strings"
	"sync"

	"golift.io/unifi"
)

var (
	inputs    []*InputPlugin
	inputSync sync.Mutex
)

// Input plugins must implement this interface.
type Input interface {
	Initialize(Logger) error                     // Called once on startup to initialize the plugin.
	Metrics() (*Metrics, bool, error)            // Called every time new metrics are requested.
	MetricsFrom(*Filter) (*Metrics, bool, error) // Called every time new metrics are requested.
	RawMetrics(*Filter) ([]byte, error)
}

// InputPlugin describes an input plugin's consumable interface.
type InputPlugin struct {
	Name   string
	Config interface{} // Each config is passed into an unmarshaller later.
	Input
}

// Filter is used for metrics filters. Many fields for lots of expansion.
type Filter struct {
	Type string
	Term string
	Name string
	Tags string
	Role string
	Kind string
	Path string
	Area int
	Item int
	Unit int
	Sign int64
	Mass int64
	Rate float64
	Cost float64
	Free bool
	True bool
	Done bool
	Stop bool
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

// InitializeInputs runs the passed-in initializer method for each input plugin.
func (u *UnifiPoller) InitializeInputs() error {
	inputSync.Lock()
	defer inputSync.Unlock()

	for _, input := range inputs {
		// This must return, or the app locks up here.
		if err := input.Initialize(u); err != nil {
			return err
		}
	}

	return nil
}

// Metrics aggregates all the measurements from all configured inputs and returns them.
func (u *UnifiPoller) Metrics() (*Metrics, bool, error) {
	errs := []string{}
	metrics := &Metrics{}
	ok := false

	for _, input := range inputs {
		m, _, err := input.Metrics()
		if err != nil {
			errs = append(errs, err.Error())
		}

		if m == nil {
			continue
		}

		ok = true
		metrics = AppendMetrics(metrics, m)
	}

	var err error

	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, ", "))
	}

	return metrics, ok, err
}

// MetricsFrom aggregates all the measurements from filtered inputs and returns them.
func (u *UnifiPoller) MetricsFrom(filter *Filter) (*Metrics, bool, error) {
	errs := []string{}
	metrics := &Metrics{}
	ok := false

	for _, input := range inputs {
		if !strings.EqualFold(input.Name, filter.Name) {
			continue
		}

		m, _, err := input.MetricsFrom(filter)
		if err != nil {
			errs = append(errs, err.Error())
		}

		if m == nil {
			continue
		}

		ok = true
		metrics = AppendMetrics(metrics, m)
	}

	var err error

	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, ", "))
	}

	return metrics, ok, err
}

// AppendMetrics combined the metrics from two sources.
func AppendMetrics(existing *Metrics, m *Metrics) *Metrics {
	existing.SitesDPI = append(existing.SitesDPI, m.SitesDPI...)
	existing.Sites = append(existing.Sites, m.Sites...)
	existing.ClientsDPI = append(existing.ClientsDPI, m.ClientsDPI...)
	existing.Clients = append(existing.Clients, m.Clients...)
	existing.IDSList = append(existing.IDSList, m.IDSList...)

	if m.Devices == nil {
		return existing
	}

	if existing.Devices == nil {
		existing.Devices = &unifi.Devices{}
	}

	existing.UAPs = append(existing.UAPs, m.UAPs...)
	existing.USGs = append(existing.USGs, m.USGs...)
	existing.USWs = append(existing.USWs, m.USWs...)
	existing.UDMs = append(existing.UDMs, m.UDMs...)

	return existing
}
