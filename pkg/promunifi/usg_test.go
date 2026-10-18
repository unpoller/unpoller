//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

// Network 11 omits the stat object. Uplink and speedtest-status remain.
const network11GatewayJSON = `{
	"adopted": true,
	"type": "udm",
	"model": "UDRULT",
	"name": "Gateway",
	"uplink": {
		"name": "wan",
		"latency": 17000,
		"speed": 1000,
		"max_speed": 2500,
		"uptime": 3600
	},
	"speedtest-status": {
		"latency": 15000,
		"runtime": 9,
		"rundate": 1791344670,
		"xput_download": 646,
		"xput_upload": 41,
		"source_interface": "wan"
	}
}`

func TestExportUSGstatsWithoutStatKeepsUplinkAndSpeedtest(t *testing.T) {
	t.Parallel()

	var udm unifi.UDM

	require.NoError(t, json.Unmarshal([]byte(network11GatewayJSON), &udm))
	require.Nil(t, udm.Stat.Gw, "Network 11 omits stat.gw")

	udm.SiteName = "Default"
	udm.SourceName = "https://ucg.example"

	r := &fakeReport{}
	testGatewayProm().exportUDM(r, &udm)

	got := metricsByName(t, r.sent)
	a := assert.New(t)

	uplink, ok := got["unpoller_device_uplink_latency_seconds"]
	require.True(t, ok, "uplink latency must still be exported")
	a.InDelta(17, uplink.value, 0.0001)
	a.Equal([]string{"wan", "Default", "Gateway", "https://ucg.example", ""}, uplink.labels)

	speed, ok := got["unpoller_device_uplink_speed_mbps"]
	require.True(t, ok)
	a.Equal(1000.0, speed.value)

	maxSpeed, ok := got["unpoller_device_uplink_max_speed_mbps"]
	require.True(t, ok)
	a.Equal(2500.0, maxSpeed.value)

	uptime, ok := got["unpoller_device_uplink_uptime_seconds"]
	require.True(t, ok)
	a.Equal(3600.0, uptime.value)

	latency, ok := got["unpoller_device_speedtest_latency_seconds"]
	require.True(t, ok)
	a.InDelta(15, latency.value, 0.0001)

	runtime, ok := got["unpoller_device_speedtest_runtime_seconds"]
	require.True(t, ok)
	a.Equal(9.0, runtime.value)

	rundate, ok := got["unpoller_device_speedtest_rundate_seconds"]
	require.True(t, ok)
	a.Equal(1791344670.0, rundate.value)

	download, ok := got["unpoller_device_speedtest_download"]
	require.True(t, ok)
	a.Equal(646.0, download.value)

	upload, ok := got["unpoller_device_speedtest_upload"]
	require.True(t, ok)
	a.Equal(41.0, upload.value)

	_, hasLAN := got["unpoller_device_lan_receive_bytes_total"]
	a.False(hasLAN, "LAN counters come from stat.gw and stay absent")
}

func TestExportUSGstatsKeepsLANWhenStatPresent(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	labels := []string{"usg", "Default", "Gateway", "ctrl", ""}
	gw := &unifi.Gw{LanRxBytes: unifi.FlexInt{Val: 12345}}

	testGatewayProm().exportUSGstats(r, labels, gw, unifi.SpeedtestStatus{}, unifi.Uplink{})

	got := metricsByName(t, r.sent)
	lan, ok := got["unpoller_device_lan_receive_bytes_total"]
	require.True(t, ok)
	assert.Equal(t, 12345.0, lan.value)
	assert.Equal(t, []string{"lan", "Default", "Gateway", "ctrl", ""}, lan.labels)

	_, hasUplink := got["unpoller_device_uplink_latency_seconds"]
	assert.True(t, hasUplink)
}

func TestExportUSGWithoutStatKeepsUplink(t *testing.T) {
	t.Parallel()

	raw := strings.Replace(network11GatewayJSON, `"type": "udm"`, `"type": "usg"`, 1)

	var usg unifi.USG

	require.NoError(t, json.Unmarshal([]byte(raw), &usg))
	require.Nil(t, usg.Stat.Gw)

	usg.SiteName = "Default"
	usg.SourceName = "https://usg.example"

	r := &fakeReport{}
	testGatewayProm().exportUSG(r, &usg)

	got := metricsByName(t, r.sent)
	uplink, ok := got["unpoller_device_uplink_latency_seconds"]
	require.True(t, ok)
	assert.InDelta(t, 17, uplink.value, 0.0001)
	assert.Equal(t, []string{"wan", "Default", "Gateway", "https://usg.example", ""}, uplink.labels)
}

