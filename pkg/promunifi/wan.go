package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type wan struct {
	// WAN Configuration metrics
	FailoverPriority     *prometheus.Desc
	LoadBalanceWeight    *prometheus.Desc
	ProviderDownloadKbps *prometheus.Desc
	ProviderUploadKbps   *prometheus.Desc
	SmartQEnabled        *prometheus.Desc
	MagicEnabled         *prometheus.Desc
	VlanEnabled          *prometheus.Desc
	// WAN Statistics metrics
	UptimePercentage    *prometheus.Desc
	PeakDownloadPercent *prometheus.Desc
	PeakUploadPercent   *prometheus.Desc
	MaxRxBytesR         *prometheus.Desc
	MaxTxBytesR         *prometheus.Desc
	// WAN Service Provider metrics
	ServiceProviderASN *prometheus.Desc
	// WAN Creation timestamp
	CreationTimestamp *prometheus.Desc
}

func descWAN(ns string) *wan {
	labels := []string{
		"wan_id",
		"wan_name",
		"wan_networkgroup",
		"wan_type",
		"wan_load_balance_type",
		"site_name",
		"source",
	}

	providerLabels := []string{
		"wan_id",
		"wan_name",
		"wan_networkgroup",
		"isp_name",
		"isp_city",
		"site_name",
		"source",
	}

	nd := prometheus.NewDesc

	return &wan{
		// Configuration
		FailoverPriority:     nd(ns+"wan_failover_priority", "WAN failover priority (lower is higher priority)", labels, nil),
		LoadBalanceWeight:    nd(ns+"wan_load_balance_weight", "WAN load balancing weight", labels, nil),
		ProviderDownloadKbps: nd(ns+"wan_provider_download_kbps", "Configured ISP download speed in Kbps", labels, nil),
		ProviderUploadKbps:   nd(ns+"wan_provider_upload_kbps", "Configured ISP upload speed in Kbps", labels, nil),
		SmartQEnabled:        nd(ns+"wan_smartq_enabled", "SmartQueue QoS enabled (1) or disabled (0)", labels, nil),
		MagicEnabled:         nd(ns+"wan_magic_enabled", "Magic WAN enabled (1) or disabled (0)", labels, nil),
		VlanEnabled:          nd(ns+"wan_vlan_enabled", "VLAN enabled for WAN (1) or disabled (0)", labels, nil),
		// Statistics
		UptimePercentage:    nd(ns+"wan_uptime_percentage", "WAN uptime percentage", labels, nil),
		PeakDownloadPercent: nd(ns+"wan_peak_download_percent", "Peak download usage as percentage of configured capacity", labels, nil),
		PeakUploadPercent:   nd(ns+"wan_peak_upload_percent", "Peak upload usage as percentage of configured capacity", labels, nil),
		MaxRxBytesR:         nd(ns+"wan_max_rx_bytes_rate", "Maximum receive bytes rate", labels, nil),
		MaxTxBytesR:         nd(ns+"wan_max_tx_bytes_rate", "Maximum transmit bytes rate", labels, nil),
		// Service Provider
		ServiceProviderASN: nd(ns+"wan_service_provider_asn", "Service provider autonomous system number", providerLabels, nil),
		// Creation
		CreationTimestamp: nd(ns+"wan_creation_timestamp", "WAN configuration creation timestamp", labels, nil),
	}
}

func (u *promUnifi) exportWAN(r report, w *unifi.WANEnrichedConfiguration) {
	if w == nil {
		return
	}

	cfg := w.Configuration
	stats := w.Statistics
	details := w.Details

	// Base labels
	labels := []string{
		cfg.ID,
		cfg.Name,
		cfg.WANNetworkgroup,
		cfg.WANType,
		cfg.WANLoadBalanceType,
		"", // site_name - will be set by caller if available
		"", // source - will be set by caller if available
	}

	// Convert boolean FlexBool values to float64
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

	metrics := []*metric{
		{u.WAN.FailoverPriority, gauge, cfg.WANFailoverPriority.Val, labels},
		{u.WAN.LoadBalanceWeight, gauge, cfg.WANLoadBalanceWeight.Val, labels},
		{u.WAN.ProviderDownloadKbps, gauge, cfg.WANProviderCapabilities.DownloadKbps.Val, labels},
		{u.WAN.ProviderUploadKbps, gauge, cfg.WANProviderCapabilities.UploadKbps.Val, labels},
		{u.WAN.SmartQEnabled, gauge, smartQEnabled, labels},
		{u.WAN.MagicEnabled, gauge, magicEnabled, labels},
		{u.WAN.VlanEnabled, gauge, vlanEnabled, labels},
		{u.WAN.UptimePercentage, gauge, stats.UptimePercentage, labels},
		{u.WAN.PeakDownloadPercent, gauge, stats.PeakUsage.DownloadPercentage, labels},
		{u.WAN.PeakUploadPercent, gauge, stats.PeakUsage.UploadPercentage, labels},
		{u.WAN.MaxRxBytesR, gauge, stats.PeakUsage.MaxRxBytesR.Val, labels},
		{u.WAN.MaxTxBytesR, gauge, stats.PeakUsage.MaxTxBytesR.Val, labels},
		{u.WAN.CreationTimestamp, gauge, details.CreationTimestamp.Val, labels},
	}

	// Service provider info (uses different labels)
	providerLabels := []string{
		cfg.ID,
		cfg.Name,
		cfg.WANNetworkgroup,
		details.ServiceProvider.Name,
		details.ServiceProvider.City,
		"", // site_name
		"", // source
	}

	metrics = append(metrics, &metric{u.WAN.ServiceProviderASN, gauge, details.ServiceProvider.ASN.Val, providerLabels})

	r.send(metrics)
}
