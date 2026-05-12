//nolint:testpackage // white-box tests exercise unexported cache + fetch routing.
package promunifi

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	dto "github.com/prometheus/client_model/go"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unpoller/pkg/poller"
	"golift.io/cnfg"
)

// stubCollect is a minimal poller.Collect implementation for cache tests.
type stubCollect struct {
	mu       sync.Mutex
	calls    atomic.Int64
	delay    time.Duration
	metrics  *poller.Metrics
	err      error
	panicMsg string
	errLogs  atomic.Int64
}

func (s *stubCollect) Metrics(_ *poller.Filter) (*poller.Metrics, error) {
	s.calls.Add(1)

	s.mu.Lock()
	delay := s.delay
	m := s.metrics
	err := s.err
	panicMsg := s.panicMsg
	s.mu.Unlock()

	if panicMsg != "" {
		panic(panicMsg)
	}

	if delay > 0 {
		time.Sleep(delay)
	}

	return m, err
}

func (s *stubCollect) Events(_ *poller.Filter) (*poller.Events, error) { return &poller.Events{}, nil }
func (s *stubCollect) Poller() poller.Poller                           { return poller.Poller{} }
func (s *stubCollect) Inputs() []string                                { return nil }
func (s *stubCollect) Outputs() []string                               { return nil }
func (s *stubCollect) Logf(string, ...any)                             {}
func (s *stubCollect) LogErrorf(string, ...any)                        { s.errLogs.Add(1) }
func (s *stubCollect) LogDebugf(string, ...any)                        {}

func TestMetricsCacheSetKeepsLastGoodOnError(t *testing.T) {
	t.Parallel()

	c := &metricsCache{}
	good := &poller.Metrics{}

	c.set(good, nil)
	c.set(nil, errors.New("transient 429"))

	m, fetchedAt, err := c.get()
	require.NotNil(t, m, "cache should retain last good snapshot on error")
	assert.Same(t, good, m)
	assert.False(t, fetchedAt.IsZero())
	assert.EqualError(t, err, "transient 429")
}

func TestMetricsCacheSetReplacesOnSuccess(t *testing.T) {
	t.Parallel()

	c := &metricsCache{}
	first := &poller.Metrics{}
	second := &poller.Metrics{}

	c.set(first, nil)
	c.set(second, nil)

	m, _, err := c.get()
	assert.NoError(t, err)
	assert.Same(t, second, m)
}

func TestFetchMetricsServesCacheForGlobalScrape(t *testing.T) {
	t.Parallel()

	stub := &stubCollect{}
	cached := &poller.Metrics{}

	u := &promUnifi{Config: &Config{}, Collector: stub, cache: &metricsCache{}}
	u.cache.set(cached, nil)

	got, err := u.fetchMetrics(nil)
	require.NoError(t, err)
	assert.Same(t, cached, got)
	assert.Zero(t, stub.calls.Load(), "global scrape must not hit upstream when cache has data")
}

func TestFetchMetricsReturnsErrorWhenCacheEmpty(t *testing.T) {
	t.Parallel()

	stub := &stubCollect{}
	u := &promUnifi{Config: &Config{}, Collector: stub, cache: &metricsCache{}}

	got, err := u.fetchMetrics(nil)
	assert.Nil(t, got)
	assert.Error(t, err)
}

// TestFetchMetricsBypassesNilCache covers the nil-cache defensive branch.
// Run() always initializes the cache in production; this branch exists only
// for tests that exercise fetchMetrics directly without invoking Run.
func TestFetchMetricsBypassesNilCache(t *testing.T) {
	t.Parallel()

	stub := &stubCollect{metrics: &poller.Metrics{}}
	u := &promUnifi{Config: &Config{}, Collector: stub}

	got, err := u.fetchMetrics(nil)
	require.NoError(t, err)
	assert.Same(t, stub.metrics, got)
	assert.EqualValues(t, 1, stub.calls.Load())
}

func TestFetchMetricsSingleflightCoalescesScrapes(t *testing.T) {
	t.Parallel()

	stub := &stubCollect{metrics: &poller.Metrics{}, delay: 50 * time.Millisecond}
	u := &promUnifi{Config: &Config{}, Collector: stub}

	const concurrent = 20

	var wg sync.WaitGroup

	wg.Add(concurrent)

	for i := 0; i < concurrent; i++ {
		go func() {
			defer wg.Done()

			_, _ = u.fetchMetrics(&poller.Filter{Path: "https://controller.example/"})
		}()
	}

	wg.Wait()

	assert.LessOrEqual(t, stub.calls.Load(), int64(2),
		"singleflight should coalesce concurrent scrapes to ~1 upstream call")
}

func TestFetchMetricsScrapeKeyedByFilterPath(t *testing.T) {
	t.Parallel()

	stub := &stubCollect{metrics: &poller.Metrics{}, delay: 50 * time.Millisecond}
	u := &promUnifi{Config: &Config{}, Collector: stub}

	var wg sync.WaitGroup

	wg.Add(2)

	for i, path := range []string{"https://a.example/", "https://b.example/"} {
		i, path := i, path

		go func() {
			defer wg.Done()

			_, _ = u.fetchMetrics(&poller.Filter{Path: path, Name: fmt.Sprintf("t%d", i)})
		}()
	}

	wg.Wait()

	assert.EqualValues(t, 2, stub.calls.Load(),
		"distinct targets should not be coalesced")
}

