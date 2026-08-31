//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

func testWANConfig() *unifi.WANEnrichedConfiguration {
	return &unifi.WANEnrichedConfiguration{
		SiteName:   "Default (default)",
		SourceName: "https://udr.example",
		Configuration: unifi.WANConfiguration{
			ID:                 "a1",
			Name:               "Internet 1",
			WANNetworkgroup:    "WAN",
			WANType:            "dhcp",
			WANLoadBalanceType: "weighted",
		},
		Details: unifi.WANDetails{
			ServiceProvider: unifi.WANServiceProvider{
				Name: "SpaceX Starlink",
				City: "Brussels",
			},
		},
	}
}

// TestExportWANIsAttributed is the point of this change. Every unpoller_wan_*
// series used to ship with empty site_name and source, which makes the metrics
// of two controllers polled by the same instance indistinguishable. The only
// way to attribute them downstream was to hardcode a mapping in the scrape
// config and hope no second controller ever gained a gateway.
func TestExportWANIsAttributed(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	u := &promUnifi{WAN: descWAN("unifi_")}
	u.exportWAN(r, testWANConfig())

	require.NotEmpty(t, r.sent, "a WAN configuration must produce metrics")

	a := assert.New(t)

	for _, m := range r.sent {
		require.GreaterOrEqual(t, len(m.Labels), 7, "every WAN metric carries the base label set")
		// site_name and source are the last two of the base label set; the
		// provider descriptors substitute isp_name/isp_city earlier but keep
		// the same trailing pair.
		a.Equal("Default (default)", m.Labels[len(m.Labels)-2], "site_name must be populated")
		a.Equal("https://udr.example", m.Labels[len(m.Labels)-1], "source must be populated")
		a.NotContains(m.Labels[:len(m.Labels)-2], "", "no other label should be blank in this fixture")
	}
}

// TestExportWANNilIsSafe pins the existing guard: the collector iterates over
// whatever the input plugin produced, and a nil entry must not panic a poll.
func TestExportWANNilIsSafe(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	u := &promUnifi{WAN: descWAN("unifi_")}
	u.exportWAN(r, nil)

	assert.Empty(t, r.sent, "a nil WAN configuration produces nothing and does not panic")
}
