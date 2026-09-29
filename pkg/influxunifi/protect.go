package influxunifi

import "github.com/unpoller/unifi/v6"

// batchProtectDevices emits one measurement per Protect device kind (sensor, camera, light,
// bridge, link station, NVR). Each kind is guarded independently because GetProtectDevices
// leaves the field for a failing endpoint nil/empty rather than failing the whole poll.
func (u *InfluxUnifi) batchProtectDevices(r report, d *unifi.ProtectDevices) {
	if d == nil {
		return
	}

	if d.NVR != nil {
		nvr := d.NVR
		fields := map[string]any{
			"state": unifi.ProtectStateValue(nvr.State),
		}

		if nvr.ArmMode != nil {
			fields["arm_status"] = nvr.ArmMode.Status
		}

		r.send(&metric{
			Table: "protect_nvr",
			Tags: map[string]string{
				"source": d.SourceName,
				"name":   nvr.Name,
				"mac":    nvr.MAC,
				"guid":   nvr.GUID,
			},
			Fields: fields,
		})
	}

	for _, s := range d.Sensors {
		if s == nil {
			continue
		}

		fields := map[string]any{
			"state":              unifi.ProtectStateValue(s.State),
			"is_opened":          s.IsOpened.Val,
			"is_motion_detected": s.IsMotionDetected.Val,
		}

		if s.BatteryStatus != nil {
			fields["battery_percent"] = s.BatteryStatus.Percentage.Val
			fields["battery_low"] = s.BatteryStatus.IsLow.Val
		}

		if temp, ok := s.TemperatureReading(); ok {
			fields["temperature"] = temp
		}

		if humidity, ok := s.HumidityReading(); ok {
			fields["humidity"] = humidity
		}

		if light, ok := s.LightReading(); ok {
			fields["light"] = light
		}

		if s.AirQuality != nil {
			aq := s.AirQuality
			channels := []struct {
				name    string
				reading *unifi.ProtectSensorStatValue
			}{
				{"aqi", aq.AQI},
				{"co2_ppm", aq.CO2},
				{"tvoc", aq.TVOC},
				{"voc", aq.VOC},
				{"vape", aq.Vape},
				{"pm1_0", aq.PM1p0},
				{"pm2_5", aq.PM2p5},
				{"pm4_0", aq.PM4p0},
				{"pm10", aq.PM10p0},
			}

			for _, ch := range channels {
				value, ok := ch.reading.Reading()
				if !ok {
					continue
				}

				fields[ch.name] = value
			}
		}

		r.send(&metric{
			Table: "protect_sensor",
			Tags: map[string]string{
				"source": d.SourceName,
				"name":   s.Name,
				"mac":    s.MAC,
				"guid":   s.GUID,
			},
			Fields: fields,
		})
	}

	for _, c := range d.Cameras {
		if c == nil {
			continue
		}

		r.send(&metric{
			Table: "protect_camera",
			Tags: map[string]string{
				"source": d.SourceName,
				"name":   c.Name,
				"mac":    c.MAC,
				"guid":   c.GUID,
			},
			Fields: map[string]any{
				"state":           unifi.ProtectStateValue(c.State),
				"mic_enabled":     c.IsMicEnabled.Val,
				"video_mode":      c.VideoMode,
				"has_package_cam": c.HasPackageCamera.Val,
			},
		})
	}

	for _, l := range d.Lights {
		if l == nil {
			continue
		}

		r.send(&metric{
			Table: "protect_light",
			Tags: map[string]string{
				"source": d.SourceName,
				"name":   l.Name,
				"mac":    l.MAC,
				"guid":   l.GUID,
			},
			Fields: map[string]any{
				"state":                  unifi.ProtectStateValue(l.State),
				"is_light_on":            l.IsLightOn.Val,
				"is_dark":                l.IsDark.Val,
				"is_pir_motion_detected": l.IsPirMotionDetected.Val,
			},
		})
	}

	for _, b := range d.Bridges {
		if b == nil {
			continue
		}

		r.send(&metric{
			Table: "protect_bridge",
			Tags: map[string]string{
				"source": d.SourceName,
				"name":   b.Name,
				"mac":    b.MAC,
				"guid":   b.GUID,
			},
			Fields: map[string]any{
				"state":        unifi.ProtectStateValue(b.State),
				"client_count": len(b.Clients),
				"max_clients":  b.MaxClients.Val,
			},
		})
	}

	for _, ls := range d.LinkStations {
		if ls == nil {
			continue
		}

		fields := map[string]any{
			"state":        unifi.ProtectStateValue(ls.State),
			"is_alarm_hub": ls.IsAlarmHub.Val,
		}

		if ls.AlarmHub != nil {
			fields["armed"] = ls.AlarmHub.Armed.Val
		}

		r.send(&metric{
			Table: "protect_link_station",
			Tags: map[string]string{
				"source": d.SourceName,
				"name":   ls.Name,
				"mac":    ls.MAC,
				"guid":   ls.GUID,
			},
			Fields: fields,
		})
	}
}
