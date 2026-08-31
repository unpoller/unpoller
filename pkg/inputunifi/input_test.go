package inputunifi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
	"github.com/unpoller/unpoller/pkg/inputunifi"
	"github.com/unpoller/unpoller/pkg/poller"
	"golift.io/cnfg"
	"golift.io/cnfgfile"
)

const (
	metaInfoJSON     = `{"applicationVersion":"7.2.105"}`
	camerasJSON      = `[{"id":"cam-1","modelKey":"camera","state":"CONNECTED","name":"Driveway","mac":"001122334455"}]`
	sensorsJSON      = `[{"id":"sen-1","modelKey":"sensor","state":"CONNECTED","name":"Garage","mac":"001122334466"}]`
	nvrJSON          = `{"id":"nvr-1","modelKey":"nvr","state":"CONNECTED","name":"UNVR4","mac":"001122334477"}`
	emptyListJSON    = `[]`
	protectLogsJSON  = `{"items":[{"id":"evt-1","modelKey":"event","type":"motion","camera":"cam-1","timestamp":1735689600000}]}`
	unifiOSSPA       = `<!doctype html><html lang="en"><body></body></html>`
	protectAPIKeyStr = "protect-key-abc"
)

// unvr is a fake Protect-only console. It answers the Protect Integration paths and the
// UniFi OS login, and serves the console's own SPA HTML for everything else -- which is
// exactly what a UNVR does with /proxy/network/*, there being no Network application to
// route to. That HTML is what breaks NewUnifi in unpoller/unpoller#1066.
type unvr struct {
	mu        sync.Mutex
	requested []string
}

func (s *unvr) start(t *testing.T) *httptest.Server {
	t.Helper()

	record := func(r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.requested = append(s.requested, r.URL.Path)
	}

	mux := http.NewServeMux()

	for path, body := range map[string]string{
		unifi.APIProtectMetaInfoPath:     metaInfoJSON,
		unifi.APIProtectCamerasPath:      camerasJSON,
		unifi.APIProtectSensorsPath:      sensorsJSON,
		unifi.APIProtectNVRPath:          nvrJSON,
		unifi.APIProtectLightsPath:       emptyListJSON,
		unifi.APIProtectBridgesPath:      emptyListJSON,
		unifi.APIProtectLinkStationsPath: emptyListJSON,
		unifi.APIProtectLogPath:          protectLogsJSON,
	} {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			record(r)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(body))
		})
	}

	mux.HandleFunc(unifi.APILoginPathNew, func(w http.ResponseWriter, r *http.Request) {
		record(r)
		w.Header().Set("x-csrf-token", "csrf-token")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		record(r)
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(unifiOSSPA))
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// paths returns every path requested so far, and networkPaths only those under /proxy/network.
func (s *unvr) paths() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]string(nil), s.requested...)
}

func (s *unvr) networkPaths() []string {
	var network []string

	for _, p := range s.paths() {
		if strings.HasPrefix(p, "/proxy/network") || p == unifi.APIStatusPath {
			network = append(network, p)
		}
	}

	return network
}

// newProtectOnlyInput builds an input with a single Protect-only controller. Options are
// applied before defaults are filled in, so a test can flip any flag it needs.
func newProtectOnlyInput(url string, opts ...func(*inputunifi.Controller)) *inputunifi.InputUnifi {
	enabled := true
	c := &inputunifi.Controller{
		URL:                url,
		User:               "ro-user",
		Pass:               "secret",
		ProtectAPIKey:      protectAPIKeyStr,
		DisableNetwork:     &enabled,
		SaveProtectDevices: &enabled,
	}

	for _, opt := range opts {
		opt(c)
	}

	return &inputunifi.InputUnifi{Config: &inputunifi.Config{Controllers: []*inputunifi.Controller{c}}}
}

// TestProtectOnlyControllerInitializes is the headline of unpoller/unpoller#1066: a UNVR is
// configured like any other controller and must come up. Before disable_network existed it
// died in NewUnifi and never even printed its config summary.
func TestProtectOnlyControllerInitializes(t *testing.T) {
	t.Parallel()

	fake := &unvr{}
	srv := fake.start(t)

	u := newProtectOnlyInput(srv.URL)
	require.NoError(t, u.Initialize(nil))

	require.Len(t, u.Controllers, 1)
	c := u.Controllers[0]

	a := assert.New(t)
	require.NotNil(t, c.Unifi, "controller client is nil, so initialization failed")
	// NewProtectClient reports the Protect application version, there being no Network one.
	a.Equal("7.2.105", c.Unifi.ServerVersion)
	// Site discovery is skipped, so /proxy/network is never touched at all.
	a.Empty(fake.networkPaths())
}

