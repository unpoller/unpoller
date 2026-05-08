package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

// wifiBroadcast holds Prometheus descriptors for WiFi SSID broadcast metrics.
type wifiBroadcast struct {
	Enabled *prometheus.Desc
}

func descWifiBroadcast(ns string) *wifiBroadcast {
	labels := []string{"site_name", "name", "network", "security_type"}

	return &wifiBroadcast{
		Enabled: prometheus.NewDesc(ns+"wifi_broadcast_enabled",
			"WiFi SSID broadcast enabled (1=enabled, 0=disabled)",
			labels, nil),
	}
}

func (u *promUnifi) exportWifiBroadcast(r report, wb *unifi.WifiBroadcast) {
	if wb == nil {
		return
	}

	enabled := 0.0
	if wb.Enabled {
		enabled = 1.0
	}

	labels := []string{wb.SiteName, wb.Name, wb.Network, wb.SecurityConfiguration.Type}

	r.send([]*metric{
		{u.WifiBroadcast.Enabled, gauge, enabled, labels},
	})
}

// firewallZone holds Prometheus descriptors for firewall zone metrics.
type firewallZone struct {
	NetworkCount *prometheus.Desc
}

func descFirewallZone(ns string) *firewallZone {
	labels := []string{"site_name", "name", "origin"}

	return &firewallZone{
		NetworkCount: prometheus.NewDesc(ns+"firewall_zone_network_count",
			"Number of networks assigned to the firewall zone",
			labels, nil),
	}
}

func (u *promUnifi) exportFirewallZone(r report, fz *unifi.FirewallZone) {
	if fz == nil {
		return
	}

	labels := []string{fz.SiteName, fz.Name, fz.Metadata.Origin}

	r.send([]*metric{
		{u.FirewallZone.NetworkCount, gauge, float64(len(fz.NetworkIDs)), labels},
	})
}

// aclRule holds Prometheus descriptors for ACL rule metrics.
type aclRule struct {
	Enabled *prometheus.Desc
	Index   *prometheus.Desc
}

func descACLRule(ns string) *aclRule {
	labels := []string{"site_name", "name", "action"}

	return &aclRule{
		Enabled: prometheus.NewDesc(ns+"acl_rule_enabled",
			"ACL rule enabled (1=enabled, 0=disabled)",
			labels, nil),
		Index: prometheus.NewDesc(ns+"acl_rule_index",
			"ACL rule evaluation order index",
			labels, nil),
	}
}

func (u *promUnifi) exportACLRule(r report, ar *unifi.ACLRule) {
	if ar == nil {
		return
	}

	enabled := 0.0
	if ar.Enabled {
		enabled = 1.0
	}

	labels := []string{ar.SiteName, ar.Name, ar.Action}

	r.send([]*metric{
		{u.ACLRule.Enabled, gauge, enabled, labels},
		{u.ACLRule.Index, gauge, ar.Index, labels},
	})
}

// vpnServer holds Prometheus descriptors for VPN server metrics.
type vpnServer struct {
	Enabled *prometheus.Desc
}

func descVPNServer(ns string) *vpnServer {
	labels := []string{"site_name", "name", "vpn_type", "origin"}

	return &vpnServer{
		Enabled: prometheus.NewDesc(ns+"vpn_server_enabled",
			"VPN server enabled (1=enabled, 0=disabled)",
			labels, nil),
	}
}

func (u *promUnifi) exportVPNServer(r report, vs *unifi.VPNServer) {
	if vs == nil {
		return
	}

	enabled := 0.0
	if vs.Enabled {
		enabled = 1.0
	}

	labels := []string{vs.SiteName, vs.Name, vs.Type, vs.Metadata.Origin}

	r.send([]*metric{
		{u.VPNServer.Enabled, gauge, enabled, labels},
	})
}

// siteToSiteTunnel holds Prometheus descriptors for site-to-site VPN tunnel metrics.
type siteToSiteTunnel struct {
	Presence *prometheus.Desc
}

