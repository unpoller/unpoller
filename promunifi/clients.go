package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type uclient struct {
	Anomalies         *prometheus.Desc
	BytesR            *prometheus.Desc
	CCQ               *prometheus.Desc
	Satisfaction      *prometheus.Desc
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
	TxRetries         *prometheus.Desc
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
	labels := []string{"name", "mac", "site_name", "gw_name", "sw_name", "vlan", "ip", "oui", "network", "sw_port", "ap_name", "wired"}
	labelW := append([]string{"radio_name", "radio", "radio_proto", "channel", "essid", "bssid", "radio_desc"}, labels...)
	return &uclient{
		Anomalies:      prometheus.NewDesc(ns+"anomalies_total", "Client Anomalies", labelW, nil),
		BytesR:         prometheus.NewDesc(ns+"transfer_rate_bytes", "Client Data Rate", labelW, nil),
		CCQ:            prometheus.NewDesc(ns+"ccq_percent", "Client Connection Quality", labelW, nil),
		Satisfaction:   prometheus.NewDesc(ns+"satisfaction_percent", "Client Satisfaction", labelW, nil),
		Noise:          prometheus.NewDesc(ns+"noise_db", "Client AP Noise", labelW, nil),
		RoamCount:      prometheus.NewDesc(ns+"roam_count_total", "Client Roam Counter", labelW, nil),
		RSSI:           prometheus.NewDesc(ns+"rssi_db", "Client RSSI", labelW, nil),
		RxBytes:        prometheus.NewDesc(ns+"receive_bytes_total", "Client Receive Bytes", labels, nil),
		RxBytesR:       prometheus.NewDesc(ns+"receive_rate_bytes", "Client Receive Data Rate", labels, nil),
		RxPackets:      prometheus.NewDesc(ns+"receive_packets_total", "Client Receive Packets", labels, nil),
		RxRate:         prometheus.NewDesc(ns+"radio_receive_rate_bps", "Client Receive Rate", labelW, nil),
		Signal:         prometheus.NewDesc(ns+"radio_signal_db", "Client Signal Strength", labelW, nil),
		TxBytes:        prometheus.NewDesc(ns+"transmit_bytes_total", "Client Transmit Bytes", labels, nil),
		TxBytesR:       prometheus.NewDesc(ns+"transmit_rate_bytes", "Client Transmit Data Rate", labels, nil),
		TxPackets:      prometheus.NewDesc(ns+"transmit_packets_total", "Client Transmit Packets", labels, nil),
		TxRetries:      prometheus.NewDesc(ns+"transmit_retries_total", "Client Transmit Retries", labels, nil),
		TxPower:        prometheus.NewDesc(ns+"radio_transmit_power_dbm", "Client Transmit Power", labelW, nil),
		TxRate:         prometheus.NewDesc(ns+"radio_transmit_rate_bps", "Client Transmit Rate", labelW, nil),
		WifiTxAttempts: prometheus.NewDesc(ns+"wifi_attempts_transmit_total", "Client Wifi Transmit Attempts", labelW, nil),
		Uptime:         prometheus.NewDesc(ns+"uptime_seconds", "Client Uptime", labelW, nil),
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

func (u *promUnifi) exportClient(r report, c *unifi.Client) {
	labels := []string{c.Name, c.Mac, c.SiteName, c.GwName, c.SwName, c.Vlan.Txt, c.IP, c.Oui, c.Network, c.SwPort.Txt, c.ApName, ""}
	labelW := append([]string{c.RadioName, c.Radio, c.RadioProto, c.Channel.Txt, c.Essid, c.Bssid, c.RadioDescription}, labels...)

	if c.IsWired.Val {
		labels[len(labels)-1] = "true"
		labelW[len(labelW)-1] = "true"
		r.send([]*metric{
			{u.Client.RxBytes, prometheus.CounterValue, c.WiredRxBytes, labels},
			{u.Client.RxBytesR, prometheus.GaugeValue, c.WiredRxBytesR, labels},
			{u.Client.RxPackets, prometheus.CounterValue, c.WiredRxPackets, labels},
			{u.Client.TxBytes, prometheus.CounterValue, c.WiredTxBytes, labels},
			{u.Client.TxBytesR, prometheus.GaugeValue, c.WiredTxBytesR, labels},
			{u.Client.TxPackets, prometheus.CounterValue, c.WiredTxPackets, labels},
		})
	} else {
		labels[len(labels)-1] = "false"
		labelW[len(labelW)-1] = "false"
		r.send([]*metric{
			{u.Client.Anomalies, prometheus.CounterValue, c.Anomalies, labelW},
			{u.Client.CCQ, prometheus.GaugeValue, c.Ccq / 10, labelW},
			{u.Client.Satisfaction, prometheus.GaugeValue, c.Satisfaction, labelW},
			{u.Client.Noise, prometheus.GaugeValue, c.Noise, labelW},
			{u.Client.RoamCount, prometheus.CounterValue, c.RoamCount, labelW},
			{u.Client.RSSI, prometheus.GaugeValue, c.Rssi, labelW},
			{u.Client.Signal, prometheus.GaugeValue, c.Signal, labelW},
			{u.Client.TxPower, prometheus.GaugeValue, c.TxPower, labelW},
			{u.Client.TxRate, prometheus.GaugeValue, c.TxRate * 1000, labelW},
			{u.Client.WifiTxAttempts, prometheus.CounterValue, c.WifiTxAttempts, labelW},
			{u.Client.RxRate, prometheus.GaugeValue, c.RxRate * 1000, labelW},
			{u.Client.TxRetries, prometheus.CounterValue, c.TxRetries, labels},
			{u.Client.TxBytes, prometheus.CounterValue, c.TxBytes, labels},
			{u.Client.TxBytesR, prometheus.GaugeValue, c.TxBytesR, labels},
			{u.Client.TxPackets, prometheus.CounterValue, c.TxPackets, labels},
			{u.Client.RxBytes, prometheus.CounterValue, c.RxBytes, labels},
			{u.Client.RxBytesR, prometheus.GaugeValue, c.RxBytesR, labels},
			{u.Client.RxPackets, prometheus.CounterValue, c.RxPackets, labels},
			{u.Client.BytesR, prometheus.GaugeValue, c.BytesR, labelW},
		})
	}
	r.send([]*metric{
		{u.Client.Uptime, prometheus.GaugeValue, c.Uptime, labelW},
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
