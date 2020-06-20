package influxunifi

import (
	"github.com/unifi-poller/unifi"
)

// batchIDS generates intrusion detection datapoints for InfluxDB.
func (u *InfluxUnifi) batchIDS(r report, i *unifi.IDS) { // nolint: funlen
	fields := map[string]interface{}{
		/*
			"site_id":                  i.SiteID,
			"dstipASN":                 i.DstIPASN,
			"dstipCountry":             i.DstIPCountry,
			"flow_id":                  i.FlowID,
			"inner_alert_action":       i.InnerAlertAction,
			"inner_alert_category":     i.InnerAlertCategory,
			"inner_alert_signature":    i.InnerAlertSignature,
			"inner_alert_rev":          i.InnerAlertRev,
			"inner_alert_severity":     i.InnerAlertSeverity,
			"inner_alert_gid":          i.InnerAlertGID,
			"inner_alert_signature_id": i.InnerAlertSignatureID,
			"srcipASN":                 i.SrcIPASN,
			"srcipCountry":             i.SrcIPCountry,
			"unique_alertid":           i.UniqueAlertID,
			"usgipASN":                 i.UsgIPASN,
			"usgipCountry":             i.UsgIPCountry,
		*/
		"dest_port":               i.DestPort,
		"src_port":                i.SrcPort,
		"app_proto":               i.AppProto,
		"catname":                 i.Catname,
		"dest_ip":                 i.DestIP,
		"dst_mac":                 i.DstMAC,
		"host":                    i.Host,
		"key":                     i.Key,
		"msg":                     i.Msg,
		"proto":                   i.Proto,
		"src_ip":                  i.SrcIP,
		"src_mac":                 i.SrcMAC,
		"usgip":                   i.USGIP,
		"dstipGeo_asn":            i.DestIPGeo.Asn,
		"dstipGeo_latitude":       i.DestIPGeo.Latitude,
		"dstipGeo_longitude":      i.DestIPGeo.Longitude,
		"dstipGeo_city":           i.DestIPGeo.City,
		"dstipGeo_continent_code": i.DestIPGeo.ContinentCode,
		"dstipGeo_country_code":   i.DestIPGeo.CountryCode,
		"dstipGeo_country_name":   i.DestIPGeo.CountryName,
		"dstipGeo_organization":   i.DestIPGeo.Organization,
		"srcipGeo_asn":            i.SourceIPGeo.Asn,
		"srcipGeo_latitude":       i.SourceIPGeo.Latitude,
		"srcipGeo_longitude":      i.SourceIPGeo.Longitude,
		"srcipGeo_city":           i.SourceIPGeo.City,
		"srcipGeo_continent_code": i.SourceIPGeo.ContinentCode,
		"srcipGeo_country_code":   i.SourceIPGeo.CountryCode,
		"srcipGeo_country_name":   i.SourceIPGeo.CountryName,
		"srcipGeo_organization":   i.SourceIPGeo.Organization,
	}

	r.send(&metric{
		Table:  "unifi_ids",
		TS:     i.Datetime,
		Fields: cleanFields(fields),
		Tags: cleanTags(map[string]string{
			"site_name":  i.SiteName,
			"source":     i.SourceName,
			"in_iface":   i.InIface,
			"event_type": i.EventType,
			"subsystem":  i.Subsystem,
			"archived":   i.Archived.Txt,
		}),
	})
}

