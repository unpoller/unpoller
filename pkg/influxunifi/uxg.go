package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// uxgT is used as a name for printed/logged counters.
const uxgT = item("UXG")

// batchUXG generates 10Gb Unifi Gateway datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchUXG(r report, s *unifi.UXG) { // nolint: funlen
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := map[string]string{
		"source":    s.SourceName,
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
	}

	var gw *unifi.Gw
	if s.Stat != nil {
		gw = s.Stat.Gw
	}

	var sw *unifi.Sw
	if s.Stat != nil {
		sw = s.Stat.Sw
	}

	fields := Combine(
		u.batchUDMstorage(s.Storage),
		u.batchUDMtemps(s.Temperatures),
		u.batchUSGstats(s.SpeedtestStatus, gw, s.Uplink),
		u.batchSysStats(s.SysStats, s.SystemStats),
		map[string]any{
			"source":        s.SourceName,
			"ip":            s.IP,
			"bytes":         s.Bytes.Val,
			"last_seen":     s.LastSeen.Val,
			"license_state": s.LicenseState,
			"guest-num_sta": s.GuestNumSta.Val,
			"rx_bytes":      s.RxBytes.Val,
			"tx_bytes":      s.TxBytes.Val,
			"uptime":        s.Uptime.Val,
			"state":         s.State.Val,
			"user-num_sta":  s.UserNumSta.Val,
			"version":       s.Version,
			"num_desktop":   s.NumDesktop.Val,
			"num_handheld":  s.NumHandheld.Val,
			"num_mobile":    s.NumMobile.Val,
		},
	)

	r.addCount(uxgT)
	r.send(&metric{Table: "usg", Tags: tags, Fields: fields})
	u.batchNetTable(r, tags, s.NetworkTable)
	u.batchUSGwans(r, tags, s.Wan1, s.Wan2)

	tags = map[string]string{
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"source":    s.SourceName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
	}
	fields = Combine(
		u.batchUSWstat(sw),
		map[string]any{
			"guest-num_sta": s.GuestNumSta.Val,
			"ip":            s.IP,
			"bytes":         s.Bytes.Val,
			"last_seen":     s.LastSeen.Val,
			"rx_bytes":      s.RxBytes.Val,
			"tx_bytes":      s.TxBytes.Val,
			"uptime":        s.Uptime.Val,
		})

	r.send(&metric{Table: "usw", Tags: tags, Fields: fields})
	u.batchPortTable(r, tags, s.PortTable) // udm has a usw in it.
}