// Trimmed from a UDM Pro (5.1.33, Network 10.6.106) a minute into a real WAN1
// outage: the port link stayed up, and the upstream went dark.
const udmWANOfflineJSON = `{
	"adopted": true,
	"type": "udm",
	"model": "UDMPRO",
	"name": "Gateway",
	"wan1": {"name": "eth9", "ifname": "eth9", "up": true, "latency": 0, "availability": 0},
	"wan2": {"name": "eth8", "ifname": "eth8", "up": true, "latency": 22, "availability": 100},
	"last_wan_status": {"WAN": "offline", "WAN2": "online"},
	"uptime_stats": {
		"WAN": {
			"downtime": 66,
			"monitors": [
				{"target": "www.microsoft.com", "type": "icmp", "availability": 100.0},
				{"target": "google.com", "type": "icmp", "availability": 100.0, "latency_average": 0}
			],
			"alerting_monitors": [{"target": "ping.ui.com", "type": "icmp", "availability": 0.0, "latency_average": 0}]
		},
		"WAN2": {
			"availability": 100.0,
			"latency_average": 21,
			"time_period": 86400,
			"uptime": 1429359,
			"monitors": [{"target": "www.microsoft.com", "type": "icmp", "availability": 100.0, "latency_average": 28}]
		}
	}
}`

func TestExportUDMWANHealthDuringOutage(t *testing.T) {
	t.Parallel()

	var udm unifi.UDM

	require.NoError(t, json.Unmarshal([]byte(udmWANOfflineJSON), &udm))

	udm.SiteName = "Default"
	udm.SourceName = "https://udm.example"

	r := &fakeReport{}
	testGatewayProm().exportUDM(r, &udm)

	const (
		eth9 = "eth9,Default,Gateway,https://udm.example,"
		eth8 = "eth8,Default,Gateway,https://udm.example,"
		wan  = "WAN,Default,Gateway,https://udm.example,"
		wan2 = "WAN2,Default,Gateway,https://udm.example,"
	)

	a := assert.New(t)
	a.Equal(map[string]float64{eth9: 0, eth8: 100},
		seriesByLabels(t, r.sent, "unpoller_device_wan_availability_percent"))
	a.Equal(map[string]float64{eth8: 0.022},
		seriesByLabels(t, r.sent, "unpoller_device_wan_latency_seconds"), "an offline WAN's 0 latency is skipped")
	a.Equal(map[string]float64{wan + ",offline": 0, wan2 + ",online": 1},
		seriesByLabels(t, r.sent, "unpoller_device_wan_status"))

	a.Empty(seriesByLabels(t, r.sent, "unpoller_device_wan_uptime_stats_latency_average_seconds"))
	a.Equal(map[string]float64{wan2: 1429359},
		seriesByLabels(t, r.sent, "unpoller_device_wan_uptime_stats_uptime_seconds"))
	a.Equal(map[string]float64{wan: 66},
		seriesByLabels(t, r.sent, "unpoller_device_wan_uptime_stats_downtime_seconds"))

	a.Equal(map[string]float64{
		wan + ",www.microsoft.com,icmp,false":  100,
		wan + ",google.com,icmp,false":         100,
		wan + ",ping.ui.com,icmp,true":         0,
		wan2 + ",www.microsoft.com,icmp,false": 100,
	}, seriesByLabels(t, r.sent, "unpoller_device_wan_monitor_availability_percent"))
	a.Equal(map[string]float64{wan2 + ",www.microsoft.com,icmp,false": 0.028},
		seriesByLabels(t, r.sent, "unpoller_device_wan_monitor_latency_average_seconds"),
		"monitors on the offline WAN report 0 or no latency")
}

