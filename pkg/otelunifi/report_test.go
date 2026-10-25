//nolint:testpackage // white-box: exercises the unexported client export.
package otelunifi

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
	"github.com/unpoller/unpoller/pkg/poller"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestExportClientsRecordsRadioRatesInBitsPerSecond(t *testing.T) {
	t.Parallel()

	// The controller reports tx_rate and rx_rate in kilobits per second: 866667 is an 866.7 Mbps link.
	var c unifi.Client

	require.NoError(t, json.Unmarshal(
		[]byte(`{"mac": "aa:bb:cc:dd:ee:01", "is_wired": false, "tx_rate": 866667, "rx_rate": 780000}`), &c))

	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })

	u := &OtelOutput{OtelUnifi: &OtelUnifi{Config: &Config{}}}
	u.exportClients(t.Context(), provider.Meter("test"), &poller.Metrics{Clients: []any{&c}}, &Report{})

	var rm metricdata.ResourceMetrics

	require.NoError(t, reader.Collect(t.Context(), &rm))

	got := map[string][]float64{}

	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if g, ok := m.Data.(metricdata.Gauge[float64]); ok {
				for _, dp := range g.DataPoints {
					got[m.Name] = append(got[m.Name], dp.Value)
				}
			}
		}
	}

	assert.Equal(t, []float64{866667000}, got["unifi_client_tx_rate_bps"])
	assert.Equal(t, []float64{780000000}, got["unifi_client_rx_rate_bps"])
}
