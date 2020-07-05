package inputunifi

import (
	"time"

	"github.com/pkg/errors"
	"github.com/unifi-poller/unifi"
	"github.com/unifi-poller/webserver"
)

/* Event collection. Events are also sent to the webserver for display. */

func (u *InputUnifi) collectControllerEvents(c *Controller) ([]interface{}, error) {
	var (
		logs    = []interface{}{}
		newLogs []interface{}
	)

	// Get the sites we care about.
	sites, err := u.getFilteredSites(c)
	if err != nil {
		return nil, errors.Wrap(err, "unifi.GetSites()")
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
				return logs, errors.Wrap(err, "unifi.GetAlarms()")
			}

			for _, e := range events {
				logs = append(logs, e)

				webserver.NewInputEvent(PluginName, s.Name+"_alarms", &webserver.Event{Ts: e.Datetime, Msg: e.Msg,
					Tags: map[string]string{"type": "alarm", "key": e.Key, "site_id": e.SiteID,
						"site_name": e.SiteName, "source": e.SourceName},
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
				return logs, errors.Wrap(err, "unifi.GetAnomalies()")
			}

			for _, e := range events {
				logs = append(logs, e)

				webserver.NewInputEvent(PluginName, s.Name+"_anomalies", &webserver.Event{Ts: e.Datetime, Msg: e.Anomaly,
					Tags: map[string]string{"type": "anomaly", "site_name": e.SiteName, "source": e.SourceName},
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
				return logs, errors.Wrap(err, "unifi.GetEvents()")
			}

			for _, e := range events {
				e := redactEvent(e, c.HashPII)
				logs = append(logs, e)

				webserver.NewInputEvent(PluginName, s.Name+"_events", &webserver.Event{Msg: e.Msg, Ts: e.Datetime,
					Tags: map[string]string{"type": "event", "key": e.Key, "site_id": e.SiteID,
						"site_name": e.SiteName, "source": e.SourceName},
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
				return logs, errors.Wrap(err, "unifi.GetIDS()")
			}

			for _, e := range events {
				logs = append(logs, e)

				webserver.NewInputEvent(PluginName, s.Name+"_ids", &webserver.Event{Ts: e.Datetime, Msg: e.Msg,
					Tags: map[string]string{"type": "ids", "key": e.Key, "site_id": e.SiteID,
						"site_name": e.SiteName, "source": e.SourceName},
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