// Trimmed from the same UDM Pro about 25 seconds after WAN1 came back: the
// uptime had restarted, and a monitor that hasn't answered omits latency.
func TestExportWANHealthAfterRecovery(t *testing.T) {
	t.Parallel()

	const raw = `{
		"availability": 100.0,
		"latency_average": 27,
		"time_period": 10,
		"uptime": 26,
		"monitors": [{"target": "google.com", "type": "icmp", "availability": 0.0}]
	}`

	var stats unifi.UptimeStats

	require.NoError(t, json.Unmarshal([]byte(raw), &stats))

	r := &fakeReport{}
	labels := []string{"udm", "Default", "Gateway", "ctrl", ""}
	testGatewayProm().exportWANHealth(r, labels, map[string]*unifi.UptimeStats{"WAN": &stats}, nil)

	a := assert.New(t)
	a.Equal(map[string]float64{"WAN,Default,Gateway,ctrl,": 26},
		seriesByLabels(t, r.sent, "unpoller_device_wan_uptime_stats_uptime_seconds"))
	a.Empty(seriesByLabels(t, r.sent, "unpoller_device_wan_uptime_stats_downtime_seconds"))
	a.Equal(map[string]float64{"WAN,Default,Gateway,ctrl,,google.com,icmp,false": 0},
		seriesByLabels(t, r.sent, "unpoller_device_wan_monitor_availability_percent"))
	a.Empty(seriesByLabels(t, r.sent, "unpoller_device_wan_monitor_latency_average_seconds"))
}

func TestExportUXGWANHealth(t *testing.T) {
	t.Parallel()

	raw := strings.Replace(udmWANOfflineJSON, `"type": "udm"`, `"type": "uxg"`, 1)

	var uxg unifi.UXG

	require.NoError(t, json.Unmarshal([]byte(raw), &uxg))

	uxg.SiteName = "Default"
	uxg.SourceName = "https://uxg.example"

	r := &fakeReport{}
	testGatewayProm().exportUXG(r, &uxg)

	assert.Len(t, seriesByLabels(t, r.sent, "unpoller_device_wan_status"), 2)
	assert.Len(t, seriesByLabels(t, r.sent, "unpoller_device_wan_uptime_stats_downtime_seconds"), 1)
	assert.Len(t, seriesByLabels(t, r.sent, "unpoller_device_wan_latency_seconds"), 1)
}

// A USG reports neither WAN latency nor availability; don't export them as zero.
func TestExportWANPortsWithoutHealthFields(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	labels := []string{"usg", "Default", "Gateway", "ctrl", ""}
	wan := unifi.Wan{Name: "wan", Up: unifi.FlexBool{Val: true}, RxBytes: unifi.FlexInt{Val: 1}}

	testGatewayProm().exportWANPorts(r, labels, wan)

	a := assert.New(t)
	a.NotEmpty(seriesByLabels(t, r.sent, "unpoller_device_wan_receive_bytes_total"))
	a.Empty(seriesByLabels(t, r.sent, "unpoller_device_wan_latency_seconds"))
	a.Empty(seriesByLabels(t, r.sent, "unpoller_device_wan_availability_percent"))
}

func TestExportUptimeMonitorsWithoutAvailability(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	labels := []string{"WAN", "Default", "Gateway", "ctrl", ""}
	monitors := []unifi.UptimeMonitor{
		{Target: "absent.example", Type: "icmp"},
		{Target: "zero.example", Type: "icmp", Availability: *unifi.NewFlexInt(0)},
	}

	testGatewayProm().exportUptimeMonitors(r, labels, "false", monitors)

	a := assert.New(t)
	a.Equal(map[string]float64{"WAN,Default,Gateway,ctrl,,zero.example,icmp,false": 0},
		seriesByLabels(t, r.sent, "unpoller_device_wan_monitor_availability_percent"), "a missing availability is skipped, 0 is kept")
	a.Empty(seriesByLabels(t, r.sent, "unpoller_device_wan_monitor_latency_average_seconds"))
}

// seriesByLabels returns the value of each series named name, keyed by its joined labels.
func seriesByLabels(t *testing.T, sent []*metric, name string) map[string]float64 {
	t.Helper()

	out := map[string]float64{}

	for _, m := range sent {
		if prometheusDescName(t, m.Desc) == name {
			out[strings.Join(m.Labels, ",")] = metricFloat(t, m)
		}
	}

	return out
}

func testGatewayProm() *promUnifi {
	return &promUnifi{
		Device: descDevice("unpoller_"),
		USG:    descUSG("unpoller_device_"),
	}
}
