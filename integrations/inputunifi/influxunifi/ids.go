package influxunifi

import (
	influx "github.com/influxdata/influxdb1-client/v2"
	"golift.io/unifi"
)

// IDSPoints generates intrusion detection datapoints for InfluxDB.
// These points can be passed directly to influx.
func IDSPoints(i *unifi.IDS) ([]*influx.Point, error) {
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
	pt, err := influx.NewPoint("intrusion_detect", tags, fields, i.Datetime)
	if err != nil {
		return nil, err
	}
	return []*influx.Point{pt}, nil
}
