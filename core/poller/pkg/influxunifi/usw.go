package influxunifi

import (
	"golift.io/unifi"
)

// batchUSW generates Unifi Switch datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchUSW(r report, s *unifi.USW) {
	if s.Stat.Sw == nil {
		s.Stat.Sw = &unifi.Sw{}
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
	fields := Combine(map[string]interface{}{
		"guest-num_sta":       s.GuestNumSta.Val,
		"ip":                  s.IP,
		"bytes":               s.Bytes.Val,
		"fan_level":           s.FanLevel.Val,
		"general_temperature": s.GeneralTemperature.Val,
		"last_seen":           s.LastSeen.Val,
		"rx_bytes":            s.RxBytes.Val,
		"tx_bytes":            s.TxBytes.Val,
		"uptime":              s.Uptime.Val,
		"state":               s.State.Val,
		"user-num_sta":        s.UserNumSta.Val,
		"stat_bytes":          s.Stat.Sw.Bytes.Val,
		"stat_rx_bytes":       s.Stat.Sw.RxBytes.Val,
		"stat_rx_crypts":      s.Stat.Sw.RxCrypts.Val,
		"stat_rx_dropped":     s.Stat.Sw.RxDropped.Val,
		"stat_rx_errors":      s.Stat.Sw.RxErrors.Val,
		"stat_rx_frags":       s.Stat.Sw.RxFrags.Val,
		"stat_rx_packets":     s.Stat.Sw.TxPackets.Val,
		"stat_tx_bytes":       s.Stat.Sw.TxBytes.Val,
		"stat_tx_dropped":     s.Stat.Sw.TxDropped.Val,
		"stat_tx_errors":      s.Stat.Sw.TxErrors.Val,
		"stat_tx_packets":     s.Stat.Sw.TxPackets.Val,
		"stat_tx_retries":     s.Stat.Sw.TxRetries.Val,
	}, u.batchSysStats(r, s.SysStats, s.SystemStats))
	r.send(&metric{Table: "usw", Tags: tags, Fields: fields})
	u.batchPortTable(r, tags, s.PortTable)
}

func (u *InfluxUnifi) batchPortTable(r report, tags map[string]string, pt []unifi.Port) {
	for _, p := range pt {
		if !p.Up.Val || !p.Enable.Val {
			continue // only record UP ports.
		}
		tags := map[string]string{
			"site_name":   tags["site_name"],
			"device_name": tags["name"],
			"name":        p.Name,
			"poe_mode":    p.PoeMode,
			"port_poe":    p.PortPoe.Txt,
			"port_idx":    p.PortIdx.Txt,
			"port_id":     tags["name"] + " Port " + p.PortIdx.Txt,
			"poe_enable":  p.PoeEnable.Txt,
			"flowctrl_rx": p.FlowctrlRx.Txt,
			"flowctrl_tx": p.FlowctrlTx.Txt,
			"media":       p.Media,
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
		}
		if p.PoeEnable.Val && p.PortPoe.Val {
			fields["poe_current"] = p.PoeCurrent.Val
			fields["poe_power"] = p.PoePower.Val
			fields["poe_voltage"] = p.PoeVoltage.Val
		}
		r.send(&metric{Table: "usw_ports", Tags: tags, Fields: fields})
	}
}
