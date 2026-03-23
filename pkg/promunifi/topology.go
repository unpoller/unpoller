package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type topology struct {
	// Summary metrics
	VerticesTotal    *prometheus.Desc
	EdgesTotal       *prometheus.Desc
	DevicesTotal     *prometheus.Desc
	ClientsTotal     *prometheus.Desc
	HasUnknownSwitch *prometheus.Desc

	// Connection type metrics
	ConnectionsWired    *prometheus.Desc
	ConnectionsWireless *prometheus.Desc
	ConnectionsByBand   *prometheus.Desc

	// Link quality metrics
	LinkExperienceScore *prometheus.Desc
	LinkRateMbps        *prometheus.Desc
	WiredFullDuplex     *prometheus.Desc
}

func descTopology(ns string) *topology {
	siteLabels := []string{"site_name", "source"}
	linkLabels := []string{"uplink_mac", "downlink_mac", "link_type", "site_name", "source"}
	bandLabels := []string{"band", "site_name", "source"}

	nd := prometheus.NewDesc

	return &topology{
		VerticesTotal:       nd(ns+"topology_vertices_total", "Total vertices in topology", siteLabels, nil),
		EdgesTotal:          nd(ns+"topology_edges_total", "Total edges/connections in topology", siteLabels, nil),
		DevicesTotal:        nd(ns+"topology_devices_total", "UniFi devices in topology", siteLabels, nil),
		ClientsTotal:        nd(ns+"topology_clients_total", "Clients in topology", siteLabels, nil),
		HasUnknownSwitch:    nd(ns+"topology_has_unknown_switch", "Unknown switch detected in topology (1/0)", siteLabels, nil),
		ConnectionsWired:    nd(ns+"topology_connections_wired", "Number of wired connections", siteLabels, nil),
		ConnectionsWireless: nd(ns+"topology_connections_wireless", "Number of wireless connections", siteLabels, nil),
		ConnectionsByBand:   nd(ns+"topology_connections_by_band", "Number of wireless connections by radio band", bandLabels, nil),
		LinkExperienceScore: nd(ns+"topology_link_experience_score", "Link experience score (0-100)", linkLabels, nil),
		LinkRateMbps:        nd(ns+"topology_link_rate_mbps", "Link rate in Mbps", linkLabels, nil),
		WiredFullDuplex:     nd(ns+"topology_wired_full_duplex", "Number of full-duplex wired links", siteLabels, nil),
	}
}

func (u *promUnifi) exportTopology(r report, t *unifi.Topology) {
	if t == nil {
		return
	}

	siteLabels := []string{t.SiteName, t.SourceName}

	var (
		devices     int
		clients     int
		wired       int
		wireless    int
		fullDuplex  int
		bandCounts  = make(map[string]int)
		unknownSwitch float64
	)

	if t.HasUnknownSwitch {
		unknownSwitch = 1
	}

	for i := range t.Vertices {
		switch t.Vertices[i].Type {
		case "DEVICE":
			devices++
		case "CLIENT":
			clients++
		}
	}

	for i := range t.Edges {
		e := &t.Edges[i]
		linkLabels := []string{e.UplinkMac, e.DownlinkMac, e.Type, t.SiteName, t.SourceName}

		switch e.Type {
		case "WIRED":
			wired++

			if e.Duplex == "FULL_DUPLEX" {
				fullDuplex++
			}

			if e.RateMbps.Val > 0 {
				r.send([]*metric{{u.Topology.LinkRateMbps, gauge, e.RateMbps.Val, linkLabels}})
			}
		case "WIRELESS":
			wireless++

			if e.RadioBand != "" {
				bandCounts[e.RadioBand]++
			}

			if e.ExperienceScore.Val > 0 {
				r.send([]*metric{{u.Topology.LinkExperienceScore, gauge, e.ExperienceScore.Val, linkLabels}})
			}
		}
	}

	r.send([]*metric{
		{u.Topology.VerticesTotal, gauge, float64(len(t.Vertices)), siteLabels},
		{u.Topology.EdgesTotal, gauge, float64(len(t.Edges)), siteLabels},
		{u.Topology.DevicesTotal, gauge, float64(devices), siteLabels},
		{u.Topology.ClientsTotal, gauge, float64(clients), siteLabels},
		{u.Topology.HasUnknownSwitch, gauge, unknownSwitch, siteLabels},
		{u.Topology.ConnectionsWired, gauge, float64(wired), siteLabels},
		{u.Topology.ConnectionsWireless, gauge, float64(wireless), siteLabels},
		{u.Topology.WiredFullDuplex, gauge, float64(fullDuplex), siteLabels},
	})

	for band, count := range bandCounts {
		r.send([]*metric{{u.Topology.ConnectionsByBand, gauge, float64(count), []string{band, t.SiteName, t.SourceName}}})
	}
}
