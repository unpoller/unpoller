package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v6"
)

type console struct {
	Info                       *prometheus.Desc
	UpdateAvailable            *prometheus.Desc
	ApplicationInfo            *prometheus.Desc
	ApplicationUpdateAvailable *prometheus.Desc
}

func descConsole(ns string) *console {
	return &console{
		Info: prometheus.NewDesc(ns+"console_info", "UniFi OS console firmware information (always 1)",
			[]string{"source", "name", "version", "available_version", "firmware_status"}, nil),
		UpdateAvailable: prometheus.NewDesc(ns+"console_update_available",
			"UniFi OS update available (1/0); omitted when unknown", []string{"source", "name"}, nil),
		ApplicationInfo: prometheus.NewDesc(ns+"console_application_info",
			"UniFi console application versions (always 1)", []string{"source", "application", "version", "available_version"}, nil),
		ApplicationUpdateAvailable: prometheus.NewDesc(ns+"console_application_update_available",
			"UniFi console application update offered (1); omitted without an offered version", []string{"source", "application"}, nil),
	}
}

func (u *promUnifi) exportConsoleSystem(r report, system *unifi.ConsoleSystem) {
	if system == nil {
		return
	}

	for _, firmware := range system.Devices.UniFiOS {
		if firmware == nil {
			continue
		}

		r.send([]*metric{{u.Console.Info, gauge, 1, []string{
			system.SourceName, firmware.Name, firmware.Firmware, firmware.UpdateAvailable, firmware.FirmwareStatus,
		}}})

		if available, known := firmware.UpdateStatus(); known {
			r.send([]*metric{{u.Console.UpdateAvailable, gauge, available, []string{system.SourceName, firmware.Name}}})
		}
	}

	for _, app := range system.Apps.Controllers {
		if app == nil || app.Name == "" {
			continue
		}

		r.send([]*metric{{u.Console.ApplicationInfo, gauge, 1, []string{
			system.SourceName, app.Name, app.Version, app.UpdateAvailable,
		}}})

		if app.UpdateAvailable != "" {
			r.send([]*metric{{u.Console.ApplicationUpdateAvailable, gauge, 1, []string{system.SourceName, app.Name}}})
		}
	}
}
