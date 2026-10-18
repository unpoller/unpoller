//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"encoding/json"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

// stationTags maps the station_type of every stations series in sent to its tag label,
// as the descriptor names them.
func stationTags(t *testing.T, p *promUnifi, sent []*metric) map[string]string {
	t.Helper()

	got := make(map[string]string)

	for _, m := range sent {
		if m.Desc != p.Device.Counter {
			continue
		}

		var out dto.Metric

		require.NoError(t, prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, metricFloat(t, m), m.Labels...).Write(&out))

		labels := make(map[string]string)
		for _, l := range out.GetLabel() {
			labels[l.GetName()] = l.GetValue()
		}

		got[labels["station_type"]] = labels["tag"]
	}

	return got
}

func TestExportStationCountsKeepTheirLabels(t *testing.T) {
	t.Parallel()

	p, r := testProm(), &fakeReport{}
	// The exporters pass append(baseLabels, tag): five labels with room to spare.
	labels := append(make([]string, 0, 8), "uap", "Default", "AP", "ctrl", "lobby")

	p.exportSTAcount(r, labels, unifi.FlexInt{Val: 3}, unifi.FlexInt{Val: 2},
		unifi.FlexInt{Val: 4}, unifi.FlexInt{Val: 5}, unifi.FlexInt{Val: 6})
	assert.Equal(t, map[string]string{
		"user": "lobby", "guest": "lobby", "desktop": "lobby", "mobile": "lobby", "handheld": "lobby",
	}, stationTags(t, p, r.sent))
}

func TestExportUBBLinkQualityKeepsItsLabels(t *testing.T) {
	t.Parallel()

	var d unifi.UBB

	require.NoError(t, json.Unmarshal([]byte(`{"adopted": true, "tags": ["lobby"],
		"link_quality": 11, "p2p_stats": {"rx_rate": 7}}`), &d))

	p, r := testProm(), &fakeReport{}
	p.exportUBB(r, &d)
	assert.Equal(t, map[string]string{
		"link_quality": "lobby", "link_quality_current": "lobby", "link_capacity": "lobby",
		"p2p_rx_rate": "lobby", "p2p_tx_rate": "lobby", "p2p_throughput": "lobby",
	}, stationTags(t, p, r.sent))
}

func TestExportUSGStationCountsUseTheirOwnFields(t *testing.T) {
	t.Parallel()

	p, r := testProm(), &fakeReport{}
	p.exportUSG(r, decode[unifi.USG](t, `{"adopted": true, "user-num_sta": 41, "guest-num_sta": 1,
		"num_desktop": 7, "num_mobile": 2, "num_handheld": 8}`))

	got := make(map[string]float64)

	for _, m := range r.sent {
		if m.Desc == p.Device.Counter {
			// exportSTAcount appends the station type last.
			got[m.Labels[len(m.Labels)-1]] = metricFloat(t, m)
		}
	}

	assert.Equal(t, map[string]float64{"user": 41, "guest": 1, "desktop": 7, "mobile": 2, "handheld": 8}, got)
}
