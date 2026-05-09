package promunifi

import "github.com/prometheus/client_golang/prometheus"

// NewCollectorForTesting creates a promUnifi with all metric descriptors
// initialized for the given namespace prefix, for use in tests.
func NewCollectorForTesting(namespace string) prometheus.Collector {
	ns := namespace

	u := &promUnifi{Config: &Config{Namespace: ns}}

	u.Client = descClient(ns + "_client_")
	u.Device = descDevice(ns + "_device_")
	u.UAP = descUAP(ns + "_device_")
	u.USG = descUSG(ns + "_device_")
	u.USW = descUSW(ns + "_device_")
	u.PDU = descPDU(ns + "_device_")
	u.Site = descSite(ns + "_site_")
	u.SpeedTest = descSpeedTest(ns + "_speedtest_")
	u.DHCPLease = descDHCPLease(ns + "_")
	u.WAN = descWAN(ns + "_")
	u.FirewallPolicy = descFirewallPolicy(ns + "_")
	u.Topology = descTopology(ns + "_")
	u.PortAnomaly = descPortAnomaly(ns + "_")
	u.VPNMesh = descVPNMesh(ns + "_")
	u.IntegrationDevice = descIntegrationDevice(ns + "_integration_device_")
	u.WANStatus = descWANStatus(ns + "_")
	u.PortForward = descPortForward(ns + "_")
	u.SSLCertificate = descSSLCertificate(ns + "_")
	u.UPSDevice = descUPSDevice(ns + "_")
	u.WifiBroadcast = descWifiBroadcast(ns + "_")
	u.FirewallZone = descFirewallZone(ns + "_")
	u.ACLRule = descACLRule(ns + "_")
	u.VPNServer = descVPNServer(ns + "_")
	u.SiteToSiteTunnel = descSiteToSiteTunnel(ns + "_")
	u.LAG = descLAG(ns + "_")
	u.MCLAGDomain = descMCLAGDomain(ns + "_")
	u.SwitchStack = descSwitchStack(ns + "_")
	u.DNSPolicy = descDNSPolicy(ns + "_")
	u.RADIUSProfile = descRADIUSProfile(ns + "_")
	u.TrafficMatchingList = descTrafficMatchingList(ns + "_")
	u.HotspotVoucher = descHotspotVoucher(ns + "_")
	u.DPIApplication = descDPIApplication(ns + "_")
	u.DPICategory = descDPICategory(ns + "_")
	u.PendingDevice = descPendingDevice(ns + "_")
	u.Country = descCountry(ns + "_")

	return u
}
