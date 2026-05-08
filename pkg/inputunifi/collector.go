package inputunifi

// nolint: gosec
import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
)

const (
	historySeconds = 86400
	pollDuration   = time.Second * historySeconds
)

var ErrScrapeFilterMatchFailed = fmt.Errorf("scrape filter match failed, and filter is not http URL")

func (u *InputUnifi) isNill(c *Controller) bool {
	u.RLock()
	defer u.RUnlock()

	return c.Unifi == nil
}

// newDynamicCntrlr creates and saves a controller definition for further use.
// This is called when an unconfigured controller is requested.
func (u *InputUnifi) newDynamicCntrlr(url string) (bool, *Controller) {
	u.Lock()
	defer u.Unlock()

	if c := u.dynamic[url]; c != nil {
		// it already exists.
		return false, c
	}

	ccopy := u.Default // copy defaults into new controller
	u.dynamic[url] = &ccopy
	u.dynamic[url].URL = url

	return true, u.dynamic[url]
}

func (u *InputUnifi) dynamicController(filter *poller.Filter) (*poller.Metrics, error) {
	if !strings.HasPrefix(filter.Path, "http") {
		return nil, ErrScrapeFilterMatchFailed
	}

	newCntrlr, c := u.newDynamicCntrlr(filter.Path)

	if newCntrlr {
		u.Logf("Authenticating to Dynamic UniFi Controller: %s", filter.Path)

		if err := u.getUnifi(c); err != nil {
			u.logController(c)

			return nil, fmt.Errorf("authenticating to %s: %w", filter.Path, err)
		}

		u.logController(c)
	}

	return u.collectController(c)
}

func (u *InputUnifi) collectController(c *Controller) (*poller.Metrics, error) {
	u.LogDebugf("Collecting controller data: %s (%s)", c.URL, c.ID)

	if u.isNill(c) {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return nil, fmt.Errorf("re-authenticating to %s: %w", c.URL, err)
		}
	}

	metrics, err := u.pollController(c)
	if err != nil {
		u.Logf("Re-authenticating to UniFi Controller %s (poll error: %v)", c.URL, err)

		if authErr := u.getUnifi(c); authErr != nil {
			return metrics, fmt.Errorf("re-authenticating to %s: %w", c.URL, authErr)
		}

		// Brief delay to allow controller to process new authentication
		time.Sleep(500 * time.Millisecond)

		// Retry the poll after successful re-authentication
		u.LogDebugf("Retrying poll after re-authentication: %s", c.URL)
		metrics, err = u.pollController(c)
	}

	return metrics, err
}

