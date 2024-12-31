package influxunifi

import (
	"time"

	"github.com/unpoller/unifi/v5"
)

// These constants are used as names for printed/logged counters.
const (
	eventT = item("Event")
	idsT   = item("IDs")
)

// batchIDs generates intrusion detection datapoints for InfluxDB.
func (u *InfluxUnifi) batchIDs(r report, i *unifi.IDS) { // nolint:dupl
	if time.Since(i.Datetime) > u.Interval.Duration+time.Second {
		return // The event is older than our interval, ignore it.
	}

	fields := map[string]any{
		"dest_port":            i.DestPort,
		"src_port":             i.SrcPort,
		"dest_ip":              i.DestIP,
		"dst_mac":              i.DstMAC,
		"host":                 i.Host,
		"msg":                  i.Msg,
		"src_ip":               i.SrcIP,
		"src_mac":              i.SrcMAC,
		"dstip_asn":            i.DestIPGeo.Asn,
		"dstip_latitude":       i.DestIPGeo.Latitude,
		"dstip_longitude":      i.DestIPGeo.Longitude,
		"dstip_city":           i.DestIPGeo.City,
		"dstip_continent_code": i.DestIPGeo.ContinentCode,
		"dstip_country_code":   i.DestIPGeo.CountryCode,
		"dstip_country_name":   i.DestIPGeo.CountryName,
		"dstip_organization":   i.DestIPGeo.Organization,
		"srcip_asn":            i.SourceIPGeo.Asn,
		"srcip_latitude":       i.SourceIPGeo.Latitude,
		"srcip_longitude":      i.SourceIPGeo.Longitude,
		"srcip_city":           i.SourceIPGeo.City,
		"srcip_continent_code": i.SourceIPGeo.ContinentCode,
		"srcip_country_code":   i.SourceIPGeo.CountryCode,
		"srcip_country_name":   i.SourceIPGeo.CountryName,
		"srcip_organization":   i.SourceIPGeo.Organization,
	}

	r.addCount(idsT)
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
			"usgip":      i.USGIP,
			"proto":      i.Proto,
			"key":        i.Key,
			"catname":    i.Catname.String(),
			"app_proto":  i.AppProto,
			"action":     i.InnerAlertAction,
		}),
	})
}

// batchEvents generates events from UniFi for InfluxDB.
func (u *InfluxUnifi) batchEvent(r report, i *unifi.Event) { // nolint: funlen
	if time.Since(i.Datetime) > u.Interval.Duration+time.Second {
		return // The event is older than our interval, ignore it.
	}

	fields := map[string]any{
		"msg":                  i.Msg,          // contains user[] or guest[] or admin[]
		"duration":             i.Duration.Val, // probably microseconds?
		"guest":                i.Guest,        // mac address
		"user":                 i.User,         // mac address
		"host":                 i.Host,         // usg device?
		"hostname":             i.Hostname,     // client name
		"dest_port":            i.DestPort,
		"src_port":             i.SrcPort,
		"bytes":                i.Bytes.Val,
		"dest_ip":              i.DestIP,
		"dst_mac":              i.DstMAC,
		"ip":                   i.IP,
		"src_ip":               i.SrcIP,
		"src_mac":              i.SrcMAC,
		"dstip_asn":            i.DestIPGeo.Asn,
		"dstip_latitude":       i.DestIPGeo.Latitude,
		"dstip_longitude":      i.DestIPGeo.Longitude,
		"dstip_city":           i.DestIPGeo.City,
		"dstip_continent_code": i.DestIPGeo.ContinentCode,
		"dstip_country_code":   i.DestIPGeo.CountryCode,
		"dstip_country_name":   i.DestIPGeo.CountryName,
		"dstip_organization":   i.DestIPGeo.Organization,
		"srcip_asn":            i.SourceIPGeo.Asn,
		"srcip_latitude":       i.SourceIPGeo.Latitude,
		"srcip_longitude":      i.SourceIPGeo.Longitude,
		"srcip_city":           i.SourceIPGeo.City,
		"srcip_continent_code": i.SourceIPGeo.ContinentCode,
		"srcip_country_code":   i.SourceIPGeo.CountryCode,
		"srcip_country_name":   i.SourceIPGeo.CountryName,
		"srcip_organization":   i.SourceIPGeo.Organization,
	}

	r.addCount(eventT)
	r.send(&metric{
		TS:     i.Datetime,
		Table:  "unifi_events",
		Fields: cleanFields(fields),
		Tags: cleanTags(map[string]string{
			"admin":        i.Admin, // username
			"site_name":    i.SiteName,
			"source":       i.SourceName,
			"ap_from":      i.ApFrom,
			"ap_to":        i.ApTo,
			"ap":           i.Ap,
			"ap_name":      i.ApName,
			"gw":           i.Gw,
			"gw_name":      i.GwName,
			"sw":           i.Sw,
			"sw_name":      i.SwName,
			"catname":      i.Catname.String(),
			"radio":        i.Radio,
			"radio_from":   i.RadioFrom,
			"radio_to":     i.RadioTo,
			"key":          i.Key,
			"in_iface":     i.InIface,
			"event_type":   i.EventType,
			"subsystem":    i.Subsystem,
			"ssid":         i.SSID,
			"is_admin":     i.IsAdmin.Txt,
			"channel":      i.Channel.Txt,
			"channel_from": i.ChannelFrom.Txt,
			"channel_to":   i.ChannelTo.Txt,
			"usgip":        i.USGIP,
			"network":      i.Network,
			"app_proto":    i.AppProto,
			"proto":        i.Proto,
			"action":       i.InnerAlertAction,
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
func cleanFields(fields map[string]any) map[string]any { //nolint:cyclop
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
