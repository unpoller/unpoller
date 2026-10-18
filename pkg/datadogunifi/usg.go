package datadogunifi

import (
	"maps"

	"github.com/unpoller/unifi/v6"
)

// usgT is used as a name for printed/logged counters.
const usgT = item("USG")

// batchUSG generates Unifi Gateway datapoints for Datadog.
// These points can be passed directly to datadog.
func (u *DatadogUnifi) batchUSG(r report, s *unifi.USG) {
	if !s.Adopted.Val {
		return
	}

	tags := map[string]string{
		"mac":           s.Mac,
		"site_name":     s.SiteName,
		"source":        s.SourceName,
		"name":          s.Name,
		"version":       s.Version,
		"model":         s.Model,
		"serial":        s.Serial,
		"type":          s.Type,
		"ip":            s.IP,
		"license_state": s.LicenseState,
	}
	data := CombineFloat64(
		u.batchUDMtemps(s.Temperatures),
		u.batchSysStats(s.SysStats, s.SystemStats),
		u.batchUSGstats(s.SpeedtestStatus, s.Stat.Gw, s.Uplink),
		map[string]float64{
			"bytes":         s.Bytes.Val,
			"last_seen":     s.LastSeen.Val,
			"guest_num_sta": s.GuestNumSta.Val,
			"rx_bytes":      s.RxBytes.Val,
			"tx_bytes":      s.TxBytes.Val,
			"uptime":        s.Uptime.Val,
			"state":         s.State.Val,
			"user_num_sta":  s.UserNumSta.Val,
			"num_desktop":   s.NumDesktop.Val,
			"num_handheld":  s.NumHandheld.Val,
			"num_mobile":    s.NumMobile.Val,
			"upgradeable":   boolToFloat64(s.Upgradable.Val),
		},
	)

	r.addCount(usgT)

	metricName := metricNamespace("usg")

	reportGaugeForFloat64Map(r, metricName, data, tags)

	u.batchNetTable(r, tags, s.NetworkTable)
	u.batchUSGwans(r, tags, s.Wan1, s.Wan2)
}

func (u *DatadogUnifi) batchUSGstats(ss unifi.SpeedtestStatus, gw *unifi.Gw, ul unifi.Uplink) map[string]float64 {
	// Uplink and speedtest live on the device, not in stat.gw. Network 11 omits
	// stat entirely, so a missing gateway stat must not drop those fields.
	stats := map[string]float64{
		"uplink_latency":                 ul.Latency.Val,
		"uplink_speed":                   ul.Speed.Val,
		"uplink_max_speed":               ul.MaxSpeed.Val,
		"uplink_uptime":                  ul.Uptime.Val,
		"speedtest_status_latency":       ss.Latency.Val,
		"speedtest_status_runtime":       ss.Runtime.Val,
		"speedtest_status_rundate":       ss.Rundate.Val,
		"speedtest_status_ping":          ss.StatusPing.Val,
		"speedtest_status_xput_download": ss.XputDownload.Val,
		"speedtest_status_xput_upload":   ss.XputUpload.Val,
	}

	if gw == nil {
		return stats
	}

	return CombineFloat64(stats, map[string]float64{
		"lan_rx_bytes":   gw.LanRxBytes.Val,
		"lan_rx_packets": gw.LanRxPackets.Val,
		"lan_tx_bytes":   gw.LanTxBytes.Val,
		"lan_tx_packets": gw.LanTxPackets.Val,
		"lan_rx_dropped": gw.LanRxDropped.Val,
	})
}

func (u *DatadogUnifi) batchUSGwans(r report, tags map[string]string, wans ...unifi.Wan) {
	for _, wan := range wans {
		if !wan.Up.Val {
			continue
		}

		tags := cleanTags(map[string]string{
			"device_name": tags["name"],
			"site_name":   tags["site_name"],
			"source":      tags["source"],
			"ip":          wan.IP,
			"purpose":     wan.Name,
			"mac":         wan.Mac,
			"ifname":      wan.Ifname,
			"type":        wan.Type,
			"up":          wan.Up.Txt,
			"enabled":     wan.Enable.Txt,
			"gateway":     wan.Gateway,
		})

		fullDuplex := 0.0
		if wan.FullDuplex.Val {
			fullDuplex = 1.0
		}

		data := map[string]float64{
			"bytes_r":      wan.BytesR.Val,
			"full_duplex":  fullDuplex,
			"max_speed":    wan.MaxSpeed.Val,
			"rx_bytes":     wan.RxBytes.Val,
			"rx_bytes_r":   wan.RxBytesR.Val,
			"rx_dropped":   wan.RxDropped.Val,
			"rx_errors":    wan.RxErrors.Val,
			"rx_broadcast": wan.RxBroadcast.Val,
			"rx_multicast": wan.RxMulticast.Val,
			"rx_packets":   wan.RxPackets.Val,
			"speed":        wan.Speed.Val,
			"tx_bytes":     wan.TxBytes.Val,
			"tx_bytes_r":   wan.TxBytesR.Val,
			"tx_dropped":   wan.TxDropped.Val,
			"tx_errors":    wan.TxErrors.Val,
			"tx_packets":   wan.TxPackets.Val,
			"tx_broadcast": wan.TxBroadcast.Val,
			"tx_multicast": wan.TxMulticast.Val,
		}

		// Older gateways do not measure WAN health; don't report them as 0% available.
		// A WAN that is offline reports latency as 0; don't report that as perfect latency.
		if wan.Latency.Val > 0 {
			data["latency"] = wan.Latency.Val
		}

		if wan.Availability.Txt != "" {
			data["availability"] = wan.Availability.Val
		}

		metricName := metricNamespace("usg.wan_ports")
		reportGaugeForFloat64Map(r, metricName, data, tags)
	}
}

