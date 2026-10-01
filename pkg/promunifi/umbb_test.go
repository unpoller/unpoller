//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

// Trimmed from a U5G Max (UMBBE630) on 5G SA. Identifiers and addresses are omitted.
const u5gMaxJSON = `{
	"adopted": true,
	"type": "umbb",
	"model": "UMBBE630",
	"name": "U5G Max",
	"mac": "00:00:5e:00:53:02",
	"version": "7.5.3.19533",
	"uptime": 5264387,
	"mbb": {
		"mode": "failover",
		"state": "ready",
		"sim": [
			{"active": false, "slot": 1, "esim": false, "display_state": "no-sim", "display_state_elapsed": 5264323},
			{
				"active": true, "slot": 2, "esim": true, "spn": "EE",
				"display_state": "operational", "display_state_elapsed": 5230,
				"network_reject_type": "wds", "network_reject_text": "Carrier rejected data", "network_reject_age": 45,
				"rxbytes": "3013150278176", "txbytes": "725650860811"
			}
		],
		"radio": {
			"5g_sa_mode": true, "band": "n78", "rat": "5G", "networkoperator": "EE", "current_slot": 2,
			"has_coverage": true, "roaming": false, "registration_state": 1,
			"rsrp": -96, "rsrq": -13, "snr": 17.600000381469727, "signal": 3, "signal_percent": 95,
			"ca_lte": [],
			"ca_nr": [
				{"band": 78, "dl_arfcn": 637334, "dl_bw_mhz": 40.0, "primary": true, "ul_arfcn": 637334, "ul_bw_mhz": 40.0},
				{"band": 78, "dl_arfcn": 646666, "dl_bw_mhz": 40.0, "primary": false, "ul_arfcn": 0, "ul_bw_mhz": null}
			]
		}
	}
}`

func TestExportUMBB(t *testing.T) {
	t.Parallel()

	var d unifi.UMBB

	require.NoError(t, json.Unmarshal([]byte(u5gMaxJSON), &d))

	d.SiteName = "Default"
	d.SourceName = "https://udm.example"

	r := &fakeReport{}
	u := testProm()
	u.UMBB = descUMBB("unifi_")
	u.exportUMBB(r, &d)

	a := assert.New(t)
	labels := []string{"umbb", "Default", "U5G Max", "https://udm.example", ""}

	a.InDelta(-96, findMetric(t, r.sent, "unifi_mbb_rsrp_dbm", labels).value, 0.001)
	a.InDelta(-13, findMetric(t, r.sent, "unifi_mbb_rsrq_db", labels).value, 0.001)
	a.InDelta(17.6, findMetric(t, r.sent, "unifi_mbb_snr_db", labels).value, 0.001)
	a.InDelta(3, findMetric(t, r.sent, "unifi_mbb_signal_bars", labels).value, 0.001)
	a.InDelta(95, findMetric(t, r.sent, "unifi_mbb_signal_percent", labels).value, 0.001)
	a.InDelta(1, findMetric(t, r.sent, "unifi_mbb_coverage", labels).value, 0.001)
	a.InDelta(2, findMetric(t, r.sent, "unifi_mbb_carriers", append(labels, "nr")).value, 0.001)
	a.InDelta(0, findMetric(t, r.sent, "unifi_mbb_carriers", append(labels, "lte")).value, 0.001)
	a.InDelta(80, findMetric(t, r.sent, "unifi_mbb_bandwidth_mhz", append(labels, "dl")).value, 0.001)
	a.InDelta(40, findMetric(t, r.sent, "unifi_mbb_bandwidth_mhz", append(labels, "ul")).value, 0.001,
		"a null uplink bandwidth adds nothing")
	a.InDelta(1, findMetric(t, r.sent, "unifi_mbb_info",
		append(labels, "ready", "failover", "5G", "n78", "EE", "true", "2")).value, 0.001)

	noSim := append(labels, "1", "false", "")
	a.InDelta(0, findMetric(t, r.sent, "unifi_mbb_sim_active", noSim).value, 0.001)
	a.InDelta(1, findMetric(t, r.sent, "unifi_mbb_sim_info", append(noSim, "no-sim", "", "")).value, 0.001)

	esim := append(labels, "2", "true", "EE")
	a.InDelta(1, findMetric(t, r.sent, "unifi_mbb_sim_active", esim).value, 0.001)
	a.InDelta(3013150278176, findMetric(t, r.sent, "unifi_mbb_sim_receive_bytes_total", esim).value, 1)
	a.InDelta(725650860811, findMetric(t, r.sent, "unifi_mbb_sim_transmit_bytes_total", esim).value, 1)
	a.InDelta(5230, findMetric(t, r.sent, "unifi_mbb_sim_state_seconds", esim).value, 0.001)
	a.InDelta(45, findMetric(t, r.sent, "unifi_mbb_sim_network_reject_age", esim).value, 0.001)
	a.InDelta(1, findMetric(t, r.sent, "unifi_mbb_sim_info",
		append(esim, "operational", "wds", "Carrier rejected data")).value, 0.001)
}

func TestExportUMBBUnadoptedIsSkipped(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	u := testProm()
	u.UMBB = descUMBB("unifi_")
	u.exportUMBB(r, &unifi.UMBB{Name: "U5G Max"})

	assert.Empty(t, r.sent)
}

// findMetric returns the one series with this name and exact label values.
func findMetric(t *testing.T, sent []*metric, name string, labels []string) recordedMetric {
	t.Helper()

	for _, m := range sent {
		if prometheusDescName(t, m.Desc) == name && slices.Equal(m.Labels, labels) {
			return recordedMetric{value: metricFloat(t, m), labels: m.Labels}
		}
	}

	t.Fatalf("no %s series with labels %q", name, labels)

	return recordedMetric{}
}
