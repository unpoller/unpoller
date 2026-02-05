package influxunifi

import "github.com/unpoller/unifi/v5"

// udbT is used as a name for printed/logged counters.
const udbT = item("UDB")

// batchUDB generates datapoints for UDB (UniFi Device Bridge) devices.
// UDB-Switch is a hybrid device combining switch ports with WiFi 7
// wireless bridge capability.
func (u *InfluxUnifi) batchUDB(r report, s *unifi.UDB) {
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := map[string]string{
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"source":    s.SourceName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
	}

	fields := Combine(
		u.batchUSWstat(s.Stat.Sw),
		u.batchSysStats(s.SysStats, s.SystemStats),
		map[string]any{
			"guest-num_sta":        s.GuestNumSta.Val,
			"ip":                   s.IP,
			"bytes":                s.Bytes.Val,
			"fan_level":            s.FanLevel.Val,
			"general_temperature":  s.GeneralTemperature.Val,
			"last_seen":            s.LastSeen.Val,
			"rx_bytes":             s.RxBytes.Val,
			"tx_bytes":             s.TxBytes.Val,
			"uptime":               s.Uptime.Val,
			"state":                s.State.Val,
			"user-num_sta":         s.UserNumSta.Val,
			"num_sta":              s.NumSta.Val,
			"upgradeable":          s.Upgradable.Val,
			"guest-wlan-num_sta":   s.GuestWlanNumSta.Val,
			"user-wlan-num_sta":    s.UserWlanNumSta.Val,
			"satisfaction":         s.Satisfaction.Val,
			"total_max_power":      s.TotalMaxPower.Val,
			"uplink_speed":         s.Uplink.Speed.Val,
			"uplink_max_speed":     s.Uplink.MaxSpeed.Val,
			"uplink_latency":       s.Uplink.Latency.Val,
			"uplink_uptime":        s.Uplink.Uptime.Val,
		})

	r.addCount(udbT)
	r.send(&metric{Table: "udb", Tags: tags, Fields: fields})

	// Port table (reuse USW function)
	u.batchPortTable(r, tags, s.PortTable)

	// Radio table (reuse UAP functions)
	u.processRadTable(r, tags, s.RadioTable, s.RadioTableStats)

	// VAP table (reuse UAP function)
	u.processVAPTable(r, tags, s.VapTable)
}
