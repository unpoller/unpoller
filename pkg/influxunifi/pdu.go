package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// pduT is used as a name for printed/logged counters.
const pduT = item("PDU")

// batchPDU generates Unifi PDU data points for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchPDU(r report, s *unifi.PDU) {
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := map[string]string{
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"source":    s.SourceName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
	}
	fields := Combine(
		u.batchPDUstat(s.Stat.Sw),
		u.batchSysStats(s.SysStats, s.SystemStats),
		map[string]any{
			"guest-num_sta":               s.GuestNumSta.Val,
			"ip":                          s.IP,
			"bytes":                       s.Bytes.Val,
			"last_seen":                   s.LastSeen.Val,
			"rx_bytes":                    s.RxBytes.Val,
			"tx_bytes":                    s.TxBytes.Val,
			"uptime":                      s.Uptime.Val,
			"state":                       s.State.Val,
			"user-num_sta":                s.UserNumSta.Val,
			"upgradeable":                 s.Upgradeable.Val,
			"outlet_ac_power_budget":      s.OutletACPowerBudget.Val,
			"outlet_ac_power_consumption": s.OutletACPowerConsumption.Val,
			"outlet_enabled":              s.OutletEnabled.Val,
			"overheating":                 s.Overheating.Val,
			"power_source":                s.PowerSource.Val,
			"total_max_power":             s.TotalMaxPower.Val,
		})

	r.addCount(pduT)
	r.send(&metric{Table: "pdu", Tags: tags, Fields: fields})
	u.batchPortTable(r, tags, s.PortTable)

	for _, oo := range s.OutletOverrides {
		oot := cleanTags(map[string]string{
			"mac":          s.Mac,
			"site_name":    s.SiteName,
			"source":       s.SourceName,
			"name":         s.Name,
			"version":      s.Version,
			"model":        s.Model,
			"serial":       s.Serial,
			"type":         s.Type,
			"ip":           s.IP,
			"outlet_index": oo.Index.Txt,
			"outlet_name":  oo.Name,
		})
		ood := map[string]any{
			"cycle_enabled": oo.CycleEnabled.Val,
			"relay_state":   oo.RelayState.Val,
		}
		r.send(&metric{Table: "pdu.outlet_overrides", Tags: oot, Fields: ood})
	}

	for _, ot := range s.OutletTable {
		ott := cleanTags(map[string]string{
			"mac":          s.Mac,
			"site_name":    s.SiteName,
			"source":       s.SourceName,
			"name":         s.Name,
			"version":      s.Version,
			"model":        s.Model,
			"serial":       s.Serial,
			"type":         s.Type,
			"ip":           s.IP,
			"outlet_index": ot.Index.Txt,
			"outlet_name":  ot.Name,
		})
		otd := map[string]any{
			"cycle_enabled":       ot.CycleEnabled.Val,
			"relay_state":         ot.RelayState.Val,
			"outlet_caps":         ot.OutletCaps.Val,
			"outlet_power_factor": ot.OutletPowerFactor.Val,
			"outlet_current":      ot.OutletCurrent.Val,
			"outlet_power":        ot.OutletPower.Val,
			"outlet_voltage":      ot.OutletVoltage.Val,
		}
		r.send(&metric{Table: "pdu.outlet_table", Tags: ott, Fields: otd})
	}
}

func (u *InfluxUnifi) batchPDUstat(sw *unifi.Sw) map[string]any {
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
