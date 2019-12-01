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
	fields := Combine(u.processUAPstats(r, s.Stat.Ap), u.batchSysStats(r, s.SysStats, s.SystemStats))
	fields["ip"] = s.IP
	fields["bytes"] = s.Bytes.Val
	fields["last_seen"] = s.LastSeen.Val
	fields["rx_bytes"] = s.RxBytes.Val
	fields["tx_bytes"] = s.TxBytes.Val
	fields["uptime"] = s.Uptime.Val
	fields["state"] = s.State
	fields["user-num_sta"] = int(s.UserNumSta.Val)
	fields["guest-num_sta"] = int(s.GuestNumSta.Val)
	fields["num_sta"] = s.NumSta.Val
	r.send(&metric{Table: "uap", Tags: tags, Fields: fields})
	u.processVAPs(r, tags, s.VapTable, s.RadioTable, s.RadioTableStats)
}

func (u *InfluxUnifi) processUAPstats(r report, ap *unifi.Ap) map[string]interface{} {
	// Accumulative Statistics.
	return map[string]interface{}{
		"stat_user-rx_packets":  ap.UserRxPackets.Val,
		"stat_guest-rx_packets": ap.GuestRxPackets.Val,
		"stat_rx_packets":       ap.RxPackets.Val,
		"stat_user-rx_bytes":    ap.UserRxBytes.Val,
		"stat_guest-rx_bytes":   ap.GuestRxBytes.Val,
		"stat_rx_bytes":         ap.RxBytes.Val,
		"stat_user-rx_errors":   ap.UserRxErrors.Val,
		"stat_guest-rx_errors":  ap.GuestRxErrors.Val,
		"stat_rx_errors":        ap.RxErrors.Val,
		"stat_user-rx_dropped":  ap.UserRxDropped.Val,
		"stat_guest-rx_dropped": ap.GuestRxDropped.Val,
		"stat_rx_dropped":       ap.RxDropped.Val,
		"stat_user-rx_crypts":   ap.UserRxCrypts.Val,
		"stat_guest-rx_crypts":  ap.GuestRxCrypts.Val,
		"stat_rx_crypts":        ap.RxCrypts.Val,
		"stat_user-rx_frags":    ap.UserRxFrags.Val,
		"stat_guest-rx_frags":   ap.GuestRxFrags.Val,
		"stat_rx_frags":         ap.RxFrags.Val,
		"stat_user-tx_packets":  ap.UserTxPackets.Val,
		"stat_guest-tx_packets": ap.GuestTxPackets.Val,
		"stat_tx_packets":       ap.TxPackets.Val,
		"stat_user-tx_bytes":    ap.UserTxBytes.Val,
		"stat_guest-tx_bytes":   ap.GuestTxBytes.Val,
		"stat_tx_bytes":         ap.TxBytes.Val,
		"stat_user-tx_errors":   ap.UserTxErrors.Val,
		"stat_guest-tx_errors":  ap.GuestTxErrors.Val,
		"stat_tx_errors":        ap.TxErrors.Val,
		"stat_user-tx_dropped":  ap.UserTxDropped.Val,
		"stat_guest-tx_dropped": ap.GuestTxDropped.Val,
		"stat_tx_dropped":       ap.TxDropped.Val,
		"stat_user-tx_retries":  ap.UserTxRetries.Val,
		"stat_guest-tx_retries": ap.GuestTxRetries.Val,
	}
}

