package promunifi_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
	"github.com/unpoller/unpoller/pkg/poller"
	"github.com/unpoller/unpoller/pkg/promunifi"
)

func consoleSystem(t *testing.T, source, body string) *unifi.ConsoleSystem {
	t.Helper()

	var system unifi.ConsoleSystem

	require.NoError(t, json.Unmarshal([]byte(body), &system))
	system.SourceName = source

	return &system
}

func TestConsoleMetrics(t *testing.T) {
	t.Parallel()

	metrics := &poller.Metrics{ConsoleSystems: []any{
		consoleSystem(t, "home", `{"devices":{"unifiOS":[{"name":"UCG Fiber","firmware":"5.1.33",
			"firmwareStatus":"updateAvailable","updateAvailable":"5.1.34+4d109cd"}]}}`),
		consoleSystem(t, "office", `{"devices":{"unifiOS":[{"name":"UCG Fiber","firmware":"5.1.34",
			"firmwareStatus":"upToDate","updateAvailable":null}]}}`),
	}}
	collector := promunifi.NewConsoleCollectorForTest(metrics)
	registry := prometheus.NewPedanticRegistry()
	require.NoError(t, registry.Register(collector))
	require.NoError(t, testutil.GatherAndCompare(registry, strings.NewReader(`
# HELP unpoller_console_info UniFi OS console firmware information (always 1)
# TYPE unpoller_console_info gauge
unpoller_console_info{available_version="",firmware_status="upToDate",name="UCG Fiber",source="office",version="5.1.34"} 1
unpoller_console_info{available_version="5.1.34+4d109cd",firmware_status="updateAvailable",name="UCG Fiber",source="home",version="5.1.33"} 1
# HELP unpoller_console_update_available UniFi OS update available (1/0); omitted when unknown
# TYPE unpoller_console_update_available gauge
unpoller_console_update_available{name="UCG Fiber",source="home"} 1
unpoller_console_update_available{name="UCG Fiber",source="office"} 0
`)))
}

func TestConsoleUnknownStatusOmitsUpdateGauge(t *testing.T) {
	t.Parallel()

	for _, body := range []string{
		`{"devices":{"unifiOS":[{"name":"Console","firmware":"5.1.33"}]}}`,
		`{"devices":{"unifiOS":[{"name":"Console","firmwareStatus":"checking","updateAvailable":null}]}}`,
		`{"devices":{"unifiOS":[null]}}`,
		`{}`,
	} {
		m := &poller.Metrics{ConsoleSystems: []any{consoleSystem(t, "home", body), (*unifi.ConsoleSystem)(nil)}}
		collector := promunifi.NewConsoleCollectorForTest(m)
		count := testutil.CollectAndCount(collector, "unpoller_console_update_available")
		assert.Zero(t, count, "unknown data must not become a no-update result")
	}

	count := testutil.CollectAndCount(promunifi.NewConsoleCollectorForTest(&poller.Metrics{}))
	assert.Zero(t, count, "a failed or unsupported collection has no console metrics")
}
