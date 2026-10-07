//nolint:testpackage // white-box: exercises unexported batchUSGstats.
package influxunifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

func TestBatchUSGstatsWithoutStatKeepsUplinkAndSpeedtest(t *testing.T) {
	t.Parallel()

	stats := (&InfluxUnifi{}).batchUSGstats(
		unifi.SpeedtestStatus{
			Latency:      unifi.FlexInt{Val: 15},
			XputDownload: unifi.FlexInt{Val: 646},
		},
		nil,
		unifi.Uplink{
			Name:    "wan",
			Latency: unifi.FlexInt{Val: 17},
			Speed:   unifi.FlexInt{Val: 1000},
		},
	)

	require.Equal(t, 17.0, stats["uplink_latency"])
	assert.Equal(t, "wan", stats["uplink_name"])
	assert.Equal(t, 1000.0, stats["uplink_speed"])
	assert.Equal(t, 15.0, stats["speedtest-status_latency"])
	assert.Equal(t, 646.0, stats["speedtest-status_xput_download"])
	_, hasLAN := stats["lan-rx_bytes"]
	assert.False(t, hasLAN)
}

func TestBatchUSGstatsKeepsLANWhenStatPresent(t *testing.T) {
	t.Parallel()

	stats := (&InfluxUnifi{}).batchUSGstats(
		unifi.SpeedtestStatus{},
		&unifi.Gw{LanRxBytes: unifi.FlexInt{Val: 12345}},
		unifi.Uplink{},
	)

	assert.Equal(t, 12345.0, stats["lan-rx_bytes"])
	_, hasUplink := stats["uplink_latency"]
	assert.True(t, hasUplink)
}
