//nolint:testpackage // white-box: exercises the unexported batch functions.
package datadogunifi

import (
	"encoding/json"
	"slices"
	"strings"
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

// gaugeReport records gauges; other report methods are not used by the WAN batches.
type gaugeReport struct {
	report

	gauges map[string]recordedGauge
}

type recordedGauge struct {
	value float64
	tags  []string
}

func (g *gaugeReport) reportGauge(name string, value float64, tags []string) error {
	g.gauges[name+" "+strings.Join(sortedTags(tags), ",")] = recordedGauge{value: value, tags: tags}

	return nil
}

func sortedTags(tags []string) []string {
	out := slices.Clone(tags)
	slices.Sort(out)

	return out
}

func TestBatchUSGwansOmitsHealthWhenAbsent(t *testing.T) {
	t.Parallel()

	r := &gaugeReport{gauges: map[string]recordedGauge{}}
	(&DatadogUnifi{}).batchUSGwans(r, map[string]string{"name": "Gateway"},
		unifi.Wan{Name: "wan", Up: unifi.FlexBool{Val: true}, Latency: *unifi.NewFlexInt(0), Availability: *unifi.NewFlexInt(0)},
		unifi.Wan{Name: "wan2", Up: unifi.FlexBool{Val: true}})

	a := assert.New(t)
	a.Contains(r.gauges, "unifi.usg.wan_ports.rx_bytes device_name:Gateway,purpose:wan2")
	a.Contains(r.gauges, "unifi.usg.wan_ports.availability device_name:Gateway,purpose:wan")
	a.NotContains(r.gauges, "unifi.usg.wan_ports.latency device_name:Gateway,purpose:wan", "an offline WAN's 0 latency is skipped")
	a.NotContains(r.gauges, "unifi.usg.wan_ports.latency device_name:Gateway,purpose:wan2", "older gateways omit WAN health")
	a.NotContains(r.gauges, "unifi.usg.wan_ports.availability device_name:Gateway,purpose:wan2")
}

func TestBatchUptimeMonitorsWithoutAvailability(t *testing.T) {
	t.Parallel()

	r := &gaugeReport{gauges: map[string]recordedGauge{}}
	(&DatadogUnifi{}).batchUptimeMonitors(r, map[string]string{"wan_networkgroup": "WAN"}, "false", []unifi.UptimeMonitor{
		{Target: "absent.example", Type: "icmp"},
		{Target: "zero.example", Type: "icmp", Availability: *unifi.NewFlexInt(0)},
	})

	got := map[string]float64{}
	for k, g := range r.gauges {
		got[k] = g.value
	}

	assert.Equal(t, map[string]float64{
		"unifi.usg.wan_monitors.availability alerting:false,target:zero.example,type:icmp,wan_networkgroup:WAN": 0,
	}, got, "a missing availability is skipped, 0 is kept")
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

	r := &gaugeReport{gauges: map[string]recordedGauge{}}
	(&DatadogUnifi{}).batchWANHealth(r, map[string]string{"name": "Gateway"}, stats, status)

	got := map[string]float64{}
	for k, g := range r.gauges {
		got[k] = g.value
	}

	assert.Equal(t, map[string]float64{
		"unifi.usg.wan_status.online device_name:Gateway,status:offline,wan_networkgroup:WAN":                                     0,
		"unifi.usg.wan_status.online device_name:Gateway,status:online,wan_networkgroup:WAN2":                                     1,
		"unifi.usg.wan_uptime_stats.downtime device_name:Gateway,wan_networkgroup:WAN":                                            66,
		"unifi.usg.wan_uptime_stats.uptime device_name:Gateway,wan_networkgroup:WAN2":                                             1429359,
		"unifi.usg.wan_monitors.availability alerting:false,device_name:Gateway,target:google.com,type:icmp,wan_networkgroup:WAN": 100,
		"unifi.usg.wan_monitors.availability alerting:true,device_name:Gateway,target:ping.ui.com,type:icmp,wan_networkgroup:WAN": 0,
	}, got, "0 and missing latencies are skipped")
}

func TestBatchSysStatsPrefixesTemperatures(t *testing.T) {
	t.Parallel()

	// The sensor map from the comment on unifi.SystemStats, decoded the way a controller sends it.
	const raw = `{"Board (CPU)":"51 C","Board (PHY)":"51 C","CPU":"72 C","PHY":"77 C"}`

	var temps unifi.TempStatusByName

	require.NoError(t, json.Unmarshal([]byte(raw), &temps))

	stats := (&DatadogUnifi{}).batchSysStats(unifi.SysStats{}, unifi.SystemStats{CPU: unifi.FlexInt{Val: 12}, Temps: temps})

	require.Equal(t, 12.0, stats["cpu"], "the CPU sensor must not replace CPU utilization")
	assert.Equal(t, 72.0, stats["temp_cpu"])
	assert.Equal(t, 77.0, stats["temp_phy"])
	assert.Equal(t, 51.0, stats["temp_board__cpu_"])
	assert.Equal(t, 51.0, stats["temp_board__phy_"])
	assert.NotContains(t, stats, "phy")
	assert.NotContains(t, stats, "board__cpu_")
	assert.NotContains(t, stats, "board__phy_")
}

func TestBatchSysStatsSkipsNullTemperature(t *testing.T) {
	t.Parallel()

	// A sensor reported as null decodes to a nil *FlexTemp; the other sensors must still be written.
	const raw = `{"CPU":null,"PHY":"77 C"}`

	var temps unifi.TempStatusByName

	require.NoError(t, json.Unmarshal([]byte(raw), &temps))

	var stats map[string]float64

	require.NotPanics(t, func() {
		stats = (&DatadogUnifi{}).batchSysStats(unifi.SysStats{}, unifi.SystemStats{CPU: unifi.FlexInt{Val: 12}, Temps: temps})
	})

	assert.Equal(t, 12.0, stats["cpu"])
	assert.Equal(t, 77.0, stats["temp_phy"])
	assert.NotContains(t, stats, "temp_cpu")
}
