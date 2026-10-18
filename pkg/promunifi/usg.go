package promunifi

import (
	"slices"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v6"
)

type usg struct {
	WanRxPackets   *prometheus.Desc
	WanRxBytes     *prometheus.Desc
	WanRxDropped   *prometheus.Desc
	WanRxErrors    *prometheus.Desc
	WanTxPackets   *prometheus.Desc
	WanTxBytes     *prometheus.Desc
	LanRxPackets   *prometheus.Desc
	LanRxBytes     *prometheus.Desc
	LanRxDropped   *prometheus.Desc
	LanTxPackets   *prometheus.Desc
	LanTxBytes     *prometheus.Desc
	WanRxBroadcast *prometheus.Desc
	WanRxBytesR    *prometheus.Desc
	WanRxMulticast *prometheus.Desc
	WanSpeed       *prometheus.Desc
	WanTxBroadcast *prometheus.Desc
	WanTxBytesR    *prometheus.Desc
	WanTxDropped   *prometheus.Desc
	WanTxErrors    *prometheus.Desc
	WanTxMulticast *prometheus.Desc
	WanBytesR      *prometheus.Desc
	WanLatency     *prometheus.Desc
	WanAvailable   *prometheus.Desc
	WanStatus      *prometheus.Desc
	StatsUptime    *prometheus.Desc
	StatsDowntime  *prometheus.Desc
	MonAvailable   *prometheus.Desc
	MonLatency     *prometheus.Desc
	Latency        *prometheus.Desc
	UplinkLatency  *prometheus.Desc
	UplinkSpeed    *prometheus.Desc
	UplinkMaxSpeed *prometheus.Desc
	UplinkUptime   *prometheus.Desc
	Runtime        *prometheus.Desc
	Rundate        *prometheus.Desc
	XputDownload   *prometheus.Desc
	XputUpload     *prometheus.Desc
}