func descSiteToSiteTunnel(ns string) *siteToSiteTunnel {
	labels := []string{"site_name", "name", "tunnel_type", "origin"}

	return &siteToSiteTunnel{
		Presence: prometheus.NewDesc(ns+"site_to_site_tunnel_present",
			"Site-to-site VPN tunnel configured (always 1 when present)",
			labels, nil),
	}
}

func (u *promUnifi) exportSiteToSiteTunnel(r report, t *unifi.SiteToSiteTunnel) {
	if t == nil {
		return
	}

	labels := []string{t.SiteName, t.Name, t.Type, t.Metadata.Origin}

	r.send([]*metric{
		{u.SiteToSiteTunnel.Presence, gauge, 1.0, labels},
	})
}

// lag holds Prometheus descriptors for Link Aggregation Group metrics.
type lag struct {
	MemberCount *prometheus.Desc
}

func descLAG(ns string) *lag {
	labels := []string{"site_name", "lag_id", "lag_type", "origin"}

	return &lag{
		MemberCount: prometheus.NewDesc(ns+"lag_member_count",
			"Number of member entries in the LAG",
			labels, nil),
	}
}

func (u *promUnifi) exportLAG(r report, l *unifi.LAG) {
	if l == nil {
		return
	}

	labels := []string{l.SiteName, l.ID, l.Type, l.Metadata.Origin}

	r.send([]*metric{
		{u.LAG.MemberCount, gauge, float64(len(l.Members)), labels},
	})
}

// mclagDomain holds Prometheus descriptors for MC-LAG domain metrics.
type mclagDomain struct {
	LAGCount  *prometheus.Desc
	PeerCount *prometheus.Desc
}

func descMCLAGDomain(ns string) *mclagDomain {
	labels := []string{"site_name", "name", "origin"}

	return &mclagDomain{
		LAGCount: prometheus.NewDesc(ns+"mclag_domain_lag_count",
			"Number of LAGs in the MC-LAG domain",
			labels, nil),
		PeerCount: prometheus.NewDesc(ns+"mclag_domain_peer_count",
			"Number of peer devices in the MC-LAG domain",
			labels, nil),
	}
}

func (u *promUnifi) exportMCLAGDomain(r report, d *unifi.MCLAGDomain) {
	if d == nil {
		return
	}

	labels := []string{d.SiteName, d.Name, d.Metadata.Origin}

	r.send([]*metric{
		{u.MCLAGDomain.LAGCount, gauge, float64(len(d.LAGs)), labels},
		{u.MCLAGDomain.PeerCount, gauge, float64(len(d.Peers)), labels},
	})
}

// switchStack holds Prometheus descriptors for switch stack metrics.
type switchStack struct {
	MemberCount *prometheus.Desc
}

func descSwitchStack(ns string) *switchStack {
	labels := []string{"site_name", "name", "origin"}

	return &switchStack{
		MemberCount: prometheus.NewDesc(ns+"switch_stack_member_count",
			"Number of member devices in the switch stack",
			labels, nil),
	}
}

func (u *promUnifi) exportSwitchStack(r report, s *unifi.SwitchStack) {
	if s == nil {
		return
	}

	labels := []string{s.SiteName, s.Name, s.Metadata.Origin}

	r.send([]*metric{
		{u.SwitchStack.MemberCount, gauge, float64(len(s.Members)), labels},
	})
}

// dnsPolicy holds Prometheus descriptors for DNS policy metrics.
type dnsPolicy struct {
	Enabled *prometheus.Desc
}

func descDNSPolicy(ns string) *dnsPolicy {
	labels := []string{"site_name", "domain", "policy_type"}

	return &dnsPolicy{
		Enabled: prometheus.NewDesc(ns+"dns_policy_enabled",
			"DNS policy enabled (1=enabled, 0=disabled)",
			labels, nil),
	}
}

