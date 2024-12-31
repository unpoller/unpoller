package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// pduT is used as a name for printed/logged counters.
const pduT = item("PDU")

// batchPDU generates Unifi PDU datapoints for Datadog.
// These points can be passed directly to datadog.
func (u *DatadogUnifi) batchPDU(r report, s *unifi.PDU) {
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := cleanTags(map[string]string{
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
	data := CombineFloat64(
		u.batchUSWstat(s.Stat.Sw),
		u.batchSysStats(s.SysStats, s.SystemStats),
		map[string]float64{
			"guest_num_sta":               s.GuestNumSta.Val,
			"bytes":                       s.Bytes.Val,
			"outlet_ac_power_budget":      s.OutletACPowerBudget.Val,
			"outlet_ac_power_consumption": s.OutletACPowerConsumption.Val,
			"outlet_enabled":              boolToFloat64(s.OutletEnabled.Val),
			"overheating":                 boolToFloat64(s.Overheating.Val),
			"power_source":                s.PowerSource.Val,
			"total_max_power":             s.TotalMaxPower.Val,
			"last_seen":                   s.LastSeen.Val,
			"rx_bytes":                    s.RxBytes.Val,
			"tx_bytes":                    s.TxBytes.Val,
			"uptime":                      s.Uptime.Val,
			"state":                       s.State.Val,
			"user_num_sta":                s.UserNumSta.Val,
			"upgradeable":                 boolToFloat64(s.Upgradeable.Val),
		})

	r.addCount(pduT)

	metricName := metricNamespace("pdu")
	reportGaugeForFloat64Map(r, metricName, data, tags)

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
		ood := map[string]float64{
			"cycle_enabled": boolToFloat64(oo.CycleEnabled.Val),
			"relay_state":   boolToFloat64(oo.RelayState.Val),
		}
		reportGaugeForFloat64Map(r, metricNamespace("pdu.outlet_overrides"), ood, oot)
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
		otd := map[string]float64{
			"cycle_enabled":       boolToFloat64(ot.CycleEnabled.Val),
			"relay_state":         boolToFloat64(ot.RelayState.Val),
			"outlet_caps":         ot.OutletCaps.Val,
			"outlet_power_factor": ot.OutletPowerFactor.Val,
			"outlet_current":      ot.OutletCurrent.Val,
			"outlet_power":        ot.OutletPower.Val,
			"outlet_voltage":      ot.OutletVoltage.Val,
		}
		reportGaugeForFloat64Map(r, metricNamespace("pdu.outlet_table"), otd, ott)
	}
}
