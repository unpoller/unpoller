package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
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
	Latency        *prometheus.Desc
	Runtime        *prometheus.Desc
	XputDownload   *prometheus.Desc
	XputUpload     *prometheus.Desc
}

func descUSG(ns string) *usg {
	labels := []string{"ip", "type", "version", "site_name", "mac", "model", "name", "serial"}
	labelWan := append([]string{"port"}, labels...)
	return &usg{
		WanRxPackets:   prometheus.NewDesc(ns+"wan_receive_packets_total", "WAN Receive Packets Total", labelWan, nil),
		WanRxBytes:     prometheus.NewDesc(ns+"wan_receive_bytes_total", "WAN Receive Bytes Total", labelWan, nil),
		WanRxDropped:   prometheus.NewDesc(ns+"wan_receive_dropped_total", "WAN Receive Dropped Total", labelWan, nil),
		WanRxErrors:    prometheus.NewDesc(ns+"wan_receive_errors_total", "WAN Receive Errors Total", labelWan, nil),
		WanTxPackets:   prometheus.NewDesc(ns+"wan_transmit_packets_total", "WAN Transmit Packets Total", labelWan, nil),
		WanTxBytes:     prometheus.NewDesc(ns+"wan_transmit_bytes_total", "WAN Transmit Bytes Total", labelWan, nil),
		WanRxBroadcast: prometheus.NewDesc(ns+"wan_receive_broadcast_total", "WAN Receive Broadcast Total", labelWan, nil),
		WanRxBytesR:    prometheus.NewDesc(ns+"wan_receive_rate_bytes", "WAN Receive Bytes Rate", labelWan, nil),
		WanRxMulticast: prometheus.NewDesc(ns+"wan_receive_multicast_total", "WAN Receive Multicast Total", labelWan, nil),
		WanSpeed:       prometheus.NewDesc(ns+"wan_speed_bps", "WAN Speed", labelWan, nil),
		WanTxBroadcast: prometheus.NewDesc(ns+"wan_transmit_broadcast_total", "WAN Transmit Broadcast Total", labelWan, nil),
		WanTxBytesR:    prometheus.NewDesc(ns+"wan_transmit_rate_bytes", "WAN Transmit Bytes Rate", labelWan, nil),
		WanTxDropped:   prometheus.NewDesc(ns+"wan_transmit_dropped_total", "WAN Transmit Dropped Total", labelWan, nil),
		WanTxErrors:    prometheus.NewDesc(ns+"wan_transmit_errors_total", "WAN Transmit Errors Total", labelWan, nil),
		WanTxMulticast: prometheus.NewDesc(ns+"wan_transmit_multicast_total", "WAN Transmit Multicast Total", labelWan, nil),
		WanBytesR:      prometheus.NewDesc(ns+"wan_rate_bytes", "WAN Transfer Rate", labelWan, nil),
		LanRxPackets:   prometheus.NewDesc(ns+"lan_receive_packets_total", "LAN Receive Packets Total", labels, nil),
		LanRxBytes:     prometheus.NewDesc(ns+"lan_receive_bytes_total", "LAN Receive Bytes Total", labels, nil),
		LanRxDropped:   prometheus.NewDesc(ns+"lan_receive_dropped_total", "LAN Receive Dropped Total", labels, nil),
		LanTxPackets:   prometheus.NewDesc(ns+"lan_transmit_packets_total", "LAN Transmit Packets Total", labels, nil),
		LanTxBytes:     prometheus.NewDesc(ns+"lan_transmit_bytes_total", "LAN Transmit Bytes Total", labels, nil),
		Latency:        prometheus.NewDesc(ns+"speedtest_latency_seconds", "Speedtest Latency", labels, nil),
		Runtime:        prometheus.NewDesc(ns+"speedtest_runtime", "Speedtest Run Time", labels, nil),
		XputDownload:   prometheus.NewDesc(ns+"speedtest_download", "Speedtest Download Rate", labels, nil),
		XputUpload:     prometheus.NewDesc(ns+"speedtest_upload", "Speedtest Upload Rate", labels, nil),
	}
}

func (u *unifiCollector) exportUSGs(r report) {
	if r.metrics() == nil || r.metrics().Devices == nil || len(r.metrics().Devices.USGs) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, d := range r.metrics().Devices.USGs {
			u.exportUSG(r, d)
		}
	}()
}

