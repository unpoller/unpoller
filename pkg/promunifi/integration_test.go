package promunifi_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unpoller/pkg/promunifi"
)

// TestDescriptorRegistration verifies that all Prometheus metric descriptors
// can be registered without name or label conflicts.
//
// Regression: the Integration/v1 device metrics introduced in PR#999 used
// the same "_device_" namespace prefix as the existing device metrics, causing
// "inconsistent label names" panics on startup and preventing any metrics from
// being collected.
func TestDescriptorRegistration(t *testing.T) {
	reg := prometheus.NewRegistry()
	collector := promunifi.NewCollectorForTesting("unpoller")
	err := reg.Register(collector)
	require.NoError(t, err, "all metric descriptors must have unique fully-qualified names and label sets")
}
