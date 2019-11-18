package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type usw struct {
	Uptime        *prometheus.Desc
	Temperature   *prometheus.Desc
	TotalMaxPower *prometheus.Desc
	FanLevel      *prometheus.Desc
	TotalTxBytes  *prometheus.Desc
	TotalRxBytes  *prometheus.Desc
	TotalBytes    *prometheus.Desc
	NumSta        *prometheus.Desc
	UserNumSta    *prometheus.Desc
	GuestNumSta   *prometheus.Desc
	// System Stats
	Loadavg1  *prometheus.Desc
	Loadavg5  *prometheus.Desc
	Loadavg15 *prometheus.Desc
	MemBuffer *prometheus.Desc
	MemTotal  *prometheus.Desc
	MemUsed   *prometheus.Desc
	CPU       *prometheus.Desc
	Mem       *prometheus.Desc
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
	if ns += "_usw_"; ns == "_usw_" {
		ns = "usw_"
	}
	pns := ns + "port_"
	// The first five labels for switch are shared with (the same as) switch ports.
	labels := []string{"type", "version", "device_id", "ip",
		"site_name", "mac", "model", "name", "serial", "site_id"}
	// Copy labels, and replace first four with different names.
	labelP := append([]string{"port_num", "port_name", "port_mac", "port_ip"}, labels[5:]...)

	return &usw{
		// switch data
		Uptime:        prometheus.NewDesc(ns+"uptime", "Uptime", labels, nil),
		Temperature:   prometheus.NewDesc(ns+"temperature", "Temperature", labels, nil),
		TotalMaxPower: prometheus.NewDesc(ns+"max_power_total", "Total Max Power", labels, nil),
		FanLevel:      prometheus.NewDesc(ns+"fan_level", "Fan Level", labels, nil),
		TotalTxBytes:  prometheus.NewDesc(ns+"bytes_tx_total", "Total Transmitted Bytes", labels, nil),
		TotalRxBytes:  prometheus.NewDesc(ns+"bytes_rx_total", "Total Received Bytes", labels, nil),
		TotalBytes:    prometheus.NewDesc(ns+"bytes_total", "Total Bytes Transfered", labels, nil),
		NumSta:        prometheus.NewDesc(ns+"stations_total", "Number of Stations", labels, nil),
		UserNumSta:    prometheus.NewDesc(ns+"stations_user_total", "Number of User Stations", labels, nil),
		GuestNumSta:   prometheus.NewDesc(ns+"stations_guest_total", "Number of Guest Stations", labels, nil),
		Loadavg1:      prometheus.NewDesc(ns+"load_average_1", "System Load Average 1 Minute", labels, nil),
		Loadavg5:      prometheus.NewDesc(ns+"load_average_5", "System Load Average 5 Minutes", labels, nil),
		Loadavg15:     prometheus.NewDesc(ns+"load_average_15", "System Load Average 15 Minutes", labels, nil),
		MemUsed:       prometheus.NewDesc(ns+"memory_utilization", "System Memory Used", labels, nil),
		MemTotal:      prometheus.NewDesc(ns+"memory_installed", "System Installed Memory", labels, nil),
		MemBuffer:     prometheus.NewDesc(ns+"memory_buffer", "System Memory Buffer", labels, nil),
		CPU:           prometheus.NewDesc(ns+"cpu_utilization", "System CPU % Utilized", labels, nil),
		Mem:           prometheus.NewDesc(ns+"memory_utilization", "System Memory % Utilized", labels, nil), // this may not be %.

		SwRxPackets:   prometheus.NewDesc(ns+"switch_packets_rx_total", "Switch Packets Received Total", labels, nil),
		SwRxBytes:     prometheus.NewDesc(ns+"switch_bytes_rx_total", "Switch Bytes Received Total", labels, nil),
		SwRxErrors:    prometheus.NewDesc(ns+"switch_errors_rx_total", "Switch Errors Received Total", labels, nil),
		SwRxDropped:   prometheus.NewDesc(ns+"switch_dropped_rx_total", "Switch Dropped Received Total", labels, nil),
		SwRxCrypts:    prometheus.NewDesc(ns+"switch_crypts_rx_total", "Switch Crypts Received Total", labels, nil),
		SwRxFrags:     prometheus.NewDesc(ns+"switch_frags_rx_total", "Switch Frags Received Total", labels, nil),
		SwTxPackets:   prometheus.NewDesc(ns+"switch_packets_tx_total", "Switch Packets Transmit Total", labels, nil),
		SwTxBytes:     prometheus.NewDesc(ns+"switch_bytes_tx_total", "Switch Bytes Transmit Total", labels, nil),
		SwTxErrors:    prometheus.NewDesc(ns+"switch_errors_tx_total", "Switch Errors Transmit Total", labels, nil),
		SwTxDropped:   prometheus.NewDesc(ns+"switch_dropped_tx_total", "Switch Dropped Transmit Total", labels, nil),
		SwTxRetries:   prometheus.NewDesc(ns+"switch_retries_tx_total", "Switch Retries Transmit Total", labels, nil),
		SwRxMulticast: prometheus.NewDesc(ns+"switch_multicast_rx_total", "Switch Multicast Receive Total", labels, nil),
		SwRxBroadcast: prometheus.NewDesc(ns+"switch_broadcast_rx_total", "Switch Broadcast Receive Total", labels, nil),
		SwTxMulticast: prometheus.NewDesc(ns+"switch_multicast_tx_total", "Switch Multicast Transmit Total", labels, nil),
		SwTxBroadcast: prometheus.NewDesc(ns+"switch_broadcast_tx_total", "Switch Broadcast Transmit Total", labels, nil),
		SwBytes:       prometheus.NewDesc(ns+"switch_bytes_total", "Switch Bytes Transfered Total", labels, nil),

		// per-port data
		PoeCurrent:   prometheus.NewDesc(pns+"poe_current", "POE Current", labelP, nil),
		PoePower:     prometheus.NewDesc(pns+"poe_power", "POE Power", labelP, nil),
		PoeVoltage:   prometheus.NewDesc(pns+"poe_voltage", "POE Voltage", labelP, nil),
		RxBroadcast:  prometheus.NewDesc(pns+"broadcast_rx_total", "Receive Broadcast", labelP, nil),
		RxBytes:      prometheus.NewDesc(pns+"bytes_rx_total", "Total Receive Bytes", labelP, nil),
		RxBytesR:     prometheus.NewDesc(pns+"bytes_rx_rate", "Receive Bytes Rate", labelP, nil),
		RxDropped:    prometheus.NewDesc(pns+"dropped_rx_total", "Total Receive Dropped", labelP, nil),
		RxErrors:     prometheus.NewDesc(pns+"errors_rx_total", "Total Receive Errors", labelP, nil),
		RxMulticast:  prometheus.NewDesc(pns+"multicast_rx_total", "Total Receive Multicast", labelP, nil),
		RxPackets:    prometheus.NewDesc(pns+"packets_rx_total", "Total Receive Packets", labelP, nil),
		Satisfaction: prometheus.NewDesc(pns+"satisfaction", "Satisfaction", labelP, nil),
		Speed:        prometheus.NewDesc(pns+"speed", "Speed", labelP, nil),
		TxBroadcast:  prometheus.NewDesc(pns+"broadcast_tx_total", "Total Transmit Broadcast", labelP, nil),
		TxBytes:      prometheus.NewDesc(pns+"bytes_tx_total", "Total Transmit Bytes", labelP, nil),
		TxBytesR:     prometheus.NewDesc(pns+"bytes_tx_rate", "Transmit Bytes Rate", labelP, nil),
		TxDropped:    prometheus.NewDesc(pns+"dropped_tx_total", "Total Transmit Dropped", labelP, nil),
		TxErrors:     prometheus.NewDesc(pns+"errors_tx_total", "Total Transmit Errors", labelP, nil),
		TxMulticast:  prometheus.NewDesc(pns+"multicast_tx_total", "Total Tranmist Multicast", labelP, nil),
		TxPackets:    prometheus.NewDesc(pns+"packets_tx_total", "Total Transmit Packets", labelP, nil),
	}
}

