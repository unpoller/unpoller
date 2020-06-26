package poller

import (
	"fmt"
	"sync"
)

var (
	// These are used to keep track of loaded output plugins.
	outputs             []*Output    // nolint: gochecknoglobals
	outputSync          sync.RWMutex // nolint: gochecknoglobals
	errNoOutputPlugins  = fmt.Errorf("no output plugins imported")
	errAllOutputStopped = fmt.Errorf("all output plugins have stopped, or none enabled")
)

// Collect is passed into output packages so they may collect metrics to output.
type Collect interface {
	Logger
	Metrics(*Filter) (*Metrics, error)
	Events(*Filter) (*Events, error)
	// These get used by the webserver output plugin.
	Poller() Poller
	Inputs() []string
	Outputs() []string
}

// Output defines the output data for a metric exporter like influx or prometheus.
// Output packages should call NewOutput with this struct in init().
type Output struct {
	Name   string
	Config interface{}         // Each config is passed into an unmarshaller later.
	Method func(Collect) error // Called on startup for each configured output.
}

// NewOutput should be called by each output package's init function.
func NewOutput(o *Output) {
	outputSync.Lock()
	defer outputSync.Unlock()

	if o == nil || o.Method == nil {
		panic("nil output or method passed to poller.NewOutput")
	}

	outputs = append(outputs, o)
}

// Poller returns the poller config.
func (u *UnifiPoller) Poller() Poller {
	return *u.Config.Poller
}

// InitializeOutputs runs all the configured output plugins.
// If none exist, or they all exit an error is returned.
func (u *UnifiPoller) InitializeOutputs() error {
	outputSync.RLock()
	defer outputSync.RUnlock()

	v := make(chan error)
	defer close(v)

	var count int

	for _, o := range outputs {
		count++

		go func(o *Output) {
			v <- o.Method(u)
		}(o)
	}

	if count < 1 {
		return errNoOutputPlugins
	}

	for err := range v {
		if err != nil {
			return err
		}

		if count--; count == 0 {
			return errAllOutputStopped
		}
	}

	return nil
}

// Outputs allows other output plugins to see the list of loaded output plugins.
func (u *UnifiPoller) Outputs() (names []string) {
	outputSync.RLock()
	defer outputSync.RUnlock()

	for i := range outputs {
		names = append(names, outputs[i].Name)
	}

	return names
}
