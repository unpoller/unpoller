//nolint:testpackage // white-box: exercises the unexported descriptors and export path.
package promunifi

import (
	"reflect"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

func testProtectDevices() *unifi.ProtectDevices {
	return &unifi.ProtectDevices{
		SourceName: "https://protect.example",
		NVR: &unifi.ProtectNVR{
			ProtectDeviceIdentity: unifi.ProtectDeviceIdentity{Name: "NVR", Type: "nvr", ModelKey: "nvr", State: "CONNECTED"},
		},
		Sensors: []*unifi.ProtectSensor{{
			ProtectDeviceIdentity: unifi.ProtectDeviceIdentity{Name: "Sensor", Type: "sensor", ModelKey: "sensor", State: "CONNECTED"},
			IsOpened:              unifi.FlexBool{Val: true},
			IsMotionDetected:      unifi.FlexBool{Val: true},
			BatteryStatus: &unifi.ProtectBatteryStatus{
				Percentage: unifi.FlexFloat{Val: 88},
				IsLow:      unifi.FlexBool{Val: false},
			},
			Stats: &unifi.ProtectSensorStats{
				Temperature: &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 21.5}},
				Humidity:    &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 40}},
				Light:       &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 300}},
			},
		}},
		Cameras: []*unifi.ProtectCamera{{
			ProtectDeviceIdentity: unifi.ProtectDeviceIdentity{Name: "Camera", Type: "camera", ModelKey: "camera", State: "CONNECTED"},
			IsMicEnabled:          unifi.FlexBool{Val: true},
		}},
		Lights: []*unifi.ProtectLight{{
			ProtectDeviceIdentity: unifi.ProtectDeviceIdentity{Name: "Light", Type: "light", ModelKey: "light", State: "CONNECTED"},
			IsLightOn:             unifi.FlexBool{Val: true},
			IsDark:                unifi.FlexBool{Val: true},
			IsPirMotionDetected:   unifi.FlexBool{Val: false},
		}},
		Bridges: []*unifi.ProtectBridge{{
			ProtectDeviceIdentity: unifi.ProtectDeviceIdentity{Name: "Bridge", Type: "bridge", ModelKey: "bridge", State: "CONNECTED"},
			Clients:               []string{"a", "b"},
			MaxClients:            unifi.FlexInt{Val: 8},
		}},
		LinkStations: []*unifi.ProtectLinkStation{{
			ProtectDeviceIdentity: unifi.ProtectDeviceIdentity{Name: "LinkStation", Type: "linkstation", ModelKey: "linkstation", State: "CONNECTED"},
			IsAlarmHub:            unifi.FlexBool{Val: true},
			AlarmHub:              &unifi.ProtectAlarmHub{Armed: unifi.FlexBool{Val: true}},
		}},
	}
}

// TestProtectDevicesDescriptorsAreReachableByDescribe replicates the exact reflect loop
// Describe uses. That loop skips any member which is not a *prometheus.Desc, so a nested
// sub-struct would be dropped silently rather than noisily -- and an unassigned field would
// be a nil descriptor.
func TestProtectDevicesDescriptorsAreReachableByDescribe(t *testing.T) {
	t.Parallel()

	descs := descProtectDevices("unifi_")
	v := reflect.Indirect(reflect.ValueOf(descs))
	got := []*prometheus.Desc{}

	for i := 0; i < v.NumField(); i++ {
		desc, ok := v.Field(i).Interface().(*prometheus.Desc)
		require.True(t, ok, "field %d of protectDevices is not a *prometheus.Desc; Describe would skip it",
			i)
		require.NotNil(t, desc, "field %d of protectDevices was never assigned by descProtectDevices", i)

		got = append(got, desc)
	}

	assert.Len(t, got, v.NumField())
}

func TestExportProtectDevices(t *testing.T) {
	t.Parallel()

	r := &fakeReport{}
	u := &promUnifi{ProtectDevices: descProtectDevices("unifi_")}
	u.exportProtectDevices(r, testProtectDevices())

	// nvr identity(2) + sensor identity(2)+battery(2)+stats(3)+opened/motion(2) + camera
	// identity(2)+mic(1) + light identity(2)+3 bools(3) + bridge identity(2)+2 fields(2) +
	// link station identity(2)+is_alarm_hub(1)+armed(1). This is the flattened total across
	// every send() call -- it equals the descriptor struct's field count only because this
	// fixture fires each descriptor exactly once. Do not "simplify" it to reflect.NumField():
	// add a second sensor or camera and they diverge.
	assert.Len(t, r.sent, 27, "every family is emitted for a fully populated device set")

	for i, m := range r.sent {
		require.NotNil(t, m.Desc, "metric %d has no descriptor", i)
		require.NotEmpty(t, m.Labels, "metric %d has no labels", i)
	}
}

// GetProtectDevices leaves the field for a failing endpoint nil/empty instead of failing the
// whole poll, so every device kind must be independently guarded.
func TestExportProtectDevicesTolerantOfNilSections(t *testing.T) {
	t.Parallel()

	u := &promUnifi{ProtectDevices: descProtectDevices("unifi_")}

	empty := &fakeReport{}
	u.exportProtectDevices(empty, nil)
	assert.Empty(t, empty.sent)

	r := &fakeReport{}
	u.exportProtectDevices(r, &unifi.ProtectDevices{SourceName: "https://protect.example"})
	assert.Empty(t, r.sent, "a device set with no readable endpoints reports nothing")
}

