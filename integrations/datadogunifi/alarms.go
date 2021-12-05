package datadogunifi

import (
	"fmt"
	"strconv"
	"time"

	"github.com/unpoller/unifi"
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
		"dest_port":            strconv.Itoa(event.DestPort),
		"src_port":             strconv.Itoa(event.SrcPort),
		"dest_ip":              event.DestIP,
		"dst_mac":              event.DstMAC,
		"host":                 event.Host,
		"msg":                  event.Msg,
		"src_ip":               event.SrcIP,
		"src_mac":              event.SrcMAC,
		"dstip_asn":            fmt.Sprintf("%d", event.DestIPGeo.Asn),
		"dstip_latitude":       fmt.Sprintf("%0.6f", event.DestIPGeo.Latitude),
		"dstip_longitude":      fmt.Sprintf("%0.6f", event.DestIPGeo.Longitude),
		"dstip_city":           event.DestIPGeo.City,
		"dstip_continent_code": event.DestIPGeo.ContinentCode,
		"dstip_country_code":   event.DestIPGeo.CountryCode,
		"dstip_country_name":   event.DestIPGeo.CountryName,
		"dstip_organization":   event.DestIPGeo.Organization,
		"srcip_asn":            fmt.Sprintf("%d", event.SourceIPGeo.Asn),
		"srcip_latitude":       fmt.Sprintf("%0.6f", event.SourceIPGeo.Latitude),
		"srcip_longitude":      fmt.Sprintf("%0.6f", event.SourceIPGeo.Longitude),
		"srcip_city":           event.SourceIPGeo.City,
		"srcip_continent_code": event.SourceIPGeo.ContinentCode,
		"srcip_country_code":   event.SourceIPGeo.CountryCode,
		"srcip_country_name":   event.SourceIPGeo.CountryName,
		"srcip_organization":   event.SourceIPGeo.Organization,
		"site_name":            event.SiteName,
		"source":               event.SourceName,
		"in_iface":             event.InIface,
		"event_type":           event.EventType,
		"subsystem":            event.Subsystem,
		"archived":             event.Archived.Txt,
		"usgip":                event.USGIP,
		"proto":                event.Proto,
		"key":                  event.Key,
		"catname":              event.Catname,
		"app_proto":            event.AppProto,
		"action":               event.InnerAlertAction,
	}
	r.addCount(alarmT)

	tagMap = cleanTags(tagMap)
	tags := tagMapToTags(tagMap)
	title := fmt.Sprintf("[%s][%s] Alarm at %s from %s", event.EventType, event.Catname, event.SiteName, event.SourceName)
	r.reportEvent(title, event.Datetime, event.Msg, tags)
	r.reportWarnLog(fmt.Sprintf("[%d] %s: %s", event.Datetime.Unix(), title, event.Msg), tagMapToZapFields(tagMap))
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
	r.reportEvent(title, event.Datetime, event.Anomaly, tags)
	r.reportWarnLog(fmt.Sprintf("[%d] %s: %s", event.Datetime.Unix(), title, event.Anomaly), tagMapToZapFields(tagMap))
}
