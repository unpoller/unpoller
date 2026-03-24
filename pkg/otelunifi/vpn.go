package otelunifi

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
)

// exportVPNMeshes emits Site Magic site-to-site VPN mesh metrics.
func (u *OtelOutput) exportVPNMeshes(ctx context.Context, meter metric.Meter, m *poller.Metrics, r *Report) {
	for _, item := range m.VPNMeshes {
		mesh, ok := item.(*unifi.MagicSiteToSiteVPN)
		if !ok {
			continue
		}

		meshAttrs := attribute.NewSet(
			attribute.String("site_name", mesh.SiteName),
			attribute.String("source", mesh.SourceName),
			attribute.String("mesh_name", mesh.Name),
		)

		paused := 0.0
		if mesh.Pause.Val {
			paused = 1.0
		}

		u.recordGauge(ctx, meter, r, "unifi_vpn_mesh_paused",
			"Site Magic VPN mesh paused (1/0)", paused, meshAttrs)
		u.recordGauge(ctx, meter, r, "unifi_vpn_mesh_connections_total",
			"Total connections in Site Magic VPN mesh", float64(len(mesh.Connections)), meshAttrs)
		u.recordGauge(ctx, meter, r, "unifi_vpn_mesh_devices_total",
			"Total devices in Site Magic VPN mesh", float64(len(mesh.Devices)), meshAttrs)

		for i := range mesh.Status {
			s := &mesh.Status[i]

			statusAttrs := attribute.NewSet(
				attribute.String("site_name", mesh.SiteName),
				attribute.String("source", mesh.SourceName),
				attribute.String("mesh_name", mesh.Name),
				attribute.String("status_site", s.SiteID),
			)

			u.recordGauge(ctx, meter, r, "unifi_vpn_mesh_status_errors",
				"Number of errors for a site in a Site Magic VPN mesh", float64(len(s.Errors)), statusAttrs)
			u.recordGauge(ctx, meter, r, "unifi_vpn_mesh_status_warnings",
				"Number of warnings for a site in a Site Magic VPN mesh", float64(len(s.Warnings)), statusAttrs)

			for j := range s.Connections {
				conn := &s.Connections[j]

				connected := 0.0
				if conn.Connected.Val {
					connected = 1.0
				}

				connAttrs := attribute.NewSet(
					attribute.String("site_name", mesh.SiteName),
					attribute.String("source", mesh.SourceName),
					attribute.String("mesh_name", mesh.Name),
					attribute.String("connection_id", conn.ConnectionID),
					attribute.String("status_site", s.SiteID),
				)

				u.recordGauge(ctx, meter, r, "unifi_vpn_tunnel_connected",
					"Site Magic VPN tunnel connection status (1=connected, 0=disconnected)", connected, connAttrs)
				u.recordGauge(ctx, meter, r, "unifi_vpn_tunnel_association_time",
					"Site Magic VPN tunnel association Unix timestamp", conn.AssociationTime.Val, connAttrs)
				u.recordGauge(ctx, meter, r, "unifi_vpn_tunnel_errors",
					"Number of errors on a Site Magic VPN tunnel connection", float64(len(conn.Errors)), connAttrs)
			}
		}
	}
}
