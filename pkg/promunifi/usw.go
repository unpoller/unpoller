package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type usw struct {
	// Switch "total" traffic stats
	SwRxPackets   *prometheus.Desc
	SwRxBytes     *prometheus.Desc
	SwRxErrors    *prometheus.Desc
	SwRxDropped   *prometheus.Desc
	SwRxCrypts    *prometheus.Desc
	SwRxFrags     *prometheus.Desc
	SwTxPackets   *prometheus.Desc
	SwTxBytes     *prometheus.Desc
	SwTxErrors    *prometheus.Desc
	SwTxDropped   *prometheus.Desc
	SwTxRetries   *prometheus.Desc
	SwRxMulticast *prometheus.Desc
	SwRxBroadcast *prometheus.Desc
	SwTxMulticast *prometheus.Desc
	SwTxBroadcast *prometheus.Desc
	SwBytes       *prometheus.Desc
	// Port data.
	PoeCurrent   *prometheus.Desc
	PoePower     *prometheus.Desc
	PoeVoltage   *prometheus.Desc
	RxBroadcast  *prometheus.Desc
	RxBytes      *prometheus.Desc
	RxBytesR     *prometheus.Desc
	RxDropped    *prometheus.Desc
	RxErrors     *prometheus.Desc
	RxMulticast  *prometheus.Desc
	RxPackets    *prometheus.Desc
	Satisfaction *prometheus.Desc
	Speed        *prometheus.Desc
	TxBroadcast  *prometheus.Desc
	TxBytes      *prometheus.Desc
	TxBytesR     *prometheus.Desc
	TxDropped    *prometheus.Desc
	TxErrors     *prometheus.Desc
	TxMulticast  *prometheus.Desc
	TxPackets    *prometheus.Desc
}

