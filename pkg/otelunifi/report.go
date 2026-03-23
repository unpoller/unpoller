package otelunifi

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
)

// Report accumulates counters that are printed to a log line.
type Report struct {
	Total   int           // Total count of metrics recorded.
	Errors  int           // Total count of errors recording metrics.
	Sites   int           // Total count of sites exported.
	Clients int           // Total count of clients exported.
	UAP     int           // Total count of UAP devices exported.
	USW     int           // Total count of USW devices exported.
	USG     int           // Total count of USG devices exported.
	UDM     int           // Total count of UDM devices exported.
	UXG     int           // Total count of UXG devices exported.
	Elapsed time.Duration // Duration elapsed collecting and exporting.
}

func (r *Report) String() string {
	return fmt.Sprintf(
		"Sites: %d, Clients: %d, UAP: %d, USW: %d, USG/UDM/UXG: %d/%d/%d, Metrics: %d, Errs: %d, Elapsed: %v",
		r.Sites, r.Clients, r.UAP, r.USW, r.USG, r.UDM, r.UXG,
		r.Total, r.Errors, r.Elapsed.Round(time.Millisecond),
	)
}

// reportMetrics converts poller.Metrics to OTel measurements.
func (u *OtelOutput) reportMetrics(m *poller.Metrics, _ *poller.Events) (*Report, error) {
	r := &Report{}
	start := time.Now()

	meter := otel.GetMeterProvider().Meter(PluginName)

	ctx := context.Background()

	u.exportSites(ctx, meter, m, r)
	u.exportClients(ctx, meter, m, r)
	u.exportDevices(ctx, meter, m, r)
	u.exportFirewallPolicies(ctx, meter, m, r)

	r.Elapsed = time.Since(start)

	return r, nil
}

// exportSites emits site-level gauge metrics.
func (u *OtelOutput) exportSites(ctx context.Context, meter metric.Meter, m *poller.Metrics, r *Report) {
	for _, item := range m.Sites {
		s, ok := item.(*unifi.Site)
		if !ok {
			continue
		}

		r.Sites++

		for _, h := range s.Health {
			attrs := attribute.NewSet(
				attribute.String("site_name", s.SiteName),
				attribute.String("source", s.SourceName),
				attribute.String("subsystem", h.Subsystem),
				attribute.String("status", h.Status),
			)

			u.recordGauge(ctx, meter, r, "unifi_site_users",
				"Number of users on the site subsystem", h.NumUser.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_guests",
				"Number of guests on the site subsystem", h.NumGuest.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_iot",
				"Number of IoT devices on the site subsystem", h.NumIot.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_aps",
				"Number of access points", h.NumAp.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_gateways",
				"Number of gateways", h.NumGw.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_switches",
				"Number of switches", h.NumSw.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_adopted",
				"Number of adopted devices", h.NumAdopted.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_disconnected",
				"Number of disconnected devices", h.NumDisconnected.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_pending",
				"Number of pending devices", h.NumPending.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_disabled",
				"Number of disabled devices", h.NumDisabled.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_latency_seconds",
				"Site WAN latency in seconds", h.Latency.Val/1000, attrs) //nolint:mnd
			u.recordGauge(ctx, meter, r, "unifi_site_uptime_seconds",
				"Site uptime in seconds", h.Uptime.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_tx_bytes_rate",
				"Site transmit bytes rate", h.TxBytesR.Val, attrs)
			u.recordGauge(ctx, meter, r, "unifi_site_rx_bytes_rate",
				"Site receive bytes rate", h.RxBytesR.Val, attrs)
		}
	}
}

