//nolint:testpackage // white-box: exercises the unexported batch functions.
package influxunifi

import (
	"encoding/json"
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

// captureReport records sent points; other report methods are not used by the WAN batches.
type captureReport struct {
	report

	sent []*metric
}

func (c *captureReport) send(m *metric) { c.sent = append(c.sent, m) }

func TestBatchUSGwansOmitsHealthWhenAbsent(t *testing.T) {
	t.Parallel()

	r := &captureReport{}
	(&InfluxUnifi{}).batchUSGwans(r, map[string]string{"name": "Gateway"},
		unifi.Wan{Name: "wan", Up: unifi.FlexBool{Val: true}},
		unifi.Wan{Name: "wan2", Up: unifi.FlexBool{Val: true}, Latency: *unifi.NewFlexInt(30), Availability: *unifi.NewFlexInt(100)},
		unifi.Wan{Name: "wan3", Up: unifi.FlexBool{Val: true}, Latency: *unifi.NewFlexInt(0), Availability: *unifi.NewFlexInt(0)},
	)

	require.Len(t, r.sent, 3)
	assert.NotContains(t, r.sent[0].Fields, "latency")
	assert.NotContains(t, r.sent[0].Fields, "availability")
	assert.Equal(t, 30.0, r.sent[1].Fields["latency"])
	assert.Equal(t, 100.0, r.sent[1].Fields["availability"])
	assert.NotContains(t, r.sent[2].Fields, "latency", "an offline WAN's 0 latency is skipped")
	assert.Equal(t, 0.0, r.sent[2].Fields["availability"])
}

func TestBatchUptimeMonitorsWithoutAvailability(t *testing.T) {
	t.Parallel()

	r := &captureReport{}
	(&InfluxUnifi{}).batchUptimeMonitors(r, map[string]string{"wan_networkgroup": "WAN"}, "false", []unifi.UptimeMonitor{
		{Target: "absent.example", Type: "icmp"},
		{Target: "zero.example", Type: "icmp", Availability: *unifi.NewFlexInt(0)},
	})

	require.Len(t, r.sent, 1, "a monitor with no fields is skipped")
	assert.Equal(t, "zero.example", r.sent[0].Tags["target"])
	assert.Equal(t, map[string]any{"availability": 0.0}, r.sent[0].Fields)
}

func TestBatchWANHealthDuringOutage(t *testing.T) {
	t.Parallel()

	// Trimmed from a UDM Pro (5.1.33, Network 10.6.106) a minute into a real WAN1 outage.
	const raw = `{
		"WAN": {
			"downtime": 66,
			"monitors": [{"target": "google.com", "type": "icmp", "availability": 100.0, "latency_average": 0}],
			"alerting_monitors": [{"target": "ping.ui.com", "type": "icmp", "availability": 0.0}]
		},
		"WAN2": {"availability": 100.0, "latency_average": 21, "time_period": 86400, "uptime": 1429359}
	}`

	var stats map[string]*unifi.UptimeStats

	require.NoError(t, json.Unmarshal([]byte(raw), &stats))

	status := map[string]string{"WAN": "offline", "WAN2": "online"}

	r := &captureReport{}
	(&InfluxUnifi{}).batchWANHealth(r, map[string]string{"name": "Gateway", "site_name": "Default"}, stats, status)

	points := map[string]*metric{}
	for _, m := range r.sent {
		points[m.Table+"/"+m.Tags["wan_networkgroup"]+"/"+m.Tags["target"]] = m
	}

	a := assert.New(t)
	require.Len(t, points, 6)
	a.Equal(map[string]any{"status": "offline"}, points["usg_wan_status/WAN/"].Fields)
	a.Equal(map[string]any{"status": "online"}, points["usg_wan_status/WAN2/"].Fields)
	a.Equal(map[string]string{"device_name": "Gateway", "site_name": "Default", "source": "", "wan_networkgroup": "WAN"},
		points["usg_wan_uptime_stats/WAN/"].Tags)
	a.Equal(map[string]any{"downtime": 66.0}, points["usg_wan_uptime_stats/WAN/"].Fields)
	a.Equal(map[string]any{"uptime": 1429359.0}, points["usg_wan_uptime_stats/WAN2/"].Fields)
	a.Equal(map[string]any{"availability": 100.0}, points["usg_wan_monitors/WAN/google.com"].Fields,
		"a monitor's 0 latency is skipped")
	a.Equal("true", points["usg_wan_monitors/WAN/ping.ui.com"].Tags["alerting"])
	a.Equal(map[string]any{"availability": 0.0}, points["usg_wan_monitors/WAN/ping.ui.com"].Fields)
}
