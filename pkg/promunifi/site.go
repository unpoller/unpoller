package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
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
	DPITxPackets          *prometheus.Desc
	DPIRxPackets          *prometheus.Desc
	DPITxBytes            *prometheus.Desc
	DPIRxBytes            *prometheus.Desc
}

func descSite(ns string) *site {
	labels := []string{"subsystem", "status", "site_name", "source"}
	labelDPI := []string{"category", "application", "site_name", "source"}
	nd := prometheus.NewDesc

	return &site{
		NumUser:               nd(ns+"users", "Number of Users", labels, nil),
		NumGuest:              nd(ns+"guests", "Number of Guests", labels, nil),
		NumIot:                nd(ns+"iots", "Number of IoT Devices", labels, nil),
		TxBytesR:              nd(ns+"transmit_rate_bytes", "Bytes Transmit Rate", labels, nil),
		RxBytesR:              nd(ns+"receive_rate_bytes", "Bytes Receive Rate", labels, nil),
		NumAp:                 nd(ns+"aps", "Access Point Count", labels, nil),
		NumAdopted:            nd(ns+"adopted", "Adoption Count", labels, nil),
		NumDisabled:           nd(ns+"disabled", "Disabled Count", labels, nil),
		NumDisconnected:       nd(ns+"disconnected", "Disconnected Count", labels, nil),
		NumPending:            nd(ns+"pending", "Pending Count", labels, nil),
		NumGw:                 nd(ns+"gateways", "Gateway Count", labels, nil),
		NumSw:                 nd(ns+"switches", "Switch Count", labels, nil),
		NumSta:                nd(ns+"stations", "Station Count", labels, nil),
		Latency:               nd(ns+"latency_seconds", "Latency", labels, nil),
		Uptime:                nd(ns+"uptime_seconds", "Uptime", labels, nil),
		Drops:                 nd(ns+"intenet_drops_total", "Internet (WAN) Disconnections", labels, nil),
		XputUp:                nd(ns+"xput_up_rate", "Speedtest Upload", labels, nil),
		XputDown:              nd(ns+"xput_down_rate", "Speedtest Download", labels, nil),
		SpeedtestPing:         nd(ns+"speedtest_ping", "Speedtest Ping", labels, nil),
		RemoteUserNumActive:   nd(ns+"remote_user_active", "Remote Users Active", labels, nil),
		RemoteUserNumInactive: nd(ns+"remote_user_inactive", "Remote Users Inactive", labels, nil),
		RemoteUserRxBytes:     nd(ns+"remote_user_receive_bytes_total", "Remote Users Receive Bytes", labels, nil),
		RemoteUserTxBytes:     nd(ns+"remote_user_transmit_bytes_total", "Remote Users Transmit Bytes", labels, nil),
		RemoteUserRxPackets:   nd(ns+"remote_user_receive_packets_total", "Remote Users Receive Packets", labels, nil),
		RemoteUserTxPackets:   nd(ns+"remote_user_transmit_packets_total", "Remote Users Transmit Packets", labels, nil),
		DPITxPackets:          nd(ns+"dpi_transmit_packets", "Site DPI Transmit Packets", labelDPI, nil),
		DPIRxPackets:          nd(ns+"dpi_receive_packets", "Site DPI Receive Packets", labelDPI, nil),
		DPITxBytes:            nd(ns+"dpi_transmit_bytes", "Site DPI Transmit Bytes", labelDPI, nil),
		DPIRxBytes:            nd(ns+"dpi_receive_bytes", "Site DPI Receive Bytes", labelDPI, nil),
	}
}

func (u *promUnifi) exportSiteDPI(r report, v any) {
	s, ok := v.(*unifi.DPITable)
	if !ok {
		u.LogErrorf("invalid type given to SiteDPI: %T", v)

		return
	}

	for _, dpi := range s.ByApp {
		labelDPI := []string{unifi.DPICats.Get(dpi.Cat.Int()), unifi.DPIApps.GetApp(dpi.Cat.Int(), dpi.App.Int()), s.SiteName, s.SourceName}

		//	log.Println(labelsDPI, dpi.Cat, dpi.App, dpi.TxBytes, dpi.RxBytes, dpi.TxPackets, dpi.RxPackets)
		r.send([]*metric{
			{u.Site.DPITxPackets, gauge, dpi.TxPackets.Val, labelDPI},
			{u.Site.DPIRxPackets, gauge, dpi.RxPackets.Val, labelDPI},
			{u.Site.DPITxBytes, gauge, dpi.TxBytes.Val, labelDPI},
			{u.Site.DPIRxBytes, gauge, dpi.RxBytes.Val, labelDPI},
		})
	}
}

func (u *promUnifi) exportSite(r report, s *unifi.Site) {
	for _, h := range s.Health {
		switch labels := []string{h.Subsystem, h.Status, s.SiteName, s.SourceName}; labels[0] {
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
