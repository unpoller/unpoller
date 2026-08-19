package inputunas_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
	"github.com/unpoller/unpoller/pkg/inputunas"
	"github.com/unpoller/unpoller/pkg/poller"
	"golift.io/cnfg"
	"golift.io/cnfgfile"
)

const (
	deviceInfoJSON = `{"name":"UNAS-Pro","model":"UNASPro","version":"4.0.6","firmwareVersion":"4.0.6",
		"status":"HEALTHY","cpu":{"currentLoad":12,"temperature":45},
		"memory":{"free":1000,"total":8000,"available":6000}}`
	storageJSON = `{"pools":[{"number":1,"id":"pool-1","type":"RAID","status":"HEALTHY",
		"capacity":1000,"usage":400,"activeRaidGroupId":"rg-1",
		"raidGroups":[{"number":1,"id":"rg-1","currentLevel":"RAID5","configLevel":"RAID5",
		"currentProtection":1,"expectedProtection":1,"progress":100}]}],
		"disks":[{"slotId":"1","poolId":"pool-1","type":"HDD","state":"HEALTHY",
		"model":"HAT5300","serial":"ABC123","temperature":38,"healthScore":100,"size":500}]}`
	drivesJSON    = `{"drives":[{"id":"d-1","name":"media","type":"SHARE","status":"HEALTHY","quota":100,"usage":50,"memberCount":3}]}`
	networkIOJSON = `{"receiveKBPS":120,"transmitKBPS":340,"timestamp":"2026-01-01T00:00:00Z"}`
)

// unasServer is a fake UNAS Pro console. failData, when set, makes every data endpoint
// return 401 -- which is how an expired session presents itself.
type unasServer struct {
	logins   atomic.Int64
	failData atomic.Bool
}

func (s *unasServer) start(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, _ *http.Request) {
		s.logins.Add(1)
		w.Header().Set("x-csrf-token", "csrf-token")
		w.WriteHeader(http.StatusOK)
	})

	for path, body := range map[string]string{
		unifi.APIUNASDeviceInfoPath: deviceInfoJSON,
		unifi.APIUNASStoragePath:    storageJSON,
		unifi.APIUNASDrivesPath:     drivesJSON,
		unifi.APIUNASNetworkIOPath:  networkIOJSON,
	} {
		mux.HandleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
			if s.failData.Load() {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(body))
		})
	}

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

func newInput(t *testing.T, urls ...string) *inputunas.InputUNAS {
	t.Helper()

	devices := make([]*inputunas.Device, len(urls))
	for i, u := range urls {
		devices[i] = &inputunas.Device{URL: u, User: "unpoller", Pass: "secret"}
	}

	return &inputunas.InputUNAS{Config: &inputunas.Config{Enable: true, Devices: devices}}
}

func TestMetricsCollectsConsole(t *testing.T) {
	t.Parallel()

	srv := (&unasServer{}).start(t)
	u := newInput(t, srv.URL)
	require.NoError(t, u.Initialize(nil))

	m, err := u.Metrics(nil)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Len(t, m.UNASDevices, 1)

	d, ok := m.UNASDevices[0].(*unifi.UNASDevice)
	require.True(t, ok, "outputs type-assert to *unifi.UNASDevice, so nothing else may be appended")

	a := assert.New(t)
	a.Equal("UNAS-Pro", d.Name())
	a.Equal("UNASPro", d.Model())
	a.Equal(srv.URL, d.SourceName)
	require.NotNil(t, d.Storage)
	a.Len(d.Storage.Pools, 1)
	a.Len(d.Storage.Disks, 1)
	a.Len(d.Drives, 1)
	require.NotNil(t, d.NetworkIO)
	a.InDelta(120, d.NetworkIO.ReceiveKBPS.Val, 0.01)
}

// A UNAS session cookie expires after roughly two hours. When it does, every endpoint fails
// at once, and the plugin must log back in and re-attach fresh cookies rather than reuse the
// dead jar -- otherwise a long-running poller reports nothing until it is restarted.
func TestMetricsReauthenticatesAfterSessionExpiry(t *testing.T) {
	t.Parallel()

	fake := &unasServer{}
	srv := fake.start(t)
	u := newInput(t, srv.URL)
	require.NoError(t, u.Initialize(nil))
	require.Equal(t, int64(1), fake.logins.Load())

	fake.failData.Store(true)

	m, err := u.Metrics(nil)
	require.Error(t, err, "a console that fails every endpoint even after re-auth is a real error")
	require.Nil(t, m)
	require.Equal(t, int64(2), fake.logins.Load(), "a total failure must trigger exactly one re-login")

	// Once the console answers again, the refreshed session works without another login.
	fake.failData.Store(false)

	m, err = u.Metrics(nil)
	require.NoError(t, err)
	require.Len(t, m.UNASDevices, 1)
	require.Equal(t, int64(2), fake.logins.Load(), "a healthy poll must not re-authenticate")
}

// One dead console must not discard the metrics of a healthy one. poller.collectMetrics
// drops the result entirely when an input returns both metrics and an error, so returning
// them together here would throw away everything the working console reported.
func TestMetricsPartialFailureStillReportsHealthyConsole(t *testing.T) {
	t.Parallel()

	good := (&unasServer{}).start(t)
	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer dead.Close()

	u := newInput(t, dead.URL, good.URL)
	require.NoError(t, u.Initialize(nil))

	m, err := u.Metrics(nil)
	require.NoError(t, err, "one bad console out of two is not an error")
	require.Len(t, m.UNASDevices, 1)
}

