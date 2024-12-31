package datadogunifi_test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unpoller/pkg/datadogunifi"
	"github.com/unpoller/unpoller/pkg/unittest"
	"golift.io/cnfg"
	"gopkg.in/yaml.v3"
)

type mockValue struct {
	value any
	tags  []string
}

// mockStatsd allows us to mock statsd.ClientInterface and collect data to ensure we're writing out
// metrics as expected with the correct types
type mockStatsd struct {
	sync.RWMutex
	gauges        map[string]mockValue
	counts        map[string]mockValue
	histograms    map[string]mockValue
	distributions map[string]mockValue
	sets          map[string]mockValue
	timings       map[string]mockValue
	events        []string
	checks        []string
}

// GaugeWithTimestamp mock interface
// nolint:all
func (m *mockStatsd) GaugeWithTimestamp(name string, value float64, tags []string, rate float64, timestamp time.Time) error {
	// not supported
	return nil
}

// CountWithTimestamp mock interface
// nolint:all
func (m *mockStatsd) CountWithTimestamp(name string, value int64, tags []string, rate float64, timestamp time.Time) error {
	// not supported
	return nil
}

// IsClosed mock interface
// nolint:all
func (m *mockStatsd) IsClosed() bool {
	return false
}

// HistogramWithTimestamp mock interface
// nolint:all
func (m *mockStatsd) HistogramWithTimestamp(name string, value float64, tags []string, rate float64, timestamp time.Time) error {
	return nil
}

// DistributionWithTimestamp mock interface
// nolint:all
func (m *mockStatsd) DistributionWithTimestamp(name string, value float64, tags []string, rate float64, timestamp time.Time) error {
	return nil
}

// SetWithTimestamp mock interface
// nolint:all
func (m *mockStatsd) SetWithTimestamp(name string, value float64, tags []string, rate float64, timestamp time.Time) error {
	return nil
}

// TimingWithTimestamp mock interface
// nolint:all
func (m *mockStatsd) TimingWithTimestamp(name string, value int64, tags []string, rate float64) error {
	return nil
}

// GetTelemetry mock interface
// nolint:all
func (m *mockStatsd) GetTelemetry() statsd.Telemetry {
	return statsd.Telemetry{}
}

// Gauge measures the value of a metric at a particular time.
func (m *mockStatsd) Gauge(name string, value float64, tags []string, _ float64) error {
	m.Lock()
	defer m.Unlock()

	m.gauges[name] = mockValue{value, tags}

	return nil
}

// Count tracks how many times something happened per second.
func (m *mockStatsd) Count(name string, value int64, tags []string, _ float64) error {
	m.Lock()
	defer m.Unlock()

	m.counts[name] = mockValue{value, tags}

	return nil
}

