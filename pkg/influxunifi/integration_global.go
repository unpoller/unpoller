package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchDPIApplication generates InfluxDB points for a DPI application catalogue entry.
// These are global (no site) reference records; we emit a presence gauge keyed by ID.
func (u *InfluxUnifi) batchDPIApplication(r report, app *unifi.DPIApplication) {
	if app == nil {
		return
	}

	tags := map[string]string{
		"app_id":   app.ID.Txt,
		"app_name": app.Name,
	}

	r.send(&metric{
		Table:  "dpi_application",
		Tags:   tags,
		Fields: map[string]any{"app_id_val": app.ID.Val},
	})
}

// batchDPICategory generates InfluxDB points for a DPI category catalogue entry.
func (u *InfluxUnifi) batchDPICategory(r report, cat *unifi.DPICategory) {
	if cat == nil {
		return
	}

	tags := map[string]string{
		"cat_id":   cat.ID.Txt,
		"cat_name": cat.Name,
	}

	r.send(&metric{
		Table:  "dpi_category",
		Tags:   tags,
		Fields: map[string]any{"cat_id_val": cat.ID.Val},
	})
}

// batchPendingDevice generates InfluxDB points for a device awaiting adoption.
func (u *InfluxUnifi) batchPendingDevice(r report, d *unifi.PendingDevice) {
	if d == nil {
		return
	}

	tags := map[string]string{
		"mac":   d.MACAddress,
		"model": d.Model,
		"state": d.State,
		"ip":    d.IPAddress,
	}

	firmwareUpdatable := 0
	if d.FirmwareUpdatable {
		firmwareUpdatable = 1
	}

	supported := 0
	if d.Supported {
		supported = 1
	}

	fields := map[string]any{
		"firmware_updatable": firmwareUpdatable,
		"supported":          supported,
		"feature_count":      len(d.Features),
	}

	r.send(&metric{Table: "pending_device", Tags: tags, Fields: fields})
}

// batchCountry generates InfluxDB points for a country entry used in geo-based firewall filters.
// Countries are global reference data with no numeric metrics; we emit a presence gauge.
func (u *InfluxUnifi) batchCountry(r report, c *unifi.Country) {
	if c == nil {
		return
	}

	tags := map[string]string{
		"code": c.Code,
		"name": c.Name,
	}

	r.send(&metric{
		Table:  "country",
		Tags:   tags,
		Fields: map[string]any{"present": 1},
	})
}
