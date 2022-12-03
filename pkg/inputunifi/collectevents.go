package inputunifi

import (
	"fmt"
	"time"

	"github.com/unpoller/unifi"
	"github.com/unpoller/unpoller/pkg/webserver"
)

/* Event collection. Events are also sent to the webserver for display. */

func (u *InputUnifi) collectControllerEvents(c *Controller) ([]interface{}, error) {
	if u.isNill(c) {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return nil, fmt.Errorf("re-authenticating to %s: %w", c.URL, err)
		}
	}

	var (
		logs    = []interface{}{}
		newLogs []interface{}
	)

	// Get the sites we care about.
	sites, err := u.getFilteredSites(c)
	if err != nil {
		return nil, fmt.Errorf("unifi.GetSites(): %w", err)
	}

	type caller func([]interface{}, []*unifi.Site, *Controller) ([]interface{}, error)

	for _, call := range []caller{u.collectIDS, u.collectAnomalies, u.collectAlarms, u.collectEvents} {
		if newLogs, err = call(logs, sites, c); err != nil {
			return logs, err
		}

		logs = append(logs, newLogs...)
	}

	return logs, nil
}

func (u *InputUnifi) collectAlarms(logs []interface{}, sites []*unifi.Site, c *Controller) ([]interface{}, error) {
	if *c.SaveAlarms {
		for _, s := range sites {
			events, err := c.Unifi.GetAlarmsSite(s)
			if err != nil {
				return logs, fmt.Errorf("unifi.GetAlarms(): %w", err)
			}

			for _, e := range events {
				logs = append(logs, e)

				webserver.NewInputEvent(PluginName, s.ID+"_alarms", &webserver.Event{
					Ts: e.Datetime, Msg: e.Msg, Tags: map[string]string{
						"type": "alarm", "key": e.Key, "site_id": e.SiteID,
						"site_name": e.SiteName, "source": e.SourceName,
					},
				})
			}
		}
	}

	return logs, nil
}

func (u *InputUnifi) collectAnomalies(logs []interface{}, sites []*unifi.Site, c *Controller) ([]interface{}, error) {
	if *c.SaveAnomal {
		for _, s := range sites {
			events, err := c.Unifi.GetAnomaliesSite(s)
			if err != nil {
				return logs, fmt.Errorf("unifi.GetAnomalies(): %w", err)
			}

			for _, e := range events {
				logs = append(logs, e)

				webserver.NewInputEvent(PluginName, s.ID+"_anomalies", &webserver.Event{
					Ts: e.Datetime, Msg: e.Anomaly, Tags: map[string]string{
						"type": "anomaly", "site_name": e.SiteName, "source": e.SourceName,
					},
				})
			}
		}
	}

	return logs, nil
}

func (u *InputUnifi) collectEvents(logs []interface{}, sites []*unifi.Site, c *Controller) ([]interface{}, error) {
	if *c.SaveEvents {
		for _, s := range sites {
			events, err := c.Unifi.GetSiteEvents(s, time.Hour)
			if err != nil {
				return logs, fmt.Errorf("unifi.GetEvents(): %w", err)
			}

			for _, e := range events {
				e := redactEvent(e, c.HashPII)
				logs = append(logs, e)

				webserver.NewInputEvent(PluginName, s.ID+"_events", &webserver.Event{
					Msg: e.Msg, Ts: e.Datetime, Tags: map[string]string{
						"type": "event", "key": e.Key, "site_id": e.SiteID,
						"site_name": e.SiteName, "source": e.SourceName,
					},
				})
			}
		}
	}

	return logs, nil
}

func (u *InputUnifi) collectIDS(logs []interface{}, sites []*unifi.Site, c *Controller) ([]interface{}, error) {
	if *c.SaveIDS {
		for _, s := range sites {
			events, err := c.Unifi.GetIDSSite(s)
			if err != nil {
				return logs, fmt.Errorf("unifi.GetIDS(): %w", err)
			}

			for _, e := range events {
				logs = append(logs, e)

				webserver.NewInputEvent(PluginName, s.ID+"_ids", &webserver.Event{
					Ts: e.Datetime, Msg: e.Msg, Tags: map[string]string{
						"type": "ids", "key": e.Key, "site_id": e.SiteID,
						"site_name": e.SiteName, "source": e.SourceName,
					},
				})
			}
		}
	}

	return logs, nil
}

// redactEvent attempts to mask personally identying information from log messages.
// This currently misses the "msg" value entirely and leaks PII information.
func redactEvent(e *unifi.Event, hash *bool) *unifi.Event {
	if !*hash {
		return e
	}

	// metrics.Events[i].Msg <-- not sure what to do here.
	e.DestIPGeo = unifi.IPGeo{}
	e.SourceIPGeo = unifi.IPGeo{}
	e.Host = RedactNamePII(e.Host, hash)
	e.Hostname = RedactNamePII(e.Hostname, hash)
	e.DstMAC = RedactMacPII(e.DstMAC, hash)
	e.SrcMAC = RedactMacPII(e.SrcMAC, hash)

	return e
}
