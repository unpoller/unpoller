package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type uclient struct {
	Anomalies         *prometheus.Desc
	BytesR            *prometheus.Desc
	CCQ               *prometheus.Desc
	Noise             *prometheus.Desc
	RoamCount         *prometheus.Desc
	RSSI              *prometheus.Desc
	RxBytes           *prometheus.Desc
	RxBytesR          *prometheus.Desc
	RxPackets         *prometheus.Desc
	RxRate            *prometheus.Desc
	Signal            *prometheus.Desc
	TxBytes           *prometheus.Desc
	TxBytesR          *prometheus.Desc
	TxPackets         *prometheus.Desc
	TxPower           *prometheus.Desc
	TxRate            *prometheus.Desc
	Uptime            *prometheus.Desc
	WifiTxAttempts    *prometheus.Desc
	WiredRxBytes      *prometheus.Desc
	WiredRxBytesR     *prometheus.Desc
	WiredRxPackets    *prometheus.Desc
	WiredTxBytes      *prometheus.Desc
	WiredTxBytesR     *prometheus.Desc
	WiredTxPackets    *prometheus.Desc
	DpiStatsApp       *prometheus.Desc
	DpiStatsCat       *prometheus.Desc
	DpiStatsRxBytes   *prometheus.Desc
	DpiStatsRxPackets *prometheus.Desc
	DpiStatsTxBytes   *prometheus.Desc
	DpiStatsTxPackets *prometheus.Desc
}

func descClient(ns string) *uclient {
	labels := []string{"name", "mac", "site_name", "gw_mac", "gw_name", "sw_mac", "sw_name", "vlan", "ip", "oui", "network", "sw_port",
		"ap_mac", "ap_name", "radio_name", "radio", "radio_proto", "channel", "essid", "bssid", "radio_desc", "wired"}
	return &uclient{
		Anomalies:      prometheus.NewDesc(ns+"anomalies_total", "Client Anomalies", labels, nil),
		BytesR:         prometheus.NewDesc(ns+"transfer_rate_bytes", "Client Data Rate", labels, nil),
		CCQ:            prometheus.NewDesc(ns+"ccq_percent", "Client Connection Quality", labels, nil),
		Noise:          prometheus.NewDesc(ns+"noise_db", "Client AP Noise", labels, nil),
		RoamCount:      prometheus.NewDesc(ns+"roam_count_total", "Client Roam Counter", labels, nil),
		RSSI:           prometheus.NewDesc(ns+"rssi_db", "Client RSSI", labels, nil),
		RxBytes:        prometheus.NewDesc(ns+"receive_bytes_total", "Client Receive Bytes", labels, nil),
		RxBytesR:       prometheus.NewDesc(ns+"receive_rate_bytes", "Client Receive Data Rate", labels, nil),
		RxPackets:      prometheus.NewDesc(ns+"receive_packets_total", "Client Receive Packets", labels, nil),
		RxRate:         prometheus.NewDesc(ns+"radio_receive_rate_bps", "Client Receive Rate", labels, nil),
		Signal:         prometheus.NewDesc(ns+"radio_signal_db", "Client Signal Strength", labels, nil),
		TxBytes:        prometheus.NewDesc(ns+"transmit_bytes_total", "Client Transmit Bytes", labels, nil),
		TxBytesR:       prometheus.NewDesc(ns+"transmit_rate_bytes", "Client Transmit Data Rate", labels, nil),
		TxPackets:      prometheus.NewDesc(ns+"transmit_packets_total", "Client Transmit Packets", labels, nil),
		TxPower:        prometheus.NewDesc(ns+"radio_transmit_power_dbm", "Client Transmit Power", labels, nil),
		TxRate:         prometheus.NewDesc(ns+"radio_transmit_rate_bps", "Client Transmit Rate", labels, nil),
		WifiTxAttempts: prometheus.NewDesc(ns+"wifi_attempts_transmit_total", "Client Wifi Transmit Attempts", labels, nil),
		Uptime:         prometheus.NewDesc(ns+"uptime_seconds", "Client Uptime", labels, nil),
		/* needs more "looking into"
		DpiStatsApp:       prometheus.NewDesc(ns+"dpi_stats_app", "Client DPI Stats App", labels, nil),
		DpiStatsCat:       prometheus.NewDesc(ns+"dpi_stats_cat", "Client DPI Stats Cat", labels, nil),
		DpiStatsRxBytes:   prometheus.NewDesc(ns+"dpi_stats_receive_bytes_total", "Client DPI Stats Receive Bytes", labels, nil),
		DpiStatsRxPackets: prometheus.NewDesc(ns+"dpi_stats_receive_packets_total", "Client DPI Stats Receive Packets", labels, nil),
		DpiStatsTxBytes:   prometheus.NewDesc(ns+"dpi_stats_transmit_bytes_total", "Client DPI Stats Transmit Bytes", labels, nil),
		DpiStatsTxPackets: prometheus.NewDesc(ns+"dpi_stats_transmit_packets_total", "Client DPI Stats Transmit Packets", labels, nil),
		*/
	}
}

