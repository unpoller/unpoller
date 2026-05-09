package promunifi

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type integrationDevice struct {
	CPUUtilizationPct    *prometheus.Desc
	MemoryUtilizationPct *prometheus.Desc
	LoadAverage1Min      *prometheus.Desc
	LoadAverage5Min      *prometheus.Desc
	LoadAverage15Min     *prometheus.Desc
	UptimeSec            *prometheus.Desc
	RadioTxRetriesPct    *prometheus.Desc
	UplinkRxRateBps      *prometheus.Desc
	UplinkTxRateBps      *prometheus.Desc
}

func descIntegrationDevice(ns string) *integrationDevice {
	labels := []string{"device_id"}
	radioLabels := []string{"device_id", "frequency_ghz"}
	uplinkLabels := []string{"device_id", "uplink_index"}

	return &integrationDevice{
		CPUUtilizationPct:    prometheus.NewDesc(ns+"integration_device_cpu_utilization_pct", "Device CPU utilization percentage (Integration/v1)", labels, nil),
		MemoryUtilizationPct: prometheus.NewDesc(ns+"integration_device_memory_utilization_pct", "Device memory utilization percentage (Integration/v1)", labels, nil),
		LoadAverage1Min:      prometheus.NewDesc(ns+"integration_device_load_average_1min", "Device 1-minute load average (Integration/v1)", labels, nil),
		LoadAverage5Min:      prometheus.NewDesc(ns+"integration_device_load_average_5min", "Device 5-minute load average (Integration/v1)", labels, nil),
		LoadAverage15Min:     prometheus.NewDesc(ns+"integration_device_load_average_15min", "Device 15-minute load average (Integration/v1)", labels, nil),
		UptimeSec:            prometheus.NewDesc(ns+"integration_device_uptime_seconds", "Device uptime in seconds (Integration/v1)", labels, nil),
		RadioTxRetriesPct:    prometheus.NewDesc(ns+"integration_device_radio_tx_retries_pct", "Per-radio TX retry percentage (Integration/v1)", radioLabels, nil),
		UplinkRxRateBps:      prometheus.NewDesc(ns+"integration_device_uplink_rx_rate_bps", "Per-uplink receive rate in bps (Integration/v1)", uplinkLabels, nil),
		UplinkTxRateBps:      prometheus.NewDesc(ns+"integration_device_uplink_tx_rate_bps", "Per-uplink transmit rate in bps (Integration/v1)", uplinkLabels, nil),
	}
}

func (u *promUnifi) exportIntegrationDeviceStats(r report, ds *unifi.IntegrationDeviceStats) {
	if ds == nil {
		return
	}

	labels := []string{ds.DeviceID}

	r.send([]*metric{
		{u.IntegrationDevice.CPUUtilizationPct, gauge, ds.CPUUtilizationPct, labels},
		{u.IntegrationDevice.MemoryUtilizationPct, gauge, ds.MemoryUtilizationPct, labels},
		{u.IntegrationDevice.LoadAverage1Min, gauge, ds.LoadAverage1Min, labels},
		{u.IntegrationDevice.LoadAverage5Min, gauge, ds.LoadAverage5Min, labels},
		{u.IntegrationDevice.LoadAverage15Min, gauge, ds.LoadAverage15Min, labels},
		{u.IntegrationDevice.UptimeSec, gauge, ds.UptimeSec, labels},
	})

	for _, radio := range ds.Radios {
		radioLabels := []string{ds.DeviceID, radio.FrequencyGHz.Txt}

		r.send([]*metric{
			{u.IntegrationDevice.RadioTxRetriesPct, gauge, radio.TxRetriesPct, radioLabels},
		})
	}

	for i, uplink := range ds.Uplinks {
		uplinkLabels := []string{ds.DeviceID, fmt.Sprint(i)}

		r.send([]*metric{
			{u.IntegrationDevice.UplinkRxRateBps, gauge, uplink.RxRateBps, uplinkLabels},
			{u.IntegrationDevice.UplinkTxRateBps, gauge, uplink.TxRateBps, uplinkLabels},
		})
	}
}