func (u *unifiCollector) exportUSG(r report, d *unifi.USG) {
	labels := []string{d.IP, d.Type, d.Version, d.SiteName, d.Mac, d.Model, d.Name, d.Serial}
	// Gateway System Data.
	r.send([]*metricExports{
		{u.Device.Uptime, prometheus.GaugeValue, d.Uptime, labels},
		{u.Device.TotalTxBytes, prometheus.CounterValue, d.TxBytes, labels},
		{u.Device.TotalRxBytes, prometheus.CounterValue, d.RxBytes, labels},
		{u.Device.TotalBytes, prometheus.CounterValue, d.Bytes, labels},
		{u.Device.NumSta, prometheus.GaugeValue, d.NumSta, labels},
		{u.Device.UserNumSta, prometheus.GaugeValue, d.UserNumSta, labels},
		{u.Device.GuestNumSta, prometheus.GaugeValue, d.GuestNumSta, labels},
		{u.Device.NumDesktop, prometheus.GaugeValue, d.NumDesktop, labels},
		{u.Device.NumMobile, prometheus.GaugeValue, d.NumMobile, labels},
		{u.Device.NumHandheld, prometheus.GaugeValue, d.NumHandheld, labels},
		{u.Device.Loadavg1, prometheus.GaugeValue, d.SysStats.Loadavg1, labels},
		{u.Device.Loadavg5, prometheus.GaugeValue, d.SysStats.Loadavg5, labels},
		{u.Device.Loadavg15, prometheus.GaugeValue, d.SysStats.Loadavg15, labels},
		{u.Device.MemUsed, prometheus.GaugeValue, d.SysStats.MemUsed, labels},
		{u.Device.MemTotal, prometheus.GaugeValue, d.SysStats.MemTotal, labels},
		{u.Device.MemBuffer, prometheus.GaugeValue, d.SysStats.MemBuffer, labels},
		{u.Device.CPU, prometheus.GaugeValue, d.SystemStats.CPU, labels},
		{u.Device.Mem, prometheus.GaugeValue, d.SystemStats.Mem, labels},
	})
	u.exportWANPorts(r, labels, d.Wan1, d.Wan2)
	u.exportUSGstats(r, labels, d.Stat.Gw, d.SpeedtestStatus)
}

func (u *unifiCollector) exportUSGstats(r report, labels []string, gw *unifi.Gw, st unifi.SpeedtestStatus) {
	labelWan := append([]string{"all"}, labels...)
	r.send([]*metricExports{
		// Combined Port Stats
		{u.USG.WanRxPackets, prometheus.CounterValue, gw.WanRxPackets, labelWan},
		{u.USG.WanRxBytes, prometheus.CounterValue, gw.WanRxBytes, labelWan},
		{u.USG.WanRxDropped, prometheus.CounterValue, gw.WanRxDropped, labelWan},
		{u.USG.WanTxPackets, prometheus.CounterValue, gw.WanTxPackets, labelWan},
		{u.USG.WanTxBytes, prometheus.CounterValue, gw.WanTxBytes, labelWan},
		{u.USG.WanRxErrors, prometheus.CounterValue, gw.WanRxErrors, labelWan},
		{u.USG.LanRxPackets, prometheus.CounterValue, gw.LanRxPackets, labels},
		{u.USG.LanRxBytes, prometheus.CounterValue, gw.LanRxBytes, labels},
		{u.USG.LanTxPackets, prometheus.CounterValue, gw.LanTxPackets, labels},
		{u.USG.LanTxBytes, prometheus.CounterValue, gw.LanTxBytes, labels},
		{u.USG.LanRxDropped, prometheus.CounterValue, gw.LanRxDropped, labels},
		// Speed Test Stats
		{u.USG.Latency, prometheus.GaugeValue, st.Latency.Val / 1000, labels},
		{u.USG.Runtime, prometheus.GaugeValue, st.Runtime, labels},
		{u.USG.XputDownload, prometheus.GaugeValue, st.XputDownload, labels},
		{u.USG.XputUpload, prometheus.GaugeValue, st.XputUpload, labels},
	})
}

func (u *unifiCollector) exportWANPorts(r report, labels []string, wans ...unifi.Wan) {
	for _, wan := range wans {
		if !wan.Up.Val {
			continue // only record UP interfaces.
		}
		l := append([]string{wan.Name}, labels...)
		r.send([]*metricExports{
			{u.USG.WanRxPackets, prometheus.CounterValue, wan.RxPackets, l},
			{u.USG.WanRxBytes, prometheus.CounterValue, wan.RxBytes, l},
			{u.USG.WanRxDropped, prometheus.CounterValue, wan.RxDropped, l},
			{u.USG.WanRxErrors, prometheus.CounterValue, wan.RxErrors, l},
			{u.USG.WanTxPackets, prometheus.CounterValue, wan.TxPackets, l},
			{u.USG.WanTxBytes, prometheus.CounterValue, wan.TxBytes, l},
			{u.USG.WanRxBroadcast, prometheus.CounterValue, wan.RxBroadcast, l},
			{u.USG.WanRxMulticast, prometheus.CounterValue, wan.RxMulticast, l},
			{u.USG.WanSpeed, prometheus.CounterValue, wan.Speed.Val * 1000000, l},
			{u.USG.WanTxBroadcast, prometheus.CounterValue, wan.TxBroadcast, l},
			{u.USG.WanTxBytesR, prometheus.CounterValue, wan.TxBytesR, l},
			{u.USG.WanTxDropped, prometheus.CounterValue, wan.TxDropped, l},
			{u.USG.WanTxErrors, prometheus.CounterValue, wan.TxErrors, l},
			{u.USG.WanTxMulticast, prometheus.CounterValue, wan.TxMulticast, l},
			{u.USG.WanBytesR, prometheus.GaugeValue, wan.BytesR, l},
		})
	}
}
