package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type speedtest struct {
	DownloadMbps *prometheus.Desc
	UploadMbps   *prometheus.Desc
	LatencyMs    *prometheus.Desc
	Timestamp    *prometheus.Desc
}

func descSpeedTest(ns string) *speedtest {
	labels := []string{"wan_interface", "wan_group", "site_name", "source"}

	return &speedtest{
		DownloadMbps: prometheus.NewDesc(ns+"download_mbps", "Speed Test Download in Mbps", labels, nil),
		UploadMbps:   prometheus.NewDesc(ns+"upload_mbps", "Speed Test Upload in Mbps", labels, nil),
		LatencyMs:    prometheus.NewDesc(ns+"latency_ms", "Speed Test Latency in milliseconds", labels, nil),
		Timestamp:    prometheus.NewDesc(ns+"timestamp_seconds", "Speed Test Timestamp (Unix epoch)", labels, nil),
	}
}

func (u *promUnifi) exportSpeedTest(r report, st *unifi.SpeedTestResult) {
	if st == nil {
		return
	}

	labels := []string{st.InterfaceName, st.WANNetworkGroup, st.SiteName, st.SourceName}

	r.send([]*metric{
		{u.SpeedTest.DownloadMbps, gauge, st.DownloadMbps, labels},
		{u.SpeedTest.UploadMbps, gauge, st.UploadMbps, labels},
		{u.SpeedTest.LatencyMs, gauge, st.LatencyMs, labels},
		{u.SpeedTest.Timestamp, gauge, st.Time, labels},
	})
}
