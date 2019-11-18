package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type usg struct {
	Uptime         *prometheus.Desc
	Temperature    *prometheus.Desc
	TotalMaxPower  *prometheus.Desc
	FanLevel       *prometheus.Desc
	TotalTxBytes   *prometheus.Desc
	TotalRxBytes   *prometheus.Desc
	TotalBytes     *prometheus.Desc
	NumSta         *prometheus.Desc
	UserNumSta     *prometheus.Desc
	GuestNumSta    *prometheus.Desc
	NumDesktop     *prometheus.Desc
	NumMobile      *prometheus.Desc
	NumHandheld    *prometheus.Desc
	Loadavg1       *prometheus.Desc
	Loadavg5       *prometheus.Desc
	Loadavg15      *prometheus.Desc
	MemBuffer      *prometheus.Desc
	MemTotal       *prometheus.Desc
	MemUsed        *prometheus.Desc
	CPU            *prometheus.Desc
	Mem            *prometheus.Desc
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
	Latency        *prometheus.Desc
	Runtime        *prometheus.Desc
	XputDownload   *prometheus.Desc
	XputUpload     *prometheus.Desc
}

func descUSG(ns string) *usg {
	if ns += "_usg_"; ns == "_usg_" {
		ns = "usg_"
	}
	labels := []string{"site_name", "mac", "model", "name", "serial", "site_id",
		"type", "version", "device_id", "ip"}
	labelWan := append([]string{"port"}, labels...)

	return &usg{
		Uptime:         prometheus.NewDesc(ns+"uptime", "Uptime", labels, nil),
		TotalTxBytes:   prometheus.NewDesc(ns+"tx_bytes_total", "Total Transmitted Bytes", labels, nil),
		TotalRxBytes:   prometheus.NewDesc(ns+"rx_bytes_total", "Total Received Bytes", labels, nil),
		TotalBytes:     prometheus.NewDesc(ns+"bytes_total", "Total Bytes Transfered", labels, nil),
		NumSta:         prometheus.NewDesc(ns+"stations_total", "Number of Stations", labels, nil),
		UserNumSta:     prometheus.NewDesc(ns+"stations_user_total", "Number of User Stations", labels, nil),
		GuestNumSta:    prometheus.NewDesc(ns+"stations_guest_total", "Number of Guest Stations", labels, nil),
		NumDesktop:     prometheus.NewDesc(ns+"desktops_total", "Number of Desktops", labels, nil),
		NumMobile:      prometheus.NewDesc(ns+"mobile_total", "Number of Mobiles", labels, nil),
		NumHandheld:    prometheus.NewDesc(ns+"handheld_total", "Number of Handhelds", labels, nil),
		Loadavg1:       prometheus.NewDesc(ns+"load_average_1", "System Load Average 1 Minute", labels, nil),
		Loadavg5:       prometheus.NewDesc(ns+"load_average_5", "System Load Average 5 Minutes", labels, nil),
		Loadavg15:      prometheus.NewDesc(ns+"load_average_15", "System Load Average 15 Minutes", labels, nil),
		MemUsed:        prometheus.NewDesc(ns+"memory_used", "System Memory Used", labels, nil),
		MemTotal:       prometheus.NewDesc(ns+"memory_installed", "System Installed Memory", labels, nil),
		MemBuffer:      prometheus.NewDesc(ns+"memory_buffer", "System Memory Buffer", labels, nil),
		CPU:            prometheus.NewDesc(ns+"cpu_utilization", "System CPU % Utilized", labels, nil),
		Mem:            prometheus.NewDesc(ns+"memory_utilization", "System Memory % Utilized", labels, nil), // this may not be %.
		WanRxPackets:   prometheus.NewDesc(ns+"wan_rx_packets_total", "WAN Receive Packets Total", labelWan, nil),
		WanRxBytes:     prometheus.NewDesc(ns+"wan_rx_bytes_total", "WAN Receive Bytes Total", labelWan, nil),
		WanRxDropped:   prometheus.NewDesc(ns+"wan_rx_dropped_total", "WAN Receive Dropped Total", labelWan, nil),
		WanRxErrors:    prometheus.NewDesc(ns+"wan_rx_errors_total", "WAN Receive Errors Total", labelWan, nil),
		WanTxPackets:   prometheus.NewDesc(ns+"wan_tx_packets_total", "WAN Transmit Packets Total", labelWan, nil),
		WanTxBytes:     prometheus.NewDesc(ns+"wan_tx_bytes_total", "WAN Transmit Bytes Total", labelWan, nil),
		WanRxBroadcast: prometheus.NewDesc(ns+"wan_rx_broadcast_total", "WAN Receive Broadcast Total", labelWan, nil),
		WanRxBytesR:    prometheus.NewDesc(ns+"wan_rx_bytes_rate", "WAN Receive Bytes Rate", labelWan, nil),
		WanRxMulticast: prometheus.NewDesc(ns+"wan_rx_multicast_total", "WAN Receive Multicast Total", labelWan, nil),
		WanSpeed:       prometheus.NewDesc(ns+"wan_speed", "WAN Speed", labelWan, nil),
		WanTxBroadcast: prometheus.NewDesc(ns+"wan_tx_broadcast_total", "WAN Transmit Broadcast Total", labelWan, nil),
		WanTxBytesR:    prometheus.NewDesc(ns+"wan_tx_bytes_rate", "WAN Transmit Bytes Rate", labelWan, nil),
		WanTxDropped:   prometheus.NewDesc(ns+"wan_tx_dropped_total", "WAN Transmit Dropped Total", labelWan, nil),
		WanTxErrors:    prometheus.NewDesc(ns+"wan_tx_errors_total", "WAN Transmit Errors Total", labelWan, nil),
		WanTxMulticast: prometheus.NewDesc(ns+"wan_tx_multicast_total", "WAN Transmit Multicast Total", labelWan, nil),
		WanBytesR:      prometheus.NewDesc(ns+"wan_bytes_rate", "WAN Transfer Rate", labelWan, nil),
		LanRxPackets:   prometheus.NewDesc(ns+"lan_rx_packets_total", "LAN Receive Packets Total", labels, nil),
		LanRxBytes:     prometheus.NewDesc(ns+"lan_rx_bytes_total", "LAN Receive Bytes Total", labels, nil),
		LanRxDropped:   prometheus.NewDesc(ns+"lan_rx_dropped_total", "LAN Receive Dropped Total", labels, nil),
		LanTxPackets:   prometheus.NewDesc(ns+"lan_tx_packets_total", "LAN Transmit Packets Total", labels, nil),
		LanTxBytes:     prometheus.NewDesc(ns+"lan_tx_bytes_total", "LAN Transmit Bytes Total", labels, nil),
		Latency:        prometheus.NewDesc(ns+"speedtest_latency", "Speedtest Latency", labels, nil),
		Runtime:        prometheus.NewDesc(ns+"speedtest_runtime", "Speedtest Run Time", labels, nil),
		XputDownload:   prometheus.NewDesc(ns+"speedtest_download_rate", "Speedtest Download Rate", labels, nil),
		XputUpload:     prometheus.NewDesc(ns+"speedtest_upload_rate", "Speedtest Upload Rate", labels, nil),
	}
}

