package inputunifi

/* This file contains the three poller.Input interface methods. */

import (
	"fmt"
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

	if u.Logger = l; u.Disable {
		u.Logf("UniFi input plugin disabled or missing configuration!")
		return nil
	}

	if u.setDefaults(&u.Default); len(u.Controllers) == 0 && !u.Dynamic {
		u.Controllers = []*Controller{&u.Default}
	}

	if len(u.Controllers) == 0 {
		u.Logf("No controllers configured. Polling dynamic controllers only! Defaults:")
		u.logController(&u.Default)
	}

	for i, c := range u.Controllers {
		switch err := u.getUnifi(u.setControllerDefaults(c)); err {
		case nil:
			if err := u.checkSites(c); err != nil {
				u.LogErrorf("checking sites on %s: %v", c.URL, err)
			}

			u.Logf("Configured UniFi Controller %d of %d:", i+1, len(u.Controllers))
		default:
			u.LogErrorf("Controller %d of %d Auth or Connection Error, retrying: %v", i+1, len(u.Controllers), err)
		}

		u.logController(c)
	}

	return nil
}

func (u *InputUnifi) logController(c *Controller) {
	u.Logf("   => URL: %s (verify SSL: %v)", c.URL, *c.VerifySSL)

	if c.Unifi != nil {
		u.Logf("   => Version: %s (%s)", c.Unifi.ServerVersion, c.Unifi.UUID)
	}

	u.Logf("   => Username: %s (has password: %v)", c.User, c.Pass != "")
	u.Logf("   => Hash PII / Poll Sites: %v / %s", *c.HashPII, strings.Join(c.Sites, ", "))
	u.Logf("   => Save Sites / Save DPI: %v / %v (metrics)", *c.SaveSites, *c.SaveDPI)
	u.Logf("   => Save Events / Save IDS: %v / %v (logs)", *c.SaveEvents, *c.SaveIDS)
	u.Logf("   => Save Alarms / Anomalies: %v / %v (logs)", *c.SaveAlarms, *c.SaveAnomal)
}

// Events allows you to pull only events (and IDS) from the UniFi Controller.
// This does not fully respect HashPII, but it may in the future!
// Use Filter.Path to pick a specific controller, otherwise poll them all!
func (u *InputUnifi) Events(filter *poller.Filter) (*poller.Events, error) {
	if u.Disable {
		return nil, nil
	}

	logs := []interface{}{}

	if filter == nil {
		filter = &poller.Filter{}
	}

	for _, c := range u.Controllers {
		if filter.Path != "" && !strings.EqualFold(c.URL, filter.Path) {
			continue
		}

		events, err := u.collectControllerEvents(c)
		if err != nil {
			return nil, err
		}

		logs = append(logs, events...)
	}

	return &poller.Events{Logs: logs}, nil
}

// Metrics grabs all the measurements from a UniFi controller and returns them.
// Set Filter.Path to a controller URL for a specific controller (or get them all).
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

		m, err := u.collectController(c)
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
		return u.getSitesJSON(c, unifi.APIDevicePath, sites)
	case "client", "clients", "c":
		return u.getSitesJSON(c, unifi.APIClientPath, sites)
	case "other", "o":
		return c.Unifi.GetJSON(filter.Path)
	default:
		return []byte{}, errNoFilterKindProvided
	}
}

func (u *InputUnifi) getSitesJSON(c *Controller, path string, sites []*unifi.Site) ([]byte, error) {
	allJSON := []byte{}

	for _, s := range sites {
		apiPath := fmt.Sprintf(path, s.Name)
		u.LogDebugf("Returning Path '%s' for site: %s (%s):\n", apiPath, s.Desc, s.Name)

		body, err := c.Unifi.GetJSON(apiPath)
		if err != nil {
			return allJSON, err
		}

		allJSON = append(allJSON, body...)
	}

	return allJSON, nil
}
