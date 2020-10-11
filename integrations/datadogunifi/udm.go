package datadogunifi

import (
	"fmt"

	"github.com/unifi-poller/unifi"
)

// reportSysStats is used by all device types.
func (u *DatadogUnifi) reportSysStats(r report, metricName func(string) string, s unifi.SysStats, ss unifi.SystemStats, tags []string) {
	data := map[string]float64{
		"loadavg_1":     s.Loadavg1.Val,
		"loadavg_5":     s.Loadavg5.Val,
		"loadavg_15":    s.Loadavg15.Val,
		"mem_used":      s.MemUsed.Val,
		"mem_buffer":    s.MemBuffer.Val,
		"mem_total":     s.MemTotal.Val,
		"cpu":           ss.CPU.Val,
		"mem":           ss.Mem.Val,
		"system_uptime": ss.Uptime.Val,
	}
	for name, value := range data {
		r.reportGauge(metricName(name), value, tags)
	}
}

func (u *DatadogUnifi) reportUDMtemps(r report, metricName func(string) string, tags []string, temps []unifi.Temperature) {
	for _, t := range temps {
		name := fmt.Sprintf("temp_%s", t.Name)
		r.reportGauge(metricName(name), t.Value, tags)
	}
}

// reportUDM generates Unifi Gateway datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *DatadogUnifi) reportUDM(r report, s *unifi.UDM) { // nolint: funlen
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	metricName := metricNamespace("usg")

	tags := []string{
		tag("source", s.SourceName),
		tag("ip", s.IP),
		tag("license_state", s.LicenseState),
		tag("mac", s.Mac),
		tag("site_name", s.SiteName),
		tag("name", s.Name),
		tag("version", s.Version),
		tag("model", s.Model),
		tag("serial", s.Serial),
		tag("type", s.Type),
	}
	u.reportUDMtemps(r, metricName, tags, s.Temperatures)
	u.reportUSGstats(r, metricName, tags, s.SpeedtestStatus, s.Stat.Gw, s.Uplink)
	u.reportSysStats(r, metricName, s.SysStats, s.SystemStats, tags)

	data := map[string]float64{
		"bytes":         s.Bytes.Val,
		"last_seen":     s.LastSeen.Val,
		"guest-num_sta": s.GuestNumSta.Val,
		"rx_bytes":      s.RxBytes.Val,
		"tx_bytes":      s.TxBytes.Val,
		"uptime":        s.Uptime.Val,
		"state":         s.State.Val,
		"user-num_sta":  s.UserNumSta.Val,
		"num_desktop":   s.NumDesktop.Val,
		"num_handheld":  s.NumHandheld.Val,
		"num_mobile":    s.NumMobile.Val,
	}
	for name, value := range data {
		r.reportGauge(metricName(name), value, tags)
	}
	u.reportNetTable(r, s.Name, s.SiteName, s.SourceName, s.NetworkTable)
	u.reportUSGwans(r, s.Name, s.SiteName, s.SourceName, s.Wan1, s.Wan2)

	tags = []string{
		tag("mac", s.Mac),
		tag("site_name", s.SiteName),
		tag("source", s.SourceName),
		tag("name", s.Name),
		tag("version", s.Version),
		tag("model", s.Model),
		tag("serial", s.Serial),
		tag("type", s.Type),
		tag("ip", s.IP),
	}
	metricName = metricNamespace("usw")
	u.reportUSWstat(r, metricName, tags, s.Stat.Sw)

	data = map[string]float64{
		"guest-num_sta": s.GuestNumSta.Val,
		"bytes":         s.Bytes.Val,
		"last_seen":     s.LastSeen.Val,
		"rx_bytes":      s.RxBytes.Val,
		"tx_bytes":      s.TxBytes.Val,
		"uptime":        s.Uptime.Val,
	}
	for name, value := range data {
		r.reportGauge(metricName(name), value, tags)
	}

	u.reportPortTable(r, s.Name, s.SiteName, s.SourceName, s.Type, s.PortTable) // udm has a usw in it.

	if s.Stat.Ap == nil {
		return // we're done now. the following code process UDM (non-pro) UAP data.
	}

	tags = []string{
		tag("mac", s.Mac),
		tag("site_name", s.SiteName),
		tag("source", s.SourceName),
		tag("name", s.Name),
		tag("version", s.Version),
		tag("model", s.Model),
		tag("serial", s.Serial),
		tag("type", s.Type),
	}

	metricName = metricNamespace("uap")
	u.reportUAPstats(s.Stat.Ap, r, metricName, tags)

	data = map[string]float64{
		"bytes":         s.Bytes.Val,
		"last_seen":     s.LastSeen.Val,
		"rx_bytes":      s.RxBytes.Val,
		"tx_bytes":      s.TxBytes.Val,
		"uptime":        s.Uptime.Val,
		"state":         s.State.Val,
		"user-num_sta":  s.UserNumSta.Val,
		"guest-num_sta": s.GuestNumSta.Val,
		"num_sta":       s.NumSta.Val,
	}
	for name, value := range data {
		r.reportGauge(metricName(name), value, tags)
	}

	u.reportRadTable(r, s.Name, s.SiteName, s.SourceName, *s.RadioTable, *s.RadioTableStats)
	u.reportVAPTable(r, s.Name, s.SiteName, s.SourceName, *s.VapTable)
}