// exportUSG Exports Security Gateway Data
// uplink and port tables structs are ignored. that data should be in other exported fields.
func (u *unifiCollector) exportUSG(s *unifi.USG) []*metricExports {
	labels := []string{s.SiteName, s.Mac, s.Model, s.Name, s.Serial, s.SiteID,
		s.Type, s.Version, s.DeviceID, s.IP}
	labelWan := append([]string{"all"}, labels...)

	// Gateway System Data.
	return append([]*metricExports{
		{u.USG.Uptime, prometheus.GaugeValue, s.Uptime, labels},
		{u.USG.TotalTxBytes, prometheus.CounterValue, s.TxBytes, labels},
		{u.USG.TotalRxBytes, prometheus.CounterValue, s.RxBytes, labels},
		{u.USG.TotalBytes, prometheus.CounterValue, s.Bytes, labels},
		{u.USG.NumSta, prometheus.GaugeValue, s.NumSta, labels},
		{u.USG.UserNumSta, prometheus.GaugeValue, s.UserNumSta, labels},
		{u.USG.GuestNumSta, prometheus.GaugeValue, s.GuestNumSta, labels},
		{u.USG.NumDesktop, prometheus.CounterValue, s.NumDesktop, labels},
		{u.USG.NumMobile, prometheus.CounterValue, s.NumMobile, labels},
		{u.USG.NumHandheld, prometheus.CounterValue, s.NumHandheld, labels},
		{u.USG.Loadavg1, prometheus.GaugeValue, s.SysStats.Loadavg1, labels},
		{u.USG.Loadavg5, prometheus.GaugeValue, s.SysStats.Loadavg5, labels},
		{u.USG.Loadavg15, prometheus.GaugeValue, s.SysStats.Loadavg15, labels},
		{u.USG.MemUsed, prometheus.GaugeValue, s.SysStats.MemUsed, labels},
		{u.USG.MemTotal, prometheus.GaugeValue, s.SysStats.MemTotal, labels},
		{u.USG.MemBuffer, prometheus.GaugeValue, s.SysStats.MemBuffer, labels},
		{u.USG.CPU, prometheus.GaugeValue, s.SystemStats.CPU, labels},
		{u.USG.Mem, prometheus.GaugeValue, s.SystemStats.Mem, labels},
		// Combined Port Stats
		{u.USG.WanRxPackets, prometheus.CounterValue, s.Stat.Gw.WanRxPackets, labelWan},
		{u.USG.WanRxBytes, prometheus.CounterValue, s.Stat.Gw.WanRxBytes, labelWan},
		{u.USG.WanRxDropped, prometheus.CounterValue, s.Stat.Gw.WanRxDropped, labelWan},
		{u.USG.WanTxPackets, prometheus.CounterValue, s.Stat.Gw.WanTxPackets, labelWan},
		{u.USG.WanTxBytes, prometheus.CounterValue, s.Stat.Gw.WanTxBytes, labelWan},
		{u.USG.WanRxErrors, prometheus.CounterValue, s.Stat.Gw.WanRxErrors, labelWan},
		{u.USG.LanRxPackets, prometheus.CounterValue, s.Stat.Gw.LanRxPackets, labels},
		{u.USG.LanRxBytes, prometheus.CounterValue, s.Stat.Gw.LanRxBytes, labels},
		{u.USG.LanTxPackets, prometheus.CounterValue, s.Stat.Gw.LanTxPackets, labels},
		{u.USG.LanTxBytes, prometheus.CounterValue, s.Stat.Gw.LanTxBytes, labels},
		{u.USG.LanRxDropped, prometheus.CounterValue, s.Stat.Gw.LanRxDropped, labels},
		// Speed Test Stats
		{u.USG.Latency, prometheus.GaugeValue, s.SpeedtestStatus.Latency, labels},
		{u.USG.Runtime, prometheus.GaugeValue, s.SpeedtestStatus.Runtime, labels},
		{u.USG.XputDownload, prometheus.GaugeValue, s.SpeedtestStatus.XputDownload, labels},
		{u.USG.XputUpload, prometheus.GaugeValue, s.SpeedtestStatus.XputUpload, labels},
	}, u.exportWANPorts(labels, s.Wan1, s.Wan2)...)
}

