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

			if err := input.Initialize(u); err != nil {
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

			e, err := input.Events(filter)

			if err != nil {
				resultChan <- eventInputResult{err: err}

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

			m, err := input.Metrics(filter)
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
