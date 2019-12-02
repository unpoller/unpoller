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
	UplinkLatency  *prometheus.Desc
	UplinkSpeed    *prometheus.Desc
	Runtime        *prometheus.Desc
	XputDownload   *prometheus.Desc
	XputUpload     *prometheus.Desc
}

func descUSG(ns string) *usg {
	//	labels := []string{"ip", "version", "model", "serial", "type", "mac", "site_name", "name"}
	//	labelWan := append([]string{"port"}, labels[6:]...)
	labels := []string{"port", "site_name", "name"}
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
		LanRxPackets:   prometheus.NewDesc(ns+"lan_receive_packets_total", "LAN Receive Packets Total", labels, nil),
		LanRxBytes:     prometheus.NewDesc(ns+"lan_receive_bytes_total", "LAN Receive Bytes Total", labels, nil),
		LanRxDropped:   prometheus.NewDesc(ns+"lan_receive_dropped_total", "LAN Receive Dropped Total", labels, nil),
		LanTxPackets:   prometheus.NewDesc(ns+"lan_transmit_packets_total", "LAN Transmit Packets Total", labels, nil),
		LanTxBytes:     prometheus.NewDesc(ns+"lan_transmit_bytes_total", "LAN Transmit Bytes Total", labels, nil),
		Latency:        prometheus.NewDesc(ns+"speedtest_latency_seconds", "Speedtest Latency", labels, nil),
		UplinkLatency:  prometheus.NewDesc(ns+"uplink_latency_seconds", "Uplink Latency", labels, nil),
		UplinkSpeed:    prometheus.NewDesc(ns+"uplink_speed_mbps", "Uplink Speed", labels, nil),
		Runtime:        prometheus.NewDesc(ns+"speedtest_runtime", "Speedtest Run Time", labels, nil),
		XputDownload:   prometheus.NewDesc(ns+"speedtest_download", "Speedtest Download Rate", labels, nil),
		XputUpload:     prometheus.NewDesc(ns+"speedtest_upload", "Speedtest Upload Rate", labels, nil),
	}
}

func (u *promUnifi) exportUSG(r report, d *unifi.USG) {
	labels := []string{d.Type, d.SiteName, d.Name}
	infoLabels := []string{d.IP, d.Version, d.Model, d.Serial, d.Mac}
	// Gateway System Data.
	r.send([]*metric{
		{u.Device.Info, prometheus.GaugeValue, 1.0, append(labels, infoLabels...)},
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
	u.exportUSGstats(r, labels, d.Stat.Gw, d.SpeedtestStatus, d.Uplink)
}

func (u *promUnifi) exportUSGstats(r report, labels []string, gw *unifi.Gw, st unifi.SpeedtestStatus, ul unifi.Uplink) {
	labelLan := []string{"lan", labels[6], labels[7]}
	labelWan := []string{"all", labels[6], labels[7]}
	r.send([]*metric{
		/* // Combined Port Stats - not really needed. sum() the others instead.
		{u.USG.WanRxPackets, prometheus.CounterValue, gw.WanRxPackets, labelWan},
		{u.USG.WanRxBytes, prometheus.CounterValue, gw.WanRxBytes, labelWan},
		{u.USG.WanRxDropped, prometheus.CounterValue, gw.WanRxDropped, labelWan},
		{u.USG.WanTxPackets, prometheus.CounterValue, gw.WanTxPackets, labelWan},
		{u.USG.WanTxBytes, prometheus.CounterValue, gw.WanTxBytes, labelWan},
		{u.USG.WanRxErrors, prometheus.CounterValue, gw.WanRxErrors, labelWan},
		*/
		{u.USG.LanRxPackets, prometheus.CounterValue, gw.LanRxPackets, labelLan},
		{u.USG.LanRxBytes, prometheus.CounterValue, gw.LanRxBytes, labelLan},
		{u.USG.LanTxPackets, prometheus.CounterValue, gw.LanTxPackets, labelLan},
		{u.USG.LanTxBytes, prometheus.CounterValue, gw.LanTxBytes, labelLan},
		{u.USG.LanRxDropped, prometheus.CounterValue, gw.LanRxDropped, labelLan},
		{u.USG.UplinkLatency, prometheus.GaugeValue, ul.Latency.Val / 1000, labelWan},
		{u.USG.UplinkSpeed, prometheus.GaugeValue, ul.Speed, labelWan},
		// Speed Test Stats
		{u.USG.Latency, prometheus.GaugeValue, st.Latency.Val / 1000, labelWan},
		{u.USG.Runtime, prometheus.GaugeValue, st.Runtime, labelWan},
		{u.USG.XputDownload, prometheus.GaugeValue, st.XputDownload, labelWan},
		{u.USG.XputUpload, prometheus.GaugeValue, st.XputUpload, labelWan},
	})
}

func (u *promUnifi) exportWANPorts(r report, labels []string, wans ...unifi.Wan) {
	for _, wan := range wans {
		if !wan.Up.Val {
			continue // only record UP interfaces.
		}
		labelWan := []string{wan.Name, labels[6], labels[7]}
		r.send([]*metric{
			{u.USG.WanRxPackets, prometheus.CounterValue, wan.RxPackets, labelWan},
			{u.USG.WanRxBytes, prometheus.CounterValue, wan.RxBytes, labelWan},
			{u.USG.WanRxDropped, prometheus.CounterValue, wan.RxDropped, labelWan},
			{u.USG.WanRxErrors, prometheus.CounterValue, wan.RxErrors, labelWan},
			{u.USG.WanTxPackets, prometheus.CounterValue, wan.TxPackets, labelWan},
			{u.USG.WanTxBytes, prometheus.CounterValue, wan.TxBytes, labelWan},
			{u.USG.WanRxBroadcast, prometheus.CounterValue, wan.RxBroadcast, labelWan},
			{u.USG.WanRxMulticast, prometheus.CounterValue, wan.RxMulticast, labelWan},
			{u.USG.WanSpeed, prometheus.CounterValue, wan.Speed.Val * 1000000, labelWan},
			{u.USG.WanTxBroadcast, prometheus.CounterValue, wan.TxBroadcast, labelWan},
			{u.USG.WanTxBytesR, prometheus.CounterValue, wan.TxBytesR, labelWan},
			{u.USG.WanTxDropped, prometheus.CounterValue, wan.TxDropped, labelWan},
			{u.USG.WanTxErrors, prometheus.CounterValue, wan.TxErrors, labelWan},
			{u.USG.WanTxMulticast, prometheus.CounterValue, wan.TxMulticast, labelWan},
			{u.USG.WanBytesR, prometheus.GaugeValue, wan.BytesR, labelWan},
		})
	}
}
