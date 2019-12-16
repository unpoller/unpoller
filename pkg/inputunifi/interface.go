package inputunifi

import (
	"fmt"
	"strings"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
	"golift.io/unifi"
)

// Metrics grabs all the measurements from a UniFi controller and returns them.
func (u *InputUnifi) Metrics() (*poller.Metrics, error) {
	if u.Config.Disable {
		return nil, nil
	}

	errs := []string{}
	metrics := &poller.Metrics{}

	for _, c := range u.Config.Controllers {
		m, err := u.collectController(c)
		if err != nil {
			errs = append(errs, err.Error())
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

	if len(errs) > 0 {
		return metrics, fmt.Errorf(strings.Join(errs, ", "))
	}

	return metrics, nil
}

// Initialize gets called one time when starting up.
// Satisfies poller.Input interface.
func (u *InputUnifi) Initialize(l poller.Logger) error {
	if u.Config.Disable {
		l.Logf("unifi input disabled")
		return nil
	}

	if len(u.Config.Controllers) < 1 {
		return fmt.Errorf("no unifi controllers defined for unifi input")
	}

	u.Logger = l

	for i, c := range u.Config.Controllers {
		if c.Name == "" {
			u.Config.Controllers[i].Name = c.URL
		}

		switch err := u.getUnifi(c); err {
		case nil:
			if err := u.checkSites(c); err != nil {
				u.LogErrorf("checking sites on %s: %v", c.Name, err)
			}

			u.Logf("Polling UniFi Controller at %s v%s as user %s. Sites: %v",
				c.URL, c.Unifi.ServerVersion, c.User, c.Sites)
		default:
			u.LogErrorf("Controller Auth or Connection failed, but continuing to retry! %s: %v", c.Name, err)
		}
	}

	return nil
}

// checkSites makes sure the list of provided sites exists on the controller.
// This only runs once during initialization.
func (u *InputUnifi) checkSites(c Controller) error {
	u.Config.RLock()
	defer u.Config.RUnlock()
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

	if poller.StringInSlice("all", c.Sites) {
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
