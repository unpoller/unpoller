//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

// The relay and cycle states of an outlet are FlexBool. exportMetrics accepts numbers and plain
// bools only, so sending the FlexBool itself counted every outlet as an error ("not a number")
// and exported neither metric. metricFloat fails the test on any value exportMetrics would reject.
func TestExportPDUOutletTableExportsRelayAndCycleState(t *testing.T) {
	t.Parallel()

	on := unifi.FlexBool{Val: true, Txt: "true"}
	labels := []string{"usw", "Default", "PDU", "https://udm.example", ""}
	table := []unifi.OutletTable{
		{Index: unifi.FlexInt{Val: 1, Txt: "1"}, Name: "Rack left", RelayState: on},
		{Index: unifi.FlexInt{Val: 2, Txt: "2"}, Name: "Rack right", CycleEnabled: on},
	}
	overrides := []unifi.OutletOverride{
		{Index: unifi.FlexInt{Val: 1, Txt: "1"}, Name: "Rack left", RelayState: on},
	}

	r := &fakeReport{}
	u := &promUnifi{PDU: descPDU("unifi_device_")}
	u.exportPDUOutletTable(r, labels, table, overrides)

	relay := map[string]float64{}
	cycle := map[string]float64{}

	for _, m := range r.sent {
		value := metricFloat(t, m)

		switch prometheusDescName(t, m.Desc) {
		case "unifi_device_outlet_relay_state":
			relay[m.Labels[0]] = value
		case "unifi_device_outlet_cycle_enabled":
			cycle[m.Labels[0]] = value
		}
	}

	require.Len(t, r.sent, 2*7+2, "seven metrics per outlet and two per override")

	a := assert.New(t)

	a.Equal(map[string]float64{
		"PDU Outlet 1":          1,
		"PDU Outlet 2":          0,
		"PDU Outlet Override 1": 1,
	}, relay)
	a.Equal(map[string]float64{
		"PDU Outlet 1":          0,
		"PDU Outlet 2":          1,
		"PDU Outlet Override 1": 0,
	}, cycle)
}
