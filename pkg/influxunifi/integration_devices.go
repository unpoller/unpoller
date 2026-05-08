package influxunifi

import (
	"fmt"

	"github.com/unpoller/unifi/v5"
)

// batchIntegrationDeviceStats generates InfluxDB points for Integration/v1 device statistics.
func (u *InfluxUnifi) batchIntegrationDeviceStats(r report, ds *unifi.IntegrationDeviceStats) {
	if ds == nil {
		return
	}

	tags := map[string]string{
		"device_id": ds.DeviceID,
	}

	fields := map[string]any{
		"cpu_utilization_pct":    ds.CPUUtilizationPct.Val,
		"memory_utilization_pct": ds.MemoryUtilizationPct.Val,
		"load_average_1min":      ds.LoadAverage1Min.Val,
		"load_average_5min":      ds.LoadAverage5Min.Val,
		"load_average_15min":     ds.LoadAverage15Min.Val,
		"uptime_sec":             ds.UptimeSec.Val,
	}

	r.send(&metric{Table: "integration_device_stats", Tags: tags, Fields: fields})

	for _, radio := range ds.Radios {
		radioTags := map[string]string{
			"device_id":     ds.DeviceID,
			"frequency_ghz": radio.FrequencyGHz.Txt,
		}

		r.send(&metric{
			Table: "integration_device_radio",
			Tags:  radioTags,
			Fields: map[string]any{
				"tx_retries_pct": radio.TxRetriesPct.Val,
			},
		})
	}

	for i, uplink := range ds.Uplinks {
		uplinkTags := map[string]string{
			"device_id":    ds.DeviceID,
			"uplink_index": fmt.Sprint(i),
		}

		r.send(&metric{
			Table: "integration_device_uplink",
			Tags:  uplinkTags,
			Fields: map[string]any{
				"rx_rate_bps": uplink.RxRateBps.Val,
				"tx_rate_bps": uplink.TxRateBps.Val,
			},
		})
	}
}
