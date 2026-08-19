package poller_test

import (
	"reflect"
	"testing"
	"time"

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

// TestAppendMetricsCoversEverySliceField walks the Metrics struct by reflection and fails if
// any slice field is not merged by AppendMetrics.
//
// This is a guard against a silent, invisible failure mode: a new metric family needs both a
// field here and an append line there, and omitting the second loses every one of those
// metrics with no error, no log line, and a passing build. SpeedTests was lost that way --
// collected by inputunifi and exported by all three output plugins, but never merged, so the
// export code was dead for every user. Asserting field-by-field would just repeat the bug;
// reflection means a new field is covered the moment it is declared.
func TestAppendMetricsCoversEverySliceField(t *testing.T) {
	t.Parallel()

	// Put one zero-valued element in every slice, so a merged field has two and a dropped
	// field has one. The element values are irrelevant -- only the count is under test.
	fill := func(m *poller.Metrics) {
		v := reflect.ValueOf(m).Elem()
		for i := 0; i < v.NumField(); i++ {
			if f := v.Field(i); f.Kind() == reflect.Slice {
				f.Set(reflect.Append(f, reflect.Zero(f.Type().Elem())))
			}
		}
	}

	existing, incoming := &poller.Metrics{}, &poller.Metrics{}
	fill(existing)
	fill(incoming)

	got := reflect.ValueOf(poller.AppendMetrics(existing, incoming)).Elem()

	for i := 0; i < got.NumField(); i++ {
		f := got.Field(i)
		if f.Kind() != reflect.Slice {
			continue
		}

		assert.Equal(t, 2, f.Len(),
			"Metrics.%s is missing an append line in AppendMetrics, so those metrics are silently dropped",
			got.Type().Field(i).Name)
	}
}

func TestAppendMetricsPropagatesTS(t *testing.T) {
	t.Parallel()

	first := time.Now().Add(-time.Minute)
	second := time.Now()

	// The aggregate starts bare, so it must adopt the first input's timestamp -- otherwise it
	// stays zero and influxunifi stamps untimed points with the zero time.
	got := poller.AppendMetrics(&poller.Metrics{}, &poller.Metrics{TS: first})
	assert.Equal(t, first, got.TS)

	// A timestamp already set wins: it marks when the batch began.
	got = poller.AppendMetrics(&poller.Metrics{TS: first}, &poller.Metrics{TS: second})
	assert.Equal(t, first, got.TS)
}
