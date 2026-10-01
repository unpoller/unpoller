package influxunifi

import (
	"math"
	"strconv"

	"github.com/unpoller/unifi/v6"
)

// umbbT is used as a name for printed/logged counters.
const umbbT = item("UMBB")

// batchUMBB generates datapoints for UMBB (mobile broadband) devices, such as the U5G Max.
// The device itself goes to "umbb" with its cellular radio state; each SIM slot goes to "umbb_sim".
func (u *InfluxUnifi) batchUMBB(r report, s *unifi.UMBB) {
	if !s.Adopted.Val {
		return
	}

	rad := s.Mbb.Radio
	dl, ul := mbbBandwidth(rad)

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
		u.batchSysStats(s.SysStats, s.SystemStats),
		map[string]any{
			"ip":                 s.IP,
			"last_seen":          s.LastSeen.Val,
			"uptime":             s.Uptime.Val,
			"state":              s.State.Val,
			"upgradeable":        s.Upgradable.Val,
			"uplink_speed":       s.Uplink.Speed.Val,
			"uplink_max_speed":   s.Uplink.MaxSpeed.Val,
			"uplink_latency":     s.Uplink.Latency.Val,
			"uplink_uptime":      s.Uplink.Uptime.Val,
			"internet":           s.Internet.Val,
			"mbb_state":          s.Mbb.State,
			"failover_mode":      s.Mbb.Mode,
			"rat":                rad.Rat,
			"band":               rad.Band,
			"operator":           rad.NetworkOperator,
			"sa_mode":            rad.SA5GMode.Val,
			"current_slot":       rad.CurrentSlot.Val,
			"primary_slot":       s.MbbOverrides.PrimarySlot.Val,
			"mcc":                rad.Mcc.Val,
			"mnc":                rad.Mnc.Val,
			"cell_id":            rad.CellID.Val,
			"pci":                rad.Pci.Val,
			"channel":            rad.Channel.Val,
			"has_coverage":       rad.HasCoverage.Val,
			"roaming":            rad.Roaming.Val,
			"registration_state": rad.RegistrationState.Val,
			"nr_carriers":        float64(len(rad.CaNr)),
			"lte_carriers":       float64(len(rad.CaLte)),
			"dl_bandwidth_mhz":   dl,
			"ul_bandwidth_mhz":   ul,
		})

	for k, v := range mbbOptional(rad) {
		fields[k] = v
	}

	r.addCount(umbbT)
	r.send(&metric{Table: "umbb", Tags: tags, Fields: fields})

	for _, sim := range s.Mbb.Sim {
		r.send(&metric{
			Table: "umbb_sim",
			Tags: map[string]string{
				"mac":       s.Mac,
				"site_name": s.SiteName,
				"source":    s.SourceName,
				"name":      s.Name,
				"model":     s.Model,
				"type":      s.Type,
				"slot":      sim.Slot.Txt,
				"esim":      strconv.FormatBool(sim.Esim.Val),
				"spn":       sim.Spn,
			},
			Fields: map[string]any{
				"active":                sim.Active.Val,
				"card_present":          sim.CardPresent.Val,
				"pin_blocked":           sim.PinBlocked.Val,
				"has_carrier":           sim.HasCarrier.Val,
				"metered":               sim.Metered.Val,
				"data_warning":          sim.DataWarning.Val,
				"data_limited":          sim.DataLimited.Val,
				"display_state":         sim.DisplayState,
				"display_state_elapsed": sim.DisplayStateElapsed.Val,
				"network_reject_type":   sim.NetworkRejectType,
				"network_reject_text":   sim.NetworkRejectText,
				"rx_bytes":              sim.RxBytes.Val,
				"tx_bytes":              sim.TxBytes.Val,
			},
		})
	}
}

// mbbOptional returns the signal readings and bitrates the modem reported, keyed by field name.
// Missing readings are left out, so they are not mistaken for 0 dBm.
func mbbOptional(rad unifi.MBBRadio) map[string]float64 {
	out := map[string]float64{}

	for k, v := range map[string]unifi.FlexInt{
		"rsrp":               rad.Rsrp,
		"rsrq":               rad.Rsrq,
		"snr":                rad.Snr,
		"nr_rsrp":            rad.RsrpNr,
		"nr_rsrq":            rad.RsrqNr,
		"nr_snr":             rad.SnrNr,
		"signal":             rad.Signal,
		"signal_percent":     rad.SignalPercent,
		"lte_max_bitrate_dl": rad.MaxBitrateDl,
		"lte_max_bitrate_ul": rad.MaxBitrateUl,
		"nr_max_bitrate_dl":  rad.MaxBitrateDlNr,
		"nr_max_bitrate_ul":  rad.MaxBitrateUlNr,
	} {
		if v.Txt != "" {
			out[k] = v.Val
		}
	}

	return out
}

// mbbBandwidth sums the downlink and uplink bandwidth of every aggregated carrier, NR and LTE.
// A carrier with no uplink reports a null bandwidth, which decodes as zero.
// A sum that overflows is reported as zero: InfluxDB rejects a whole point with an infinite field.
func mbbBandwidth(rad unifi.MBBRadio) (dl, ul float64) {
	for _, c := range rad.CaNr {
		dl += c.DlBwMhz.Val
		ul += c.UlBwMhz.Val
	}

	for _, c := range rad.CaLte {
		dl += c.DlBwMhz.Val
		ul += c.UlBwMhz.Val
	}

	return finiteOrZero(dl), finiteOrZero(ul)
}

func finiteOrZero(v float64) float64 {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return 0
	}

	return v
}
