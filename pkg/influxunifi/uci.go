package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// uciT is used as a name for printed/logged counters.
const uciT = item("UCI")

// batchUCI generates UCI datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchUCI(r report, s *unifi.UCI) { // nolint: funlen
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

	var sw *unifi.Sw
	if s.Stat != nil {
		sw = s.Stat.Sw
	}

	sysStats := unifi.SysStats{}
	if s.SysStats != nil {
		sysStats = *s.SysStats
	}

	systemStats := unifi.SystemStats{}
	if s.SystemStats != nil {
		systemStats = *s.SystemStats
	}

	fields := Combine(
		u.batchSysStats(sysStats, systemStats),
		map[string]any{
			"source":        s.SourceName,
			"ip":            s.IP,
			"bytes":         s.Bytes.Val,
			"last_seen":     s.LastSeen.Val,
			"license_state": s.LicenseState,
			"rx_bytes":      s.RxBytes.Val,
			"tx_bytes":      s.TxBytes.Val,
			"uptime":        s.Uptime.Val,
			"state":         s.State.Val,
			"version":       s.Version,
		},
	)

	r.addCount(uciT)
	r.send(&metric{Table: "uci", Tags: tags, Fields: fields})

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
			"ip":        s.IP,
			"bytes":     s.Bytes.Val,
			"last_seen": s.LastSeen.Val,
			"rx_bytes":  s.RxBytes.Val,
			"tx_bytes":  s.TxBytes.Val,
			"uptime":    s.Uptime.Val,
		})

	r.send(&metric{Table: "uci", Tags: tags, Fields: fields})
}
