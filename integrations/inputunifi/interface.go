package inputunifi

/* This file contains the three poller.Input interface methods. */

import (
	"fmt"
	"os"
	"strings"
	"time"

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
	u.Logf("   => Save Events: %v", *c.SaveEvents)
	u.Logf("   => Save Sites: %v", *c.SaveSites)
}

// Events allows you to pull only events (and IDS) from the UniFi Controller.
// This does not respect HashPII, but it may in the future!
// Use Filter.Dur to set a search duration into the past; 1 minute default.
// Set Filter.Skip to true to disable IDS collection.
func (u *InputUnifi) Events(filter *poller.Filter) (*poller.Events, error) {
	if u.Disable {
		return nil, nil
	}

	events := &poller.Events{}

	for _, c := range u.Controllers {
		if filter != nil && filter.Path != "" &&
			!strings.EqualFold(c.URL, filter.Path) {
			// continue only if we have a filter path and it doesn't match.
			continue
		}

		if filter == nil || filter.Dur == 0 {
			filter = &poller.Filter{Dur: time.Minute}
		}

		// Get the sites we care about.
		sites, err := u.getFilteredSites(c)
		if err != nil {
			return events, errors.Wrap(err, "unifi.GetSites()")
		}

		e, err := c.Unifi.GetEvents(sites, time.Now().Add(-filter.Dur))
		if err != nil {
			return events, errors.Wrap(err, "unifi.GetEvents()")
		}

		for _, l := range e {
			events.Logs = append(events.Logs, l)
		}

		if !filter.Skip {
			i, err := c.Unifi.GetIDS(sites, time.Now().Add(-filter.Dur))
			if err != nil {
				return events, errors.Wrap(err, "unifi.GetIDS()")
			}

			for _, l := range i {
				events.Logs = append(events.Logs, l)
			}
		}
	}

	return events, nil
}

// Metrics grabs all the measurements from a UniFi controller and returns them.
// Set Filter.Path to a controller URL for a specific controller (or get them all).
// Set Filter.Skip to true to Skip Events and IDS collection (Prometheus does this).
func (u *InputUnifi) Metrics(filter *poller.Filter) (*poller.Metrics, error) {
	if u.Disable {
		return nil, nil
	}

	metrics := &poller.Metrics{}

	if filter == nil {
		filter = &poller.Filter{}
	}

	// Check if the request is for an existing, configured controller (or all controllers)
	for _, c := range u.Controllers {
		if filter.Path != "" && !strings.EqualFold(c.URL, filter.Path) {
			// continue only if we have a filter path and it doesn't match.
			continue
		}

		m, err := u.collectController(c, filter)
		if err != nil {
			return metrics, err
		}

		metrics = poller.AppendMetrics(metrics, m)
	}

	if filter.Path == "" || len(metrics.Clients) != 0 {
		return metrics, nil
	}

	if !u.Dynamic {
		return nil, errDynamicLookupsDisabled
	}

	// Attempt a dynamic metrics fetch from an unconfigured controller.
	return u.dynamicController(filter)
}

// RawMetrics returns API output from the first configured UniFi controller.
// Adjust filter.Unit to pull from a controller other than the first.
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

func (u *InputUnifi) dumpSitesJSON(c *Controller, path, name string, sites []*unifi.Site) ([]byte, error) {
	allJSON := []byte{}

	for _, s := range sites {
		apiPath := fmt.Sprintf(path, s.Name)
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Dumping %s: '%s' JSON for site: %s (%s):\n", name, apiPath, s.Desc, s.Name)

		body, err := c.Unifi.GetJSON(apiPath)
		if err != nil {
			return allJSON, err
		}

		allJSON = append(allJSON, body...)
	}

	return allJSON, nil
}
