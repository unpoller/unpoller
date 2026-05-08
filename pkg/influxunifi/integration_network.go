package influxunifi

import (
	"fmt"

	"github.com/unpoller/unifi/v5"
)

// batchWifiBroadcast generates InfluxDB points for a WiFi SSID broadcast configuration.
func (u *InfluxUnifi) batchWifiBroadcast(r report, wb *unifi.WifiBroadcast) {
	if wb == nil {
		return
	}

	tags := map[string]string{
		"site_name":      wb.SiteName,
		"broadcast_id":   wb.ID,
		"broadcast_name": wb.Name,
		"network":        wb.Network,
		"security_type":  wb.SecurityConfiguration.Type,
	}

	enabled := 0
	if wb.Enabled {
		enabled = 1
	}

	fields := map[string]any{
		"enabled":                   enabled,
		"broadcasting_device_count": len(wb.BroadcastingDeviceFilter),
	}

	r.send(&metric{Table: "wifi_broadcast", Tags: tags, Fields: fields})
}

// batchFirewallZone generates InfluxDB points for a firewall zone.
func (u *InfluxUnifi) batchFirewallZone(r report, fz *unifi.FirewallZone) {
	if fz == nil {
		return
	}

	tags := map[string]string{
		"site_name": fz.SiteName,
		"zone_id":   fz.ID,
		"zone_name": fz.Name,
		"origin":    fz.Metadata.Origin,
	}

	fields := map[string]any{
		"network_count": len(fz.NetworkIDs),
	}

	r.send(&metric{Table: "firewall_zone", Tags: tags, Fields: fields})
}

// batchACLRule generates InfluxDB points for an ACL rule.
func (u *InfluxUnifi) batchACLRule(r report, rule *unifi.ACLRule) {
	if rule == nil {
		return
	}

	tags := map[string]string{
		"site_name":     rule.SiteName,
		"rule_id":       rule.ID,
		"rule_name":     rule.Name,
		"action":        rule.Action,
		"source_filter": rule.SourceFilter,
	}

	enabled := 0
	if rule.Enabled {
		enabled = 1
	}

	fields := map[string]any{
		"enabled":                enabled,
		"index":                  rule.Index.Val,
		"enforcing_device_count": len(rule.EnforcingDeviceFilter),
	}

	r.send(&metric{Table: "acl_rule", Tags: tags, Fields: fields})
}

// batchVPNServer generates InfluxDB points for a VPN server configuration.
func (u *InfluxUnifi) batchVPNServer(r report, vs *unifi.VPNServer) {
	if vs == nil {
		return
	}

	tags := map[string]string{
		"site_name":   vs.SiteName,
		"server_id":   vs.ID,
		"server_name": vs.Name,
		"vpn_type":    vs.Type,
		"origin":      vs.Metadata.Origin,
	}

	enabled := 0
	if vs.Enabled {
		enabled = 1
	}

	fields := map[string]any{
		"enabled": enabled,
	}

	r.send(&metric{Table: "vpn_server", Tags: tags, Fields: fields})
}

// batchSiteToSiteTunnel generates InfluxDB points for a site-to-site VPN tunnel.
func (u *InfluxUnifi) batchSiteToSiteTunnel(r report, t *unifi.SiteToSiteTunnel) {
	if t == nil {
		return
	}

	tags := map[string]string{
		"site_name":   t.SiteName,
		"tunnel_id":   t.ID,
		"tunnel_name": t.Name,
		"tunnel_type": t.Type,
		"origin":      t.Metadata.Origin,
	}

	r.send(&metric{
		Table:  "site_to_site_tunnel",
		Tags:   tags,
		Fields: map[string]any{"present": 1},
	})
}

// batchLAG generates InfluxDB points for a link aggregation group.
func (u *InfluxUnifi) batchLAG(r report, lag *unifi.LAG) {
	if lag == nil {
		return
	}

	tags := map[string]string{
		"site_name": lag.SiteName,
		"lag_id":    lag.ID,
		"lag_type":  lag.Type,
		"origin":    lag.Metadata.Origin,
	}

	totalPorts := 0

	for _, m := range lag.Members {
		totalPorts += len(m.PortIndexes)
	}

	fields := map[string]any{
		"member_count": len(lag.Members),
		"port_count":   totalPorts,
	}

	r.send(&metric{Table: "lag", Tags: tags, Fields: fields})
}

