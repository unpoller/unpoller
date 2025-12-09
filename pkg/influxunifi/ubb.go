package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// ubbT is used as a name for printed/logged counters.
const ubbT = item("UBB")

// batchUBB generates UBB (UniFi Building Bridge) datapoints for InfluxDB.
// UBB devices are point-to-point wireless bridges with dual radios:
//   - wifi0: 5GHz radio (802.11ac)
//   - terra2/wlan0/ad: 60GHz radio (802.11ad - Terragraph/WiGig)
func (u *InfluxUnifi) batchUBB(r report, s *unifi.UBB) { // nolint: funlen
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := map[string]string{
		"source":    s.SourceName,
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
	}

	sysStats := unifi.SysStats{}
	if s.SysStats != nil {
		sysStats = *s.SysStats
	}

	systemStats := unifi.SystemStats{}
	if s.SystemStats != nil {
		systemStats = *s.SystemStats
	}

	fields := Combine(
		u.batchSysStats(sysStats, systemStats),
		u.batchUBBstats(s.Stat),
		map[string]any{
			"source":           s.SourceName,
			"ip":               s.IP,
			"bytes":            s.Bytes.Val,
			"last_seen":        s.LastSeen.Val,
			"license_state":    s.LicenseState,
			"rx_bytes":         s.RxBytes.Val,
			"tx_bytes":         s.TxBytes.Val,
			"uptime":           s.Uptime.Val,
			"state":            s.State.Val,
			"user-num_sta":     s.UserNumSta.Val,
			"version":          s.Version,
			"uplink_speed":     s.Uplink.Speed.Val,
			"uplink_max_speed": s.Uplink.MaxSpeed.Val,
			"uplink_latency":   s.Uplink.Latency.Val,
			"uplink_uptime":    s.Uplink.Uptime.Val,
			"upgradeable":      s.Upgradable.Val,
		},
	)

	// Add UBB-specific P2P and link quality metrics
	if s.P2PStats != nil {
		fields["p2p_rx_rate"] = s.P2PStats.RXRate.Val
		fields["p2p_tx_rate"] = s.P2PStats.TXRate.Val
		fields["p2p_throughput"] = s.P2PStats.Throughput.Val
	}
	fields["link_quality"] = s.LinkQuality.Val
	fields["link_quality_current"] = s.LinkQualityCurrent.Val
	fields["link_capacity"] = s.LinkCapacity.Val

	r.addCount(ubbT)
	r.send(&metric{Table: "ubb", Tags: tags, Fields: fields})

	// Export VAP table (Virtual Access Point table - wireless interface stats)
	u.processVAPTable(r, tags, s.VapTable)

	// Export Radio tables (includes 5GHz wifi0 and 60GHz terra2/ad radios)
	u.processRadTable(r, tags, s.RadioTable, s.RadioTableStats)
}

