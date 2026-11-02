package inputunifi_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
	"github.com/unpoller/unpoller/pkg/inputunifi"
	"github.com/unpoller/unpoller/pkg/poller"
	"golift.io/cnfg"
	"golift.io/cnfgfile"
)

// tagsController is a fake Network application that counts device-tags requests and answers
// enough of the API for NewUnifi and one poll to succeed.
func tagsController(t *testing.T) (*httptest.Server, *atomic.Int32) {
	t.Helper()

	var tagRequests atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := `{"meta":{"rc":"ok"},"data":[]}`

		switch path := r.URL.Path; {
		case strings.HasSuffix(path, "/device-tags"):
			tagRequests.Add(1)
		case path == "/", path == unifi.APILoginPathNew:
		case path == "/proxy/network/status":
			body = `{"meta":{"server_version":"10.6.109"}}`
		case path == "/proxy/network/api/stat/sites":
			body = `{"data":[{"name":"default","desc":"Main"}]}`
		case strings.HasSuffix(path, "/stat/sysinfo"):
			body = `{"data":[{"hostname":"controller","version":"10.6.109","update_available":false}]}`
		case path == unifi.APIConsoleSystemPath:
			body = consoleBody
		case strings.HasSuffix(path, "/stat/device"), strings.HasSuffix(path, "/stat/sta"):
		default:
			w.WriteHeader(http.StatusNotFound)
		}

		_, err := w.Write([]byte(body))
		assert.NoError(t, err)
	}))
	t.Cleanup(srv.Close)

	return srv, &tagRequests
}

// tagsInput builds an input with one controller against srv, with skip_device_tags set at the
// defaults level and/or on the controller, exactly as a config file would.
func tagsInput(t *testing.T, srv *httptest.Server, defaults, controller *bool) *inputunifi.InputUnifi {
	t.Helper()

	disabled := false
	u := &inputunifi.InputUnifi{Config: &inputunifi.Config{
		Default: inputunifi.Controller{SkipDeviceTags: defaults},
		Controllers: []*inputunifi.Controller{{
			URL: srv.URL, User: "monitor", Pass: "test-password",
			SaveSites: &disabled, SaveDPI: &disabled, SaveRogue: &disabled,
			SaveTraffic: &disabled, SaveSpeedTest: &disabled, SkipDeviceTags: controller,
		}},
	}}
	require.NoError(t, u.Initialize(nil))
	require.NotNil(t, u.Controllers[0].Unifi)
	t.Cleanup(u.Controllers[0].Unifi.CloseIdleConnections)

	return u
}

// TestSkipDeviceTags checks that the option reaches the UniFi client and that a poll then does
// (or does not) request device-tags. Unset means fetched, as it always has been; the defaults
// level is inherited; a controller's own setting wins over the defaults.
func TestSkipDeviceTags(t *testing.T) {
	t.Parallel()

	enabled, disabled := true, false

	for _, tc := range []struct {
		name       string
		defaults   *bool
		controller *bool
		skip       bool
	}{
		{name: "unset fetches tags"},
		{name: "defaults on is inherited", defaults: &enabled, skip: true},
		{name: "controller on", controller: &enabled, skip: true},
		{name: "controller off overrides defaults on", defaults: &enabled, controller: &disabled},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv, tagRequests := tagsController(t)
			u := tagsInput(t, srv, tc.defaults, tc.controller)
			c := u.Controllers[0]

			require.NotNil(t, c.SkipDeviceTags, "defaults must fill in skip_device_tags")
			assert.Equal(t, tc.skip, *c.SkipDeviceTags)
			assert.Equal(t, tc.skip, c.Unifi.SkipDeviceTags, "the option must reach the UniFi client")

			_, err := u.Metrics(&poller.Filter{Path: c.URL})
			require.NoError(t, err)

			if tc.skip {
				assert.Zero(t, tagRequests.Load(), "a poll must not request device-tags")
			} else {
				assert.Positive(t, tagRequests.Load(), "a poll must request device-tags")
			}
		})
	}
}

func TestSkipDeviceTagsBindsInEveryConfigFormat(t *testing.T) {
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
		{"toml", "up.conf", "[unifi.defaults]\n  skip_device_tags = true\n" +
			"[[unifi.controller]]\n  url = \"https://udm\"\n  skip_device_tags = false\n"},
		{"json", "up.json", `{"unifi":{"defaults":{"skip_device_tags":true},` +
			`"controllers":[{"url":"https://udm","skip_device_tags":false}]}}`},
		{"yaml", "up.yaml", "unifi:\n  defaults:\n    skip_device_tags: true\n" +
			"  controllers:\n    - url: \"https://udm\"\n      skip_device_tags: false\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := &inputunifi.InputUnifi{Config: &inputunifi.Config{}}
			require.NoError(t, cnfgfile.Unmarshal(input, write(t, tc.file, tc.body)))
			require.NotNil(t, input.Default.SkipDeviceTags, "skip_device_tags did not bind on defaults from %s", tc.name)
			assert.True(t, *input.Default.SkipDeviceTags)
			require.Len(t, input.Controllers, 1)
			require.NotNil(t, input.Controllers[0].SkipDeviceTags, "skip_device_tags did not bind on the controller from %s", tc.name)
			assert.False(t, *input.Controllers[0].SkipDeviceTags)
		})
	}
}

// Not parallel: t.Setenv is incompatible with a parallel test or parent.
func TestSkipDeviceTagsBindsFromEnvironment(t *testing.T) {
	t.Setenv("UP_UNIFI_DEFAULT_SKIP_DEVICE_TAGS", "true")
	t.Setenv("UP_UNIFI_CONTROLLER_0_URL", "https://udm")
	t.Setenv("UP_UNIFI_CONTROLLER_0_SKIP_DEVICE_TAGS", "false")

	input := &inputunifi.InputUnifi{Config: &inputunifi.Config{}}
	_, err := cnfg.UnmarshalENV(input, "UP")
	require.NoError(t, err)
	require.NotNil(t, input.Default.SkipDeviceTags, "UP_UNIFI_DEFAULT_SKIP_DEVICE_TAGS did not bind")
	assert.True(t, *input.Default.SkipDeviceTags)
	require.Len(t, input.Controllers, 1)
	require.NotNil(t, input.Controllers[0].SkipDeviceTags, "UP_UNIFI_CONTROLLER_0_SKIP_DEVICE_TAGS did not bind")
	assert.False(t, *input.Controllers[0].SkipDeviceTags)
}

// The shipped examples must keep fetching tags, so upgrading unpoller never silently drops the
// tag labels from anyone's metrics.
func TestShippedExamplesDoNotSkipDeviceTags(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"up.conf.example", "up.json.example", "up.yaml.example",
		"remote_api_example.conf", "remote_api_example.json", "remote_api_example.yaml"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			input := &inputunifi.InputUnifi{Config: &inputunifi.Config{}}
			require.NoError(t, cnfgfile.Unmarshal(input, filepath.Join("..", "..", "examples", name)))

			if input.Default.SkipDeviceTags != nil {
				assert.False(t, *input.Default.SkipDeviceTags, "%s ships with skip_device_tags set on defaults", name)
			}

			for i, c := range input.Controllers {
				if c.SkipDeviceTags != nil {
					assert.False(t, *c.SkipDeviceTags, "%s controller %d ships with skip_device_tags set", name, i)
				}
			}
		})
	}
}
