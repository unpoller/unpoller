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
	Initialize(Logger) error    // Called once on startup to initialize the plugin.
	Metrics() (*Metrics, error) // Called every time new metrics are requested.
}

// InputPlugin describes an input plugin's consumable interface.
type InputPlugin struct {
	Config interface{} // Each config is passed into an unmarshaller later.
	Input
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
func (u *UnifiPoller) Metrics() (*Metrics, error) {
	errs := []string{}
	metrics := &Metrics{}

	for _, input := range inputs {
		m, err := input.Metrics()
		if err != nil {
			errs = append(errs, err.Error())
		}

		if m == nil {
			continue
		}

		metrics.Sites = append(metrics.Sites, m.Sites...)
		metrics.Clients = append(metrics.Clients, m.Clients...)
		metrics.IDSList = append(metrics.IDSList, m.IDSList...)

		if m.Devices == nil {
			continue
		}

		if metrics.Devices == nil {
			metrics.Devices = &unifi.Devices{}
		}

		metrics.UAPs = append(metrics.UAPs, m.UAPs...)
		metrics.USGs = append(metrics.USGs, m.USGs...)
		metrics.USWs = append(metrics.USWs, m.USWs...)
		metrics.UDMs = append(metrics.UDMs, m.UDMs...)
	}

	var err error

	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, ", "))
	}

	return metrics, err
}
