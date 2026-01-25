package inputunifi

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/webserver"
)

/* Event collection. Events are also sent to the webserver for display. */

func (u *InputUnifi) collectControllerEvents(c *Controller) ([]any, error) {
	u.LogDebugf("Collecting controller events: %s (%s)", c.URL, c.ID)

	if u.isNill(c) {
		u.Logf("Re-authenticating to UniFi Controller: %s", c.URL)

		if err := u.getUnifi(c); err != nil {
			return nil, fmt.Errorf("re-authenticating to %s: %w", c.URL, err)
		}
	}

	var (
		logs    = []any{}
		newLogs []any
	)

	// Get the sites we care about.
	sites, err := u.getFilteredSites(c)
	if err != nil {
		return nil, fmt.Errorf("unifi.GetSites(): %w", err)
	}

	type caller func([]any, []*unifi.Site, *Controller) ([]any, error)

	for _, call := range []caller{u.collectIDs, u.collectAnomalies, u.collectAlarms, u.collectEvents, u.collectSyslog, u.collectProtectLogs} {
		if newLogs, err = call(logs, sites, c); err != nil {
			return logs, err
		}

		logs = append(logs, newLogs...)
	}

	return logs, nil
}

func (u *InputUnifi) collectAlarms(logs []any, sites []*unifi.Site, c *Controller) ([]any, error) {
	if *c.SaveAlarms {
		u.LogDebugf("Collecting controller alarms: %s (%s)", c.URL, c.ID)

		// Get devices for all sites to build MAC-to-name lookup
		devices, err := c.Unifi.GetDevices(sites)
		if err != nil {
			u.LogDebugf("Failed to get devices for alarm enrichment: %v (continuing without device names)", err)
			devices = &unifi.Devices{} // Empty devices struct, alarms will not have device names
		}

		// Build MAC address to device name lookup map
		macToName := make(map[string]string)
		for _, d := range devices.UAPs {
			if d.Mac != "" && d.Name != "" {
				macToName[strings.ToLower(d.Mac)] = d.Name
			}
		}
		for _, d := range devices.USGs {
			if d.Mac != "" && d.Name != "" {
				macToName[strings.ToLower(d.Mac)] = d.Name
			}
		}
		for _, d := range devices.USWs {
			if d.Mac != "" && d.Name != "" {
				macToName[strings.ToLower(d.Mac)] = d.Name
			}
		}
		for _, d := range devices.UDMs {
			if d.Mac != "" && d.Name != "" {
				macToName[strings.ToLower(d.Mac)] = d.Name
			}
		}
		for _, d := range devices.UXGs {
			if d.Mac != "" && d.Name != "" {
				macToName[strings.ToLower(d.Mac)] = d.Name
			}
		}
		for _, d := range devices.PDUs {
			if d.Mac != "" && d.Name != "" {
				macToName[strings.ToLower(d.Mac)] = d.Name
			}
		}
		for _, d := range devices.UBBs {
			if d.Mac != "" && d.Name != "" {
				macToName[strings.ToLower(d.Mac)] = d.Name
			}
		}
		for _, d := range devices.UCIs {
			if d.Mac != "" && d.Name != "" {
				macToName[strings.ToLower(d.Mac)] = d.Name
			}
		}

		for _, s := range sites {
			events, err := c.Unifi.GetAlarmsSite(s)
			if err != nil {
				return logs, fmt.Errorf("unifi.GetAlarms(): %w", err)
			}

			for _, e := range events {
				// Try to extract MAC address from alarm message and enrich with device name
				e.DeviceName = u.extractDeviceNameFromAlarm(e, macToName)

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

func (u *InputUnifi) collectAnomalies(logs []any, sites []*unifi.Site, c *Controller) ([]any, error) {
	if *c.SaveAnomal {
		u.LogDebugf("Collecting controller anomalies: %s (%s)", c.URL, c.ID)

		for _, s := range sites {
			events, err := c.Unifi.GetAnomaliesSite(s)
			if err != nil {
				return logs, fmt.Errorf("unifi.GetAnomalies(): %w", err)
			}

			for _, e := range events {
				// Apply site name override for anomalies if configured
				if c.DefaultSiteNameOverride != "" {
					lower := strings.ToLower(e.SiteName)
					if lower == "default" || strings.Contains(lower, "default") {
						e.SiteName = c.DefaultSiteNameOverride
					}
				}
				
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

func (u *InputUnifi) collectEvents(logs []any, sites []*unifi.Site, c *Controller) ([]any, error) {
	if *c.SaveEvents {
		u.LogDebugf("Collecting controller site events (v1): %s (%s)", c.URL, c.ID)

		for _, s := range sites {
			events, err := c.Unifi.GetSiteEvents(s, time.Hour)
			if err != nil {
				return logs, fmt.Errorf("unifi.GetEvents(): %w", err)
			}

			for _, e := range events {
				e := redactEvent(e, c.HashPII, c.DropPII)
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

func (u *InputUnifi) collectSyslog(logs []any, sites []*unifi.Site, c *Controller) ([]any, error) {
	if *c.SaveSyslog {
		u.LogDebugf("Collecting controller syslog (v2): %s (%s)", c.URL, c.ID)

		// Use v2 system-log API
		req := unifi.DefaultSystemLogRequest(time.Hour)
		entries, err := c.Unifi.GetSystemLog(sites, req)
		if err != nil {
			return logs, fmt.Errorf("unifi.GetSystemLog(): %w", err)
		}

		for _, e := range entries {
			e := redactSystemLogEntry(e, c.HashPII, c.DropPII)
			logs = append(logs, e)

			webserver.NewInputEvent(PluginName, e.SiteName+"_syslog", &webserver.Event{
				Msg: e.Msg(), Ts: e.Datetime(), Tags: map[string]string{
					"type": "syslog", "key": e.Key, "event": e.Event,
					"site_name": e.SiteName, "source": e.SourceName,
					"category": e.Category, "subcategory": e.Subcategory,
					"severity": e.Severity,
				},
			})
		}
	}

	return logs, nil
}

func (u *InputUnifi) collectProtectLogs(logs []any, _ []*unifi.Site, c *Controller) ([]any, error) {
	if *c.SaveProtectLogs {
		u.LogDebugf("Collecting Protect logs: %s (%s)", c.URL, c.ID)

		req := unifi.DefaultProtectLogRequest(0) // Uses default 24-hour window
		entries, err := c.Unifi.GetProtectLogs(req)
		if err != nil {
			return logs, fmt.Errorf("unifi.GetProtectLogs(): %w", err)
		}

		for _, e := range entries {
			e := redactProtectLogEntry(e, c.HashPII, c.DropPII)

			// Fetch thumbnail if enabled and event has a camera (only camera events have real thumbnails)
			// Skip access/adminActivity events - they don't have actual camera thumbnails
			if *c.ProtectThumbnails && e.Thumbnail != "" && e.Camera != "" && hasProtectThumbnail(e.Type) {
				// Thumbnail field is like "e-69499de2037add03e4015fa8" - strip "e-" prefix
				thumbID := e.Thumbnail
				if len(thumbID) > 2 && thumbID[:2] == "e-" {
					thumbID = thumbID[2:]
				}
				if thumbData, err := c.Unifi.GetProtectEventThumbnail(thumbID); err == nil {
					e.ThumbnailBase64 = base64.StdEncoding.EncodeToString(thumbData)
				} else {
					u.LogDebugf("Failed to fetch thumbnail for event %s (thumb: %s): %v", e.ID, thumbID, err)
				}
			}

			logs = append(logs, e)

			webserver.NewInputEvent(PluginName, "protect_logs", &webserver.Event{
				Msg: e.Msg(), Ts: e.Datetime(), Tags: map[string]string{
					"type":        "protect_log",
					"event_type":  e.GetEventType(),
					"category":    e.GetCategory(),
					"subcategory": e.GetSubCategory(),
					"severity":    e.GetSeverity(),
					"camera":      e.Camera,
					"source":      e.SourceName,
				},
			})
		}
	}

	return logs, nil
}

func (u *InputUnifi) collectIDs(logs []any, sites []*unifi.Site, c *Controller) ([]any, error) {
	if *c.SaveIDs {
		u.LogDebugf("Collecting controller IDs data: %s (%s)", c.URL, c.ID)

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
func redactEvent(e *unifi.Event, hash *bool, dropPII *bool) *unifi.Event {
	if !*hash && !*dropPII {
		return e
	}

	// metrics.Events[i].Msg <-- not sure what to do here.
	e.DestIPGeo = unifi.IPGeo{}
	e.SourceIPGeo = unifi.IPGeo{}

	if *dropPII {
		e.Host = ""
		e.Hostname = ""
		e.DstMAC = ""
		e.SrcMAC = ""
	} else {
		// hash it
		e.Host = RedactNamePII(e.Host, hash, dropPII)
		e.Hostname = RedactNamePII(e.Hostname, hash, dropPII)
		e.DstMAC = RedactMacPII(e.DstMAC, hash, dropPII)
		e.SrcMAC = RedactMacPII(e.SrcMAC, hash, dropPII)
	}

	return e
}

// redactSystemLogEntry attempts to mask personally identifying information from v2 system log entries.
func redactSystemLogEntry(e *unifi.SystemLogEntry, hash *bool, dropPII *bool) *unifi.SystemLogEntry {
	if !*hash && !*dropPII {
		return e
	}

	// Redact CLIENT parameter if present
	if client, ok := e.Parameters["CLIENT"]; ok {
		if *dropPII {
			client.Hostname = ""
			client.Name = ""
			client.ID = ""
			client.IP = ""
		} else {
			client.Hostname = RedactNamePII(client.Hostname, hash, dropPII)
			client.Name = RedactNamePII(client.Name, hash, dropPII)
			client.ID = RedactMacPII(client.ID, hash, dropPII)
			client.IP = RedactIPPII(client.IP, hash, dropPII)
		}
		e.Parameters["CLIENT"] = client
	}

	// Redact IP parameter if present
	if ip, ok := e.Parameters["IP"]; ok {
		if *dropPII {
			ip.ID = ""
			ip.Name = ""
		} else {
			ip.ID = RedactIPPII(ip.ID, hash, dropPII)
			ip.Name = RedactIPPII(ip.Name, hash, dropPII)
		}
		e.Parameters["IP"] = ip
	}

	// Redact ADMIN parameter if present
	if admin, ok := e.Parameters["ADMIN"]; ok {
		if *dropPII {
			admin.Name = ""
		} else {
			admin.Name = RedactNamePII(admin.Name, hash, dropPII)
		}
		e.Parameters["ADMIN"] = admin
	}

	return e
}

// redactProtectLogEntry attempts to mask personally identifying information from Protect log entries.
func redactProtectLogEntry(e *unifi.ProtectLogEntry, hash *bool, dropPII *bool) *unifi.ProtectLogEntry {
	if !*hash && !*dropPII {
		return e
	}

	// Redact user names from message keys
	if e.Description != nil {
		for i, mk := range e.Description.MessageKeys {
			if mk.Key == "userLink" || mk.Action == "viewUsers" {
				if *dropPII {
					e.Description.MessageKeys[i].Text = ""
				} else {
					e.Description.MessageKeys[i].Text = RedactNamePII(mk.Text, hash, dropPII)
				}
			}
		}
	}

	return e
}

// hasProtectThumbnail returns true if the event type has actual camera thumbnails.
// Access and adminActivity events don't have real thumbnails (they're user activity logs).
func hasProtectThumbnail(eventType string) bool {
	switch eventType {
	case "motion", "smartDetectZone", "smartDetectLine", "ring", "sensorMotion",
		"sensorContact", "sensorAlarm", "doorbell", "package", "person", "vehicle",
		"animal", "face", "licensePlate":
		return true
	default:
		return false
	}
}

// extractDeviceNameFromAlarm attempts to extract a device name for an alarm by looking up
// MAC addresses found in the alarm message or fields. Returns empty string if no match found.
func (u *InputUnifi) extractDeviceNameFromAlarm(alarm *unifi.Alarm, macToName map[string]string) string {
	// Try to extract MAC from message like "AP[fc:ec:da:89:a6:91] was disconnected"
	// Look for pattern: [XX:XX:XX:XX:XX:XX] where X is hex digit
	msg := alarm.Msg

	// Simple regex-like search for MAC address in brackets
	start := strings.Index(msg, "[")
	end := strings.Index(msg, "]")
	if start >= 0 && end > start {
		potentialMAC := msg[start+1 : end]
		// Basic validation: should be 17 characters and contain colons
		if len(potentialMAC) == 17 && strings.Count(potentialMAC, ":") == 5 {
			if name, ok := macToName[strings.ToLower(potentialMAC)]; ok {
				return name
			}
		}
	}

	// Also try SrcMAC and DstMAC fields if present
	if alarm.SrcMAC != "" {
		if name, ok := macToName[strings.ToLower(alarm.SrcMAC)]; ok {
			return name
		}
	}

	if alarm.DstMAC != "" {
		if name, ok := macToName[strings.ToLower(alarm.DstMAC)]; ok {
			return name
		}
	}

	return ""
}
