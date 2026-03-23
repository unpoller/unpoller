package otelunifi

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
)

// exportFirewallPolicies emits per-rule and per-site aggregate firewall policy metrics.
func (u *OtelOutput) exportFirewallPolicies(ctx context.Context, meter metric.Meter, m *poller.Metrics, r *Report) {
	type siteKey struct{ site, source string }
	type siteStats struct {
		total    int
		enabled  int
		disabled int
		predef   int
		custom   int
		logging  int
		byAction map[string]int
	}

	sites := make(map[siteKey]*siteStats)

	for _, item := range m.FirewallPolicies {
		p, ok := item.(*unifi.FirewallPolicy)
		if !ok {
			continue
		}

		attrs := attribute.NewSet(
			attribute.String("rule_name", p.Name),
			attribute.String("action", p.Action),
			attribute.String("protocol", p.Protocol),
			attribute.String("ip_version", p.IPVersion),
			attribute.String("source_zone", p.Source.ZoneID),
			attribute.String("dest_zone", p.Destination.ZoneID),
			attribute.String("site_name", p.SiteName),
			attribute.String("source", p.SourceName),
		)

		enabled := 0.0
		if p.Enabled.Val {
			enabled = 1.0
		}

		u.recordGauge(ctx, meter, r, "unifi_firewall_rule_enabled",
			"Firewall rule enabled status (1=enabled, 0=disabled)", enabled, attrs)
		u.recordGauge(ctx, meter, r, "unifi_firewall_rule_index",
			"Firewall rule priority index", p.Index.Val, attrs)

		// Accumulate site-level stats
		key := siteKey{p.SiteName, p.SourceName}
		if _, ok := sites[key]; !ok {
			sites[key] = &siteStats{byAction: make(map[string]int)}
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
	}

	// Emit per-site aggregate metrics
	for key, s := range sites {
		siteAttrs := attribute.NewSet(
			attribute.String("site_name", key.site),
			attribute.String("source", key.source),
		)

		u.recordGauge(ctx, meter, r, "unifi_firewall_rules_total",
			"Total number of firewall rules", float64(s.total), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_firewall_rules_enabled",
			"Number of enabled firewall rules", float64(s.enabled), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_firewall_rules_disabled",
			"Number of disabled firewall rules", float64(s.disabled), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_firewall_rules_predefined",
			"Number of predefined firewall rules", float64(s.predef), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_firewall_rules_custom",
			"Number of custom firewall rules", float64(s.custom), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_firewall_rules_logging_enabled",
			"Number of firewall rules with logging enabled", float64(s.logging), siteAttrs)

		for action, count := range s.byAction {
			actionAttrs := attribute.NewSet(
				attribute.String("action", action),
				attribute.String("site_name", key.site),
				attribute.String("source", key.source),
			)

			u.recordGauge(ctx, meter, r, "unifi_firewall_rules_by_action",
				"Number of firewall rules grouped by action", float64(count), actionAttrs)
		}
	}
}
