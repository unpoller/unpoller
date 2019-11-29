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
	// The first five labels for switch are shared with (the same as) switch ports.
	labels := []string{"ip", "type", "version", "site_name", "mac", "model", "name", "serial"}
	// Copy labels, and replace first four with different names.
	labelP := append([]string{"port_num", "port_name", "port_mac", "port_ip"}, labels[4:]...)
	return &usw{
		SwRxPackets:   prometheus.NewDesc(ns+"switch_receive_packets_total", "Switch Packets Received Total", labels, nil),
		SwRxBytes:     prometheus.NewDesc(ns+"switch_receive_bytes_total", "Switch Bytes Received Total", labels, nil),
		SwRxErrors:    prometheus.NewDesc(ns+"switch_receive_errors_total", "Switch Errors Received Total", labels, nil),
		SwRxDropped:   prometheus.NewDesc(ns+"switch_receive_dropped_total", "Switch Dropped Received Total", labels, nil),
		SwRxCrypts:    prometheus.NewDesc(ns+"switch_receive_crypts_total", "Switch Crypts Received Total", labels, nil),
		SwRxFrags:     prometheus.NewDesc(ns+"switch_receive_frags_total", "Switch Frags Received Total", labels, nil),
		SwTxPackets:   prometheus.NewDesc(ns+"switch_transmit_packets_total", "Switch Packets Transmit Total", labels, nil),
		SwTxBytes:     prometheus.NewDesc(ns+"switch_transmit_bytes_total", "Switch Bytes Transmit Total", labels, nil),
		SwTxErrors:    prometheus.NewDesc(ns+"switch_transmit_errors_total", "Switch Errors Transmit Total", labels, nil),
		SwTxDropped:   prometheus.NewDesc(ns+"switch_transmit_dropped_total", "Switch Dropped Transmit Total", labels, nil),
		SwTxRetries:   prometheus.NewDesc(ns+"switch_transmit_retries_total", "Switch Retries Transmit Total", labels, nil),
		SwRxMulticast: prometheus.NewDesc(ns+"switch_receive_multicast_total", "Switch Multicast Receive Total", labels, nil),
		SwRxBroadcast: prometheus.NewDesc(ns+"switch_receive_broadcast_total", "Switch Broadcast Receive Total", labels, nil),
		SwTxMulticast: prometheus.NewDesc(ns+"switch_transmit_multicast_total", "Switch Multicast Transmit Total", labels, nil),
		SwTxBroadcast: prometheus.NewDesc(ns+"switch_transmit_broadcast_total", "Switch Broadcast Transmit Total", labels, nil),
		SwBytes:       prometheus.NewDesc(ns+"switch_bytes_total", "Switch Bytes Transferred Total", labels, nil),
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
		Satisfaction: prometheus.NewDesc(pns+"satisfaction_percent", "Satisfaction", labelP, nil),
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

func (u *unifiCollector) exportUSWs(r report) {
	if r.metrics() == nil || r.metrics().Devices == nil || len(r.metrics().Devices.USWs) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, d := range r.metrics().Devices.USWs {
			u.exportUSW(r, d)
		}
	}()
}

func (u *unifiCollector) exportUSW(r report, d *unifi.USW) {
	labels := []string{d.IP, d.Type, d.Version, d.SiteName, d.Mac, d.Model, d.Name, d.Serial}
	if d.HasTemperature.Val {
		r.send([]*metricExports{{u.Device.Temperature, prometheus.GaugeValue, d.GeneralTemperature, labels}})
	}
	if d.HasFan.Val {
		r.send([]*metricExports{{u.Device.FanLevel, prometheus.GaugeValue, d.FanLevel, labels}})
	}

	// Switch data.
	r.send([]*metricExports{
		{u.Device.Uptime, prometheus.GaugeValue, d.Uptime, labels},
		{u.Device.TotalMaxPower, prometheus.GaugeValue, d.TotalMaxPower, labels},
		{u.Device.TotalTxBytes, prometheus.CounterValue, d.TxBytes, labels},
		{u.Device.TotalRxBytes, prometheus.CounterValue, d.RxBytes, labels},
		{u.Device.TotalBytes, prometheus.CounterValue, d.Bytes, labels},
		{u.Device.NumSta, prometheus.GaugeValue, d.NumSta, labels},
		{u.Device.UserNumSta, prometheus.GaugeValue, d.UserNumSta, labels},
		{u.Device.GuestNumSta, prometheus.GaugeValue, d.GuestNumSta, labels},
		{u.Device.Loadavg1, prometheus.GaugeValue, d.SysStats.Loadavg1, labels},
		{u.Device.Loadavg5, prometheus.GaugeValue, d.SysStats.Loadavg5, labels},
		{u.Device.Loadavg15, prometheus.GaugeValue, d.SysStats.Loadavg15, labels},
		{u.Device.MemUsed, prometheus.GaugeValue, d.SysStats.MemUsed, labels},
		{u.Device.MemTotal, prometheus.GaugeValue, d.SysStats.MemTotal, labels},
		{u.Device.MemBuffer, prometheus.GaugeValue, d.SysStats.MemBuffer, labels},
		{u.Device.CPU, prometheus.GaugeValue, d.SystemStats.CPU, labels},
		{u.Device.Mem, prometheus.GaugeValue, d.SystemStats.Mem, labels},
	})
	u.exportPortTable(r, labels, d.PortTable)
	u.exportUSWstats(r, labels, d.Stat.Sw)
}