func descUSW(ns string) *usw {
	pns := ns + "port_"
	labelS := []string{"site_name", "name"}
	labelP := []string{"port_id", "port_num", "port_name", "port_mac", "port_ip", "site_name", "name"}
	return &usw{
		// This data may be derivable by sum()ing the port data.
		SwRxPackets:   prometheus.NewDesc(ns+"switch_receive_packets_total", "Switch Packets Received Total", labelS, nil),
		SwRxBytes:     prometheus.NewDesc(ns+"switch_receive_bytes_total", "Switch Bytes Received Total", labelS, nil),
		SwRxErrors:    prometheus.NewDesc(ns+"switch_receive_errors_total", "Switch Errors Received Total", labelS, nil),
		SwRxDropped:   prometheus.NewDesc(ns+"switch_receive_dropped_total", "Switch Dropped Received Total", labelS, nil),
		SwRxCrypts:    prometheus.NewDesc(ns+"switch_receive_crypts_total", "Switch Crypts Received Total", labelS, nil),
		SwRxFrags:     prometheus.NewDesc(ns+"switch_receive_frags_total", "Switch Frags Received Total", labelS, nil),
		SwTxPackets:   prometheus.NewDesc(ns+"switch_transmit_packets_total", "Switch Packets Transmit Total", labelS, nil),
		SwTxBytes:     prometheus.NewDesc(ns+"switch_transmit_bytes_total", "Switch Bytes Transmit Total", labelS, nil),
		SwTxErrors:    prometheus.NewDesc(ns+"switch_transmit_errors_total", "Switch Errors Transmit Total", labelS, nil),
		SwTxDropped:   prometheus.NewDesc(ns+"switch_transmit_dropped_total", "Switch Dropped Transmit Total", labelS, nil),
		SwTxRetries:   prometheus.NewDesc(ns+"switch_transmit_retries_total", "Switch Retries Transmit Total", labelS, nil),
		SwRxMulticast: prometheus.NewDesc(ns+"switch_receive_multicast_total", "Switch Multicast Receive Total", labelS, nil),
		SwRxBroadcast: prometheus.NewDesc(ns+"switch_receive_broadcast_total", "Switch Broadcast Receive Total", labelS, nil),
		SwTxMulticast: prometheus.NewDesc(ns+"switch_transmit_multicast_total", "Switch Multicast Transmit Total", labelS, nil),
		SwTxBroadcast: prometheus.NewDesc(ns+"switch_transmit_broadcast_total", "Switch Broadcast Transmit Total", labelS, nil),
		SwBytes:       prometheus.NewDesc(ns+"switch_bytes_total", "Switch Bytes Transferred Total", labelS, nil),
		// per-port data
		PoeCurrent:   prometheus.NewDesc(pns+"poe_amperes", "POE Current", labelP, nil),
		PoePower:     prometheus.NewDesc(pns+"poe_watts", "POE Power", labelP, nil),
		PoeVoltage:   prometheus.NewDesc(pns+"poe_volts", "POE Voltage", labelP, nil),
		RxBroadcast:  prometheus.NewDesc(pns+"receive_broadcast_total", "Receive Broadcast", labelP, nil),
		RxBytes:      prometheus.NewDesc(pns+"receive_bytes_total", "Total Receive Bytes", labelP, nil),
		RxBytesR:     prometheus.NewDesc(pns+"receive_rate_bytes", "Receive Bytes Rate", labelP, nil),
		RxDropped:    prometheus.NewDesc(pns+"receive_dropped_total", "Total Receive Dropped", labelP, nil),
		RxErrors:     prometheus.NewDesc(pns+"receive_errors_total", "Total Receive Errors", labelP, nil),
		RxMulticast:  prometheus.NewDesc(pns+"receive_multicast_total", "Total Receive Multicast", labelP, nil),
		RxPackets:    prometheus.NewDesc(pns+"receive_packets_total", "Total Receive Packets", labelP, nil),
		Satisfaction: prometheus.NewDesc(pns+"satisfaction_ratio", "Satisfaction", labelP, nil),
		Speed:        prometheus.NewDesc(pns+"port_speed_bps", "Speed", labelP, nil),
		TxBroadcast:  prometheus.NewDesc(pns+"transmit_broadcast_total", "Total Transmit Broadcast", labelP, nil),
		TxBytes:      prometheus.NewDesc(pns+"transmit_bytes_total", "Total Transmit Bytes", labelP, nil),
		TxBytesR:     prometheus.NewDesc(pns+"transmit_rate_bytes", "Transmit Bytes Rate", labelP, nil),
		TxDropped:    prometheus.NewDesc(pns+"transmit_dropped_total", "Total Transmit Dropped", labelP, nil),
		TxErrors:     prometheus.NewDesc(pns+"transmit_errors_total", "Total Transmit Errors", labelP, nil),
		TxMulticast:  prometheus.NewDesc(pns+"transmit_multicast_total", "Total Tranmist Multicast", labelP, nil),
		TxPackets:    prometheus.NewDesc(pns+"transmit_packets_total", "Total Transmit Packets", labelP, nil),
	}
}

func (u *promUnifi) exportUSW(r report, d *unifi.USW) {
	if !d.Adopted.Val || d.Locating.Val {
		return
	}
	labels := []string{d.Type, d.SiteName, d.Name}
	infoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID, d.Bytes.Txt, d.Uptime.Txt}
	u.exportUSWstats(r, labels, d.Stat.Sw)
	u.exportPRTtable(r, labels, d.PortTable)
	u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)
	u.exportSYSstats(r, labels, d.SysStats, d.SystemStats)
	u.exportSTAcount(r, labels, d.UserNumSta, d.GuestNumSta)
	r.send([]*metric{
		{u.Device.Info, gauge, 1.0, append(labels, infoLabels...)},
		{u.Device.Uptime, gauge, d.Uptime, labels},
	})
	// Switch System Data.
	if d.HasTemperature.Val {
		r.send([]*metric{{u.Device.Temperature, gauge, d.GeneralTemperature, labels}})
	}
	if d.HasFan.Val {
		r.send([]*metric{{u.Device.FanLevel, gauge, d.FanLevel, labels}})
	}
	if d.TotalMaxPower.Txt != "" {
		r.send([]*metric{{u.Device.TotalMaxPower, gauge, d.TotalMaxPower, labels}})
	}
}