func (u *unifiCollector) exportClients(r report) {
	if r.metrics() == nil || len(r.metrics().Clients) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, c := range r.metrics().Clients {
			u.exportClient(r, c)
		}
	}()
}

func (u *unifiCollector) exportClient(r report, c *unifi.Client) {
	labels := []string{c.Name, c.Mac, c.SiteName, c.GwMac, c.GwName, c.SwMac, c.SwName, c.Vlan.Txt, c.IP, c.Oui, c.Network, c.SwPort.Txt,
		c.ApMac, c.ApName, c.RadioName, c.Radio, c.RadioProto, c.Channel.Txt, c.Essid, c.Bssid, c.RadioDescription, "false"}

	if c.IsWired.Val {
		labels[len(labels)-1] = "true"
		r.send([]*metricExports{
			{u.Client.RxBytes, prometheus.CounterValue, c.WiredRxBytes, labels},
			{u.Client.RxBytesR, prometheus.GaugeValue, c.WiredRxBytesR, labels},
			{u.Client.RxPackets, prometheus.CounterValue, c.WiredRxPackets, labels},
			{u.Client.TxBytes, prometheus.CounterValue, c.WiredTxBytes, labels},
			{u.Client.TxBytesR, prometheus.GaugeValue, c.WiredTxBytesR, labels},
			{u.Client.TxPackets, prometheus.CounterValue, c.WiredTxPackets, labels},
		})
	} else {
		labels[len(labels)-1] = "false"
		r.send([]*metricExports{
			{u.Client.Anomalies, prometheus.CounterValue, c.Anomalies, labels},
			{u.Client.CCQ, prometheus.GaugeValue, c.Ccq / 10, labels},
			{u.Client.Noise, prometheus.GaugeValue, c.Noise, labels},
			{u.Client.RoamCount, prometheus.CounterValue, c.RoamCount, labels},
			{u.Client.RSSI, prometheus.GaugeValue, c.Rssi, labels},
			{u.Client.Signal, prometheus.GaugeValue, c.Signal, labels},
			{u.Client.TxPower, prometheus.GaugeValue, c.TxPower, labels},
			{u.Client.TxRate, prometheus.GaugeValue, c.TxRate * 1000, labels},
			{u.Client.WifiTxAttempts, prometheus.CounterValue, c.WifiTxAttempts, labels},
			{u.Client.RxRate, prometheus.GaugeValue, c.RxRate * 1000, labels},
			{u.Client.TxBytes, prometheus.CounterValue, c.TxBytes, labels},
			{u.Client.TxBytesR, prometheus.GaugeValue, c.TxBytesR, labels},
			{u.Client.TxPackets, prometheus.CounterValue, c.TxPackets, labels},
			{u.Client.RxBytes, prometheus.CounterValue, c.RxBytes, labels},
			{u.Client.RxBytesR, prometheus.GaugeValue, c.RxBytesR, labels},
			{u.Client.RxPackets, prometheus.CounterValue, c.RxPackets, labels},
			{u.Client.BytesR, prometheus.GaugeValue, c.BytesR, labels},
		})
	}
	r.send([]*metricExports{
		{u.Client.Uptime, prometheus.GaugeValue, c.Uptime, labels},
		/* needs more "looking into"
		{u.Client.DpiStatsApp, prometheus.GaugeValue, c.DpiStats.App, labels},
		{u.Client.DpiStatsCat, prometheus.GaugeValue, c.DpiStats.Cat, labels},
		{u.Client.DpiStatsRxBytes, prometheus.CounterValue, c.DpiStats.RxBytes, labels},
		{u.Client.DpiStatsRxPackets, prometheus.CounterValue, c.DpiStats.RxPackets, labels},
		{u.Client.DpiStatsTxBytes, prometheus.CounterValue, c.DpiStats.TxBytes, labels},
		{u.Client.DpiStatsTxPackets, prometheus.CounterValue, c.DpiStats.TxPackets, labels},
		*/
	})
}
