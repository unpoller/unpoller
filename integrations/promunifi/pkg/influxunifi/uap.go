package influxunifi

import (
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"golift.io/unifi"
)

// UAPPoints generates Wireless-Access-Point datapoints for InfluxDB.
// These points can be passed directly to influx.
func UAPPoints(u *unifi.UAP, now time.Time) ([]*influx.Point, error) {
	if u.Stat.Ap == nil {
		u.Stat.Ap = &unifi.Ap{}
	}
	tags := map[string]string{
		"id":                  u.ID,
		"ip":                  u.IP,
		"mac":                 u.Mac,
		"device_type":         u.Stat.Ap.O,
		"device_oid":          u.Stat.Ap.Oid,
		"device_ap":           u.Stat.Ap.Ap,
		"site_id":             u.SiteID,
		"site_name":           u.SiteName,
		"name":                u.Name,
		"adopted":             u.Adopted.Txt,
		"cfgversion":          u.Cfgversion,
		"config_network_ip":   u.ConfigNetwork.IP,
		"config_network_type": u.ConfigNetwork.Type,
		"connect_request_ip":  u.ConnectRequestIP,
		"device_id":           u.DeviceID,
		"has_eth1":            u.HasEth1.Txt,
		"inform_ip":           u.InformIP,
		"known_cfgversion":    u.KnownCfgversion,
		"model":               u.Model,
		"serial":              u.Serial,
		"type":                u.Type,
	}
	fields := map[string]interface{}{
		"ip":            u.IP,
		"bytes":         u.Bytes.Val,
		"last_seen":     u.LastSeen.Val,
		"rx_bytes":      u.RxBytes.Val,
		"tx_bytes":      u.TxBytes.Val,
		"uptime":        u.Uptime.Val,
		"state":         u.State,
		"user-num_sta":  int(u.UserNumSta.Val),
		"guest-num_sta": int(u.GuestNumSta.Val),
		"num_sta":       u.NumSta.Val,
		"version":       u.Version,
		"loadavg_1":     u.SysStats.Loadavg1.Val,
		"loadavg_5":     u.SysStats.Loadavg5.Val,
		"loadavg_15":    u.SysStats.Loadavg15.Val,
		"mem_buffer":    u.SysStats.MemBuffer.Val,
		"mem_total":     u.SysStats.MemTotal.Val,
		"mem_used":      u.SysStats.MemUsed.Val,
		"cpu":           u.SystemStats.CPU.Val,
		"mem":           u.SystemStats.Mem.Val,
		"system_uptime": u.SystemStats.Uptime.Val,
		// Accumulative Statistics.
		"stat_user-rx_packets":  u.Stat.Ap.UserRxPackets.Val,
		"stat_guest-rx_packets": u.Stat.Ap.GuestRxPackets.Val,
		"stat_rx_packets":       u.Stat.Ap.RxPackets.Val,
		"stat_user-rx_bytes":    u.Stat.Ap.UserRxBytes.Val,
		"stat_guest-rx_bytes":   u.Stat.Ap.GuestRxBytes.Val,
		"stat_rx_bytes":         u.Stat.Ap.RxBytes.Val,
		"stat_user-rx_errors":   u.Stat.Ap.UserRxErrors.Val,
		"stat_guest-rx_errors":  u.Stat.Ap.GuestRxErrors.Val,
		"stat_rx_errors":        u.Stat.Ap.RxErrors.Val,
		"stat_user-rx_dropped":  u.Stat.Ap.UserRxDropped.Val,
		"stat_guest-rx_dropped": u.Stat.Ap.GuestRxDropped.Val,
		"stat_rx_dropped":       u.Stat.Ap.RxDropped.Val,
		"stat_user-rx_crypts":   u.Stat.Ap.UserRxCrypts.Val,
		"stat_guest-rx_crypts":  u.Stat.Ap.GuestRxCrypts.Val,
		"stat_rx_crypts":        u.Stat.Ap.RxCrypts.Val,
		"stat_user-rx_frags":    u.Stat.Ap.UserRxFrags.Val,
		"stat_guest-rx_frags":   u.Stat.Ap.GuestRxFrags.Val,
		"stat_rx_frags":         u.Stat.Ap.RxFrags.Val,
		"stat_user-tx_packets":  u.Stat.Ap.UserTxPackets.Val,
		"stat_guest-tx_packets": u.Stat.Ap.GuestTxPackets.Val,
		"stat_tx_packets":       u.Stat.Ap.TxPackets.Val,
		"stat_user-tx_bytes":    u.Stat.Ap.UserTxBytes.Val,
		"stat_guest-tx_bytes":   u.Stat.Ap.GuestTxBytes.Val,
		"stat_tx_bytes":         u.Stat.Ap.TxBytes.Val,
		"stat_user-tx_errors":   u.Stat.Ap.UserTxErrors.Val,
		"stat_guest-tx_errors":  u.Stat.Ap.GuestTxErrors.Val,
		"stat_tx_errors":        u.Stat.Ap.TxErrors.Val,
		"stat_user-tx_dropped":  u.Stat.Ap.UserTxDropped.Val,
		"stat_guest-tx_dropped": u.Stat.Ap.GuestTxDropped.Val,
		"stat_tx_dropped":       u.Stat.Ap.TxDropped.Val,
		"stat_user-tx_retries":  u.Stat.Ap.UserTxRetries.Val,
		"stat_guest-tx_retries": u.Stat.Ap.GuestTxRetries.Val,
	}
	pt, err := influx.NewPoint("uap", tags, fields, now)
	if err != nil {
		return nil, err
	}
	morePoints, err := processVAPs(u.VapTable, u.RadioTable, u.RadioTableStats, u.Name, u.ID, u.Mac, u.SiteName, now)
	if err != nil {
		return nil, err
	}
	return append(morePoints, pt), nil
}

// processVAPs creates points for Wifi Radios. This works with several types of UAP-capable devices.
func processVAPs(vt unifi.VapTable, rt unifi.RadioTable, rts unifi.RadioTableStats, name, id, mac, sitename string, ts time.Time) ([]*influx.Point, error) {
	tags := make(map[string]string)
	fields := make(map[string]interface{})
	points := []*influx.Point{}

	// Loop each virtual AP (ESSID) and extract data for it
	// from radio_tables and radio_table_stats.
	for _, s := range vt {
		tags["device_name"] = name
		tags["device_id"] = id
		tags["device_mac"] = mac
		tags["site_name"] = sitename
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

		for _, p := range rt {
			if p.Name != s.RadioName {
				continue
			}
			tags["wlangroup_id"] = p.WlangroupID
			tags["channel"] = p.Channel.Txt
			tags["radio"] = p.Radio
			fields["current_antenna_gain"] = p.CurrentAntennaGain.Val
			fields["ht"] = p.Ht.Txt
			fields["max_txpower"] = p.MaxTxpower.Val
			fields["min_rssi_enabled"] = p.MinRssiEnabled.Val
			fields["min_txpower"] = p.MinTxpower.Val
			fields["nss"] = p.Nss.Val
			fields["radio_caps"] = p.RadioCaps.Val
			fields["tx_power"] = p.TxPower.Val
		}

		for _, p := range rts {
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

		pt, err := influx.NewPoint("uap_vaps", tags, fields, ts)
		if err != nil {
			return points, err
		}
		points = append(points, pt)
	}
	return points, nil
}
