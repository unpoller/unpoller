package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

// dpiApplication holds Prometheus descriptors for DPI application catalogue metrics.
// No site label — this is a global catalogue.
type dpiApplication struct {
	Presence *prometheus.Desc
}

func descDPIApplication(ns string) *dpiApplication {
	labels := []string{"app_id", "name"}

	return &dpiApplication{
		Presence: prometheus.NewDesc(ns+"dpi_application_present",
			"DPI application catalogue entry present (always 1)",
			labels, nil),
	}
}

func (u *promUnifi) exportDPIApplication(r report, app *unifi.DPIApplication) {
	if app == nil {
		return
	}

	labels := []string{app.ID.Txt, app.Name}

	r.send([]*metric{
		{u.DPIApplication.Presence, gauge, 1.0, labels},
	})
}

// dpiCategory holds Prometheus descriptors for DPI category catalogue metrics.
// No site label — this is a global catalogue.
type dpiCategory struct {
	Presence *prometheus.Desc
}

func descDPICategory(ns string) *dpiCategory {
	labels := []string{"cat_id", "name"}

	return &dpiCategory{
		Presence: prometheus.NewDesc(ns+"dpi_category_present",
			"DPI category catalogue entry present (always 1)",
			labels, nil),
	}
}

func (u *promUnifi) exportDPICategory(r report, cat *unifi.DPICategory) {
	if cat == nil {
		return
	}

	labels := []string{cat.ID.Txt, cat.Name}

	r.send([]*metric{
		{u.DPICategory.Presence, gauge, 1.0, labels},
	})
}

// pendingDevice holds Prometheus descriptors for pending-adoption device metrics.
// No site label — these are controller-global.
type pendingDevice struct {
	FirmwareUpdatable *prometheus.Desc
	Supported         *prometheus.Desc
}

func descPendingDevice(ns string) *pendingDevice {
	labels := []string{"mac_address", "model", "state", "firmware_version"}

	return &pendingDevice{
		FirmwareUpdatable: prometheus.NewDesc(ns+"pending_device_firmware_updatable",
			"Pending device has a firmware update available (1=yes, 0=no)",
			labels, nil),
		Supported: prometheus.NewDesc(ns+"pending_device_supported",
			"Pending device model is supported by the controller (1=yes, 0=no)",
			labels, nil),
	}
}

func (u *promUnifi) exportPendingDevice(r report, pd *unifi.PendingDevice) {
	if pd == nil {
		return
	}

	firmwareUpdatable := 0.0
	if pd.FirmwareUpdatable {
		firmwareUpdatable = 1.0
	}

	supported := 0.0
	if pd.Supported {
		supported = 1.0
	}

	labels := []string{pd.MACAddress, pd.Model, pd.State, pd.FirmwareVersion}

	r.send([]*metric{
		{u.PendingDevice.FirmwareUpdatable, gauge, firmwareUpdatable, labels},
		{u.PendingDevice.Supported, gauge, supported, labels},
	})
}

// country holds Prometheus descriptors for country list metrics.
// No site label — this is a global geo-filter catalogue.
type country struct {
	Presence *prometheus.Desc
}

func descCountry(ns string) *country {
	labels := []string{"code", "name"}

	return &country{
		Presence: prometheus.NewDesc(ns+"country_present",
			"Country entry present in geo-filter catalogue (always 1)",
			labels, nil),
	}
}

func (u *promUnifi) exportCountry(r report, c *unifi.Country) {
	if c == nil {
		return
	}

	labels := []string{c.Code, c.Name}

	r.send([]*metric{
		{u.Country.Presence, gauge, 1.0, labels},
	})
}
