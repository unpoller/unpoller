package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// uciT is used as a name for printed/logged counters.
const uciT = item("UCI")

// batchUCI generates UCI datapoints for Datadog.
// These points can be passed directly to datadog.
func (u *DatadogUnifi) batchUCI(r report, s *unifi.UCI) { // nolint: funlen
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

	data := CombineFloat64(
		u.batchSysStats(sysStats, systemStats),
		map[string]float64{
			"bytes":     s.Bytes.Val,
			"last_seen": s.LastSeen.Val,
			"rx_bytes":  s.RxBytes.Val,
			"tx_bytes":  s.TxBytes.Val,
			"uptime":    s.Uptime.Val,
			"state":     s.State.Val,
		},
	)

	r.addCount(uciT)

	metricName := metricNamespace("uci")
	reportGaugeForFloat64Map(r, metricName, data, tags)

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
			"bytes":     s.Bytes.Val,
			"last_seen": s.LastSeen.Val,
			"rx_bytes":  s.RxBytes.Val,
			"tx_bytes":  s.TxBytes.Val,
			"uptime":    s.Uptime.Val,
		})

	metricName = metricNamespace("uci")
	reportGaugeForFloat64Map(r, metricName, data, tags)
}
