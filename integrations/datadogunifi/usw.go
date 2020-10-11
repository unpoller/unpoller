package datadogunifi

import (
	"fmt"

	"github.com/unifi-poller/unifi"
)

// reportUSW generates Unifi Switch datapoints for Datadog.
// These points can be passed directly to datadog.
func (u *DatadogUnifi) reportUSW(r report, s *unifi.USW) {
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := []string{
		tag("mac", s.Mac),
		tag("site_name", s.SiteName),
		tag("source", s.SourceName),
		tag("name", s.Name),
		tag("version", s.Version),
		tag("model", s.Model),
		tag("serial", s.Serial),
		tag("type", s.Type),
		tag("ip", s.IP),
	}
	metricName := metricNamespace("usw")
	u.reportUSWstat(r, metricName, tags, s.Stat.Sw)
	u.reportSysStats(r, metricName, s.SysStats, s.SystemStats, tags)

	data := map[string]float64{
		"guest-num_sta":       s.GuestNumSta.Val,
		"bytes":               s.Bytes.Val,
		"fan_level":           s.FanLevel.Val,
		"general_temperature": s.GeneralTemperature.Val,
		"last_seen":           s.LastSeen.Val,
		"rx_bytes":            s.RxBytes.Val,
		"tx_bytes":            s.TxBytes.Val,
		"uptime":              s.Uptime.Val,
		"state":               s.State.Val,
		"user-num_sta":        s.UserNumSta.Val,
	}
	reportGaugeForMap(r, metricName, data, tags)
	u.reportPortTable(r, s.Name, s.SiteName, s.SourceName, s.Type, s.PortTable)
}

func (u *DatadogUnifi) reportUSWstat(r report, metricName func(string) string, tags []string, sw *unifi.Sw) {
	if sw == nil {
		return
	}

	data := map[string]float64{
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
	reportGaugeForMap(r, metricName, data, tags)
}

func (u *DatadogUnifi) reportPortTable(r report, deviceName string, siteName string, source string, typeTag string, pt []unifi.Port) {
	for _, p := range pt {
		if !p.Up.Val || !p.Enable.Val {
			continue // only record UP ports.
		}

		tags := []string{
			tag("site_name", siteName),
			tag("device_name", deviceName),
			tag("source", source),
			tag("type", typeTag),
			tag("name", p.Name),
			tag("poe_mode", p.PoeMode),
			tag("port_poe", p.PortPoe.Txt),
			tag("port_idx", p.PortIdx.Txt),
			tag("port_id", fmt.Sprintf("%s_port_%s", deviceName, p.PortIdx.Txt)),
			tag("poe_enable", p.PoeEnable.Txt),
			tag("flowctrl_rx", p.FlowctrlRx.Txt),
			tag("flowctrl_tx", p.FlowctrlTx.Txt),
			tag("media", p.Media),
		}
		data := map[string]float64{
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
			data["poe_current"] = p.PoeCurrent.Val
			data["poe_power"] = p.PoePower.Val
			data["poe_voltage"] = p.PoeVoltage.Val
		}

		metricName := metricNamespace("usw_ports")
		reportGaugeForMap(r, metricName, data, tags)
	}
}
