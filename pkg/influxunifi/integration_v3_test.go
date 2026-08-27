package influxunifi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	influxdb3 "github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unpoller/pkg/influxunifi"
	"github.com/unpoller/unpoller/pkg/unittest"
	"golift.io/cnfg"
)

type v3WriteCapture struct {
	mu    sync.Mutex
	batch []string
}

func (c *v3WriteCapture) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)

		return
	}

	if !strings.Contains(r.URL.Path, "write") {
		http.NotFound(w, r)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	c.mu.Lock()
	c.batch = append(c.batch, string(body))
	c.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func (c *v3WriteCapture) lines() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	raw := strings.Join(c.batch, "\n")

	return strings.Split(strings.TrimSpace(raw), "\n")
}

func TestInfluxV3Integration(t *testing.T) {
	capture := &v3WriteCapture{}
	server := httptest.NewServer(http.HandlerFunc(capture.handler))
	t.Cleanup(server.Close)

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     server.URL,
		Token:    "test-token",
		Database: "unpoller",
		WriteOptions: &influxdb3.WriteOptions{
			UseV2Api: false,
		},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	testRig := unittest.NewTestSetup(t)
	t.Cleanup(testRig.Close)

	u := influxunifi.InfluxUnifi{
		Collector:      testRig.Collector,
		Version:        influxunifi.InfluxV3,
		InfluxV3Client: client,
		InfluxDB: &influxunifi.InfluxDB{
			Config: &influxunifi.Config{
				Version:   3,
				Database:  "unpoller",
				AuthToken: "test-token",
				URL:       server.URL,
				Interval:  cnfg.Duration{Duration: time.Hour},
			},
		},
	}

	testRig.Initialize()
	u.Poll(time.Minute)

	lines := capture.lines()
	require.NotEmpty(t, lines)

	body := strings.Join(lines, "\n")
	assert.Contains(t, body, "subsystems,")
	assert.Contains(t, body, "clients,")
	assert.Contains(t, body, "channel_name=")
	assert.NotRegexp(t, `(?m)^clients,[^ ]* channel=`, body)
}