func (u *unifiCollector) exportUSWstats(r report, labels []string, sw *unifi.Sw) {
	r.send([]*metricExports{
		{u.USW.SwRxPackets, prometheus.CounterValue, sw.RxPackets, labels},
		{u.USW.SwRxBytes, prometheus.CounterValue, sw.RxBytes, labels},
		{u.USW.SwRxErrors, prometheus.CounterValue, sw.RxErrors, labels},
		{u.USW.SwRxDropped, prometheus.CounterValue, sw.RxDropped, labels},
		{u.USW.SwRxCrypts, prometheus.CounterValue, sw.RxCrypts, labels},
		{u.USW.SwRxFrags, prometheus.CounterValue, sw.RxFrags, labels},
		{u.USW.SwTxPackets, prometheus.CounterValue, sw.TxPackets, labels},
		{u.USW.SwTxBytes, prometheus.CounterValue, sw.TxBytes, labels},
		{u.USW.SwTxErrors, prometheus.CounterValue, sw.TxErrors, labels},
		{u.USW.SwTxDropped, prometheus.CounterValue, sw.TxDropped, labels},
		{u.USW.SwTxRetries, prometheus.CounterValue, sw.TxRetries, labels},
		{u.USW.SwRxMulticast, prometheus.CounterValue, sw.RxMulticast, labels},
		{u.USW.SwRxBroadcast, prometheus.CounterValue, sw.RxBroadcast, labels},
		{u.USW.SwTxMulticast, prometheus.CounterValue, sw.TxMulticast, labels},
		{u.USW.SwTxBroadcast, prometheus.CounterValue, sw.TxBroadcast, labels},
		{u.USW.SwBytes, prometheus.CounterValue, sw.Bytes, labels},
	})
}

func (u *unifiCollector) exportPortTable(r report, labels []string, pt []unifi.Port) {
	// Per-port data on a switch
	for _, p := range pt {
		if !p.Up.Val {
			continue
		}
		// Copy labels, and add four new ones.
		l := append([]string{p.PortIdx.Txt, p.Name, p.Mac, p.IP}, labels[4:]...)
		if p.PoeEnable.Val && p.PortPoe.Val {
			r.send([]*metricExports{
				{u.USW.PoeCurrent, prometheus.GaugeValue, p.PoeCurrent, l},
				{u.USW.PoePower, prometheus.GaugeValue, p.PoePower, l},
				{u.USW.PoeVoltage, prometheus.GaugeValue, p.PoeVoltage, l},
			})
		}

		r.send([]*metricExports{
			{u.USW.RxBroadcast, prometheus.CounterValue, p.RxBroadcast, l},
			{u.USW.RxBytes, prometheus.CounterValue, p.RxBytes, l},
			{u.USW.RxBytesR, prometheus.GaugeValue, p.RxBytesR, l},
			{u.USW.RxDropped, prometheus.CounterValue, p.RxDropped, l},
			{u.USW.RxErrors, prometheus.CounterValue, p.RxErrors, l},
			{u.USW.RxMulticast, prometheus.CounterValue, p.RxMulticast, l},
			{u.USW.RxPackets, prometheus.CounterValue, p.RxPackets, l},
			{u.USW.Satisfaction, prometheus.GaugeValue, p.Satisfaction, l},
			{u.USW.Speed, prometheus.GaugeValue, p.Speed.Val * 1000000, l},
			{u.USW.TxBroadcast, prometheus.CounterValue, p.TxBroadcast, l},
			{u.USW.TxBytes, prometheus.CounterValue, p.TxBytes, l},
			{u.USW.TxBytesR, prometheus.GaugeValue, p.TxBytesR, l},
			{u.USW.TxDropped, prometheus.CounterValue, p.TxDropped, l},
			{u.USW.TxErrors, prometheus.CounterValue, p.TxErrors, l},
			{u.USW.TxMulticast, prometheus.CounterValue, p.TxMulticast, l},
		})
	}
}
