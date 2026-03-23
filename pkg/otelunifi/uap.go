package otelunifi

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v5"
)

// exportUAP emits metrics for a wireless access point.
func (u *OtelOutput) exportUAP(ctx context.Context, meter metric.Meter, r *Report, s *unifi.UAP) {
	if !s.Adopted.Val || s.Locating.Val {
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