// batchMCLAGDomain generates InfluxDB points for a multi-chassis LAG domain.
func (u *InfluxUnifi) batchMCLAGDomain(r report, d *unifi.MCLAGDomain) {
	if d == nil {
		return
	}

	tags := map[string]string{
		"site_name":   d.SiteName,
		"domain_id":   d.ID,
		"domain_name": d.Name,
		"origin":      d.Metadata.Origin,
	}

	fields := map[string]any{
		"peer_count": len(d.Peers),
		"lag_count":  len(d.LAGs),
	}

	r.send(&metric{Table: "mclag_domain", Tags: tags, Fields: fields})

	for i, peer := range d.Peers {
		peerTags := map[string]string{
			"site_name":   d.SiteName,
			"domain_id":   d.ID,
			"domain_name": d.Name,
			"peer_index":  fmt.Sprint(i),
			"device_id":   peer.DeviceID,
			"role":        peer.Role,
		}

		r.send(&metric{
			Table:  "mclag_peer",
			Tags:   peerTags,
			Fields: map[string]any{"link_port_count": len(peer.LinkPorts)},
		})
	}
}

// batchSwitchStack generates InfluxDB points for a switch stack.
func (u *InfluxUnifi) batchSwitchStack(r report, ss *unifi.SwitchStack) {
	if ss == nil {
		return
	}

	tags := map[string]string{
		"site_name":  ss.SiteName,
		"stack_id":   ss.ID,
		"stack_name": ss.Name,
		"origin":     ss.Metadata.Origin,
	}

	fields := map[string]any{
		"member_count": len(ss.Members),
		"lag_count":    len(ss.LAGs),
	}

	r.send(&metric{Table: "switch_stack", Tags: tags, Fields: fields})
}

// batchDNSPolicy generates InfluxDB points for a DNS policy.
func (u *InfluxUnifi) batchDNSPolicy(r report, p *unifi.DNSPolicy) {
	if p == nil {
		return
	}

	tags := map[string]string{
		"site_name":   p.SiteName,
		"policy_id":   p.ID,
		"policy_type": p.Type,
		"domain":      p.Domain,
	}

	enabled := 0
	if p.Enabled {
		enabled = 1
	}

	fields := map[string]any{
		"enabled": enabled,
	}

	r.send(&metric{Table: "dns_policy", Tags: tags, Fields: fields})
}

// batchRADIUSProfile generates InfluxDB points for a RADIUS profile.
func (u *InfluxUnifi) batchRADIUSProfile(r report, p *unifi.RADIUSProfile) {
	if p == nil {
		return
	}

	tags := map[string]string{
		"site_name":    p.SiteName,
		"profile_id":   p.ID,
		"profile_name": p.Name,
		"origin":       p.Metadata.Origin,
	}

	r.send(&metric{
		Table:  "radius_profile",
		Tags:   tags,
		Fields: map[string]any{"present": 1},
	})
}

// batchTrafficMatchingList generates InfluxDB points for a traffic matching list.
func (u *InfluxUnifi) batchTrafficMatchingList(r report, l *unifi.TrafficMatchingList) {
	if l == nil {
		return
	}

	tags := map[string]string{
		"site_name": l.SiteName,
		"list_id":   l.ID,
		"list_name": l.Name,
		"list_type": l.Type,
	}

	r.send(&metric{
		Table:  "traffic_matching_list",
		Tags:   tags,
		Fields: map[string]any{"present": 1},
	})
}

// batchHotspotVoucher generates InfluxDB points for a hotspot voucher.
func (u *InfluxUnifi) batchHotspotVoucher(r report, v *unifi.HotspotVoucher) {
	if v == nil {
		return
	}

	tags := map[string]string{
		"site_name":    v.SiteName,
		"voucher_id":   v.ID,
		"voucher_name": v.Name,
	}

	fields := map[string]any{
		"authorized_guest_count":  v.AuthorizedGuestCount.Val,
		"authorized_guest_limit":  v.AuthorizedGuestLimit.Val,
		"data_usage_limit_mbytes": v.DataUsageLimitMBytes.Val,
		"time_limit_minutes":      v.TimeLimitMinutes.Val,
	}

	r.send(&metric{Table: "hotspot_voucher", Tags: tags, Fields: fields})
}
