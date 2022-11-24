package poller

import (
	"strings"
	"sync"
	"time"
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

// InitializeInputs runs the passed-in initializer method for each input plugin.
func (u *UnifiPoller) InitializeInputs() error {
	inputSync.RLock()
	defer inputSync.RUnlock()

	for _, input := range inputs {
		// This must return, or the app locks up here.
		if err := input.Initialize(u); err != nil {
			return err
		}
	}

	return nil
}

// Events aggregates log messages (events) from one or more sources.
func (u *UnifiPoller) Events(filter *Filter) (*Events, error) {
	inputSync.RLock()
	defer inputSync.RUnlock()

	events := Events{}

	for _, input := range inputs {
		if filter != nil &&
			filter.Name != "" &&
			!strings.EqualFold(input.Name, filter.Name) {
			continue
		}

		e, err := input.Events(filter)
		if err != nil {
			return &events, err
		}

		// Logs is the only member to extend at this time.
		events.Logs = append(events.Logs, e.Logs...)
	}

	return &events, nil
}

// Metrics aggregates all the measurements from filtered inputs and returns them.
// Passing a null filter returns everything!
func (u *UnifiPoller) Metrics(filter *Filter) (*Metrics, error) {
	inputSync.RLock()
	defer inputSync.RUnlock()

	metrics := &Metrics{}

	for _, input := range inputs {
		if filter != nil &&
			filter.Name != "" &&
			!strings.EqualFold(input.Name, filter.Name) {
			continue
		}

		m, err := input.Metrics(filter)
		if err != nil {
			return metrics, err
		}

		metrics = AppendMetrics(metrics, m)
	}

	return metrics, nil
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
