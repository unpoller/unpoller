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
	if !d.Adopted.Val || d.Locating.Val {
		return
	}
	labels := []string{d.Type, d.SiteName, d.Name}
	infoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID, d.Bytes.Txt, d.Uptime.Txt}
	// Gateway System Data.
	u.exportWANPorts(r, labels, d.Wan1, d.Wan2)
	u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)
	u.exportSYSstats(r, labels, d.SysStats, d.SystemStats)
	u.exportUSGstats(r, labels, d.Stat.Gw, d.SpeedtestStatus, d.Uplink)
	u.exportSTAcount(r, labels, d.UserNumSta, d.GuestNumSta, d.NumDesktop, d.UserNumSta, d.GuestNumSta)
	r.send([]*metric{
		{u.Device.Info, gauge, 1.0, append(labels, infoLabels...)},
		{u.Device.Uptime, gauge, d.Uptime, labels},
	})
}

// Gateway States
func (u *promUnifi) exportUSGstats(r report, labels []string, gw *unifi.Gw, st unifi.SpeedtestStatus, ul unifi.Uplink) {
	if gw == nil {
		return
	}
	labelLan := []string{"lan", labels[1], labels[2]}
	labelWan := []string{"all", labels[1], labels[2]}
	r.send([]*metric{
		{u.USG.LanRxPackets, counter, gw.LanRxPackets, labelLan},
		{u.USG.LanRxBytes, counter, gw.LanRxBytes, labelLan},
		{u.USG.LanTxPackets, counter, gw.LanTxPackets, labelLan},
		{u.USG.LanTxBytes, counter, gw.LanTxBytes, labelLan},
		{u.USG.LanRxDropped, counter, gw.LanRxDropped, labelLan},
		{u.USG.UplinkLatency, gauge, ul.Latency.Val / 1000, labelWan},
		{u.USG.UplinkSpeed, gauge, ul.Speed, labelWan},
		// Speed Test Stats
		{u.USG.Latency, gauge, st.Latency.Val / 1000, labelWan},
		{u.USG.Runtime, gauge, st.Runtime, labelWan},
		{u.USG.XputDownload, gauge, st.XputDownload, labelWan},
		{u.USG.XputUpload, gauge, st.XputUpload, labelWan},
	})
}

// WAN Stats
func (u *promUnifi) exportWANPorts(r report, labels []string, wans ...unifi.Wan) {
	for _, wan := range wans {
		if !wan.Up.Val {
			continue // only record UP interfaces.
		}
		labelWan := []string{wan.Name, labels[1], labels[2]}
		r.send([]*metric{
			{u.USG.WanRxPackets, counter, wan.RxPackets, labelWan},
			{u.USG.WanRxBytes, counter, wan.RxBytes, labelWan},
			{u.USG.WanRxDropped, counter, wan.RxDropped, labelWan},
			{u.USG.WanRxErrors, counter, wan.RxErrors, labelWan},
			{u.USG.WanTxPackets, counter, wan.TxPackets, labelWan},
			{u.USG.WanTxBytes, counter, wan.TxBytes, labelWan},
			{u.USG.WanRxBroadcast, counter, wan.RxBroadcast, labelWan},
			{u.USG.WanRxMulticast, counter, wan.RxMulticast, labelWan},
			{u.USG.WanSpeed, counter, wan.Speed.Val * 1000000, labelWan},
			{u.USG.WanTxBroadcast, counter, wan.TxBroadcast, labelWan},
			{u.USG.WanTxBytesR, counter, wan.TxBytesR, labelWan},
			{u.USG.WanTxDropped, counter, wan.TxDropped, labelWan},
			{u.USG.WanTxErrors, counter, wan.TxErrors, labelWan},
			{u.USG.WanTxMulticast, counter, wan.TxMulticast, labelWan},
			{u.USG.WanBytesR, gauge, wan.BytesR, labelWan},
		})
	}
}
