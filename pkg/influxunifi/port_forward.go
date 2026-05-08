package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchPortForward generates InfluxDB points for a port forwarding rule.
func (u *InfluxUnifi) batchPortForward(r report, pf *unifi.PortForward) {
	if pf == nil {
		return
	}

	tags := map[string]string{
		"site_name": pf.SiteName,
		"source":    pf.SourceName,
		"rule_id":   pf.ID,
		"rule_name": pf.Name,
		"proto":     pf.Proto,
	}

	enabled := 0
	if pf.Enabled.Val {
		enabled = 1
	}

	logged := 0
	if pf.Log.Val {
		logged = 1
	}

	fields := map[string]any{
		"enabled": enabled,
		"logged":  logged,
	}

	r.send(&metric{Table: "port_forward", Tags: tags, Fields: fields})
}

// batchSSLCertificate generates InfluxDB points for an SSL certificate.
func (u *InfluxUnifi) batchSSLCertificate(r report, cert *unifi.SSLCertificate) {
	if cert == nil || cert.ID == "" {
		return
	}

	tags := map[string]string{
		"site_name":   cert.SiteName,
		"cert_id":     cert.ID,
		"cert_type":   cert.CertType,
		"status":      cert.Status,
		"fingerprint": cert.Fingerprint,
	}

	isActive := 0
	if cert.IsActive.Val {
		isActive = 1
	}

	isValid := 0
	if cert.IsValid.Val {
		isValid = 1
	}

	fields := map[string]any{
		"is_active":  isActive,
		"is_valid":   isValid,
		"valid_from": cert.ValidFrom.Val,
		"valid_to":   cert.ValidTo.Val,
		"chain_len":  len(cert.Chain),
	}

	r.send(&metric{Table: "ssl_certificate", Tags: tags, Fields: fields})
}
