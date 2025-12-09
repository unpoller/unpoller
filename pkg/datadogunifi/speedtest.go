package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchSpeedTest generates Unifi Speed Test datapoints for Datadog.
// These points can be passed directly to Datadog.
func (u *DatadogUnifi) batchSpeedTest(r report, st *unifi.SpeedTestResult) {
	if st == nil {
		return
	}

	metricName := metricNamespace("speedtest")

	tags := []string{
		tag("site_name", st.SiteName),
		tag("source", st.SourceName),
		tag("wan_interface", st.InterfaceName),
		tag("wan_group", st.WANNetworkGroup),
		tag("network_conf_id", st.NetworkConfID),
	}

	data := map[string]float64{
		"download_mbps": st.DownloadMbps.Val,
		"upload_mbps":   st.UploadMbps.Val,
		"latency_ms":    st.LatencyMs.Val,
		"timestamp":     st.Time.Val,
	}

	if st.WANProviderCapabilities != nil {
		data["provider_download_kbps"] = st.WANProviderCapabilities.DownloadKbps.Val
		data["provider_upload_kbps"] = st.WANProviderCapabilities.UploadKbps.Val
	}

	for name, value := range data {
		_ = r.reportGauge(metricName(name), value, tags)
	}
}
