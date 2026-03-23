package otelunifi

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v5"
)

// exportUDM emits metrics for a UniFi Dream Machine (all variants).
func (u *OtelOutput) exportUDM(ctx context.Context, meter metric.Meter, r *Report, s *unifi.UDM) {
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

	up := 0.0
	if s.State.Val == 1 {
		up = 1.0
	}

	u.recordGauge(ctx, meter, r, "unifi_device_udm_up",
		"Whether UDM is up (1) or down (0)", up, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_udm_uptime_seconds",
		"UDM uptime in seconds", s.Uptime.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_udm_cpu_utilization",
		"UDM CPU utilization percentage", s.SystemStats.CPU.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_udm_mem_utilization",
		"UDM memory utilization percentage", s.SystemStats.Mem.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_udm_load_avg_1",
		"UDM load average 1-minute", s.SysStats.Loadavg1.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_udm_load_avg_5",
		"UDM load average 5-minute", s.SysStats.Loadavg5.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_udm_load_avg_15",
		"UDM load average 15-minute", s.SysStats.Loadavg15.Val, attrs)
}

// exportUXG emits metrics for a UniFi Next-Gen Gateway.
func (u *OtelOutput) exportUXG(ctx context.Context, meter metric.Meter, r *Report, s *unifi.UXG) {
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

	up := 0.0
	if s.State.Val == 1 {
		up = 1.0
	}

	u.recordGauge(ctx, meter, r, "unifi_device_uxg_up",
		"Whether UXG is up (1) or down (0)", up, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uxg_uptime_seconds",
		"UXG uptime in seconds", s.Uptime.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uxg_cpu_utilization",
		"UXG CPU utilization percentage", s.SystemStats.CPU.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uxg_mem_utilization",
		"UXG memory utilization percentage", s.SystemStats.Mem.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_uxg_load_avg_1",
		"UXG load average 1-minute", s.SysStats.Loadavg1.Val, attrs)
}
