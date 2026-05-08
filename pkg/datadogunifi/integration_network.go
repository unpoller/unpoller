package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchWifiBroadcast generates WifiBroadcast (SSID) datapoints for Datadog.
func (u *DatadogUnifi) batchWifiBroadcast(r report, wb *unifi.WifiBroadcast) {
	if wb == nil {
		return
	}

	metricName := metricNamespace("wifi_broadcast")

	tags := cleanTags(map[string]string{
		"site_name":     wb.SiteName,
		"id":            wb.ID,
		"name":          wb.Name,
		"network":       wb.Network,
		"security_type": wb.SecurityConfiguration.Type,
	})

	_ = r.reportGauge(metricName("enabled"), boolToFloat64(wb.Enabled), tagMapToTags(tags))
}

// batchFirewallZone generates FirewallZone datapoints for Datadog.
func (u *DatadogUnifi) batchFirewallZone(r report, fz *unifi.FirewallZone) {
	if fz == nil {
		return
	}

	metricName := metricNamespace("firewall_zone")

	tags := cleanTags(map[string]string{
		"site_name": fz.SiteName,
		"id":        fz.ID,
		"name":      fz.Name,
		"origin":    fz.Metadata.Origin,
	})

	_ = r.reportGauge(metricName("network_count"), float64(len(fz.NetworkIDs)), tagMapToTags(tags))
}

// batchACLRule generates ACLRule datapoints for Datadog.
func (u *DatadogUnifi) batchACLRule(r report, acl *unifi.ACLRule) {
	if acl == nil {
		return
	}

	metricName := metricNamespace("acl_rule")

	tags := cleanTags(map[string]string{
		"site_name":     acl.SiteName,
		"id":            acl.ID,
		"name":          acl.Name,
		"action":        acl.Action,
		"source_filter": acl.SourceFilter,
	})

	_ = r.reportGauge(metricName("enabled"), boolToFloat64(acl.Enabled), tagMapToTags(tags))
	_ = r.reportGauge(metricName("index"), acl.Index.Val, tagMapToTags(tags))
	_ = r.reportGauge(metricName("enforcing_device_count"), float64(len(acl.EnforcingDeviceFilter)), tagMapToTags(tags))
}

// batchVPNServer generates VPNServer datapoints for Datadog.
func (u *DatadogUnifi) batchVPNServer(r report, vs *unifi.VPNServer) {
	if vs == nil {
		return
	}

	metricName := metricNamespace("vpn_server")

	tags := cleanTags(map[string]string{
		"site_name": vs.SiteName,
		"id":        vs.ID,
		"name":      vs.Name,
		"type":      vs.Type,
		"origin":    vs.Metadata.Origin,
	})

	_ = r.reportGauge(metricName("enabled"), boolToFloat64(vs.Enabled), tagMapToTags(tags))
}

// batchSiteToSiteTunnel generates SiteToSiteTunnel datapoints for Datadog.
func (u *DatadogUnifi) batchSiteToSiteTunnel(r report, tun *unifi.SiteToSiteTunnel) {
	if tun == nil {
		return
	}

	metricName := metricNamespace("site_to_site_tunnel")

	tags := cleanTags(map[string]string{
		"site_name": tun.SiteName,
		"id":        tun.ID,
		"name":      tun.Name,
		"type":      tun.Type,
		"origin":    tun.Metadata.Origin,
	})

	// Emit a presence gauge (1.0 = tunnel is configured).
	_ = r.reportGauge(metricName("present"), 1.0, tagMapToTags(tags))
}

// batchLAG generates LAG (link aggregation group) datapoints for Datadog.
func (u *DatadogUnifi) batchLAG(r report, lag *unifi.LAG) {
	if lag == nil {
		return
	}

	metricName := metricNamespace("lag")

	tags := cleanTags(map[string]string{
		"site_name": lag.SiteName,
		"id":        lag.ID,
		"type":      lag.Type,
		"origin":    lag.Metadata.Origin,
	})

	_ = r.reportGauge(metricName("member_count"), float64(len(lag.Members)), tagMapToTags(tags))
}

