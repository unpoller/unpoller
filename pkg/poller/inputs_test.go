package poller_test

import (
	"testing"

	"github.com/unpoller/unpoller/pkg/poller"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// panicInput is an Input that panics on Metrics and Events, simulating a
// malformed controller response crashing an input plugin. See issue #1030.
type panicInput struct{}

func (panicInput) Initialize(poller.Logger) error { return nil }

func (panicInput) Metrics(*poller.Filter) (*poller.Metrics, error) {
	panic("simulated aggregated-dashboard panic")
}

func (panicInput) Events(*poller.Filter) (*poller.Events, error) {
	panic("simulated aggregated-dashboard panic")
}

func (panicInput) RawMetrics(*poller.Filter) ([]byte, error) { return nil, nil }

func (panicInput) DebugInput() (bool, error) { return false, nil }

func TestCollectMetricsRecoversPanickingInput(t *testing.T) {
	t.Parallel()

	collector := poller.NewTestCollector(t)
	collector.AddInput(&poller.InputPlugin{Name: "panic-input", Input: panicInput{}})

	var metrics *poller.Metrics

	var err error

	require.NotPanics(t, func() {
		metrics, err = collector.Metrics(nil)
	})

	assert.NotNil(t, metrics)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "panic-input")
}

// nilEventsInput simulates a disabled input plugin, which returns (nil, nil)
// from Events. See https://github.com/unpoller/unpoller/issues/1030.
type nilEventsInput struct{}

func (nilEventsInput) Initialize(poller.Logger) error { return nil }

func (nilEventsInput) Metrics(*poller.Filter) (*poller.Metrics, error) { return nil, nil }

func (nilEventsInput) Events(*poller.Filter) (*poller.Events, error) { return nil, nil }

func (nilEventsInput) RawMetrics(*poller.Filter) ([]byte, error) { return nil, nil }

func (nilEventsInput) DebugInput() (bool, error) { return false, nil }

func TestCollectEventsHandlesNilEventsResult(t *testing.T) {
	t.Parallel()

	collector := poller.NewTestCollector(t)
	collector.AddInput(&poller.InputPlugin{Name: "nil-events-input", Input: nilEventsInput{}})

	var events *poller.Events

	var err error

	require.NotPanics(t, func() {
		events, err = collector.Events(nil)
	})

	assert.NotNil(t, events)
	require.NoError(t, err)
}

func TestCollectEventsRecoversPanickingInput(t *testing.T) {
	t.Parallel()

	collector := poller.NewTestCollector(t)
	collector.AddInput(&poller.InputPlugin{Name: "panic-input", Input: panicInput{}})

	var events *poller.Events

	var err error

	require.NotPanics(t, func() {
		events, err = collector.Events(nil)
	})

	assert.NotNil(t, events)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "panic-input")
}