func (u *unifiCollector) exportWANPorts(labels []string, wans ...unifi.Wan) []*metricExports {
	var metrics []*metricExports
	for _, wan := range wans {
		if !wan.Up.Val {
			continue // only record UP interfaces.
		}
		l := append([]string{wan.Name}, labels...)

		metrics = append(metrics, []*metricExports{
			{u.USG.WanRxPackets, prometheus.CounterValue, wan.RxPackets, l},
			{u.USG.WanRxBytes, prometheus.CounterValue, wan.RxBytes, l},
			{u.USG.WanRxDropped, prometheus.CounterValue, wan.RxDropped, l},
			{u.USG.WanRxErrors, prometheus.CounterValue, wan.RxErrors, l},
			{u.USG.WanTxPackets, prometheus.CounterValue, wan.TxPackets, l},
			{u.USG.WanTxBytes, prometheus.CounterValue, wan.TxBytes, l},
			{u.USG.WanRxBroadcast, prometheus.CounterValue, wan.RxBroadcast, l},
			{u.USG.WanRxMulticast, prometheus.CounterValue, wan.RxMulticast, l},
			{u.USG.WanSpeed, prometheus.CounterValue, wan.Speed, l},
			{u.USG.WanTxBroadcast, prometheus.CounterValue, wan.TxBroadcast, l},
			{u.USG.WanTxBytesR, prometheus.CounterValue, wan.TxBytesR, l},
			{u.USG.WanTxDropped, prometheus.CounterValue, wan.TxDropped, l},
			{u.USG.WanTxErrors, prometheus.CounterValue, wan.TxErrors, l},
			{u.USG.WanTxMulticast, prometheus.CounterValue, wan.TxMulticast, l},
			{u.USG.WanBytesR, prometheus.GaugeValue, wan.BytesR, l},
		}...)
	}

	return metrics
}
