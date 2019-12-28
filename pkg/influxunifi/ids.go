package influxunifi

import (
	"golift.io/unifi"
)

// batchIDS generates intrusion detection datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchIDS(r report, i *unifi.IDS) {
	tags := map[string]string{
		"in_iface":       i.InIface,
		"event_type":     i.EventType,
		"proto":          i.Proto,
		"app_proto":      i.AppProto,
		"usgip":          i.Usgip,
		"country_code":   i.SrcipGeo.CountryCode,
		"country_name":   i.SrcipGeo.CountryName,
		"region":         i.SrcipGeo.Region,
		"city":           i.SrcipGeo.City,
		"postal_code":    i.SrcipGeo.PostalCode,
		"srcipASN":       i.SrcipASN,
		"usgipASN":       i.UsgipASN,
		"alert_category": i.InnerAlertCategory,
		"subsystem":      i.Subsystem,
		"catname":        i.Catname,
	}
	fields := map[string]interface{}{
		"event_type":   i.EventType,
		"proto":        i.Proto,
		"app_proto":    i.AppProto,
		"usgip":        i.Usgip,
		"country_name": i.SrcipGeo.CountryName,
		"city":         i.SrcipGeo.City,
		"postal_code":  i.SrcipGeo.PostalCode,
		"srcipASN":     i.SrcipASN,
		"usgipASN":     i.UsgipASN,
	}
	r.send(&metric{Table: "intrusion_detect", Tags: tags, Fields: fields})
}
