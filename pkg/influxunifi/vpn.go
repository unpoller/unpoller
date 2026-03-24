package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchMagicSiteToSiteVPN generates Site Magic VPN datapoints for InfluxDB.
func (u *InfluxUnifi) batchMagicSiteToSiteVPN(r report, m *unifi.MagicSiteToSiteVPN) {
	if m == nil {
		return
	}

	meshTags := map[string]string{
		"site_name": m.SiteName,
		"source":    m.SourceName,
		"mesh_id":   m.ID,
		"mesh_name": m.Name,
	}

	paused := 0.0
	if m.Pause.Val {
		paused = 1.0
	}

	meshFields := map[string]any{
		"paused":            paused,
		"connections_total": len(m.Connections),
		"devices_total":     len(m.Devices),
	}

	r.send(&metric{Table: "vpn_mesh", Tags: meshTags, Fields: meshFields})

	for i := range m.Status {
		s := &m.Status[i]

		for j := range s.Connections {
			conn := &s.Connections[j]

			connected := 0.0
			if conn.Connected.Val {
				connected = 1.0
			}

			connTags := map[string]string{
				"site_name":     m.SiteName,
				"source":        m.SourceName,
				"mesh_name":     m.Name,
				"connection_id": conn.ConnectionID,
				"status_site":   s.SiteID,
			}

			connFields := map[string]any{
				"connected":        connected,
				"association_time": conn.AssociationTime.Val,
				"errors":           len(conn.Errors),
			}

			r.send(&metric{Table: "vpn_mesh_connection", Tags: connTags, Fields: connFields})
		}

		statusTags := map[string]string{
			"site_name":   m.SiteName,
			"source":      m.SourceName,
			"mesh_name":   m.Name,
			"status_site": s.SiteID,
		}

		statusFields := map[string]any{
			"errors":   len(s.Errors),
			"warnings": len(s.Warnings),
		}

		r.send(&metric{Table: "vpn_mesh_status", Tags: statusTags, Fields: statusFields})
	}
}
