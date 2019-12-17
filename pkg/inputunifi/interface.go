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

// Metrics grabs all the measurements from a UniFi controller and returns them.
func (u *InputUnifi) Metrics() (*poller.Metrics, bool, error) {
	return u.MetricsFrom(poller.Filter{})
}

// MetricsFrom grabs all the measurements from a UniFi controller and returns them.
func (u *InputUnifi) MetricsFrom(filter poller.Filter) (*poller.Metrics, bool, error) {
	if u.Config.Disable {
		return nil, false, nil
	}

	errs := []string{}
	metrics := &poller.Metrics{}
	ok := false

	for _, c := range u.Config.Controllers {
		if filter.Term != "" && c.Name != filter.Term {
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
		return metrics, ok, fmt.Errorf(strings.Join(errs, ", "))
	}

	return metrics, ok, nil
}

// RawMetrics returns API output from the first configured unifi controller.
func (u *InputUnifi) RawMetrics(filter poller.Filter) ([]byte, error) {
	c := u.Config.Controllers[0] // We could pull the controller number from the filter.
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

	switch filter.Type {
	case "d", "device", "devices":
		return u.dumpSitesJSON(c, unifi.APIDevicePath, "Devices", sites)
	case "client", "clients", "c":
		return u.dumpSitesJSON(c, unifi.APIClientPath, "Clients", sites)
	case "other", "o":
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Dumping Path '%s':\n", filter.Term)
		return c.Unifi.GetJSON(filter.Term)
	default:
		return []byte{}, fmt.Errorf("must provide filter: devices, clients, other")
	}
}
