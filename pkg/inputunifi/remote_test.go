package inputunifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unpoller/unifi/v6"
)

func TestRemoteSitePollNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		sites []unifi.RemoteSite
		want  []string
	}{
		{
			name: "multi-site uses internalReference not display name",
			sites: []unifi.RemoteSite{
				{ID: "uuid-1", InternalReference: "default", Name: "Default"},
				{ID: "uuid-2", InternalReference: "abc1def2", Name: "Office"},
			},
			want: []string{"default", "abc1def2"},
		},
		{
			name: "falls back to Name when internalReference is empty",
			sites: []unifi.RemoteSite{
				{ID: "uuid-1", Name: "default"},
			},
			want: []string{"default"},
		},
		{
			name: "skips sites with neither field",
			sites: []unifi.RemoteSite{
				{ID: "uuid-1"},
				{ID: "uuid-2", InternalReference: "site2", Name: "Two"},
			},
			want: []string{"site2"},
		},
		{
			name:  "empty list",
			sites: nil,
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, remoteSitePollNames(tt.sites))
		})
	}
}

func TestFormatRemoteSites(t *testing.T) {
	t.Parallel()

	sites := []unifi.RemoteSite{
		{ID: "uuid-1", InternalReference: "default", Name: "Default"},
		{ID: "uuid-2", InternalReference: "abc1def2", Name: "Office"},
		{ID: "uuid-3", InternalReference: "only-id"},
		{ID: "uuid-4", Name: "legacy-name-only"},
	}

	assert.Equal(t, []string{
		"default",
		"Office (abc1def2)",
		"only-id",
		"legacy-name-only",
	}, formatRemoteSites(sites))
}