// batchMCLAGDomain generates MCLAGDomain datapoints for Datadog.
func (u *DatadogUnifi) batchMCLAGDomain(r report, mcd *unifi.MCLAGDomain) {
	if mcd == nil {
		return
	}

	metricName := metricNamespace("mclag_domain")

	tags := cleanTags(map[string]string{
		"site_name": mcd.SiteName,
		"id":        mcd.ID,
		"name":      mcd.Name,
		"origin":    mcd.Metadata.Origin,
	})

	_ = r.reportGauge(metricName("peer_count"), float64(len(mcd.Peers)), tagMapToTags(tags))
	_ = r.reportGauge(metricName("lag_count"), float64(len(mcd.LAGs)), tagMapToTags(tags))
}

// batchSwitchStack generates SwitchStack datapoints for Datadog.
func (u *DatadogUnifi) batchSwitchStack(r report, ss *unifi.SwitchStack) {
	if ss == nil {
		return
	}

	metricName := metricNamespace("switch_stack")

	tags := cleanTags(map[string]string{
		"site_name": ss.SiteName,
		"id":        ss.ID,
		"name":      ss.Name,
		"origin":    ss.Metadata.Origin,
	})

	_ = r.reportGauge(metricName("member_count"), float64(len(ss.Members)), tagMapToTags(tags))
	_ = r.reportGauge(metricName("lag_count"), float64(len(ss.LAGs)), tagMapToTags(tags))
}

// batchDNSPolicy generates DNSPolicy datapoints for Datadog.
func (u *DatadogUnifi) batchDNSPolicy(r report, dp *unifi.DNSPolicy) {
	if dp == nil {
		return
	}

	metricName := metricNamespace("dns_policy")

	tags := cleanTags(map[string]string{
		"site_name": dp.SiteName,
		"id":        dp.ID,
		"type":      dp.Type,
		"domain":    dp.Domain,
	})

	_ = r.reportGauge(metricName("enabled"), boolToFloat64(dp.Enabled), tagMapToTags(tags))
}

// batchRADIUSProfile generates RADIUSProfile datapoints for Datadog.
func (u *DatadogUnifi) batchRADIUSProfile(r report, rp *unifi.RADIUSProfile) {
	if rp == nil {
		return
	}

	metricName := metricNamespace("radius_profile")

	tags := cleanTags(map[string]string{
		"site_name": rp.SiteName,
		"id":        rp.ID,
		"name":      rp.Name,
		"origin":    rp.Metadata.Origin,
	})

	// Emit a presence gauge (1.0 = profile exists).
	_ = r.reportGauge(metricName("present"), 1.0, tagMapToTags(tags))
}

// batchTrafficMatchingList generates TrafficMatchingList datapoints for Datadog.
func (u *DatadogUnifi) batchTrafficMatchingList(r report, tml *unifi.TrafficMatchingList) {
	if tml == nil {
		return
	}

	metricName := metricNamespace("traffic_matching_list")

	tags := cleanTags(map[string]string{
		"site_name": tml.SiteName,
		"id":        tml.ID,
		"name":      tml.Name,
		"type":      tml.Type,
	})

	// Emit a presence gauge (1.0 = list exists).
	_ = r.reportGauge(metricName("present"), 1.0, tagMapToTags(tags))
}

// batchHotspotVoucher generates HotspotVoucher datapoints for Datadog.
func (u *DatadogUnifi) batchHotspotVoucher(r report, hv *unifi.HotspotVoucher) {
	if hv == nil {
		return
	}

	metricName := metricNamespace("hotspot_voucher")

	tags := cleanTags(map[string]string{
		"site_name":  hv.SiteName,
		"id":         hv.ID,
		"name":       hv.Name,
		"expires_at": hv.ExpiresAt,
	})

	_ = r.reportGauge(metricName("authorized_guest_count"), hv.AuthorizedGuestCount.Val, tagMapToTags(tags))
	_ = r.reportGauge(metricName("authorized_guest_limit"), hv.AuthorizedGuestLimit.Val, tagMapToTags(tags))
	_ = r.reportGauge(metricName("data_usage_limit_mbytes"), hv.DataUsageLimitMBytes.Val, tagMapToTags(tags))
	_ = r.reportGauge(metricName("time_limit_minutes"), hv.TimeLimitMinutes.Val, tagMapToTags(tags))
}
