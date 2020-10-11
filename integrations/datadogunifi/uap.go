package datadogunifi

import (
	"github.com/unifi-poller/unifi"
)

// reportUAP generates Wireless-Access-Point datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *DatadogUnifi) reportUAP(r report, s *unifi.UAP) {
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := []string{
		tag("ip", s.IP),
		tag("mac", s.Mac),
		tag("site_name", s.SiteName),
		tag("source", s.SourceName),
		tag("name", s.Name),
		tag("version", s.Version),
		tag("model", s.Model),
		tag("serial", s.Serial),
		tag("type", s.Type),
	}

	metricName := metricNamespace("uap")

	u.reportUAPstats(s.Stat.Ap, r, metricName, tags)
	u.reportSysStats(r, metricName, s.SysStats, s.SystemStats, tags)

	data := map[string]float64{
		"bytes":         s.Bytes.Val,
		"last_seen":     s.LastSeen.Val,
		"rx_bytes":      s.RxBytes.Val,
		"tx_bytes":      s.TxBytes.Val,
		"uptime":        s.Uptime.Val,
		"user-num_sta":  s.UserNumSta.Val,
		"guest-num_sta": s.GuestNumSta.Val,
		"num_sta":       s.NumSta.Val,
	}
	reportGaugeForMap(r, metricName, data, tags)

	u.reportRadTable(r, s.Name, s.SiteName, s.SourceName, s.RadioTable, s.RadioTableStats)
	u.reportVAPTable(r, s.Name, s.SiteName, s.SourceName, s.VapTable)
	u.reportPortTable(r, s.Name, s.SiteName, s.SourceName, s.Type, s.PortTable)
}

func (u *DatadogUnifi) reportUAPstats(ap *unifi.Ap, r report, metricName func(string) string, tags []string) {
	if ap == nil {
		return
	}

	// Accumulative Statistics.
	data := map[string]float64{
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
	reportGaugeForMap(r, metricName, data, tags)
}

// reportVAPTable creates points for Wifi Radios. This works with several types of UAP-capable devices.
func (u *DatadogUnifi) reportVAPTable(r report, deviceName string, siteName string, source string, vt unifi.VapTable) { // nolint: funlen
	for _, s := range vt {
		tags := []string{
			tag("device_name", deviceName),
			tag("site_name", siteName),
			tag("source", source),
			tag("ap_mac", s.ApMac),
			tag("bssid", s.Bssid),
			tag("id", s.ID),
			tag("name", s.Name),
			tag("radio_name", s.RadioName),
			tag("radio", s.Radio),
			tag("essid", s.Essid),
			tag("site_id", s.SiteID),
			tag("usage", s.Usage),
			tag("state", s.State),
			tag("is_guest", s.IsGuest.Txt),
		}
		data := map[string]float64{
			"ccq":                       float64(s.Ccq),
			"mac_filter_rejections":     float64(s.MacFilterRejections),
			"num_satisfaction_sta":      s.NumSatisfactionSta.Val,
			"avg_client_signal":         s.AvgClientSignal.Val,
			"satisfaction":              s.Satisfaction.Val,
			"satisfaction_now":          s.SatisfactionNow.Val,
			"num_sta":                   float64(s.NumSta),
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

		metricName := metricNamespace("uap_vaps")
		reportGaugeForMap(r, metricName, data, tags)
	}
}

func (u *DatadogUnifi) reportRadTable(r report, deviceName string, siteName string, source string, rt unifi.RadioTable, rts unifi.RadioTableStats) {
	for _, p := range rt {
		tags := []string{
			tag("device_name", deviceName),
			tag("site_name", siteName),
			tag("source", source),
			tag("channel", p.Channel.Txt),
			tag("radio", p.Radio),
		}
		data := map[string]float64{
			"current_antenna_gain": p.CurrentAntennaGain.Val,
			"ht":                   p.Ht.Val,
			"max_txpower":          p.MaxTxpower.Val,
			"min_txpower":          p.MinTxpower.Val,
			"nss":                  p.Nss.Val,
			"radio_caps":           p.RadioCaps.Val,
		}

		for _, t := range rts {
			if t.Name == p.Name {
				data["ast_be_xmit"] = t.AstBeXmit.Val
				data["channel"] = t.Channel.Val
				data["cu_self_rx"] = t.CuSelfRx.Val
				data["cu_self_tx"] = t.CuSelfTx.Val
				data["cu_total"] = t.CuTotal.Val
				data["extchannel"] = t.Extchannel.Val
				data["gain"] = t.Gain.Val
				data["guest-num_sta"] = t.GuestNumSta.Val
				data["num_sta"] = t.NumSta.Val
				data["tx_packets"] = t.TxPackets.Val
				data["tx_power"] = t.TxPower.Val
				data["tx_retries"] = t.TxRetries.Val
				data["user-num_sta"] = t.UserNumSta.Val

				break
			}
		}

		metricName := metricNamespace("uap_radios")

		reportGaugeForMap(r, metricName, data, tags)
	}
}
