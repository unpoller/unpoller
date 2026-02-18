package inputunifi

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/unpoller/unifi/v5"
	"github.com/unpoller/unpoller/pkg/webserver"
)

/* This code reformats our data to be displayed on the built-in web interface. */

// controllerID returns the controller UUID for display, or "" if the client is nil (e.g. after 429 re-auth failure).
// Avoids SIGSEGV when updateWeb runs while c.Unifi is nil (see unpoller/unpoller#943).
func controllerID(c *Controller) string {
	if c == nil {
		return ""
	}
	client := c.Unifi
	if client == nil {
		return ""
	}
	return client.UUID
}

func updateWeb(c *Controller, metrics *Metrics) {
	defer func() {
		if r := recover(); r != nil {
			// Log and swallow panic so one bad controller doesn't kill the process (e.g. nil c/c.Unifi race).
			log.Printf("[ERROR] updateWeb panic recovered: %v", r)
		}
	}()
	if c == nil || metrics == nil {
		return
	}
	if c.Unifi == nil {
		return
	}
	webserver.UpdateInput(&webserver.Input{
		Name:    PluginName, // Forgetting this leads to 3 hours of head scratching.
		Sites:   formatSites(c, metrics.Sites),
		Clients: formatClients(c, metrics.Clients),
		Devices: formatDevices(c, metrics.Devices),
	})
}

func formatConfig(config *Config) *Config {
	return &Config{
		Default:     *formatControllers([]*Controller{&config.Default})[0],
		Disable:     config.Disable,
		Dynamic:     config.Dynamic,
		Controllers: formatControllers(config.Controllers),
	}
}

func formatControllers(controllers []*Controller) []*Controller {
	fixed := []*Controller{}

	for _, c := range controllers {
		id := ""
		if c.Unifi != nil {
			id = c.Unifi.UUID
		}

		fixed = append(fixed, &Controller{
			VerifySSL:       c.VerifySSL,
			SaveAnomal:      c.SaveAnomal,
			SaveAlarms:      c.SaveAlarms,
			SaveRogue:       c.SaveRogue,
			SaveEvents:      c.SaveEvents,
			SaveSyslog:      c.SaveSyslog,
			SaveProtectLogs: c.SaveProtectLogs,
			SaveIDs:         c.SaveIDs,
			SaveDPI:         c.SaveDPI,
			HashPII:         c.HashPII,
			DropPII:         c.DropPII,
			SaveSites:       c.SaveSites,
			SaveTraffic:     c.SaveTraffic,
			User:            c.User,
			Pass:            strconv.FormatBool(c.Pass != ""),
			APIKey:          strconv.FormatBool(c.APIKey != ""),
			URL:             c.URL,
			Sites:           c.Sites,
			ID:              id,
		})
	}

	return fixed
}

func formatSites(c *Controller, sites []*unifi.Site) (s webserver.Sites) {
	id := controllerID(c)

	for _, site := range sites {
		if site == nil {
			continue
		}
		s = append(s, &webserver.Site{
			ID:         site.ID,
			Name:       site.Name,
			Desc:       site.Desc,
			Source:     site.SourceName,
			Controller: id,
		})
	}

	return s
}

func formatClients(c *Controller, clients []*unifi.Client) (d webserver.Clients) {
	for _, client := range clients {
		if client == nil {
			continue
		}
		clientType, deviceMAC := "unknown", "unknown"
		if client.ApMac != "" {
			clientType = "wireless"
			deviceMAC = client.ApMac
		} else if client.SwMac != "" {
			clientType = "wired"
			deviceMAC = client.SwMac
		}

		if deviceMAC == "" {
			deviceMAC = client.GwMac
		}

		d = append(d, &webserver.Client{
			Name:       client.Name,
			SiteID:     client.SiteID,
			Source:     client.SourceName,
			Controller: controllerID(c),
			MAC:        client.Mac,
			IP:         client.IP,
			Type:       clientType,
			DeviceMAC:  deviceMAC,
			Rx:         client.RxBytes.Int64(),
			Tx:         client.TxBytes.Int64(),
			Since:      time.Unix(client.FirstSeen.Int64(), 0),
			Last:       time.Unix(client.LastSeen.Int64(), 0),
		})
	}

	return d
}

func formatDevices(c *Controller, devices *unifi.Devices) (d webserver.Devices) { // nolint: funlen
	if devices == nil {
		return d
	}

	id := controllerID(c)

	for _, device := range devices.UAPs {
		if device == nil {
			continue
		}
		d = append(d, &webserver.Device{
			Name:       device.Name,
			SiteID:     device.SiteID,
			Source:     device.SourceName,
			Controller: id,
			MAC:        device.Mac,
			IP:         device.IP,
			Type:       device.Type,
			Model:      device.Model,
			Version:    device.Version,
			Uptime:     int(device.Uptime.Val),
			Clients:    int(device.NumSta.Val),
			Config:     nil,
		})
	}

	for _, device := range devices.UDMs {
		if device == nil {
			continue
		}
		d = append(d, &webserver.Device{
			Name:       device.Name,
			SiteID:     device.SiteID,
			Source:     device.SourceName,
			Controller: id,
			MAC:        device.Mac,
			IP:         device.IP,
			Type:       device.Type,
			Model:      device.Model,
			Version:    device.Version,
			Uptime:     int(device.Uptime.Val),
			Clients:    int(device.NumSta.Val),
			Config:     nil,
		})
	}

	for _, device := range devices.USWs {
		if device == nil {
			continue
		}
		d = append(d, &webserver.Device{
			Name:       device.Name,
			SiteID:     device.SiteID,
			Source:     device.SourceName,
			Controller: id,
			MAC:        device.Mac,
			IP:         device.IP,
			Type:       device.Type,
			Model:      device.Model,
			Version:    device.Version,
			Uptime:     int(device.Uptime.Val),
			Clients:    int(device.NumSta.Val),
			Config:     nil,
		})
	}

	for _, device := range devices.USGs {
		if device == nil {
			continue
		}
		d = append(d, &webserver.Device{
			Name:       device.Name,
			SiteID:     device.SiteID,
			Source:     device.SourceName,
			Controller: id,
			MAC:        device.Mac,
			IP:         device.IP,
			Type:       device.Type,
			Model:      device.Model,
			Version:    device.Version,
			Uptime:     int(device.Uptime.Val),
			Clients:    int(device.NumSta.Val),
			Config:     nil,
		})
	}

	return d
}

// Logf logs a message.
func (u *InputUnifi) Logf(msg string, v ...any) {
	webserver.NewInputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "info"},
	})

	if u.Logger != nil {
		u.Logger.Logf(msg, v...)
	}
}

// LogErrorf logs an error message.
func (u *InputUnifi) LogErrorf(msg string, v ...any) {
	webserver.NewInputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "error"},
	})

	if u.Logger != nil {
		u.Logger.LogErrorf(msg, v...)
	}
}

// LogDebugf logs a debug message.
func (u *InputUnifi) LogDebugf(msg string, v ...any) {
	webserver.NewInputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "debug"},
	})

	if u.Logger != nil {
		u.Logger.LogDebugf(msg, v...)
	}
}
