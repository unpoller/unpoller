package unifi

import (
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
)

// Points generates Unifi Switch datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u USW) Points() ([]*influx.Point, error) {
	tags := map[string]string{
		"id":                     u.ID,
		"mac":                    u.Mac,
		"device_type":            u.Stat.O,
		"device_oid":             u.Stat.Oid,
		"site_id":                u.SiteID,
		"site_name":              u.SiteName,
		"name":                   u.Name,
		"adopted":                u.Adopted.Txt,
		"adopt_ip":               u.AdoptIP,
		"adopt_url":              u.AdoptURL,
		"cfgversion":             u.Cfgversion,
		"config_network_ip":      u.ConfigNetwork.IP,
		"config_network_type":    u.ConfigNetwork.Type,
		"connect_request_ip":     u.ConnectRequestIP,
		"connect_request_port":   u.ConnectRequestPort,
		"default":                u.Default.Txt,
		"device_id":              u.DeviceID,
		"discovered_via":         u.DiscoveredVia,
		"inform_ip":              u.InformIP,
		"last_uplink_mac":        u.LastUplink.UplinkMac,
		"known_cfgversion":       u.KnownCfgversion,
		"led_override":           u.LedOverride,
		"locating":               u.Locating.Txt,
		"model":                  u.Model,
		"outdoor_mode_override":  u.OutdoorModeOverride,
		"serial":                 u.Serial,
		"type":                   u.Type,
		"version_incompatible":   u.VersionIncompatible.Txt,
		"dot1x_portctrl_enabled": u.Dot1XPortctrlEnabled.Txt,
		"flowctrl_enabled":       u.FlowctrlEnabled.Txt,
		"has_fan":                u.HasFan.Txt,
		"has_temperature":        u.HasTemperature.Txt,
		"jumboframe_enabled":     u.JumboframeEnabled.Txt,
		"stp_priority":           u.StpPriority,
		"stp_version":            u.StpVersion,
	}
	fields := map[string]interface{}{
		"fw_caps":               u.FwCaps,
		"guest-num_sta":         u.GuestNumSta,
		"ip":                    u.IP,
		"bytes":                 u.Bytes,
		"fan_level":             u.FanLevel,
		"general_temperature":   u.GeneralTemperature,
		"last_seen":             u.LastSeen,
		"license_state":         u.LicenseState,
		"overheating":           u.Overheating.Val,
		"rx_bytes":              u.RxBytes,
		"tx_bytes":              u.TxBytes,
		"uptime":                u.Uptime,
		"considered_lost_at":    u.ConsideredLostAt,
		"next_heartbeat_at":     u.NextHeartbeatAt,
		"roll_upgrade":          u.Rollupgrade.Val,
		"state":                 u.State,
		"upgradable":            u.Upgradable.Val,
		"user-num_sta":          u.UserNumSta,
		"version":               u.Version,
		"loadavg_1":             u.SysStats.Loadavg1,
		"loadavg_5":             u.SysStats.Loadavg5,
		"loadavg_15":            u.SysStats.Loadavg15,
		"mem_buffer":            u.SysStats.MemBuffer,
		"mem_used":              u.SysStats.MemUsed,
		"mem_total":             u.SysStats.MemTotal,
		"cpu":                   u.SystemStats.CPU,
		"mem":                   u.SystemStats.Mem,
		"system_uptime":         u.SystemStats.Uptime,
		"stat_bytes":            u.Stat.Bytes,
		"stat_duration":         u.Stat.Duration,
		"stat_guest-rx_bytes":   u.Stat.RxBytes,
		"stat_guest-rx_crypts":  u.Stat.RxCrypts,
		"stat_guest-rx_dropped": u.Stat.RxDropped,
		"stat_guest-rx_errors":  u.Stat.RxErrors,
		"stat_guest-rx_frags":   u.Stat.RxFrags,
		"stat_guest-rx_packets": u.Stat.RxPackets,
		"stat_guest-tx_bytes":   u.Stat.TxBytes,
		"stat_guest-tx_dropped": u.Stat.TxDropped,
		"stat_guest-tx_errors":  u.Stat.TxErrors,
		"stat_guest-tx_packets": u.Stat.TxPackets,
		"stat_guest-tx_retries": u.Stat.TxRetries,
		"stat_rx_bytes":         u.Stat.RxBytes,
		"stat_rx_crypts":        u.Stat.RxCrypts,
		"stat_rx_dropped":       u.Stat.RxDropped,
		"stat_rx_errors":        u.Stat.RxErrors,
		"stat_rx_frags":         u.Stat.RxFrags,
		"stat_rx_packets":       u.Stat.TxPackets,
		"stat_tx_bytes":         u.Stat.TxBytes,
		"stat_tx_dropped":       u.Stat.TxDropped,
		"stat_tx_errors":        u.Stat.TxErrors,
		"stat_tx_packets":       u.Stat.TxPackets,
		"stat_tx_retries":       u.Stat.TxRetries,
		"uplink_depth":          u.UplinkDepth.Txt,
		// Add the port stats too.
	}
	pt, err := influx.NewPoint("usw", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}
	points := []*influx.Point{pt}
	for _, p := range u.PortTable {
		tags := map[string]string{
			"site_id":       u.SiteID,
			"site_name":     u.SiteName,
			"device_name":   u.Name,
			"name":          p.Name,
			"enable":        p.Enable.Txt,
			"is_uplink":     p.IsUplink.Txt,
			"up":            p.Up.Txt,
			"portconf_id":   p.PortconfID,
			"dot1x_mode":    p.Dot1XMode,
			"dot1x_status":  p.Dot1XStatus,
			"stp_state":     p.StpState,
			"sfp_found":     p.SfpFound.Txt,
			"op_mode":       p.OpMode,
			"poe_mode":      p.PoeMode,
			"port_poe":      p.PortPoe.Txt,
			"port_idx":      p.PortIdx.Txt,
			"port_id":       u.Name + " Port " + p.PortIdx.Txt,
			"poe_enable":    p.PoeEnable.Txt,
			"flowctrl_rx":   p.FlowctrlRx.Txt,
			"flowctrl_tx":   p.FlowctrlTx.Txt,
			"autoneg":       p.Autoneg.Txt,
			"full_duplex":   p.FullDuplex.Txt,
			"jumbo":         p.Jumbo.Txt,
			"masked":        p.Masked.Txt,
			"poe_good":      p.PoeGood.Txt,
			"media":         p.Media,
			"poe_class":     p.PoeClass,
			"poe_caps":      p.PoeCaps.Txt,
			"aggregated_by": p.AggregatedBy.Txt,
		}
		fields := map[string]interface{}{
			"dbytes_r":     p.BytesR.Val,
			"rx_broadcast": p.RxBroadcast.Val,
			"rx_bytes":     p.RxBytes.Val,
			"rx_bytes-r":   p.RxBytesR.Val,
			"rx_dropped":   p.RxDropped.Val,
			"rx_errors":    p.RxErrors.Val,
			"rx_multicast": p.RxMulticast.Val,
			"rx_packets":   p.RxPackets.Val,
			"speed":        p.Speed.Val,
			"stp_pathcost": p.StpPathcost.Val,
			"tx_broadcast": p.TxBroadcast.Val,
			"tx_bytes":     p.TxBytes.Val,
			"tx_bytes-r":   p.TxBytesR.Val,
			"tx_dropped":   p.TxDropped.Val,
			"tx_errors":    p.TxErrors.Val,
			"tx_multicast": p.TxMulticast.Val,
			"tx_packets":   p.TxPackets.Val,
			"poe_current":  p.PoeCurrent.Val,
			"poe_power":    p.PoePower.Val,
			"poe_voltage":  p.PoeVoltage.Val,
			"full_duplex":  p.FullDuplex.Val,
		}
		pt, err = influx.NewPoint("usw_ports", tags, fields, time.Now())
		if err != nil {
			return points, err
		}
		points = append(points, pt)
	}
	return points, nil
}
