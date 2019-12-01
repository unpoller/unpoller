package influxunifi

import (
	"golift.io/unifi"
)

// batchUAP generates Wireless-Access-Point datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchUAP(r report, s *unifi.UAP) {
	if s.Stat.Ap == nil {
		s.Stat.Ap = &unifi.Ap{}
	}
	tags := map[string]string{
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
	}
	fields := map[string]interface{}{
		"ip":            s.IP,
		"bytes":         s.Bytes.Val,
		"last_seen":     s.LastSeen.Val,
		"rx_bytes":      s.RxBytes.Val,
		"tx_bytes":      s.TxBytes.Val,
		"uptime":        s.Uptime.Val,
		"state":         s.State,
		"user-num_sta":  int(s.UserNumSta.Val),
		"guest-num_sta": int(s.GuestNumSta.Val),
		"num_sta":       s.NumSta.Val,
		"loadavg_1":     s.SysStats.Loadavg1.Val,
		"loadavg_5":     s.SysStats.Loadavg5.Val,
		"loadavg_15":    s.SysStats.Loadavg15.Val,
		"mem_buffer":    s.SysStats.MemBuffer.Val,
		"mem_total":     s.SysStats.MemTotal.Val,
		"mem_used":      s.SysStats.MemUsed.Val,
		"cpu":           s.SystemStats.CPU.Val,
		"mem":           s.SystemStats.Mem.Val,
		"system_uptime": s.SystemStats.Uptime.Val,
		// Accumulative Statistics.
		"stat_user-rx_packets":  s.Stat.Ap.UserRxPackets.Val,
		"stat_guest-rx_packets": s.Stat.Ap.GuestRxPackets.Val,
		"stat_rx_packets":       s.Stat.Ap.RxPackets.Val,
		"stat_user-rx_bytes":    s.Stat.Ap.UserRxBytes.Val,
		"stat_guest-rx_bytes":   s.Stat.Ap.GuestRxBytes.Val,
		"stat_rx_bytes":         s.Stat.Ap.RxBytes.Val,
		"stat_user-rx_errors":   s.Stat.Ap.UserRxErrors.Val,
		"stat_guest-rx_errors":  s.Stat.Ap.GuestRxErrors.Val,
		"stat_rx_errors":        s.Stat.Ap.RxErrors.Val,
		"stat_user-rx_dropped":  s.Stat.Ap.UserRxDropped.Val,
		"stat_guest-rx_dropped": s.Stat.Ap.GuestRxDropped.Val,
		"stat_rx_dropped":       s.Stat.Ap.RxDropped.Val,
		"stat_user-rx_crypts":   s.Stat.Ap.UserRxCrypts.Val,
		"stat_guest-rx_crypts":  s.Stat.Ap.GuestRxCrypts.Val,
		"stat_rx_crypts":        s.Stat.Ap.RxCrypts.Val,
		"stat_user-rx_frags":    s.Stat.Ap.UserRxFrags.Val,
		"stat_guest-rx_frags":   s.Stat.Ap.GuestRxFrags.Val,
		"stat_rx_frags":         s.Stat.Ap.RxFrags.Val,
		"stat_user-tx_packets":  s.Stat.Ap.UserTxPackets.Val,
		"stat_guest-tx_packets": s.Stat.Ap.GuestTxPackets.Val,
		"stat_tx_packets":       s.Stat.Ap.TxPackets.Val,
		"stat_user-tx_bytes":    s.Stat.Ap.UserTxBytes.Val,
		"stat_guest-tx_bytes":   s.Stat.Ap.GuestTxBytes.Val,
		"stat_tx_bytes":         s.Stat.Ap.TxBytes.Val,
		"stat_user-tx_errors":   s.Stat.Ap.UserTxErrors.Val,
		"stat_guest-tx_errors":  s.Stat.Ap.GuestTxErrors.Val,
		"stat_tx_errors":        s.Stat.Ap.TxErrors.Val,
		"stat_user-tx_dropped":  s.Stat.Ap.UserTxDropped.Val,
		"stat_guest-tx_dropped": s.Stat.Ap.GuestTxDropped.Val,
		"stat_tx_dropped":       s.Stat.Ap.TxDropped.Val,
		"stat_user-tx_retries":  s.Stat.Ap.UserTxRetries.Val,
		"stat_guest-tx_retries": s.Stat.Ap.GuestTxRetries.Val,
	}
	r.send(&metric{Table: "uap", Tags: tags, Fields: fields})
	u.processVAPs(r, s.VapTable, s.RadioTable, s.RadioTableStats, s.Name, s.Mac, s.SiteName)
}

// processVAPs creates points for Wifi Radios. This works with several types of UAP-capable devices.
func (u *InfluxUnifi) processVAPs(r report, vt unifi.VapTable, rt unifi.RadioTable, rts unifi.RadioTableStats, name, mac, sitename string) {
	// Loop each virtual AP (ESSID) and extract data for it
	// from radio_tables and radio_table_stats.
	for _, s := range vt {
		tags := make(map[string]string)
		fields := make(map[string]interface{})
		tags["device_name"] = name
		tags["device_mac"] = mac
		tags["site_name"] = sitename
		tags["ap_mac"] = s.ApMac
		tags["bssid"] = s.Bssid
		tags["id"] = s.ID
		tags["name"] = s.Name
		tags["radio_name"] = s.RadioName
		tags["essid"] = s.Essid
		tags["site_id"] = s.SiteID
		tags["usage"] = s.Usage
		tags["state"] = s.State
		tags["is_guest"] = s.IsGuest.Txt

		fields["ccq"] = s.Ccq
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
			tags["channel"] = p.Channel.Txt
			tags["radio"] = p.Radio
			fields["current_antenna_gain"] = p.CurrentAntennaGain.Val
			fields["ht"] = p.Ht.Txt
			fields["max_txpower"] = p.MaxTxpower.Val
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
		r.send(&metric{Table: "uap_vaps", Tags: tags, Fields: fields})
	}
}
