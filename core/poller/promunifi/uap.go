package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type uap struct {
	Uptime       *prometheus.Desc
	TotalTxBytes *prometheus.Desc
	TotalRxBytes *prometheus.Desc
	TotalBytes   *prometheus.Desc
	BytesD       *prometheus.Desc
	TxBytesD     *prometheus.Desc
	RxBytesD     *prometheus.Desc
	BytesR       *prometheus.Desc
	NumSta       *prometheus.Desc
	UserNumSta   *prometheus.Desc
	GuestNumSta  *prometheus.Desc
	// System Stats
	Loadavg1  *prometheus.Desc
	Loadavg5  *prometheus.Desc
	Loadavg15 *prometheus.Desc
	MemBuffer *prometheus.Desc
	MemTotal  *prometheus.Desc
	MemUsed   *prometheus.Desc
	CPU       *prometheus.Desc
	Mem       *prometheus.Desc
	// Ap Traffic Stats -- not sure about these yet.
	ApBytes                  *prometheus.Desc
	ApWifiTxDropped          *prometheus.Desc
	ApRxErrors               *prometheus.Desc
	ApRxDropped              *prometheus.Desc
	ApRxFrags                *prometheus.Desc
	ApRxCrypts               *prometheus.Desc
	ApTxPackets              *prometheus.Desc
	ApTxBytes                *prometheus.Desc
	ApTxErrors               *prometheus.Desc
	ApTxDropped              *prometheus.Desc
	ApTxRetries              *prometheus.Desc
	ApRxPackets              *prometheus.Desc
	ApRxBytes                *prometheus.Desc
	UserRxDropped            *prometheus.Desc
	GuestRxDropped           *prometheus.Desc
	UserRxErrors             *prometheus.Desc
	GuestRxErrors            *prometheus.Desc
	UserRxPackets            *prometheus.Desc
	GuestRxPackets           *prometheus.Desc
	UserRxBytes              *prometheus.Desc
	GuestRxBytes             *prometheus.Desc
	UserRxCrypts             *prometheus.Desc
	GuestRxCrypts            *prometheus.Desc
	UserRxFrags              *prometheus.Desc
	GuestRxFrags             *prometheus.Desc
	UserTxPackets            *prometheus.Desc
	GuestTxPackets           *prometheus.Desc
	UserTxBytes              *prometheus.Desc
	GuestTxBytes             *prometheus.Desc
	UserTxErrors             *prometheus.Desc
	GuestTxErrors            *prometheus.Desc
	UserTxDropped            *prometheus.Desc
	GuestTxDropped           *prometheus.Desc
	UserTxRetries            *prometheus.Desc
	GuestTxRetries           *prometheus.Desc
	MacFilterRejections      *prometheus.Desc
	UserMacFilterRejections  *prometheus.Desc
	GuestMacFilterRejections *prometheus.Desc
	WifiTxAttempts           *prometheus.Desc
	UserWifiTxDropped        *prometheus.Desc
	GuestWifiTxDropped       *prometheus.Desc
	UserWifiTxAttempts       *prometheus.Desc
	GuestWifiTxAttempts      *prometheus.Desc
}

