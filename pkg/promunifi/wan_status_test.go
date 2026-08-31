//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

// TestExportWANStatusIsAttributed is the point of this change. Every UniFi OS
// console names its only site "default", so site_name alone cannot tell two
// controllers apart: an instance polling several of them emitted WAN status
// that was indistinguishable, and any downstream attribution had to be guessed.
func TestExportWANStatusIsAttributed(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	u := &promUnifi{WANStatus: descWANStatus("unifi_")}
	u.exportWANStatus(r, &unifi.WANStatus{
		SiteName:   "Default (default)",
		SourceName: "https://udr.example",
		WANInterfaces: []unifi.WANStatusInterface{
			{Name: "Internet 1", State: "ACTIVE", WANNetworkgroup: "WAN"},
			{Name: "Cellular", State: "BACKUP", WANNetworkgroup: "WAN3"},
		},
	})

	require.Len(t, r.sent, 2, "one metric per WAN interface")

	a := assert.New(t)

	for _, m := range r.sent {
		require.Len(t, m.Labels, 5)
		a.Equal("Default (default)", m.Labels[0], "site_name")
		a.Equal("https://udr.example", m.Labels[1], "source must identify the controller")
	}

	a.Equal(float64(1), r.sent[0].Value, "ACTIVE is 1")
	a.Equal(float64(0), r.sent[1].Value, "anything else is 0")
	a.Equal("BACKUP", r.sent[1].Labels[4], "the raw state is kept as a label")
}

// TestExportWANStatusNilIsSafe pins the existing guard: a site with no gateway
// yields a nil status, and that must not panic a poll.
func TestExportWANStatusNilIsSafe(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	u := &promUnifi{WANStatus: descWANStatus("unifi_")}
	u.exportWANStatus(r, nil)

	assert.Empty(t, r.sent)
}
