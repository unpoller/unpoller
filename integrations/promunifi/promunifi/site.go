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
	Uptime                *prometheus.Desc
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
	labels := []string{"subsystem", "status", "site_name"}
	return &site{
		NumUser:               prometheus.NewDesc(ns+"users", "Number of Users", labels, nil),
		NumGuest:              prometheus.NewDesc(ns+"guests", "Number of Guests", labels, nil),
		NumIot:                prometheus.NewDesc(ns+"iots", "Number of IoT Devices", labels, nil),
		TxBytesR:              prometheus.NewDesc(ns+"transmit_rate_bytes", "Bytes Transmit Rate", labels, nil),
		RxBytesR:              prometheus.NewDesc(ns+"receive_rate_bytes", "Bytes Receive Rate", labels, nil),
		NumAp:                 prometheus.NewDesc(ns+"aps", "Access Point Count", labels, nil),
		NumAdopted:            prometheus.NewDesc(ns+"adopted", "Adoption Count", labels, nil),
		NumDisabled:           prometheus.NewDesc(ns+"disabled", "Disabled Count", labels, nil),
		NumDisconnected:       prometheus.NewDesc(ns+"disconnected", "Disconnected Count", labels, nil),
		NumPending:            prometheus.NewDesc(ns+"pending", "Pending Count", labels, nil),
		NumGw:                 prometheus.NewDesc(ns+"gateways", "Gateway Count", labels, nil),
		NumSw:                 prometheus.NewDesc(ns+"switches", "Switch Count", labels, nil),
		NumSta:                prometheus.NewDesc(ns+"stations", "Station Count", labels, nil),
		Latency:               prometheus.NewDesc(ns+"latency_seconds", "Latency", labels, nil),
		Uptime:                prometheus.NewDesc(ns+"uptime_seconds", "Uptime", labels, nil),
		Drops:                 prometheus.NewDesc(ns+"intenet_drops_total", "Internet (WAN) Disconnections", labels, nil),
		XputUp:                prometheus.NewDesc(ns+"xput_up_rate", "Speedtest Upload", labels, nil),
		XputDown:              prometheus.NewDesc(ns+"xput_down_rate", "Speedtest Download", labels, nil),
		SpeedtestPing:         prometheus.NewDesc(ns+"speedtest_ping", "Speedtest Ping", labels, nil),
		RemoteUserNumActive:   prometheus.NewDesc(ns+"remote_user_active", "Remote Users Active", labels, nil),
		RemoteUserNumInactive: prometheus.NewDesc(ns+"remote_user_inactive", "Remote Users Inactive", labels, nil),
		RemoteUserRxBytes:     prometheus.NewDesc(ns+"remote_user_receive_bytes_total", "Remote Users Receive Bytes", labels, nil),
		RemoteUserTxBytes:     prometheus.NewDesc(ns+"remote_user_transmit_bytes_total", "Remote Users Transmit Bytes", labels, nil),
		RemoteUserRxPackets:   prometheus.NewDesc(ns+"remote_user_receive_packets_total", "Remote Users Receive Packets", labels, nil),
		RemoteUserTxPackets:   prometheus.NewDesc(ns+"remote_user_transmit_packets_total", "Remote Users Transmit Packets", labels, nil),
	}
}

func (u *unifiCollector) exportSite(r report, s *unifi.Site) {
	for _, h := range s.Health {
		labels := []string{h.Subsystem, h.Status, s.SiteName}
		switch h.Subsystem {
		case "www":
			r.send([]*metricExports{
				{u.Site.TxBytesR, prometheus.GaugeValue, h.TxBytesR, labels},
				{u.Site.RxBytesR, prometheus.GaugeValue, h.RxBytesR, labels},
				{u.Site.Uptime, prometheus.GaugeValue, h.Latency, labels},
				{u.Site.Latency, prometheus.GaugeValue, h.Latency.Val / 1000, labels},
				{u.Site.XputUp, prometheus.GaugeValue, h.XputUp, labels},
				{u.Site.XputDown, prometheus.GaugeValue, h.XputDown, labels},
				{u.Site.SpeedtestPing, prometheus.GaugeValue, h.SpeedtestPing, labels},
				{u.Site.Drops, prometheus.CounterValue, h.Drops, labels},
			})

		case "wlan":
			r.send([]*metricExports{
				{u.Site.TxBytesR, prometheus.GaugeValue, h.TxBytesR, labels},
				{u.Site.RxBytesR, prometheus.GaugeValue, h.RxBytesR, labels},
				{u.Site.NumAdopted, prometheus.GaugeValue, h.NumAdopted, labels},
				{u.Site.NumDisconnected, prometheus.GaugeValue, h.NumDisconnected, labels},
				{u.Site.NumPending, prometheus.GaugeValue, h.NumPending, labels},
				{u.Site.NumUser, prometheus.GaugeValue, h.NumUser, labels},
				{u.Site.NumGuest, prometheus.GaugeValue, h.NumGuest, labels},
				{u.Site.NumIot, prometheus.GaugeValue, h.NumIot, labels},
				{u.Site.NumAp, prometheus.GaugeValue, h.NumAp, labels},
				{u.Site.NumDisabled, prometheus.GaugeValue, h.NumDisabled, labels},
			})

		case "wan":
			r.send([]*metricExports{
				{u.Site.TxBytesR, prometheus.GaugeValue, h.TxBytesR, labels},
				{u.Site.RxBytesR, prometheus.GaugeValue, h.RxBytesR, labels},
				{u.Site.NumAdopted, prometheus.GaugeValue, h.NumAdopted, labels},
				{u.Site.NumDisconnected, prometheus.GaugeValue, h.NumDisconnected, labels},
				{u.Site.NumPending, prometheus.GaugeValue, h.NumPending, labels},
				{u.Site.NumGw, prometheus.GaugeValue, h.NumGw, labels},
				{u.Site.NumSta, prometheus.GaugeValue, h.NumSta, labels},
			})

		case "lan":
			r.send([]*metricExports{
				{u.Site.TxBytesR, prometheus.GaugeValue, h.TxBytesR, labels},
				{u.Site.RxBytesR, prometheus.GaugeValue, h.RxBytesR, labels},
				{u.Site.NumAdopted, prometheus.GaugeValue, h.NumAdopted, labels},
				{u.Site.NumDisconnected, prometheus.GaugeValue, h.NumDisconnected, labels},
				{u.Site.NumPending, prometheus.GaugeValue, h.NumPending, labels},
				{u.Site.NumUser, prometheus.GaugeValue, h.NumUser, labels},
				{u.Site.NumGuest, prometheus.GaugeValue, h.NumGuest, labels},
				{u.Site.NumIot, prometheus.GaugeValue, h.NumIot, labels},
				{u.Site.NumSw, prometheus.GaugeValue, h.NumSw, labels},
			})

		case "vpn":
			r.send([]*metricExports{
				{u.Site.RemoteUserNumActive, prometheus.GaugeValue, h.RemoteUserNumActive, labels},
				{u.Site.RemoteUserNumInactive, prometheus.GaugeValue, h.RemoteUserNumInactive, labels},
				{u.Site.RemoteUserRxBytes, prometheus.CounterValue, h.RemoteUserRxBytes, labels},
				{u.Site.RemoteUserTxBytes, prometheus.CounterValue, h.RemoteUserTxBytes, labels},
				{u.Site.RemoteUserRxPackets, prometheus.CounterValue, h.RemoteUserRxPackets, labels},
				{u.Site.RemoteUserTxPackets, prometheus.CounterValue, h.RemoteUserTxPackets, labels},
			})
		}
	}
}
