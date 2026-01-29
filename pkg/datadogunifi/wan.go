package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchWAN generates WAN configuration datapoints for Datadog.
// These points can be passed directly to Datadog.
func (u *DatadogUnifi) batchWAN(r report, w *unifi.WANEnrichedConfiguration) {
	if w == nil {
		return
	}

	metricName := metricNamespace("wan")

	cfg := w.Configuration
	stats := w.Statistics
	details := w.Details

	tags := []string{
		tag("wan_id", cfg.ID),
		tag("wan_name", cfg.Name),
		tag("wan_networkgroup", cfg.WANNetworkgroup),
		tag("wan_type", cfg.WANType),
		tag("wan_load_balance_type", cfg.WANLoadBalanceType),
		tag("isp_name", details.ServiceProvider.Name),
		tag("isp_city", details.ServiceProvider.City),
	}

	// Convert boolean FlexBool values to float64 for Datadog
	smartQEnabled := 0.0
	if cfg.WANSmartqEnabled.Val {
		smartQEnabled = 1.0
	}

	magicEnabled := 0.0
	if cfg.WANMagicEnabled.Val {
		magicEnabled = 1.0
	}

	vlanEnabled := 0.0
	if cfg.WANVlanEnabled.Val {
		vlanEnabled = 1.0
	}

	data := map[string]float64{
		// Configuration
		"failover_priority":      cfg.WANFailoverPriority.Val,
		"load_balance_weight":    cfg.WANLoadBalanceWeight.Val,
		"provider_download_kbps": cfg.WANProviderCapabilities.DownloadKbps.Val,
		"provider_upload_kbps":   cfg.WANProviderCapabilities.UploadKbps.Val,
		"smartq_enabled":         smartQEnabled,
		"magic_enabled":          magicEnabled,
		"vlan_enabled":           vlanEnabled,
		// Statistics
		"uptime_percentage":     stats.UptimePercentage,
		"peak_download_percent": stats.PeakUsage.DownloadPercentage,
		"peak_upload_percent":   stats.PeakUsage.UploadPercentage,
		"max_rx_bytes_rate":     stats.PeakUsage.MaxRxBytesR.Val,
		"max_tx_bytes_rate":     stats.PeakUsage.MaxTxBytesR.Val,
		// Service Provider
		"service_provider_asn": details.ServiceProvider.ASN.Val,
		// Metadata
		"creation_timestamp": details.CreationTimestamp.Val,
	}

	for name, value := range data {
		_ = r.reportGauge(metricName(name), value, tags)
	}
}
