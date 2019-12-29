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
	if u.Disable {
		l.Logf("UniFi input plugin disabled!")
		return nil
	}

	if u.setDefaults(&u.Default); len(u.Controllers) < 1 && !u.Dynamic {
		new := u.Default // copy defaults.
		u.Controllers = []*Controller{&new}
	}

	if len(u.Controllers) < 1 {
		l.Logf("No controllers configured. Polling dynamic controllers only!")
	}

	u.dynamic = make(map[string]*Controller)
	u.Logger = l

	for _, c := range u.Controllers {
		u.setDefaults(c)

		switch err := u.getUnifi(c); err {
		case nil:
			if err := u.checkSites(c); err != nil {
				u.LogErrorf("checking sites on %s: %v", c.Role, err)
			}

			u.Logf("Configured UniFi Controller at %s v%s as user %s. Sites: %v",
				c.URL, c.Unifi.ServerVersion, c.User, c.Sites)
		default:
			u.LogErrorf("Controller Auth or Connection failed, but continuing to retry! %s: %v", c.Role, err)
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
	if u.Disable {
		return nil, false, nil
	}

	errs := []string{}
	metrics := &poller.Metrics{}
	ok := false

	if filter != nil && filter.Path != "" {
		if !u.Dynamic {
			return metrics, false, fmt.Errorf("filter path requested but dynamic lookups disabled")
		}

		// Attempt a dynamic metrics fetch from an unconfigured controller.
		m, err := u.dynamicController(filter.Path)

		return m, err == nil && m != nil, err
	}

	// Check if the request is for an existing, configured controller.
	for _, c := range u.Controllers {
		if filter != nil && !strings.EqualFold(c.Role, filter.Role) {
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

	return metrics, ok, nil
}

// RawMetrics returns API output from the first configured unifi controller.
func (u *InputUnifi) RawMetrics(filter *poller.Filter) ([]byte, error) {
	if l := len(u.Controllers); filter.Unit >= l {
		return nil, fmt.Errorf("control number %d not found, %d controller(s) configured (0 index)", filter.Unit, l)
	}

	c := u.Controllers[filter.Unit]
	if u.isNill(c) {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return nil, fmt.Errorf("re-authenticating to %s: %v", c.Role, err)
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
