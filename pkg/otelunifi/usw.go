package otelunifi

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v5"
)

// exportUSW emits metrics for a UniFi switch.
func (u *OtelOutput) exportUSW(ctx context.Context, meter metric.Meter, r *Report, s *unifi.USW) {
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

	u.recordGauge(ctx, meter, r, "unifi_device_usw_up",
		"Whether USW is up (1) or down (0)", up, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_uptime_seconds",
		"USW uptime in seconds", s.Uptime.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_cpu_utilization",
		"USW CPU utilization percentage", s.SystemStats.CPU.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_mem_utilization",
		"USW memory utilization percentage", s.SystemStats.Mem.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_load_avg_1",
		"USW load average 1-minute", s.SysStats.Loadavg1.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_rx_bytes",
		"USW total receive bytes", s.Stat.RxBytes.Val, attrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_tx_bytes",
		"USW total transmit bytes", s.Stat.TxBytes.Val, attrs)

	if !u.DeadPorts {
		for _, p := range s.PortTable {
			if !p.Up.Val || !p.Enable.Val {
				continue
			}

			u.exportUSWPort(ctx, meter, r, s, p)
		}
	} else {
		for _, p := range s.PortTable {
			u.exportUSWPort(ctx, meter, r, s, p)
		}
	}
}

// exportUSWPort emits metrics for a single switch port.
func (u *OtelOutput) exportUSWPort(
	ctx context.Context,
	meter metric.Meter,
	r *Report,
	s *unifi.USW,
	p unifi.Port,
) {
	portAttrs := attribute.NewSet(
		attribute.String("mac", s.Mac),
		attribute.String("site_name", s.SiteName),
		attribute.String("source", s.SourceName),
		attribute.String("name", s.Name),
		attribute.String("port_name", p.Name),
		attribute.Int64("port_num", int64(p.PortIdx.Val)),
		attribute.String("port_mac", p.Mac),
		attribute.String("port_ip", p.IP),
	)

	portUp := 0.0
	if p.Up.Val {
		portUp = 1.0
	}

	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_up",
		"Whether switch port is up (1) or down (0)", portUp, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_speed_mbps",
		"Switch port speed in Mbps", p.Speed.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_rx_bytes",
		"Switch port receive bytes total", p.RxBytes.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_tx_bytes",
		"Switch port transmit bytes total", p.TxBytes.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_rx_bytes_rate",
		"Switch port receive bytes rate", p.RxBytesR.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_tx_bytes_rate",
		"Switch port transmit bytes rate", p.TxBytesR.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_rx_packets",
		"Switch port receive packets total", p.RxPackets.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_tx_packets",
		"Switch port transmit packets total", p.TxPackets.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_rx_errors",
		"Switch port receive errors total", p.RxErrors.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_tx_errors",
		"Switch port transmit errors total", p.TxErrors.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_rx_dropped",
		"Switch port receive dropped total", p.RxDropped.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_tx_dropped",
		"Switch port transmit dropped total", p.TxDropped.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_poe_current_amps",
		"Switch port PoE current in amps", p.PoeCurrent.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_poe_power_watts",
		"Switch port PoE power in watts", p.PoePower.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_poe_voltage",
		"Switch port PoE voltage", p.PoeVoltage.Val, portAttrs)
	u.recordGauge(ctx, meter, r, "unifi_device_usw_port_satisfaction",
		"Switch port satisfaction score", p.Satisfaction.Val, portAttrs)
}
