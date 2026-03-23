package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchTopology generates topology datapoints for InfluxDB.
func (u *InfluxUnifi) batchTopology(r report, t *unifi.Topology) {
	if t == nil {
		return
	}

	var (
		devices    int
		clients    int
		wired      int
		wireless   int
		fullDuplex int
	)

	unknownSwitch := 0
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

		edgeTags := map[string]string{
			"uplink_mac":   e.UplinkMac,
			"downlink_mac": e.DownlinkMac,
			"link_type":    e.Type,
			"site_name":    t.SiteName,
			"source":       t.SourceName,
		}

		switch e.Type {
		case "WIRED":
			wired++

			if e.Duplex == "FULL_DUPLEX" {
				fullDuplex++
			}

			edgeFields := map[string]any{
				"rate_mbps": e.RateMbps.Val,
			}

			r.send(&metric{Table: "topology_edge", Tags: edgeTags, Fields: edgeFields})

		case "WIRELESS":
			wireless++

			edgeTags["essid"] = e.Essid
			edgeTags["radio_band"] = e.RadioBand
			edgeTags["protocol"] = e.Protocol

			edgeFields := map[string]any{
				"experience_score": e.ExperienceScore.Val,
				"channel":          e.Channel.Val,
			}

			r.send(&metric{Table: "topology_edge", Tags: edgeTags, Fields: edgeFields})
		}
	}

	summaryTags := map[string]string{
		"site_name": t.SiteName,
		"source":    t.SourceName,
	}

	summaryFields := map[string]any{
		"vertices_total":     len(t.Vertices),
		"edges_total":        len(t.Edges),
		"devices_total":      devices,
		"clients_total":      clients,
		"connections_wired":  wired,
		"connections_wireless": wireless,
		"wired_full_duplex":  fullDuplex,
		"has_unknown_switch": unknownSwitch,
	}

	r.send(&metric{Table: "topology_summary", Tags: summaryTags, Fields: summaryFields})
}
