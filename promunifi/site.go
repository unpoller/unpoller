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
	if ns += "_site_"; ns == "_site_" {
		ns = "site_"
	}
	labels := []string{"subsystem", "status", "name", "desc", "site_name"}

	return &site{
		NumUser:               prometheus.NewDesc(ns+"num_user_total", "Number of Users", labels, nil),
		NumGuest:              prometheus.NewDesc(ns+"num_guest_total", "Number of Guests", labels, nil),
		NumIot:                prometheus.NewDesc(ns+"num_iot_total", "Number of IoT Devices", labels, nil),
		TxBytesR:              prometheus.NewDesc(ns+"transmit_rate_bytes", "Bytes Transmit Rate", labels, nil),
		RxBytesR:              prometheus.NewDesc(ns+"receive_rate_bytes", "Bytes Receive Rate", labels, nil),
		NumAp:                 prometheus.NewDesc(ns+"num_ap_total", "Access Point Count", labels, nil),
		NumAdopted:            prometheus.NewDesc(ns+"num_adopted_total", "Adoption Count", labels, nil),
		NumDisabled:           prometheus.NewDesc(ns+"num_disabled_total", "Disabled Count", labels, nil),
		NumDisconnected:       prometheus.NewDesc(ns+"num_disconnected_total", "Disconnected Count", labels, nil),
		NumPending:            prometheus.NewDesc(ns+"num_pending_total", "Pending Count", labels, nil),
		NumGw:                 prometheus.NewDesc(ns+"num_gateways_total", "Gateway Count", labels, nil),
		NumSw:                 prometheus.NewDesc(ns+"num_switches_total", "Switch Count", labels, nil),
		NumSta:                prometheus.NewDesc(ns+"num_stations_total", "Station Count", labels, nil),
		Latency:               prometheus.NewDesc(ns+"latency_ms", "Latency", labels, nil),
		Uptime:                prometheus.NewDesc(ns+"uptime_seconds", "Uptime", labels, nil),
		Drops:                 prometheus.NewDesc(ns+"intenet_drops_total", "Internet (WAN) Disconnections", labels, nil),
		XputUp:                prometheus.NewDesc(ns+"xput_up_rate", "Speedtest Upload", labels, nil),
		XputDown:              prometheus.NewDesc(ns+"xput_down_rate", "Speedtest Download", labels, nil),
		SpeedtestPing:         prometheus.NewDesc(ns+"speedtest_ping", "Speedtest Ping", labels, nil),
		RemoteUserNumActive:   prometheus.NewDesc(ns+"num_remote_user_active_total", "Remote Users Active", labels, nil),
		RemoteUserNumInactive: prometheus.NewDesc(ns+"num_remote_user_inactive_total", "Remote Users Inactive", labels, nil),
		RemoteUserRxBytes:     prometheus.NewDesc(ns+"remote_user_receive_bytes_total", "Remote Users Receive Bytes", labels, nil),
		RemoteUserTxBytes:     prometheus.NewDesc(ns+"remote_user_transmit_bytes_total", "Remote Users Transmit Bytes", labels, nil),
		RemoteUserRxPackets:   prometheus.NewDesc(ns+"remote_user_receive_packets_total", "Remote Users Receive Packets", labels, nil),
		RemoteUserTxPackets:   prometheus.NewDesc(ns+"remote_user_transmit_packets_total", "Remote Users Transmit Packets", labels, nil),
	}
}

func (u *unifiCollector) exportSites(sites unifi.Sites, r *Report) {
	for _, s := range sites {
		metrics := []*metricExports{}
		labels := []string{s.Name, s.Desc, s.SiteName}
		for _, h := range s.Health {
			l := append([]string{h.Subsystem, h.Status}, labels...)

			if h.Subsystem != "vpn" {
				metrics = append(metrics, []*metricExports{
					{u.Site.TxBytesR, prometheus.GaugeValue, h.TxBytesR.Val, l},
					{u.Site.RxBytesR, prometheus.GaugeValue, h.RxBytesR.Val, l},
				}...)

			} else {
				metrics = append(metrics, []*metricExports{
					{u.Site.RemoteUserNumActive, prometheus.CounterValue, h.RemoteUserNumActive.Val, l},
					{u.Site.RemoteUserNumInactive, prometheus.CounterValue, h.RemoteUserNumInactive.Val, l},
					{u.Site.RemoteUserRxBytes, prometheus.CounterValue, h.RemoteUserRxBytes.Val, l},
					{u.Site.RemoteUserTxBytes, prometheus.CounterValue, h.RemoteUserTxBytes.Val, l},
					{u.Site.RemoteUserRxPackets, prometheus.CounterValue, h.RemoteUserRxPackets.Val, l},
					{u.Site.RemoteUserTxPackets, prometheus.CounterValue, h.RemoteUserTxPackets.Val, l},
				}...)
			}

			if h.Subsystem == "lan" || h.Subsystem == "wlan" || h.Subsystem == "wan" {
				metrics = append(metrics, []*metricExports{
					{u.Site.NumAdopted, prometheus.CounterValue, h.NumAdopted.Val, l},
					{u.Site.NumDisconnected, prometheus.CounterValue, h.NumDisconnected.Val, l},
					{u.Site.NumPending, prometheus.CounterValue, h.NumPending.Val, l},
				}...)
			}

			if h.Subsystem == "lan" || h.Subsystem == "wlan" {
				metrics = append(metrics, []*metricExports{
					{u.Site.NumUser, prometheus.CounterValue, h.NumUser.Val, l},
					{u.Site.NumGuest, prometheus.CounterValue, h.NumGuest.Val, l},
					{u.Site.NumIot, prometheus.CounterValue, h.NumIot.Val, l},
				}...)
			}

			if h.Subsystem == "wlan" {
				metrics = append(metrics, []*metricExports{
					{u.Site.NumAp, prometheus.CounterValue, h.NumAp.Val, l},
					{u.Site.NumDisabled, prometheus.CounterValue, h.NumDisabled.Val, l},
				}...)
			}

			if h.Subsystem == "wan" {
				metrics = append(metrics, []*metricExports{
					{u.Site.NumGw, prometheus.CounterValue, h.NumGw.Val, l},
					{u.Site.NumSta, prometheus.CounterValue, h.NumSta.Val, l},
				}...)
			}

			if h.Subsystem == "lan" {
				metrics = append(metrics, []*metricExports{
					{u.Site.NumSw, prometheus.CounterValue, h.NumSw.Val, l},
				}...)
			}

			if h.Subsystem == "www" {
				metrics = append(metrics, []*metricExports{
					{u.Site.Uptime, prometheus.GaugeValue, h.Latency.Val, l},
					{u.Site.Latency, prometheus.GaugeValue, h.Latency.Val, l},
					{u.Site.XputUp, prometheus.GaugeValue, h.XputUp.Val, l},
					{u.Site.XputDown, prometheus.GaugeValue, h.XputDown.Val, l},
					{u.Site.SpeedtestPing, prometheus.GaugeValue, h.SpeedtestPing.Val, l},
					{u.Site.Drops, prometheus.CounterValue, h.Drops.Val, l},
				}...)
			}
		}
		r.ch <- metrics
	}
}
