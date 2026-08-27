package datadogunifi

import (
	"github.com/unpoller/unifi/v6"
)

// batchProtectDevices generates UniFi Protect device datapoints for Datadog.
//
// Each device kind is guarded independently because GetProtectDevices leaves the field for a
// failing endpoint nil/empty rather than failing the whole poll.
func (u *DatadogUnifi) batchProtectDevices(r report, d *unifi.ProtectDevices) {
	if d == nil {
		return
	}

	if d.NVR != nil {
		u.batchProtectNVR(r, d)
	}

	for _, s := range d.Sensors {
		if s != nil {
			u.batchProtectSensor(r, d, s)
		}
	}

	for _, c := range d.Cameras {
		if c != nil {
			u.batchProtectCamera(r, d, c)
		}
	}

	for _, l := range d.Lights {
		if l != nil {
			u.batchProtectLight(r, d, l)
		}
	}

	for _, b := range d.Bridges {
		if b != nil {
			u.batchProtectBridge(r, d, b)
		}
	}

	for _, ls := range d.LinkStations {
		if ls != nil {
			u.batchProtectLinkStation(r, d, ls)
		}
	}
}

func (u *DatadogUnifi) batchProtectNVR(r report, d *unifi.ProtectDevices) {
	metricName := metricNamespace("protect_nvr")
	nvr := d.NVR

	tags := tagMapToTags(cleanTags(map[string]string{
		"source": d.SourceName,
		"name":   nvr.Name,
		"mac":    nvr.MAC,
		"guid":   nvr.GUID,
	}))

	_ = r.reportGauge(metricName("present"), 1.0, tags)
	_ = r.reportGauge(metricName("state"), unifi.ProtectStateValue(nvr.State), tags)
}

func (u *DatadogUnifi) batchProtectSensor(r report, d *unifi.ProtectDevices, s *unifi.ProtectSensor) {
	metricName := metricNamespace("protect_sensor")

	tags := tagMapToTags(cleanTags(map[string]string{
		"source": d.SourceName,
		"name":   s.Name,
		"mac":    s.MAC,
		"guid":   s.GUID,
	}))

	_ = r.reportGauge(metricName("present"), 1.0, tags)
	_ = r.reportGauge(metricName("state"), unifi.ProtectStateValue(s.State), tags)
	_ = r.reportGauge(metricName("is_opened"), boolToFloat64(s.IsOpened.Val), tags)
	_ = r.reportGauge(metricName("is_motion_detected"), boolToFloat64(s.IsMotionDetected.Val), tags)

	if s.BatteryStatus != nil {
		_ = r.reportGauge(metricName("battery_percent"), s.BatteryStatus.Percentage.Val, tags)
		_ = r.reportGauge(metricName("battery_low"), boolToFloat64(s.BatteryStatus.IsLow.Val), tags)
	}

	if s.Stats != nil {
		if s.Stats.Temperature != nil {
			_ = r.reportGauge(metricName("temperature"), s.Stats.Temperature.Value.Val, tags)
		}

		if s.Stats.Humidity != nil {
			_ = r.reportGauge(metricName("humidity"), s.Stats.Humidity.Value.Val, tags)
		}

		if s.Stats.Light != nil {
			_ = r.reportGauge(metricName("light"), s.Stats.Light.Value.Val, tags)
		}
	}
}

func (u *DatadogUnifi) batchProtectCamera(r report, d *unifi.ProtectDevices, c *unifi.ProtectCamera) {
	metricName := metricNamespace("protect_camera")

	tags := tagMapToTags(cleanTags(map[string]string{
		"source": d.SourceName,
		"name":   c.Name,
		"mac":    c.MAC,
		"guid":   c.GUID,
	}))

	_ = r.reportGauge(metricName("present"), 1.0, tags)
	_ = r.reportGauge(metricName("state"), unifi.ProtectStateValue(c.State), tags)
	_ = r.reportGauge(metricName("mic_enabled"), boolToFloat64(c.IsMicEnabled.Val), tags)
}

func (u *DatadogUnifi) batchProtectLight(r report, d *unifi.ProtectDevices, l *unifi.ProtectLight) {
	metricName := metricNamespace("protect_light")

	tags := tagMapToTags(cleanTags(map[string]string{
		"source": d.SourceName,
		"name":   l.Name,
		"mac":    l.MAC,
		"guid":   l.GUID,
	}))

	_ = r.reportGauge(metricName("present"), 1.0, tags)
	_ = r.reportGauge(metricName("state"), unifi.ProtectStateValue(l.State), tags)
	_ = r.reportGauge(metricName("is_light_on"), boolToFloat64(l.IsLightOn.Val), tags)
	_ = r.reportGauge(metricName("is_dark"), boolToFloat64(l.IsDark.Val), tags)
	_ = r.reportGauge(metricName("is_pir_motion_detected"), boolToFloat64(l.IsPirMotionDetected.Val), tags)
}

func (u *DatadogUnifi) batchProtectBridge(r report, d *unifi.ProtectDevices, b *unifi.ProtectBridge) {
	metricName := metricNamespace("protect_bridge")

	tags := tagMapToTags(cleanTags(map[string]string{
		"source": d.SourceName,
		"name":   b.Name,
		"mac":    b.MAC,
		"guid":   b.GUID,
	}))

	_ = r.reportGauge(metricName("present"), 1.0, tags)
	_ = r.reportGauge(metricName("state"), unifi.ProtectStateValue(b.State), tags)
	_ = r.reportGauge(metricName("client_count"), float64(len(b.Clients)), tags)
	_ = r.reportGauge(metricName("max_clients"), b.MaxClients.Val, tags)
}

func (u *DatadogUnifi) batchProtectLinkStation(r report, d *unifi.ProtectDevices, ls *unifi.ProtectLinkStation) {
	metricName := metricNamespace("protect_link_station")

	tags := tagMapToTags(cleanTags(map[string]string{
		"source": d.SourceName,
		"name":   ls.Name,
		"mac":    ls.MAC,
		"guid":   ls.GUID,
	}))

	_ = r.reportGauge(metricName("present"), 1.0, tags)
	_ = r.reportGauge(metricName("state"), unifi.ProtectStateValue(ls.State), tags)
	_ = r.reportGauge(metricName("is_alarm_hub"), boolToFloat64(ls.IsAlarmHub.Val), tags)

	if ls.AlarmHub != nil {
		_ = r.reportGauge(metricName("armed"), boolToFloat64(ls.AlarmHub.Armed.Val), tags)
	}
}
