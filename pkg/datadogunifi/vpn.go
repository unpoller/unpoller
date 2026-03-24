package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchMagicSiteToSiteVPN generates Site Magic VPN datapoints for Datadog.
func (u *DatadogUnifi) batchMagicSiteToSiteVPN(r report, m *unifi.MagicSiteToSiteVPN) {
	if m == nil {
		return
	}

	meshMetric := metricNamespace("vpn_mesh")

	meshTags := []string{
		tag("site_name", m.SiteName),
		tag("source", m.SourceName),
		tag("mesh_name", m.Name),
	}

	paused := 0.0
	if m.Pause.Val {
		paused = 1.0
	}

	_ = r.reportGauge(meshMetric("paused"), paused, meshTags)
	_ = r.reportGauge(meshMetric("connections_total"), float64(len(m.Connections)), meshTags)
	_ = r.reportGauge(meshMetric("devices_total"), float64(len(m.Devices)), meshTags)

	tunnelMetric := metricNamespace("vpn_tunnel")
	statusMetric := metricNamespace("vpn_mesh_status")

	for i := range m.Status {
		s := &m.Status[i]

		statusTags := []string{
			tag("site_name", m.SiteName),
			tag("source", m.SourceName),
			tag("mesh_name", m.Name),
			tag("status_site", s.SiteID),
		}

		_ = r.reportGauge(statusMetric("errors"), float64(len(s.Errors)), statusTags)
		_ = r.reportGauge(statusMetric("warnings"), float64(len(s.Warnings)), statusTags)

		for j := range s.Connections {
			conn := &s.Connections[j]

			connected := 0.0
			if conn.Connected.Val {
				connected = 1.0
			}

			connTags := []string{
				tag("site_name", m.SiteName),
				tag("source", m.SourceName),
				tag("mesh_name", m.Name),
				tag("connection_id", conn.ConnectionID),
				tag("status_site", s.SiteID),
			}

			_ = r.reportGauge(tunnelMetric("connected"), connected, connTags)
			_ = r.reportGauge(tunnelMetric("association_time"), conn.AssociationTime.Val, connTags)
			_ = r.reportGauge(tunnelMetric("errors"), float64(len(conn.Errors)), connTags)
		}
	}
}