func descUAP(ns string) *uap {
	if ns += "_uap_"; ns == "_uap_" {
		ns = "uap_"
	}
	labels := []string{"site_name", "mac", "model", "name", "serial", "site_id",
		"type", "version", "device_id", "ip"}

	return &uap{
		Uptime:       prometheus.NewDesc(ns+"uptime", "Uptime", labels, nil),
		TotalTxBytes: prometheus.NewDesc(ns+"tx_bytes_total", "Total Transmitted Bytes", labels, nil),
		TotalRxBytes: prometheus.NewDesc(ns+"rx_bytes_total", "Total Received Bytes", labels, nil),
		TotalBytes:   prometheus.NewDesc(ns+"bytes_total", "Total Bytes Transfered", labels, nil),
		BytesD:       prometheus.NewDesc(ns+"bytes_d_total", "Total Bytes D???", labels, nil),
		TxBytesD:     prometheus.NewDesc(ns+"bytes_d_tx", "Transmit Bytes D???", labels, nil),
		RxBytesD:     prometheus.NewDesc(ns+"bytes_d_rx", "Receive Bytes D???", labels, nil),
		BytesR:       prometheus.NewDesc(ns+"bytes_rate", "Transfer Rate", labels, nil),
		NumSta:       prometheus.NewDesc(ns+"stations_total", "Number of Stations", labels, nil),
		UserNumSta:   prometheus.NewDesc(ns+"stations_user_total", "Number of User Stations", labels, nil),
		GuestNumSta:  prometheus.NewDesc(ns+"stations_guest_total", "Number of Guest Stations", labels, nil),
		Loadavg1:     prometheus.NewDesc(ns+"load_average_1", "System Load Average 1 Minute", labels, nil),
		Loadavg5:     prometheus.NewDesc(ns+"load_average_5", "System Load Average 5 Minutes", labels, nil),
		Loadavg15:    prometheus.NewDesc(ns+"load_average_15", "System Load Average 15 Minutes", labels, nil),
		MemUsed:      prometheus.NewDesc(ns+"memory_utilization", "System Memory Used", labels, nil),
		MemTotal:     prometheus.NewDesc(ns+"memory_installed", "System Installed Memory", labels, nil),
		MemBuffer:    prometheus.NewDesc(ns+"memory_buffer", "System Memory Buffer", labels, nil),
		CPU:          prometheus.NewDesc(ns+"cpu_utilization", "System CPU % Utilized", labels, nil),
		Mem:          prometheus.NewDesc(ns+"memory", "System Memory % Utilized", labels, nil),
	}
}

// exportUAP exports Access Point Data
func (u *unifiCollector) exportUAP(s *unifi.UAP) []*metricExports {
	labels := []string{s.SiteName, s.Mac, s.Model, s.Name, s.Serial, s.SiteID,
		s.Type, s.Version, s.DeviceID, s.IP}

	// Switch data.
	m := []*metricExports{
		{u.UAP.Uptime, prometheus.GaugeValue, s.Uptime, labels},
		{u.UAP.TotalTxBytes, prometheus.CounterValue, s.TxBytes, labels},
		{u.UAP.TotalRxBytes, prometheus.CounterValue, s.RxBytes, labels},
		{u.UAP.TotalBytes, prometheus.CounterValue, s.Bytes, labels},
		{u.UAP.BytesD, prometheus.CounterValue, s.BytesD, labels}, // not sure if these 3 Ds are counters or gauges.
		{u.UAP.TxBytesD, prometheus.CounterValue, s.TxBytesD, labels},
		{u.UAP.RxBytesD, prometheus.CounterValue, s.RxBytesD, labels},
		{u.UAP.BytesR, prometheus.GaugeValue, s.BytesR, labels},
		{u.UAP.NumSta, prometheus.GaugeValue, s.NumSta, labels},
		{u.UAP.UserNumSta, prometheus.GaugeValue, s.UserNumSta, labels},
		{u.UAP.GuestNumSta, prometheus.GaugeValue, s.GuestNumSta, labels},
		{u.UAP.Loadavg1, prometheus.GaugeValue, s.SysStats.Loadavg1, labels},
		{u.UAP.Loadavg5, prometheus.GaugeValue, s.SysStats.Loadavg5, labels},
		{u.UAP.Loadavg15, prometheus.GaugeValue, s.SysStats.Loadavg15, labels},
		{u.UAP.MemUsed, prometheus.GaugeValue, s.SysStats.MemUsed, labels},
		{u.UAP.MemTotal, prometheus.GaugeValue, s.SysStats.MemTotal, labels},
		{u.UAP.MemBuffer, prometheus.GaugeValue, s.SysStats.MemBuffer, labels},
		{u.UAP.CPU, prometheus.GaugeValue, s.SystemStats.CPU, labels},
		{u.UAP.Mem, prometheus.GaugeValue, s.SystemStats.Mem, labels},
	}
	return m
}
