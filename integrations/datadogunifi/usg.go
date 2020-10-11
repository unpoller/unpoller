package datadogunifi

import (
	"github.com/unifi-poller/unifi"
)

// reportUSG generates Unifi Gateway datapoints for Datadog.
// These points can be passed directly to datadog.
func (u *DatadogUnifi) reportUSG(r report, s *unifi.USG) {
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
		tag("license_state", s.LicenseState),
	}
	metricName := metricNamespace("usg")
	u.reportSysStats(r, metricName, s.SysStats, s.SystemStats, tags)
	u.reportUSGstats(r, metricName, tags, s.SpeedtestStatus, s.Stat.Gw, s.Uplink)

	data := map[string]float64{
		"bytes":         s.Bytes.Val,
		"last_seen":     s.LastSeen.Val,
		"guest-num_sta": s.GuestNumSta.Val,
		"rx_bytes":      s.RxBytes.Val,
		"tx_bytes":      s.TxBytes.Val,
		"uptime":        s.Uptime.Val,
		"state":         s.State.Val,
		"user-num_sta":  s.UserNumSta.Val,
		"num_desktop":   s.NumDesktop.Val,
		"num_handheld":  s.NumHandheld.Val,
		"num_mobile":    s.NumMobile.Val,
	}
	reportGaugeForMap(r, metricName, data, tags)

	u.reportNetTable(r, s.Name, s.SiteName, s.SourceName, s.NetworkTable)
	u.reportUSGwans(r, s.Name, s.SiteName, s.SourceName, s.Wan1, s.Wan2)
}

func (u *DatadogUnifi) reportUSGstats(r report, metricName func(string) string, tags []string, ss unifi.SpeedtestStatus, gw *unifi.Gw, ul unifi.Uplink) {
	if gw == nil {
		return
	}
	data := map[string]float64{
		"uplink_latency":                 ul.Latency.Val,
		"uplink_speed":                   ul.Speed.Val,
		"speedtest-status_latency":       ss.Latency.Val,
		"speedtest-status_runtime":       ss.Runtime.Val,
		"speedtest-status_rundate":       ss.Rundate.Val,
		"speedtest-status_ping":          ss.StatusPing.Val,
		"speedtest-status_xput_download": ss.XputDownload.Val,
		"speedtest-status_xput_upload":   ss.XputUpload.Val,
		"lan-rx_bytes":                   gw.LanRxBytes.Val,
		"lan-rx_packets":                 gw.LanRxPackets.Val,
		"lan-tx_bytes":                   gw.LanTxBytes.Val,
		"lan-tx_packets":                 gw.LanTxPackets.Val,
		"lan-rx_dropped":                 gw.LanRxDropped.Val,
	}
	reportGaugeForMap(r, metricName, data, tags)
}

func (u *DatadogUnifi) reportUSGwans(r report, deviceName string, siteName string, source string, wans ...unifi.Wan) {
	for _, wan := range wans {
		if !wan.Up.Val {
			continue
		}

		tags := []string{
			tag("device_name", deviceName),
			tag("site_name", siteName),
			tag("source", source),
			tag("ip", wan.IP),
			tag("purpose", wan.Name),
			tag("mac", wan.Mac),
			tag("ifname", wan.Ifname),
			tag("type", wan.Type),
			tag("up", wan.Up.Txt),
			tag("enabled", wan.Enable.Txt),
			tag("gateway", wan.Gateway),
		}
		fullDuplex := float64(0)
		if wan.FullDuplex.Val {
			fullDuplex = 1
		}

		data := map[string]float64{
			"bytes-r":      wan.BytesR.Val,
			"full_duplex":  fullDuplex,
			"max_speed":    wan.MaxSpeed.Val,
			"rx_bytes":     wan.RxBytes.Val,
			"rx_bytes-r":   wan.RxBytesR.Val,
			"rx_dropped":   wan.RxDropped.Val,
			"rx_errors":    wan.RxErrors.Val,
			"rx_broadcast": wan.RxBroadcast.Val,
			"rx_multicast": wan.RxMulticast.Val,
			"rx_packets":   wan.RxPackets.Val,
			"speed":        wan.Speed.Val,
			"tx_bytes":     wan.TxBytes.Val,
			"tx_bytes-r":   wan.TxBytesR.Val,
			"tx_dropped":   wan.TxDropped.Val,
			"tx_errors":    wan.TxErrors.Val,
			"tx_packets":   wan.TxPackets.Val,
			"tx_broadcast": wan.TxBroadcast.Val,
			"tx_multicast": wan.TxMulticast.Val,
		}
		metricName := metricNamespace("usg_wan_ports")
		reportGaugeForMap(r, metricName, data, tags)
	}
}

func (u *DatadogUnifi) reportNetTable(r report, deviceName string, siteName string, source string, nt unifi.NetworkTable) {
	for _, p := range nt {
		tags := []string{
			tag("device_name", deviceName),
			tag("site_name", siteName),
			tag("source", source),
			tag("up", p.Up.Txt),
			tag("enabled", p.Enabled.Txt),
			tag("ip", p.IP),
			tag("mac", p.Mac),
			tag("name", p.Name),
			tag("domain_name", p.DomainName),
			tag("purpose", p.Purpose),
			tag("is_guest", p.IsGuest.Txt),
		}
		data := map[string]float64{
			"num_sta":    p.NumSta.Val,
			"rx_bytes":   p.RxBytes.Val,
			"rx_packets": p.RxPackets.Val,
			"tx_bytes":   p.TxBytes.Val,
			"tx_packets": p.TxPackets.Val,
		}
		metricName := metricNamespace("usg_networks")
		reportGaugeForMap(r, metricName, data, tags)
	}
}