func (u *promUnifi) exportDNSPolicy(r report, dp *unifi.DNSPolicy) {
	if dp == nil {
		return
	}

	enabled := 0.0
	if dp.Enabled {
		enabled = 1.0
	}

	labels := []string{dp.SiteName, dp.Domain, dp.Type}

	r.send([]*metric{
		{u.DNSPolicy.Enabled, gauge, enabled, labels},
	})
}

// radiusProfile holds Prometheus descriptors for RADIUS profile metrics.
type radiusProfile struct {
	Presence *prometheus.Desc
}

func descRADIUSProfile(ns string) *radiusProfile {
	labels := []string{"site_name", "name", "origin"}

	return &radiusProfile{
		Presence: prometheus.NewDesc(ns+"radius_profile_present",
			"RADIUS profile configured (always 1 when present)",
			labels, nil),
	}
}

func (u *promUnifi) exportRADIUSProfile(r report, rp *unifi.RADIUSProfile) {
	if rp == nil {
		return
	}

	labels := []string{rp.SiteName, rp.Name, rp.Metadata.Origin}

	r.send([]*metric{
		{u.RADIUSProfile.Presence, gauge, 1.0, labels},
	})
}

// trafficMatchingList holds Prometheus descriptors for traffic matching list metrics.
type trafficMatchingList struct {
	Presence *prometheus.Desc
}

func descTrafficMatchingList(ns string) *trafficMatchingList {
	labels := []string{"site_name", "name", "list_type"}

	return &trafficMatchingList{
		Presence: prometheus.NewDesc(ns+"traffic_matching_list_present",
			"Traffic matching list configured (always 1 when present)",
			labels, nil),
	}
}

func (u *promUnifi) exportTrafficMatchingList(r report, tml *unifi.TrafficMatchingList) {
	if tml == nil {
		return
	}

	labels := []string{tml.SiteName, tml.Name, tml.Type}

	r.send([]*metric{
		{u.TrafficMatchingList.Presence, gauge, 1.0, labels},
	})
}

// hotspotVoucher holds Prometheus descriptors for hotspot voucher metrics.
type hotspotVoucher struct {
	AuthorizedGuestCount *prometheus.Desc
	AuthorizedGuestLimit *prometheus.Desc
	DataUsageLimitMBytes *prometheus.Desc
	TimeLimitMinutes     *prometheus.Desc
}

func descHotspotVoucher(ns string) *hotspotVoucher {
	labels := []string{"site_name", "name", "code"}

	return &hotspotVoucher{
		AuthorizedGuestCount: prometheus.NewDesc(ns+"hotspot_voucher_authorized_guest_count",
			"Number of guests currently authorized with this voucher",
			labels, nil),
		AuthorizedGuestLimit: prometheus.NewDesc(ns+"hotspot_voucher_authorized_guest_limit",
			"Maximum number of guests allowed with this voucher (0=unlimited)",
			labels, nil),
		DataUsageLimitMBytes: prometheus.NewDesc(ns+"hotspot_voucher_data_usage_limit_mbytes",
			"Data usage cap for the voucher in MBytes (0=no limit)",
			labels, nil),
		TimeLimitMinutes: prometheus.NewDesc(ns+"hotspot_voucher_time_limit_minutes",
			"Time limit for the voucher in minutes (0=no limit)",
			labels, nil),
	}
}

func (u *promUnifi) exportHotspotVoucher(r report, hv *unifi.HotspotVoucher) {
	if hv == nil {
		return
	}

	labels := []string{hv.SiteName, hv.Name, hv.Code}

	r.send([]*metric{
		{u.HotspotVoucher.AuthorizedGuestCount, gauge, hv.AuthorizedGuestCount, labels},
		{u.HotspotVoucher.AuthorizedGuestLimit, gauge, hv.AuthorizedGuestLimit, labels},
		{u.HotspotVoucher.DataUsageLimitMBytes, gauge, hv.DataUsageLimitMBytes, labels},
		{u.HotspotVoucher.TimeLimitMinutes, gauge, hv.TimeLimitMinutes, labels},
	})
}
