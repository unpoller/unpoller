package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type portForward struct {
	Enabled *prometheus.Desc
	Log     *prometheus.Desc
}

func descPortForward(ns string) *portForward {
	labels := []string{"site_name", "name", "proto", "fwd_ip", "fwd_port", "dst_port"}

	return &portForward{
		Enabled: prometheus.NewDesc(ns+"port_forward_enabled",
			"Port forward rule enabled (1=enabled, 0=disabled)",
			labels, nil),
		Log: prometheus.NewDesc(ns+"port_forward_log",
			"Port forward rule logging enabled (1=on, 0=off)",
			labels, nil),
	}
}

func (u *promUnifi) exportPortForward(r report, pf *unifi.PortForward) {
	if pf == nil {
		return
	}

	labels := []string{pf.SiteName, pf.Name, pf.Proto, pf.FwdIP, pf.FwdPort, pf.DstPort}

	r.send([]*metric{
		{u.PortForward.Enabled, gauge, pf.Enabled.Val, labels},
		{u.PortForward.Log, gauge, pf.Log.Val, labels},
	})
}

// sslCertificate holds Prometheus descriptors for SSL certificate metrics.
type sslCertificate struct {
	IsActive  *prometheus.Desc
	IsValid   *prometheus.Desc
	ValidFrom *prometheus.Desc
	ValidTo   *prometheus.Desc
}

func descSSLCertificate(ns string) *sslCertificate {
	labels := []string{"site_name", "cert_type", "subject", "issuer", "status"}

	return &sslCertificate{
		IsActive: prometheus.NewDesc(ns+"ssl_cert_active",
			"SSL certificate is the active certificate (1=active, 0=inactive)",
			labels, nil),
		IsValid: prometheus.NewDesc(ns+"ssl_cert_valid",
			"SSL certificate passes validity checks (1=valid, 0=invalid)",
			labels, nil),
		ValidFrom: prometheus.NewDesc(ns+"ssl_cert_valid_from_seconds",
			"SSL certificate validity start time (Unix epoch)",
			labels, nil),
		ValidTo: prometheus.NewDesc(ns+"ssl_cert_valid_to_seconds",
			"SSL certificate expiry time (Unix epoch)",
			labels, nil),
	}
}

func (u *promUnifi) exportSSLCertificate(r report, cert *unifi.SSLCertificate) {
	if cert == nil || cert.ID == "" {
		return
	}

	labels := []string{cert.SiteName, cert.CertType, cert.Subject, cert.Issuer, cert.Status}

	r.send([]*metric{
		{u.SSLCertificate.IsActive, gauge, cert.IsActive.Val, labels},
		{u.SSLCertificate.IsValid, gauge, cert.IsValid.Val, labels},
		{u.SSLCertificate.ValidFrom, gauge, cert.ValidFrom, labels},
		{u.SSLCertificate.ValidTo, gauge, cert.ValidTo, labels},
	})
}
