package otelunifi

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
)

// exportPortAnomalies emits per-port anomaly metrics.
func (u *OtelOutput) exportPortAnomalies(ctx context.Context, meter metric.Meter, m *poller.Metrics, r *Report) {
	for _, item := range m.PortAnomalies {
		a, ok := item.(*unifi.PortAnomaly)
		if !ok {
			continue
		}

		attrs := attribute.NewSet(
			attribute.String("site_name", a.SiteName),
			attribute.String("source", a.SourceName),
			attribute.String("device_mac", a.DeviceMAC),
			attribute.String("port_idx", a.PortIdx.Txt),
			attribute.String("anomaly_type", a.AnomalyType),
		)

		u.recordGauge(ctx, meter, r, "unifi_port_anomaly_count",
			"Number of anomaly events on this port", a.Count.Val, attrs)
		u.recordGauge(ctx, meter, r, "unifi_port_anomaly_last_seen",
			"Unix timestamp of the last anomaly event on this port", a.LastSeen.Val, attrs)
	}
}
