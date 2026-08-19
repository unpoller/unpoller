package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v6"
)

// protectDevices holds Prometheus descriptors for UniFi Protect device metrics.
//
// The struct is deliberately flat even though the metrics come from six different device kinds
// (NVR, sensor, camera, light, bridge, link station): Describe reflects over the members of this
// struct and sends each one to the channel, so a nested sub-struct would be sent as something
// that is not a *prometheus.Desc.
type protectDevices struct {
	// Shared across every device kind.
	Presence *prometheus.Desc
	State    *prometheus.Desc
	// Sensors.
	SensorBatteryPercent *prometheus.Desc
	SensorBatteryLow     *prometheus.Desc
	SensorTemperature    *prometheus.Desc
	SensorHumidity       *prometheus.Desc
	SensorLight          *prometheus.Desc
	SensorIsOpened       *prometheus.Desc
	SensorIsMotion       *prometheus.Desc
	// Cameras.
	CameraMicEnabled *prometheus.Desc
	// Lights.
	LightIsOn           *prometheus.Desc
	LightIsDark         *prometheus.Desc
	LightIsMotionActive *prometheus.Desc
	// Bridges.
	BridgeClientCount *prometheus.Desc
	BridgeMaxClients  *prometheus.Desc
	// Link stations.
	LinkStationIsAlarmHub *prometheus.Desc
	LinkStationArmed      *prometheus.Desc
}

func descProtectDevices(ns string) *protectDevices { // nolint: funlen
	ns += "protect_"
	device := []string{"source", "name", "type", "model_key"}

	return &protectDevices{
		Presence: prometheus.NewDesc(ns+"device_present",
			"Protect device reachable (always 1 when present)", device, nil),
		State: prometheus.NewDesc(ns+"device_state",
			"Protect device connection state: 2=connected, 1=connecting, 0=disconnected, -1=unknown", device, nil),
		SensorBatteryPercent: prometheus.NewDesc(ns+"sensor_battery_percent",
			"Protect sensor battery charge percentage", device, nil),
		SensorBatteryLow: prometheus.NewDesc(ns+"sensor_battery_low",
			"Protect sensor battery is low (1) or not (0)", device, nil),
		SensorTemperature: prometheus.NewDesc(ns+"sensor_temperature_celsius",
			"Protect sensor temperature reading", device, nil),
		SensorHumidity: prometheus.NewDesc(ns+"sensor_humidity_percent",
			"Protect sensor humidity reading", device, nil),
		SensorLight: prometheus.NewDesc(ns+"sensor_light_lux",
			"Protect sensor light reading", device, nil),
		SensorIsOpened: prometheus.NewDesc(ns+"sensor_is_opened",
			"Protect door/window sensor is opened (1) or closed (0)", device, nil),
		SensorIsMotion: prometheus.NewDesc(ns+"sensor_is_motion_detected",
			"Protect sensor currently detects motion (1) or not (0)", device, nil),
		CameraMicEnabled: prometheus.NewDesc(ns+"camera_mic_enabled",
			"Protect camera microphone is enabled (1) or not (0)", device, nil),
		LightIsOn: prometheus.NewDesc(ns+"light_is_on",
			"Protect light is currently on (1) or off (0)", device, nil),
		LightIsDark: prometheus.NewDesc(ns+"light_is_dark",
			"Protect light's ambient light sensor reports dark (1) or not (0)", device, nil),
		LightIsMotionActive: prometheus.NewDesc(ns+"light_is_pir_motion_detected",
			"Protect light PIR sensor currently detects motion (1) or not (0)", device, nil),
		BridgeClientCount: prometheus.NewDesc(ns+"bridge_client_count",
			"Protect bridge current connected client count", device, nil),
		BridgeMaxClients: prometheus.NewDesc(ns+"bridge_max_clients",
			"Protect bridge maximum supported client count", device, nil),
		LinkStationIsAlarmHub: prometheus.NewDesc(ns+"link_station_is_alarm_hub",
			"Protect link station is acting as the alarm hub (1) or not (0)", device, nil),
		LinkStationArmed: prometheus.NewDesc(ns+"link_station_armed",
			"Protect link station alarm hub is armed (1) or disarmed (0)", device, nil),
	}
}

