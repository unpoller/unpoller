//nolint:testpackage // white-box: exercises the unexported export skip path.
package promunifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unpoller/unifi/v6"
)

// Locate/identify only blinks the LED. An adopted device in that state must
// still scrape; dropping it was the silent Grafana gap in #1075.
func TestExportUCILocatingStillExports(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	u := &promUnifi{Device: descDevice("unifi_")}
	u.exportUCI(r, &unifi.UCI{
		Adopted:    unifi.FlexBool{Val: true},
		Locating:   unifi.FlexBool{Val: true},
		Type:       "uci",
		SiteName:   "default",
		Name:       "UCI-1",
		SourceName: "https://controller.example",
	})

	assert.NotEmpty(t, r.sent, "an adopted device in locate mode must still export metrics")
}

func TestExportUCIUnadoptedIsSkipped(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	u := &promUnifi{Device: descDevice("unifi_")}
	u.exportUCI(r, &unifi.UCI{
		Adopted:    unifi.FlexBool{Val: false},
		Locating:   unifi.FlexBool{Val: true},
		Type:       "uci",
		SiteName:   "default",
		Name:       "UCI-1",
		SourceName: "https://controller.example",
	})

	assert.Empty(t, r.sent, "unadopted devices remain excluded")
}
