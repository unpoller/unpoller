package poller

import (
	"fmt"
	"strings"
	"time"

	"github.com/davidnewhall/unifi-poller/pkg/metrics"
	"golift.io/unifi"
)

// GetUnifi returns a UniFi controller interface.
func (u *UnifiPoller) GetUnifi(c Controller) error {
	var err error

	u.Lock()
	defer u.Unlock()

	if c.Unifi != nil {
		c.Unifi.CloseIdleConnections()
	}
	// Create an authenticated session to the Unifi Controller.
	c.Unifi, err = unifi.NewUnifi(&unifi.Config{
		User:      c.User,
		Pass:      c.Pass,
		URL:       c.URL,
		VerifySSL: c.VerifySSL,
		ErrorLog:  u.LogErrorf, // Log all errors.
		DebugLog:  u.LogDebugf, // Log debug messages.
	})

	if err != nil {
		c.Unifi = nil
		return fmt.Errorf("unifi controller: %v", err)
	}

	u.LogDebugf("Authenticated with controller successfully, %s", c.URL)

	return u.CheckSites(c)
}

// CheckSites makes sure the list of provided sites exists on the controller.
// This does not run in Lambda (run-once) mode.
func (u *UnifiPoller) CheckSites(c Controller) error {
	if strings.Contains(strings.ToLower(u.Config.Mode), "lambda") {
		return nil // Skip this in lambda mode.
	}

	u.LogDebugf("Checking Controller Sites List")

	sites, err := c.Unifi.GetSites()
	if err != nil {
		return err
	}

	msg := []string{}

	for _, site := range sites {
		msg = append(msg, site.Name+" ("+site.Desc+")")
	}
	u.Logf("Found %d site(s) on controller: %v", len(msg), strings.Join(msg, ", "))

	if StringInSlice("all", c.Sites) {
		c.Sites = []string{"all"}
		return nil
	}

FIRST:
	for _, s := range c.Sites {
		for _, site := range sites {
			if s == site.Name {
				continue FIRST
			}
		}
		return fmt.Errorf("configured site not found on controller: %v", s)
	}

	return nil
}

// CollectMetrics grabs all the measurements from a UniFi controller and returns them.
func (u *UnifiPoller) CollectMetrics() (*metrics.Metrics, error) {
	errs := []string{}
	metrics := &metrics.Metrics{}

	for _, c := range u.Config.Controllers {
		m, err := u.collectController(c)
		if err != nil {
			u.LogErrorf("collecting metrics from %s: %v", c.URL, err)
			u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

			if err := u.GetUnifi(c); err != nil {
				u.LogErrorf("re-authenticating to %s: %v", c.URL, err)
				errs = append(errs, err.Error())
			} else if m, err = u.collectController(c); err != nil {
				u.LogErrorf("collecting metrics from %s: %v", c.URL, err)
				errs = append(errs, err.Error())
			}
		}

		if m == nil {
			continue
		}

		metrics.Sites = append(metrics.Sites, m.Sites...)
		metrics.Clients = append(metrics.Clients, m.Clients...)
		metrics.IDSList = append(metrics.IDSList, m.IDSList...)

		if m.Devices == nil {
			continue
		}

		if metrics.Devices == nil {
			metrics.Devices = &unifi.Devices{}
		}

		metrics.UAPs = append(metrics.UAPs, m.UAPs...)
		metrics.USGs = append(metrics.USGs, m.USGs...)
		metrics.USWs = append(metrics.USWs, m.USWs...)
		metrics.UDMs = append(metrics.UDMs, m.UDMs...)
	}

	var err error

	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, ", "))
	}

	return metrics, err
}

func (u *UnifiPoller) collectController(c Controller) (*metrics.Metrics, error) {
	var err error

	if c.Unifi == nil {
		// Some users need to re-auth every interval because the cookie times out.
		// Sometimes we hit this path when the controller dies.
		u.Logf("Re-authenticating to UniFi Controller: %v", c.URL)

		if err := u.GetUnifi(c); err != nil {
			return nil, err
		}
	}

	m := &metrics.Metrics{TS: u.LastCheck} // At this point, it's the Current Check.

	// Get the sites we care about.
	if m.Sites, err = u.GetFilteredSites(c); err != nil {
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
func (u *UnifiPoller) augmentMetrics(c Controller, metrics *metrics.Metrics) *metrics.Metrics {
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

// GetFilteredSites returns a list of sites to fetch data for.
// Omits requested but unconfigured sites. Grabs the full list from the
// controller and returns the sites provided in the config file.
func (u *UnifiPoller) GetFilteredSites(c Controller) (unifi.Sites, error) {
	var i int

	sites, err := c.Unifi.GetSites()
	if err != nil {
		return nil, err
	} else if len(c.Sites) < 1 || StringInSlice("all", c.Sites) {
		return sites, nil
	}

	for _, s := range sites {
		// Only include valid sites in the request filter.
		if StringInSlice(s.Name, c.Sites) {
			sites[i] = s
			i++
		}
	}

	return sites[:i], nil
}
