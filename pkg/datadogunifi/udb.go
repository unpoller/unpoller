package datadogunifi

import "github.com/unpoller/unifi/v5"

// udbT is used as a name for printed/logged counters.
const udbT = item("UDB")

// batchUDB generates datapoints for UDB (UniFi Device Bridge) devices.
// UDB-Switch is a hybrid device combining switch ports with WiFi 7
// wireless bridge capability.
func (u *DatadogUnifi) batchUDB(r report, s *unifi.UDB) {
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := cleanTags(map[string]string{
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"source":    s.SourceName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
		"ip":        s.IP,
	})

	data := CombineFloat64(
		u.batchUSWstat(s.Stat.Sw),
		u.batchSysStats(s.SysStats, s.SystemStats),
		map[string]float64{
			"guest_num_sta":        s.GuestNumSta.Val,
			"bytes":                s.Bytes.Val,
			"fan_level":            s.FanLevel.Val,
			"general_temperature":  s.GeneralTemperature.Val,
			"last_seen":            s.LastSeen.Val,
			"rx_bytes":             s.RxBytes.Val,
			"tx_bytes":             s.TxBytes.Val,
			"uptime":               s.Uptime.Val,
			"state":                s.State.Val,
			"user_num_sta":         s.UserNumSta.Val,
			"num_sta":              s.NumSta.Val,
			"upgradeable":          boolToFloat64(s.Upgradable.Val),
			"guest_wlan_num_sta":   s.GuestWlanNumSta.Val,
			"user_wlan_num_sta":    s.UserWlanNumSta.Val,
			"satisfaction":         s.Satisfaction.Val,
			"total_max_power":      s.TotalMaxPower.Val,
			"uplink_speed":         s.Uplink.Speed.Val,
			"uplink_max_speed":     s.Uplink.MaxSpeed.Val,
			"uplink_latency":       s.Uplink.Latency.Val,
			"uplink_uptime":        s.Uplink.Uptime.Val,
		})

	r.addCount(udbT)

	metricName := metricNamespace("udb")

	reportGaugeForFloat64Map(r, metricName, data, tags)

	// Port table (reuse USW function)
	u.batchPortTable(r, tags, s.PortTable)

	// Radio table (reuse UAP functions)
	u.processRadTable(r, tags, s.RadioTable, s.RadioTableStats)

	// VAP table (reuse UAP function)
	u.processVAPTable(r, tags, s.VapTable)
}
