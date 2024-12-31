package datadogunifi

import (
	"fmt"
	"strconv"
	"time"

	"github.com/unpoller/unifi/v5"
)

const (
	alarmT   = item("Alarm")
	anomalyT = item("Anomaly")
)

// batchAlarms generates alarm events and logs for Datadog.
func (u *DatadogUnifi) batchAlarms(r report, event *unifi.Alarm) { // nolint:dupl
	if time.Since(event.Datetime) > u.Interval.Duration+time.Second {
		return // The event is older than our interval, ignore it.
	}

	tagMap := map[string]string{
		"dst_port":              strconv.Itoa(event.DestPort),
		"src_port":              strconv.Itoa(event.SrcPort),
		"dest_ip":               event.DestIP,
		"dst_mac":               event.DstMAC,
		"host":                  event.Host,
		"msg":                   event.Msg,
		"src_ip":                event.SrcIP,
		"src_mac":               event.SrcMAC,
		"dst_ip_asn":            fmt.Sprintf("%d", event.DestIPGeo.Asn),
		"dst_ip_latitude":       fmt.Sprintf("%0.6f", event.DestIPGeo.Latitude),
		"dst_ip_longitude":      fmt.Sprintf("%0.6f", event.DestIPGeo.Longitude),
		"dst_ip_city":           event.DestIPGeo.City,
		"dst_ip_continent_code": event.DestIPGeo.ContinentCode,
		"dst_ip_country_code":   event.DestIPGeo.CountryCode,
		"dst_ip_country_name":   event.DestIPGeo.CountryName,
		"dst_ip_organization":   event.DestIPGeo.Organization,
		"src_ip_asn":            fmt.Sprintf("%d", event.SourceIPGeo.Asn),
		"src_ip_latitude":       fmt.Sprintf("%0.6f", event.SourceIPGeo.Latitude),
		"src_ip_longitude":      fmt.Sprintf("%0.6f", event.SourceIPGeo.Longitude),
		"src_ip_city":           event.SourceIPGeo.City,
		"src_ip_continent_code": event.SourceIPGeo.ContinentCode,
		"src_ip_country_code":   event.SourceIPGeo.CountryCode,
		"src_ip_country_name":   event.SourceIPGeo.CountryName,
		"src_ip_organization":   event.SourceIPGeo.Organization,
		"site_name":             event.SiteName,
		"source":                event.SourceName,
		"in_iface":              event.InIface,
		"event_type":            event.EventType,
		"subsystem":             event.Subsystem,
		"archived":              event.Archived.Txt,
		"usg_ip":                event.USGIP,
		"proto":                 event.Proto,
		"key":                   event.Key,
		"catname":               event.Catname.String(),
		"app_proto":             event.AppProto,
		"action":                event.InnerAlertAction,
	}

	r.addCount(alarmT)

	tagMap = cleanTags(tagMap)
	tags := tagMapToTags(tagMap)
	title := fmt.Sprintf("[%s][%s] Alarm at %s from %s", event.EventType, event.Catname, event.SiteName, event.SourceName)
	_ = r.reportEvent(title, event.Datetime, event.Msg, tags)
	r.reportWarnLog(fmt.Sprintf("[%d] %s: %s - %s", event.Datetime.Unix(), title, event.Msg, tagMapToSimpleStrings(tagMap)))
}

// batchAnomaly generates Anomalies from UniFi for Datadog.
func (u *DatadogUnifi) batchAnomaly(r report, event *unifi.Anomaly) {
	if time.Since(event.Datetime) > u.Interval.Duration+time.Second {
		return // The event is older than our interval, ignore it.
	}

	r.addCount(anomalyT)

	tagMap := cleanTags(map[string]string{
		"application": "unifi_anomaly",
		"source":      event.SourceName,
		"site_name":   event.SiteName,
		"device_mac":  event.DeviceMAC,
	})
	tags := tagMapToTags(tagMap)

	title := fmt.Sprintf("Anomaly detected at %s from %s", event.SiteName, event.SourceName)
	_ = r.reportEvent(title, event.Datetime, event.Anomaly, tags)
	r.reportWarnLog(fmt.Sprintf("[%d] %s: %s - %s", event.Datetime.Unix(), title, event.Anomaly, tagMapToSimpleStrings(tagMap)))
}