//nolint:cyclop
func (u *InputUnifi) pollController(c *Controller) (*poller.Metrics, error) {
	u.RLock()
	defer u.RUnlock()

	if c == nil {
		return nil, fmt.Errorf("controller is nil")
	}

	if c.Unifi == nil {
		return nil, fmt.Errorf("controller client is nil (e.g. after 429 or auth failure): %s", c.URL)
	}

	u.LogDebugf("Polling controller: %s (%s)", c.URL, c.ID)

	// Get the sites we care about.
	sites, err := u.getFilteredSites(c)
	if err != nil {
		return nil, fmt.Errorf("unifi.GetSites(): %w", err)
	}

	m := &Metrics{TS: time.Now(), Sites: sites}

	// FIXME needs to be last poll time maybe
	st := m.TS.Add(-1 * pollDuration)
	tp := unifi.EpochMillisTimePeriod{StartEpochMillis: st.UnixMilli(), EndEpochMillis: m.TS.UnixMilli()}

	if c.SaveRogue != nil && *c.SaveRogue {
		if m.RogueAPs, err = c.Unifi.GetRogueAPs(sites); err != nil {
			return nil, fmt.Errorf("unifi.GetRogueAPs(%s): %w", c.URL, err)
		}

		u.LogDebugf("Found %d RogueAPs entries", len(m.RogueAPs))
	}

	if c.SaveDPI != nil && *c.SaveDPI {
		if m.SitesDPI, err = c.Unifi.GetSiteDPI(sites); err != nil {
			return nil, fmt.Errorf("unifi.GetSiteDPI(%s): %w", c.URL, err)
		}

		u.LogDebugf("Found %d SitesDPI entries", len(m.SitesDPI))

		if m.ClientsDPI, err = c.Unifi.GetClientsDPI(sites); err != nil {
			return nil, fmt.Errorf("unifi.GetClientsDPI(%s): %w", c.URL, err)
		}

		u.LogDebugf("Found %d ClientsDPI entries", len(m.ClientsDPI))
	}

	if c.SaveTraffic != nil && *c.SaveTraffic {
		if m.CountryTraffic, err = c.Unifi.GetCountryTraffic(sites, &tp); err != nil {
			return nil, fmt.Errorf("unifi.GetCountryTraffic(%s): %w", c.URL, err)
		}

		u.LogDebugf("Found %d CountryTraffic entries", len(m.CountryTraffic))
	}

	if c.SaveDPI != nil && *c.SaveDPI {
		// Supplement DPI data with the v2 traffic API, which works on newer firmware
		// (Network 9.1+) where the legacy /stat/stadpi and /stat/sitedpi endpoints
		// return empty results. GetClientTraffic is called regardless of SaveTraffic
		// because it provides DPI-equivalent per-client app/category breakdowns.
		clientUsageByApp, err := c.Unifi.GetClientTraffic(sites, &tp, true)
		if err != nil {
			u.LogDebugf("unifi.GetClientTraffic(%s): %v (legacy DPI endpoints will be used if available)", c.URL, err)
		} else {
			u.LogDebugf("Found %d ClientUsageByApp entries", len(clientUsageByApp))
			b4 := len(m.ClientsDPI)
			u.convertToClientDPI(clientUsageByApp, m)
			u.LogDebugf("Added %d ClientDPI entries from v2 traffic API for a total of %d", len(m.ClientsDPI)-b4, len(m.ClientsDPI))
			b4Sites := len(m.SitesDPI)
			u.convertToSiteDPI(clientUsageByApp, m)
			u.LogDebugf("Added %d SitesDPI entries from v2 traffic API for a total of %d", len(m.SitesDPI)-b4Sites, len(m.SitesDPI))
		}
	}

	// Get all the points.
	if m.Clients, err = c.Unifi.GetClients(sites); err != nil {
		return nil, fmt.Errorf("unifi.GetClients(%s): %w", c.URL, err)
	}

	u.LogDebugf("Found %d Clients entries", len(m.Clients))

	if m.Devices, err = c.Unifi.GetDevices(sites); err != nil {
		return nil, fmt.Errorf("unifi.GetDevices(%s): %w", c.URL, err)
	}

	u.LogDebugf("Found %d UBB, %d UXG, %d PDU, %d UCI, %d UDB, %d UAP %d USG %d USW %d UDM devices",
		len(m.Devices.UBBs), len(m.Devices.UXGs),
		len(m.Devices.PDUs), len(m.Devices.UCIs),
		len(m.Devices.UDBs), len(m.Devices.UAPs), len(m.Devices.USGs),
		len(m.Devices.USWs), len(m.Devices.UDMs))

	// Get speed test results for all WANs
	if m.SpeedTests, err = c.Unifi.GetSpeedTests(sites, historySeconds); err != nil {
		// Don't fail collection if speed tests fail - older controllers may not have this endpoint
		u.LogDebugf("unifi.GetSpeedTests(%s): %v (continuing)", c.URL, err)
	} else {
		u.LogDebugf("Found %d SpeedTests entries", len(m.SpeedTests))
	}

	// Get DHCP leases with associations.
	// Wrapped in recover so a nil-pointer panic in the library (e.g. when a 401 causes nil devices)
	// never crashes the poller. See https://github.com/unpoller/unpoller/issues/965
	func() {
		defer func() {
			if r := recover(); r != nil {
				u.LogErrorf("GetActiveDHCPLeasesWithAssociations panic recovered (see issue #965): %v", r)
			}
		}()

		if m.DHCPLeases, err = c.Unifi.GetActiveDHCPLeasesWithAssociations(sites); err != nil {
			// Don't fail collection if DHCP leases fail - older controllers may not have this endpoint
			u.LogDebugf("unifi.GetActiveDHCPLeasesWithAssociations(%s): %v (continuing)", c.URL, err)
		} else {
			u.LogDebugf("Found %d DHCPLeases entries", len(m.DHCPLeases))
		}
	}()

	// Get WAN enriched configuration
	if m.WANConfigs, err = c.Unifi.GetWANEnrichedConfiguration(sites); err != nil {
		// Don't fail collection if WAN config fails - older controllers may not have this endpoint
		u.LogDebugf("unifi.GetWANEnrichedConfiguration(%s): %v (continuing)", c.URL, err)
	} else {
		u.LogDebugf("Found %d WAN configuration entries", len(m.WANConfigs))
	}

	// Get firewall policies
	if m.FirewallPolicies, err = c.Unifi.GetFirewallPolicies(sites); err != nil {
		// Don't fail collection if firewall policies fail - older controllers may not have this endpoint
		u.LogDebugf("unifi.GetFirewallPolicies(%s): %v (continuing)", c.URL, err)
	} else {
		u.LogDebugf("Found %d FirewallPolicies entries", len(m.FirewallPolicies))
	}

	// Get controller system info (UniFi OS only)
	if m.Sysinfos, err = c.Unifi.GetSysinfo(sites); err != nil {
		// Don't fail collection if sysinfo fails - older controllers may not have this endpoint
		u.LogDebugf("unifi.GetSysinfo(%s): %v (continuing)", c.URL, err)
	} else {
		u.LogDebugf("Found %d Sysinfo entries", len(m.Sysinfos))
	}

	// Get network topology
	if m.Topologies, err = c.Unifi.GetTopology(sites); err != nil {
		// Don't fail collection if topology fails - older controllers may not have this endpoint
		u.LogDebugf("unifi.GetTopology(%s): %v (continuing)", c.URL, err)
	} else {
		u.LogDebugf("Found %d Topology entries", len(m.Topologies))
	}

	// Get port anomalies
	if m.PortAnomalies, err = c.Unifi.GetPortAnomalies(sites); err != nil {
		// Don't fail collection if port anomalies fail - older controllers may not have this endpoint
		u.LogDebugf("unifi.GetPortAnomalies(%s): %v (continuing)", c.URL, err)
	} else {
		u.LogDebugf("Found %d PortAnomalies entries", len(m.PortAnomalies))
	}

	// Get Site Magic site-to-site VPN mesh data
	if m.VPNMeshes, err = c.Unifi.GetMagicSiteToSiteVPN(sites); err != nil {
		// Don't fail collection if VPN data fails - older controllers may not have this endpoint
		u.LogDebugf("unifi.GetMagicSiteToSiteVPN(%s): %v (continuing)", c.URL, err)
	} else {
		u.LogDebugf("Found %d VPNMeshes entries", len(m.VPNMeshes))
	}

	// Legacy API additions (v5.26.0) — available on most firmware, no API key required.
	u.collectLegacyPerSite(c, sites, m)

	// Integration/v1 API additions (v5.26.0) — require API key and Network 9.3.43+.
	if c.APIKey != "" {
		u.collectIntegrationV1(c, sites, m)
	}

	// Update web UI only on success; call explicitly so we never run with nil c/c.Unifi (no defer).
	// Recover so a panic in updateWeb (e.g. old image, race) never kills the poller.
	if c != nil && c.Unifi != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					u.LogErrorf("updateWeb panic recovered (upgrade image if this persists): %v", r)
				}
			}()

			updateWeb(c, m)
		}()
	}

	return u.augmentMetrics(c, m), nil
}

// FIXME this would be better implemented on FlexInt itself
func (u *InputUnifi) intToFlexInt(i int) unifi.FlexInt {
	return unifi.FlexInt{
		Val: float64(i),
		Txt: fmt.Sprintf("%d", i),
	}
}

// FIXME this would be better implemented on FlexInt itself
func (u *InputUnifi) int64ToFlexInt(i int64) unifi.FlexInt {
	return unifi.FlexInt{
		Val: float64(i),
		Txt: fmt.Sprintf("%d", i),
	}
}

