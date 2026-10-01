package otelunifi

import (
	"context"
	"math"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v6"
)

// exportUMBB emits metrics for a mobile broadband (cellular) device, such as the U5G Max.
func (u *OtelOutput) exportUMBB(ctx context.Context, meter metric.Meter, r *Report, s *unifi.UMBB) {
	if !s.Adopted.Val {
		return
	}

	rad := s.Mbb.Radio
	attrs := attribute.NewSet(
		attribute.String("mac", s.Mac),
		attribute.String("site_name", s.SiteName),
		attribute.String("source", s.SourceName),
		attribute.String("name", s.Name),
		attribute.String("model", s.Model),
		attribute.String("version", s.Version),
		attribute.String("type", s.Type),
		attribute.String("ip", s.IP),
		attribute.String("mbb_state", s.Mbb.State),
		attribute.String("mbb_mode", s.Mbb.Mode),
		attribute.String("rat", rad.Rat),
		attribute.String("band", rad.Band),
		attribute.String("operator", rad.NetworkOperator),
	)

	up := 0.0
	if s.State.Val == 1 {
		up = 1.0
	}

	u.recordGauge(ctx, meter, r, "unifi_device_umbb_up",
		"Whether UMBB is up (1) or down (0)", up, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_uptime_seconds",
		"UMBB uptime in seconds", s.Uptime.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_cpu_utilization",
		"UMBB CPU utilization percentage", s.SystemStats.CPU.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_mem_utilization",
		"UMBB memory utilization percentage", s.SystemStats.Mem.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_rsrp_dbm",
		"Cellular RSRP in dBm", rad.Rsrp.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_rsrq_db",
		"Cellular RSRQ in dB", rad.Rsrq.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_snr_db",
		"Cellular signal-to-noise ratio in dB", rad.Snr.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_signal_bars",
		"Cellular signal strength in bars (0-5)", rad.Signal.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_signal_percent",
		"Cellular signal strength percentage", rad.SignalPercent.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_coverage",
		"Modem has cellular coverage (1) or not (0)", rad.HasCoverage.Float64(), attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_roaming",
		"Modem is roaming (1) or on its home network (0)", rad.Roaming.Float64(), attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_registration_state",
		"Cellular network registration state code", rad.RegistrationState.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_nr_carriers",
		"Aggregated 5G NR component carriers in use", float64(len(rad.CaNr)), attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_lte_carriers",
		"Aggregated LTE component carriers in use", float64(len(rad.CaLte)), attrs)

	dl, ul := mbbBandwidth(rad)

	u.recordGauge(ctx, meter, r, "unifi_device_umbb_dl_bandwidth_mhz",
		"Total aggregated downlink carrier bandwidth in MHz", dl, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_umbb_ul_bandwidth_mhz",
		"Total aggregated uplink carrier bandwidth in MHz", ul, attrs)

	for _, sim := range s.Mbb.Sim {
		simAttrs := attribute.NewSet(
			attribute.String("mac", s.Mac),
			attribute.String("site_name", s.SiteName),
			attribute.String("source", s.SourceName),
			attribute.String("name", s.Name),
			attribute.String("slot", sim.Slot.Txt),
			attribute.String("esim", strconv.FormatBool(sim.Esim.Val)),
			attribute.String("spn", sim.Spn),
			attribute.String("display_state", sim.DisplayState),
			attribute.String("network_reject_type", sim.NetworkRejectType),
			attribute.String("network_reject_text", sim.NetworkRejectText),
		)

		u.recordGauge(ctx, meter, r, "unifi_device_umbb_sim_active",
			"SIM slot is the active data SIM (1) or not (0)", sim.Active.Float64(), simAttrs)
		u.recordGauge(ctx, meter, r, "unifi_device_umbb_sim_state_seconds",
			"Seconds the SIM has been in its current display state", sim.DisplayStateElapsed.Val, simAttrs)
		u.recordGauge(ctx, meter, r, "unifi_device_umbb_sim_network_reject_age",
			"Age of the last carrier network reject, as reported by the controller", sim.NetworkRejectAge.Val, simAttrs)
		u.recordGauge(ctx, meter, r, "unifi_device_umbb_sim_receive_bytes_total",
			"Bytes received over this SIM", sim.RxBytes.Val, simAttrs)
		u.recordGauge(ctx, meter, r, "unifi_device_umbb_sim_transmit_bytes_total",
			"Bytes transmitted over this SIM", sim.TxBytes.Val, simAttrs)
	}
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
