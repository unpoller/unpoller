package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// uswT is used as a name for printed/logged counters.
const uswT = item("USW")

// batchUSW generates Unifi Switch datapoints for Datadog.
// These points can be passed directly to datadog.
func (u *DatadogUnifi) batchUSW(r report, s *unifi.USW) {
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := cleanTags(map[string]string{
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"source":    s.SourceName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
		"ip":        s.IP,
	})
	data := CombineFloat64(
		u.batchUSWstat(s.Stat.Sw),
		u.batchSysStats(s.SysStats, s.SystemStats),
		map[string]float64{
			"guest_num_sta":       s.GuestNumSta.Val,
			"bytes":               s.Bytes.Val,
			"fan_level":           s.FanLevel.Val,
			"general_temperature": s.GeneralTemperature.Val,
			"last_seen":           s.LastSeen.Val,
			"rx_bytes":            s.RxBytes.Val,
			"tx_bytes":            s.TxBytes.Val,
			"uptime":              s.Uptime.Val,
			"state":               s.State.Val,
			"user_num_sta":        s.UserNumSta.Val,
			"upgradeable":         boolToFloat64(s.Upgradeable.Val),
			"uplink_speed":        s.Uplink.Speed.Val,
			"uplink_max_speed":    s.Uplink.MaxSpeed.Val,
			"uplink_latency":      s.Uplink.Latency.Val,
			"uplink_uptime":       s.Uplink.Uptime.Val,
		})

	r.addCount(uswT)

	metricName := metricNamespace("usw")

	reportGaugeForFloat64Map(r, metricName, data, tags)

	u.batchPortTable(r, tags, s.PortTable)
}

func (u *DatadogUnifi) batchUSWstat(sw *unifi.Sw) map[string]float64 {
	if sw == nil {
		return map[string]float64{}
	}

	return map[string]float64{
		"stat_bytes":      sw.Bytes.Val,
		"stat_rx_bytes":   sw.RxBytes.Val,
		"stat_rx_crypts":  sw.RxCrypts.Val,
		"stat_rx_dropped": sw.RxDropped.Val,
		"stat_rx_errors":  sw.RxErrors.Val,
		"stat_rx_frags":   sw.RxFrags.Val,
		"stat_rx_packets": sw.TxPackets.Val,
		"stat_tx_bytes":   sw.TxBytes.Val,
		"stat_tx_dropped": sw.TxDropped.Val,
		"stat_tx_errors":  sw.TxErrors.Val,
		"stat_tx_packets": sw.TxPackets.Val,
		"stat_tx_retries": sw.TxRetries.Val,
	}
}

//nolint:funlen
func (u *DatadogUnifi) batchPortTable(r report, t map[string]string, pt []unifi.Port) {
	for _, p := range pt {
		if !u.DeadPorts && (!p.Up.Val || !p.Enable.Val) {
			continue // only record UP ports.
		}

		tags := cleanTags(map[string]string{
			"site_name":      t["site_name"],
			"device_name":    t["name"],
			"source":         t["source"],
			"type":           t["type"],
			"name":           p.Name,
			"poe_mode":       p.PoeMode,
			"port_poe":       p.PortPoe.Txt,
			"port_idx":       p.PortIdx.Txt,
			"port_id":        t["name"] + " Port " + p.PortIdx.Txt,
			"poe_enable":     p.PoeEnable.Txt,
			"flow_ctrl_rx":   p.FlowctrlRx.Txt,
			"flow_ctrl_tx":   p.FlowctrlTx.Txt,
			"media":          p.Media,
			"has_sfp":        p.SFPFound.Txt,
			"sfp_compliance": p.SFPCompliance,
			"sfp_serial":     p.SFPSerial,
			"sfp_vendor":     p.SFPVendor,
			"sfp_part":       p.SFPPart,
		})
		data := map[string]float64{
			"bytes_r":       p.BytesR.Val,
			"rx_broadcast":  p.RxBroadcast.Val,
			"rx_bytes":      p.RxBytes.Val,
			"rx_bytes_r":    p.RxBytesR.Val,
			"rx_dropped":    p.RxDropped.Val,
			"rx_errors":     p.RxErrors.Val,
			"rx_multicast":  p.RxMulticast.Val,
			"rx_packets":    p.RxPackets.Val,
			"speed":         p.Speed.Val,
			"stp_path_cost": p.StpPathcost.Val,
			"tx_broadcast":  p.TxBroadcast.Val,
			"tx_bytes":      p.TxBytes.Val,
			"tx_bytes_r":    p.TxBytesR.Val,
			"tx_dropped":    p.TxDropped.Val,
			"tx_errors":     p.TxErrors.Val,
			"tx_multicast":  p.TxMulticast.Val,
			"tx_packets":    p.TxPackets.Val,
		}

		if p.PoeEnable.Val && p.PortPoe.Val {
			data["poe_current"] = p.PoeCurrent.Val
			data["poe_power"] = p.PoePower.Val
			data["poe_voltage"] = p.PoeVoltage.Val
		}

		if p.SFPFound.Val {
			data["sfp_current"] = p.SFPCurrent.Val
			data["sfp_voltage"] = p.SFPVoltage.Val
			data["sfp_temperature"] = p.SFPTemperature.Val
			data["sfp_tx_power"] = p.SFPTxpower.Val
			data["sfp_rx_power"] = p.SFPRxpower.Val
		}

		metricName := metricNamespace("usw.ports")
		reportGaugeForFloat64Map(r, metricName, data, tags)
	}
}
