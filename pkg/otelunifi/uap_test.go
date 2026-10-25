//nolint:testpackage // white-box: exercises the unexported access point export.
package otelunifi

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestExportUAPRecordsOperatingChannelAndTxPower(t *testing.T) {
	t.Parallel()

	// radio_table is the configuration: with the default "auto" settings its channel decodes to 0 and
	// its tx_power is 0. The channel and power in use are in radio_table_stats, which lists the radios
	// in the opposite order here, and the two tables are matched by radio name. wifi2 has no stats entry.
	var ap unifi.UAP

	require.NoError(t, json.Unmarshal([]byte(`{"adopted": true, "mac": "aa:bb:cc:00:00:01",
		"radio_table": [
			{"name": "wifi0", "radio": "ng", "channel": "auto", "tx_power": "0", "tx_power_mode": "auto"},
			{"name": "wifi1", "radio": "na", "channel": "auto", "tx_power": "0", "tx_power_mode": "high"},
			{"name": "wifi2", "radio": "6e", "channel": "auto", "tx_power": "0", "tx_power_mode": "auto"}],
		"radio_table_stats": [
			{"name": "wifi1", "radio": "na", "channel": 149, "tx_power": 22},
			{"name": "wifi0", "radio": "ng", "channel": 6, "tx_power": 17}]}`), &ap))

	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })

	u := &OtelOutput{OtelUnifi: &OtelUnifi{Config: &Config{}}}
	u.exportUAP(t.Context(), provider.Meter("test"), &Report{}, &ap)

	var rm metricdata.ResourceMetrics

	require.NoError(t, reader.Collect(t.Context(), &rm))

	got := map[string]map[string]float64{}

	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			g, ok := m.Data.(metricdata.Gauge[float64])
			if !ok {
				continue
			}

			for _, dp := range g.DataPoints {
				radio, _ := dp.Attributes.Value(attribute.Key("radio_name"))

				if got[m.Name] == nil {
					got[m.Name] = map[string]float64{}
				}

				got[m.Name][radio.AsString()] = dp.Value
			}
		}
	}

	assert.Equal(t, map[string]float64{"wifi0": 6, "wifi1": 149}, got["unifi_device_uap_radio_channel"])
	assert.Equal(t, map[string]float64{"wifi0": 17, "wifi1": 22}, got["unifi_device_uap_radio_tx_power_dbm"])
}