// exportClients emits per-client gauge metrics.
func (u *OtelOutput) exportClients(ctx context.Context, meter metric.Meter, m *poller.Metrics, r *Report) {
	for _, item := range m.Clients {
		c, ok := item.(*unifi.Client)
		if !ok {
			continue
		}

		r.Clients++

		attrs := attribute.NewSet(
			attribute.String("mac", c.Mac),
			attribute.String("site_name", c.SiteName),
			attribute.String("source", c.SourceName),
			attribute.String("name", c.Name),
			attribute.String("ip", c.IP),
			attribute.String("oui", c.Oui),
			attribute.String("network", c.Network),
			attribute.String("ap_name", c.ApName),
			attribute.String("sw_name", c.SwName),
			attribute.Bool("wired", c.IsWired.Val),
		)

		u.recordGauge(ctx, meter, r, "unifi_client_uptime_seconds",
			"Client uptime in seconds", c.Uptime.Val, attrs)
		u.recordGauge(ctx, meter, r, "unifi_client_rx_bytes",
			"Client total bytes received", c.RxBytes.Val, attrs)
		u.recordGauge(ctx, meter, r, "unifi_client_tx_bytes",
			"Client total bytes transmitted", c.TxBytes.Val, attrs)
		u.recordGauge(ctx, meter, r, "unifi_client_rx_bytes_rate",
			"Client receive bytes rate", c.RxBytesR.Val, attrs)
		u.recordGauge(ctx, meter, r, "unifi_client_tx_bytes_rate",
			"Client transmit bytes rate", c.TxBytesR.Val, attrs)

		if !c.IsWired.Val {
			wifiAttrs := attribute.NewSet(
				attribute.String("mac", c.Mac),
				attribute.String("site_name", c.SiteName),
				attribute.String("source", c.SourceName),
				attribute.String("name", c.Name),
				attribute.String("ip", c.IP),
				attribute.String("oui", c.Oui),
				attribute.String("network", c.Network),
				attribute.String("ap_name", c.ApName),
				attribute.String("sw_name", c.SwName),
				attribute.Bool("wired", false),
				attribute.String("essid", c.Essid),
				attribute.String("radio", c.Radio),
				attribute.String("radio_proto", c.RadioProto),
			)

			u.recordGauge(ctx, meter, r, "unifi_client_signal_db",
				"Client signal strength in dBm", c.Signal.Val, wifiAttrs)
			u.recordGauge(ctx, meter, r, "unifi_client_noise_db",
				"Client AP noise floor in dBm", c.Noise.Val, wifiAttrs)
			u.recordGauge(ctx, meter, r, "unifi_client_rssi_db",
				"Client RSSI in dBm", c.Rssi.Val, wifiAttrs)
			u.recordGauge(ctx, meter, r, "unifi_client_tx_rate_bps",
				"Client transmit rate in bps", c.TxRate.Val, wifiAttrs)
			u.recordGauge(ctx, meter, r, "unifi_client_rx_rate_bps",
				"Client receive rate in bps", c.RxRate.Val, wifiAttrs)
		}
	}
}

// exportDevices routes each device to its type-specific exporter.
func (u *OtelOutput) exportDevices(ctx context.Context, meter metric.Meter, m *poller.Metrics, r *Report) {
	for _, item := range m.Devices {
		switch d := item.(type) {
		case *unifi.UAP:
			r.UAP++
			u.exportUAP(ctx, meter, r, d)

		case *unifi.USW:
			r.USW++
			u.exportUSW(ctx, meter, r, d)

		case *unifi.USG:
			r.USG++
			u.exportUSG(ctx, meter, r, d)

		case *unifi.UDM:
			r.UDM++
			u.exportUDM(ctx, meter, r, d)

		case *unifi.UXG:
			r.UXG++
			u.exportUXG(ctx, meter, r, d)

		default:
			if u.Collector.Poller().LogUnknownTypes {
				u.LogDebugf("otel: unknown device type: %T", item)
			}
		}
	}
}

// recordGauge is a helper that records a single float64 gauge observation.
func (u *OtelOutput) recordGauge(
	_ context.Context,
	meter metric.Meter,
	r *Report,
	name, description string,
	value float64,
	attrs attribute.Set,
) {
	g, err := meter.Float64ObservableGauge(name, metric.WithDescription(description))
	if err != nil {
		r.Errors++
		u.LogDebugf("otel: creating gauge %s: %v", name, err)

		return
	}

	_, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveFloat64(g, value, metric.WithAttributeSet(attrs))

		return nil
	}, g)
	if err != nil {
		r.Errors++
		u.LogDebugf("otel: registering callback for %s: %v", name, err)

		return
	}

	r.Total++
}