func descUSG(ns string) *usg {
	labels := []string{"port", "site_name", "name", "source", "tag"}
	statsLabels := []string{"wan_networkgroup", "site_name", "name", "source", "tag"}
	statusLabels := []string{"wan_networkgroup", "site_name", "name", "source", "tag", "status"}
	monLabels := []string{"wan_networkgroup", "site_name", "name", "source", "tag", "target", "monitor_type", "alerting"}

	return &usg{
		WanRxPackets:   prometheus.NewDesc(ns+"wan_receive_packets_total", "WAN Receive Packets Total", labels, nil),
		WanRxBytes:     prometheus.NewDesc(ns+"wan_receive_bytes_total", "WAN Receive Bytes Total", labels, nil),
		WanRxDropped:   prometheus.NewDesc(ns+"wan_receive_dropped_total", "WAN Receive Dropped Total", labels, nil),
		WanRxErrors:    prometheus.NewDesc(ns+"wan_receive_errors_total", "WAN Receive Errors Total", labels, nil),
		WanTxPackets:   prometheus.NewDesc(ns+"wan_transmit_packets_total", "WAN Transmit Packets Total", labels, nil),
		WanTxBytes:     prometheus.NewDesc(ns+"wan_transmit_bytes_total", "WAN Transmit Bytes Total", labels, nil),
		WanRxBroadcast: prometheus.NewDesc(ns+"wan_receive_broadcast_total", "WAN Receive Broadcast Total", labels, nil),
		WanRxBytesR:    prometheus.NewDesc(ns+"wan_receive_rate_bytes", "WAN Receive Bytes Rate", labels, nil),
		WanRxMulticast: prometheus.NewDesc(ns+"wan_receive_multicast_total", "WAN Receive Multicast Total", labels, nil),
		WanSpeed:       prometheus.NewDesc(ns+"wan_speed_bps", "WAN Speed", labels, nil),
		WanTxBroadcast: prometheus.NewDesc(ns+"wan_transmit_broadcast_total", "WAN Transmit Broadcast Total", labels, nil),
		WanTxBytesR:    prometheus.NewDesc(ns+"wan_transmit_rate_bytes", "WAN Transmit Bytes Rate", labels, nil),
		WanTxDropped:   prometheus.NewDesc(ns+"wan_transmit_dropped_total", "WAN Transmit Dropped Total", labels, nil),
		WanTxErrors:    prometheus.NewDesc(ns+"wan_transmit_errors_total", "WAN Transmit Errors Total", labels, nil),
		WanTxMulticast: prometheus.NewDesc(ns+"wan_transmit_multicast_total", "WAN Transmit Multicast Total", labels, nil),
		WanBytesR:      prometheus.NewDesc(ns+"wan_rate_bytes", "WAN Transfer Rate", labels, nil),
		WanLatency:     prometheus.NewDesc(ns+"wan_latency_seconds", "WAN Latency", labels, nil),
		WanAvailable:   prometheus.NewDesc(ns+"wan_availability_percent", "WAN Availability", labels, nil),
		StatsUptime:    prometheus.NewDesc(ns+"wan_uptime_stats_uptime_seconds", "WAN Uptime", statsLabels, nil),
		StatsDowntime:  prometheus.NewDesc(ns+"wan_uptime_stats_downtime_seconds", "WAN Downtime", statsLabels, nil),
		WanStatus:      prometheus.NewDesc(ns+"wan_status", "WAN Status: 1=online, 0=other", statusLabels, nil),
		MonAvailable:   prometheus.NewDesc(ns+"wan_monitor_availability_percent", "WAN Monitor Availability", monLabels, nil),
		MonLatency:     prometheus.NewDesc(ns+"wan_monitor_latency_average_seconds", "WAN Monitor Average Latency", monLabels, nil),
		LanRxPackets:   prometheus.NewDesc(ns+"lan_receive_packets_total", "LAN Receive Packets Total", labels, nil),
		LanRxBytes:     prometheus.NewDesc(ns+"lan_receive_bytes_total", "LAN Receive Bytes Total", labels, nil),
		LanRxDropped:   prometheus.NewDesc(ns+"lan_receive_dropped_total", "LAN Receive Dropped Total", labels, nil),
		LanTxPackets:   prometheus.NewDesc(ns+"lan_transmit_packets_total", "LAN Transmit Packets Total", labels, nil),
		LanTxBytes:     prometheus.NewDesc(ns+"lan_transmit_bytes_total", "LAN Transmit Bytes Total", labels, nil),
		UplinkLatency:  prometheus.NewDesc(ns+"uplink_latency_seconds", "Uplink Latency", labels, nil),
		UplinkSpeed:    prometheus.NewDesc(ns+"uplink_speed_mbps", "Uplink Speed", labels, nil),
		UplinkMaxSpeed: prometheus.NewDesc(ns+"uplink_max_speed_mbps", "Uplink Max Speed", labels, nil),
		UplinkUptime:   prometheus.NewDesc(ns+"uplink_uptime_seconds", "Uplink Uptime", labels, nil),
		Latency:        prometheus.NewDesc(ns+"speedtest_latency_seconds", "Speedtest Latency", labels, nil),
		Runtime:        prometheus.NewDesc(ns+"speedtest_runtime_seconds", "Speedtest Run Time", labels, nil),
		Rundate:        prometheus.NewDesc(ns+"speedtest_rundate_seconds", "Speedtest Run Date", labels, nil),
		XputDownload:   prometheus.NewDesc(ns+"speedtest_download", "Speedtest Download Rate", labels, nil),
		XputUpload:     prometheus.NewDesc(ns+"speedtest_upload", "Speedtest Upload Rate", labels, nil),
	}
}