func TestNormalizeInterval(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   time.Duration
		want time.Duration
	}{
		{"unset uses default", 0, defaultInterval},
		{"negative uses default", -1, defaultInterval},
		{"below minimum clamps", time.Second, minimumInterval},
		{"exact minimum unchanged", minimumInterval, minimumInterval},
		{"above minimum unchanged", 2 * time.Minute, 2 * time.Minute},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u := &promUnifi{Config: &Config{Interval: cnfg.Duration{Duration: tc.in}}}
			u.normalizeInterval()
			assert.Equal(t, tc.want, u.Interval.Duration)
		})
	}
}

func TestRefreshCachePreservesLastGoodAcrossError(t *testing.T) {
	t.Parallel()

	good := &poller.Metrics{}
	stub := &stubCollect{metrics: good}
	u := &promUnifi{Config: &Config{}, Collector: stub, cache: &metricsCache{}}

	u.refreshCache()
	stub.mu.Lock()
	stub.metrics = nil
	stub.err = errors.New("429 too many requests")
	stub.mu.Unlock()
	u.refreshCache()

	m, fetchedAt, err := u.cache.get()
	require.Same(t, good, m, "good snapshot must survive a failed refresh")
	assert.False(t, fetchedAt.IsZero())
	assert.EqualError(t, err, "429 too many requests")
	assert.EqualValues(t, 1, stub.errLogs.Load(), "failed refresh should log via LogErrorf")
}

func TestCacheAgeGaugeReportsUnpopulated(t *testing.T) {
	t.Parallel()

	u := &promUnifi{Config: &Config{Namespace: "test"}, cache: &metricsCache{}}
	gauge := u.cacheAgeGauge()

	require.Equal(t, -1.0, currentGaugeValue(t, gauge),
		"unpopulated cache should report -1 sentinel")
}

func TestCacheAgeGaugeReportsFreshness(t *testing.T) {
	t.Parallel()

	u := &promUnifi{Config: &Config{Namespace: "test"}, cache: &metricsCache{}}
	u.cache.set(&poller.Metrics{}, nil)

	gauge := u.cacheAgeGauge()
	age := currentGaugeValue(t, gauge)
	assert.GreaterOrEqual(t, age, 0.0)
	assert.Less(t, age, 5.0, "freshly populated cache should report a small age")
}

func TestSafeRefreshRecoversFromPanic(t *testing.T) {
	t.Parallel()

	good := &poller.Metrics{}
	stub := &stubCollect{metrics: good}
	u := &promUnifi{Config: &Config{}, Collector: stub, cache: &metricsCache{}}

	// First a successful refresh so we have a good snapshot to preserve.
	u.safeRefresh()
	stub.mu.Lock()
	stub.metrics = nil
	stub.panicMsg = "simulated input plugin panic"
	stub.mu.Unlock()

	require.NotPanics(t, u.safeRefresh, "safeRefresh must swallow upstream panics")

	m, fetchedAt, err := u.cache.get()
	assert.Same(t, good, m, "panic must not blank out the cache")
	assert.False(t, fetchedAt.IsZero())
	assert.ErrorContains(t, err, "refresh panicked")
	assert.EqualValues(t, 1, stub.errLogs.Load(), "panic should be logged once")
}

func TestRefreshFailuresCounterIncrementsOnErrorAndPanic(t *testing.T) {
	t.Parallel()

	stub := &stubCollect{err: errors.New("upstream down")}
	counter := prometheus.NewCounter(prometheus.CounterOpts{Name: "test_refresh_failures_total"})
	u := &promUnifi{
		Config:          &Config{},
		Collector:       stub,
		cache:           &metricsCache{},
		refreshFailures: counter,
	}

	u.refreshCache()
	u.refreshCache()
	assert.Equal(t, 2.0, currentCounterValue(t, counter), "error path should increment counter")

	stub.mu.Lock()
	stub.err = nil
	stub.panicMsg = "boom"
	stub.mu.Unlock()
	u.safeRefresh()
	assert.Equal(t, 3.0, currentCounterValue(t, counter), "panic path should increment counter")
}

func TestFetchMetricsScrapeReturnsUpstreamError(t *testing.T) {
	t.Parallel()

	stub := &stubCollect{err: errors.New("controller unreachable")}
	u := &promUnifi{Config: &Config{}, Collector: stub}

	got, err := u.fetchMetrics(&poller.Filter{Path: "https://x.example/"})
	assert.Nil(t, got)
	assert.EqualError(t, err, "controller unreachable")
}

// currentGaugeValue reads the current value out of a GaugeFunc by collecting
// it through Prometheus' standard interface — avoids depending on internal
// client_golang APIs.
func currentGaugeValue(t *testing.T, c prometheus.Collector) float64 {
	t.Helper()

	ch := make(chan prometheus.Metric, 1)
	c.Collect(ch)
	close(ch)

	m := <-ch

	var pb dto.Metric

	require.NoError(t, m.Write(&pb))

	return pb.GetGauge().GetValue()
}

func currentCounterValue(t *testing.T, c prometheus.Collector) float64 {
	t.Helper()

	ch := make(chan prometheus.Metric, 1)
	c.Collect(ch)
	close(ch)

	m := <-ch

	var pb dto.Metric

	require.NoError(t, m.Write(&pb))

	return pb.GetCounter().GetValue()
}