// The control case: the same console without disable_network still fails, so this test
// pins that the flag -- not some unrelated change -- is what makes the difference.
func TestProtectOnlyControllerFailsWithoutFlag(t *testing.T) {
	t.Parallel()

	srv := (&unvr{}).start(t)

	disabled := false
	u := newProtectOnlyInput(srv.URL, func(c *inputunifi.Controller) { c.DisableNetwork = &disabled })
	require.NoError(t, u.Initialize(nil))

	assert.Nil(t, u.Controllers[0].Unifi, "NewUnifi should not survive a console with no Network application")
}

func TestProtectOnlyControllerCollectsMetrics(t *testing.T) {
	t.Parallel()

	fake := &unvr{}
	srv := fake.start(t)

	u := newProtectOnlyInput(srv.URL)
	require.NoError(t, u.Initialize(nil))

	m, err := u.Metrics(nil)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Len(t, m.ProtectDevices, 1)

	pd, ok := m.ProtectDevices[0].(*unifi.ProtectDevices)
	require.True(t, ok, "outputs type-assert to *unifi.ProtectDevices, so nothing else may be appended")

	a := assert.New(t)
	a.Equal("7.2.105", pd.Version)
	a.Len(pd.Cameras, 1)
	a.Len(pd.Sensors, 1)
	require.NotNil(t, pd.NVR)
	a.Equal("UNVR4", pd.NVR.Name)

	// Nothing Network-shaped is collected, and nothing Network-shaped is requested.
	a.Empty(m.Devices)
	a.Empty(m.Clients)
	a.Empty(m.Sites)
	a.Empty(fake.networkPaths())

	require.Len(t, m.ControllerStatuses, 1)
	a.True(m.ControllerStatuses[0].Up)
}

// A Prometheus scrape of one target passes filter.Path. That path counts a poll as successful
// only if it produced devices or clients, and a Protect-only console produces neither -- so
// without ProtectDevices in that check a working console reports ErrDynamicLookupsDisabled.
func TestProtectOnlyControllerSucceedsOnFilteredScrape(t *testing.T) {
	t.Parallel()

	srv := (&unvr{}).start(t)

	u := newProtectOnlyInput(srv.URL)
	require.NoError(t, u.Initialize(nil))

	m, err := u.Metrics(&poller.Filter{Path: srv.URL})
	require.NoError(t, err)
	require.NotNil(t, m)
	assert.Len(t, m.ProtectDevices, 1)
}

// Protect event logs come from the legacy Protect endpoint, which is site-independent. Every
// other event collector is per-site, so all of them have to be skipped for this to work.
func TestProtectOnlyControllerCollectsEvents(t *testing.T) {
	t.Parallel()

	enabled := true
	fake := &unvr{}
	srv := fake.start(t)

	u := newProtectOnlyInput(srv.URL, func(c *inputunifi.Controller) { c.SaveProtectLogs = &enabled })
	require.NoError(t, u.Initialize(nil))

	e, err := u.Events(nil)
	require.NoError(t, err)
	require.NotNil(t, e)
	assert.Len(t, e.Logs, 1)
	assert.Empty(t, fake.networkPaths())
}

// RawMetrics' device/client kinds are all site-scoped. Returning an empty result for them
// would look like "this console has no devices"; the error says why instead.
func TestRawMetricsRejectsNetworkKindsWhenDisabled(t *testing.T) {
	t.Parallel()

	srv := (&unvr{}).start(t)

	u := newProtectOnlyInput(srv.URL)
	require.NoError(t, u.Initialize(nil))

	_, err := u.RawMetrics(&poller.Filter{Kind: "devices"})
	require.ErrorIs(t, err, inputunifi.ErrNetworkDisabled)

	// The raw-path kind still works: it is the only way to inspect a Protect-only console.
	body, err := u.RawMetrics(&poller.Filter{Kind: "other", Path: unifi.APIProtectMetaInfoPath})
	require.NoError(t, err)
	assert.JSONEq(t, metaInfoJSON, string(body))
}

// A Protect-only controller that saves nothing is always a mistake, and silently collecting
// nothing is the failure mode hardest to spot in a log, so it must be called out at startup.
func TestProtectOnlyWarnsWhenNothingWillBeCollected(t *testing.T) {
	t.Parallel()

	srv := (&unvr{}).start(t)

	disabled := false
	logger := &captureLogger{}
	u := newProtectOnlyInput(srv.URL, func(c *inputunifi.Controller) { c.SaveProtectDevices = &disabled })
	require.NoError(t, u.Initialize(logger))

	assert.Contains(t, logger.errors(), "nothing will be collected from it")
}

// The same, for a console that wants Protect devices but has no key to ask for them with.
func TestProtectOnlyWarnsWithoutAPIKey(t *testing.T) {
	t.Parallel()

	srv := (&unvr{}).start(t)

	logger := &captureLogger{}
	u := newProtectOnlyInput(srv.URL, func(c *inputunifi.Controller) { c.ProtectAPIKey = "" })
	require.NoError(t, u.Initialize(logger))

	assert.Contains(t, logger.errors(), "no protect_api_key")
}

