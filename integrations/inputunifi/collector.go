package inputunifi

import (
	"fmt"
	"strings"
	"time"

	"github.com/unifi-poller/poller"
	"github.com/unifi-poller/unifi"
)

func (u *InputUnifi) isNill(c *Controller) bool {
	u.RLock()
	defer u.RUnlock()

	return c.Unifi == nil
}

// newDynamicCntrlr creates and saves a controller (with auth cookie) for repeated use.
// This is called when an unconfigured controller is requested.
func (u *InputUnifi) newDynamicCntrlr(url string) (bool, *Controller) {
	u.Lock()
	defer u.Unlock()

	c := u.dynamic[url]
	if c != nil {
		// it already exists.
		return false, c
	}

	ccopy := u.Default // copy defaults into new controller
	c = &ccopy
	u.dynamic[url] = c
	c.Role = url
	c.URL = url

	return true, c
}

func (u *InputUnifi) dynamicController(url string) (*poller.Metrics, error) {
	if !strings.HasPrefix(url, "http") {
		return nil, fmt.Errorf("scrape filter match failed, and filter is not http URL")
	}

	new, c := u.newDynamicCntrlr(url)

	if new {
		u.Logf("Authenticating to Dynamic UniFi Controller: %s", url)

		if err := u.getUnifi(c); err != nil {
			return nil, fmt.Errorf("authenticating to %s: %v", url, err)
		}
	}

	return u.collectController(c)
}

func (u *InputUnifi) collectController(c *Controller) (*poller.Metrics, error) {
	if u.isNill(c) {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return nil, fmt.Errorf("re-authenticating to %s: %v", c.Role, err)
		}
	}

	metrics, err := u.pollController(c)
	if err != nil {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return metrics, fmt.Errorf("re-authenticating to %s: %v", c.Role, err)
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
		return nil, fmt.Errorf("unifi.GetSites(): %v", err)
	}

	if c.SaveDPI {
		if m.SitesDPI, err = c.Unifi.GetSiteDPI(m.Sites); err != nil {
			return nil, fmt.Errorf("unifi.GetSiteDPI(%v): %v", c.URL, err)
		}

		if m.ClientsDPI, err = c.Unifi.GetClientsDPI(m.Sites); err != nil {
			return nil, fmt.Errorf("unifi.GetClientsDPI(%v): %v", c.URL, err)
		}
	}

	if c.SaveIDS {
		m.IDSList, err = c.Unifi.GetIDS(m.Sites, time.Now().Add(time.Minute), time.Now())
		if err != nil {
			return nil, fmt.Errorf("unifi.GetIDS(%v): %v", c.URL, err)
		}
	}

	// Get all the points.
	if m.Clients, err = c.Unifi.GetClients(m.Sites); err != nil {
		return nil, fmt.Errorf("unifi.GetClients(%v): %v", c.URL, err)
	}

	if m.Devices, err = c.Unifi.GetDevices(m.Sites); err != nil {
		return nil, fmt.Errorf("unifi.GetDevices(%v): %v", c.URL, err)
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
	for i, c := range metrics.Clients {
		if devices[c.Mac] = c.Name; c.Name == "" {
			devices[c.Mac] = c.Hostname
		}

		metrics.Clients[i].SwName = devices[c.SwMac]
		metrics.Clients[i].ApName = devices[c.ApMac]
		metrics.Clients[i].GwName = devices[c.GwMac]
		metrics.Clients[i].RadioDescription = bssdIDs[metrics.Clients[i].Bssid] + metrics.Clients[i].RadioProto
	}

	for i := range metrics.ClientsDPI {
		// Name on Client DPI data also comes blank, find it based on MAC address.
		metrics.ClientsDPI[i].Name = devices[metrics.ClientsDPI[i].MAC]
		if metrics.ClientsDPI[i].Name == "" {
			metrics.ClientsDPI[i].Name = metrics.ClientsDPI[i].MAC
		}
	}

	if !*c.SaveSites {
		metrics.Sites = nil
	}

	return metrics
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