// Switch Stats
func (u *promUnifi) exportUSWstats(r report, labels []string, sw *unifi.Sw) {
	if sw == nil {
		return
	}
	labelS := labels[1:]
	r.send([]*metric{
		{u.USW.SwRxPackets, counter, sw.RxPackets, labelS},
		{u.USW.SwRxBytes, counter, sw.RxBytes, labelS},
		{u.USW.SwRxErrors, counter, sw.RxErrors, labelS},
		{u.USW.SwRxDropped, counter, sw.RxDropped, labelS},
		{u.USW.SwRxCrypts, counter, sw.RxCrypts, labelS},
		{u.USW.SwRxFrags, counter, sw.RxFrags, labelS},
		{u.USW.SwTxPackets, counter, sw.TxPackets, labelS},
		{u.USW.SwTxBytes, counter, sw.TxBytes, labelS},
		{u.USW.SwTxErrors, counter, sw.TxErrors, labelS},
		{u.USW.SwTxDropped, counter, sw.TxDropped, labelS},
		{u.USW.SwTxRetries, counter, sw.TxRetries, labelS},
		{u.USW.SwRxMulticast, counter, sw.RxMulticast, labelS},
		{u.USW.SwRxBroadcast, counter, sw.RxBroadcast, labelS},
		{u.USW.SwTxMulticast, counter, sw.TxMulticast, labelS},
		{u.USW.SwTxBroadcast, counter, sw.TxBroadcast, labelS},
		{u.USW.SwBytes, counter, sw.Bytes, labelS},
	})
}

// Switch Port Table
func (u *promUnifi) exportPRTtable(r report, labels []string, pt []unifi.Port) {
	// Per-port data on a switch
	for _, p := range pt {
		if !p.Up.Val || !p.Enable.Val {
			continue
		}
		// Copy labels, and add four new ones.
		labelP := []string{labels[2] + " Port " + p.PortIdx.Txt, p.PortIdx.Txt, p.Name, p.Mac, p.IP, labels[1], labels[2]}
		if p.PoeEnable.Val && p.PortPoe.Val {
			r.send([]*metric{
				{u.USW.PoeCurrent, gauge, p.PoeCurrent, labelP},
				{u.USW.PoePower, gauge, p.PoePower, labelP},
				{u.USW.PoeVoltage, gauge, p.PoeVoltage, labelP},
			})
		}

		r.send([]*metric{
			{u.USW.RxBroadcast, counter, p.RxBroadcast, labelP},
			{u.USW.RxBytes, counter, p.RxBytes, labelP},
			{u.USW.RxBytesR, gauge, p.RxBytesR, labelP},
			{u.USW.RxDropped, counter, p.RxDropped, labelP},
			{u.USW.RxErrors, counter, p.RxErrors, labelP},
			{u.USW.RxMulticast, counter, p.RxMulticast, labelP},
			{u.USW.RxPackets, counter, p.RxPackets, labelP},
			{u.USW.Satisfaction, gauge, p.Satisfaction.Val / 100.0, labelP},
			{u.USW.Speed, gauge, p.Speed.Val * 1000000, labelP},
			{u.USW.TxBroadcast, counter, p.TxBroadcast, labelP},
			{u.USW.TxBytes, counter, p.TxBytes, labelP},
			{u.USW.TxBytesR, gauge, p.TxBytesR, labelP},
			{u.USW.TxDropped, counter, p.TxDropped, labelP},
			{u.USW.TxErrors, counter, p.TxErrors, labelP},
			{u.USW.TxMulticast, counter, p.TxMulticast, labelP},
			{u.USW.TxPackets, counter, p.TxPackets, labelP},
		})
	}
}
