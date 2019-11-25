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
		TotalTxBytes: prometheus.NewDesc(ns+"bytes_tx_total", "Total Transmitted Bytes", labels, nil),
		TotalRxBytes: prometheus.NewDesc(ns+"bytes_rx_total", "Total Received Bytes", labels, nil),
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
		MemUsed:      prometheus.NewDesc(ns+"memory_used_bytes", "System Memory Used", labels, nil),
		MemTotal:     prometheus.NewDesc(ns+"memory_installed_bytes", "System Installed Memory", labels, nil),
		MemBuffer:    prometheus.NewDesc(ns+"memory_buffer_bytes", "System Memory Buffer", labels, nil),
		CPU:          prometheus.NewDesc(ns+"cpu_utilization", "System CPU % Utilized", labels, nil),
		Mem:          prometheus.NewDesc(ns+"memory_utilization", "System Memory % Utilized", labels, nil),
	}
}

func (u *unifiCollector) exportUAPs(uaps []*unifi.UAP, ch chan []*metricExports) {
	for _, a := range uaps {
		ch <- u.exportUAP(a)
	}
}

// exportUAP exports Access Point Data
func (u *unifiCollector) exportUAP(a *unifi.UAP) []*metricExports {
	labels := []string{a.SiteName, a.Mac, a.Model, a.Name, a.Serial, a.SiteID,
		a.Type, a.Version, a.DeviceID, a.IP}

	// Switch data.
	metrics := []*metricExports{
		{u.UAP.Uptime, prometheus.GaugeValue, a.Uptime, labels},
		{u.UAP.TotalTxBytes, prometheus.CounterValue, a.TxBytes, labels},
		{u.UAP.TotalRxBytes, prometheus.CounterValue, a.RxBytes, labels},
		{u.UAP.TotalBytes, prometheus.CounterValue, a.Bytes, labels},
		{u.UAP.BytesD, prometheus.CounterValue, a.BytesD, labels}, // not sure if these 3 Ds are counters or gauges.
		{u.UAP.TxBytesD, prometheus.CounterValue, a.TxBytesD, labels},
		{u.UAP.RxBytesD, prometheus.CounterValue, a.RxBytesD, labels},
		{u.UAP.BytesR, prometheus.GaugeValue, a.BytesR, labels},
		{u.UAP.NumSta, prometheus.GaugeValue, a.NumSta, labels},
		{u.UAP.UserNumSta, prometheus.GaugeValue, a.UserNumSta, labels},
		{u.UAP.GuestNumSta, prometheus.GaugeValue, a.GuestNumSta, labels},
		{u.UAP.Loadavg1, prometheus.GaugeValue, a.SysStats.Loadavg1, labels},
		{u.UAP.Loadavg5, prometheus.GaugeValue, a.SysStats.Loadavg5, labels},
		{u.UAP.Loadavg15, prometheus.GaugeValue, a.SysStats.Loadavg15, labels},
		{u.UAP.MemUsed, prometheus.GaugeValue, a.SysStats.MemUsed, labels},
		{u.UAP.MemTotal, prometheus.GaugeValue, a.SysStats.MemTotal, labels},
		{u.UAP.MemBuffer, prometheus.GaugeValue, a.SysStats.MemBuffer, labels},
		{u.UAP.CPU, prometheus.GaugeValue, a.SystemStats.CPU, labels},
		{u.UAP.Mem, prometheus.GaugeValue, a.SystemStats.Mem, labels},
	}
	return metrics
}
