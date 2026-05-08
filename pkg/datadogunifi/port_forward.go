package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchPortForward generates PortForward datapoints for Datadog.
func (u *DatadogUnifi) batchPortForward(r report, pf *unifi.PortForward) {
	if pf == nil {
		return
	}

	metricName := metricNamespace("port_forward")

	tags := cleanTags(map[string]string{
		"site_name": pf.SiteName,
		"source":    pf.SourceName,
		"id":        pf.ID,
		"name":      pf.Name,
		"proto":     pf.Proto,
		"fwd_ip":    pf.FwdIP,
		"fwd_port":  pf.FwdPort,
		"dst_port":  pf.DstPort,
		"pf_iface":  pf.PfwdPf,
	})

	_ = r.reportGauge(metricName("enabled"), boolToFloat64(pf.Enabled.Val), tagMapToTags(tags))
	_ = r.reportGauge(metricName("log"), boolToFloat64(pf.Log.Val), tagMapToTags(tags))
}

// batchSSLCertificate generates SSLCertificate datapoints for Datadog.
func (u *DatadogUnifi) batchSSLCertificate(r report, cert *unifi.SSLCertificate) {
	if cert == nil || cert.ID == "" {
		return
	}

	metricName := metricNamespace("ssl_cert")

	tags := cleanTags(map[string]string{
		"site_name":   cert.SiteName,
		"id":          cert.ID,
		"cert_type":   cert.CertType,
		"status":      cert.Status,
		"issuer":      cert.Issuer,
		"subject":     cert.Subject,
		"fingerprint": cert.Fingerprint,
	})

	_ = r.reportGauge(metricName("is_active"), boolToFloat64(cert.IsActive.Val), tagMapToTags(tags))
	_ = r.reportGauge(metricName("is_valid"), boolToFloat64(cert.IsValid.Val), tagMapToTags(tags))
	_ = r.reportGauge(metricName("valid_from"), cert.ValidFrom.Val, tagMapToTags(tags))
	_ = r.reportGauge(metricName("valid_to"), cert.ValidTo.Val, tagMapToTags(tags))
	_ = r.reportGauge(metricName("chain_length"), float64(len(cert.Chain)), tagMapToTags(tags))
}
