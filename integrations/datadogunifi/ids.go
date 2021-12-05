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
		tag("usgip", i.Usgip),
		tag("country_code", i.SrcipGeo.CountryCode),
		tag("country_name", i.SrcipGeo.CountryName),
		tag("city", i.SrcipGeo.City),
		tag("srcipASN", i.SrcipASN),
		tag("usgipASN", i.UsgipASN),
		tag("alert_category", i.InnerAlertCategory),
		tag("subsystem", i.Subsystem),
		tag("catname", i.Catname),
	}

	metricName := metricNamespace("intrusion_detect")
	r.reportCount(metricName("count"), 1, tags)
}
