package datadogunifi

import (
	"fmt"
	"strconv"
	"time"

	"github.com/unpoller/unifi"
)

// These constants are used as names for printed/logged counters.
const (
	eventT = item("Event")
	idsT   = item("IDS")
)

// batchIDS generates intrusion detection datapoints for Datadog.
func (u *DatadogUnifi) batchIDS(r report, i *unifi.IDS) { // nolint:dupl
	if time.Since(i.Datetime) > u.Interval.Duration+time.Second {
		return // The event is older than our interval, ignore it.
	}

	tagMap := map[string]string{
		"dest_port":            strconv.Itoa(i.DestPort),
		"src_port":             strconv.Itoa(i.SrcPort),
		"dest_ip":              i.DestIP,
		"dst_mac":              i.DstMAC,
		"host":                 i.Host,
		"msg":                  i.Msg,
		"src_ip":               i.SrcIP,
		"src_mac":              i.SrcMAC,
		"dstip_asn":            fmt.Sprintf("%d", i.DestIPGeo.Asn),
		"dstip_latitude":       fmt.Sprintf("%0.6f", i.DestIPGeo.Latitude),
		"dstip_longitude":      fmt.Sprintf("%0.6f", i.DestIPGeo.Longitude),
		"dstip_city":           i.DestIPGeo.City,
		"dstip_continent_code": i.DestIPGeo.ContinentCode,
		"dstip_country_code":   i.DestIPGeo.CountryCode,
		"dstip_country_name":   i.DestIPGeo.CountryName,
		"dstip_organization":   i.DestIPGeo.Organization,
		"srcip_asn":            fmt.Sprintf("%d", i.SourceIPGeo.Asn),
		"srcip_latitude":       fmt.Sprintf("%0.6f", i.SourceIPGeo.Latitude),
		"srcip_longitude":      fmt.Sprintf("%0.6f", i.SourceIPGeo.Longitude),
		"srcip_city":           i.SourceIPGeo.City,
		"srcip_continent_code": i.SourceIPGeo.ContinentCode,
		"srcip_country_code":   i.SourceIPGeo.CountryCode,
		"srcip_country_name":   i.SourceIPGeo.CountryName,
		"srcip_organization":   i.SourceIPGeo.Organization,
		"site_name":            i.SiteName,
		"source":               i.SourceName,
		"in_iface":             i.InIface,
		"event_type":           i.EventType,
		"subsystem":            i.Subsystem,
		"archived":             i.Archived.Txt,
		"usgip":                i.USGIP,
		"proto":                i.Proto,
		"key":                  i.Key,
		"catname":              i.Catname,
		"app_proto":            i.AppProto,
		"action":               i.InnerAlertAction,
	}

	r.addCount(idsT)

	tagMap = cleanTags(tagMap)
	tags := tagMapToTags(tagMap)
	title := fmt.Sprintf("Intrusion Detection at %s from %s", i.SiteName, i.SourceName)
	r.reportEvent(title, i.Datetime, i.Msg, tags)
	r.reportWarnLog(fmt.Sprintf("[%d] %s: %s", i.Datetime.Unix(), title, i.Msg), tagMapToZapFields(tagMap))
}

// batchEvents generates events from UniFi for Datadog.
func (u *DatadogUnifi) batchEvent(r report, i *unifi.Event) { // nolint: funlen
	if time.Since(i.Datetime) > u.Interval.Duration+time.Second {
		return // The event is older than our interval, ignore it.
	}

	tagMap := map[string]string{
		"guest":                i.Guest,    // mac address
		"user":                 i.User,     // mac address
		"host":                 i.Host,     // usg device?
		"hostname":             i.Hostname, // client name
		"dest_port":            strconv.Itoa(i.DestPort),
		"src_port":             strconv.Itoa(i.SrcPort),
		"dest_ip":              i.DestIP,
		"dst_mac":              i.DstMAC,
		"ip":                   i.IP,
		"src_ip":               i.SrcIP,
		"src_mac":              i.SrcMAC,
		"dstip_asn":            fmt.Sprintf("%d", i.DestIPGeo.Asn),
		"dstip_latitude":       fmt.Sprintf("%0.6f", i.DestIPGeo.Latitude),
		"dstip_longitude":      fmt.Sprintf("%0.6f", i.DestIPGeo.Longitude),
		"dstip_city":           i.DestIPGeo.City,
		"dstip_continent_code": i.DestIPGeo.ContinentCode,
		"dstip_country_code":   i.DestIPGeo.CountryCode,
		"dstip_country_name":   i.DestIPGeo.CountryName,
		"dstip_organization":   i.DestIPGeo.Organization,
		"srcip_asn":            fmt.Sprintf("%d", i.SourceIPGeo.Asn),
		"srcip_latitude":       fmt.Sprintf("%0.6f", i.SourceIPGeo.Latitude),
		"srcip_longitude":      fmt.Sprintf("%0.6f", i.SourceIPGeo.Longitude),
		"srcip_city":           i.SourceIPGeo.City,
		"srcip_continent_code": i.SourceIPGeo.ContinentCode,
		"srcip_country_code":   i.SourceIPGeo.CountryCode,
		"srcip_country_name":   i.SourceIPGeo.CountryName,
		"srcip_organization":   i.SourceIPGeo.Organization,
		"admin":                i.Admin, // username
		"site_name":            i.SiteName,
		"source":               i.SourceName,
		"ap_from":              i.ApFrom,
		"ap_to":                i.ApTo,
		"ap":                   i.Ap,
		"ap_name":              i.ApName,
		"gw":                   i.Gw,
		"gw_name":              i.GwName,
		"sw":                   i.Sw,
		"sw_name":              i.SwName,
		"catname":              i.Catname,
		"radio":                i.Radio,
		"radio_from":           i.RadioFrom,
		"radio_to":             i.RadioTo,
		"key":                  i.Key,
		"in_iface":             i.InIface,
		"event_type":           i.EventType,
		"subsystem":            i.Subsystem,
		"ssid":                 i.SSID,
		"is_admin":             i.IsAdmin.Txt,
		"channel":              i.Channel.Txt,
		"channel_from":         i.ChannelFrom.Txt,
		"channel_to":           i.ChannelTo.Txt,
		"usgip":                i.USGIP,
		"network":              i.Network,
		"app_proto":            i.AppProto,
		"proto":                i.Proto,
		"action":               i.InnerAlertAction,
	}

	r.addCount(eventT)

	tagMap = cleanTags(tagMap)
	tags := tagMapToTags(tagMap)
	title := fmt.Sprintf("Unifi Event at %s from %s", i.SiteName, i.SourceName)
	r.reportEvent(title, i.Datetime, i.Msg, tags)
	r.reportInfoLog(fmt.Sprintf("[%d] %s: %s", i.Datetime.Unix(), title, i.Msg), tagMapToZapFields(tagMap))
}
