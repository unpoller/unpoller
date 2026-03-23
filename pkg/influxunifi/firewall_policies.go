package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchFirewallPolicy generates a firewall policy datapoint for InfluxDB.
func (u *InfluxUnifi) batchFirewallPolicy(r report, p *unifi.FirewallPolicy) {
	if p == nil {
		return
	}

	tags := map[string]string{
		"rule_name":   p.Name,
		"action":      p.Action,
		"protocol":    p.Protocol,
		"ip_version":  p.IPVersion,
		"source_zone": p.Source.ZoneID,
		"dest_zone":   p.Destination.ZoneID,
		"site_name":   p.SiteName,
		"source":      p.SourceName,
	}

	enabled := 0
	if p.Enabled.Val {
		enabled = 1
	}

	predefined := 0
	if p.Predefined.Val {
		predefined = 1
	}

	logging := 0
	if p.Logging.Val {
		logging = 1
	}

	fields := map[string]any{
		"enabled":    enabled,
		"index":      p.Index.Val,
		"predefined": predefined,
		"logging":    logging,
	}

	r.send(&metric{Table: "firewall_policy", Tags: tags, Fields: fields})
}
