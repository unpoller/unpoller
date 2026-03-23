package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchTopology generates topology datapoints for Datadog.
func (u *DatadogUnifi) batchTopology(r report, t *unifi.Topology) {
	if t == nil {
		return
	}

	metricName := metricNamespace("topology")

	siteTags := []string{
		tag("site_name", t.SiteName),
		tag("source", t.SourceName),
	}

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

		edgeTags := []string{
			tag("uplink_mac", e.UplinkMac),
			tag("downlink_mac", e.DownlinkMac),
			tag("link_type", e.Type),
			tag("site_name", t.SiteName),
			tag("source", t.SourceName),
		}

		switch e.Type {
		case "WIRED":
			wired++

			if e.Duplex == "FULL_DUPLEX" {
				fullDuplex++
			}

			_ = r.reportGauge(metricName("link_rate_mbps"), e.RateMbps.Val, edgeTags)

		case "WIRELESS":
			wireless++

			if e.RadioBand != "" {
				bandCounts[e.RadioBand]++
			}

			if e.ExperienceScore.Val > 0 {
				_ = r.reportGauge(metricName("link_experience_score"), e.ExperienceScore.Val, edgeTags)
			}
		}
	}

	summary := map[string]float64{
		"vertices_total":       float64(len(t.Vertices)),
		"edges_total":          float64(len(t.Edges)),
		"devices_total":        float64(devices),
		"clients_total":        float64(clients),
		"connections_wired":    float64(wired),
		"connections_wireless": float64(wireless),
		"wired_full_duplex":    float64(fullDuplex),
		"has_unknown_switch":   unknownSwitch,
	}

	for name, value := range summary {
		_ = r.reportGauge(metricName(name), value, siteTags)
	}

	for band, count := range bandCounts {
		bandTags := []string{
			tag("band", band),
			tag("site_name", t.SiteName),
			tag("source", t.SourceName),
		}
		_ = r.reportGauge(metricName("connections_by_band"), float64(count), bandTags)
	}
}
