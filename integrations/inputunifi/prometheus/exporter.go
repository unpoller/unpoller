package prometheus

import "github.com/davidnewhall/unifi-poller/metrics"

// Metrics contains all the data from the controller.
type Metrics struct {
	*metrics.Metrics
}

// ProcessExports turns the data into exported data.
func (m *Metrics) ProcessExports() []error {
	return nil
}