// processVAPs creates points for Wifi Radios. This works with several types of UAP-capable devices.
func (u *InfluxUnifi) processVAPs(r report, tags map[string]string, vt unifi.VapTable, rt unifi.RadioTable, rts unifi.RadioTableStats) {
	// Loop each virtual AP (ESSID) and extract data for it
	// from radio_tables and radio_table_stats.
	for _, s := range vt {
		t := make(map[string]string)      // tags
		f := make(map[string]interface{}) // fields
		t["device_name"] = tags["name"]
		t["site_name"] = tags["site_name"]
		t["ap_mac"] = s.ApMac
		t["bssid"] = s.Bssid
		t["id"] = s.ID
		t["name"] = s.Name
		t["radio_name"] = s.RadioName
		t["essid"] = s.Essid
		t["site_id"] = s.SiteID
		t["usage"] = s.Usage
		t["state"] = s.State
		t["is_guest"] = s.IsGuest.Txt

		f["ccq"] = s.Ccq
		f["mac_filter_rejections"] = s.MacFilterRejections
		f["num_satisfaction_sta"] = s.NumSatisfactionSta.Val
		f["avg_client_signal"] = s.AvgClientSignal.Val
		f["satisfaction"] = s.Satisfaction.Val
		f["satisfaction_now"] = s.SatisfactionNow.Val
		f["rx_bytes"] = s.RxBytes.Val
		f["rx_crypts"] = s.RxCrypts.Val
		f["rx_dropped"] = s.RxDropped.Val
		f["rx_errors"] = s.RxErrors.Val
		f["rx_frags"] = s.RxFrags.Val
		f["rx_nwids"] = s.RxNwids.Val
		f["rx_packets"] = s.RxPackets.Val
		f["tx_bytes"] = s.TxBytes.Val
		f["tx_dropped"] = s.TxDropped.Val
		f["tx_errors"] = s.TxErrors.Val
		f["tx_packets"] = s.TxPackets.Val
		f["tx_power"] = s.TxPower.Val
		f["tx_retries"] = s.TxRetries.Val
		f["tx_combined_retries"] = s.TxCombinedRetries.Val
		f["tx_data_mpdu_bytes"] = s.TxDataMpduBytes.Val
		f["tx_rts_retries"] = s.TxRtsRetries.Val
		f["tx_success"] = s.TxSuccess.Val
		f["tx_total"] = s.TxTotal.Val
		f["tx_tcp_goodbytes"] = s.TxTCPStats.Goodbytes.Val
		f["tx_tcp_lat_avg"] = s.TxTCPStats.LatAvg.Val
		f["tx_tcp_lat_max"] = s.TxTCPStats.LatMax.Val
		f["tx_tcp_lat_min"] = s.TxTCPStats.LatMin.Val
		f["rx_tcp_goodbytes"] = s.RxTCPStats.Goodbytes.Val
		f["rx_tcp_lat_avg"] = s.RxTCPStats.LatAvg.Val
		f["rx_tcp_lat_max"] = s.RxTCPStats.LatMax.Val
		f["rx_tcp_lat_min"] = s.RxTCPStats.LatMin.Val
		f["wifi_tx_latency_mov_avg"] = s.WifiTxLatencyMov.Avg.Val
		f["wifi_tx_latency_mov_max"] = s.WifiTxLatencyMov.Max.Val
		f["wifi_tx_latency_mov_min"] = s.WifiTxLatencyMov.Min.Val
		f["wifi_tx_latency_mov_total"] = s.WifiTxLatencyMov.Total.Val
		f["wifi_tx_latency_mov_cuont"] = s.WifiTxLatencyMov.TotalCount.Val

		// XXX: This is busted. It needs its own table....
		for _, p := range rt {
			if p.Name != s.RadioName {
				continue
			}
			t["channel"] = p.Channel.Txt
			t["radio"] = p.Radio
			f["current_antenna_gain"] = p.CurrentAntennaGain.Val
			f["ht"] = p.Ht.Txt
			f["max_txpower"] = p.MaxTxpower.Val
			f["min_txpower"] = p.MinTxpower.Val
			f["nss"] = p.Nss.Val
			f["radio_caps"] = p.RadioCaps.Val
			f["tx_power"] = p.TxPower.Val
		}

		for _, p := range rts {
			if p.Name != s.RadioName {
				continue
			}
			f["ast_be_xmit"] = p.AstBeXmit.Val
			f["channel"] = p.Channel.Val
			f["cu_self_rx"] = p.CuSelfRx.Val
			f["cu_self_tx"] = p.CuSelfTx.Val
			f["cu_total"] = p.CuTotal.Val
			f["extchannel"] = p.Extchannel.Val
			f["gain"] = p.Gain.Val
			f["guest-num_sta"] = p.GuestNumSta.Val
			f["num_sta"] = p.NumSta.Val
			f["radio"] = p.Radio
			f["tx_packets"] = p.TxPackets.Val
			f["tx_power"] = p.TxPower.Val
			f["tx_retries"] = p.TxRetries.Val
			f["user-num_sta"] = p.UserNumSta.Val
		}
		r.send(&metric{Table: "uap_vaps", Tags: t, Fields: f})
	}
}
