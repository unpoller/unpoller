package otelunifi

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v6"
)

// exportUAP emits metrics for a wireless access point.
func (u *OtelOutput) exportUAP(ctx context.Context, meter metric.Meter, r *Report, s *unifi.UAP) {
	if !s.Adopted.Val {
		return
	}

	attrs := attribute.NewSet(
		attribute.String("mac", s.Mac),
		attribute.String("site_name", s.SiteName),
		attribute.String("source", s.SourceName),
		attribute.String("name", s.Name),
		attribute.String("model", s.Model),
		attribute.String("version", s.Version),
		attribute.String("type", s.Type),
		attribute.String("ip", s.IP),
	)

	u.recordGauge(ctx, meter, r, "unifi_device_uap_uptime_seconds",
		"UAP uptime in seconds", s.Uptime.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_cpu_utilization",
		"UAP CPU utilization percentage", s.SystemStats.CPU.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_mem_utilization",
		"UAP memory utilization percentage", s.SystemStats.Mem.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_load_avg_1",
		"UAP load average 1-minute", s.SysStats.Loadavg1.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_load_avg_5",
		"UAP load average 5-minute", s.SysStats.Loadavg5.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_load_avg_15",
		"UAP load average 15-minute", s.SysStats.Loadavg15.Val, attrs)

	up := 0.0
	if s.State.Val == 1 {
		up = 1.0
	}

	u.recordGauge(ctx, meter, r, "unifi_device_uap_up",
		"Whether UAP is up (1) or down (0)", up, attrs)

	u.exportLTE(ctx, meter, r, s)

	for _, radio := range s.RadioTable {
		radioAttrs := attribute.NewSet(
			attribute.String("mac", s.Mac),
			attribute.String("site_name", s.SiteName),
			attribute.String("source", s.SourceName),
			attribute.String("name", s.Name),
			attribute.String("radio", radio.Radio),
			attribute.String("radio_name", radio.Name),
		)

		u.recordGauge(ctx, meter, r, "unifi_device_uap_radio_channel",
			"UAP radio channel", float64(radio.Channel.Val), radioAttrs)
		u.recordGauge(ctx, meter, r, "unifi_device_uap_radio_tx_power_dbm",
			"UAP radio transmit power in dBm", radio.TxPower.Val, radioAttrs)
	}

	for _, vap := range s.VapTable {
		vapAttrs := attribute.NewSet(
			attribute.String("mac", s.Mac),
			attribute.String("site_name", s.SiteName),
			attribute.String("source", s.SourceName),
			attribute.String("name", s.Name),
			attribute.String("essid", vap.Essid),
			attribute.String("bssid", vap.Bssid),
			attribute.String("radio", vap.Radio),
		)

		// NumSta is a plain int in the unifi library
		u.recordGauge(ctx, meter, r, "unifi_device_uap_vap_num_stations",
			"UAP VAP connected station count", float64(vap.NumSta), vapAttrs)
		u.recordGauge(ctx, meter, r, "unifi_device_uap_vap_satisfaction",
			"UAP VAP client satisfaction score", vap.Satisfaction.Val, vapAttrs)
		u.recordGauge(ctx, meter, r, "unifi_device_uap_vap_rx_bytes",
			"UAP VAP receive bytes total", vap.RxBytes.Val, vapAttrs)
		u.recordGauge(ctx, meter, r, "unifi_device_uap_vap_tx_bytes",
			"UAP VAP transmit bytes total", vap.TxBytes.Val, vapAttrs)
	}
}

// hasLTEStatus reports whether this access point is a cellular backup module.
// U-LTE and U-LTE-Pro are typed as UAPs; their live modem status arrives as lte_* fields.
func hasLTEStatus(d *unifi.UAP) bool {
	return d.LteState.Txt != "" || d.LteFailoverMode != "" || d.LteRat != "" ||
		d.LteNetworkOperator != "" || d.LteSignal != "" || d.LteConnected.Txt != ""
}

// exportLTE emits live cellular backup status for a U-LTE or U-LTE-Pro module.
func (u *OtelOutput) exportLTE(ctx context.Context, meter metric.Meter, r *Report, s *unifi.UAP) {
	if !hasLTEStatus(s) {
		return
	}

	attrs := attribute.NewSet(
		attribute.String("mac", s.Mac),
		attribute.String("site_name", s.SiteName),
		attribute.String("source", s.SourceName),
		attribute.String("name", s.Name),
		attribute.String("model", s.Model),
		attribute.String("type", s.Type),
		attribute.String("lte_state", s.LteState.Txt),
		attribute.String("lte_failover_mode", s.LteFailoverMode),
		attribute.String("lte_signal", s.LteSignal),
		attribute.String("lte_rat", s.LteRat),
		attribute.String("lte_mode", s.LteMode),
		attribute.String("lte_band", s.LteBand),
		attribute.String("lte_pdp_type", s.LtePdpType),
		attribute.String("lte_operator", s.LteNetworkOperator),
	)

	u.recordGauge(ctx, meter, r, "unifi_device_uap_lte_connected",
		"LTE modem connected to the carrier (1) or not (0)", s.LteConnected.Float64(), attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_lte_failover",
		"LTE failover carrying traffic (1) or standing by (0)", s.LteFailover.Float64(), attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_lte_rssi_dbm",
		"LTE RSSI in dBm", s.LteRssi.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_lte_rsrp_dbm",
		"LTE RSRP in dBm", s.LteRsrp.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_lte_rsrq_db",
		"LTE RSRQ in dB", s.LteRsrq.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_lte_rx_channel",
		"LTE receive channel", s.LteRxChannel.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uap_lte_tx_channel",
		"LTE transmit channel", s.LteTxChannel.Val, attrs)
}
