//nolint:testpackage // white-box: exercises the unexported switch port export.
package otelunifi

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestExportUSWPortRecordsPoECurrentInAmps(t *testing.T) {
	t.Parallel()

	// The controller reports poe_current in milliamps: 4.84 W at 52.88 V is 0.0915 A.
	var port unifi.Port

	require.NoError(t, json.Unmarshal(
		[]byte(`{"port_idx": 3, "poe_current": "91.55", "poe_power": "4.84", "poe_voltage": "52.88"}`), &port))

	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })

	u := &OtelOutput{OtelUnifi: &OtelUnifi{Config: &Config{}}}
	u.exportUSWPort(t.Context(), provider.Meter("test"), &Report{}, &unifi.USW{}, port)

	var rm metricdata.ResourceMetrics

	require.NoError(t, reader.Collect(t.Context(), &rm))

	var got []float64

	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if g, ok := m.Data.(metricdata.Gauge[float64]); ok && m.Name == "unifi_device_usw_port_poe_current_amps" {
				for _, dp := range g.DataPoints {
					got = append(got, dp.Value)
				}
			}
		}
	}

	require.Len(t, got, 1)
	assert.InDelta(t, 0.09155, got[0], 1e-9)
	assert.InDelta(t, port.PoePower.Val/port.PoeVoltage.Val, got[0], 0.001)
}
