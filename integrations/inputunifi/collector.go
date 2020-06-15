package inputunifi

import (
	"crypto/md5" // nolint: gosec
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/unifi-poller/poller"
	"github.com/unifi-poller/unifi"
)

var (
	errScrapeFilterMatchFailed = fmt.Errorf("scrape filter match failed, and filter is not http URL")
)

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
	u.dynamic[url].Role = url
	u.dynamic[url].URL = url

	return true, u.dynamic[url]
}

func (u *InputUnifi) dynamicController(url string) (*poller.Metrics, error) {
	if !strings.HasPrefix(url, "http") {
		return nil, errScrapeFilterMatchFailed
	}

	new, c := u.newDynamicCntrlr(url)

	if new {
		u.Logf("Authenticating to Dynamic UniFi Controller: %s", url)

		if err := u.getUnifi(c); err != nil {
			return nil, errors.Wrapf(err, "authenticating to %s", url)
		}
	}

	return u.collectController(c)
}

func (u *InputUnifi) collectController(c *Controller) (*poller.Metrics, error) {
	if u.isNill(c) {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return nil, errors.Wrapf(err, "re-authenticating to %s", c.Role)
		}
	}

	metrics, err := u.pollController(c)
	if err != nil {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return metrics, errors.Wrapf(err, "re-authenticating to %s", c.Role)
		}
	}

	return metrics, err
}

func (u *InputUnifi) pollController(c *Controller) (*poller.Metrics, error) {
	var err error

	u.RLock()
	defer u.RUnlock()

	m := &poller.Metrics{TS: time.Now()} // At this point, it's the Current Check.

	// Get the sites we care about.
	if m.Sites, err = u.getFilteredSites(c); err != nil {
		return nil, errors.Wrap(err, "unifi.GetSites()")
	}

	if c.SaveDPI != nil && *c.SaveDPI {
		if m.SitesDPI, err = c.Unifi.GetSiteDPI(m.Sites); err != nil {
			return nil, errors.Wrapf(err, "unifi.GetSiteDPI(%s)", c.URL)
		}

		if m.ClientsDPI, err = c.Unifi.GetClientsDPI(m.Sites); err != nil {
			return nil, errors.Wrapf(err, "unifi.GetClientsDPI(%s)", c.URL)
		}
	}

	if c.SaveIDS != nil && *c.SaveIDS {
		m.IDSList, err = c.Unifi.GetIDS(m.Sites, time.Now().Add(time.Minute), time.Now())
		if err != nil {
			return nil, errors.Wrapf(err, "unifi.GetIDS(%s)", c.URL)
		}
	}

	// Get all the points.
	if m.Clients, err = c.Unifi.GetClients(m.Sites); err != nil {
		return nil, errors.Wrapf(err, "unifi.GetClients(%s)", c.URL)
	}

	if m.Devices, err = c.Unifi.GetDevices(m.Sites); err != nil {
		return nil, errors.Wrapf(err, "unifi.GetDevices(%s)", c.URL)
	}

	return u.augmentMetrics(c, m), nil
}

// augmentMetrics is our middleware layer between collecting metrics and writing them.
// This is where we can manipuate the returned data or make arbitrary decisions.
// This function currently adds parent device names to client metrics.
func (u *InputUnifi) augmentMetrics(c *Controller, metrics *poller.Metrics) *poller.Metrics {
	if metrics == nil || metrics.Devices == nil || metrics.Clients == nil {
		return metrics
	}

	devices := make(map[string]string)
	bssdIDs := make(map[string]string)

	for _, r := range metrics.UAPs {
		devices[r.Mac] = r.Name

		for _, v := range r.VapTable {
			bssdIDs[v.Bssid] = fmt.Sprintf("%s %s %s:", r.Name, v.Radio, v.RadioName)
		}
	}

	for _, r := range metrics.USGs {
		devices[r.Mac] = r.Name
	}

	for _, r := range metrics.USWs {
		devices[r.Mac] = r.Name
	}

	for _, r := range metrics.UDMs {
		devices[r.Mac] = r.Name
	}

	// These come blank, so set them here.
	for i, client := range metrics.Clients {
		if devices[client.Mac] = client.Name; client.Name == "" {
			devices[client.Mac] = client.Hostname
		}

		metrics.Clients[i].Mac = RedactMacPII(metrics.Clients[i].Mac, c.HashPII)
		metrics.Clients[i].Name = RedactNamePII(metrics.Clients[i].Name, c.HashPII)
		metrics.Clients[i].Hostname = RedactNamePII(metrics.Clients[i].Hostname, c.HashPII)
		metrics.Clients[i].SwName = devices[client.SwMac]
		metrics.Clients[i].ApName = devices[client.ApMac]
		metrics.Clients[i].GwName = devices[client.GwMac]
		metrics.Clients[i].RadioDescription = bssdIDs[metrics.Clients[i].Bssid] + metrics.Clients[i].RadioProto
	}

	for i := range metrics.ClientsDPI {
		// Name on Client DPI data also comes blank, find it based on MAC address.
		metrics.ClientsDPI[i].Name = devices[metrics.ClientsDPI[i].MAC]
		if metrics.ClientsDPI[i].Name == "" {
			metrics.ClientsDPI[i].Name = metrics.ClientsDPI[i].MAC
		}

		metrics.ClientsDPI[i].Name = RedactNamePII(metrics.ClientsDPI[i].Name, c.HashPII)
		metrics.ClientsDPI[i].MAC = RedactMacPII(metrics.ClientsDPI[i].MAC, c.HashPII)
	}

	if !*c.SaveSites {
		metrics.Sites = nil
	}

	return metrics
}

// RedactNamePII converts a name string to an md5 hash (first 24 chars only).
// Useful for maskiing out personally identifying information.
func RedactNamePII(pii string, hash *bool) string {
	if hash == nil || !*hash {
		return pii
	}

	s := fmt.Sprintf("%x", md5.Sum([]byte(pii))) // nolint: gosec
	// instead of 32 characters, only use 24.
	return s[:24]
}

// RedactMacPII converts a MAC address to an md5 hashed version (first 14 chars only).
// Useful for maskiing out personally identifying information.
func RedactMacPII(pii string, hash *bool) (output string) {
	if hash == nil || !*hash {
		return pii
	}

	s := fmt.Sprintf("%x", md5.Sum([]byte(pii))) // nolint: gosec
	// This formats a "fake" mac address looking string.
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", s[:2], s[2:4], s[4:6], s[6:8], s[8:10], s[10:12], s[12:14])
}

// getFilteredSites returns a list of sites to fetch data for.
// Omits requested but unconfigured sites. Grabs the full list from the
// controller and returns the sites provided in the config file.
func (u *InputUnifi) getFilteredSites(c *Controller) (unifi.Sites, error) {
	u.RLock()
	defer u.RUnlock()

	sites, err := c.Unifi.GetSites()
	if err != nil {
		return nil, err
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