// enable defaults to false, so a config that names consoles but never sets it must not poll
// them. Configuring the devices here is the point: an empty config would pass even if the
// switch were ignored entirely.
func TestNotEnabled(t *testing.T) {
	t.Parallel()

	srv := (&unasServer{}).start(t)
	u := &inputunas.InputUNAS{Config: &inputunas.Config{
		Devices: []*inputunas.Device{{URL: srv.URL, User: "unpoller", Pass: "secret"}},
	}}
	require.NoError(t, u.Initialize(nil))

	m, err := u.Metrics(nil)
	require.NoError(t, err)
	require.Nil(t, m, "enable is false, so nothing is polled")

	ok, err := u.DebugInput()
	require.NoError(t, err)
	require.True(t, ok, "--debugio must not reach a console the operator has not enabled")
}

// With no [unas] block the plugin must be silent and inert: that is what makes UNAS support
// opt-in. It must not synthesize a default device the way inputunifi does with localhost.
func TestUnconfiguredIsInert(t *testing.T) {
	t.Parallel()

	u := &inputunas.InputUNAS{}
	require.NoError(t, u.Initialize(nil))

	m, err := u.Metrics(nil)
	require.NoError(t, err)
	require.Nil(t, m, "an unconfigured plugin stays off")

	ok, err := u.DebugInput()
	require.NoError(t, err)
	require.True(t, ok)
}

func TestDeviceWithNoURLIsIgnored(t *testing.T) {
	t.Parallel()

	u := &inputunas.InputUNAS{Config: &inputunas.Config{Enable: true, Devices: []*inputunas.Device{{}, nil}}}
	require.NoError(t, u.Initialize(nil))

	m, err := u.Metrics(nil)
	require.NoError(t, err)
	require.Empty(t, m.UNASDevices)
}

func TestEventsAreEmpty(t *testing.T) {
	t.Parallel()

	e, err := (&inputunas.InputUNAS{}).Events(&poller.Filter{})
	require.NoError(t, err)
	require.NotNil(t, e)
	require.Empty(t, e.Logs)
}

func TestRawMetrics(t *testing.T) {
	t.Parallel()

	srv := (&unasServer{}).start(t)
	u := newInput(t, srv.URL)
	require.NoError(t, u.Initialize(nil))

	body, err := u.RawMetrics(&poller.Filter{Kind: "storage"})
	require.NoError(t, err)
	require.Contains(t, string(body), "pool-1")

	_, err = u.RawMetrics(&poller.Filter{Unit: 9})
	require.ErrorIs(t, err, inputunas.ErrNoDevices)
}

// The enable switch has to survive four independent binding paths -- toml, json, yaml and the
// UP_ environment -- each driven by its own struct tag. A typo in any one tag leaves the
// plugin silently off for users of that format, which is the same class of invisible failure
// as a missing tag anywhere else in this config. Note the env name derives from the xml tag,
// not the json one.
func TestEnableBindsInEveryConfigFormat(t *testing.T) {
	t.Parallel()

	write := func(t *testing.T, name, body string) string {
		t.Helper()

		path := filepath.Join(t.TempDir(), name)
		require.NoError(t, os.WriteFile(path, []byte(body), 0o600))

		return path
	}

	for _, tc := range []struct {
		name string
		file string
		body string
	}{
		{"toml", "up.conf", "[unas]\n  enable = true\n[[unas.device]]\n  url = \"https://nas\"\n"},
		{"json", "up.json", `{"unas":{"enable":true,"devices":[{"url":"https://nas"}]}}`},
		{"yaml", "up.yaml", "unas:\n  enable: true\n  devices:\n    - url: \"https://nas\"\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := &inputunas.InputUNAS{Config: &inputunas.Config{}}
			require.NoError(t, cnfgfile.Unmarshal(input, write(t, tc.file, tc.body)))
			assert.True(t, input.Enable, "enable did not bind from %s", tc.name)
			require.Len(t, input.Devices, 1)
			assert.Equal(t, "https://nas", input.Devices[0].URL)
		})
	}
}

// Not parallel: t.Setenv is incompatible with a parallel test or parent.
//
// The env variable name comes from the xml tag, not the json one, which is why this is worth
// asserting rather than assuming from the toml/json/yaml cases above.
func TestEnableBindsFromEnvironment(t *testing.T) {
	t.Setenv("UP_UNAS_ENABLE", "true")

	input := &inputunas.InputUNAS{Config: &inputunas.Config{}}
	_, err := cnfg.UnmarshalENV(input, "UP")
	require.NoError(t, err)
	assert.True(t, input.Enable, "UP_UNAS_ENABLE did not bind")
}

// The shipped examples must all default to off, so that installing unpoller never starts
// polling a storage console nobody asked for.
func TestShippedExamplesDoNotEnableUNAS(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"up.conf.example", "up.json.example", "up.yaml.example"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			input := &inputunas.InputUNAS{Config: &inputunas.Config{}}
			require.NoError(t, cnfgfile.Unmarshal(input, filepath.Join("..", "..", "examples", name)))
			assert.False(t, input.Enable, "%s ships with UNAS enabled", name)
		})
	}
}
