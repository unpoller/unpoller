package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// ubbT is used as a name for printed/logged counters.
const ubbT = item("UBB")

// batchUBB generates UBB datapoints for Datadog.
// These points can be passed directly to datadog.
func (u *DatadogUnifi) batchUBB(r report, s *unifi.UBB) { // nolint: funlen
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

	data := CombineFloat64(
		u.batchSysStats(sysStats, systemStats),
		map[string]float64{
			"bytes":            s.Bytes.Val,
			"last_seen":        s.LastSeen.Val,
			"rx_bytes":         s.RxBytes.Val,
			"tx_bytes":         s.TxBytes.Val,
			"uptime":           s.Uptime.Val,
			"state":            s.State.Val,
			"user_num_sta":     s.UserNumSta.Val,
			"uplink_speed":     s.Uplink.Speed.Val,
			"uplink_max_speed": s.Uplink.MaxSpeed.Val,
			"uplink_latency":   s.Uplink.Latency.Val,
			"uplink_uptime":    s.Uplink.Uptime.Val,
		},
	)

	r.addCount(ubbT)

	metricName := metricNamespace("ubb")
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
		u.batchUBBstat(sw),
		map[string]float64{
			"bytes":     s.Bytes.Val,
			"last_seen": s.LastSeen.Val,
			"rx_bytes":  s.RxBytes.Val,
			"tx_bytes":  s.TxBytes.Val,
			"uptime":    s.Uptime.Val,
		})

	metricName = metricNamespace("ubb")
	reportGaugeForFloat64Map(r, metricName, data, tags)
}

func (u *DatadogUnifi) batchUBBstat(sw *unifi.Bb) map[string]float64 {
	if sw == nil {
		return map[string]float64{}
	}

	return map[string]float64{
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
