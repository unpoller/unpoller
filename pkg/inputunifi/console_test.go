package inputunifi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
	"github.com/unpoller/unpoller/pkg/inputunifi"
	"github.com/unpoller/unpoller/pkg/poller"
)

const consoleBody = `{"devices":{"unifiOS":[{"name":"UCG Fiber","firmware":"5.1.33",
	"firmwareStatus":"updateAvailable","updateAvailable":"5.1.34+4d109cd"}]}}`

func consoleInput(t *testing.T, handler http.HandlerFunc) *inputunifi.InputUnifi {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := `{"meta":{"rc":"ok"},"data":[]}`

		switch path := r.URL.Path; {
		case path == unifi.APIConsoleSystemPath:
			handler(w, r)

			return
		case path == "/", path == unifi.APILoginPathNew:
		case path == "/proxy/network/status":
			body = `{"meta":{"server_version":"10.6.109"}}`
		case path == "/proxy/network/api/stat/sites":
			body = `{"data":[{"name":"default","desc":"Main"},{"name":"office","desc":"Office"}]}`
		case strings.HasSuffix(path, "/stat/sysinfo"):
			body = `{"data":[{"hostname":"controller","version":"10.6.109","update_available":false}]}`
		case strings.HasSuffix(path, "/stat/device"), strings.HasSuffix(path, "/stat/sta"):
		default:
			w.WriteHeader(http.StatusNotFound)
		}

		_, err := w.Write([]byte(body))
		assert.NoError(t, err)
	}))
	t.Cleanup(srv.Close)

	disabled := false
	u := &inputunifi.InputUnifi{Config: &inputunifi.Config{
		Controllers: []*inputunifi.Controller{{
			URL: srv.URL, User: "monitor", Pass: "test-password",
			SaveSites: &disabled, SaveDPI: &disabled, SaveRogue: &disabled,
			SaveTraffic: &disabled, SaveSpeedTest: &disabled,
		}},
	}}
	require.NoError(t, u.Initialize(nil))
	require.NotNil(t, u.Controllers[0].Unifi)
	t.Cleanup(u.Controllers[0].Unifi.CloseIdleConnections)

	return u
}

func TestConsoleUpdatesCollectedOncePerController(t *testing.T) {
	t.Parallel()

	var requests atomic.Int32

	u := consoleInput(t, func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)

		_, err := w.Write([]byte(consoleBody))
		assert.NoError(t, err)
	})
	c := u.Controllers[0]
	c.ID = "home"
	metrics, err := u.Metrics(&poller.Filter{Path: c.URL})
	require.NoError(t, err)
	require.Len(t, metrics.ConsoleSystems, 1)
	assert.EqualValues(t, 1, requests.Load(), "two sites must share one console request")
	require.Len(t, metrics.Sysinfos, 2)
	assert.False(t, metrics.Sysinfos[0].(*unifi.Sysinfo).UpdateAvail, "Network update flags retain their meaning")
	system, ok := metrics.ConsoleSystems[0].(*unifi.ConsoleSystem)
	require.True(t, ok)
	assert.Equal(t, "home", system.SourceName)
	assert.Equal(t, "5.1.34+4d109cd", system.Devices.UniFiOS[0].UpdateAvailable)
	assert.True(t, metrics.ControllerStatuses[0].Up)

	// Exercise the generic aggregation used between inputs and outputs as well.
	aggregate := poller.AppendMetrics(&poller.Metrics{}, metrics)
	require.Len(t, aggregate.ConsoleSystems, 1)
	assert.Same(t, system, aggregate.ConsoleSystems[0])
}

func TestConsoleFailureDoesNotFailPollOrKeepOldReading(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		filtered bool
		dynamic  bool
	}{
		{name: "global"},
		{name: "target", filtered: true},
		{name: "target with dynamic enabled", filtered: true, dynamic: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				forbidden atomic.Bool
				requests  atomic.Int32
			)

			u := consoleInput(t, func(w http.ResponseWriter, _ *http.Request) {
				requests.Add(1)

				if forbidden.Load() {
					w.WriteHeader(http.StatusForbidden)

					return
				}

				_, err := w.Write([]byte(consoleBody))
				assert.NoError(t, err)
			})
			u.Dynamic = tc.dynamic

			var filter *poller.Filter

			if tc.filtered {
				filter = &poller.Filter{Path: u.Controllers[0].URL}
			}

			first, err := u.Metrics(filter)
			require.NoError(t, err)
			require.Len(t, first.ConsoleSystems, 1)
			forbidden.Store(true)

			second, err := u.Metrics(filter)
			require.NoError(t, err)
			assert.Empty(t, second.ConsoleSystems, "failed collection must not reuse an earlier update reading")
			assert.Len(t, second.Sysinfos, 2)
			assert.True(t, second.ControllerStatuses[0].Up)
			assert.EqualValues(t, 2, requests.Load(), "a configured target must not be polled again as a dynamic target")
		})
	}
}

func TestConsoleUpdatesSkipUnvalidatedAuthentication(t *testing.T) {
	t.Parallel()

	for _, remote := range []bool{false, true} {
		t.Run(map[bool]string{false: "API key", true: "remote"}[remote], func(t *testing.T) {
			t.Parallel()

			u := consoleInput(t, func(_ http.ResponseWriter, _ *http.Request) {
				t.Error("console endpoint must not be requested")
			})
			if remote {
				u.Controllers[0].Remote = true
			} else {
				u.Controllers[0].APIKey = "test-key"
			}

			m, err := u.Metrics(nil)
			require.NoError(t, err)
			assert.Empty(t, m.ConsoleSystems)
			assert.True(t, m.ControllerStatuses[0].Up)
		})
	}
}

func TestConsoleNameRespectsPIISettings(t *testing.T) {
	t.Parallel()

	for _, drop := range []bool{false, true} {
		t.Run(map[bool]string{false: "hash", true: "drop"}[drop], func(t *testing.T) {
			t.Parallel()

			u := consoleInput(t, func(w http.ResponseWriter, _ *http.Request) {
				_, err := w.Write([]byte(consoleBody))
				assert.NoError(t, err)
			})
			hash := !drop
			u.Controllers[0].HashPII = &hash
			u.Controllers[0].DropPII = &drop
			m, err := u.Metrics(nil)
			require.NoError(t, err)
			require.Len(t, m.ConsoleSystems, 1)
			firmware := m.ConsoleSystems[0].(*unifi.ConsoleSystem).Devices.UniFiOS[0]
			assert.Equal(t, inputunifi.RedactNamePII("UCG Fiber", &hash, &drop), firmware.Name)
			assert.NotEqual(t, "UCG Fiber", firmware.Name)
		})
	}
}
