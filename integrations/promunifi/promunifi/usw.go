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
	labels := []string{"site_name", "mac", "model", "name", "serial", "site_id",
		"type", "version", "device_id", "oid"}
	// Copy labels, and replace last four with different names.
	labelP := append(append([]string{}, labels[:6]...),
		"port_num", "port_name", "port_mac", "port_ip")

	return &usw{
		// switch data
		Uptime:        prometheus.NewDesc(ns+"uptime", "Uptime", labels, nil),
		Temperature:   prometheus.NewDesc(ns+"temperature", "Temperature", labels, nil),
		TotalMaxPower: prometheus.NewDesc(ns+"max_power_total", "Total Max Power", labels, nil),
		FanLevel:      prometheus.NewDesc(ns+"fan_level", "Fan Level", labels, nil),
		TotalTxBytes:  prometheus.NewDesc(ns+"tx_bytes_total", "Total Transmitted Bytes", labels, nil),
		TotalRxBytes:  prometheus.NewDesc(ns+"rx_bytes_total", "Total Received Bytes", labels, nil),
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
		Mem:           prometheus.NewDesc(ns+"memory", "System Memory % Utilized", labels, nil), // this may not be %.
		// per-port data
		PoeCurrent:   prometheus.NewDesc(pns+"poe_current", "POE Current", labelP, nil),
		PoePower:     prometheus.NewDesc(pns+"poe_power", "POE Power", labelP, nil),
		PoeVoltage:   prometheus.NewDesc(pns+"poe_voltage", "POE Voltage", labelP, nil),
		RxBroadcast:  prometheus.NewDesc(pns+"rx_broadcast_total", "Receive Broadcast", labelP, nil),
		RxBytes:      prometheus.NewDesc(pns+"rx_bytes_total", "Total Receive Bytes", labelP, nil),
		RxBytesR:     prometheus.NewDesc(pns+"rx_bytes_rate", "Receive Bytes Rate", labelP, nil),
		RxDropped:    prometheus.NewDesc(pns+"rx_dropped_total", "Total Receive Dropped", labelP, nil),
		RxErrors:     prometheus.NewDesc(pns+"rx_errors_total", "Total Receive Errors", labelP, nil),
		RxMulticast:  prometheus.NewDesc(pns+"rx_multicast_total", "Total Receive Multicast", labelP, nil),
		RxPackets:    prometheus.NewDesc(pns+"rx_packets_total", "Total Receive Packets", labelP, nil),
		Satisfaction: prometheus.NewDesc(pns+"satisfaction", "Satisfaction", labelP, nil),
		Speed:        prometheus.NewDesc(pns+"speed", "Speed", labelP, nil),
		TxBroadcast:  prometheus.NewDesc(pns+"tx_broadcast_total", "Total Transmit Broadcast", labelP, nil),
		TxBytes:      prometheus.NewDesc(pns+"tx_bytes_total", "Total Transmit Bytes", labelP, nil),
		TxBytesR:     prometheus.NewDesc(pns+"tx_bytes_rate", "Transmit Bytes Rate", labelP, nil),
		TxDropped:    prometheus.NewDesc(pns+"tx_dropped_total", "Total Transmit Dropped", labelP, nil),
		TxErrors:     prometheus.NewDesc(pns+"tx_errors_total", "Total Transmit Errors", labelP, nil),
		TxMulticast:  prometheus.NewDesc(pns+"tx_multicast_total", "Total Tranmist Multicast", labelP, nil),
		TxPackets:    prometheus.NewDesc(pns+"tx_packets_total", "Total Transmit Packets", labelP, nil),
	}
}

// exportUSW exports Network Switch Data
func (u *unifiCollector) exportUSW(s *unifi.USW) []*metricExports {
	labels := []string{s.SiteName, s.Mac, s.Model, s.Name, s.Serial, s.SiteID,
		s.Type, s.Version, s.DeviceID, s.Stat.Oid}

	// Switch data.
	m := []*metricExports{
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
	}

	// Per-port data on the switch
	for _, p := range s.PortTable {
		// Copy labels, and replace last four with different data.
		l := append(append([]string{}, labels[:6]...), p.PortIdx.Txt, p.Name, p.Mac, p.IP)
		m = append(m, []*metricExports{
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
	return m
}
