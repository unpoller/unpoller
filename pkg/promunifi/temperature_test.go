//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

// series lists every series sent for desc as the values of the named labels, as the descriptor
// names them, joined by ",".
func series(t *testing.T, sent []*metric, desc *prometheus.Desc, names ...string) []string {
	t.Helper()

	var got []string

	for _, m := range sent {
		if m.Desc != desc {
			continue
		}

		var out dto.Metric

		require.NoError(t, prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 0, m.Labels...).Write(&out))

		labels := make(map[string]string)
		for _, l := range out.GetLabel() {
			labels[l.GetName()] = l.GetValue()
		}

		values := make([]string, 0, len(names))
		for _, n := range names {
			values = append(values, labels[n])
		}

		got = append(got, strings.Join(values, ","))
	}

	return got
}

func decode[D any](t *testing.T, payload string) *D {
	t.Helper()

	d := new(D)

	require.NoError(t, json.Unmarshal([]byte(payload), d))

	return d
}

func TestExportTemperatureAndStorageKeepTheirLabels(t *testing.T) {
	t.Parallel()

	const (
		sensors = `"temperatures": [{"name": "Local", "type": "board", "value": 41}, {"name": "CPU", "type": "cpu", "value": 62}]`
		disks   = `"storage": [{"mount_point": "/data", "name": "nvme0", "size": 1000, "used": 400}, {"mount_point": "/boot", "name": "emmc", "size": 200, "used": 50}]`
	)

	p, r := testProm(), &fakeReport{}

	p.exportUXG(r, decode[unifi.UXG](t, `{"name": "uxg", "adopted": true, "tags": ["lobby"], `+sensors+`, `+disks+`}`))
	p.exportUSG(r, decode[unifi.USG](t, `{"name": "usg", "adopted": true, "tags": ["lobby"], `+sensors+`,
		"system-stats": {"temps": {"Board (CPU)": "51 C"}}}`))
	p.exportUSW(r, decode[unifi.USW](t, `{"name": "usw", "adopted": true, "has_temperature": true, "general_temperature": 38}`))
	p.exportUDB(r, decode[unifi.UDB](t, `{"name": "udb", "adopted": true, "tags": ["lobby"], "has_temperature": true, "general_temperature": 38}`))
	p.exportUBB(r, decode[unifi.UBB](t, `{"name": "ubb", "adopted": true, "tags": ["lobby"], "general_temperature": 33}`))
	p.exportUDM(r, decode[unifi.UDM](t, `{"name": "udm", "adopted": true, "tags": ["lobby"], `+sensors+`, `+disks+`}`))

	assert.Equal(t, []string{
		"uxg,Local,board,lobby", "uxg,CPU,cpu,lobby",
		"usg,Local,board,lobby", "usg,CPU,cpu,lobby", "usg,board_cpu,board_cpu,lobby",
		"usw,general,board,", "udb,general,board,lobby", "ubb,ubb,general,lobby",
		"udm,Local,board,lobby", "udm,CPU,cpu,lobby",
	}, series(t, r.sent, p.Device.Temperature, "name", "temp_area", "temp_type", "tag"))
	assert.Equal(t, []string{
		"uxg,/data,nvme0,size,lobby", "uxg,/data,nvme0,used,lobby", "uxg,/boot,emmc,size,lobby", "uxg,/boot,emmc,used,lobby",
		"udm,/data,nvme0,size,lobby", "udm,/data,nvme0,used,lobby", "udm,/boot,emmc,size,lobby", "udm,/boot,emmc,used,lobby",
	}, series(t, r.sent, p.Device.Storage, "name", "mountpoint", "storage_name", "storage_reading", "tag"))
}
