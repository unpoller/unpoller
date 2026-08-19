//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

// fakeReport collects the metrics an export function sends, and nothing else. It embeds
// *Report (nil) only to satisfy the rest of the report interface; the export path under test
// calls nothing but send.
type fakeReport struct {
	*Report

	sent []*metric
}

func (f *fakeReport) send(m []*metric) { f.sent = append(f.sent, m...) }

func testUNASDevice() *unifi.UNASDevice {
	return &unifi.UNASDevice{
		SourceName: "https://unas.example",
		DeviceInfo: &unifi.UNASDeviceInfo{
			Name:   "UNAS-Pro",
			Model:  "UNASPro",
			CPU:    unifi.UNASCPU{CurrentLoad: unifi.FlexInt{Val: 12}, Temperature: unifi.FlexInt{Val: 45}},
			Memory: unifi.UNASMemory{Total: unifi.FlexInt{Val: 8000}},
		},
		NetworkIO: &unifi.UNASNetworkIO{ReceiveKBPS: unifi.FlexInt{Val: 120}},
		Storage: &unifi.UNASStorage{
			Pools: []unifi.UNASPool{{
				ID: "pool-1", RaidGroups: []unifi.UNASRaidGroup{{ID: "rg-1"}},
			}},
			Disks: []unifi.UNASDisk{{SlotID: "1", Temperature: unifi.FlexInt{Val: 38}}},
		},
		Drives: []*unifi.UNASDrive{{ID: "d-1", Name: "media"}, nil},
	}
}

// TestUNASDescriptorsAreReachableByDescribe replicates the exact reflect loop Describe uses.
// That loop skips any member which is not a *prometheus.Desc, so a nested sub-struct would be
// dropped silently rather than noisily -- and an unassigned field would be a nil descriptor.
// The loop is duplicated rather than driven through Describe itself because Describe walks
// every metric family and the descriptors are assigned inline in Run, which also starts an
// HTTP listener; there is no way to build a collector holding only this family.
func TestUNASDescriptorsAreReachableByDescribe(t *testing.T) {
	t.Parallel()

	descs := descUNASDevice("unifi_")
	v := reflect.Indirect(reflect.ValueOf(descs))
	got := []*prometheus.Desc{}

	for i := 0; i < v.NumField(); i++ {
		desc, ok := v.Field(i).Interface().(*prometheus.Desc)
		require.True(t, ok, "field %d of unasDevice is not a *prometheus.Desc; Describe would skip it",
			i)
		require.NotNil(t, desc, "field %d of unasDevice was never assigned by descUNASDevice", i)

		got = append(got, desc)
	}

	assert.Len(t, got, v.NumField())
}

func TestExportUNASDevice(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	u := &promUnifi{UNASDevice: descUNASDevice("unifi_")}
	u.exportUNASDevice(r, testUNASDevice())

	// presence(1) + console(5) + net io(2) + pool(2) + raid(3) + disk(11) + one share(3).
	// This is the flattened total across the seven send() calls -- it equals the descriptor
	// struct's field count only because this fixture fires each descriptor exactly once. Do
	// not "simplify" it to reflect.NumField(): add a second pool or disk and they diverge.
	assert.Len(t, r.sent, 27, "every family is emitted for a fully populated console")

	for i, m := range r.sent {
		require.NotNil(t, m.Desc, "metric %d has no descriptor", i)
		require.NotEmpty(t, m.Labels, "metric %d has no labels", i)
	}
}

// GetUNASDevice leaves the field for a failing endpoint nil instead of failing the whole
// poll, so every section must be independently guarded.
func TestExportUNASDeviceTolerantOfNilSections(t *testing.T) {
	t.Parallel()

	u := &promUnifi{UNASDevice: descUNASDevice("unifi_")}

	r := &fakeReport{}
	u.exportUNASDevice(r, &unifi.UNASDevice{SourceName: "https://unas.example"})
	assert.Len(t, r.sent, 1, "a console with no readable endpoints still reports its presence")

	empty := &fakeReport{}
	u.exportUNASDevice(empty, nil)
	assert.Empty(t, empty.sent)
}
