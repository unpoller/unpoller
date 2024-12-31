package influxunifi

import (
	"time"

	"github.com/unpoller/unifi/v5"
)

const (
	alarmT   = item("Alarm")
	anomalyT = item("Anomaly")
)

// batchAlarms generates alarm datapoints for InfluxDB.
func (u *InfluxUnifi) batchAlarms(r report, event *unifi.Alarm) { // nolint:dupl
	if time.Since(event.Datetime) > u.Interval.Duration+time.Second {
		return // The event is older than our interval, ignore it.
	}

	fields := map[string]any{
		"dest_port":            event.DestPort,
		"src_port":             event.SrcPort,
		"dest_ip":              event.DestIP,
		"dst_mac":              event.DstMAC,
		"host":                 event.Host,
		"msg":                  event.Msg,
		"src_ip":               event.SrcIP,
		"src_mac":              event.SrcMAC,
		"dstip_asn":            event.DestIPGeo.Asn,
		"dstip_latitude":       event.DestIPGeo.Latitude,
		"dstip_longitude":      event.DestIPGeo.Longitude,
		"dstip_city":           event.DestIPGeo.City,
		"dstip_continent_code": event.DestIPGeo.ContinentCode,
		"dstip_country_code":   event.DestIPGeo.CountryCode,
		"dstip_country_name":   event.DestIPGeo.CountryName,
		"dstip_organization":   event.DestIPGeo.Organization,
		"srcip_asn":            event.SourceIPGeo.Asn,
		"srcip_latitude":       event.SourceIPGeo.Latitude,
		"srcip_longitude":      event.SourceIPGeo.Longitude,
		"srcip_city":           event.SourceIPGeo.City,
		"srcip_continent_code": event.SourceIPGeo.ContinentCode,
		"srcip_country_code":   event.SourceIPGeo.CountryCode,
		"srcip_country_name":   event.SourceIPGeo.CountryName,
		"srcip_organization":   event.SourceIPGeo.Organization,
	}

	r.addCount(alarmT)
	r.send(&metric{
		Table:  "unifi_alarm",
		TS:     event.Datetime,
		Fields: cleanFields(fields),
		Tags: cleanTags(map[string]string{
			"site_name":  event.SiteName,
			"source":     event.SourceName,
			"in_iface":   event.InIface,
			"event_type": event.EventType,
			"subsystem":  event.Subsystem,
			"archived":   event.Archived.Txt,
			"usgip":      event.USGIP,
			"proto":      event.Proto,
			"key":        event.Key,
			"catname":    event.Catname.String(),
			"app_proto":  event.AppProto,
			"action":     event.InnerAlertAction,
		}),
	})
}

// batchAnomaly generates Anomalies from UniFi for InfluxDB.
func (u *InfluxUnifi) batchAnomaly(r report, event *unifi.Anomaly) {
	if time.Since(event.Datetime) > u.Interval.Duration+time.Second {
		return // The event is older than our interval, ignore it.
	}

	r.addCount(anomalyT)
	r.send(&metric{
		TS:     event.Datetime,
		Table:  "unifi_anomaly",
		Fields: map[string]any{"msg": event.Anomaly},
		Tags: cleanTags(map[string]string{
			"application": "unifi_anomaly",
			"source":      event.SourceName,
			"site_name":   event.SiteName,
			"device_mac":  event.DeviceMAC,
		}),
	})
}
