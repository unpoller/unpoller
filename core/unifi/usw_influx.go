package unifi

import (
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
)

// Points generates Unifi Switch datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *USW) Points() ([]*influx.Point, error) {
	return u.PointsAt(time.Now())
}

// PointsAt generates Unifi Switch datapoints for InfluxDB.
// These points can be passed directly to influx.
// This is just like Points(), but specify when points were created.
func (u *USW) PointsAt(now time.Time) ([]*influx.Point, error) {
	if u.Stat.sw == nil {
		// Disabled devices lack stats.
		u.Stat.sw = &sw{}
	}
	tags := map[string]string{
		"id":                     u.ID,
		"mac":                    u.Mac,
		"device_oid":             u.Stat.Oid,
		"site_id":                u.SiteID,
		"site_name":              u.SiteName,
		"name":                   u.Name,
		"adopted":                u.Adopted.Txt,
		"cfgversion":             u.Cfgversion,
		"config_network_ip":      u.ConfigNetwork.IP,
		"config_network_type":    u.ConfigNetwork.Type,
		"device_id":              u.DeviceID,
		"inform_ip":              u.InformIP,
		"known_cfgversion":       u.KnownCfgversion,
		"locating":               u.Locating.Txt,
		"model":                  u.Model,
		"serial":                 u.Serial,
		"type":                   u.Type,
		"dot1x_portctrl_enabled": u.Dot1XPortctrlEnabled.Txt,
		"flowctrl_enabled":       u.FlowctrlEnabled.Txt,
		"has_fan":                u.HasFan.Txt,
		"has_temperature":        u.HasTemperature.Txt,
		"jumboframe_enabled":     u.JumboframeEnabled.Txt,
		"stp_priority":           u.StpPriority,
		"stp_version":            u.StpVersion,
	}
	fields := map[string]interface{}{
		"fw_caps":             u.FwCaps.Val,
		"guest-num_sta":       u.GuestNumSta.Val,
		"ip":                  u.IP,
		"bytes":               u.Bytes.Val,
		"fan_level":           u.FanLevel.Val,
		"general_temperature": u.GeneralTemperature.Val,
		"last_seen":           u.LastSeen.Val,
		"license_state":       u.LicenseState,
		"overheating":         u.Overheating.Val,
		"rx_bytes":            u.RxBytes.Val,
		"tx_bytes":            u.TxBytes.Val,
		"uptime":              u.Uptime.Val,
		"state":               u.State.Val,
		"user-num_sta":        u.UserNumSta.Val,
		"version":             u.Version,
		"loadavg_1":           u.SysStats.Loadavg1.Val,
		"loadavg_5":           u.SysStats.Loadavg5.Val,
		"loadavg_15":          u.SysStats.Loadavg15.Val,
		"mem_buffer":          u.SysStats.MemBuffer.Val,
		"mem_used":            u.SysStats.MemUsed.Val,
		"mem_total":           u.SysStats.MemTotal.Val,
		"cpu":                 u.SystemStats.CPU.Val,
		"mem":                 u.SystemStats.Mem.Val,
		"system_uptime":       u.SystemStats.Uptime.Val,
		"stat_bytes":          u.Stat.Bytes.Val,
		"stat_rx_bytes":       u.Stat.RxBytes.Val,
		"stat_rx_crypts":      u.Stat.RxCrypts.Val,
		"stat_rx_dropped":     u.Stat.RxDropped.Val,
		"stat_rx_errors":      u.Stat.RxErrors.Val,
		"stat_rx_frags":       u.Stat.RxFrags.Val,
		"stat_rx_packets":     u.Stat.TxPackets.Val,
		"stat_tx_bytes":       u.Stat.TxBytes.Val,
		"stat_tx_dropped":     u.Stat.TxDropped.Val,
		"stat_tx_errors":      u.Stat.TxErrors.Val,
		"stat_tx_packets":     u.Stat.TxPackets.Val,
		"stat_tx_retries":     u.Stat.TxRetries.Val,
		"uplink_depth":        u.UplinkDepth.Txt,
	}
	pt, err := influx.NewPoint("usw", tags, fields, now)
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
		pt, err = influx.NewPoint("usw_ports", tags, fields, now)
		if err != nil {
			return points, err
		}
		points = append(points, pt)
	}
	return points, nil
}