// exportUSW exports Network Switch Data
func (u *unifiCollector) exportUSW(s *unifi.USW) []*metricExports {
	labels := []string{s.Type, s.Version, s.DeviceID, s.IP,
		s.SiteName, s.Mac, s.Model, s.Name, s.Serial, s.SiteID}

	// Switch data.
	return append([]*metricExports{
		{u.USW.Uptime, prometheus.GaugeValue, s.Uptime, labels},
		{u.USW.Temperature, prometheus.GaugeValue, s.GeneralTemperature, labels},
		{u.USW.TotalMaxPower, prometheus.GaugeValue, s.TotalMaxPower, labels},
		{u.USW.FanLevel, prometheus.GaugeValue, s.FanLevel, labels},
		{u.USW.TotalTxBytes, prometheus.CounterValue, s.TxBytes, labels},
		{u.USW.TotalRxBytes, prometheus.CounterValue, s.RxBytes, labels},
		{u.USW.TotalBytes, prometheus.CounterValue, s.Bytes, labels},
		{u.USW.NumSta, prometheus.GaugeValue, s.NumSta, labels},
		{u.USW.UserNumSta, prometheus.GaugeValue, s.UserNumSta, labels},
		{u.USW.GuestNumSta, prometheus.GaugeValue, s.GuestNumSta, labels},
		{u.USW.Loadavg1, prometheus.GaugeValue, s.SysStats.Loadavg1, labels},
		{u.USW.Loadavg5, prometheus.GaugeValue, s.SysStats.Loadavg5, labels},
		{u.USW.Loadavg15, prometheus.GaugeValue, s.SysStats.Loadavg15, labels},
		{u.USW.MemUsed, prometheus.GaugeValue, s.SysStats.MemUsed, labels},
		{u.USW.MemTotal, prometheus.GaugeValue, s.SysStats.MemTotal, labels},
		{u.USW.MemBuffer, prometheus.GaugeValue, s.SysStats.MemBuffer, labels},
		{u.USW.CPU, prometheus.GaugeValue, s.SystemStats.CPU, labels},
		{u.USW.Mem, prometheus.GaugeValue, s.SystemStats.Mem, labels},
		{u.USW.SwRxPackets, prometheus.CounterValue, s.Stat.Sw.RxPackets, labels},
		{u.USW.SwRxBytes, prometheus.CounterValue, s.Stat.Sw.RxBytes, labels},
		{u.USW.SwRxErrors, prometheus.CounterValue, s.Stat.Sw.RxErrors, labels},
		{u.USW.SwRxDropped, prometheus.CounterValue, s.Stat.Sw.RxDropped, labels},
		{u.USW.SwRxCrypts, prometheus.CounterValue, s.Stat.Sw.RxCrypts, labels},
		{u.USW.SwRxFrags, prometheus.CounterValue, s.Stat.Sw.RxFrags, labels},
		{u.USW.SwTxPackets, prometheus.CounterValue, s.Stat.Sw.TxPackets, labels},
		{u.USW.SwTxBytes, prometheus.CounterValue, s.Stat.Sw.TxBytes, labels},
		{u.USW.SwTxErrors, prometheus.CounterValue, s.Stat.Sw.TxErrors, labels},
		{u.USW.SwTxDropped, prometheus.CounterValue, s.Stat.Sw.TxDropped, labels},
		{u.USW.SwTxRetries, prometheus.CounterValue, s.Stat.Sw.TxRetries, labels},
		{u.USW.SwRxMulticast, prometheus.CounterValue, s.Stat.Sw.RxMulticast, labels},
		{u.USW.SwRxBroadcast, prometheus.CounterValue, s.Stat.Sw.RxBroadcast, labels},
		{u.USW.SwTxMulticast, prometheus.CounterValue, s.Stat.Sw.TxMulticast, labels},
		{u.USW.SwTxBroadcast, prometheus.CounterValue, s.Stat.Sw.TxBroadcast, labels},
		{u.USW.SwBytes, prometheus.CounterValue, s.Stat.Sw.Bytes, labels},
	}, u.exportPortTable(s.PortTable, labels[5:])...)
}

