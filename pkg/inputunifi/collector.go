package inputunifi

// nolint: gosec
import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/poller"
)

const (
	history_seconds = 86400
	poll_duration   = time.Second * history_seconds
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
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return metrics, fmt.Errorf("re-authenticating to %s: %w", c.URL, err)
		}
	}

	return metrics, err
}

//nolint:cyclop
func (u *InputUnifi) pollController(c *Controller) (*poller.Metrics, error) {
	u.RLock()
	defer u.RUnlock()

	u.LogDebugf("Polling controller: %s (%s)", c.URL, c.ID)

	// Get the sites we care about.
	sites, err := u.getFilteredSites(c)
	if err != nil {
		return nil, fmt.Errorf("unifi.GetSites(): %w", err)
	}

	m := &Metrics{TS: time.Now(), Sites: sites}
	defer updateWeb(c, m)

	// FIXME needs to be last poll time maybe
	st := m.TS.Add(-1 * poll_duration)
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

	if c.SaveTraffic != nil && *c.SaveTraffic && c.SaveDPI != nil && *c.SaveDPI {
		clientUsageByApp, err := c.Unifi.GetClientTraffic(sites, &tp, true)
		if err != nil {
			return nil, fmt.Errorf("unifi.GetClientTraffic(%s): %w", c.URL, err)
		}
		u.LogDebugf("Found %d ClientUsageByApp entries", len(clientUsageByApp))
		b4 := len(m.ClientsDPI)
		u.convertToClientDPI(clientUsageByApp, m)
		u.LogDebugf("Added %d ClientDPI entries for a total of %d", len(m.ClientsDPI)-b4, len(m.ClientsDPI))
	}

	// Get all the points.
	if m.Clients, err = c.Unifi.GetClients(sites); err != nil {
		return nil, fmt.Errorf("unifi.GetClients(%s): %w", c.URL, err)
	}
	u.LogDebugf("Found %d Clients entries", len(m.Clients))

	if m.Devices, err = c.Unifi.GetDevices(sites); err != nil {
		return nil, fmt.Errorf("unifi.GetDevices(%s): %w", c.URL, err)
	}
	u.LogDebugf("Found %d UBB, %d UXG, %d PDU, %d UCI, %d UAP %d USG %d USW %d UDM devices",
		len(m.Devices.UBBs), len(m.Devices.UXGs),
		len(m.Devices.PDUs), len(m.Devices.UCIs),
		len(m.Devices.UAPs), len(m.Devices.USGs),
		len(m.Devices.USWs), len(m.Devices.UDMs))

	// Get speed test results for all WANs
	if m.SpeedTests, err = c.Unifi.GetSpeedTests(sites, history_seconds); err != nil {
		// Don't fail collection if speed tests fail - older controllers may not have this endpoint
		u.LogDebugf("unifi.GetSpeedTests(%s): %v (continuing)", c.URL, err)
	} else {
		u.LogDebugf("Found %d SpeedTests entries", len(m.SpeedTests))
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
		m.ClientsDPI = append(m.ClientsDPI, client)
	}

	for _, ap := range metrics.RogueAPs {
		// XXX: do we need augment this data?
		m.RogueAPs = append(m.RogueAPs, ap)
	}

	if *c.SaveSites {
		for _, site := range metrics.Sites {
			m.Sites = append(m.Sites, site)
		}

		for _, site := range metrics.SitesDPI {
			m.SitesDPI = append(m.SitesDPI, site)
		}
	}

	for _, speedTest := range metrics.SpeedTests {
		m.SpeedTests = append(m.SpeedTests, speedTest)
	}

	for _, traffic := range metrics.CountryTraffic {
		m.CountryTraffic = append(m.CountryTraffic, traffic)
	}

	return m
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

	// Apply the default_site_name_override to the first site in the list, if configured.
	if len(sites) > 0 && c.DefaultSiteNameOverride != "" {
		sites[0].Name = c.DefaultSiteNameOverride
	}

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
