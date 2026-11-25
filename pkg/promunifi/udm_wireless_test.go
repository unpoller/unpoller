//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

// UniFi Network 11 UDR7: stat has only "wireless", not "ap", but the radio and VAP tables are still top-level.
const network11UDR7WirelessJSON = `{
	"adopted": true,
	"type": "udm",
	"model": "UDMA67A",
	"name": "Gateway",
	"mac": "aa:bb:cc:dd:ee:ff",
	"stat": {
		"wireless": {
			"wifi_tx_retry_packets": 26415,
			"wifi_tx_attempt_packets": 245244
		}
	},
	"radio_table": [
		{"name": "wifi0", "radio": "ng", "ht": 20, "max_txpower": 23, "min_txpower": 1, "nss": 2, "radio_caps": 1}
	],
	"radio_table_stats": [
		{"name": "wifi0", "radio": "ng", "channel": 13, "cu_total": 8, "num_sta": 13, "tx_power": 15, "tx_retries": 9}
	],
	"vap_table": [
		{"name": "wifi0", "bssid": "aa:bb:cc:dd:ee:01", "radio": "ng", "radio_name": "wifi0", "essid": "home",
		 "usage": "user", "up": true, "ccq": 0, "num_satisfaction_sta": 12, "avg_client_signal": -54, "satisfaction": 100}
	]
}`

func TestExportUDMWirelessWithoutStatAp(t *testing.T) {
	t.Parallel()

	var udm unifi.UDM

	require.NoError(t, json.Unmarshal([]byte(network11UDR7WirelessJSON), &udm))
	require.Nil(t, udm.Stat.Ap, "Network 11 omits stat.ap")

	udm.SiteName = "Default"
	udm.SourceName = "https://udr7.example"

	u := testGatewayProm()
	u.UAP = descUAP("unpoller_device_")

	r := &fakeReport{}
	u.exportUDM(r, &udm)

	var vap, radio int
	for _, m := range r.sent {
		name := prometheusDescName(t, m.Desc)
		switch {
		case strings.HasPrefix(name, "unpoller_device_vap_"):
			vap++
		case strings.HasPrefix(name, "unpoller_device_radio_"):
			radio++
		}
	}
	require.NotZero(t, vap, "VAP (WiFi network) metrics must be exported without stat.ap")
	require.NotZero(t, radio, "radio metrics must be exported without stat.ap")
}