func (u *promUnifi) exportUSG(r report, d *unifi.USG) {
	if !d.Adopted.Val {
		return
	}

	baseLabels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	baseInfoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}

	u.exportWithTags(r, d.Tags, func(tagLabels []string) {
		tag := tagLabels[0]
		labels := slices.Clip(append(baseLabels, tag))
		infoLabels := append(baseInfoLabels, tag)

		for _, t := range d.Temperatures {
			r.send([]*metric{{u.Device.Temperature, gauge, t.Value, append(labels, t.Name, t.Type)}})
		}

		for k, v := range d.SystemStats.Temps {
			temp := v.CelsiusInt64()
			k = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(k, " ", "_"), ")", ""), "(", "")

			if k = strings.ToLower(k); temp != 0 && k != "" {
				r.send([]*metric{{u.Device.Temperature, gauge, temp, append(labels, k, k)}})
			}
		}

		// Gateway System Data.
		u.exportWANPorts(r, labels, d.Wan1, d.Wan2)
		u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)
		u.exportSYSstats(r, labels, d.SysStats, d.SystemStats)
		u.exportUSGstats(r, labels, d.Stat.Gw, d.SpeedtestStatus, d.Uplink)
		u.exportSTAcount(r, labels, d.UserNumSta, d.GuestNumSta, d.NumDesktop, d.NumMobile, d.NumHandheld)
		r.send([]*metric{
			{u.Device.Info, gauge, 1.0, append(baseLabels, infoLabels...)},
			{u.Device.Uptime, gauge, d.Uptime, labels},
			{u.Device.Upgradeable, gauge, d.Upgradable.Val, labels},
		})
	})
}

// exportDeviceUplink emits the shared uplink metrics for any device that
// carries a unifi.Uplink (USW/UBB/UDB). The uplink port name is used as the
// first label, matching the convention used by exportUSGstats.
func (u *promUnifi) exportDeviceUplink(r report, labels []string, ul unifi.Uplink) {
	port := ul.Name
	if port == "" {
		port = "uplink"
	}

	labelUL := []string{port, labels[1], labels[2], labels[3], labels[4]}

	r.send([]*metric{
		{u.USG.UplinkLatency, gauge, ul.Latency.Val / 1000, labelUL},
		{u.USG.UplinkSpeed, gauge, ul.Speed, labelUL},
		{u.USG.UplinkMaxSpeed, gauge, ul.MaxSpeed, labelUL},
		{u.USG.UplinkUptime, gauge, ul.Uptime, labelUL},
	})
}

// Gateway Stats.
func (u *promUnifi) exportUSGstats(r report, labels []string, gw *unifi.Gw, st unifi.SpeedtestStatus, ul unifi.Uplink) {
	var sourceInterface string

	if st.SourceInterface != "" {
		sourceInterface = st.SourceInterface
	} else {
		sourceInterface = "all"
	}

	// Uplink and speedtest live on the device, not in stat.gw. Network 11 omits
	// stat entirely, so a missing gateway stat must not drop those series.
	labelWan := []string{sourceInterface, labels[1], labels[2], labels[3], labels[4]}

	r.send([]*metric{
		{u.USG.UplinkLatency, gauge, ul.Latency.Val / 1000, labelWan},
		{u.USG.UplinkSpeed, gauge, ul.Speed, labelWan},
		{u.USG.UplinkMaxSpeed, gauge, ul.MaxSpeed, labelWan},
		{u.USG.UplinkUptime, gauge, ul.Uptime, labelWan},
		{u.USG.Latency, gauge, st.Latency.Val / 1000, labelWan},
		{u.USG.Runtime, gauge, st.Runtime, labelWan},
		{u.USG.Rundate, gauge, st.Rundate, labelWan},
		{u.USG.XputDownload, gauge, st.XputDownload, labelWan},
		{u.USG.XputUpload, gauge, st.XputUpload, labelWan},
	})

	if gw == nil {
		return
	}

	labelLan := []string{"lan", labels[1], labels[2], labels[3], labels[4]}

	r.send([]*metric{
		{u.USG.LanRxPackets, counter, gw.LanRxPackets, labelLan},
		{u.USG.LanRxBytes, counter, gw.LanRxBytes, labelLan},
		{u.USG.LanTxPackets, counter, gw.LanTxPackets, labelLan},
		{u.USG.LanTxBytes, counter, gw.LanTxBytes, labelLan},
		{u.USG.LanRxDropped, counter, gw.LanRxDropped, labelLan},
	})
}

