package inputunifi

import (
	"fmt"
	"time"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
	"golift.io/unifi"
)

func (u *InputUnifi) isNill(c Controller) bool {
	u.Config.RLock()
	defer u.Config.RUnlock()

	return c.Unifi == nil
}

func (u *InputUnifi) collectController(c Controller) (*poller.Metrics, error) {
	if u.isNill(c) {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return nil, fmt.Errorf("re-authenticating to %s: %v", c.Name, err)
		}
	}

	m, err := u.pollController(c)
	if err == nil {
		return m, nil
	}

	return u.pollController(c)
}

func (u *InputUnifi) pollController(c Controller) (*poller.Metrics, error) {
	var err error

	u.Config.RLock()
	defer u.Config.RUnlock()

	m := &poller.Metrics{TS: time.Now()} // At this point, it's the Current Check.

	// Get the sites we care about.
	if m.Sites, err = u.getFilteredSites(c); err != nil {
		return m, fmt.Errorf("unifi.GetSites(%v): %v", c.URL, err)
	}

	if c.SaveIDS {
		m.IDSList, err = c.Unifi.GetIDS(m.Sites, time.Now().Add(2*time.Minute), time.Now())
		if err != nil {
			return m, fmt.Errorf("unifi.GetIDS(%v): %v", c.URL, err)
		}
	}

	// Get all the points.
	if m.Clients, err = c.Unifi.GetClients(m.Sites); err != nil {
		return m, fmt.Errorf("unifi.GetClients(%v): %v", c.URL, err)
	}

	if m.Devices, err = c.Unifi.GetDevices(m.Sites); err != nil {
		return m, fmt.Errorf("unifi.GetDevices(%v): %v", c.URL, err)
	}

	return u.augmentMetrics(c, m), nil
}

// augmentMetrics is our middleware layer between collecting metrics and writing them.
// This is where we can manipuate the returned data or make arbitrary decisions.
// This function currently adds parent device names to client metrics.
func (u *InputUnifi) augmentMetrics(c Controller, metrics *poller.Metrics) *poller.Metrics {
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
		metrics.Clients[i].SwName = devices[c.SwMac]
		metrics.Clients[i].ApName = devices[c.ApMac]
		metrics.Clients[i].GwName = devices[c.GwMac]
		metrics.Clients[i].RadioDescription = bssdIDs[metrics.Clients[i].Bssid] + metrics.Clients[i].RadioProto
	}

	if !c.SaveSites {
		metrics.Sites = nil
	}

	return metrics
}

// getFilteredSites returns a list of sites to fetch data for.
// Omits requested but unconfigured sites. Grabs the full list from the
// controller and returns the sites provided in the config file.
func (u *InputUnifi) getFilteredSites(c Controller) (unifi.Sites, error) {
	u.Config.RLock()
	defer u.Config.RUnlock()

	sites, err := c.Unifi.GetSites()
	if err != nil {
		return nil, err
	} else if len(c.Sites) < 1 || poller.StringInSlice("all", c.Sites) {
		return sites, nil
	}

	var i int

	for _, s := range sites {
		// Only include valid sites in the request filter.
		if poller.StringInSlice(s.Name, c.Sites) {
			sites[i] = s
			i++
		}
	}

	return sites[:i], nil
}