func (u *InputUnifi) convertToClientDPI(clientUsageByApp []*unifi.ClientUsageByApp, metrics *Metrics) {
	for _, client := range clientUsageByApp {
		byApp := make([]unifi.DPIData, 0)
		byCat := make([]unifi.DPIData, 0)

		type catCount struct {
			BytesReceived    int64
			BytesTransmitted int64
		}

		byCatMap := make(map[int]catCount)
		dpiClients := make([]*unifi.DPIClient, 0)
		// TODO create cat table
		for _, app := range client.UsageByApp {
			dpiData := unifi.DPIData{
				App:          u.intToFlexInt(app.Application),
				Cat:          u.intToFlexInt(app.Category),
				Clients:      dpiClients,
				KnownClients: u.intToFlexInt(0),
				RxBytes:      u.int64ToFlexInt(app.BytesReceived),
				RxPackets:    u.int64ToFlexInt(0), // We don't have packets from Unifi Controller
				TxBytes:      u.int64ToFlexInt(app.BytesTransmitted),
				TxPackets:    u.int64ToFlexInt(0), // We don't have packets from Unifi Controller
			}

			cat, ok := byCatMap[app.Category]
			if ok {
				cat.BytesReceived += app.BytesReceived
				cat.BytesTransmitted += app.BytesTransmitted
			} else {
				cat = catCount{
					BytesReceived:    app.BytesReceived,
					BytesTransmitted: app.BytesTransmitted,
				}
				byCatMap[app.Category] = cat
			}

			byApp = append(byApp, dpiData)
		}

		if len(byApp) <= 1 {
			byCat = byApp
		} else {
			for category, cat := range byCatMap {
				dpiData := unifi.DPIData{
					App:          u.intToFlexInt(16777215), // Unknown
					Cat:          u.intToFlexInt(category),
					Clients:      dpiClients,
					KnownClients: u.intToFlexInt(0),
					RxBytes:      u.int64ToFlexInt(cat.BytesReceived),
					RxPackets:    u.int64ToFlexInt(0), // We don't have packets from Unifi Controller
					TxBytes:      u.int64ToFlexInt(cat.BytesTransmitted),
					TxPackets:    u.int64ToFlexInt(0), // We don't have packets from Unifi Controller
				}
				byCat = append(byCat, dpiData)
			}
		}

		dpiTable := unifi.DPITable{
			ByApp:      byApp,
			ByCat:      byCat,
			MAC:        client.Client.Mac,
			Name:       client.Client.Name,
			SiteName:   client.TrafficSite.SiteName,
			SourceName: client.TrafficSite.SourceName,
		}
		metrics.ClientsDPI = append(metrics.ClientsDPI, &dpiTable)
	}
}

// convertToSiteDPI aggregates v2 client traffic data into per-site DPITable entries.
// It only adds a site entry if the site doesn't already have one from the legacy API,
// so old-firmware users are unaffected.
func (u *InputUnifi) convertToSiteDPI(clientUsageByApp []*unifi.ClientUsageByApp, metrics *Metrics) {
	// Build a set of sites already covered by the legacy API.
	existing := make(map[string]bool)

	for _, s := range metrics.SitesDPI {
		existing[s.SiteName] = true
	}

	type appKey struct {
		App int
		Cat int
	}

	type siteAgg struct {
		byApp      map[appKey]*unifi.DPIData
		byCat      map[int]*unifi.DPIData
		sourceName string
	}

	siteMap := make(map[string]*siteAgg)

	for _, client := range clientUsageByApp {
		siteName := client.TrafficSite.SiteName
		if existing[siteName] {
			continue
		}

		agg, ok := siteMap[siteName]
		if !ok {
			agg = &siteAgg{
				byApp:      make(map[appKey]*unifi.DPIData),
				byCat:      make(map[int]*unifi.DPIData),
				sourceName: client.TrafficSite.SourceName,
			}
			siteMap[siteName] = agg
		}

		for _, app := range client.UsageByApp {
			k := appKey{App: app.Application, Cat: app.Category}

			if d, ok := agg.byApp[k]; ok {
				d.RxBytes.Val += float64(app.BytesReceived)
				d.TxBytes.Val += float64(app.BytesTransmitted)
			} else {
				agg.byApp[k] = &unifi.DPIData{
					App:       u.intToFlexInt(app.Application),
					Cat:       u.intToFlexInt(app.Category),
					RxBytes:   u.int64ToFlexInt(app.BytesReceived),
					RxPackets: u.int64ToFlexInt(0),
					TxBytes:   u.int64ToFlexInt(app.BytesTransmitted),
					TxPackets: u.int64ToFlexInt(0),
				}
			}

			if d, ok := agg.byCat[app.Category]; ok {
				d.RxBytes.Val += float64(app.BytesReceived)
				d.TxBytes.Val += float64(app.BytesTransmitted)
			} else {
				agg.byCat[app.Category] = &unifi.DPIData{
					App:       u.intToFlexInt(16777215), // unknown app — category aggregate
					Cat:       u.intToFlexInt(app.Category),
					RxBytes:   u.int64ToFlexInt(app.BytesReceived),
					RxPackets: u.int64ToFlexInt(0),
					TxBytes:   u.int64ToFlexInt(app.BytesTransmitted),
					TxPackets: u.int64ToFlexInt(0),
				}
			}
		}
	}

	for siteName, agg := range siteMap {
		byApp := make([]unifi.DPIData, 0, len(agg.byApp))
		for _, d := range agg.byApp {
			byApp = append(byApp, *d)
		}

		byCat := make([]unifi.DPIData, 0, len(agg.byCat))
		for _, d := range agg.byCat {
			byCat = append(byCat, *d)
		}

		metrics.SitesDPI = append(metrics.SitesDPI, &unifi.DPITable{
			ByApp:      byApp,
			ByCat:      byCat,
			SiteName:   siteName,
			SourceName: agg.sourceName,
		})
	}
}

