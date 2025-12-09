package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchSpeedTest generates Unifi Speed Test datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchSpeedTest(r report, st *unifi.SpeedTestResult) {
	if st == nil {
		return
	}

	tags := map[string]string{
		"site_name":       st.SiteName,
		"source":          st.SourceName,
		"wan_interface":   st.InterfaceName,
		"wan_group":       st.WANNetworkGroup,
		"network_conf_id": st.NetworkConfID,
	}

	fields := map[string]any{
		"download_mbps": st.DownloadMbps.Val,
		"upload_mbps":   st.UploadMbps.Val,
		"latency_ms":    st.LatencyMs.Val,
		"timestamp":     st.Time.Val,
	}

	if st.WANProviderCapabilities != nil {
		fields["provider_download_kbps"] = st.WANProviderCapabilities.DownloadKbps.Val
		fields["provider_upload_kbps"] = st.WANProviderCapabilities.UploadKbps.Val
	}

	r.send(&metric{Table: "speedtest", Tags: tags, Fields: fields})
}
