package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type site struct {
	NumUser               *prometheus.Desc
	NumGuest              *prometheus.Desc
	NumIot                *prometheus.Desc
	TxBytesR              *prometheus.Desc
	RxBytesR              *prometheus.Desc
	NumAp                 *prometheus.Desc
	NumAdopted            *prometheus.Desc
	NumDisabled           *prometheus.Desc
	NumDisconnected       *prometheus.Desc
	NumPending            *prometheus.Desc
	NumGw                 *prometheus.Desc
	NumSw                 *prometheus.Desc
	NumSta                *prometheus.Desc
	Latency               *prometheus.Desc
	Drops                 *prometheus.Desc
	XputUp                *prometheus.Desc
	XputDown              *prometheus.Desc
	SpeedtestPing         *prometheus.Desc
	RemoteUserNumActive   *prometheus.Desc
	RemoteUserNumInactive *prometheus.Desc
	RemoteUserRxBytes     *prometheus.Desc
	RemoteUserTxBytes     *prometheus.Desc
	RemoteUserRxPackets   *prometheus.Desc
	RemoteUserTxPackets   *prometheus.Desc
}

func descSite(ns string) *site {
	if ns += "_site_"; ns == "_site_" {
		ns = "site_"
	}
	labels := []string{"subsystem", "status", "gwversion", "name", "desc", "site_name"}

	return &site{
		NumUser:               prometheus.NewDesc(ns+"user_total", "Number of Users", labels, nil),
		NumGuest:              prometheus.NewDesc(ns+"guest_total", "Number of Guests", labels, nil),
		NumIot:                prometheus.NewDesc(ns+"iot_total", "Number of IoT Devices", labels, nil),
		TxBytesR:              prometheus.NewDesc(ns+"bytes_tx_rate", "Bytes Transmit Rate", labels, nil),
		RxBytesR:              prometheus.NewDesc(ns+"bytes_rx_rate", "Bytes Receive Rate", labels, nil),
		NumAp:                 prometheus.NewDesc(ns+"ap_total", "Access Point Count", labels, nil),
		NumAdopted:            prometheus.NewDesc(ns+"adopted_total", "Adoption Count", labels, nil),
		NumDisabled:           prometheus.NewDesc(ns+"disabled_total", "Disabled Count", labels, nil),
		NumDisconnected:       prometheus.NewDesc(ns+"disconnected_total", "Disconnected Count", labels, nil),
		NumPending:            prometheus.NewDesc(ns+"pending_total", "Pending Count", labels, nil),
		NumGw:                 prometheus.NewDesc(ns+"gateways_total", "Gateway Count", labels, nil),
		NumSw:                 prometheus.NewDesc(ns+"switches_total", "Switch Count", labels, nil),
		NumSta:                prometheus.NewDesc(ns+"stations_total", "Station Count", labels, nil),
		Latency:               prometheus.NewDesc(ns+"latency", "Latency", labels, nil),
		Drops:                 prometheus.NewDesc(ns+"drops_total", "Drops", labels, nil),
		XputUp:                prometheus.NewDesc(ns+"xput_up_rate", "Speedtest Upload", labels, nil),
		XputDown:              prometheus.NewDesc(ns+"xput_down_rate", "Speedtest Download", labels, nil),
		SpeedtestPing:         prometheus.NewDesc(ns+"speedtest_ping", "Speedtest Ping", labels, nil),
		RemoteUserNumActive:   prometheus.NewDesc(ns+"remote_user_active_total", "Remote Users Active", labels, nil),
		RemoteUserNumInactive: prometheus.NewDesc(ns+"remote_user_inactive_total", "Remote Users Inactive", labels, nil),
		RemoteUserRxBytes:     prometheus.NewDesc(ns+"remote_user_rx_bytes_total", "Remote Users Receive Bytes", labels, nil),
		RemoteUserTxBytes:     prometheus.NewDesc(ns+"remote_user_tx_bytes_total", "Remote Users Transmit Bytes", labels, nil),
		RemoteUserRxPackets:   prometheus.NewDesc(ns+"remote_user_rx_packets_total", "Remote Users Receive Packets", labels, nil),
		RemoteUserTxPackets:   prometheus.NewDesc(ns+"remote_user_tx_packets_total", "Remote Users Transmit Packets", labels, nil),
	}
}

func (u *unifiCollector) exportSites(sites unifi.Sites, r *Report) {
	for _, s := range sites {
		metrics := []*metricExports{}
		labels := []string{s.Name, s.Desc, s.SiteName}
		for _, h := range s.Health {
			l := append([]string{h.Subsystem, h.Status, h.GwVersion}, labels...)

			// XXX: More of these are subsystem specific (like the vpn/remote user stuff below)
			metrics = append(metrics, []*metricExports{
				{u.Site.NumUser, prometheus.CounterValue, h.NumUser.Val, l},
				{u.Site.NumGuest, prometheus.CounterValue, h.NumGuest.Val, l},
				{u.Site.NumIot, prometheus.CounterValue, h.NumIot.Val, l},
				{u.Site.TxBytesR, prometheus.GaugeValue, h.TxBytesR.Val, l},
				{u.Site.RxBytesR, prometheus.GaugeValue, h.RxBytesR.Val, l},
				{u.Site.NumAp, prometheus.CounterValue, h.NumAp.Val, l},
				{u.Site.NumAdopted, prometheus.CounterValue, h.NumAdopted.Val, l},
				{u.Site.NumDisabled, prometheus.CounterValue, h.NumDisabled.Val, l},
				{u.Site.NumDisconnected, prometheus.CounterValue, h.NumDisconnected.Val, l},
				{u.Site.NumPending, prometheus.CounterValue, h.NumPending.Val, l},
				{u.Site.NumGw, prometheus.CounterValue, h.NumGw.Val, l},
				{u.Site.NumSw, prometheus.CounterValue, h.NumSw.Val, l},
				{u.Site.NumSta, prometheus.CounterValue, h.NumSta.Val, l},
				{u.Site.Latency, prometheus.GaugeValue, h.Latency.Val, l},
				{u.Site.Drops, prometheus.CounterValue, h.Drops.Val, l},
				{u.Site.XputUp, prometheus.GaugeValue, h.XputUp.Val, l},
				{u.Site.XputDown, prometheus.GaugeValue, h.XputDown.Val, l},
				{u.Site.SpeedtestPing, prometheus.GaugeValue, h.SpeedtestPing.Val, l},
			}...)

			if h.Subsystem == "vpn" {
				metrics = append(metrics, []*metricExports{
					{u.Site.RemoteUserNumActive, prometheus.CounterValue, h.RemoteUserNumActive.Val, l},
					{u.Site.RemoteUserNumInactive, prometheus.CounterValue, h.RemoteUserNumInactive.Val, l},
					{u.Site.RemoteUserRxBytes, prometheus.CounterValue, h.RemoteUserRxBytes.Val, l},
					{u.Site.RemoteUserTxBytes, prometheus.CounterValue, h.RemoteUserTxBytes.Val, l},
					{u.Site.RemoteUserRxPackets, prometheus.CounterValue, h.RemoteUserRxPackets.Val, l},
					{u.Site.RemoteUserTxPackets, prometheus.CounterValue, h.RemoteUserTxPackets.Val, l},
				}...)
			}
		}
		r.ch <- metrics
	}
}
