//nolint:testpackage // white-box: exercises unexported batchUSGstats.
package datadogunifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

func TestBatchUSGstatsWithoutStatKeepsUplinkAndSpeedtest(t *testing.T) {
	t.Parallel()

	stats := (&DatadogUnifi{}).batchUSGstats(
		unifi.SpeedtestStatus{
			Latency:      unifi.FlexInt{Val: 15},
			XputDownload: unifi.FlexInt{Val: 646},
		},
		nil,
		unifi.Uplink{
			Latency: unifi.FlexInt{Val: 17},
			Speed:   unifi.FlexInt{Val: 1000},
		},
	)

	require.Equal(t, 17.0, stats["uplink_latency"])
	assert.Equal(t, 1000.0, stats["uplink_speed"])
	assert.Equal(t, 15.0, stats["speedtest_status_latency"])
	assert.Equal(t, 646.0, stats["speedtest_status_xput_download"])
	_, hasLAN := stats["lan_rx_bytes"]
	assert.False(t, hasLAN)
}

func TestBatchUSGstatsKeepsLANWhenStatPresent(t *testing.T) {
	t.Parallel()

	stats := (&DatadogUnifi{}).batchUSGstats(
		unifi.SpeedtestStatus{},
		&unifi.Gw{LanRxBytes: unifi.FlexInt{Val: 12345}},
		unifi.Uplink{},
	)

	assert.Equal(t, 12345.0, stats["lan_rx_bytes"])
	_, hasUplink := stats["uplink_latency"]
	assert.True(t, hasUplink)
}
