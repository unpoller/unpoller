//nolint:testpackage // white-box: exercises the unexported output plugin.
package promunifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unpoller/pkg/poller"
)

func TestDebugOutputAcceptsUnsetHTTPListen(t *testing.T) {
	poller.SetHealthCheckMode(true)
	t.Cleanup(func() { poller.SetHealthCheckMode(false) })

	u := &promUnifi{Config: &Config{}}

	ok, err := u.DebugOutput()
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, defaultHTTPListen, u.HTTPListen)
}

func TestDebugOutputAcceptsAnIPv6HTTPListen(t *testing.T) {
	poller.SetHealthCheckMode(true)
	t.Cleanup(func() { poller.SetHealthCheckMode(false) })

	tests := []struct {
		listen string
		valid  bool
	}{
		{"0.0.0.0:9130", true},
		{"[::]:9130", true},
		{":9130", true},
		{"0.0.0.0", false},
		{"0.0.0.0:9130:9131", false},
	}

	for _, test := range tests {
		u := &promUnifi{Config: &Config{HTTPListen: test.listen}}

		ok, err := u.DebugOutput()
		assert.Equal(t, test.valid, ok, test.listen)

		if test.valid {
			require.NoError(t, err, test.listen)
		} else {
			require.Error(t, err, test.listen)
		}
	}
}
