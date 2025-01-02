package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// uxgT is used as a name for printed/logged counters.
const uxgT = item("UXG")

// batchUXG generates 10Gb Unifi Gateway datapoints for Datadog.
// These points can be passed directly to datadog.
func (u *DatadogUnifi) batchUXG(r report, s *unifi.UXG) { // nolint: funlen
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := cleanTags(map[string]string{
		"source":        s.SourceName,
		"mac":           s.Mac,
		"site_name":     s.SiteName,
		"name":          s.Name,
		"version":       s.Version,
		"model":         s.Model,
		"serial":        s.Serial,
		"type":          s.Type,
		"ip":            s.IP,
		"license_state": s.LicenseState,
	})

	var gw *unifi.Gw
	if s.Stat != nil {
		gw = s.Stat.Gw
	}

	var sw *unifi.Sw
	if s.Stat != nil {
		sw = s.Stat.Sw
	}

	data := CombineFloat64(
		u.batchUDMstorage(s.Storage),
		u.batchUDMtemps(s.Temperatures),
		u.batchUSGstats(s.SpeedtestStatus, gw, s.Uplink),
		u.batchSysStats(s.SysStats, s.SystemStats),
		map[string]float64{
			"bytes":            s.Bytes.Val,
			"last_seen":        s.LastSeen.Val,
			"guest_num_sta":    s.GuestNumSta.Val,
			"rx_bytes":         s.RxBytes.Val,
			"tx_bytes":         s.TxBytes.Val,
			"uptime":           s.Uptime.Val,
			"state":            s.State.Val,
			"user_num_sta":     s.UserNumSta.Val,
			"num_desktop":      s.NumDesktop.Val,
			"num_handheld":     s.NumHandheld.Val,
			"num_mobile":       s.NumMobile.Val,
			"uplink_speed":     s.Uplink.Speed.Val,
			"uplink_max_speed": s.Uplink.MaxSpeed.Val,
			"uplink_latency":   s.Uplink.Latency.Val,
			"uplink_uptime":    s.Uplink.Uptime.Val,
		},
	)

	r.addCount(uxgT)

	metricName := metricNamespace("usg")
	reportGaugeForFloat64Map(r, metricName, data, tags)

	u.batchNetTable(r, tags, s.NetworkTable)
	u.batchUSGwans(r, tags, s.Wan1, s.Wan2)

	tags = cleanTags(map[string]string{
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"source":    s.SourceName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
		"ip":        s.IP,
	})
	data = CombineFloat64(
		u.batchUSWstat(sw),
		map[string]float64{
			"guest_num_sta": s.GuestNumSta.Val,
			"bytes":         s.Bytes.Val,
			"last_seen":     s.LastSeen.Val,
			"rx_bytes":      s.RxBytes.Val,
			"tx_bytes":      s.TxBytes.Val,
			"uptime":        s.Uptime.Val,
		})

	metricName = metricNamespace("usw")
	reportGaugeForFloat64Map(r, metricName, data, tags)

	u.batchPortTable(r, tags, s.PortTable) // udm has a usw in it.
}