// batchEvents generates events from UniFi for InfluxDB.
func (u *InfluxUnifi) batchEvent(r report, i *unifi.Event) { // nolint: funlen
	fields := map[string]interface{}{
		/*
			"site_id":                  i.SiteID,
			"flow_id":                  i.FlowID,
			"inner_alert_signature":    i.InnerAlertSignature,
			"inner_alert_gid":          i.InnerAlertGID,
			"inner_alert_rev":          i.InnerAlertRev,
			"inner_alert_severity":     i.InnerAlertSeverity,
			"inner_alert_signature_id": i.InnerAlertSignatureID,
			"unique_alertid":           i.UniqueAlertID,
			"usgipASN":                 i.UsgIPASN,
			"usgipCountry":             i.UsgIPCountry,
			"srcipASN":                 i.SrcIPASN,
			"srcipCountry":             i.SrcIPCountry,
		*/
		"dest_port":               i.DestPort,
		"src_port":                i.SrcPort,
		"bytes":                   i.Bytes,
		"duration":                i.Duration,
		"admin":                   i.Admin,
		"ap":                      i.Ap,
		"ap_from":                 i.ApFrom,
		"ap_name":                 i.ApName,
		"ap_to":                   i.ApTo,
		"app_proto":               i.AppProto,
		"catname":                 i.Catname,
		"channel":                 i.Channel,
		"channel_from":            i.ChannelFrom,
		"channel_to":              i.ChannelTo,
		"dest_ip":                 i.DestIP,
		"dst_mac":                 i.DstMAC,
		"guest":                   i.Guest,
		"gw":                      i.Gw,
		"gw_name":                 i.GwName,
		"host":                    i.Host,
		"hostname":                i.Hostname,
		"ip":                      i.IP,
		"inner_alert_action":      i.InnerAlertAction,
		"inner_alert_category":    i.InnerAlertCategory,
		"key":                     i.Key,
		"msg":                     i.Msg,
		"network":                 i.Network,
		"proto":                   i.Proto,
		"radio":                   i.Radio,
		"radio_from":              i.RadioFrom,
		"radio_to":                i.RadioTo,
		"src_ip":                  i.SrcIP,
		"src_mac":                 i.SrcMAC,
		"ssid":                    i.SSID,
		"sw":                      i.Sw,
		"sw_name":                 i.SwName,
		"user":                    i.User,
		"usgip":                   i.USGIP,
		"dstipGeo_asn":            i.DestIPGeo.Asn,
		"dstipGeo_latitude":       i.DestIPGeo.Latitude,
		"dstipGeo_longitude":      i.DestIPGeo.Longitude,
		"dstipGeo_city":           i.DestIPGeo.City,
		"dstipGeo_continent_code": i.DestIPGeo.ContinentCode,
		"dstipGeo_country_code":   i.DestIPGeo.CountryCode,
		"dstipGeo_country_name":   i.DestIPGeo.CountryName,
		"dstipGeo_organization":   i.DestIPGeo.Organization,
		"srcipGeo_asn":            i.SourceIPGeo.Asn,
		"srcipGeo_latitude":       i.SourceIPGeo.Latitude,
		"srcipGeo_longitude":      i.SourceIPGeo.Longitude,
		"srcipGeo_city":           i.SourceIPGeo.City,
		"srcipGeo_continent_code": i.SourceIPGeo.ContinentCode,
		"srcipGeo_country_code":   i.SourceIPGeo.CountryCode,
		"srcipGeo_country_name":   i.SourceIPGeo.CountryName,
		"srcipGeo_organization":   i.SourceIPGeo.Organization,
	}

	r.send(&metric{
		TS:     i.Datetime,
		Table:  "unifi_events",
		Fields: cleanFields(fields),
		Tags: cleanTags(map[string]string{
			"site_name":  i.SiteName,
			"source":     i.SourceName,
			"in_iface":   i.InIface,
			"event_type": i.EventType,
			"subsystem":  i.Subsystem,
			"is_admin":   i.IsAdmin.Txt,
			"gw_name":    i.GwName, // also field
			"ap_name":    i.ApName, // also field
			"sw_name":    i.SwName, // also field
			"ssid":       i.SSID,   // also field
		}),
	})
}

// cleanTags removes any tag that is empty.
func cleanTags(tags map[string]string) map[string]string {
	for i := range tags {
		if tags[i] == "" {
			delete(tags, i)
		}
	}

	return tags
}

// cleanFields removes any field with a default (or empty) value.
func cleanFields(fields map[string]interface{}) map[string]interface{} {
	for s := range fields {
		switch v := fields[s].(type) {
		case nil:
			delete(fields, s)
		case int, int64, float64:
			if v == 0 {
				delete(fields, s)
			}
		case unifi.FlexBool:
			if v.Txt == "" {
				delete(fields, s)
			}
		case unifi.FlexInt:
			if v.Txt == "" {
				delete(fields, s)
			}
		case string:
			if v == "" {
				delete(fields, s)
			}
		}
	}

	return fields
}