// A sensor with no battery and no stats must still report presence, state, and its two
// always-present boolean gauges -- and nothing else.
func TestExportProtectDevicesSensorTolerantOfNilBatteryAndStats(t *testing.T) {
	t.Parallel()

	u := &promUnifi{ProtectDevices: descProtectDevices("unifi_")}
	r := &fakeReport{}

	u.exportProtectDevices(r, &unifi.ProtectDevices{
		SourceName: "https://protect.example",
		Sensors: []*unifi.ProtectSensor{{
			ProtectDeviceIdentity: unifi.ProtectDeviceIdentity{Name: "Sensor", Type: "sensor", ModelKey: "sensor", State: "CONNECTED"},
		}},
	})

	// identity(2) + is_opened/is_motion(2), no battery or stats series.
	assert.Len(t, r.sent, 4)
}

func TestExportProtectSensorAirQuality(t *testing.T) {
	t.Parallel()

	u := &promUnifi{ProtectDevices: descProtectDevices("unifi_")}

	unknownOnly := &fakeReport{}
	u.exportProtectDevices(unknownOnly, &unifi.ProtectDevices{
		SourceName: "https://protect.example",
		Sensors: []*unifi.ProtectSensor{{
			ProtectDeviceIdentity: unifi.ProtectDeviceIdentity{Name: "AQ", Type: "UP-AirQuality", ModelKey: "sensor", State: "CONNECTED"},
			Stats: &unifi.ProtectSensorStats{
				Temperature: &unifi.ProtectSensorStatValue{Status: "unknown"},
				Humidity:    &unifi.ProtectSensorStatValue{Status: "unknown"},
				Light:       &unifi.ProtectSensorStatValue{Status: "unknown"},
			},
		}},
	})

	unknown := readingsByName(t, unknownOnly.sent)
	assert.NotContains(t, unknown, "unifi_protect_sensor_temperature_celsius")
	assert.NotContains(t, unknown, "unifi_protect_sensor_humidity_percent")
	assert.NotContains(t, unknown, "unifi_protect_sensor_light_lux")
	assert.NotContains(t, unknown, "unifi_protect_sensor_aqi")

	withAir := &fakeReport{}
	u.exportProtectDevices(withAir, &unifi.ProtectDevices{
		SourceName: "https://protect.example",
		Sensors: []*unifi.ProtectSensor{{
			ProtectDeviceIdentity: unifi.ProtectDeviceIdentity{Name: "AQ", Type: "UP-AirQuality", ModelKey: "sensor", State: "CONNECTED"},
			Stats: &unifi.ProtectSensorStats{
				Temperature: &unifi.ProtectSensorStatValue{Status: "unknown"},
				Humidity:    &unifi.ProtectSensorStatValue{Status: "unknown"},
			},
			AirQuality: &unifi.ProtectAirQuality{
				AQI:         &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 7}, Status: "neutral"},
				Vape:        &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 0}, Status: "safe"},
				CO2:         &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 452}, Status: "neutral"},
				PM2p5:       &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 1.79}, Status: "neutral"},
				Humidity:    &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 51}, Status: "neutral"},
				Temperature: &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 25.4}, Status: "neutral"},
			},
		}},
	})

	got := readingsByName(t, withAir.sent)
	assert.InDelta(t, 25.4, got["unifi_protect_sensor_temperature_celsius"], 0.001)
	assert.InDelta(t, 51, got["unifi_protect_sensor_humidity_percent"], 0.001)
	assert.InDelta(t, 7, got["unifi_protect_sensor_aqi"], 0.001)
	assert.InDelta(t, 0, got["unifi_protect_sensor_vape"], 0.001)
	assert.InDelta(t, 452, got["unifi_protect_sensor_co2_ppm"], 0.001)
	assert.InDelta(t, 1.79, got["unifi_protect_sensor_pm2_5"], 0.001)
	assert.NotContains(t, got, "unifi_protect_sensor_tvoc")

	statsWin := &fakeReport{}
	u.exportProtectDevices(statsWin, &unifi.ProtectDevices{
		SourceName: "https://protect.example",
		Sensors: []*unifi.ProtectSensor{{
			ProtectDeviceIdentity: unifi.ProtectDeviceIdentity{Name: "Both", Type: "sensor", ModelKey: "sensor", State: "CONNECTED"},
			Stats: &unifi.ProtectSensorStats{
				Temperature: &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 21.5}, Status: "safe"},
			},
			AirQuality: &unifi.ProtectAirQuality{
				Temperature: &unifi.ProtectSensorStatValue{Value: unifi.FlexFloat{Val: 99}, Status: "neutral"},
			},
		}},
	})

	assert.InDelta(t, 21.5, readingsByName(t, statsWin.sent)["unifi_protect_sensor_temperature_celsius"], 0.001)
}

func readingsByName(t *testing.T, sent []*metric) map[string]float64 {
	t.Helper()

	out := make(map[string]float64, len(sent))
	for _, m := range sent {
		name := prometheusDescName(t, m.Desc)
		value, ok := m.Value.(float64)
		require.True(t, ok, "metric %s value is %T", name, m.Value)

		out[name] = value
	}

	return out
}

func prometheusDescName(t *testing.T, desc *prometheus.Desc) string {
	t.Helper()

	const key = `fqName: "`

	text := desc.String()
	start := strings.Index(text, key)
	require.GreaterOrEqual(t, start, 0, "descriptor has no fqName: %s", text)

	rest := text[start+len(key):]
	end := strings.Index(rest, `"`)
	require.GreaterOrEqual(t, end, 0, "descriptor fqName is unclosed: %s", text)

	return rest[:end]
}