func (u *unifiCollector) exportPortTable(pt []unifi.Port, labels []string) []*metricExports {
	var metrics []*metricExports
	// Per-port data on a switch
	for _, p := range pt {
		// Copy labels, and add four new ones.
		l := append([]string{p.PortIdx.Txt, p.Name, p.Mac, p.IP}, labels...)

		metrics = append(metrics, []*metricExports{
			{u.USW.PoeCurrent, prometheus.GaugeValue, p.PoeCurrent, l},
			{u.USW.PoePower, prometheus.GaugeValue, p.PoePower, l},
			{u.USW.PoeVoltage, prometheus.GaugeValue, p.PoeVoltage, l},
			{u.USW.RxBroadcast, prometheus.CounterValue, p.RxBroadcast, l},
			{u.USW.RxBytes, prometheus.CounterValue, p.RxBytes, l},
			{u.USW.RxBytesR, prometheus.GaugeValue, p.RxBytesR, l},
			{u.USW.RxDropped, prometheus.CounterValue, p.RxDropped, l},
			{u.USW.RxErrors, prometheus.CounterValue, p.RxErrors, l},
			{u.USW.RxMulticast, prometheus.CounterValue, p.RxMulticast, l},
			{u.USW.RxPackets, prometheus.CounterValue, p.RxPackets, l},
			{u.USW.Satisfaction, prometheus.GaugeValue, p.Satisfaction, l},
			{u.USW.Speed, prometheus.GaugeValue, p.Speed, l},
			{u.USW.TxBroadcast, prometheus.CounterValue, p.TxBroadcast, l},
			{u.USW.TxBytes, prometheus.CounterValue, p.TxBytes, l},
			{u.USW.TxBytesR, prometheus.GaugeValue, p.TxBytesR, l},
			{u.USW.TxDropped, prometheus.CounterValue, p.TxDropped, l},
			{u.USW.TxErrors, prometheus.CounterValue, p.TxErrors, l},
			{u.USW.TxMulticast, prometheus.CounterValue, p.TxMulticast, l},
		}...)
	}

	return metrics
}