// A key-only config is what an operator actually writes for a UNVR, and setDefaults fills the
// unset user with "unifipoller". Sending that placeholder made the console answer 403 and
// killed the controller outright, so the credentials must not be sent when no session can be
// used -- and the startup summary must not claim an authentication that never happens.
func TestProtectOnlyDoesNotSendUnusedCredentials(t *testing.T) {
	t.Parallel()

	fake := &unvr{}
	srv := fake.start(t)

	logger := &captureLogger{}
	u := newProtectOnlyInput(srv.URL, func(c *inputunifi.Controller) {
		c.User = ""
		c.Pass = ""
	})
	require.NoError(t, u.Initialize(logger))

	a := assert.New(t)
	require.NotNil(t, u.Controllers[0].Unifi)
	a.NotContains(fake.paths(), unifi.APILoginPathNew, "no session is usable, so none should be attempted")
	a.Contains(logger.infoLines(), "Protect API key only (no session needed)")
	a.NotContains(logger.infoLines(), "unifipoller")
}

// With Protect logs enabled the session is genuinely needed, so the credentials are sent.
func TestProtectOnlyLogsInForProtectLogs(t *testing.T) {
	t.Parallel()

	enabled := true
	fake := &unvr{}
	srv := fake.start(t)

	u := newProtectOnlyInput(srv.URL, func(c *inputunifi.Controller) { c.SaveProtectLogs = &enabled })
	require.NoError(t, u.Initialize(nil))

	assert.Contains(t, fake.paths(), unifi.APILoginPathNew)
}

// disable_network has to survive four independent binding paths -- toml, json, yaml and the
// UP_ environment -- each driven by its own struct tag. A typo in any one tag leaves the
// option silently inert for users of that format. Note the env name derives from the xml tag.
func TestDisableNetworkBindsInEveryConfigFormat(t *testing.T) {
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
		{"toml", "up.conf", "[[unifi.controller]]\n  url = \"https://unvr\"\n  disable_network = true\n"},
		{"json", "up.json", `{"unifi":{"controllers":[{"url":"https://unvr","disable_network":true}]}}`},
		{"yaml", "up.yaml", "unifi:\n  controllers:\n    - url: \"https://unvr\"\n      disable_network: true\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := &inputunifi.InputUnifi{Config: &inputunifi.Config{}}
			require.NoError(t, cnfgfile.Unmarshal(input, write(t, tc.file, tc.body)))
			require.Len(t, input.Controllers, 1)
			require.NotNil(t, input.Controllers[0].DisableNetwork, "disable_network did not bind from %s", tc.name)
			assert.True(t, *input.Controllers[0].DisableNetwork)
		})
	}
}

// Not parallel: t.Setenv is incompatible with a parallel test or parent.
//
// The env variable name comes from the xml tag, not the json one, which is why this is worth
// asserting rather than assuming from the toml/json/yaml cases above.
func TestDisableNetworkBindsFromEnvironment(t *testing.T) {
	t.Setenv("UP_UNIFI_CONTROLLER_0_URL", "https://unvr")
	t.Setenv("UP_UNIFI_CONTROLLER_0_DISABLE_NETWORK", "true")

	input := &inputunifi.InputUnifi{Config: &inputunifi.Config{}}
	_, err := cnfg.UnmarshalENV(input, "UP")
	require.NoError(t, err)
	require.Len(t, input.Controllers, 1)
	require.NotNil(t, input.Controllers[0].DisableNetwork, "UP_UNIFI_CONTROLLER_0_DISABLE_NETWORK did not bind")
	assert.True(t, *input.Controllers[0].DisableNetwork)
}

// The shipped examples must keep the Network application enabled, so upgrading unpoller never
// silently stops collecting from a working controller.
func TestShippedExamplesDoNotDisableNetwork(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"up.conf.example", "up.json.example", "up.yaml.example"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			input := &inputunifi.InputUnifi{Config: &inputunifi.Config{}}
			require.NoError(t, cnfgfile.Unmarshal(input, filepath.Join("..", "..", "examples", name)))

			for i, c := range input.Controllers {
				if c.DisableNetwork != nil {
					assert.False(t, *c.DisableNetwork, "%s controller %d ships with disable_network set", name, i)
				}
			}
		})
	}
}

// captureLogger records what Initialize logs, so the startup warnings can be asserted on.
type captureLogger struct {
	mu    sync.Mutex
	errs  []string
	infos []string
}

func (l *captureLogger) LogDebugf(string, ...any) {}

func (l *captureLogger) Logf(msg string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.infos = append(l.infos, fmt.Sprintf(msg, v...))
}

func (l *captureLogger) LogErrorf(msg string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.errs = append(l.errs, fmt.Sprintf(msg, v...))
}

func (l *captureLogger) errors() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	return strings.Join(l.errs, "\n")
}

func (l *captureLogger) infoLines() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	return strings.Join(l.infos, "\n")
}
