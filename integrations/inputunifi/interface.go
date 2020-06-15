package inputunifi

/* This file contains the three poller.Input interface methods. */

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/unifi-poller/poller"
	"github.com/unifi-poller/unifi"
)

var (
	errDynamicLookupsDisabled = fmt.Errorf("filter path requested but dynamic lookups disabled")
	errControllerNumNotFound  = fmt.Errorf("controller number not found")
	errNoFilterKindProvided   = fmt.Errorf("must provide filter: devices, clients, other")
)

// Initialize gets called one time when starting up.
// Satisfies poller.Input interface.
func (u *InputUnifi) Initialize(l poller.Logger) error {
	if u.Config == nil {
		u.Config = &Config{Disable: true}
	}

	if u.Disable {
		l.Logf("UniFi input plugin disabled or missing configuration!")
		return nil
	}

	if u.setDefaults(&u.Default); len(u.Controllers) == 0 && !u.Dynamic {
		u.Controllers = []*Controller{&u.Default}
	}

	if len(u.Controllers) == 0 {
		l.Logf("No controllers configured. Polling dynamic controllers only!")
	}

	u.dynamic = make(map[string]*Controller)
	u.Logger = l

	for i, c := range u.Controllers {
		switch err := u.getUnifi(u.setControllerDefaults(c)); err {
		case nil:
			if err := u.checkSites(c); err != nil {
				u.LogErrorf("checking sites on %s: %v", c.URL, err)
			}

			u.Logf("Configured UniFi Controller %d:", i+1)
		default:
			u.LogErrorf("Controller %d Auth or Connection Error, retrying: %v", i+1, err)
		}

		u.logController(c)
	}

	return nil
}

func (u *InputUnifi) logController(c *Controller) {
	u.Logf("   => URL: %s", c.URL)

	if c.Unifi != nil {
		u.Logf("   => Version: %s", c.Unifi.ServerVersion)
	}

	u.Logf("   => Username: %s (has password: %v)", c.User, c.Pass != "")
	u.Logf("   => Hash PII: %v", *c.HashPII)
	u.Logf("   => Verify SSL: %v", *c.VerifySSL)
	u.Logf("   => Save DPI: %v", *c.SaveDPI)
	u.Logf("   => Save IDS: %v", *c.SaveIDS)
	u.Logf("   => Save Sites: %v", *c.SaveSites)
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

	metrics := &poller.Metrics{}

	// Check if the request is for an existing, configured controller (or all controllers)
	for _, c := range u.Controllers {
		if filter != nil && !strings.EqualFold(c.URL, filter.Path) {
			continue
		}

		m, err := u.collectController(c)
		if err != nil {
			return metrics, false, err
		}

		metrics = poller.AppendMetrics(metrics, m)
	}

	if filter == nil || len(metrics.Clients) != 0 {
		return metrics, true, nil
	}

	if !u.Dynamic {
		return nil, false, errDynamicLookupsDisabled
	}

	// Attempt a dynamic metrics fetch from an unconfigured controller.
	m, err := u.dynamicController(filter.Path)

	return m, err == nil && m != nil, err
}

// RawMetrics returns API output from the first configured unifi controller.
func (u *InputUnifi) RawMetrics(filter *poller.Filter) ([]byte, error) {
	if l := len(u.Controllers); filter.Unit >= l {
		return nil, errors.Wrapf(errControllerNumNotFound, "%d controller(s) configured, '%d'", l, filter.Unit)
	}

	c := u.Controllers[filter.Unit]
	if u.isNill(c) {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return nil, errors.Wrapf(err, "re-authenticating to %s", c.URL)
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
		return []byte{}, errNoFilterKindProvided
	}
}
