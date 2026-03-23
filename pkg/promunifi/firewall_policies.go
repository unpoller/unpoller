package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type firewallpolicy struct {
	// Per-rule metrics
	RuleEnabled *prometheus.Desc
	RuleIndex   *prometheus.Desc
	// Aggregate metrics
	RulesTotal          *prometheus.Desc
	RulesEnabled        *prometheus.Desc
	RulesDisabled       *prometheus.Desc
	RulesByAction       *prometheus.Desc
	RulesPredefined     *prometheus.Desc
	RulesCustom         *prometheus.Desc
	RulesLoggingEnabled *prometheus.Desc
}

func descFirewallPolicy(ns string) *firewallpolicy {
	// Per-rule labels
	ruleLabels := []string{
		"rule_name",
		"action",
		"protocol",
		"ip_version",
		"source_zone",
		"dest_zone",
		"site_name",
		"source",
	}
	// Site-level labels
	siteLabels := []string{"site_name", "source"}
	// Action-level labels
	actionLabels := []string{"action", "site_name", "source"}

	nd := prometheus.NewDesc

	return &firewallpolicy{
		RuleEnabled:         nd(ns+"firewall_rule_enabled", "Firewall rule enabled status (1=enabled, 0=disabled)", ruleLabels, nil),
		RuleIndex:           nd(ns+"firewall_rule_index", "Firewall rule priority index", ruleLabels, nil),
		RulesTotal:          nd(ns+"firewall_rules_total", "Total number of firewall rules", siteLabels, nil),
		RulesEnabled:        nd(ns+"firewall_rules_enabled", "Number of enabled firewall rules", siteLabels, nil),
		RulesDisabled:       nd(ns+"firewall_rules_disabled", "Number of disabled firewall rules", siteLabels, nil),
		RulesByAction:       nd(ns+"firewall_rules_by_action", "Number of firewall rules grouped by action", actionLabels, nil),
		RulesPredefined:     nd(ns+"firewall_rules_predefined", "Number of predefined firewall rules", siteLabels, nil),
		RulesCustom:         nd(ns+"firewall_rules_custom", "Number of custom firewall rules", siteLabels, nil),
		RulesLoggingEnabled: nd(ns+"firewall_rules_logging_enabled", "Number of firewall rules with logging enabled", siteLabels, nil),
	}
}

func (u *promUnifi) exportFirewallPolicies(r report, policies []*unifi.FirewallPolicy) {
	if len(policies) == 0 {
		return
	}

	// Per-site aggregate counters, keyed by "siteName|source"
	type siteKey struct{ site, source string }
	type siteStats struct {
		total    int
		enabled  int
		disabled int
		predef   int
		custom   int
		logging  int
		site     string
		source   string
		byAction map[string]int
	}

	sites := make(map[siteKey]*siteStats)

	for _, p := range policies {
		key := siteKey{p.SiteName, p.SourceName}
		if _, ok := sites[key]; !ok {
			sites[key] = &siteStats{
				site:     p.SiteName,
				source:   p.SourceName,
				byAction: make(map[string]int),
			}
		}

		s := sites[key]
		s.total++

		if p.Enabled.Val {
			s.enabled++
		} else {
			s.disabled++
		}

		if p.Predefined.Val {
			s.predef++
		} else {
			s.custom++
		}

		if p.Logging.Val {
			s.logging++
		}

		if p.Action != "" {
			s.byAction[p.Action]++
		}

		// Per-rule metrics
		ruleLabels := []string{
			p.Name,
			p.Action,
			p.Protocol,
			p.IPVersion,
			p.Source.ZoneID,
			p.Destination.ZoneID,
			p.SiteName,
			p.SourceName,
		}

		enabledVal := 0.0
		if p.Enabled.Val {
			enabledVal = 1.0
		}

		r.send([]*metric{
			{u.FirewallPolicy.RuleEnabled, gauge, enabledVal, ruleLabels},
			{u.FirewallPolicy.RuleIndex, gauge, p.Index.Val, ruleLabels},
		})
	}

	// Site-level aggregate metrics
	for _, s := range sites {
		siteLabels := []string{s.site, s.source}

		r.send([]*metric{
			{u.FirewallPolicy.RulesTotal, gauge, float64(s.total), siteLabels},
			{u.FirewallPolicy.RulesEnabled, gauge, float64(s.enabled), siteLabels},
			{u.FirewallPolicy.RulesDisabled, gauge, float64(s.disabled), siteLabels},
			{u.FirewallPolicy.RulesPredefined, gauge, float64(s.predef), siteLabels},
			{u.FirewallPolicy.RulesCustom, gauge, float64(s.custom), siteLabels},
			{u.FirewallPolicy.RulesLoggingEnabled, gauge, float64(s.logging), siteLabels},
		})

		for action, count := range s.byAction {
			r.send([]*metric{
				{u.FirewallPolicy.RulesByAction, gauge, float64(count), []string{action, s.site, s.source}},
			})
		}
	}
}
