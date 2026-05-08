package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchDPIApplication generates DPI application catalogue datapoints for Datadog.
// DPIApplications are global (no site scope).
func (u *DatadogUnifi) batchDPIApplication(r report, da *unifi.DPIApplication) {
	if da == nil {
		return
	}

	metricName := metricNamespace("dpi_application")

	tags := cleanTags(map[string]string{
		"name": da.Name,
	})

	_ = r.reportGauge(metricName("id"), da.ID.Val, tagMapToTags(tags))
}

// batchDPICategory generates DPI category catalogue datapoints for Datadog.
// DPICategories are global (no site scope).
func (u *DatadogUnifi) batchDPICategory(r report, dc *unifi.DPICategory) {
	if dc == nil {
		return
	}

	metricName := metricNamespace("dpi_category")

	tags := cleanTags(map[string]string{
		"name": dc.Name,
	})

	_ = r.reportGauge(metricName("id"), dc.ID.Val, tagMapToTags(tags))
}

// batchPendingDevice generates PendingDevice (adoption queue) datapoints for Datadog.
// PendingDevices are global (no site scope).
func (u *DatadogUnifi) batchPendingDevice(r report, pd *unifi.PendingDevice) {
	if pd == nil {
		return
	}

	metricName := metricNamespace("pending_device")

	tags := cleanTags(map[string]string{
		"mac_address": pd.MACAddress,
		"ip_address":  pd.IPAddress,
		"model":       pd.Model,
		"state":       pd.State,
	})

	_ = r.reportGauge(metricName("supported"), boolToFloat64(pd.Supported), tagMapToTags(tags))
	_ = r.reportGauge(metricName("firmware_updatable"), boolToFloat64(pd.FirmwareUpdatable), tagMapToTags(tags))
}

// batchCountry generates Country datapoints for Datadog.
// Countries are global (no site scope).
func (u *DatadogUnifi) batchCountry(r report, co *unifi.Country) {
	if co == nil {
		return
	}

	metricName := metricNamespace("country")

	tags := cleanTags(map[string]string{
		"code": co.Code,
		"name": co.Name,
	})

	// Emit a presence gauge (1.0 = country entry exists in the catalogue).
	_ = r.reportGauge(metricName("present"), 1.0, tagMapToTags(tags))
}
