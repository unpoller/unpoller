package influxunifi

import (
	"golift.io/unifi"
)

// batchUAP generates Wireless-Access-Point datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchUAP(r report, s *unifi.UAP) {
	if !s.Adopted.Val || s.Locating.Val {
		return
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
	fields := Combine(u.processUAPstats(s.Stat.Ap), u.batchSysStats(s.SysStats, s.SystemStats))
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
	u.processRadTable(r, tags, s.RadioTable, s.RadioTableStats)
	u.processVAPTable(r, tags, s.VapTable)
}

func (u *InfluxUnifi) processUAPstats(ap *unifi.Ap) map[string]interface{} {
	if ap == nil {
		return map[string]interface{}{}
	}
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

// processVAPTable creates points for Wifi Radios. This works with several types of UAP-capable devices.
func (u *InfluxUnifi) processVAPTable(r report, t map[string]string, vt unifi.VapTable) {
	for _, s := range vt {
		tags := map[string]string{
			"device_name": t["name"],
			"site_name":   t["site_name"],
			"ap_mac":      s.ApMac,
			"bssid":       s.Bssid,
			"id":          s.ID,
			"name":        s.Name,
			"radio_name":  s.RadioName,
			"radio":       s.Radio,
			"essid":       s.Essid,
			"site_id":     s.SiteID,
			"usage":       s.Usage,
			"state":       s.State,
			"is_guest":    s.IsGuest.Txt,
		}
		fields := map[string]interface{}{
			"ccq":                       s.Ccq,
			"mac_filter_rejections":     s.MacFilterRejections,
			"num_satisfaction_sta":      s.NumSatisfactionSta.Val,
			"avg_client_signal":         s.AvgClientSignal.Val,
			"satisfaction":              s.Satisfaction.Val,
			"satisfaction_now":          s.SatisfactionNow.Val,
			"num_sta":                   s.NumSta,
			"channel":                   s.Channel.Val,
			"rx_bytes":                  s.RxBytes.Val,
			"rx_crypts":                 s.RxCrypts.Val,
			"rx_dropped":                s.RxDropped.Val,
			"rx_errors":                 s.RxErrors.Val,
			"rx_frags":                  s.RxFrags.Val,
			"rx_nwids":                  s.RxNwids.Val,
			"rx_packets":                s.RxPackets.Val,
			"tx_bytes":                  s.TxBytes.Val,
			"tx_dropped":                s.TxDropped.Val,
			"tx_errors":                 s.TxErrors.Val,
			"tx_packets":                s.TxPackets.Val,
			"tx_power":                  s.TxPower.Val,
			"tx_retries":                s.TxRetries.Val,
			"tx_combined_retries":       s.TxCombinedRetries.Val,
			"tx_data_mpdu_bytes":        s.TxDataMpduBytes.Val,
			"tx_rts_retries":            s.TxRtsRetries.Val,
			"tx_success":                s.TxSuccess.Val,
			"tx_total":                  s.TxTotal.Val,
			"tx_tcp_goodbytes":          s.TxTCPStats.Goodbytes.Val,
			"tx_tcp_lat_avg":            s.TxTCPStats.LatAvg.Val,
			"tx_tcp_lat_max":            s.TxTCPStats.LatMax.Val,
			"tx_tcp_lat_min":            s.TxTCPStats.LatMin.Val,
			"rx_tcp_goodbytes":          s.RxTCPStats.Goodbytes.Val,
			"rx_tcp_lat_avg":            s.RxTCPStats.LatAvg.Val,
			"rx_tcp_lat_max":            s.RxTCPStats.LatMax.Val,
			"rx_tcp_lat_min":            s.RxTCPStats.LatMin.Val,
			"wifi_tx_latency_mov_avg":   s.WifiTxLatencyMov.Avg.Val,
			"wifi_tx_latency_mov_max":   s.WifiTxLatencyMov.Max.Val,
			"wifi_tx_latency_mov_min":   s.WifiTxLatencyMov.Min.Val,
			"wifi_tx_latency_mov_total": s.WifiTxLatencyMov.Total.Val,
			"wifi_tx_latency_mov_cuont": s.WifiTxLatencyMov.TotalCount.Val,
		}
		r.send(&metric{Table: "uap_vaps", Tags: tags, Fields: fields})
	}
}

func (u *InfluxUnifi) processRadTable(r report, t map[string]string, rt unifi.RadioTable, rts unifi.RadioTableStats) {
	for _, p := range rt {
		tags := map[string]string{
			"device_name": t["name"],
			"site_name":   t["site_name"],
			"channel":     p.Channel.Txt,
			"radio":       p.Radio,
		}
		fields := map[string]interface{}{
			"current_antenna_gain": p.CurrentAntennaGain.Val,
			"ht":                   p.Ht.Txt,
			"max_txpower":          p.MaxTxpower.Val,
			"min_txpower":          p.MinTxpower.Val,
			"nss":                  p.Nss.Val,
			"radio_caps":           p.RadioCaps.Val,
		}
		for _, t := range rts {
			if t.Name == p.Name {
				fields["ast_be_xmit"] = t.AstBeXmit.Val
				fields["channel"] = t.Channel.Val
				fields["cu_self_rx"] = t.CuSelfRx.Val
				fields["cu_self_tx"] = t.CuSelfTx.Val
				fields["cu_total"] = t.CuTotal.Val
				fields["extchannel"] = t.Extchannel.Val
				fields["gain"] = t.Gain.Val
				fields["guest-num_sta"] = t.GuestNumSta.Val
				fields["num_sta"] = t.NumSta.Val
				fields["radio"] = t.Radio
				fields["tx_packets"] = t.TxPackets.Val
				fields["tx_power"] = t.TxPower.Val
				fields["tx_retries"] = t.TxRetries.Val
				fields["user-num_sta"] = t.UserNumSta.Val
				break
			}
		}
		r.send(&metric{Table: "uap_radios", Tags: tags, Fields: fields})
	}
}
