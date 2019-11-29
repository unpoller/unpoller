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
	if ns += "_client_"; ns == "_client_" {
		ns = "client_"
	}

	labels := []string{"name", "mac", "site_name", "gw_mac", "gw_name", "sw_mac", "sw_name", "vlan", "ip", "oui", "network"}
	labelWired := append([]string{"sw_port"}, labels...)
	labelWireless := append([]string{"ap_mac", "ap_name", "radio_name", "radio", "radio_proto", "channel", "essid", "bssid", "radio_desc"}, labels...)
	wireless := prometheus.Labels{"wired": "false"}
	wired := prometheus.Labels{"wired": "true"}

	return &uclient{
		Anomalies:      prometheus.NewDesc(ns+"anomalies_total", "Client Anomalies", labelWireless, wireless),
		BytesR:         prometheus.NewDesc(ns+"transfer_rate_bytes", "Client Data Rate", labelWireless, wireless),
		CCQ:            prometheus.NewDesc(ns+"ccq_percent", "Client Connection Quality", labelWireless, wireless),
		Noise:          prometheus.NewDesc(ns+"noise_db", "Client AP Noise", labelWireless, wireless),
		RoamCount:      prometheus.NewDesc(ns+"roam_count_total", "Client Roam Counter", labelWireless, wireless),
		RSSI:           prometheus.NewDesc(ns+"rssi_db", "Client RSSI", labelWireless, wireless),
		RxBytes:        prometheus.NewDesc(ns+"receive_bytes_total", "Client Receive Bytes", labelWireless, wireless),
		RxBytesR:       prometheus.NewDesc(ns+"receive_rate_bytes", "Client Receive Data Rate", labelWireless, wireless),
		RxPackets:      prometheus.NewDesc(ns+"receive_packets_total", "Client Receive Packets", labelWireless, wireless),
		RxRate:         prometheus.NewDesc(ns+"radio_receive_rate_bps", "Client Receive Rate", labelWireless, wireless),
		Signal:         prometheus.NewDesc(ns+"radio_signal_db", "Client Signal Strength", labelWireless, wireless),
		TxBytes:        prometheus.NewDesc(ns+"transmit_bytes_total", "Client Transmit Bytes", labelWireless, wireless),
		TxBytesR:       prometheus.NewDesc(ns+"transmit_rate_bytes", "Client Transmit Data Rate", labelWireless, wireless),
		TxPackets:      prometheus.NewDesc(ns+"transmit_packets_total", "Client Transmit Packets", labelWireless, wireless),
		TxPower:        prometheus.NewDesc(ns+"radio_transmit_power_dbm", "Client Transmit Power", labelWireless, wireless),
		TxRate:         prometheus.NewDesc(ns+"radio_transmit_rate_bps", "Client Transmit Rate", labelWireless, wireless),
		WifiTxAttempts: prometheus.NewDesc(ns+"wifi_attempts_transmit_total", "Client Wifi Transmit Attempts", labelWireless, wireless),

		WiredRxBytes:   prometheus.NewDesc(ns+"wired_receive_bytes_total", "Client Wired Receive Bytes", labelWired, wired),
		WiredRxBytesR:  prometheus.NewDesc(ns+"wired_receive_rate_bytes", "Client Wired Receive Data Rate", labelWired, wired),
		WiredRxPackets: prometheus.NewDesc(ns+"wired_receive_packets_total", "Client Wired Receive Packets", labelWired, wired),
		WiredTxBytes:   prometheus.NewDesc(ns+"wired_transmit_bytes_total", "Client Wired Transmit Bytes", labelWired, wired),
		WiredTxBytesR:  prometheus.NewDesc(ns+"wired_transmit_rate_bytes", "Client Wired Data Rate", labelWired, wired),
		WiredTxPackets: prometheus.NewDesc(ns+"wired_transmit_packets_total", "Client Wired Transmit Packets", labelWired, wired),

		Uptime: prometheus.NewDesc(ns+"uptime_seconds", "Client Uptime", labels, nil),
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
	labels := []string{c.Name, c.Mac, c.SiteName, c.GwMac, c.GwName, c.SwMac, c.SwName, c.Vlan.Txt, c.IP, c.Oui, c.Network}
	labelWired := append([]string{c.SwPort.Txt}, labels...)
	labelWireless := append([]string{c.ApMac, c.ApName, c.RadioName, c.Radio, c.RadioProto, c.Channel.Txt, c.Essid, c.Bssid, c.RadioDescription}, labels...)

	if c.IsWired.Val {
		r.send([]*metricExports{
			{u.Client.WiredRxBytes, prometheus.CounterValue, c.WiredRxBytes, labelWired},
			{u.Client.WiredRxBytesR, prometheus.GaugeValue, c.WiredRxBytesR, labelWired},
			{u.Client.WiredRxPackets, prometheus.CounterValue, c.WiredRxPackets, labelWired},
			{u.Client.WiredTxBytes, prometheus.CounterValue, c.WiredTxBytes, labelWired},
			{u.Client.WiredTxBytesR, prometheus.GaugeValue, c.WiredTxBytesR, labelWired},
			{u.Client.WiredTxPackets, prometheus.CounterValue, c.WiredTxPackets, labelWired},
		})
	} else {
		r.send([]*metricExports{
			{u.Client.Anomalies, prometheus.CounterValue, c.Anomalies, labelWireless},
			{u.Client.CCQ, prometheus.GaugeValue, c.Ccq, labelWireless},
			{u.Client.Noise, prometheus.GaugeValue, c.Noise, labelWireless},
			{u.Client.RoamCount, prometheus.CounterValue, c.RoamCount, labelWireless},
			{u.Client.RSSI, prometheus.GaugeValue, c.Rssi, labelWireless},
			{u.Client.Signal, prometheus.GaugeValue, c.Signal, labelWireless},
			{u.Client.TxPower, prometheus.GaugeValue, c.TxPower, labelWireless},
			{u.Client.TxRate, prometheus.GaugeValue, c.TxRate * 1000, labelWireless},
			{u.Client.WifiTxAttempts, prometheus.CounterValue, c.WifiTxAttempts, labelWireless},
			{u.Client.RxRate, prometheus.GaugeValue, c.RxRate * 1000, labelWireless},
			{u.Client.TxBytes, prometheus.CounterValue, c.TxBytes, labelWireless},
			{u.Client.TxBytesR, prometheus.GaugeValue, c.TxBytesR, labelWireless},
			{u.Client.TxPackets, prometheus.CounterValue, c.TxPackets, labelWireless},
			{u.Client.RxBytes, prometheus.CounterValue, c.RxBytes, labelWireless},
			{u.Client.RxBytesR, prometheus.GaugeValue, c.RxBytesR, labelWireless},
			{u.Client.RxPackets, prometheus.CounterValue, c.RxPackets, labelWireless},
			{u.Client.BytesR, prometheus.GaugeValue, c.BytesR, labelWireless},
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
