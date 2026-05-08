package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchIntegrationDevStats generates integration device statistics datapoints for Datadog.
func (u *DatadogUnifi) batchIntegrationDevStats(r report, ids *unifi.IntegrationDeviceStats) {
	if ids == nil {
		return
	}

	metricName := metricNamespace("integration_device")

	tags := cleanTags(map[string]string{
		"device_id": ids.DeviceID,
	})

	tagSlice := tagMapToTags(tags)

	_ = r.reportGauge(metricName("cpu_utilization_pct"), ids.CPUUtilizationPct.Val, tagSlice)
	_ = r.reportGauge(metricName("memory_utilization_pct"), ids.MemoryUtilizationPct.Val, tagSlice)
	_ = r.reportGauge(metricName("uptime_sec"), ids.UptimeSec.Val, tagSlice)
	_ = r.reportGauge(metricName("load_average_1min"), ids.LoadAverage1Min.Val, tagSlice)
	_ = r.reportGauge(metricName("load_average_5min"), ids.LoadAverage5Min.Val, tagSlice)
	_ = r.reportGauge(metricName("load_average_15min"), ids.LoadAverage15Min.Val, tagSlice)

	radioMetric := metricNamespace("integration_device_radio")

	for i := range ids.Radios {
		radio := &ids.Radios[i]

		radioTags := append(tagSlice,
			tag("frequency_ghz", radio.FrequencyGHz.Val),
		)

		_ = r.reportGauge(radioMetric("tx_retries_pct"), radio.TxRetriesPct.Val, radioTags)
	}

	uplinkMetric := metricNamespace("integration_device_uplink")

	for i := range ids.Uplinks {
		uplink := &ids.Uplinks[i]

		uplinkTags := append(tagSlice,
			tag("uplink_index", i),
		)

		_ = r.reportGauge(uplinkMetric("rx_rate_bps"), uplink.RxRateBps.Val, uplinkTags)
		_ = r.reportGauge(uplinkMetric("tx_rate_bps"), uplink.TxRateBps.Val, uplinkTags)
	}
}