// batchWANHealth reports per-WAN health keyed by WAN network group (WAN, WAN2).
// Unlike batchUSGwans, it does not skip a WAN whose port is down.
// uptime_stats availability is left out: it is the same statistic as the wan uptime_percentage gauge.
// uptime_stats latency_average is left out too: its averaging window varies in length.
func (u *DatadogUnifi) batchWANHealth(
	r report, tags map[string]string, stats map[string]*unifi.UptimeStats, status map[string]string,
) {
	for wan, s := range status {
		tags := wanHealthTags(tags, wan)
		tags["status"] = s

		online := 0.0
		if s == "online" {
			online = 1.0
		}

		_ = r.reportGauge(metricNamespace("usg.wan_status")("online"), online, tagMapToTags(cleanTags(tags)))
	}

	for wan, s := range stats {
		if s == nil {
			continue
		}

		tags := cleanTags(wanHealthTags(tags, wan))
		data := map[string]float64{}

		// While a WAN is offline the controller omits uptime and reports downtime instead.
		if s.Uptime.Txt != "" {
			data["uptime"] = s.Uptime.Val
		}

		if s.Downtime.Txt != "" {
			data["downtime"] = s.Downtime.Val
		}

		reportGaugeForFloat64Map(r, metricNamespace("usg.wan_uptime_stats"), data, tags)
		u.batchUptimeMonitors(r, tags, "false", s.Monitors)
		u.batchUptimeMonitors(r, tags, "true", s.AlertingMonitors)
	}
}

func (u *DatadogUnifi) batchUptimeMonitors(r report, tags map[string]string, alerting string, monitors []unifi.UptimeMonitor) {
	for _, m := range monitors {
		tags := maps.Clone(tags)
		tags["target"] = m.Target
		tags["type"] = m.Type
		tags["alerting"] = alerting

		data := map[string]float64{}

		if m.Availability.Txt != "" {
			data["availability"] = m.Availability.Val
		}

		// A failing monitor reports latency as 0 or omits it.
		if m.LatencyAverage.Val > 0 {
			data["latency_average"] = m.LatencyAverage.Val
		}

		reportGaugeForFloat64Map(r, metricNamespace("usg.wan_monitors"), data, cleanTags(tags))
	}
}

func wanHealthTags(tags map[string]string, wan string) map[string]string {
	return map[string]string{
		"device_name":      tags["name"],
		"site_name":        tags["site_name"],
		"source":           tags["source"],
		"wan_networkgroup": wan,
	}
}

func (u *DatadogUnifi) batchNetTable(r report, tags map[string]string, nt unifi.NetworkTable) {
	for _, p := range nt {
		tags := cleanTags(map[string]string{
			"device_name": tags["name"],
			"site_name":   tags["site_name"],
			"source":      tags["source"],
			"up":          p.Up.Txt,
			"enabled":     p.Enabled.Txt,
			"ip":          p.IP,
			"mac":         p.Mac,
			"name":        p.Name,
			"domain_name": p.DomainName,
			"purpose":     p.Purpose,
			"is_guest":    p.IsGuest.Txt,
		})
		data := map[string]float64{
			"num_sta":    p.NumSta.Val,
			"rx_bytes":   p.RxBytes.Val,
			"rx_packets": p.RxPackets.Val,
			"tx_bytes":   p.TxBytes.Val,
			"tx_packets": p.TxPackets.Val,
		}

		metricName := metricNamespace("usg.networks")
		reportGaugeForFloat64Map(r, metricName, data, tags)
	}
}
