package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

// upsDevice holds Prometheus descriptors for UPS device selector metrics.
// UPSDeviceSelector has no numeric fields — we emit a presence gauge so the
// device appears in Prometheus at all.
type upsDevice struct {
	Presence *prometheus.Desc
}

func descUPSDevice(ns string) *upsDevice {
	labels := []string{"site_name", "mac", "label"}

	return &upsDevice{
		Presence: prometheus.NewDesc(ns+"ups_device_present",
			"UPS device detected on site (always 1 when present)",
			labels, nil),
	}
}

func (u *promUnifi) exportUPSDevice(r report, d *unifi.UPSDeviceSelector) {
	if d == nil {
		return
	}

	labels := []string{d.SiteName, d.MAC, d.Label}

	r.send([]*metric{
		{u.UPSDevice.Presence, gauge, 1.0, labels},
	})
}
