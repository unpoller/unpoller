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
	labels := []string{"type", "version", "ip", "site_name", "mac", "model", "name", "serial"}
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
		TotalBytes:    prometheus.NewDesc(ns+"bytes_total", "Total Bytes Transferred", labels, nil),
		NumSta:        prometheus.NewDesc(ns+"num_stations", "Number of Stations", labels, nil),
		UserNumSta:    prometheus.NewDesc(ns+"num_stations_user", "Number of User Stations", labels, nil),
		GuestNumSta:   prometheus.NewDesc(ns+"num_stations_guest", "Number of Guest Stations", labels, nil),
		Loadavg1:      prometheus.NewDesc(ns+"load_average_1", "System Load Average 1 Minute", labels, nil),
		Loadavg5:      prometheus.NewDesc(ns+"load_average_5", "System Load Average 5 Minutes", labels, nil),
		Loadavg15:     prometheus.NewDesc(ns+"load_average_15", "System Load Average 15 Minutes", labels, nil),
		MemUsed:       prometheus.NewDesc(ns+"memory_used_bytes", "System Memory Used", labels, nil),
		MemTotal:      prometheus.NewDesc(ns+"memory_installed_bytes", "System Installed Memory", labels, nil),
		MemBuffer:     prometheus.NewDesc(ns+"memory_buffer_bytes", "System Memory Buffer", labels, nil),
		CPU:           prometheus.NewDesc(ns+"cpu_utilization_percent", "System CPU % Utilized", labels, nil),
		Mem:           prometheus.NewDesc(ns+"memory_utilization_percent", "System Memory % Utilized", labels, nil),

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
		PoeCurrent:   prometheus.NewDesc(pns+"poe_current", "POE Current", labelP, nil),
		PoePower:     prometheus.NewDesc(pns+"poe_power", "POE Power", labelP, nil),
		PoeVoltage:   prometheus.NewDesc(pns+"poe_voltage", "POE Voltage", labelP, nil),
		RxBroadcast:  prometheus.NewDesc(pns+"receive_broadcast_total", "Receive Broadcast", labelP, nil),
		RxBytes:      prometheus.NewDesc(pns+"receive_bytes_total", "Total Receive Bytes", labelP, nil),
		RxBytesR:     prometheus.NewDesc(pns+"receive_rate_bytes", "Receive Bytes Rate", labelP, nil),
		RxDropped:    prometheus.NewDesc(pns+"receive_dropped_total", "Total Receive Dropped", labelP, nil),
		RxErrors:     prometheus.NewDesc(pns+"receive_errors_total", "Total Receive Errors", labelP, nil),
		RxMulticast:  prometheus.NewDesc(pns+"receive_multicast_total", "Total Receive Multicast", labelP, nil),
		RxPackets:    prometheus.NewDesc(pns+"receive_packets_total", "Total Receive Packets", labelP, nil),
		Satisfaction: prometheus.NewDesc(pns+"satisfaction_percent", "Satisfaction", labelP, nil),
		Speed:        prometheus.NewDesc(pns+"port_speed_mbps", "Speed", labelP, nil),
		TxBroadcast:  prometheus.NewDesc(pns+"transmit_broadcast_total", "Total Transmit Broadcast", labelP, nil),
		TxBytes:      prometheus.NewDesc(pns+"transmit_bytes_total", "Total Transmit Bytes", labelP, nil),
		TxBytesR:     prometheus.NewDesc(pns+"transmit_rate_bytes", "Transmit Bytes Rate", labelP, nil),
		TxDropped:    prometheus.NewDesc(pns+"transmit_dropped_total", "Total Transmit Dropped", labelP, nil),
		TxErrors:     prometheus.NewDesc(pns+"transmit_errors_total", "Total Transmit Errors", labelP, nil),
		TxMulticast:  prometheus.NewDesc(pns+"transmit_multicast_total", "Total Tranmist Multicast", labelP, nil),
		TxPackets:    prometheus.NewDesc(pns+"transmit_packets_total", "Total Transmit Packets", labelP, nil),
	}
}

func (u *unifiCollector) exportUSWs(r *Report) {
	if r.Metrics == nil || r.Metrics.Devices == nil || len(r.Metrics.Devices.USWs) < 1 {
		return
	}
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for _, s := range r.Metrics.Devices.USWs {
			u.exportUSW(r, s)
		}
	}()
}

func (u *unifiCollector) exportUSW(r *Report, s *unifi.USW) {
	labels := []string{s.Type, s.Version, s.IP, s.SiteName, s.Mac, s.Model, s.Name, s.Serial}

	if s.HasTemperature.Val {
		r.send([]*metricExports{{u.USW.Temperature, prometheus.GaugeValue, s.GeneralTemperature, labels}})
	}
	if s.HasFan.Val {
		r.send([]*metricExports{{u.USW.FanLevel, prometheus.GaugeValue, s.FanLevel, labels}})
	}

	// Switch data.
	r.send([]*metricExports{
		{u.USW.Uptime, prometheus.GaugeValue, s.Uptime, labels},
		{u.USW.TotalMaxPower, prometheus.GaugeValue, s.TotalMaxPower, labels},
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
	})
	u.exportPortTable(r, s.PortTable, labels[5:])
}

func (u *unifiCollector) exportPortTable(r *Report, pt []unifi.Port, labels []string) {
	// Per-port data on a switch
	for _, p := range pt {
		if !p.Up.Val {
			continue
		}
		// Copy labels, and add four new ones.
		l := append([]string{p.PortIdx.Txt, p.Name, p.Mac, p.IP}, labels...)
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
			{u.USW.Speed, prometheus.GaugeValue, p.Speed, l},
			{u.USW.TxBroadcast, prometheus.CounterValue, p.TxBroadcast, l},
			{u.USW.TxBytes, prometheus.CounterValue, p.TxBytes, l},
			{u.USW.TxBytesR, prometheus.GaugeValue, p.TxBytesR, l},
			{u.USW.TxDropped, prometheus.CounterValue, p.TxDropped, l},
			{u.USW.TxErrors, prometheus.CounterValue, p.TxErrors, l},
			{u.USW.TxMulticast, prometheus.CounterValue, p.TxMulticast, l},
		})
	}
}
