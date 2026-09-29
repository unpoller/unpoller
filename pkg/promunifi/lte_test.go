//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

// Captured from a U-LTE-Pro (model ULTEPEU) while the module was in
// "LTE Failover Ready" / SIM standby. IMEI, ICCID, and cell identity are omitted.
const ulteStandbyJSON = `{
	"adopted": true,
	"type": "uap",
	"model": "ULTEPEU",
	"name": "U-LTE-Pro",
	"mac": "aa:bb:cc:dd:ee:ff",
	"lte_state": "ready",
	"lte_failover": false,
	"lte_failover_mode": "failover",
	"lte_connected": "yes",
	"lte_networkoperator": "o2 - de",
	"lte_signal": "Good signal strength (3)",
	"lte_rssi": "-58",
	"lte_rsrp": "-95",
	"lte_rsrq": "-16",
	"lte_rat": "LTE",
	"lte_mode": "LTE",
	"lte_band": "eutran-3",
	"lte_rx_chan": "1600",
	"lte_tx_chan": "19600",
	"lte_pdptype": "IPV4"
}`

func testProm() *promUnifi {
	return &promUnifi{
		Device: descDevice("unifi_"),
		UAP:    descUAP("unifi_"),
		USG:    descUSG("unifi_"),
	}
}

func unmarshalULTE(t *testing.T, payload string) *unifi.UAP {
	t.Helper()

	var d unifi.UAP

	require.NoError(t, json.Unmarshal([]byte(payload), &d))

	d.SiteName = "Default"
	d.SourceName = "https://udm.example"

	return &d
}

func TestExportULTEStandbyStatus(t *testing.T) {
	t.Parallel()

	d := unmarshalULTE(t, ulteStandbyJSON)
	r := &fakeReport{}
	testProm().exportUAP(r, d)

	got := metricsByName(t, r.sent)
	a := assert.New(t)

	a.Equal(1.0, got["unifi_lte_connected"].value)
	a.Equal(0.0, got["unifi_lte_failover"].value)
	a.InDelta(-58, got["unifi_lte_rssi_dbm"].value, 0.001)
	a.InDelta(-95, got["unifi_lte_rsrp_dbm"].value, 0.001)
	a.InDelta(-16, got["unifi_lte_rsrq_db"].value, 0.001)
	a.InDelta(1600, got["unifi_lte_rx_channel"].value, 0.001)
	a.InDelta(19600, got["unifi_lte_tx_channel"].value, 0.001)
	a.Equal(1.0, got["unifi_lte_info"].value)
	a.Equal([]string{
		"uap", "Default", "U-LTE-Pro", "https://udm.example", "",
		"ready", "failover", "Good signal strength (3)", "LTE", "LTE", "eutran-3", "IPV4", "o2 - de",
	}, got["unifi_lte_info"].labels)
}

func TestExportULTEFailoverActive(t *testing.T) {
	t.Parallel()

	d := unmarshalULTE(t, ulteStandbyJSON)
	require.NoError(t, json.Unmarshal([]byte(`true`), &d.LteFailover))

	r := &fakeReport{}
	testProm().exportUAP(r, d)

	got := metricsByName(t, r.sent)
	assert.Equal(t, 1.0, got["unifi_lte_failover"].value)
	assert.Equal(t, 1.0, got["unifi_lte_connected"].value)
}

func TestExportUAPWithoutLTEOmitsModemMetrics(t *testing.T) {
	t.Parallel()

	d := &unifi.UAP{
		Adopted:    unifi.FlexBool{Val: true, Txt: "true"},
		Type:       "uap",
		Name:       "Hallway",
		SiteName:   "Default",
		SourceName: "https://udm.example",
	}
	r := &fakeReport{}
	testProm().exportUAP(r, d)

	for name := range metricsByName(t, r.sent) {
		assert.NotContains(t, name, "lte", "a wireless AP must not emit modem metrics")
	}
}

type recordedMetric struct {
	value  float64
	labels []string
}

func metricsByName(t *testing.T, sent []*metric) map[string]recordedMetric {
	t.Helper()

	out := make(map[string]recordedMetric, len(sent))

	for _, m := range sent {
		name := prometheusDescName(t, m.Desc)
		out[name] = recordedMetric{value: metricFloat(t, m), labels: m.Labels}
	}

	return out
}

func metricFloat(t *testing.T, m *metric) float64 {
	t.Helper()

	switch v := m.Value.(type) {
	case float64:
		return v
	case unifi.FlexInt:
		return v.Val
	case bool:
		if v {
			return 1
		}

		return 0
	default:
		t.Fatalf("metric value is %T", m.Value)

		return 0
	}
}