// exportProtectDevices emits every metric family for one Protect controller's device set.
//
// Each device kind is guarded independently because GetProtectDevices leaves fields for a
// failing endpoint nil/empty rather than failing the whole poll.
func (u *promUnifi) exportProtectDevices(r report, d *unifi.ProtectDevices) {
	if d == nil {
		return
	}

	if d.NVR != nil {
		u.exportProtectIdentity(r, d.SourceName, &d.NVR.ProtectDeviceIdentity)
	}

	for _, s := range d.Sensors {
		if s == nil {
			continue
		}

		labels := u.exportProtectIdentity(r, d.SourceName, &s.ProtectDeviceIdentity)

		if s.BatteryStatus != nil {
			r.send([]*metric{
				{u.ProtectDevices.SensorBatteryPercent, gauge, s.BatteryStatus.Percentage.Val, labels},
			})

			isLow := 0.0
			if s.BatteryStatus.IsLow.Val {
				isLow = 1.0
			}

			r.send([]*metric{{u.ProtectDevices.SensorBatteryLow, gauge, isLow, labels}})
		}

		if s.Stats != nil {
			if s.Stats.Temperature != nil {
				r.send([]*metric{{u.ProtectDevices.SensorTemperature, gauge, s.Stats.Temperature.Value.Val, labels}})
			}

			if s.Stats.Humidity != nil {
				r.send([]*metric{{u.ProtectDevices.SensorHumidity, gauge, s.Stats.Humidity.Value.Val, labels}})
			}

			if s.Stats.Light != nil {
				r.send([]*metric{{u.ProtectDevices.SensorLight, gauge, s.Stats.Light.Value.Val, labels}})
			}
		}

		isOpened := 0.0
		if s.IsOpened.Val {
			isOpened = 1.0
		}

		isMotion := 0.0
		if s.IsMotionDetected.Val {
			isMotion = 1.0
		}

		r.send([]*metric{
			{u.ProtectDevices.SensorIsOpened, gauge, isOpened, labels},
			{u.ProtectDevices.SensorIsMotion, gauge, isMotion, labels},
		})
	}

	for _, c := range d.Cameras {
		if c == nil {
			continue
		}

		labels := u.exportProtectIdentity(r, d.SourceName, &c.ProtectDeviceIdentity)

		micEnabled := 0.0
		if c.IsMicEnabled.Val {
			micEnabled = 1.0
		}

		r.send([]*metric{{u.ProtectDevices.CameraMicEnabled, gauge, micEnabled, labels}})
	}

	for _, l := range d.Lights {
		if l == nil {
			continue
		}

		labels := u.exportProtectIdentity(r, d.SourceName, &l.ProtectDeviceIdentity)

		isOn, isDark, isMotion := 0.0, 0.0, 0.0
		if l.IsLightOn.Val {
			isOn = 1.0
		}

		if l.IsDark.Val {
			isDark = 1.0
		}

		if l.IsPirMotionDetected.Val {
			isMotion = 1.0
		}

		r.send([]*metric{
			{u.ProtectDevices.LightIsOn, gauge, isOn, labels},
			{u.ProtectDevices.LightIsDark, gauge, isDark, labels},
			{u.ProtectDevices.LightIsMotionActive, gauge, isMotion, labels},
		})
	}

	for _, b := range d.Bridges {
		if b == nil {
			continue
		}

		labels := u.exportProtectIdentity(r, d.SourceName, &b.ProtectDeviceIdentity)

		r.send([]*metric{
			{u.ProtectDevices.BridgeClientCount, gauge, float64(len(b.Clients)), labels},
			{u.ProtectDevices.BridgeMaxClients, gauge, b.MaxClients.Val, labels},
		})
	}

	for _, ls := range d.LinkStations {
		if ls == nil {
			continue
		}

		labels := u.exportProtectIdentity(r, d.SourceName, &ls.ProtectDeviceIdentity)

		isHub := 0.0
		if ls.IsAlarmHub.Val {
			isHub = 1.0
		}

		r.send([]*metric{{u.ProtectDevices.LinkStationIsAlarmHub, gauge, isHub, labels}})

		if ls.AlarmHub != nil {
			armed := 0.0
			if ls.AlarmHub.Armed.Val {
				armed = 1.0
			}

			r.send([]*metric{{u.ProtectDevices.LinkStationArmed, gauge, armed, labels}})
		}
	}
}

// exportProtectIdentity emits the presence and state gauges shared by every Protect device kind
// and returns the label set used for that device's other metrics.
func (u *promUnifi) exportProtectIdentity(r report, sourceName string, id *unifi.ProtectDeviceIdentity) []string {
	labels := []string{sourceName, id.Name, id.Type, id.ModelKey}

	r.send([]*metric{
		{u.ProtectDevices.Presence, gauge, 1.0, labels},
		{u.ProtectDevices.State, gauge, unifi.ProtectStateValue(id.State), labels},
	})

	return labels
}
