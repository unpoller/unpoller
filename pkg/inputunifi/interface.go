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

	if u.setDefaults(&u.Config.Default); len(u.Config.Controllers) < 1 {
		new := u.Config.Default // copy defaults.
		u.Config.Controllers = []*Controller{&new}
	}

	u.Logger = l

	for _, c := range u.Config.Controllers {
		u.setDefaults(c)

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
	return u.MetricsFrom(nil)
}

// MetricsFrom grabs all the measurements from a UniFi controller and returns them.
func (u *InputUnifi) MetricsFrom(filter *poller.Filter) (*poller.Metrics, bool, error) {
	if u.Config.Disable || filter == nil || filter.Term == "" {
		return nil, false, nil
	}

	errs := []string{}
	metrics := &poller.Metrics{}
	ok := false

	// Check if the request is for an existing, configured controller.
	for _, c := range u.Config.Controllers {
		if !strings.EqualFold(c.Name, filter.Term) {
			continue
		}

		exists, err := u.appendController(c, metrics)
		if err != nil {
			errs = append(errs, err.Error())
		}

		if exists {
			ok = true
		}
	}

	if len(errs) > 0 {
		return metrics, ok, fmt.Errorf(strings.Join(errs, ", "))
	}

	if u.Config.Dynamic && !ok && strings.HasPrefix(filter.Term, "http") {
		// Attempt to a dynamic metrics fetch from an unconfigured controller.
		return u.dynamicController(filter.Term)
	}

	return metrics, ok, nil
}

// RawMetrics returns API output from the first configured unifi controller.
func (u *InputUnifi) RawMetrics(filter *poller.Filter) ([]byte, error) {
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
