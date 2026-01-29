package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type dhcplease struct {
	// Network-level pool metrics (exported once per network)
	ActiveLeases       *prometheus.Desc
	PoolSize           *prometheus.Desc
	UtilizationPercent *prometheus.Desc
	FreePercent        *prometheus.Desc
	AvailableIPs       *prometheus.Desc
	// Per-lease metrics
	LeaseStart *prometheus.Desc
	LeaseEnd   *prometheus.Desc
	LeaseTime  *prometheus.Desc
	IsStatic   *prometheus.Desc
}

func descDHCPLease(ns string) *dhcplease {
	// Network-level labels (for pool metrics)
	networkLabels := []string{
		"network",
		"network_id",
		"site_name",
		"source",
	}
	// Per-lease labels
	leaseLabels := []string{
		"ip",
		"mac",
		"hostname",
		"network",
		"network_id",
		"client_name",
		"site_name",
		"source",
	}
	nd := prometheus.NewDesc

	return &dhcplease{
		ActiveLeases:       nd(ns+"dhcp_active_leases", "Number of active DHCP leases for this network", networkLabels, nil),
		PoolSize:           nd(ns+"dhcp_pool_size", "Total number of IPs in DHCP pool range", networkLabels, nil),
		UtilizationPercent: nd(ns+"dhcp_utilization_percent", "DHCP pool utilization percentage (used)", networkLabels, nil),
		FreePercent:        nd(ns+"dhcp_free_percent", "DHCP pool free percentage (available)", networkLabels, nil),
		AvailableIPs:       nd(ns+"dhcp_available_ips", "Number of available IPs in DHCP pool", networkLabels, nil),
		LeaseStart:         nd(ns+"dhcp_lease_start", "DHCP lease start timestamp", leaseLabels, nil),
		LeaseEnd:           nd(ns+"dhcp_lease_end", "DHCP lease end timestamp", leaseLabels, nil),
		LeaseTime:          nd(ns+"dhcp_lease_time", "DHCP lease duration in seconds", leaseLabels, nil),
		IsStatic:           nd(ns+"dhcp_is_static", "Whether this is a static DHCP lease (1) or dynamic (0)", leaseLabels, nil),
	}
}

func (u *promUnifi) exportDHCPLease(r report, l *unifi.DHCPLease) {
	// Per-lease labels
	leaseLabels := []string{
		l.IP,
		l.Mac,
		l.Hostname,
		l.Network,
		l.NetworkID,
		l.ClientName,
		l.SiteName,
		l.SourceName,
	}

	// Convert FlexBool to float64 (1.0 for true, 0.0 for false)
	isStaticVal := 0.0
	if l.IsStatic.Val {
		isStaticVal = 1.0
	}

	metrics := []*metric{
		{u.DHCPLease.IsStatic, gauge, isStaticVal, leaseLabels},
	}

	// Add lease time metrics if available
	if l.LeaseStart.Val > 0 {
		metrics = append(metrics, &metric{u.DHCPLease.LeaseStart, gauge, l.LeaseStart.Val, leaseLabels})
	}

	if l.LeaseEnd.Val > 0 {
		metrics = append(metrics, &metric{u.DHCPLease.LeaseEnd, gauge, l.LeaseEnd.Val, leaseLabels})
	}

	if l.LeaseTime.Val > 0 {
		metrics = append(metrics, &metric{u.DHCPLease.LeaseTime, gauge, l.LeaseTime.Val, leaseLabels})
	}

	r.send(metrics)
}

// exportDHCPNetworkPool exports network-level DHCP pool metrics (once per network).
func (u *promUnifi) exportDHCPNetworkPool(r report, leases []*unifi.DHCPLease) {
	// Group leases by network_id to export pool metrics once per network
	networkMetrics := make(map[string]*networkPoolData)

	for _, lease := range leases {
		if lease.NetworkTableEntry == nil {
			continue
		}

		networkID := lease.NetworkID
		if networkID == "" {
			continue
		}

		// Use the first lease for each network to get pool data
		if _, exists := networkMetrics[networkID]; !exists {
			poolSize := lease.GetPoolSize()
			if poolSize > 0 {
				networkMetrics[networkID] = &networkPoolData{
					Network:      lease.Network,
					NetworkID:    networkID,
					SiteName:     lease.SiteName,
					SourceName:   lease.SourceName,
					PoolSize:     poolSize,
					ActiveLeases: lease.GetActiveLeaseCount(),
					Utilization:  lease.GetUtilizationPercentage(),
					FreePercent:  100.0 - lease.GetUtilizationPercentage(),
					AvailableIPs: lease.GetAvailableIPs(),
				}
			}
		}
	}

	// Export metrics for each unique network
	for _, data := range networkMetrics {
		networkLabels := []string{
			data.Network,
			data.NetworkID,
			data.SiteName,
			data.SourceName,
		}

		r.send([]*metric{
			{u.DHCPLease.PoolSize, gauge, float64(data.PoolSize), networkLabels},
			{u.DHCPLease.ActiveLeases, gauge, float64(data.ActiveLeases), networkLabels},
			{u.DHCPLease.UtilizationPercent, gauge, data.Utilization, networkLabels},
			{u.DHCPLease.FreePercent, gauge, data.FreePercent, networkLabels},
			{u.DHCPLease.AvailableIPs, gauge, float64(data.AvailableIPs), networkLabels},
		})
	}
}

type networkPoolData struct {
	Network      string
	NetworkID    string
	SiteName     string
	SourceName   string
	PoolSize     int
	ActiveLeases int
	Utilization  float64
	FreePercent  float64
	AvailableIPs int
}