// Histogram tracks the statistical distribution of a set of values on each host.
func (m *mockStatsd) Histogram(name string, value float64, tags []string, _ float64) error {
	m.Lock()
	defer m.Unlock()

	m.histograms[name] = mockValue{value, tags}

	return nil
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func (m *mockStatsd) Distribution(name string, value float64, tags []string, _ float64) error {
	m.Lock()
	defer m.Unlock()

	m.distributions[name] = mockValue{value, tags}

	return nil
}

// Decr is just Count of -1
func (m *mockStatsd) Decr(name string, tags []string, _ float64) error {
	m.Lock()
	defer m.Unlock()

	m.counts[name] = mockValue{-1, tags}

	return nil
}

// Incr is just Count of 1
func (m *mockStatsd) Incr(name string, tags []string, _ float64) error {
	m.Lock()
	defer m.Unlock()

	m.counts[name] = mockValue{1, tags}

	return nil
}

// Set counts the number of unique elements in a group.
func (m *mockStatsd) Set(name string, value string, tags []string, _ float64) error {
	m.Lock()
	defer m.Unlock()

	m.sets[name] = mockValue{value, tags}

	return nil
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func (m *mockStatsd) Timing(name string, value time.Duration, tags []string, _ float64) error {
	m.Lock()
	defer m.Unlock()

	m.timings[name] = mockValue{value, tags}

	return nil
}

// TimeInMilliseconds sends timing information in milliseconds.
// It is flushed by statsd with percentiles, mean and other info (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing)
func (m *mockStatsd) TimeInMilliseconds(name string, value float64, tags []string, _ float64) error {
	m.Lock()
	defer m.Unlock()

	m.timings[name] = mockValue{value, tags}

	return nil
}

// Event sends the provided Event.
func (m *mockStatsd) Event(e *statsd.Event) error {
	m.Lock()
	defer m.Unlock()

	m.events = append(m.events, e.Title)

	return nil
}

// SimpleEvent sends an event with the provided title and text.
func (m *mockStatsd) SimpleEvent(title, _ string) error {
	m.Lock()
	defer m.Unlock()

	m.events = append(m.events, title)

	return nil
}

// ServiceCheck sends the provided ServiceCheck.
func (m *mockStatsd) ServiceCheck(sc *statsd.ServiceCheck) error {
	m.Lock()
	defer m.Unlock()

	m.checks = append(m.checks, sc.Name)

	return nil
}

// SimpleServiceCheck sends an serviceCheck with the provided name and status.
func (m *mockStatsd) SimpleServiceCheck(name string, _ statsd.ServiceCheckStatus) error {
	m.Lock()
	defer m.Unlock()

	m.checks = append(m.checks, name)

	return nil
}

// Close the client connection.
func (m *mockStatsd) Close() error {
	return nil
}

// Flush forces a flush of all the queued dogstatsd payloads.
func (m *mockStatsd) Flush() error {
	return nil
}

// SetWriteTimeout allows the user to set a custom write timeout.
func (m *mockStatsd) SetWriteTimeout(_ time.Duration) error {
	return nil
}

type testExpectations struct {
	Gauges        []string `yaml:"gauges"`
	Counts        []string `yaml:"counts"`
	Timings       []string `yaml:"timings"`
	Sets          []string `yaml:"sets"`
	Histograms    []string `yaml:"histograms"`
	Distributions []string `yaml:"distributions"`
	ServiceChecks []string `yaml:"service_checks"`
}

func TestDataDogUnifiIntegration(t *testing.T) {
	// load test expectations file
	yamlFile, err := os.ReadFile("integration_test_expectations.yaml")
	require.NoError(t, err)

	var testExpectationsData testExpectations
	err = yaml.Unmarshal(yamlFile, &testExpectationsData)
	require.NoError(t, err)

	testRig := unittest.NewTestSetup(t)
	defer testRig.Close()

	mockCapture := &mockStatsd{
		gauges:        make(map[string]mockValue, 0),
		counts:        make(map[string]mockValue, 0),
		histograms:    make(map[string]mockValue, 0),
		distributions: make(map[string]mockValue, 0),
		sets:          make(map[string]mockValue, 0),
		timings:       make(map[string]mockValue, 0),
		events:        make([]string, 0),
		checks:        make([]string, 0),
	}

	u := datadogunifi.DatadogUnifi{
		Datadog: &datadogunifi.Datadog{
			Config: &datadogunifi.Config{
				Enable:   unittest.PBool(true),
				Interval: cnfg.Duration{Duration: time.Hour},
			},
		},
		Statsd: mockCapture,
	}

	testRig.Initialize()

	u.Collector = testRig.Collector
	u.Collect(time.Minute)
	mockCapture.RLock()
	defer mockCapture.RUnlock()

	// gauges
	assert.Equal(t, len(testExpectationsData.Gauges), len(mockCapture.gauges))

	expectedKeys := unittest.NewSetFromSlice[string](testExpectationsData.Gauges)
	foundKeys := unittest.NewSetFromMap[string](mockCapture.gauges)
	additions, deletions := expectedKeys.Difference(foundKeys)
	assert.Len(t, additions, 0)
	assert.Len(t, deletions, 0)

	// counts
	assert.Len(t, mockCapture.counts, 12)

	expectedKeys = unittest.NewSetFromSlice[string](testExpectationsData.Counts)
	foundKeys = unittest.NewSetFromMap[string](mockCapture.counts)
	additions, deletions = expectedKeys.Difference(foundKeys)
	assert.Len(t, additions, 0)
	assert.Len(t, deletions, 0)

	// timings
	assert.Len(t, mockCapture.timings, 2)

	expectedKeys = unittest.NewSetFromSlice[string](testExpectationsData.Timings)
	foundKeys = unittest.NewSetFromMap[string](mockCapture.timings)
	additions, deletions = expectedKeys.Difference(foundKeys)
	assert.Len(t, additions, 0)
	assert.Len(t, deletions, 0)

	// histograms
	assert.Len(t, mockCapture.histograms, 0)

	expectedKeys = unittest.NewSetFromSlice[string](testExpectationsData.Histograms)
	foundKeys = unittest.NewSetFromMap[string](mockCapture.histograms)
	additions, deletions = expectedKeys.Difference(foundKeys)
	assert.Len(t, additions, 0)
	assert.Len(t, deletions, 0)

	// distributions
	assert.Len(t, mockCapture.distributions, 0)

	expectedKeys = unittest.NewSetFromSlice[string](testExpectationsData.Distributions)
	foundKeys = unittest.NewSetFromMap[string](mockCapture.distributions)
	additions, deletions = expectedKeys.Difference(foundKeys)
	assert.Len(t, additions, 0)
	assert.Len(t, deletions, 0)

	// sets
	assert.Len(t, mockCapture.sets, 0)

	expectedKeys = unittest.NewSetFromSlice[string](testExpectationsData.Sets)
	foundKeys = unittest.NewSetFromMap[string](mockCapture.sets)
	additions, deletions = expectedKeys.Difference(foundKeys)
	assert.Len(t, additions, 0)
	assert.Len(t, deletions, 0)

	// events
	// at least one event from an alarm should happen
	assert.GreaterOrEqual(t, len(mockCapture.events), 1)

	// service checks
	assert.Len(t, mockCapture.checks, 0)

	expectedKeys = unittest.NewSetFromSlice[string](testExpectationsData.ServiceChecks)
	foundKeys = unittest.NewSetFromSlice[string](mockCapture.checks)
	additions, deletions = expectedKeys.Difference(foundKeys)
	assert.Len(t, additions, 0)
	assert.Len(t, deletions, 0)
}
