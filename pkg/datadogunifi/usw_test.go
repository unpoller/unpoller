//nolint:testpackage // white-box: exercises the unexported batch functions.
package datadogunifi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

func TestBatchUSWUpgradeableFollowsUpgradableKey(t *testing.T) {
	t.Parallel()

	// The controller sends "upgradable", true while a firmware update is waiting.
	const raw = `[
		{"name": "Waiting", "adopted": true, "upgradable": true},
		{"name": "Current", "adopted": true, "upgradable": false}
	]`

	var switches []*unifi.USW

	require.NoError(t, json.Unmarshal([]byte(raw), &switches))

	r := &gaugeReport{report: &Report{Counts: &Counts{Val: map[item]int{}}}, gauges: map[string]recordedGauge{}}
	for _, sw := range switches {
		(&DatadogUnifi{}).batchUSW(r, sw)
	}

	require.Contains(t, r.gauges, "unifi.usw.upgradeable name:Current")
	assert.Equal(t, 1.0, r.gauges["unifi.usw.upgradeable name:Waiting"].value, "a switch with an update waiting reports 1, as USG, UAP and UDB do")
	assert.Equal(t, 0.0, r.gauges["unifi.usw.upgradeable name:Current"].value)
}
