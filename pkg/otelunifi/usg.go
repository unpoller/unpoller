package otelunifi

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v5"
)

// exportUSG emits metrics for a UniFi Security Gateway.
func (u *OtelOutput) exportUSG(ctx context.Context, meter metric.Meter, r *Report, s *unifi.USG) {
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

	u.recordGauge(ctx, meter, r, "unifi_device_usg_up",
		"Whether USG is up (1) or down (0)", up, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usg_uptime_seconds",
		"USG uptime in seconds", s.Uptime.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usg_cpu_utilization",
		"USG CPU utilization percentage", s.SystemStats.CPU.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usg_mem_utilization",
		"USG memory utilization percentage", s.SystemStats.Mem.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usg_load_avg_1",
		"USG load average 1-minute", s.SysStats.Loadavg1.Val, attrs)

	// Export WAN1 metrics
	u.exportUSGWan(ctx, meter, r, s, s.Wan1, "wan1")
	// Export WAN2 metrics if present
	u.exportUSGWan(ctx, meter, r, s, s.Wan2, "wan2")
}

// exportUSGWan emits metrics for a single WAN interface on a USG.
func (u *OtelOutput) exportUSGWan(
	ctx context.Context,
	meter metric.Meter,
	r *Report,
	s *unifi.USG,
	wan unifi.Wan,
	ifaceName string,
) {
	if wan.IP == "" {
		return
	}

	wanAttrs := attribute.NewSet(
		attribute.String("mac", s.Mac),
		attribute.String("site_name", s.SiteName),
		attribute.String("source", s.SourceName),
		attribute.String("name", s.Name),
		attribute.String("iface", ifaceName),
		attribute.String("ip", wan.IP),
	)

	u.recordGauge(ctx, meter, r, "unifi_device_usg_wan_rx_bytes",
		"USG WAN interface receive bytes total", wan.RxBytes.Val, wanAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usg_wan_tx_bytes",
		"USG WAN interface transmit bytes total", wan.TxBytes.Val, wanAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usg_wan_rx_packets",
		"USG WAN interface receive packets total", wan.RxPackets.Val, wanAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usg_wan_tx_packets",
		"USG WAN interface transmit packets total", wan.TxPackets.Val, wanAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usg_wan_rx_errors",
		"USG WAN interface receive errors total", wan.RxErrors.Val, wanAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usg_wan_tx_errors",
		"USG WAN interface transmit errors total", wan.TxErrors.Val, wanAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usg_wan_speed_mbps",
		"USG WAN interface link speed in Mbps", wan.Speed.Val, wanAttrs)
}
