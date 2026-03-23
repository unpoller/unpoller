package otelunifi

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
)

// exportTopology emits network topology metrics.
func (u *OtelOutput) exportTopology(ctx context.Context, meter metric.Meter, m *poller.Metrics, r *Report) {
	for _, item := range m.Topologies {
		t, ok := item.(*unifi.Topology)
		if !ok {
			continue
		}

		siteAttrs := attribute.NewSet(
			attribute.String("site_name", t.SiteName),
			attribute.String("source", t.SourceName),
		)

		var (
			devices    int
			clients    int
			wired      int
			wireless   int
			fullDuplex int
		)

		unknownSwitch := 0.0
		if t.HasUnknownSwitch {
			unknownSwitch = 1.0
		}

		for i := range t.Vertices {
			switch t.Vertices[i].Type {
			case "DEVICE":
				devices++
			case "CLIENT":
				clients++
			}
		}

		bandCounts := make(map[string]int)

		for i := range t.Edges {
			e := &t.Edges[i]

			edgeAttrs := attribute.NewSet(
				attribute.String("uplink_mac", e.UplinkMac),
				attribute.String("downlink_mac", e.DownlinkMac),
				attribute.String("link_type", e.Type),
				attribute.String("site_name", t.SiteName),
				attribute.String("source", t.SourceName),
			)

			switch e.Type {
			case "WIRED":
				wired++

				if e.Duplex == "FULL_DUPLEX" {
					fullDuplex++
				}

				u.recordGauge(ctx, meter, r, "unifi_topology_link_rate_mbps",
					"Wired link rate in Mbps", e.RateMbps.Val, edgeAttrs)

			case "WIRELESS":
				wireless++

				if e.RadioBand != "" {
					bandCounts[e.RadioBand]++
				}

				if e.ExperienceScore.Val > 0 {
					u.recordGauge(ctx, meter, r, "unifi_topology_link_experience_score",
						"Wireless link experience score (0-100)", e.ExperienceScore.Val, edgeAttrs)
				}
			}
		}

		u.recordGauge(ctx, meter, r, "unifi_topology_vertices_total",
			"Total vertices in topology", float64(len(t.Vertices)), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_topology_edges_total",
			"Total edges/connections in topology", float64(len(t.Edges)), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_topology_devices_total",
			"UniFi devices in topology", float64(devices), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_topology_clients_total",
			"Clients in topology", float64(clients), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_topology_has_unknown_switch",
			"Unknown switch detected in topology (1/0)", unknownSwitch, siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_topology_connections_wired",
			"Number of wired connections", float64(wired), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_topology_connections_wireless",
			"Number of wireless connections", float64(wireless), siteAttrs)
		u.recordGauge(ctx, meter, r, "unifi_topology_wired_full_duplex",
			"Number of full-duplex wired links", float64(fullDuplex), siteAttrs)

		for band, count := range bandCounts {
			bandAttrs := attribute.NewSet(
				attribute.String("band", band),
				attribute.String("site_name", t.SiteName),
				attribute.String("source", t.SourceName),
			)

			u.recordGauge(ctx, meter, r, "unifi_topology_connections_by_band",
				"Number of wireless connections by radio band", float64(count), bandAttrs)
		}
	}
}
