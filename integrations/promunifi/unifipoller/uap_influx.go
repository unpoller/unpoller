package unifipoller

import (
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"golift.io/unifi"
)

// UAPPoints generates Wireless-Access-Point datapoints for InfluxDB.
// These points can be passed directly to influx.
func UAPPoints(u *unifi.UAP, now time.Time) ([]*influx.Point, error) {
	tags := map[string]string{
		"id":                    u.ID,
		"mac":                   u.Mac,
		"device_type":           u.Stat.O,
		"device_oid":            u.Stat.Oid,
		"device_ap":             u.Stat.Ap,
		"site_id":               u.SiteID,
		"site_name":             u.SiteName,
		"name":                  u.Name,
		"adopted":               u.Adopted.Txt,
		"cfgversion":            u.Cfgversion,
		"config_network_ip":     u.ConfigNetwork.IP,
		"config_network_type":   u.ConfigNetwork.Type,
		"connect_request_ip":    u.ConnectRequestIP,
		"device_id":             u.DeviceID,
		"has_eth1":              u.HasEth1.Txt,
		"inform_ip":             u.InformIP,
		"isolated":              u.Isolated.Txt,
		"known_cfgversion":      u.KnownCfgversion,
		"model":                 u.Model,
		"outdoor_mode_override": u.OutdoorModeOverride,
		"serial":                u.Serial,
		"type":                  u.Type,
		"vwireEnabled":          u.VwireEnabled.Txt,
	}
	fields := map[string]interface{}{
		"ip":                          u.IP,
		"bytes":                       u.Bytes.Val,
		"bytes_d":                     u.BytesD.Val,
		"bytes_r":                     u.BytesR.Val,
		"last_seen":                   u.LastSeen.Val,
		"rx_bytes":                    u.RxBytes.Val,
		"rx_bytes-d":                  u.RxBytesD.Val,
		"tx_bytes":                    u.TxBytes.Val,
		"tx_bytes-d":                  u.TxBytesD.Val,
		"uptime":                      u.Uptime.Val,
		"scanning":                    u.Scanning.Val,
		"spectrum_scanning":           u.SpectrumScanning.Val,
		"roll_upgrade":                u.Rollupgrade.Val,
		"state":                       u.State,
		"upgradable":                  u.Upgradable.Val,
		"user-num_sta":                u.UserNumSta,
		"guest-num_sta":               u.GuestNumSta,
		"version":                     u.Version,
		"loadavg_1":                   u.SysStats.Loadavg1.Val,
		"loadavg_5":                   u.SysStats.Loadavg5.Val,
		"loadavg_15":                  u.SysStats.Loadavg15.Val,
		"mem_buffer":                  u.SysStats.MemBuffer.Val,
		"mem_total":                   u.SysStats.MemTotal.Val,
		"mem_used":                    u.SysStats.MemUsed.Val,
		"cpu":                         u.SystemStats.CPU.Val,
		"mem":                         u.SystemStats.Mem.Val,
		"system_uptime":               u.SystemStats.Uptime.Val,
		"stat_guest-wifi0-rx_packets": u.Stat.GuestWifi0RxPackets.Val,
		"stat_guest-wifi1-rx_packets": u.Stat.GuestWifi1RxPackets.Val,
		"stat_user-wifi1-rx_packets":  u.Stat.UserWifi1RxPackets.Val,
		"stat_user-wifi0-rx_packets":  u.Stat.UserWifi0RxPackets.Val,
		"stat_user-rx_packets":        u.Stat.UserRxPackets.Val,
		"stat_guest-rx_packets":       u.Stat.GuestRxPackets.Val,
		"stat_wifi0-rx_packets":       u.Stat.Wifi0RxPackets.Val,
		"stat_wifi1-rx_packets":       u.Stat.Wifi1RxPackets.Val,
		"stat_rx_packets":             u.Stat.RxPackets.Val,
		"stat_guest-wifi0-rx_bytes":   u.Stat.GuestWifi0RxBytes.Val,
		"stat_guest-wifi1-rx_bytes":   u.Stat.GuestWifi1RxBytes.Val,
		"stat_user-wifi1-rx_bytes":    u.Stat.UserWifi1RxBytes.Val,
		"stat_user-wifi0-rx_bytes":    u.Stat.UserWifi0RxBytes.Val,
		"stat_user-rx_bytes":          u.Stat.UserRxBytes.Val,
		"stat_guest-rx_bytes":         u.Stat.GuestRxBytes.Val,
		"stat_wifi0-rx_bytes":         u.Stat.Wifi0RxBytes.Val,
		"stat_wifi1-rx_bytes":         u.Stat.Wifi1RxBytes.Val,
		"stat_rx_bytes":               u.Stat.RxBytes.Val,
		"stat_guest-wifi0-rx_errors":  u.Stat.GuestWifi0RxErrors.Val,
		"stat_guest-wifi1-rx_errors":  u.Stat.GuestWifi1RxErrors.Val,
		"stat_user-wifi1-rx_errors":   u.Stat.UserWifi1RxErrors.Val,
		"stat_user-wifi0-rx_errors":   u.Stat.UserWifi0RxErrors.Val,
		"stat_user-rx_errors":         u.Stat.UserRxErrors.Val,
		"stat_guest-rx_errors":        u.Stat.GuestRxErrors.Val,
		"stat_wifi0-rx_errors":        u.Stat.Wifi0RxErrors.Val,
		"stat_wifi1-rx_errors":        u.Stat.Wifi1RxErrors.Val,
		"stat_rx_errors":              u.Stat.RxErrors.Val,
		"stat_guest-wifi0-rx_dropped": u.Stat.GuestWifi0RxDropped.Val,
		"stat_guest-wifi1-rx_dropped": u.Stat.GuestWifi1RxDropped.Val,
		"stat_user-wifi1-rx_dropped":  u.Stat.UserWifi1RxDropped.Val,
		"stat_user-wifi0-rx_dropped":  u.Stat.UserWifi0RxDropped.Val,
		"stat_user-rx_dropped":        u.Stat.UserRxDropped.Val,
		"stat_guest-rx_dropped":       u.Stat.GuestRxDropped.Val,
		"stat_wifi0-rx_dropped":       u.Stat.Wifi0RxDropped.Val,
		"stat_wifi1-rx_dropped":       u.Stat.Wifi1RxDropped.Val,
		"stat_rx_dropped":             u.Stat.RxDropped.Val,
		"stat_guest-wifi0-rx_crypts":  u.Stat.GuestWifi0RxCrypts.Val,
		"stat_guest-wifi1-rx_crypts":  u.Stat.GuestWifi1RxCrypts.Val,
		"stat_user-wifi1-rx_crypts":   u.Stat.UserWifi1RxCrypts.Val,
		"stat_user-wifi0-rx_crypts":   u.Stat.UserWifi0RxCrypts.Val,
		"stat_user-rx_crypts":         u.Stat.UserRxCrypts.Val,
		"stat_guest-rx_crypts":        u.Stat.GuestRxCrypts.Val,
		"stat_wifi0-rx_crypts":        u.Stat.Wifi0RxCrypts.Val,
		"stat_wifi1-rx_crypts":        u.Stat.Wifi1RxCrypts.Val,
		"stat_rx_crypts":              u.Stat.RxCrypts.Val,
		"stat_guest-wifi0-rx_frags":   u.Stat.GuestWifi0RxFrags.Val,
		"stat_guest-wifi1-rx_frags":   u.Stat.GuestWifi1RxFrags.Val,
		"stat_user-wifi1-rx_frags":    u.Stat.UserWifi1RxFrags.Val,
		"stat_user-wifi0-rx_frags":    u.Stat.UserWifi0RxFrags.Val,
		"stat_user-rx_frags":          u.Stat.UserRxFrags.Val,
		"stat_guest-rx_frags":         u.Stat.GuestRxFrags.Val,
		"stat_wifi0-rx_frags":         u.Stat.Wifi0RxFrags.Val,
		"stat_wifi1-rx_frags":         u.Stat.Wifi1RxFrags.Val,
		"stat_rx_frags":               u.Stat.RxFrags.Val,
		"stat_guest-wifi0-tx_packets": u.Stat.GuestWifi0TxPackets.Val,
		"stat_guest-wifi1-tx_packets": u.Stat.GuestWifi1TxPackets.Val,
		"stat_user-wifi1-tx_packets":  u.Stat.UserWifi1TxPackets.Val,
		"stat_user-wifi0-tx_packets":  u.Stat.UserWifi0TxPackets.Val,
		"stat_user-tx_packets":        u.Stat.UserTxPackets.Val,
		"stat_guest-tx_packets":       u.Stat.GuestTxPackets.Val,
		"stat_wifi0-tx_packets":       u.Stat.Wifi0TxPackets.Val,
		"stat_wifi1-tx_packets":       u.Stat.Wifi1TxPackets.Val,
		"stat_tx_packets":             u.Stat.TxPackets.Val,
		"stat_guest-wifi0-tx_bytes":   u.Stat.GuestWifi0TxBytes.Val,
		"stat_guest-wifi1-tx_bytes":   u.Stat.GuestWifi1TxBytes.Val,
		"stat_user-wifi1-tx_bytes":    u.Stat.UserWifi1TxBytes.Val,
		"stat_user-wifi0-tx_bytes":    u.Stat.UserWifi0TxBytes.Val,
		"stat_user-tx_bytes":          u.Stat.UserTxBytes.Val,
		"stat_guest-tx_bytes":         u.Stat.GuestTxBytes.Val,
		"stat_wifi0-tx_bytes":         u.Stat.Wifi0TxBytes.Val,
		"stat_wifi1-tx_bytes":         u.Stat.Wifi1TxBytes.Val,
		"stat_tx_bytes":               u.Stat.TxBytes.Val,
		"stat_guest-wifi0-tx_errors":  u.Stat.GuestWifi0TxErrors.Val,
		"stat_guest-wifi1-tx_errors":  u.Stat.GuestWifi1TxErrors.Val,
		"stat_user-wifi1-tx_errors":   u.Stat.UserWifi1TxErrors.Val,
		"stat_user-wifi0-tx_errors":   u.Stat.UserWifi0TxErrors.Val,
		"stat_user-tx_errors":         u.Stat.UserTxErrors.Val,
		"stat_guest-tx_errors":        u.Stat.GuestTxErrors.Val,
		"stat_wifi0-tx_errors":        u.Stat.Wifi0TxErrors.Val,
		"stat_wifi1-tx_errors":        u.Stat.Wifi1TxErrors.Val,
		"stat_tx_errors":              u.Stat.TxErrors.Val,
		"stat_guest-wifi0-tx_dropped": u.Stat.GuestWifi0TxDropped.Val,
		"stat_guest-wifi1-tx_dropped": u.Stat.GuestWifi1TxDropped.Val,
		"stat_user-wifi1-tx_dropped":  u.Stat.UserWifi1TxDropped.Val,
		"stat_user-wifi0-tx_dropped":  u.Stat.UserWifi0TxDropped.Val,
		"stat_user-tx_dropped":        u.Stat.UserTxDropped.Val,
		"stat_guest-tx_dropped":       u.Stat.GuestTxDropped.Val,
		"stat_wifi0-tx_dropped":       u.Stat.Wifi0TxDropped.Val,
		"stat_wifi1-tx_dropped":       u.Stat.Wifi1TxDropped.Val,
		"stat_tx_dropped":             u.Stat.TxDropped.Val,
		"stat_guest-wifi0-tx_retries": u.Stat.GuestWifi0TxRetries.Val,
		"stat_guest-wifi1-tx_retries": u.Stat.GuestWifi1TxRetries.Val,
		"stat_user-wifi1-tx_retries":  u.Stat.UserWifi1TxRetries.Val,
		"stat_user-wifi0-tx_retries":  u.Stat.UserWifi0TxRetries.Val,
		"stat_user-tx_retries":        u.Stat.UserTxRetries.Val,
		"stat_guest-tx_retries":       u.Stat.GuestTxRetries.Val,
		"stat_wifi0-tx_retries":       u.Stat.Wifi0TxRetries.Val,
		"stat_wifi1-tx_retries":       u.Stat.Wifi1TxRetries.Val,
	}
	pt, err := influx.NewPoint("uap", tags, fields, now)
	if err != nil {
		return nil, err
	}
	points := []*influx.Point{pt}

	tags = make(map[string]string)
	fields = make(map[string]interface{})
	// Loop each virtual AP (ESSID) and extract data for it
	// from radio_tables and radio_table_stats.
	for _, s := range u.VapTable {
		tags["device_name"] = u.Name
		tags["device_id"] = u.ID
		tags["device_mac"] = u.Mac
		tags["site_name"] = u.SiteName
		tags["ap_mac"] = s.ApMac
		tags["bssid"] = s.Bssid
		tags["id"] = s.ID
		tags["name"] = s.Name
		tags["radio_name"] = s.RadioName
		tags["wlanconf_id"] = s.WlanconfID
		tags["essid"] = s.Essid
		tags["site_id"] = s.SiteID
		tags["usage"] = s.Usage
		tags["state"] = s.State
		tags["is_guest"] = s.IsGuest.Txt
		tags["is_wep"] = s.IsWep.Txt

		fields["ccq"] = s.Ccq
		fields["extchannel"] = s.Extchannel
		fields["mac_filter_rejections"] = s.MacFilterRejections
		fields["num_satisfaction_sta"] = s.NumSatisfactionSta.Val
		fields["avg_client_signal"] = s.AvgClientSignal.Val
		fields["satisfaction"] = s.Satisfaction.Val
		fields["satisfaction_now"] = s.SatisfactionNow.Val
		fields["rx_bytes"] = s.RxBytes.Val
		fields["rx_crypts"] = s.RxCrypts.Val
		fields["rx_dropped"] = s.RxDropped.Val
		fields["rx_errors"] = s.RxErrors.Val
		fields["rx_frags"] = s.RxFrags.Val
		fields["rx_nwids"] = s.RxNwids.Val
		fields["rx_packets"] = s.RxPackets.Val
		fields["tx_bytes"] = s.TxBytes.Val
		fields["tx_dropped"] = s.TxDropped.Val
		fields["tx_errors"] = s.TxErrors.Val
		fields["tx_packets"] = s.TxPackets.Val
		fields["tx_power"] = s.TxPower.Val
		fields["tx_retries"] = s.TxRetries.Val
		fields["tx_combined_retries"] = s.TxCombinedRetries.Val
		fields["tx_data_mpdu_bytes"] = s.TxDataMpduBytes.Val
		fields["tx_rts_retries"] = s.TxRtsRetries.Val
		fields["tx_success"] = s.TxSuccess.Val
		fields["tx_total"] = s.TxTotal.Val
		fields["tx_tcp_goodbytes"] = s.TxTCPStats.Goodbytes.Val
		fields["tx_tcp_lat_avg"] = s.TxTCPStats.LatAvg.Val
		fields["tx_tcp_lat_max"] = s.TxTCPStats.LatMax.Val
		fields["tx_tcp_lat_min"] = s.TxTCPStats.LatMin.Val
		fields["rx_tcp_goodbytes"] = s.RxTCPStats.Goodbytes.Val
		fields["rx_tcp_lat_avg"] = s.RxTCPStats.LatAvg.Val
		fields["rx_tcp_lat_max"] = s.RxTCPStats.LatMax.Val
		fields["rx_tcp_lat_min"] = s.RxTCPStats.LatMin.Val
		fields["wifi_tx_latency_mov_avg"] = s.WifiTxLatencyMov.Avg.Val
		fields["wifi_tx_latency_mov_max"] = s.WifiTxLatencyMov.Max.Val
		fields["wifi_tx_latency_mov_min"] = s.WifiTxLatencyMov.Min.Val
		fields["wifi_tx_latency_mov_total"] = s.WifiTxLatencyMov.Total.Val
		fields["wifi_tx_latency_mov_cuont"] = s.WifiTxLatencyMov.TotalCount.Val

		for _, p := range u.RadioTable {
			if p.Name != s.RadioName {
				continue
			}
			tags["wlangroup_id"] = p.WlangroupID
			tags["channel"] = p.Channel.Txt
			tags["radio"] = p.Radio
			fields["current_antenna_gain"] = p.CurrentAntennaGain.Val
			fields["ht"] = p.Ht
			fields["max_txpower"] = p.MaxTxpower.Val
			fields["min_rssi_enabled"] = p.MinRssiEnabled.Val
			fields["min_txpower"] = p.MinTxpower.Val
			fields["nss"] = p.Nss.Val
			fields["radio_caps"] = p.RadioCaps.Val
			fields["tx_power"] = p.TxPower.Val
		}

		for _, p := range u.RadioTableStats {
			if p.Name != s.RadioName {
				continue
			}
			fields["ast_be_xmit"] = p.AstBeXmit.Val
			fields["channel"] = p.Channel.Val
			fields["cu_self_rx"] = p.CuSelfRx.Val
			fields["cu_self_tx"] = p.CuSelfTx.Val
			fields["cu_total"] = p.CuTotal.Val
			fields["extchannel"] = p.Extchannel.Val
			fields["gain"] = p.Gain.Val
			fields["guest-num_sta"] = p.GuestNumSta.Val
			fields["num_sta"] = p.NumSta.Val
			fields["radio"] = p.Radio
			fields["tx_packets"] = p.TxPackets.Val
			fields["tx_power"] = p.TxPower.Val
			fields["tx_retries"] = p.TxRetries.Val
			fields["user-num_sta"] = p.UserNumSta.Val
		}

		pt, err := influx.NewPoint("uap_vaps", tags, fields, now)
		if err != nil {
			return points, err
		}
		points = append(points, pt)
	}
	return points, nil
}
