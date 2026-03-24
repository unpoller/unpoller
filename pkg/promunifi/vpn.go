package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type vpnmesh struct {
	MeshPaused          *prometheus.Desc
	MeshConnectionsTotal *prometheus.Desc
	MeshDevicesTotal    *prometheus.Desc
	TunnelConnected     *prometheus.Desc
	TunnelAssociationTime *prometheus.Desc
	TunnelErrors        *prometheus.Desc
	StatusErrors        *prometheus.Desc
	StatusWarnings      *prometheus.Desc
}

func descVPNMesh(ns string) *vpnmesh {
	meshLabels := []string{"site_name", "source", "mesh_name"}
	connLabels := []string{"site_name", "source", "mesh_name", "connection_id", "status_site"}
	statusLabels := []string{"site_name", "source", "mesh_name", "status_site"}

	nd := prometheus.NewDesc

	return &vpnmesh{
		MeshPaused:            nd(ns+"vpn_mesh_paused", "Site Magic VPN mesh paused (1/0)", meshLabels, nil),
		MeshConnectionsTotal:  nd(ns+"vpn_mesh_connections_total", "Total connections in Site Magic VPN mesh", meshLabels, nil),
		MeshDevicesTotal:      nd(ns+"vpn_mesh_devices_total", "Total devices in Site Magic VPN mesh", meshLabels, nil),
		TunnelConnected:       nd(ns+"vpn_tunnel_connected", "Site Magic VPN tunnel connection status (1=connected, 0=disconnected)", connLabels, nil),
		TunnelAssociationTime: nd(ns+"vpn_tunnel_association_time", "Site Magic VPN tunnel association Unix timestamp", connLabels, nil),
		TunnelErrors:          nd(ns+"vpn_tunnel_errors", "Number of errors on a Site Magic VPN tunnel connection", connLabels, nil),
		StatusErrors:          nd(ns+"vpn_mesh_status_errors", "Number of errors for a site in a Site Magic VPN mesh", statusLabels, nil),
		StatusWarnings:        nd(ns+"vpn_mesh_status_warnings", "Number of warnings for a site in a Site Magic VPN mesh", statusLabels, nil),
	}
}

func (u *promUnifi) exportVPNMesh(r report, m *unifi.MagicSiteToSiteVPN) {
	if m == nil {
		return
	}

	meshLabels := []string{m.SiteName, m.SourceName, m.Name}

	paused := 0.0
	if m.Pause.Val {
		paused = 1.0
	}

	r.send([]*metric{
		{u.VPNMesh.MeshPaused, gauge, paused, meshLabels},
		{u.VPNMesh.MeshConnectionsTotal, gauge, float64(len(m.Connections)), meshLabels},
		{u.VPNMesh.MeshDevicesTotal, gauge, float64(len(m.Devices)), meshLabels},
	})

	for i := range m.Status {
		s := &m.Status[i]

		statusLabels := []string{m.SiteName, m.SourceName, m.Name, s.SiteID}

		r.send([]*metric{
			{u.VPNMesh.StatusErrors, gauge, float64(len(s.Errors)), statusLabels},
			{u.VPNMesh.StatusWarnings, gauge, float64(len(s.Warnings)), statusLabels},
		})

		for j := range s.Connections {
			conn := &s.Connections[j]

			connected := 0.0
			if conn.Connected.Val {
				connected = 1.0
			}

			connLabels := []string{m.SiteName, m.SourceName, m.Name, conn.ConnectionID, s.SiteID}

			r.send([]*metric{
				{u.VPNMesh.TunnelConnected, gauge, connected, connLabels},
				{u.VPNMesh.TunnelAssociationTime, gauge, conn.AssociationTime.Val, connLabels},
				{u.VPNMesh.TunnelErrors, gauge, float64(len(conn.Errors)), connLabels},
			})
		}
	}
}
