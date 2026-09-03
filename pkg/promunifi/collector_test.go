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
