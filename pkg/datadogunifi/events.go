package datadogunifi

import (
	"fmt"
	"strconv"
	"time"

	"github.com/unpoller/unifi/v5"
)

// These constants are used as names for printed/logged counters.
const (
	eventT = item("Event")
	idsT   = item("IDs")
)

// batchIDs generates intrusion detection datapoints for Datadog.
func (u *DatadogUnifi) batchIDs(r report, i *unifi.IDS) { // nolint:dupl
	if time.Since(i.Datetime) > u.Interval.Duration+time.Second {
		return // The event is older than our interval, ignore it.
	}

	tagMap := map[string]string{
		"dest_port":             strconv.Itoa(i.DestPort.Int()),
		"src_port":              strconv.Itoa(i.SrcPort.Int()),
		"dest_ip":               i.DestIP,
		"dst_mac":               i.DstMAC,
		"host":                  i.Host,
		"msg":                   i.Msg,
		"src_ip":                i.SrcIP,
		"src_mac":               i.SrcMAC,
		"dst_ip_asn":            fmt.Sprintf("%d", i.DestIPGeo.Asn),
		"dst_ip_latitude":       fmt.Sprintf("%0.6f", i.DestIPGeo.Latitude),
		"dst_ip_longitude":      fmt.Sprintf("%0.6f", i.DestIPGeo.Longitude),
		"dst_ip_city":           i.DestIPGeo.City,
		"dst_ip_continent_code": i.DestIPGeo.ContinentCode,
		"dst_ip_country_code":   i.DestIPGeo.CountryCode,
		"dst_ip_country_name":   i.DestIPGeo.CountryName,
		"dst_ip_organization":   i.DestIPGeo.Organization,
		"src_ip_asn":            fmt.Sprintf("%d", i.SourceIPGeo.Asn),
		"src_ip_latitude":       fmt.Sprintf("%0.6f", i.SourceIPGeo.Latitude),
		"src_ip_longitude":      fmt.Sprintf("%0.6f", i.SourceIPGeo.Longitude),
		"src_ip_city":           i.SourceIPGeo.City,
		"src_ip_continent_code": i.SourceIPGeo.ContinentCode,
		"src_ip_country_code":   i.SourceIPGeo.CountryCode,
		"src_ip_country_name":   i.SourceIPGeo.CountryName,
		"src_ip_organization":   i.SourceIPGeo.Organization,
		"site_name":             i.SiteName,
		"source":                i.SourceName,
		"in_iface":              i.InIface,
		"event_type":            i.EventType,
		"subsystem":             i.Subsystem,
		"archived":              i.Archived.Txt,
		"usg_ip":                i.USGIP,
		"proto":                 i.Proto,
		"key":                   i.Key,
		"catname":               i.Catname.String(),
		"app_proto":             i.AppProto,
		"action":                i.InnerAlertAction,
	}

	r.addCount(idsT)

	tagMap = cleanTags(tagMap)
	tags := tagMapToTags(tagMap)
	title := fmt.Sprintf("Intrusion Detection at %s from %s", i.SiteName, i.SourceName)
	_ = r.reportEvent(title, i.Datetime, i.Msg, tags)
	r.reportWarnLog(fmt.Sprintf("[%d] %s: %s - %s", i.Datetime.Unix(), title, i.Msg, tagMapToSimpleStrings(tagMap)))
}

// batchEvents generates events from UniFi for Datadog.
func (u *DatadogUnifi) batchEvent(r report, i *unifi.Event) { // nolint: funlen
	if time.Since(i.Datetime) > u.Interval.Duration+time.Second {
		return // The event is older than our interval, ignore it.
	}

	tagMap := map[string]string{
		"guest":                 i.Guest,    // mac address
		"user":                  i.User,     // mac address
		"host":                  i.Host,     // usg device?
		"hostname":              i.Hostname, // client name
		"dest_port":             strconv.Itoa(i.DestPort),
		"src_port":              strconv.Itoa(i.SrcPort),
		"dst_ip":                i.DestIP,
		"dst_mac":               i.DstMAC,
		"ip":                    i.IP,
		"src_ip":                i.SrcIP,
		"src_mac":               i.SrcMAC,
		"dst_ip_asn":            fmt.Sprintf("%d", i.DestIPGeo.Asn),
		"dst_ip_latitude":       fmt.Sprintf("%0.6f", i.DestIPGeo.Latitude),
		"dst_ip_longitude":      fmt.Sprintf("%0.6f", i.DestIPGeo.Longitude),
		"dst_ip_city":           i.DestIPGeo.City,
		"dst_ip_continent_code": i.DestIPGeo.ContinentCode,
		"dst_ip_country_code":   i.DestIPGeo.CountryCode,
		"dst_ip_country_name":   i.DestIPGeo.CountryName,
		"dst_ip_organization":   i.DestIPGeo.Organization,
		"src_ip_asn":            fmt.Sprintf("%d", i.SourceIPGeo.Asn),
		"src_ip_latitude":       fmt.Sprintf("%0.6f", i.SourceIPGeo.Latitude),
		"src_ip_longitude":      fmt.Sprintf("%0.6f", i.SourceIPGeo.Longitude),
		"src_ip_city":           i.SourceIPGeo.City,
		"src_ip_continent_code": i.SourceIPGeo.ContinentCode,
		"src_ip_country_code":   i.SourceIPGeo.CountryCode,
		"src_ip_country_name":   i.SourceIPGeo.CountryName,
		"src_ip_organization":   i.SourceIPGeo.Organization,
		"admin":                 i.Admin, // username
		"site_name":             i.SiteName,
		"source":                i.SourceName,
		"ap_from":               i.ApFrom,
		"ap_to":                 i.ApTo,
		"ap":                    i.Ap,
		"ap_name":               i.ApName,
		"gw":                    i.Gw,
		"gw_name":               i.GwName,
		"sw":                    i.Sw,
		"sw_name":               i.SwName,
		"catname":               i.Catname.String(),
		"radio":                 i.Radio,
		"radio_from":            i.RadioFrom,
		"radio_to":              i.RadioTo,
		"key":                   i.Key,
		"in_iface":              i.InIface,
		"event_type":            i.EventType,
		"subsystem":             i.Subsystem,
		"ssid":                  i.SSID,
		"is_admin":              i.IsAdmin.Txt,
		"channel":               i.Channel.Txt,
		"channel_from":          i.ChannelFrom.Txt,
		"channel_to":            i.ChannelTo.Txt,
		"usg_ip":                i.USGIP,
		"network":               i.Network,
		"app_proto":             i.AppProto,
		"proto":                 i.Proto,
		"action":                i.InnerAlertAction,
	}

	r.addCount(eventT)

	tagMap = cleanTags(tagMap)
	tags := tagMapToTags(tagMap)
	title := fmt.Sprintf("Unifi Event at %s from %s", i.SiteName, i.SourceName)
	_ = r.reportEvent(title, i.Datetime, i.Msg, tags)
	r.reportInfoLog(fmt.Sprintf("[%d] %s: %s - %s", i.Datetime.Unix(), title, i.Msg, tagMapToSimpleStrings(tagMap)))
}
