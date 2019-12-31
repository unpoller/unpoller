package poller

import (
	"fmt"
	"sync"
)

var (
	outputs    []*Output
	outputSync sync.Mutex
)

// Collect is passed into output packages so they may collect metrics to output.
// Output packages must implement this interface.
type Collect interface {
	Metrics() (*Metrics, bool, error)
	MetricsFrom(*Filter) (*Metrics, bool, error)
	Logger
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

// InitializeOutputs runs all the configured output plugins.
// If none exist, or they all exit an error is returned.
func (u *UnifiPoller) InitializeOutputs() error {
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
		return fmt.Errorf("no output plugins imported")
	}

	for err := range v {
		if err != nil {
			return err
		}

		if count--; count == 0 {
			return fmt.Errorf("all output plugins have stopped, or none enabled")
		}
	}

	return nil
}
