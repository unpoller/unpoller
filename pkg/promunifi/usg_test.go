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

func testGatewayProm() *promUnifi {
	return &promUnifi{
		Device: descDevice("unpoller_"),
		USG:    descUSG("unpoller_device_"),
	}
}
