package influxunifi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unpoller/pkg/influxunifi"
	"golift.io/cnfg"
)

func TestInfluxVersionDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  influxunifi.Config
		version influxunifi.InfluxVersion
	}{
		{
			name:    "v1 default without token",
			config:  influxunifi.Config{DB: "unifi"},
			version: influxunifi.InfluxV1,
		},
		{
			name:    "v2 default with token",
			config:  influxunifi.Config{AuthToken: "secret"},
			version: influxunifi.InfluxV2,
		},
		{
			name:    "explicit v3",
			config:  influxunifi.Config{Version: 3, AuthToken: "secret", Database: "unifi"},
			version: influxunifi.InfluxV3,
		},
		{
			name:    "explicit v1 with token present",
			config:  influxunifi.Config{Version: 1, AuthToken: "ignored", DB: "unifi"},
			version: influxunifi.InfluxV1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			t.Cleanup(server.Close)

			u := &influxunifi.InfluxUnifi{
				InfluxDB: &influxunifi.InfluxDB{Config: &tc.config},
			}
			u.Config.URL = server.URL

			ok, err := u.DebugOutput()
			require.NoError(t, err)
			assert.True(t, ok)
			assert.Equal(t, tc.version, u.Version)
		})
	}
}

func TestInfluxV3RequiresAuthToken(t *testing.T) {
	t.Parallel()

	u := &influxunifi.InfluxUnifi{
		InfluxDB: &influxunifi.InfluxDB{
			Config: &influxunifi.Config{
				Version:  3,
				URL:      "http://127.0.0.1:8181",
				Database: "unifi",
				Interval: cnfg.Duration{},
			},
		},
	}

	ok, err := u.DebugOutput()
	require.Error(t, err)
	assert.False(t, ok)
	assert.Contains(t, err.Error(), "auth_token")
}
