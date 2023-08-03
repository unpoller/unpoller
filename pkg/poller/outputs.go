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

type OutputPlugin interface {
	Run(Collect) error
	Enabled() bool
	DebugOutput() (bool, error)
}

// Output defines the output data for a metric exporter like influx or prometheus.
// Output packages should call NewOutput with this struct in init().
type Output struct {
	Name   string
	Config any // Each config is passed into an unmarshaller later.
	OutputPlugin
}

// NewOutput should be called by each output package's init function.
func NewOutput(o *Output) {
	outputSync.Lock()
	defer outputSync.Unlock()

	if o == nil {
		panic("nil output passed to poller.NewOutput")
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
	count, errChan := u.runOutputMethods()
	defer close(errChan)

	if count == 0 {
		return errNoOutputPlugins
	}

	// Wait for and return an error from any output plugin.
	for err := range errChan {
		if err != nil {
			return err
		}

		if count--; count == 0 {
			return errAllOutputStopped
		}
	}

	return nil
}

func (u *UnifiPoller) runOutputMethods() (int, chan error) {
	outputSync.RLock()
	defer outputSync.RUnlock()

	return runOutputMethods(outputs, u, u)
}

func runOutputMethods(outputs []*Output, l Logger, c Collect) (int, chan error) {
	// Output plugin errors go into this channel.
	err := make(chan error)

	for _, o := range outputs {
		if o != nil && o.Enabled() {
			l.LogDebugf("output plugin enabled, starting run loop for %s", o.Name)
			
			go func(o *Output) {
				err <- o.Run(c) // Run each output plugin
			}(o)
		} else {
			l.LogDebugf("output plugin disabled for %s", o.Name)
		}
	}

	return len(outputs), err
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
