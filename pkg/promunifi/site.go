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

func (u *promUnifi) exportSite(r report, s *unifi.Site) {
	for _, h := range s.Health {
		labels := []string{h.Subsystem, h.Status, s.SiteName}
		switch h.Subsystem {
		case "www":
			r.send([]*metric{
				{u.Site.TxBytesR, gauge, h.TxBytesR, labels},
				{u.Site.RxBytesR, gauge, h.RxBytesR, labels},
				{u.Site.Uptime, gauge, h.Uptime, labels},
				{u.Site.Latency, gauge, h.Latency.Val / 1000, labels},
				{u.Site.XputUp, gauge, h.XputUp, labels},
				{u.Site.XputDown, gauge, h.XputDown, labels},
				{u.Site.SpeedtestPing, gauge, h.SpeedtestPing, labels},
				{u.Site.Drops, counter, h.Drops, labels},
			})

		case "wlan":
			r.send([]*metric{
				{u.Site.TxBytesR, gauge, h.TxBytesR, labels},
				{u.Site.RxBytesR, gauge, h.RxBytesR, labels},
				{u.Site.NumAdopted, gauge, h.NumAdopted, labels},
				{u.Site.NumDisconnected, gauge, h.NumDisconnected, labels},
				{u.Site.NumPending, gauge, h.NumPending, labels},
				{u.Site.NumUser, gauge, h.NumUser, labels},
				{u.Site.NumGuest, gauge, h.NumGuest, labels},
				{u.Site.NumIot, gauge, h.NumIot, labels},
				{u.Site.NumAp, gauge, h.NumAp, labels},
				{u.Site.NumDisabled, gauge, h.NumDisabled, labels},
			})

		case "wan":
			r.send([]*metric{
				{u.Site.TxBytesR, gauge, h.TxBytesR, labels},
				{u.Site.RxBytesR, gauge, h.RxBytesR, labels},
				{u.Site.NumAdopted, gauge, h.NumAdopted, labels},
				{u.Site.NumDisconnected, gauge, h.NumDisconnected, labels},
				{u.Site.NumPending, gauge, h.NumPending, labels},
				{u.Site.NumGw, gauge, h.NumGw, labels},
				{u.Site.NumSta, gauge, h.NumSta, labels},
			})

		case "lan":
			r.send([]*metric{
				{u.Site.TxBytesR, gauge, h.TxBytesR, labels},
				{u.Site.RxBytesR, gauge, h.RxBytesR, labels},
				{u.Site.NumAdopted, gauge, h.NumAdopted, labels},
				{u.Site.NumDisconnected, gauge, h.NumDisconnected, labels},
				{u.Site.NumPending, gauge, h.NumPending, labels},
				{u.Site.NumUser, gauge, h.NumUser, labels},
				{u.Site.NumGuest, gauge, h.NumGuest, labels},
				{u.Site.NumIot, gauge, h.NumIot, labels},
				{u.Site.NumSw, gauge, h.NumSw, labels},
			})

		case "vpn":
			r.send([]*metric{
				{u.Site.RemoteUserNumActive, gauge, h.RemoteUserNumActive, labels},
				{u.Site.RemoteUserNumInactive, gauge, h.RemoteUserNumInactive, labels},
				{u.Site.RemoteUserRxBytes, counter, h.RemoteUserRxBytes, labels},
				{u.Site.RemoteUserTxBytes, counter, h.RemoteUserTxBytes, labels},
				{u.Site.RemoteUserRxPackets, counter, h.RemoteUserRxPackets, labels},
				{u.Site.RemoteUserTxPackets, counter, h.RemoteUserTxPackets, labels},
			})
		}
	}
}
