package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchWAN generates WAN configuration datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchWAN(r report, w *unifi.WANEnrichedConfiguration) {
	if w == nil {
		return
	}

	cfg := w.Configuration
	stats := w.Statistics
	details := w.Details

	tags := map[string]string{
		"wan_id":               cfg.ID,
		"wan_name":             cfg.Name,
		"wan_networkgroup":     cfg.WANNetworkgroup,
		"wan_type":             cfg.WANType,
		"wan_load_balance_type": cfg.WANLoadBalanceType,
		"isp_name":             details.ServiceProvider.Name,
		"isp_city":             details.ServiceProvider.City,
	}

	// Convert boolean FlexBool values to int for InfluxDB
	smartQEnabled := 0
	if cfg.WANSmartqEnabled.Val {
		smartQEnabled = 1
	}

	magicEnabled := 0
	if cfg.WANMagicEnabled.Val {
		magicEnabled = 1
	}

	vlanEnabled := 0
	if cfg.WANVlanEnabled.Val {
		vlanEnabled = 1
	}

	fields := map[string]any{
		// Configuration
		"failover_priority":      cfg.WANFailoverPriority.Val,
		"load_balance_weight":    cfg.WANLoadBalanceWeight.Val,
		"provider_download_kbps": cfg.WANProviderCapabilities.DownloadKbps.Val,
		"provider_upload_kbps":   cfg.WANProviderCapabilities.UploadKbps.Val,
		"smartq_enabled":         smartQEnabled,
		"magic_enabled":          magicEnabled,
		"vlan_enabled":           vlanEnabled,
		// Statistics
		"uptime_percentage":      stats.UptimePercentage,
		"peak_download_percent":  stats.PeakUsage.DownloadPercentage,
		"peak_upload_percent":    stats.PeakUsage.UploadPercentage,
		"max_rx_bytes_rate":      stats.PeakUsage.MaxRxBytesR.Val,
		"max_tx_bytes_rate":      stats.PeakUsage.MaxTxBytesR.Val,
		// Service Provider
		"service_provider_asn":   details.ServiceProvider.ASN.Val,
		// Metadata
		"creation_timestamp":     details.CreationTimestamp.Val,
	}

	r.send(&metric{Table: "wan", Tags: tags, Fields: fields})
}