// WAN Stats.
func (u *promUnifi) exportWANPorts(r report, labels []string, wans ...unifi.Wan) {
	for _, wan := range wans {
		if !wan.Up.Val {
			continue // only record UP interfaces.
		}

		labelWan := []string{wan.Name, labels[1], labels[2], labels[3], labels[4]}

		r.send([]*metric{
			{u.USG.WanRxPackets, counter, wan.RxPackets, labelWan},
			{u.USG.WanRxBytes, counter, wan.RxBytes, labelWan},
			{u.USG.WanRxDropped, counter, wan.RxDropped, labelWan},
			{u.USG.WanRxErrors, counter, wan.RxErrors, labelWan},
			{u.USG.WanTxPackets, counter, wan.TxPackets, labelWan},
			{u.USG.WanTxBytes, counter, wan.TxBytes, labelWan},
			{u.USG.WanRxBroadcast, counter, wan.RxBroadcast, labelWan},
			{u.USG.WanRxMulticast, counter, wan.RxMulticast, labelWan},
			{u.USG.WanSpeed, gauge, wan.Speed.Val * 1000000, labelWan},
			{u.USG.WanTxBroadcast, counter, wan.TxBroadcast, labelWan},
			{u.USG.WanTxBytesR, gauge, wan.TxBytesR, labelWan},
			{u.USG.WanRxBytesR, gauge, wan.RxBytesR, labelWan},
			{u.USG.WanTxDropped, counter, wan.TxDropped, labelWan},
			{u.USG.WanTxErrors, counter, wan.TxErrors, labelWan},
			{u.USG.WanTxMulticast, counter, wan.TxMulticast, labelWan},
			{u.USG.WanBytesR, gauge, wan.BytesR, labelWan},
		})

		// Older gateways do not measure WAN health; don't report them as 0% available.
		// A WAN that is offline reports latency as 0; don't report that as perfect latency.
		if wan.Latency.Val > 0 {
			r.send([]*metric{{u.USG.WanLatency, gauge, wan.Latency.Val / 1000, labelWan}})
		}

		if wan.Availability.Txt != "" {
			r.send([]*metric{{u.USG.WanAvailable, gauge, wan.Availability, labelWan}})
		}
	}
}

// exportWANHealth emits per-WAN health keyed by WAN network group (WAN, WAN2).
// Unlike exportWANPorts, it does not skip a WAN whose port is down.
// uptime_stats availability is left out: it is the same statistic as unpoller_wan_uptime_percentage.
// uptime_stats latency_average is left out too: its averaging window varies in length.
func (u *promUnifi) exportWANHealth(r report, labels []string, stats map[string]*unifi.UptimeStats, status map[string]string) {
	for wan, s := range status {
		labelStatus := []string{wan, labels[1], labels[2], labels[3], labels[4], s}
		r.send([]*metric{{u.USG.WanStatus, gauge, s == "online", labelStatus}})
	}

	for wan, s := range stats {
		if s == nil {
			continue
		}

		labelWan := []string{wan, labels[1], labels[2], labels[3], labels[4]}

		// While a WAN is offline the controller omits uptime and reports downtime instead.
		if s.Uptime.Txt != "" {
			r.send([]*metric{{u.USG.StatsUptime, gauge, s.Uptime, labelWan}})
		}

		if s.Downtime.Txt != "" {
			r.send([]*metric{{u.USG.StatsDowntime, gauge, s.Downtime, labelWan}})
		}

		u.exportUptimeMonitors(r, labelWan, "false", s.Monitors)
		u.exportUptimeMonitors(r, labelWan, "true", s.AlertingMonitors)
	}
}

func (u *promUnifi) exportUptimeMonitors(r report, labelWan []string, alerting string, monitors []unifi.UptimeMonitor) {
	for _, m := range monitors {
		labelMon := slices.Concat(labelWan, []string{m.Target, m.Type, alerting})

		if m.Availability.Txt != "" {
			r.send([]*metric{{u.USG.MonAvailable, gauge, m.Availability, labelMon}})
		}

		// A failing monitor reports latency as 0 or omits it.
		if m.LatencyAverage.Val > 0 {
			r.send([]*metric{{u.USG.MonLatency, gauge, m.LatencyAverage.Val / 1000, labelMon}})
		}
	}
}
