package datadogunifi

import (
	"github.com/unpoller/unifi"
)

// reportIDS generates intrusion detection datapoints for Datadog.
// These points can be passed directly to datadog.
func (u *DatadogUnifi) reportIDS(r report, i *unifi.IDS) {
	tags := []string{
		tag("site_name", i.SiteName),
		tag("source", i.SourceName),
		tag("in_iface", i.InIface),
		tag("event_type", i.EventType),
		tag("proto", i.Proto),
		tag("app_proto", i.AppProto),
		tag("usg_ip", i.USGIP),
		tag("country_code", i.SourceIPGeo.CountryCode),
		tag("country_name", i.SourceIPGeo.CountryName),
		tag("city", i.SourceIPGeo.City),
		tag("src_ip_ASN", i.SrcIPASN),
		tag("usg_ip_ASN", i.USGIPASN),
		tag("alert_category", i.InnerAlertCategory),
		tag("subsystem", i.Subsystem),
		tag("catname", i.Catname),
	}

	metricName := metricNamespace("intrusion_detect")
	r.reportCount(metricName("count"), 1, tags)
}
