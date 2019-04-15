package unifi

import (
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
)

// Points generates Wireless-Access-Point datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u UAP) Points() ([]*influx.Point, error) {
	/* I generally suck at InfluxDB, so if I got the tags/fields wrong,
	   please send me a PR or open an Issue to address my faults. Thanks!
	*/
	var points []*influx.Point
	tags := map[string]string{
		"id":                      u.ID,
		"mac":                     u.Mac,
		"device_type":             u.Stat.O,
		"device_oid":              u.Stat.Oid,
		"device_ap":               u.Stat.Ap,
		"site_id":                 u.SiteID,
		"name":                    u.Name,
		"addopted":                strconv.FormatBool(u.Adopted),
		"bandsteering_mode":       u.BandsteeringMode,
		"board_rev":               strconv.Itoa(u.BoardRev),
		"cfgversion":              u.Cfgversion,
		"config_network_ip":       u.ConfigNetwork.IP,
		"config_network_type":     u.ConfigNetwork.Type,
		"connect_request_ip":      u.ConnectRequestIP,
		"connect_request_port":    u.ConnectRequestPort,
		"default":                 strconv.FormatBool(u.Default),
		"device_id":               u.DeviceID,
		"discovered_via":          u.DiscoveredVia,
		"fw_caps":                 strconv.Itoa(u.FwCaps),
		"guest-num_sta":           strconv.Itoa(u.GuestNumSta),
		"guest_token":             u.GuestToken,
		"has_eth1":                strconv.FormatBool(u.HasEth1),
		"has_speaker":             strconv.FormatBool(u.HasSpeaker),
		"inform_ip":               u.InformIP,
		"isolated":                strconv.FormatBool(u.Isolated),
		"last_uplink_mac":         u.LastUplink.UplinkMac,
		"last_uplink_remote_port": strconv.Itoa(u.LastUplink.UplinkRemotePort),
		"known_cfgversion":        u.KnownCfgversion,
		"led_override":            u.LedOverride,
		"locating":                strconv.FormatBool(u.Locating),
		"model":                   u.Model,
		"outdoor_mode_override":   u.OutdoorModeOverride,
		"serial":                  u.Serial,
		"type":                    u.Type,
		"version_incompatible":    strconv.FormatBool(u.VersionIncompatible),
		"vwireEnabled":            strconv.FormatBool(u.VwireEnabled),
		"wifi_caps":               strconv.Itoa(u.WifiCaps),
	}
	fields := map[string]interface{}{
		"ip":                         u.IP,
		"bytes":                      u.Bytes,
		"bytes_d":                    u.BytesD,
		"bytes_r":                    u.BytesR,
		"last_seen":                  u.LastSeen,
		"rx_bytes":                   u.RxBytes,
		"rx_bytes-d":                 u.RxBytesD,
		"tx_bytes":                   u.TxBytes,
		"tx_bytes-d":                 u.TxBytesD,
		"uptime":                     u.Uptime.Number,
		"considered_lost_at":         u.ConsideredLostAt,
		"next_heartbeat_at":          u.NextHeartbeatAt,
		"scanning":                   u.Scanning,
		"spectrum_scanning":          u.SpectrumScanning,
		"roll_upgrade":               u.Rollupgrade,
		"state":                      u.State,
		"upgradable":                 u.Upgradable,
		"user-num_sta":               u.UserNumSta,
		"version":                    u.Version,
		"loadavg_1":                  u.SysStats.Loadavg1,
		"loadavg_5":                  u.SysStats.Loadavg5,
		"loadavg_15":                 u.SysStats.Loadavg15,
		"mem_buffer":                 u.SysStats.MemBuffer,
		"mem_total":                  u.SysStats.MemTotal,
		"mem_used":                   u.SysStats.MemUsed,
		"cpu":                        u.SystemStats.CPU,
		"mem":                        u.SystemStats.Mem,
		"system_uptime":              u.SystemStats.Uptime,
		"stat_bytes":                 u.Stat.Bytes,
		"stat_duration":              u.Stat.Duration,
		"stat_guest-rx_bytes":        u.Stat.RxBytes,
		"stat_guest-rx_crypts":       u.Stat.RxCrypts,
		"stat_guest-rx_dropped":      u.Stat.RxDropped,
		"stat_guest-rx_errors":       u.Stat.RxErrors,
		"stat_guest-rx_frags":        u.Stat.RxFrags,
		"stat_guest-rx_packets":      u.Stat.RxPackets,
		"stat_guest-tx_bytes":        u.Stat.TxBytes,
		"stat_guest-tx_dropped":      u.Stat.TxDropped,
		"stat_guest-tx_errors":       u.Stat.TxErrors,
		"stat_guest-tx_packets":      u.Stat.TxPackets,
		"stat_guest-tx_retries":      u.Stat.TxRetries,
		"stat_port_1-rx_broadcast":   u.Stat.Port1RxBroadcast,
		"stat_port_1-rx_bytes":       u.Stat.Port1RxBytes,
		"stat_port_1-rx_multicast":   u.Stat.Port1RxMulticast,
		"stat_port_1-rx_packets":     u.Stat.Port1RxPackets,
		"stat_port_1-tx_broadcast":   u.Stat.Port1TxBroadcast,
		"stat_port_1-tx_bytes":       u.Stat.Port1TxBytes,
		"stat_port_1-tx_multicast":   u.Stat.Port1TxMulticast,
		"stat_port_1-tx_packets":     u.Stat.Port1TxPackets,
		"stat_rx_bytes":              u.Stat.RxBytes,
		"stat_rx_crypts":             u.Stat.RxCrypts,
		"stat_rx_dropped":            u.Stat.RxDropped,
		"stat_rx_errors":             u.Stat.RxErrors,
		"stat_rx_frags":              u.Stat.RxFrags,
		"stat_rx_packets":            u.Stat.TxPackets,
		"stat_tx_bytes":              u.Stat.TxBytes,
		"stat_tx_dropped":            u.Stat.TxDropped,
		"stat_tx_errors":             u.Stat.TxErrors,
		"stat_tx_packets":            u.Stat.TxPackets,
		"stat_tx_retries":            u.Stat.TxRetries,
		"stat_user-rx_bytes":         u.Stat.UserRxBytes,
		"stat_user-rx_crypts":        u.Stat.UserRxCrypts,
		"stat_user-rx_dropped":       u.Stat.UserRxDropped,
		"stat_user-rx_errors":        u.Stat.UserRxErrors,
		"stat_user-rx_frags":         u.Stat.UserRxFrags,
		"stat_user-rx_packets":       u.Stat.UserRxPackets,
		"stat_user-tx_bytes":         u.Stat.UserTxBytes,
		"stat_user-tx_dropped":       u.Stat.UserTxDropped,
		"stat_user-tx_errors":        u.Stat.UserTxErrors,
		"stat_user-tx_packets":       u.Stat.UserTxPackets,
		"stat_user-tx_retries":       u.Stat.UserTxRetries,
		"stat_user-wifi0-rx_bytes":   u.Stat.UserWifi0RxBytes,
		"stat_user-wifi0-rx_crypts":  u.Stat.UserWifi0RxCrypts,
		"stat_user-wifi0-rx_dropped": u.Stat.UserWifi0RxDropped,
		"stat_user-wifi0-rx_errors":  u.Stat.UserWifi0RxErrors,
		"stat_user-wifi0-rx_frags":   u.Stat.UserWifi0RxFrags,
		"stat_user-wifi0-rx_packets": u.Stat.UserWifi0RxPackets,
		"stat_user-wifi0-tx_bytes":   u.Stat.UserWifi0TxBytes,
		"stat_user-wifi0-tx_dropped": u.Stat.UserWifi0TxDropped,
		"stat_user-wifi0-tx_errors":  u.Stat.UserWifi0TxErrors,
		"stat_user-wifi0-tx_packets": u.Stat.UserWifi0TxPackets,
		"stat_user-wifi0-tx_retries": u.Stat.UserWifi0TxRetries,
		"stat_user-wifi1-rx_bytes":   u.Stat.UserWifi1RxBytes,
		"stat_user-wifi1-rx_crypts":  u.Stat.UserWifi1RxCrypts,
		"stat_user-wifi1-rx_dropped": u.Stat.UserWifi1RxDropped,
		"stat_user-wifi1-rx_errors":  u.Stat.UserWifi1RxErrors,
		"stat_user-wifi1-rx_frags":   u.Stat.UserWifi1RxFrags,
		"stat_user-wifi1-rx_packets": u.Stat.UserWifi1RxPackets,
		"stat_user-wifi1-tx_bytes":   u.Stat.UserWifi1TxBytes,
		"stat_user-wifi1-tx_dropped": u.Stat.UserWifi1TxDropped,
		"stat_user-wifi1-tx_errors":  u.Stat.UserWifi1TxErrors,
		"stat_user-wifi1-tx_packets": u.Stat.UserWifi1TxPackets,
		"stat_user-wifi1-tx_retries": u.Stat.UserWifi1TxRetries,
		"stat_wifi0-rx_bytes":        u.Stat.Wifi0RxBytes,
		"stat_wifi0-rx_crypts":       u.Stat.Wifi0RxCrypts,
		"stat_wifi0-rx_dropped":      u.Stat.Wifi0RxDropped,
		"stat_wifi0-rx_errors":       u.Stat.Wifi0RxErrors,
		"stat_wifi0-rx_frags":        u.Stat.Wifi0RxFrags,
		"stat_wifi0-rx_packets":      u.Stat.Wifi0RxPackets,
		"stat_wifi0-tx_bytes":        u.Stat.Wifi0TxBytes,
		"stat_wifi0-tx_dropped":      u.Stat.Wifi0TxDropped,
		"stat_wifi0-tx_errors":       u.Stat.Wifi0TxErrors,
		"stat_wifi0-tx_packets":      u.Stat.Wifi0TxPackets,
		"stat_wifi0-tx_retries":      u.Stat.Wifi0TxRetries,
		"stat_wifi1-rx_bytes":        u.Stat.Wifi1RxBytes,
		"stat_wifi1-rx_crypts":       u.Stat.Wifi1RxCrypts,
		"stat_wifi1-rx_dropped":      u.Stat.Wifi1RxDropped,
		"stat_wifi1-rx_errors":       u.Stat.Wifi1RxErrors,
		"stat_wifi1-rx_frags":        u.Stat.Wifi1RxFrags,
		"stat_wifi1-rx_packets":      u.Stat.Wifi1RxPackets,
		"stat_wifi1-tx_bytes":        u.Stat.Wifi1TxBytes,
		"stat_wifi1-tx_dropped":      u.Stat.Wifi1TxDropped,
		"stat_wifi1-tx_errors":       u.Stat.Wifi1TxErrors,
		"stat_wifi1-tx_packets":      u.Stat.Wifi1TxPackets,
		"stat_wifi1-tx_retries":      u.Stat.Wifi1TxRetries,
	}
	pt, err := influx.NewPoint("uap", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}
	points = append(points, pt)
	for _, p := range u.RadioTable {
		tags := map[string]string{
			"device_name":  u.Name,
			"device_id":    u.ID,
			"device_mac":   u.Mac,
			"name":         p.Name,
			"wlangroup_id": p.WlangroupID,
			"channel":      p.Channel.String,
			"radio":        p.Radio,
		}
		fields := map[string]interface{}{
			"builtin_ant_gain":     p.BuiltinAntGain,
			"current_antenna_gain": p.CurrentAntennaGain,
			"has_dfs":              p.HasDfs,
			"has_fccdfs":           p.HasFccdfs,
			"ht":                   p.Ht,
			"is_11ac":              p.Is11Ac,
			"max_txpower":          p.MaxTxpower,
			"min_rssi_enabled":     p.MinRssiEnabled,
			"min_txpower":          p.MinTxpower,
			"nss":                  p.Nss,
			"radio_caps":           p.RadioCaps,
			"tx_power":             p.TxPower.Number,
			"tx_power_mode":        p.TxPowerMode,
		}

		for _, s := range u.RadioTableStats {
			// This may be a tad slower but it allows putting
			// all the radio stats into one table.
			if p.Name == s.Name {
				fields["ast_be_xmit"] = s.AstBeXmit
				fields["ast_cst"] = s.AstCst
				fields["channel"] = s.Channel
				fields["ast_txto"] = s.AstTxto
				fields["cu_self_rx"] = s.CuSelfRx
				fields["cu_self_tx"] = s.CuSelfTx
				fields["cu_total"] = s.CuTotal
				fields["extchannel"] = s.Extchannel
				fields["gain"] = s.Gain
				fields["guest-num_sta"] = s.GuestNumSta
				fields["num_sta"] = s.NumSta
				fields["radio"] = s.Radio
				fields["state"] = s.State
				fields["radio_tx_packets"] = s.TxPackets
				fields["radio_tx_power"] = s.TxPower
				fields["radio_tx_retries"] = s.TxRetries
				fields["user-num_sta"] = s.UserNumSta
				break
			}
		}
		for _, s := range u.VapTable {
			if p.Name == s.RadioName {
				tags["ap_mac"] = s.ApMac
				tags["bssid"] = s.Bssid
				tags["vap_id"] = s.ID
				tags["vap_name"] = s.Name
				tags["wlanconf_id"] = s.WlanconfID
				fields["ccq"] = s.Ccq
				fields["essid"] = s.Essid
				fields["extchannel"] = s.Extchannel
				fields["is_guest"] = s.IsGuest
				fields["is_wep"] = s.IsWep
				fields["mac_filter_rejections"] = s.MacFilterRejections
				fields["map_id"] = s.MapID
				fields["vap_rx_bytes"] = s.RxBytes
				fields["vap_rx_crypts"] = s.RxCrypts
				fields["vap_rx_dropped"] = s.RxDropped
				fields["vap_rx_errors"] = s.RxErrors
				fields["vap_rx_frags"] = s.RxFrags
				fields["vap_rx_nwids"] = s.RxNwids
				fields["vap_rx_packets"] = s.RxPackets
				fields["vap_tx_bytes"] = s.TxBytes
				fields["vap_tx_dropped"] = s.TxDropped
				fields["vap_tx_errors"] = s.TxErrors
				fields["vap_tx_latency_avg"] = s.TxLatencyAvg
				fields["vap_tx_latency_max"] = s.TxLatencyMax
				fields["vap_tx_latency_min"] = s.TxLatencyMin
				fields["vap_tx_packets"] = s.TxPackets
				fields["vap_tx_power"] = s.TxPower
				fields["vap_tx_retries"] = s.TxRetries
				fields["usage"] = s.Usage
				break
			}
		}
		pt, err := influx.NewPoint("uap_radios", tags, fields, time.Now())
		if err != nil {
			return points, err
		}
		points = append(points, pt)
	}
	return points, nil
}
