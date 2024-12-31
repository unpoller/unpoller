package influxunifi

import (
	"strings"

	"github.com/unpoller/unifi/v5"
)

// udmT is used as a name for printed/logged counters.
const udmT = item("UDM")

// Combine concatenates N maps. This will delete things if not used with caution.
func Combine(in ...map[string]any) map[string]any {
	out := make(map[string]any)

	for i := range in {
		for k := range in[i] {
			out[k] = in[i][k]
		}
	}

	return out
}

func sanitizeName(v string) string {
	return strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(v, " ", "_"), ")", ""), "(", ""))
}

// batchSysStats is used by all device types.
func (u *InfluxUnifi) batchSysStats(s unifi.SysStats, ss unifi.SystemStats) map[string]any {
	m := map[string]any{
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

	for k, v := range ss.Temps {
		temp := v.CelsiusInt64()

		if temp != 0 && k != "" {
			m["temp_"+sanitizeName(k)] = temp
		}
	}

	return m
}

func (u *InfluxUnifi) batchUDMtemps(temps []unifi.Temperature) map[string]any {
	output := make(map[string]any)

	for _, t := range temps {
		output["temp_"+sanitizeName(t.Name)] = t.Value
	}

	return output
}

func (u *InfluxUnifi) batchUDMstorage(storage []*unifi.Storage) map[string]any {
	output := make(map[string]any)

	for _, t := range storage {
		output["storage_"+sanitizeName(t.Name)+"_size"] = t.Size.Val
		output["storage_"+sanitizeName(t.Name)+"_used"] = t.Used.Val

		if t.Size.Val != 0 && t.Used.Val != 0 && t.Used.Val < t.Size.Val {
			output["storage_"+sanitizeName(t.Name)+"_pct"] = t.Used.Val / t.Size.Val * 100 //nolint:gomnd
		} else {
			output["storage_"+sanitizeName(t.Name)+"_pct"] = float64(0)
		}
	}

	return output
}

// batchUDM generates Unifi Gateway datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchUDM(r report, s *unifi.UDM) { // nolint: funlen
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
	fields := Combine(
		u.batchUDMstorage(s.Storage),
		u.batchUDMtemps(s.Temperatures),
		u.batchUSGstats(s.SpeedtestStatus, s.Stat.Gw, s.Uplink),
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
			"upgradeable":   s.Upgradeable.Val,
		},
	)

	r.addCount(udmT)
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
		u.batchUSWstat(s.Stat.Sw),
		map[string]any{
			"guest-num_sta": s.GuestNumSta.Val,
			"ip":            s.IP,
			"bytes":         s.Bytes.Val,
			"last_seen":     s.LastSeen.Val,
			"rx_bytes":      s.RxBytes.Val,
			"tx_bytes":      s.TxBytes.Val,
			"uptime":        s.Uptime.Val,
			"upgradeable":   s.Upgradeable.Val,
		})

	r.send(&metric{Table: "usw", Tags: tags, Fields: fields})
	u.batchPortTable(r, tags, s.PortTable) // udm has a usw in it.

	if s.Stat.Ap == nil {
		return // we're done now. the following code process UDM (non-pro) UAP data.
	}

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
	fields = u.processUAPstats(s.Stat.Ap)
	fields["ip"] = s.IP
	fields["bytes"] = s.Bytes.Val
	fields["last_seen"] = s.LastSeen.Val
	fields["rx_bytes"] = s.RxBytes.Val
	fields["tx_bytes"] = s.TxBytes.Val
	fields["uptime"] = s.Uptime.Val
	fields["state"] = s.State
	fields["user-num_sta"] = int(s.UserNumSta.Val)
	fields["guest-num_sta"] = int(s.GuestNumSta.Val)
	fields["num_sta"] = s.NumSta.Val

	r.send(&metric{Table: "uap", Tags: tags, Fields: fields})
	u.processRadTable(r, tags, *s.RadioTable, *s.RadioTableStats)
	u.processVAPTable(r, tags, *s.VapTable)
}