// collectLegacyPerSite collects v5.26.0 additions that use the legacy API (no API key needed).
// Failures are non-fatal: older firmware may not expose these endpoints.
func (u *InputUnifi) collectLegacyPerSite(c *Controller, sites []*unifi.Site, m *Metrics) {
	for _, site := range sites {
		if wan, err := c.Unifi.GetWANStatus(site); err != nil {
			u.LogDebugf("unifi.GetWANStatus(%s, %s): %v (continuing)", c.URL, site.Name, err)
		} else {
			m.WANStatuses = append(m.WANStatuses, wan)
		}

		if forwards, err := c.Unifi.GetPortForwards(site); err != nil {
			u.LogDebugf("unifi.GetPortForwards(%s, %s): %v (continuing)", c.URL, site.Name, err)
		} else {
			m.PortForwards = append(m.PortForwards, forwards...)
		}

		if cert, err := c.Unifi.GetSSLCertificate(site); err != nil {
			u.LogDebugf("unifi.GetSSLCertificate(%s, %s): %v (continuing)", c.URL, site.Name, err)
		} else if cert.ID != "" {
			m.SSLCertificates = append(m.SSLCertificates, cert)
		}

		if upsList, err := c.Unifi.GetUPSDeviceList(site); err != nil {
			u.LogDebugf("unifi.GetUPSDeviceList(%s, %s): %v (continuing)", c.URL, site.Name, err)
		} else {
			m.UPSDevices = append(m.UPSDevices, upsList...)
		}
	}
}

