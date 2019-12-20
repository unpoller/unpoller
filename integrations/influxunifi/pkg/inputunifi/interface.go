package inputunifi

/* This file contains the three poller.Input interface methods. */

import (
	"fmt"
	"os"
	"strings"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
	"golift.io/unifi"
)

// Initialize gets called one time when starting up.
// Satisfies poller.Input interface.
func (u *InputUnifi) Initialize(l poller.Logger) error {
	if u.Config.Disable {
		l.Logf("UniFi input plugin disabled!")
		return nil
	}

	if u.setDefaults(&u.Config.Default); len(u.Config.Controllers) < 1 && !u.Config.Dynamic {
		new := u.Config.Default // copy defaults.
		u.Config.Controllers = []*Controller{&new}
	}

	if len(u.Config.Controllers) < 1 {
		l.Logf("No controllers configured. Polling dynamic controllers only!")
	}

	u.dynamic = make(map[string]*Controller)
	u.Logger = l

	for _, c := range u.Config.Controllers {
		u.setDefaults(c)

		switch err := u.getUnifi(c); err {
		case nil:
			if err := u.checkSites(c); err != nil {
				u.LogErrorf("checking sites on %s: %v", c.Name, err)
			}

			u.Logf("Configured UniFi Controller at %s v%s as user %s. Sites: %v",
				c.URL, c.Unifi.ServerVersion, c.User, c.Sites)
		default:
			u.LogErrorf("Controller Auth or Connection failed, but continuing to retry! %s: %v", c.Name, err)
		}
	}

	return nil
}

// Metrics grabs all the measurements from a UniFi controller and returns them.
func (u *InputUnifi) Metrics() (*poller.Metrics, bool, error) {
	return u.MetricsFrom(nil)
}

// MetricsFrom grabs all the measurements from a UniFi controller and returns them.
func (u *InputUnifi) MetricsFrom(filter *poller.Filter) (*poller.Metrics, bool, error) {
	if u.Config.Disable {
		return nil, false, nil
	}

	errs := []string{}
	metrics := &poller.Metrics{}
	ok := false

	// Check if the request is for an existing, configured controller.
	for _, c := range u.Config.Controllers {
		if filter != nil && !strings.EqualFold(c.Name, filter.Term) {
			continue
		}

		m, err := u.collectController(c)
		if err != nil {
			errs = append(errs, err.Error())
		}

		if m == nil {
			continue
		}

		ok = true
		metrics = poller.AppendMetrics(metrics, m)
	}

	if len(errs) > 0 {
		return metrics, ok, fmt.Errorf(strings.Join(errs, ", "))
	}

	if ok {
		return metrics, true, nil
	}

	if filter != nil && !u.Config.Dynamic {
		return metrics, false, fmt.Errorf("scrape filter match failed and dynamic lookups disabled")
	}

	// Attempt a dynamic metrics fetch from an unconfigured controller.
	m, err := u.dynamicController(filter.Term)

	return m, err == nil && m != nil, err
}

// RawMetrics returns API output from the first configured unifi controller.
func (u *InputUnifi) RawMetrics(filter *poller.Filter) ([]byte, error) {
	if l := len(u.Config.Controllers); filter.Unit >= l {
		return nil, fmt.Errorf("control number %d not found, %d controller(s) configured (0 index)", filter.Unit, l)
	}

	c := u.Config.Controllers[filter.Unit]
	if u.isNill(c) {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return nil, fmt.Errorf("re-authenticating to %s: %v", c.Name, err)
		}
	}

	if err := u.checkSites(c); err != nil {
		return nil, err
	}

	sites, err := u.getFilteredSites(c)
	if err != nil {
		return nil, err
	}

	switch filter.Kind {
	case "d", "device", "devices":
		return u.dumpSitesJSON(c, unifi.APIDevicePath, "Devices", sites)
	case "client", "clients", "c":
		return u.dumpSitesJSON(c, unifi.APIClientPath, "Clients", sites)
	case "other", "o":
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Dumping Path '%s':\n", filter.Path)
		return c.Unifi.GetJSON(filter.Path)
	default:
		return []byte{}, fmt.Errorf("must provide filter: devices, clients, other")
	}
}