// batchUBBstats generates UBB-specific statistics separated by radio.
// This includes metrics for total, wifi0 (5GHz), terra2 (60GHz), and user-specific stats.
func (u *InfluxUnifi) batchUBBstats(stat *unifi.UBBStat) map[string]any {
	if stat == nil || stat.Bb == nil {
		return map[string]any{}
	}

	bb := stat.Bb

	// Total aggregated stats across both radios
	return map[string]any{
		"stat_bytes":                  bb.Bytes.Val,
		"stat_duration":               bb.Duration.Val,
		"stat_rx_packets":             bb.RxPackets.Val,
		"stat_rx_bytes":               bb.RxBytes.Val,
		"stat_rx_errors":              bb.RxErrors.Val,
		"stat_rx_dropped":             bb.RxDropped.Val,
		"stat_rx_crypts":              bb.RxCrypts.Val,
		"stat_rx_frags":               bb.RxFrags.Val,
		"stat_tx_packets":             bb.TxPackets.Val,
		"stat_tx_bytes":               bb.TxBytes.Val,
		"stat_tx_errors":              bb.TxErrors.Val,
		"stat_tx_dropped":             bb.TxDropped.Val,
		"stat_tx_retries":             bb.TxRetries.Val,
		"stat_mac_filter_rejections":  bb.MacFilterRejections.Val,
		"stat_wifi_tx_attempts":       bb.WifiTxAttempts.Val,
		"stat_wifi_tx_dropped":        bb.WifiTxDropped.Val,
		// User aggregated stats
		"stat_user-rx_packets":        bb.UserRxPackets.Val,
		"stat_user-rx_bytes":          bb.UserRxBytes.Val,
		"stat_user-rx_errors":         bb.UserRxErrors.Val,
		"stat_user-rx_dropped":        bb.UserRxDropped.Val,
		"stat_user-rx_crypts":         bb.UserRxCrypts.Val,
		"stat_user-rx_frags":          bb.UserRxFrags.Val,
		"stat_user-tx_packets":        bb.UserTxPackets.Val,
		"stat_user-tx_bytes":          bb.UserTxBytes.Val,
		"stat_user-tx_errors":         bb.UserTxErrors.Val,
		"stat_user-tx_dropped":        bb.UserTxDropped.Val,
		"stat_user-tx_retries":        bb.UserTxRetries.Val,
		"stat_user-mac_filter_rejections": bb.UserMacFilterRejections.Val,
		"stat_user-wifi_tx_attempts":      bb.UserWifiTxAttempts.Val,
		"stat_user-wifi_tx_dropped":       bb.UserWifiTxDropped.Val,
		// wifi0 radio stats (5GHz)
		"stat_wifi0-rx_packets":             bb.Wifi0RxPackets.Val,
		"stat_wifi0-rx_bytes":               bb.Wifi0RxBytes.Val,
		"stat_wifi0-rx_errors":              bb.Wifi0RxErrors.Val,
		"stat_wifi0-rx_dropped":             bb.Wifi0RxDropped.Val,
		"stat_wifi0-rx_crypts":              bb.Wifi0RxCrypts.Val,
		"stat_wifi0-rx_frags":               bb.Wifi0RxFrags.Val,
		"stat_wifi0-tx_packets":             bb.Wifi0TxPackets.Val,
		"stat_wifi0-tx_bytes":               bb.Wifi0TxBytes.Val,
		"stat_wifi0-tx_errors":              bb.Wifi0TxErrors.Val,
		"stat_wifi0-tx_dropped":             bb.Wifi0TxDropped.Val,
		"stat_wifi0-tx_retries":             bb.Wifi0TxRetries.Val,
		"stat_wifi0-mac_filter_rejections":  bb.Wifi0MacFilterRejections.Val,
		"stat_wifi0-wifi_tx_attempts":       bb.Wifi0WifiTxAttempts.Val,
		"stat_wifi0-wifi_tx_dropped":        bb.Wifi0WifiTxDropped.Val,
		// terra2 radio stats (60GHz - 802.11ad)
		"stat_terra2-rx_packets":            bb.Terra2RxPackets.Val,
		"stat_terra2-rx_bytes":              bb.Terra2RxBytes.Val,
		"stat_terra2-rx_errors":             bb.Terra2RxErrors.Val,
		"stat_terra2-rx_dropped":            bb.Terra2RxDropped.Val,
		"stat_terra2-rx_crypts":             bb.Terra2RxCrypts.Val,
		"stat_terra2-rx_frags":              bb.Terra2RxFrags.Val,
		"stat_terra2-tx_packets":            bb.Terra2TxPackets.Val,
		"stat_terra2-tx_bytes":              bb.Terra2TxBytes.Val,
		"stat_terra2-tx_errors":             bb.Terra2TxErrors.Val,
		"stat_terra2-tx_dropped":            bb.Terra2TxDropped.Val,
		"stat_terra2-tx_retries":            bb.Terra2TxRetries.Val,
		"stat_terra2-mac_filter_rejections": bb.Terra2MacFilterRejections.Val,
		"stat_terra2-wifi_tx_attempts":      bb.Terra2WifiTxAttempts.Val,
		"stat_terra2-wifi_tx_dropped":       bb.Terra2WifiTxDropped.Val,
		// User wifi0 stats
		"stat_user-wifi0-rx_packets":             bb.UserWifi0RxPackets.Val,
		"stat_user-wifi0-rx_bytes":               bb.UserWifi0RxBytes.Val,
		"stat_user-wifi0-rx_errors":              bb.UserWifi0RxErrors.Val,
		"stat_user-wifi0-rx_dropped":             bb.UserWifi0RxDropped.Val,
		"stat_user-wifi0-rx_crypts":              bb.UserWifi0RxCrypts.Val,
		"stat_user-wifi0-rx_frags":               bb.UserWifi0RxFrags.Val,
		"stat_user-wifi0-tx_packets":             bb.UserWifi0TxPackets.Val,
		"stat_user-wifi0-tx_bytes":               bb.UserWifi0TxBytes.Val,
		"stat_user-wifi0-tx_errors":              bb.UserWifi0TxErrors.Val,
		"stat_user-wifi0-tx_dropped":             bb.UserWifi0TxDropped.Val,
		"stat_user-wifi0-tx_retries":             bb.UserWifi0TxRetries.Val,
		"stat_user-wifi0-mac_filter_rejections":  bb.UserWifi0MacFilterRejections.Val,
		"stat_user-wifi0-wifi_tx_attempts":       bb.UserWifi0WifiTxAttempts.Val,
		"stat_user-wifi0-wifi_tx_dropped":        bb.UserWifi0WifiTxDropped.Val,
		// User terra2 stats (60GHz)
		"stat_user-terra2-rx_packets":            bb.UserTerra2RxPackets.Val,
		"stat_user-terra2-rx_bytes":              bb.UserTerra2RxBytes.Val,
		"stat_user-terra2-rx_errors":             bb.UserTerra2RxErrors.Val,
		"stat_user-terra2-rx_dropped":            bb.UserTerra2RxDropped.Val,
		"stat_user-terra2-rx_crypts":             bb.UserTerra2RxCrypts.Val,
		"stat_user-terra2-rx_frags":              bb.UserTerra2RxFrags.Val,
		"stat_user-terra2-tx_packets":            bb.UserTerra2TxPackets.Val,
		"stat_user-terra2-tx_bytes":              bb.UserTerra2TxBytes.Val,
		"stat_user-terra2-tx_errors":             bb.UserTerra2TxErrors.Val,
		"stat_user-terra2-tx_dropped":            bb.UserTerra2TxDropped.Val,
		"stat_user-terra2-tx_retries":            bb.UserTerra2TxRetries.Val,
		"stat_user-terra2-mac_filter_rejections": bb.UserTerra2MacFilterRejections.Val,
		"stat_user-terra2-wifi_tx_attempts":      bb.UserTerra2WifiTxAttempts.Val,
		"stat_user-terra2-wifi_tx_dropped":       bb.UserTerra2WifiTxDropped.Val,
		// Interface-specific stats
		"stat_user-wifi0-ath0-rx_packets": bb.UserWifi0Ath0RxPackets.Val,
		"stat_user-wifi0-ath0-rx_bytes":   bb.UserWifi0Ath0RxBytes.Val,
		"stat_user-wifi0-ath0-tx_packets": bb.UserWifi0Ath0TxPackets.Val,
		"stat_user-wifi0-ath0-tx_bytes":   bb.UserWifi0Ath0TxBytes.Val,
		"stat_user-terra2-wlan0-rx_packets": bb.UserTerra2Wlan0RxPackets.Val,
		"stat_user-terra2-wlan0-rx_bytes":   bb.UserTerra2Wlan0RxBytes.Val,
		"stat_user-terra2-wlan0-tx_packets": bb.UserTerra2Wlan0TxPackets.Val,
		"stat_user-terra2-wlan0-tx_bytes":   bb.UserTerra2Wlan0TxBytes.Val,
		"stat_user-terra2-wlan0-tx_dropped": bb.UserTerra2Wlan0TxDropped.Val,
		"stat_user-terra2-wlan0-rx_errors":  bb.UserTerra2Wlan0RxErrors.Val,
		"stat_user-terra2-wlan0-tx_errors":  bb.UserTerra2Wlan0TxErrors.Val,
	}
}
