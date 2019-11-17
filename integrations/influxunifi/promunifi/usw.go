package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type usw struct {
	Uptime        *prometheus.Desc `json:"uptime"`
	Temperature   *prometheus.Desc `json:"general_temperature"`
	TotalMaxPower *prometheus.Desc `json:"total_max_power"`
	FanLevel      *prometheus.Desc `json:"fan_level"`
	TotalTxBytes  *prometheus.Desc `json:"total_tx_bytes"`
	TotalRxBytes  *prometheus.Desc `json:"total_rx_bytes"`
	TotalBytes    *prometheus.Desc `json:"bytes"`
	NumSta        *prometheus.Desc `json:"num_sta"`
	UserNumSta    *prometheus.Desc `json:"user-num_sta"`
	GuestNumSta   *prometheus.Desc `json:"guest-num_sta"`
	// Port data.
	PoeCurrent   *prometheus.Desc `json:"poe_current,omitempty"`
	PoePower     *prometheus.Desc `json:"poe_power,omitempty"`
	PoeVoltage   *prometheus.Desc `json:"poe_voltage,omitempty"`
	RxBroadcast  *prometheus.Desc `json:"rx_broadcast"`
	RxBytes      *prometheus.Desc `json:"rx_bytes"`
	RxBytesR     *prometheus.Desc `json:"rx_bytes-r"`
	RxDropped    *prometheus.Desc `json:"rx_dropped"`
	RxErrors     *prometheus.Desc `json:"rx_errors"`
	RxMulticast  *prometheus.Desc `json:"rx_multicast"`
	RxPackets    *prometheus.Desc `json:"rx_packets"`
	Satisfaction *prometheus.Desc `json:"satisfaction,omitempty"`
	Speed        *prometheus.Desc `json:"speed"`
	TxBroadcast  *prometheus.Desc `json:"tx_broadcast"`
	TxBytes      *prometheus.Desc `json:"tx_bytes"`
	TxBytesR     *prometheus.Desc `json:"tx_bytes-r"`
	TxDropped    *prometheus.Desc `json:"tx_dropped"`
	TxErrors     *prometheus.Desc `json:"tx_errors"`
	TxMulticast  *prometheus.Desc `json:"tx_multicast"`
	TxPackets    *prometheus.Desc `json:"tx_packets"`
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
		Uptime:        prometheus.NewDesc(ns+"Uptime", "Uptime", labels, nil),
		Temperature:   prometheus.NewDesc(ns+"Temperature", "Temperature", labels, nil),
		TotalMaxPower: prometheus.NewDesc(ns+"TotalMaxPower", "TotalMaxPower", labels, nil),
		FanLevel:      prometheus.NewDesc(ns+"FanLevel", "FanLevel", labels, nil),
		TotalTxBytes:  prometheus.NewDesc(ns+"TxBytes", "TxBytes", labels, nil),
		TotalRxBytes:  prometheus.NewDesc(ns+"RxBytes", "RxBytes", labels, nil),
		TotalBytes:    prometheus.NewDesc(ns+"Bytes", "Bytes", labels, nil),
		NumSta:        prometheus.NewDesc(ns+"NumSta", "NumSta", labels, nil),
		UserNumSta:    prometheus.NewDesc(ns+"UserNumSta", "UserNumSta", labels, nil),
		GuestNumSta:   prometheus.NewDesc(ns+"GuestNumSta", "GuestNumSta", labels, nil),
		// per-port data
		PoeCurrent:   prometheus.NewDesc(pns+"PoeCurrent", "PoeCurrent", labelP, nil),
		PoePower:     prometheus.NewDesc(pns+"PoePower", "PoePower", labelP, nil),
		PoeVoltage:   prometheus.NewDesc(pns+"PoeVoltage", "PoeVoltage", labelP, nil),
		RxBroadcast:  prometheus.NewDesc(pns+"RxBroadcast", "RxBroadcast", labelP, nil),
		RxBytes:      prometheus.NewDesc(pns+"RxBytes", "RxBytes", labelP, nil),
		RxBytesR:     prometheus.NewDesc(pns+"RxBytesR", "RxBytesR", labelP, nil),
		RxDropped:    prometheus.NewDesc(pns+"RxDropped", "RxDropped", labelP, nil),
		RxErrors:     prometheus.NewDesc(pns+"RxErrors", "RxErrors", labelP, nil),
		RxMulticast:  prometheus.NewDesc(pns+"RxMulticast", "RxMulticast", labelP, nil),
		RxPackets:    prometheus.NewDesc(pns+"RxPackets", "RxPackets", labelP, nil),
		Satisfaction: prometheus.NewDesc(pns+"Satisfaction", "Satisfaction", labelP, nil),
		Speed:        prometheus.NewDesc(pns+"Speed", "Speed", labelP, nil),
		TxBroadcast:  prometheus.NewDesc(pns+"TxBroadcast", "TxBroadcast", labelP, nil),
		TxBytes:      prometheus.NewDesc(pns+"TxBytes", "TxBytes", labelP, nil),
		TxBytesR:     prometheus.NewDesc(pns+"TxBytesR", "TxBytesR", labelP, nil),
		TxDropped:    prometheus.NewDesc(pns+"TxDropped", "TxDropped", labelP, nil),
		TxErrors:     prometheus.NewDesc(pns+"TxErrors", "TxErrors", labelP, nil),
		TxMulticast:  prometheus.NewDesc(pns+"TxMulticast", "TxMulticast", labelP, nil),
		TxPackets:    prometheus.NewDesc(pns+"TxPackets", "TxPackets", labelP, nil),
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
	}

	// Per-port data on the switch
	for _, p := range s.PortTable {
		// Copy labels, and replace last four with different data.
		l := append(append([]string{}, labels[:6]...), p.PortIdx.Txt, p.Name, p.Mac, p.IP)
		m = append(m, &metricExports{u.USW.PoeCurrent, prometheus.GaugeValue, p.PoeCurrent, l})
		m = append(m, &metricExports{u.USW.PoePower, prometheus.GaugeValue, p.PoePower, l})
		m = append(m, &metricExports{u.USW.PoeVoltage, prometheus.GaugeValue, p.PoeVoltage, l})
		m = append(m, &metricExports{u.USW.RxBroadcast, prometheus.CounterValue, p.RxBroadcast, l})
		m = append(m, &metricExports{u.USW.RxBytes, prometheus.CounterValue, p.RxBytes, l})
		m = append(m, &metricExports{u.USW.RxBytesR, prometheus.GaugeValue, p.RxBytesR, l})
		m = append(m, &metricExports{u.USW.RxDropped, prometheus.CounterValue, p.RxDropped, l})
		m = append(m, &metricExports{u.USW.RxErrors, prometheus.CounterValue, p.RxErrors, l})
		m = append(m, &metricExports{u.USW.RxMulticast, prometheus.CounterValue, p.RxMulticast, l})
		m = append(m, &metricExports{u.USW.RxPackets, prometheus.CounterValue, p.RxPackets, l})
		m = append(m, &metricExports{u.USW.Satisfaction, prometheus.GaugeValue, p.Satisfaction, l})
		m = append(m, &metricExports{u.USW.Speed, prometheus.GaugeValue, p.Speed, l})
		m = append(m, &metricExports{u.USW.TxBroadcast, prometheus.CounterValue, p.TxBroadcast, l})
		m = append(m, &metricExports{u.USW.TxBytes, prometheus.CounterValue, p.TxBytes, l})
		m = append(m, &metricExports{u.USW.TxBytesR, prometheus.GaugeValue, p.TxBytesR, l})
		m = append(m, &metricExports{u.USW.TxDropped, prometheus.CounterValue, p.TxDropped, l})
		m = append(m, &metricExports{u.USW.TxErrors, prometheus.CounterValue, p.TxErrors, l})
		m = append(m, &metricExports{u.USW.TxMulticast, prometheus.CounterValue, p.TxMulticast, l})
	}

	return m
}
