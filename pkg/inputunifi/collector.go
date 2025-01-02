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

	if c.SaveRogue != nil && *c.SaveRogue {
		if m.RogueAPs, err = c.Unifi.GetRogueAPs(sites); err != nil {
			return nil, fmt.Errorf("unifi.GetRogueAPs(%s): %w", c.URL, err)
		}
	}

	if c.SaveDPI != nil && *c.SaveDPI {
		if m.SitesDPI, err = c.Unifi.GetSiteDPI(sites); err != nil {
			return nil, fmt.Errorf("unifi.GetSiteDPI(%s): %w", c.URL, err)
		}

		if m.ClientsDPI, err = c.Unifi.GetClientsDPI(sites); err != nil {
			return nil, fmt.Errorf("unifi.GetClientsDPI(%s): %w", c.URL, err)
		}
	}

	// Get all the points.
	if m.Clients, err = c.Unifi.GetClients(sites); err != nil {
		return nil, fmt.Errorf("unifi.GetClients(%s): %w", c.URL, err)
	}

	if m.Devices, err = c.Unifi.GetDevices(sites); err != nil {
		return nil, fmt.Errorf("unifi.GetDevices(%s): %w", c.URL, err)
	}

	return u.augmentMetrics(c, m), nil
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

// getFilteredSites returns a list of sites to fetch data for.
// Omits requested but unconfigured sites. Grabs the full list from the
// controller and returns the sites provided in the config file.
func (u *InputUnifi) getFilteredSites(c *Controller) ([]*unifi.Site, error) {
	u.RLock()
	defer u.RUnlock()

	sites, err := c.Unifi.GetSites()
	if err != nil {
		return nil, fmt.Errorf("controller: %w", err)
	} else if len(c.Sites) == 0 || StringInSlice("all", c.Sites) {
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
