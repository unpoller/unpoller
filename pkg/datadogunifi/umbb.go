package datadogunifi

import (
	"math"
	"strconv"

	"github.com/unpoller/unifi/v6"
)

// umbbT is used as a name for printed/logged counters.
const umbbT = item("UMBB")

// batchUMBB generates datapoints for UMBB (mobile broadband) devices, such as the U5G Max.
// The device itself goes to "umbb" with its cellular radio state; each SIM slot goes to "umbb_sim".
func (u *DatadogUnifi) batchUMBB(r report, s *unifi.UMBB) {
	if !s.Adopted.Val {
		return
	}

	rad := s.Mbb.Radio
	dl, ul := mbbBandwidth(rad)

	tags := cleanTags(map[string]string{
		"mac":           s.Mac,
		"site_name":     s.SiteName,
		"source":        s.SourceName,
		"name":          s.Name,
		"version":       s.Version,
		"model":         s.Model,
		"serial":        s.Serial,
		"type":          s.Type,
		"ip":            s.IP,
		"mbb_state":     s.Mbb.State,
		"failover_mode": s.Mbb.Mode,
		"rat":           rad.Rat,
		"band":          rad.Band,
		"operator":      rad.NetworkOperator,
	})

	data := CombineFloat64(
		u.batchSysStats(s.SysStats, s.SystemStats),
		map[string]float64{
			"last_seen":          s.LastSeen.Val,
			"uptime":             s.Uptime.Val,
			"state":              s.State.Val,
			"upgradeable":        boolToFloat64(s.Upgradable.Val),
			"uplink_speed":       s.Uplink.Speed.Val,
			"uplink_max_speed":   s.Uplink.MaxSpeed.Val,
			"uplink_latency":     s.Uplink.Latency.Val,
			"uplink_uptime":      s.Uplink.Uptime.Val,
			"sa_mode":            rad.SA5GMode.Float64(),
			"current_slot":       rad.CurrentSlot.Val,
			"primary_slot":       s.MbbOverrides.PrimarySlot.Val,
			"internet":           s.Internet.Float64(),
			"has_coverage":       rad.HasCoverage.Float64(),
			"roaming":            rad.Roaming.Float64(),
			"registration_state": rad.RegistrationState.Val,
			"nr_carriers":        float64(len(rad.CaNr)),
			"lte_carriers":       float64(len(rad.CaLte)),
			"dl_bandwidth_mhz":   dl,
			"ul_bandwidth_mhz":   ul,
		},
		mbbOptional(rad))

	r.addCount(umbbT)
	reportGaugeForFloat64Map(r, metricNamespace("umbb"), data, tags)

	for _, sim := range s.Mbb.Sim {
		simTags := cleanTags(map[string]string{
			"mac":                 s.Mac,
			"site_name":           s.SiteName,
			"source":              s.SourceName,
			"name":                s.Name,
			"model":               s.Model,
			"type":                s.Type,
			"slot":                sim.Slot.Txt,
			"esim":                strconv.FormatBool(sim.Esim.Val),
			"spn":                 sim.Spn,
			"display_state":       sim.DisplayState,
			"has_carrier":         strconv.FormatBool(sim.HasCarrier.Val),
			"metered":             strconv.FormatBool(sim.Metered.Val),
			"network_reject_type": sim.NetworkRejectType,
			"network_reject_text": sim.NetworkRejectText,
		})
		simData := map[string]float64{
			"active":                sim.Active.Float64(),
			"card_present":          sim.CardPresent.Float64(),
			"pin_blocked":           sim.PinBlocked.Float64(),
			"data_warning":          sim.DataWarning.Float64(),
			"data_limited":          sim.DataLimited.Float64(),
			"display_state_elapsed": sim.DisplayStateElapsed.Val,
			"rx_bytes":              sim.RxBytes.Val,
			"tx_bytes":              sim.TxBytes.Val,
		}

		reportGaugeForFloat64Map(r, metricNamespace("umbb_sim"), simData, simTags)
	}
}

// mbbOptional returns the signal readings and bitrates the modem reported, keyed by field name.
// Missing readings are left out, so they are not mistaken for 0 dBm.
func mbbOptional(rad unifi.MBBRadio) map[string]float64 {
	out := map[string]float64{}

	for k, v := range map[string]unifi.FlexInt{
		"rsrp":               rad.Rsrp,
		"rsrq":               rad.Rsrq,
		"nr_rsrp":            rad.RsrpNr,
		"nr_rsrq":            rad.RsrqNr,
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

	// SNR is a fractional dB reading, so the library decodes it as FlexFloat.
	for k, v := range map[string]unifi.FlexFloat{
		"snr":    rad.Snr,
		"nr_snr": rad.SnrNr,
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
