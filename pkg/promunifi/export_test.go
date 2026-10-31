package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unpoller/pkg/poller"
)

// NewConsoleCollectorForTest exposes the real collection/export path to external
// tests without starting Run's HTTP server or registering global collectors.
func NewConsoleCollectorForTest(metrics *poller.Metrics) prometheus.Collector {
	return &consoleTestCollector{&promUnifi{
		Config:    &Config{Buffer: 50, ReportErrors: true},
		Collector: &stubCollect{metrics: metrics},
		Console:   descConsole("unpoller_"),
	}}
}

type consoleTestCollector struct {
	*promUnifi
}

func (c *consoleTestCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Console.Info

	ch <- c.Console.UpdateAvailable
}
