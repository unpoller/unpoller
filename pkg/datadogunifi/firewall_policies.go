package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchFirewallPolicy generates firewall policy datapoints for Datadog.
func (u *DatadogUnifi) batchFirewallPolicy(r report, p *unifi.FirewallPolicy) {
	if p == nil {
		return
	}

	metricName := metricNamespace("firewall_policy")

	tags := []string{
		tag("rule_name", p.Name),
		tag("action", p.Action),
		tag("protocol", p.Protocol),
		tag("ip_version", p.IPVersion),
		tag("source_zone", p.Source.ZoneID),
		tag("dest_zone", p.Destination.ZoneID),
		tag("site_name", p.SiteName),
		tag("source", p.SourceName),
	}

	enabled := 0.0
	if p.Enabled.Val {
		enabled = 1.0
	}

	predefined := 0.0
	if p.Predefined.Val {
		predefined = 1.0
	}

	logging := 0.0
	if p.Logging.Val {
		logging = 1.0
	}

	data := map[string]float64{
		"enabled":    enabled,
		"index":      p.Index.Val,
		"predefined": predefined,
		"logging":    logging,
	}

	for name, value := range data {
		_ = r.reportGauge(metricName(name), value, tags)
	}
}
