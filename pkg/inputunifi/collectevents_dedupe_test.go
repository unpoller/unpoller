package inputunifi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
	"github.com/unpoller/unpoller/pkg/inputunifi"
)

const (
	networkStatusJSON = `{"meta":{"rc":"ok","server_version":"9.3.45","up":true}}`
	networkSitesJSON  = `{"data":[{"_id":"site-1","name":"default","desc":"Default"}]}`
	// Two distinct entries, one page. Any duplication in the output is ours.
	networkSyslogJSON = `{"data":[` +
		`{"id":"log-1","key":"EVT_AP_Connected","category":"DEVICE","event":"AP connected","severity":"INFO","timestamp":1735689600000,"title_raw":"AP connected"},` +
		`{"id":"log-2","key":"EVT_SW_Connected","category":"DEVICE","event":"Switch connected","severity":"INFO","timestamp":1735689660000,"title_raw":"Switch connected"}` +
		`],"page_number":0,"total_element_count":2,"total_page_count":1}`
)

// networkConsole is a fake UniFi OS console WITH a Network application: it logs
// in, reports one site, and serves two v2 system-log entries. Every other
// collector endpoint 404s, which is what a Network 10.x controller does for the
// legacy stat/event, stat/alarm and stat/ips paths.
type networkConsole struct {
	syslogPolls atomic.Int32
	mu          sync.Mutex
	requested   []string
}

func (s *networkConsole) start(t *testing.T) *httptest.Server {
	t.Helper()

	record := func(r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.requested = append(s.requested, r.URL.Path)
	}

	serve := func(body string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			record(r)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(body))
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc(unifi.APIPrefixNew+unifi.APIStatusPath, serve(networkStatusJSON))
	mux.HandleFunc(unifi.APIPrefixNew+unifi.APISiteList, serve(networkSitesJSON))
	mux.HandleFunc(unifi.APIPrefixNew+fmt.Sprintf(unifi.APISystemLogPath, "default"),
		func(w http.ResponseWriter, r *http.Request) {
			s.syslogPolls.Add(1)
			serve(networkSyslogJSON)(w, r)
		})
	mux.HandleFunc(unifi.APILoginPathNew, func(w http.ResponseWriter, r *http.Request) {
		record(r)
		w.Header().Set("x-csrf-token", "csrf-token")
		w.WriteHeader(http.StatusOK)
	})
	// "/" answers 200 so the library picks the UniFi OS (/proxy/network) paths;
	// anything else under it is an endpoint this console does not have.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		record(r)

		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(unifiOSSPA))

			return
		}

		http.NotFound(w, r)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// TestControllerEventsAreNotDuplicated pins the contract between
// collectControllerEvents and its collectors. Each collector returns the
// accumulated slice, and the loop used to APPEND that return value onto the
// slice it already held -- so an entry was stored once more for every collector
// that ran after the one that found it. With only save_syslog on, the two
// entries below came out as four; with alarms on they would have come out
// eight times each.
func TestControllerEventsAreNotDuplicated(t *testing.T) {
	t.Parallel()

	enabled := true
	fake := &networkConsole{}
	srv := fake.start(t)

	c := &inputunifi.Controller{
		URL:        srv.URL,
		User:       "ro-user",
		Pass:       "secret",
		SaveSyslog: &enabled,
	}
	u := &inputunifi.InputUnifi{Config: &inputunifi.Config{Controllers: []*inputunifi.Controller{c}}}
	require.NoError(t, u.Initialize(nil))

	e, err := u.Events(nil)
	require.NoError(t, err)
	require.NotNil(t, e)

	// The console was asked exactly once, so anything beyond two is our own doing.
	assert.Equal(t, int32(1), fake.syslogPolls.Load(), "system log polled more than once")
	require.Len(t, e.Logs, 2, "each system-log entry must be collected exactly once")

	ids := make([]string, 0, len(e.Logs))

	for _, l := range e.Logs {
		entry, ok := l.(*unifi.SystemLogEntry)
		require.True(t, ok, "unexpected event type %T", l)

		ids = append(ids, entry.ID)
	}

	assert.ElementsMatch(t, []string{"log-1", "log-2"}, ids)
}