// collectIntegrationV1 collects all Integration/v1 endpoints (Network 9.3.43+, API key required).
// Only called when c.APIKey != "", so ErrAPIKeyRequired will not be returned.
// ErrEndpointNotFound is expected on firmware older than Network 9.3.43.
//
//nolint:cyclop,funlen
func (u *InputUnifi) collectIntegrationV1(c *Controller, sites []*unifi.Site, m *Metrics) {
	// Fetch integration sites — required for all per-site Integration/v1 calls.
	integrationSites, err := c.Unifi.GetIntegrationSites()
	if err != nil {
		if errors.Is(err, unifi.ErrEndpointNotFound) {
			// Integration/v1 requires Network 9.3.43+. Controllers below that return 404.
			u.LogDebugf("unifi.GetIntegrationSites(%s): Integration/v1 not available (Network 9.3.43+ required)", c.URL)
		} else {
			// Unexpected failure (auth expiry, network error, 500) while an API key is configured.
			// All per-site Integration/v1 data will be absent until this resolves.
			u.Logf("unifi.GetIntegrationSites(%s): %v (skipping Integration/v1 per-site collection)", c.URL, err)
		}

		return
	}

	u.LogDebugf("Found %d IntegrationSites", len(integrationSites))

	// Build a map from legacy site name → IntegrationSite to match user-configured sites.
	// IntegrationSite.InternalReference is the same short name used in the legacy API (e.g. "default").
	intSiteByName := make(map[string]*unifi.IntegrationSite, len(integrationSites))
	for _, is := range integrationSites {
		intSiteByName[is.InternalReference] = is
	}

	// Per-site Integration/v1 collections — only for user-configured sites.
	for _, site := range sites {
		is, ok := intSiteByName[site.Name]
		if !ok {
			continue
		}

		if devStats, err := c.Unifi.GetAllIntegrationDeviceStats(is); err != nil {
			u.LogDebugf("unifi.GetAllIntegrationDeviceStats(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.IntegrationDevStats = append(m.IntegrationDevStats, devStats...)
		}

		if broadcasts, err := c.Unifi.GetWifiBroadcasts(is); err != nil {
			u.LogDebugf("unifi.GetWifiBroadcasts(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.WifiBroadcasts = append(m.WifiBroadcasts, broadcasts...)
		}

		if zones, err := c.Unifi.GetFirewallZones(is); err != nil {
			u.LogDebugf("unifi.GetFirewallZones(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.FirewallZones = append(m.FirewallZones, zones...)
		}

		if rules, err := c.Unifi.GetACLRules(is); err != nil {
			u.LogDebugf("unifi.GetACLRules(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.ACLRules = append(m.ACLRules, rules...)
		}

		if servers, err := c.Unifi.GetVPNServers(is); err != nil {
			u.LogDebugf("unifi.GetVPNServers(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.VPNServers = append(m.VPNServers, servers...)
		}

		if tunnels, err := c.Unifi.GetSiteToSiteTunnels(is); err != nil {
			u.LogDebugf("unifi.GetSiteToSiteTunnels(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.SiteToSiteTunnels = append(m.SiteToSiteTunnels, tunnels...)
		}

		if lags, err := c.Unifi.GetLAGs(is); err != nil {
			u.LogDebugf("unifi.GetLAGs(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.LAGs = append(m.LAGs, lags...)
		}

		if mclags, err := c.Unifi.GetMCLAGDomains(is); err != nil {
			u.LogDebugf("unifi.GetMCLAGDomains(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.MCLAGDomains = append(m.MCLAGDomains, mclags...)
		}

		if stacks, err := c.Unifi.GetSwitchStacks(is); err != nil {
			u.LogDebugf("unifi.GetSwitchStacks(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.SwitchStacks = append(m.SwitchStacks, stacks...)
		}

		if policies, err := c.Unifi.GetDNSPolicies(is); err != nil {
			u.LogDebugf("unifi.GetDNSPolicies(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.DNSPolicies = append(m.DNSPolicies, policies...)
		}

		if profiles, err := c.Unifi.GetRADIUSProfiles(is); err != nil {
			u.LogDebugf("unifi.GetRADIUSProfiles(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.RADIUSProfiles = append(m.RADIUSProfiles, profiles...)
		}

		if lists, err := c.Unifi.GetTrafficMatchingLists(is); err != nil {
			u.LogDebugf("unifi.GetTrafficMatchingLists(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.TrafficMatchingLists = append(m.TrafficMatchingLists, lists...)
		}

		if vouchers, err := c.Unifi.GetHotspotVouchers(is); err != nil {
			u.LogDebugf("unifi.GetHotspotVouchers(%s, %s): %v (continuing)", c.URL, is.Name, err)
		} else {
			m.HotspotVouchers = append(m.HotspotVouchers, vouchers...)
		}
	}

	// Global Integration/v1 collections (not per-site).
	if apps, err := c.Unifi.GetDPIApplications(); err != nil {
		u.LogDebugf("unifi.GetDPIApplications(%s): %v (continuing)", c.URL, err)
	} else {
		m.DPIApplications = append(m.DPIApplications, apps...)
		u.LogDebugf("Found %d DPIApplications", len(apps))
	}

	if cats, err := c.Unifi.GetDPICategories(); err != nil {
		u.LogDebugf("unifi.GetDPICategories(%s): %v (continuing)", c.URL, err)
	} else {
		m.DPICategories = append(m.DPICategories, cats...)
		u.LogDebugf("Found %d DPICategories", len(cats))
	}

	if pending, err := c.Unifi.GetPendingDevices(); err != nil {
		u.LogDebugf("unifi.GetPendingDevices(%s): %v (continuing)", c.URL, err)
	} else {
		m.PendingDevices = append(m.PendingDevices, pending...)
		u.LogDebugf("Found %d PendingDevices", len(pending))
	}

	if countries, err := c.Unifi.GetCountries(); err != nil {
		u.LogDebugf("unifi.GetCountries(%s): %v (continuing)", c.URL, err)
	} else {
		m.Countries = append(m.Countries, countries...)
		u.LogDebugf("Found %d Countries", len(countries))
	}
}

// augmentMetrics is our middleware layer between collecting metrics and writing them.
// This is where we can manipuate the returned data or make arbitrary decisions.
// This method currently adds parent device names to client metrics and hashes PII.
// This method also converts our local *Metrics type into a slice of interfaces for poller.
func (u *InputUnifi) augmentMetrics(c *Controller, metrics *Metrics) *poller.Metrics {
	if metrics == nil {
		return nil
	}

	m, devices, bssdIDs := extractDevices(metrics)

	// These come blank, so set them here.
	for _, client := range metrics.Clients {
		if devices[client.Mac] = client.Name; client.Name == "" {
			devices[client.Mac] = client.Hostname
		}

		client.Mac = RedactMacPII(client.Mac, c.HashPII, c.DropPII)
		client.Name = RedactNamePII(client.Name, c.HashPII, c.DropPII)
		client.Hostname = RedactNamePII(client.Hostname, c.HashPII, c.DropPII)
		client.SwName = devices[client.SwMac]
		client.ApName = devices[client.ApMac]
		client.GwName = devices[client.GwMac]
		client.RadioDescription = bssdIDs[client.Bssid] + client.RadioProto

		// Apply site name override for clients if configured
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(client.SiteName) {
			client.SiteName = c.DefaultSiteNameOverride
		}

		m.Clients = append(m.Clients, client)
	}

	for _, client := range metrics.ClientsDPI {
		// Name on Client DPI data also comes blank, find it based on MAC address.
		client.Name = devices[client.MAC]
		if client.Name == "" {
			client.Name = client.MAC
		}

		client.Name = RedactNamePII(client.Name, c.HashPII, c.DropPII)
		client.MAC = RedactMacPII(client.MAC, c.HashPII, c.DropPII)

		// Apply site name override for DPI clients if configured
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(client.SiteName) {
			client.SiteName = c.DefaultSiteNameOverride
		}

		m.ClientsDPI = append(m.ClientsDPI, client)
	}

	for _, ap := range metrics.RogueAPs {
		// XXX: do we need augment this data?
		m.RogueAPs = append(m.RogueAPs, ap)
	}

	if *c.SaveSites {
		for _, site := range metrics.Sites {
			// Apply site name override for sites if configured
			if c.DefaultSiteNameOverride != "" {
				if isDefaultSiteName(site.Name) {
					site.Name = c.DefaultSiteNameOverride
				}

				if isDefaultSiteName(site.SiteName) {
					site.SiteName = c.DefaultSiteNameOverride
				}
			}

			m.Sites = append(m.Sites, site)
		}

		for _, site := range metrics.SitesDPI {
			// Apply site name override for DPI sites if configured
			if c.DefaultSiteNameOverride != "" && isDefaultSiteName(site.SiteName) {
				site.SiteName = c.DefaultSiteNameOverride
			}

			m.SitesDPI = append(m.SitesDPI, site)
		}
	}

	for _, speedTest := range metrics.SpeedTests {
		// Apply site name override for speed tests if configured
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(speedTest.SiteName) {
			speedTest.SiteName = c.DefaultSiteNameOverride
		}

		m.SpeedTests = append(m.SpeedTests, speedTest)
	}

	for _, traffic := range metrics.CountryTraffic {
		// Apply site name override for country traffic if configured
		// UsageByCountry has TrafficSite.SiteName, not SiteName directly
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(traffic.TrafficSite.SiteName) {
			traffic.TrafficSite.SiteName = c.DefaultSiteNameOverride
		}

		m.CountryTraffic = append(m.CountryTraffic, traffic)
	}

	for _, lease := range metrics.DHCPLeases {
		// Apply site name override for DHCP leases if configured
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(lease.SiteName) {
			lease.SiteName = c.DefaultSiteNameOverride
		}

		m.DHCPLeases = append(m.DHCPLeases, lease)
	}

	for _, wanConfig := range metrics.WANConfigs {
		// WANEnrichedConfiguration doesn't have a SiteName field by default
		// The site context is preserved via the collector's site list
		m.WANConfigs = append(m.WANConfigs, wanConfig)
	}

	for _, sysinfo := range metrics.Sysinfos {
		m.Sysinfos = append(m.Sysinfos, sysinfo)
	}

	for _, policy := range metrics.FirewallPolicies {
		// Apply site name override for firewall policies if configured
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(policy.SiteName) {
			policy.SiteName = c.DefaultSiteNameOverride
		}

		m.FirewallPolicies = append(m.FirewallPolicies, policy)
	}

	for _, topo := range metrics.Topologies {
		// Apply site name override for topology if configured
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(topo.SiteName) {
			topo.SiteName = c.DefaultSiteNameOverride
		}

		m.Topologies = append(m.Topologies, topo)
	}

	for _, anomaly := range metrics.PortAnomalies {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(anomaly.SiteName) {
			anomaly.SiteName = c.DefaultSiteNameOverride
		}

		m.PortAnomalies = append(m.PortAnomalies, anomaly)
	}

	for _, mesh := range metrics.VPNMeshes {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(mesh.SiteName) {
			mesh.SiteName = c.DefaultSiteNameOverride
		}

		m.VPNMeshes = append(m.VPNMeshes, mesh)
	}

	// v5.26.0 additions — pass through with site name override applied.
	for _, ws := range metrics.WANStatuses {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(ws.SiteName) {
			ws.SiteName = c.DefaultSiteNameOverride
		}

		m.WANStatuses = append(m.WANStatuses, ws)
	}

	for _, pf := range metrics.PortForwards {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(pf.SiteName) {
			pf.SiteName = c.DefaultSiteNameOverride
		}

		m.PortForwards = append(m.PortForwards, pf)
	}

	for _, cert := range metrics.SSLCertificates {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(cert.SiteName) {
			cert.SiteName = c.DefaultSiteNameOverride
		}

		m.SSLCertificates = append(m.SSLCertificates, cert)
	}

	for _, ups := range metrics.UPSDevices {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(ups.SiteName) {
			ups.SiteName = c.DefaultSiteNameOverride
		}

		m.UPSDevices = append(m.UPSDevices, ups)
	}

	for _, ds := range metrics.IntegrationDevStats {
		m.IntegrationDevStats = append(m.IntegrationDevStats, ds)
	}

	for _, wb := range metrics.WifiBroadcasts {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(wb.SiteName) {
			wb.SiteName = c.DefaultSiteNameOverride
		}

		m.WifiBroadcasts = append(m.WifiBroadcasts, wb)
	}

	for _, fz := range metrics.FirewallZones {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(fz.SiteName) {
			fz.SiteName = c.DefaultSiteNameOverride
		}

		m.FirewallZones = append(m.FirewallZones, fz)
	}

	for _, rule := range metrics.ACLRules {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(rule.SiteName) {
			rule.SiteName = c.DefaultSiteNameOverride
		}

		m.ACLRules = append(m.ACLRules, rule)
	}

	for _, vs := range metrics.VPNServers {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(vs.SiteName) {
			vs.SiteName = c.DefaultSiteNameOverride
		}

		m.VPNServers = append(m.VPNServers, vs)
	}

	for _, t := range metrics.SiteToSiteTunnels {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(t.SiteName) {
			t.SiteName = c.DefaultSiteNameOverride
		}

		m.SiteToSiteTunnels = append(m.SiteToSiteTunnels, t)
	}

	for _, lag := range metrics.LAGs {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(lag.SiteName) {
			lag.SiteName = c.DefaultSiteNameOverride
		}

		m.LAGs = append(m.LAGs, lag)
	}

	for _, mc := range metrics.MCLAGDomains {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(mc.SiteName) {
			mc.SiteName = c.DefaultSiteNameOverride
		}

		m.MCLAGDomains = append(m.MCLAGDomains, mc)
	}

	for _, ss := range metrics.SwitchStacks {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(ss.SiteName) {
			ss.SiteName = c.DefaultSiteNameOverride
		}

		m.SwitchStacks = append(m.SwitchStacks, ss)
	}

	for _, dp := range metrics.DNSPolicies {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(dp.SiteName) {
			dp.SiteName = c.DefaultSiteNameOverride
		}

		m.DNSPolicies = append(m.DNSPolicies, dp)
	}

	for _, rp := range metrics.RADIUSProfiles {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(rp.SiteName) {
			rp.SiteName = c.DefaultSiteNameOverride
		}

		m.RADIUSProfiles = append(m.RADIUSProfiles, rp)
	}

	for _, tl := range metrics.TrafficMatchingLists {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(tl.SiteName) {
			tl.SiteName = c.DefaultSiteNameOverride
		}

		m.TrafficMatchingLists = append(m.TrafficMatchingLists, tl)
	}

	for _, hv := range metrics.HotspotVouchers {
		if c.DefaultSiteNameOverride != "" && isDefaultSiteName(hv.SiteName) {
			hv.SiteName = c.DefaultSiteNameOverride
		}

		m.HotspotVouchers = append(m.HotspotVouchers, hv)
	}

	// Global types — no site name to override.
	for _, app := range metrics.DPIApplications {
		m.DPIApplications = append(m.DPIApplications, app)
	}

	for _, cat := range metrics.DPICategories {
		m.DPICategories = append(m.DPICategories, cat)
	}

	for _, pd := range metrics.PendingDevices {
		m.PendingDevices = append(m.PendingDevices, pd)
	}

	for _, co := range metrics.Countries {
		m.Countries = append(m.Countries, co)
	}

	// Apply default_site_name_override to all metrics if configured.
	// This must be done AFTER all metrics are added to m, so everything is included.
	// This allows us to use the console name for Cloud Gateways while keeping
	// the actual site name ("default") for API calls.
	if c.DefaultSiteNameOverride != "" {
		applySiteNameOverride(m, c.DefaultSiteNameOverride)
	}

	return m
}

// isDefaultSiteName checks if a site name represents a "default" site.
// This handles variations like "default", "Default", "Default (default)", etc.
func isDefaultSiteName(siteName string) bool {
	if siteName == "" {
		return false
	}

	lower := strings.ToLower(siteName)
	// Check for exact match or if it contains "default" as a word
	return lower == "default" || strings.Contains(lower, "default")
}

// applySiteNameOverride replaces "default" site names with the override name
// in all devices, clients, and sites. This allows us to use console names
// for Cloud Gateways in metrics while keeping "default" for API calls.
// This makes metrics more compatible with existing dashboards that expect
// meaningful site names instead of "Default" or "Default (default)".
func applySiteNameOverride(m *poller.Metrics, overrideName string) {
	// Apply to all devices - use type switch for known device types
	for i := range m.Devices {
		switch d := m.Devices[i].(type) {
		case *unifi.UAP:
			if isDefaultSiteName(d.SiteName) {
				d.SiteName = overrideName
			}
		case *unifi.USG:
			if isDefaultSiteName(d.SiteName) {
				d.SiteName = overrideName
			}
		case *unifi.USW:
			if isDefaultSiteName(d.SiteName) {
				d.SiteName = overrideName
			}
		case *unifi.UDM:
			if isDefaultSiteName(d.SiteName) {
				d.SiteName = overrideName
			}
		case *unifi.UXG:
			if isDefaultSiteName(d.SiteName) {
				d.SiteName = overrideName
			}
		case *unifi.UBB:
			if isDefaultSiteName(d.SiteName) {
				d.SiteName = overrideName
			}
		case *unifi.UCI:
			if isDefaultSiteName(d.SiteName) {
				d.SiteName = overrideName
			}
		case *unifi.UDB:
			if isDefaultSiteName(d.SiteName) {
				d.SiteName = overrideName
			}
		case *unifi.PDU:
			if isDefaultSiteName(d.SiteName) {
				d.SiteName = overrideName
			}
		}
	}

	// Apply to all clients
	for i := range m.Clients {
		if client, ok := m.Clients[i].(*unifi.Client); ok {
			if isDefaultSiteName(client.SiteName) {
				client.SiteName = overrideName
			}
		}
	}

	// Apply to sites - check both Name and SiteName fields
	for i := range m.Sites {
		if site, ok := m.Sites[i].(*unifi.Site); ok {
			if isDefaultSiteName(site.Name) {
				site.Name = overrideName
			}

			if isDefaultSiteName(site.SiteName) {
				site.SiteName = overrideName
			}
		}
	}

	// Apply to rogue APs
	for i := range m.RogueAPs {
		if ap, ok := m.RogueAPs[i].(*unifi.RogueAP); ok {
			if isDefaultSiteName(ap.SiteName) {
				ap.SiteName = overrideName
			}
		}
	}

	// Apply to DHCP leases
	for i := range m.DHCPLeases {
		if lease, ok := m.DHCPLeases[i].(*unifi.DHCPLease); ok {
			if isDefaultSiteName(lease.SiteName) {
				lease.SiteName = overrideName
			}
		}
	}

	// Apply to sysinfo (controller metrics)
	for i := range m.Sysinfos {
		if s, ok := m.Sysinfos[i].(*unifi.Sysinfo); ok {
			if isDefaultSiteName(s.SiteName) {
				s.SiteName = overrideName
			}
		}
	}

	// Apply to WAN configs
	for i := range m.WANConfigs {
		if wanConfig, ok := m.WANConfigs[i].(*unifi.WANEnrichedConfiguration); ok {
			// WAN configs don't have SiteName field, but we'll add it in the exporter
			_ = wanConfig
		}
	}

	// Apply to firewall policies
	for i := range m.FirewallPolicies {
		if policy, ok := m.FirewallPolicies[i].(*unifi.FirewallPolicy); ok {
			if isDefaultSiteName(policy.SiteName) {
				policy.SiteName = overrideName
			}
		}
	}

	for i := range m.Topologies {
		if topo, ok := m.Topologies[i].(*unifi.Topology); ok {
			if isDefaultSiteName(topo.SiteName) {
				topo.SiteName = overrideName
			}
		}
	}

	for i := range m.PortAnomalies {
		if anomaly, ok := m.PortAnomalies[i].(*unifi.PortAnomaly); ok {
			if isDefaultSiteName(anomaly.SiteName) {
				anomaly.SiteName = overrideName
			}
		}
	}

	for i := range m.VPNMeshes {
		if mesh, ok := m.VPNMeshes[i].(*unifi.MagicSiteToSiteVPN); ok {
			if isDefaultSiteName(mesh.SiteName) {
				mesh.SiteName = overrideName
			}
		}
	}

	// v5.26.0 additions.
	for i := range m.WANStatuses {
		if ws, ok := m.WANStatuses[i].(*unifi.WANStatus); ok && isDefaultSiteName(ws.SiteName) {
			ws.SiteName = overrideName
		}
	}

	for i := range m.PortForwards {
		if pf, ok := m.PortForwards[i].(*unifi.PortForward); ok && isDefaultSiteName(pf.SiteName) {
			pf.SiteName = overrideName
		}
	}

	for i := range m.SSLCertificates {
		if cert, ok := m.SSLCertificates[i].(*unifi.SSLCertificate); ok && isDefaultSiteName(cert.SiteName) {
			cert.SiteName = overrideName
		}
	}

	for i := range m.UPSDevices {
		if ups, ok := m.UPSDevices[i].(*unifi.UPSDeviceSelector); ok && isDefaultSiteName(ups.SiteName) {
			ups.SiteName = overrideName
		}
	}

	for i := range m.WifiBroadcasts {
		if wb, ok := m.WifiBroadcasts[i].(*unifi.WifiBroadcast); ok && isDefaultSiteName(wb.SiteName) {
			wb.SiteName = overrideName
		}
	}

	for i := range m.FirewallZones {
		if fz, ok := m.FirewallZones[i].(*unifi.FirewallZone); ok && isDefaultSiteName(fz.SiteName) {
			fz.SiteName = overrideName
		}
	}

	for i := range m.ACLRules {
		if r, ok := m.ACLRules[i].(*unifi.ACLRule); ok && isDefaultSiteName(r.SiteName) {
			r.SiteName = overrideName
		}
	}

	for i := range m.VPNServers {
		if vs, ok := m.VPNServers[i].(*unifi.VPNServer); ok && isDefaultSiteName(vs.SiteName) {
			vs.SiteName = overrideName
		}
	}

	for i := range m.SiteToSiteTunnels {
		if t, ok := m.SiteToSiteTunnels[i].(*unifi.SiteToSiteTunnel); ok && isDefaultSiteName(t.SiteName) {
			t.SiteName = overrideName
		}
	}

	for i := range m.LAGs {
		if lag, ok := m.LAGs[i].(*unifi.LAG); ok && isDefaultSiteName(lag.SiteName) {
			lag.SiteName = overrideName
		}
	}

	for i := range m.MCLAGDomains {
		if mc, ok := m.MCLAGDomains[i].(*unifi.MCLAGDomain); ok && isDefaultSiteName(mc.SiteName) {
			mc.SiteName = overrideName
		}
	}

	for i := range m.SwitchStacks {
		if ss, ok := m.SwitchStacks[i].(*unifi.SwitchStack); ok && isDefaultSiteName(ss.SiteName) {
			ss.SiteName = overrideName
		}
	}

	for i := range m.DNSPolicies {
		if dp, ok := m.DNSPolicies[i].(*unifi.DNSPolicy); ok && isDefaultSiteName(dp.SiteName) {
			dp.SiteName = overrideName
		}
	}

	for i := range m.RADIUSProfiles {
		if rp, ok := m.RADIUSProfiles[i].(*unifi.RADIUSProfile); ok && isDefaultSiteName(rp.SiteName) {
			rp.SiteName = overrideName
		}
	}

	for i := range m.TrafficMatchingLists {
		if tl, ok := m.TrafficMatchingLists[i].(*unifi.TrafficMatchingList); ok && isDefaultSiteName(tl.SiteName) {
			tl.SiteName = overrideName
		}
	}

	for i := range m.HotspotVouchers {
		if hv, ok := m.HotspotVouchers[i].(*unifi.HotspotVoucher); ok && isDefaultSiteName(hv.SiteName) {
			hv.SiteName = overrideName
		}
	}
}

// this is a helper function for augmentMetrics.
func extractDevices(metrics *Metrics) (*poller.Metrics, map[string]string, map[string]string) {
	m := &poller.Metrics{TS: metrics.TS}
	devices := make(map[string]string)
	bssdIDs := make(map[string]string)

	for _, r := range metrics.Devices.UAPs {
		devices[r.Mac] = r.Name
		m.Devices = append(m.Devices, r)

		for _, v := range r.VapTable {
			bssdIDs[v.Bssid] = fmt.Sprintf("%s %s %s:", r.Name, v.Radio, v.RadioName)
		}
	}

	for _, r := range metrics.Devices.USGs {
		devices[r.Mac] = r.Name
		m.Devices = append(m.Devices, r)
	}

	for _, r := range metrics.Devices.USWs {
		devices[r.Mac] = r.Name
		m.Devices = append(m.Devices, r)
	}

	for _, r := range metrics.Devices.UDMs {
		devices[r.Mac] = r.Name
		m.Devices = append(m.Devices, r)
	}

	for _, r := range metrics.Devices.UXGs {
		devices[r.Mac] = r.Name
		m.Devices = append(m.Devices, r)
	}

	for _, r := range metrics.Devices.UBBs {
		devices[r.Mac] = r.Name
		m.Devices = append(m.Devices, r)
	}

	for _, r := range metrics.Devices.UCIs {
		devices[r.Mac] = r.Name
		m.Devices = append(m.Devices, r)
	}

	for _, r := range metrics.Devices.UDBs {
		devices[r.Mac] = r.Name
		m.Devices = append(m.Devices, r)

		for _, v := range r.VapTable {
			bssdIDs[v.Bssid] = fmt.Sprintf("%s %s %s:", r.Name, v.Radio, v.RadioName)
		}
	}

	for _, r := range metrics.Devices.PDUs {
		devices[r.Mac] = r.Name
		m.Devices = append(m.Devices, r)
	}

	return m, devices, bssdIDs
}

// RedactNamePII converts a name string to an md5 hash (first 24 chars only).
// Useful for maskiing out personally identifying information.
func RedactNamePII(pii string, hash *bool, dropPII *bool) string {
	if dropPII != nil && *dropPII {
		return ""
	}

	if hash == nil || !*hash || pii == "" {
		return pii
	}

	s := fmt.Sprintf("%x", md5.Sum([]byte(pii))) // nolint: gosec
	// instead of 32 characters, only use 24.
	return s[:24]
}

// RedactMacPII converts a MAC address to an md5 hashed version (first 14 chars only).
// Useful for maskiing out personally identifying information.
func RedactMacPII(pii string, hash *bool, dropPII *bool) (output string) {
	if dropPII != nil && *dropPII {
		return ""
	}

	if hash == nil || !*hash || pii == "" {
		return pii
	}

	s := fmt.Sprintf("%x", md5.Sum([]byte(pii))) // nolint: gosec
	// This formats a "fake" mac address looking string.
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", s[:2], s[2:4], s[4:6], s[6:8], s[8:10], s[10:12], s[12:14])
}

// RedactIPPII converts an IP address to an md5 hashed version (first 12 chars only).
// Useful for maskiing out personally identifying information.
func RedactIPPII(pii string, hash *bool, dropPII *bool) string {
	if dropPII != nil && *dropPII {
		return ""
	}

	if hash == nil || !*hash || pii == "" {
		return pii
	}

	s := fmt.Sprintf("%x", md5.Sum([]byte(pii))) // nolint: gosec
	// Format as a "fake" IP-like string.
	return fmt.Sprintf("%s.%s.%s", s[:4], s[4:8], s[8:12])
}

// getFilteredSites returns a list of sites to fetch data for.
// Omits requested but unconfigured sites. Grabs the full list from the
// controller and returns the sites provided in the config file.
func (u *InputUnifi) getFilteredSites(c *Controller) ([]*unifi.Site, error) {
	u.RLock()
	defer u.RUnlock()

	sites, err := c.Unifi.GetSites()
	if err != nil {
		return nil, fmt.Errorf("controller: %w", err)
	}

	// Note: We do NOT override the site name here because it's used in API calls.
	// The API expects the actual site name (e.g., "default"), not the override.
	// The override will be applied later when augmenting metrics for display purposes.

	if len(c.Sites) == 0 || StringInSlice("all", c.Sites) {
		return sites, nil
	}

	i := 0

	for _, s := range sites {
		// Only include valid sites in the request filter.
		if StringInSlice(s.Name, c.Sites) {
			sites[i] = s
			i++
		}
	}

	return sites[:i], nil
}
