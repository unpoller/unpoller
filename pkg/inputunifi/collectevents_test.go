package inputunifi // nolint: testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOverrideSiteName is what this change exists for. Every UniFi OS console
// calls its only site "default", so a poller watching several of them ships log
// entries that cannot be told apart — which is what default_site_name_override
// is meant to solve, and did not, on this path.
func TestOverrideSiteName(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		override string
		siteName string
		want     string
	}{
		"stock name from a UniFi OS console": {
			override: "Les_Solidarites",
			siteName: "Default (default)",
			want:     "Les_Solidarites",
		},
		"bare default": {
			override: "Les_Solidarites",
			siteName: "default",
			want:     "Les_Solidarites",
		},
		"a site the operator already named is left alone": {
			override: "Les_Solidarites",
			siteName: "CharlHot - SweetHome",
			want:     "CharlHot - SweetHome",
		},
		"no override configured is a no-op": {
			override: "",
			siteName: "Default (default)",
			want:     "Default (default)",
		},
		"empty site name is not a default": {
			override: "Les_Solidarites",
			siteName: "",
			want:     "",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c := &Controller{DefaultSiteNameOverride: tc.override}
			assert.Equal(t, tc.want, overrideSiteName(c, tc.siteName))
		})
	}
}
