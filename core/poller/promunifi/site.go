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

// XXX: The help values can be more verbose.
func descSite(ns string) *site {
	if ns += "_site_"; ns == "_site_" {
		ns = "site_"
	}

	labels := []string{"name", "desc", "site_name", "subsystem", "status", "gwversion"}

	return &site{
		NumUser:               prometheus.NewDesc(ns+"NumUser", "NumUser", labels, nil),
		NumGuest:              prometheus.NewDesc(ns+"NumGuest", "NumGuest", labels, nil),
		NumIot:                prometheus.NewDesc(ns+"NumIot", "NumIot", labels, nil),
		TxBytesR:              prometheus.NewDesc(ns+"TxBytesR", "TxBytesR", labels, nil),
		RxBytesR:              prometheus.NewDesc(ns+"RxBytesR", "RxBytesR", labels, nil),
		NumAp:                 prometheus.NewDesc(ns+"NumAp", "NumAp", labels, nil),
		NumAdopted:            prometheus.NewDesc(ns+"NumAdopted", "NumAdopted", labels, nil),
		NumDisabled:           prometheus.NewDesc(ns+"NumDisabled", "NumDisabled", labels, nil),
		NumDisconnected:       prometheus.NewDesc(ns+"NumDisconnected", "NumDisconnected", labels, nil),
		NumPending:            prometheus.NewDesc(ns+"NumPending", "NumPending", labels, nil),
		NumGw:                 prometheus.NewDesc(ns+"NumGw", "NumGw", labels, nil),
		NumSw:                 prometheus.NewDesc(ns+"NumSw", "NumSw", labels, nil),
		NumSta:                prometheus.NewDesc(ns+"NumSta", "NumSta", labels, nil),
		Latency:               prometheus.NewDesc(ns+"Latency", "Latency", labels, nil),
		Drops:                 prometheus.NewDesc(ns+"Drops", "Drops", labels, nil),
		XputUp:                prometheus.NewDesc(ns+"XputUp", "XputUp", labels, nil),
		XputDown:              prometheus.NewDesc(ns+"XputDown", "XputDown", labels, nil),
		SpeedtestPing:         prometheus.NewDesc(ns+"SpeedtestPing", "SpeedtestPing", labels, nil),
		RemoteUserNumActive:   prometheus.NewDesc(ns+"RemoteUserNumActive", "RemoteUserNumActive", labels, nil),
		RemoteUserNumInactive: prometheus.NewDesc(ns+"RemoteUserNumInactive", "RemoteUserNumInactive", labels, nil),
		RemoteUserRxBytes:     prometheus.NewDesc(ns+"RemoteUserRxBytes", "RemoteUserRxBytes", labels, nil),
		RemoteUserTxBytes:     prometheus.NewDesc(ns+"RemoteUserTxBytes", "RemoteUserTxBytes", labels, nil),
		RemoteUserRxPackets:   prometheus.NewDesc(ns+"RemoteUserRxPackets", "RemoteUserRxPackets", labels, nil),
		RemoteUserTxPackets:   prometheus.NewDesc(ns+"RemoteUserTxPackets", "RemoteUserTxPackets", labels, nil),
	}
}

// exportSite exports Network Site Data
func (u *unifiCollector) exportSite(s *unifi.Site) []*metricExports {
	labels := []string{s.Name, s.Desc, s.SiteName}
	var m []*metricExports
	for _, h := range s.Health {
		l := append(labels, h.Subsystem, h.Status, h.GwVersion)
		m = append(m, &metricExports{u.Site.NumUser, prometheus.CounterValue, h.NumUser.Val, l})
		m = append(m, &metricExports{u.Site.NumGuest, prometheus.CounterValue, h.NumGuest.Val, l})
		m = append(m, &metricExports{u.Site.NumIot, prometheus.CounterValue, h.NumIot.Val, l})
		m = append(m, &metricExports{u.Site.TxBytesR, prometheus.GaugeValue, h.TxBytesR.Val, l})
		m = append(m, &metricExports{u.Site.RxBytesR, prometheus.GaugeValue, h.RxBytesR.Val, l})
		m = append(m, &metricExports{u.Site.NumAp, prometheus.CounterValue, h.NumAp.Val, l})
		m = append(m, &metricExports{u.Site.NumAdopted, prometheus.CounterValue, h.NumAdopted.Val, l})
		m = append(m, &metricExports{u.Site.NumDisabled, prometheus.CounterValue, h.NumDisabled.Val, l})
		m = append(m, &metricExports{u.Site.NumDisconnected, prometheus.CounterValue, h.NumDisconnected.Val, l})
		m = append(m, &metricExports{u.Site.NumPending, prometheus.CounterValue, h.NumPending.Val, l})
		m = append(m, &metricExports{u.Site.NumGw, prometheus.CounterValue, h.NumGw.Val, l})
		m = append(m, &metricExports{u.Site.NumSw, prometheus.CounterValue, h.NumSw.Val, l})
		m = append(m, &metricExports{u.Site.NumSta, prometheus.CounterValue, h.NumSta.Val, l})
		m = append(m, &metricExports{u.Site.Latency, prometheus.GaugeValue, h.Latency.Val, l})
		m = append(m, &metricExports{u.Site.Drops, prometheus.CounterValue, h.Drops.Val, l})
		m = append(m, &metricExports{u.Site.XputUp, prometheus.GaugeValue, h.XputUp.Val, l})
		m = append(m, &metricExports{u.Site.XputDown, prometheus.GaugeValue, h.XputDown.Val, l})
		m = append(m, &metricExports{u.Site.SpeedtestPing, prometheus.GaugeValue, h.SpeedtestPing.Val, l})
		if h.Subsystem == "vpn" {
			m = append(m, &metricExports{u.Site.RemoteUserNumActive, prometheus.CounterValue, h.RemoteUserNumActive.Val, l})
			m = append(m, &metricExports{u.Site.RemoteUserNumInactive, prometheus.CounterValue, h.RemoteUserNumInactive.Val, l})
			m = append(m, &metricExports{u.Site.RemoteUserRxBytes, prometheus.CounterValue, h.RemoteUserRxBytes.Val, l})
			m = append(m, &metricExports{u.Site.RemoteUserTxBytes, prometheus.CounterValue, h.RemoteUserTxBytes.Val, l})
			m = append(m, &metricExports{u.Site.RemoteUserRxPackets, prometheus.CounterValue, h.RemoteUserRxPackets.Val, l})
			m = append(m, &metricExports{u.Site.RemoteUserTxPackets, prometheus.CounterValue, h.RemoteUserTxPackets.Val, l})
		}
	}
	return m
}
