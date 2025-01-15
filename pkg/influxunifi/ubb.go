package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// ubbT is used as a name for printed/logged counters.
const ubbT = item("UBB")

// batchUXG generates UBB datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchUBB(r report, s *unifi.UBB) { // nolint: funlen
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

	var sw *unifi.Bb
	if s.Stat != nil {
		sw = s.Stat.Bb
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
			"source":           s.SourceName,
			"ip":               s.IP,
			"bytes":            s.Bytes.Val,
			"last_seen":        s.LastSeen.Val,
			"license_state":    s.LicenseState,
			"rx_bytes":         s.RxBytes.Val,
			"tx_bytes":         s.TxBytes.Val,
			"uptime":           s.Uptime.Val,
			"state":            s.State.Val,
			"user-num_sta":     s.UserNumSta.Val,
			"version":          s.Version,
			"uplink_speed":     s.Uplink.Speed.Val,
			"uplink_max_speed": s.Uplink.MaxSpeed.Val,
			"uplink_latency":   s.Uplink.Latency.Val,
			"uplink_uptime":    s.Uplink.Uptime.Val,
		},
	)

	r.addCount(ubbT)
	r.send(&metric{Table: "ubb", Tags: tags, Fields: fields})

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
		u.batchUBBstat(sw),
		map[string]any{
			"ip":        s.IP,
			"bytes":     s.Bytes.Val,
			"last_seen": s.LastSeen.Val,
			"rx_bytes":  s.RxBytes.Val,
			"tx_bytes":  s.TxBytes.Val,
			"uptime":    s.Uptime.Val,
		})

	r.send(&metric{Table: "ubb", Tags: tags, Fields: fields})
}

func (u *InfluxUnifi) batchUBBstat(sw *unifi.Bb) map[string]any {
	if sw == nil {
		return map[string]any{}
	}

	return map[string]any{
		"stat_bytes":      sw.Bytes.Val,
		"stat_rx_bytes":   sw.RxBytes.Val,
		"stat_rx_crypts":  sw.RxCrypts.Val,
		"stat_rx_dropped": sw.RxDropped.Val,
		"stat_rx_errors":  sw.RxErrors.Val,
		"stat_rx_frags":   sw.RxFrags.Val,
		"stat_rx_packets": sw.TxPackets.Val,
		"stat_tx_bytes":   sw.TxBytes.Val,
		"stat_tx_dropped": sw.TxDropped.Val,
		"stat_tx_errors":  sw.TxErrors.Val,
		"stat_tx_packets": sw.TxPackets.Val,
		"stat_tx_retries": sw.TxRetries.Val,
	}
}
