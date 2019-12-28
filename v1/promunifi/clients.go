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
		Anomalies:      prometheus.NewDesc(ns+"anomalies", "Client Anomalies", labelW, nil),
		BytesR:         prometheus.NewDesc(ns+"transfer_rate_bytes", "Client Data Rate", labelW, nil),
		CCQ:            prometheus.NewDesc(ns+"ccq_ratio", "Client Connection Quality", labelW, nil),
		Satisfaction:   prometheus.NewDesc(ns+"satisfaction_ratio", "Client Satisfaction", labelW, nil),
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
		Uptime:         prometheus.NewDesc(ns+"uptime_seconds", "Client Uptime", labelW, nil), // XXX: re-purpose for info tags.
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
			{u.Client.RxBytes, counter, c.WiredRxBytes, labels},
			{u.Client.RxBytesR, gauge, c.WiredRxBytesR, labels},
			{u.Client.RxPackets, counter, c.WiredRxPackets, labels},
			{u.Client.TxBytes, counter, c.WiredTxBytes, labels},
			{u.Client.TxBytesR, gauge, c.WiredTxBytesR, labels},
			{u.Client.TxPackets, counter, c.WiredTxPackets, labels},
		})
	} else {
		labels[len(labels)-1] = "false"
		labelW[len(labelW)-1] = "false"
		r.send([]*metric{
			{u.Client.Anomalies, counter, c.Anomalies, labelW},
			{u.Client.CCQ, gauge, float64(c.Ccq) / 1000.0, labelW},
			{u.Client.Satisfaction, gauge, c.Satisfaction.Val / 100.0, labelW},
			{u.Client.Noise, gauge, c.Noise, labelW},
			{u.Client.RoamCount, counter, c.RoamCount, labelW},
			{u.Client.RSSI, gauge, c.Rssi, labelW},
			{u.Client.Signal, gauge, c.Signal, labelW},
			{u.Client.TxPower, gauge, c.TxPower, labelW},
			{u.Client.TxRate, gauge, c.TxRate * 1000, labelW},
			{u.Client.WifiTxAttempts, counter, c.WifiTxAttempts, labelW},
			{u.Client.RxRate, gauge, c.RxRate * 1000, labelW},
			{u.Client.TxRetries, counter, c.TxRetries, labels},
			{u.Client.TxBytes, counter, c.TxBytes, labels},
			{u.Client.TxBytesR, gauge, c.TxBytesR, labels},
			{u.Client.TxPackets, counter, c.TxPackets, labels},
			{u.Client.RxBytes, counter, c.RxBytes, labels},
			{u.Client.RxBytesR, gauge, c.RxBytesR, labels},
			{u.Client.RxPackets, counter, c.RxPackets, labels},
			{u.Client.BytesR, gauge, c.BytesR, labelW},
		})
	}
	r.send([]*metric{{u.Client.Uptime, gauge, c.Uptime, labelW}})
	/* needs more "looking into"
	{u.Client.DpiStatsApp, gauge, c.DpiStats.App, labels},
	{u.Client.DpiStatsCat, gauge, c.DpiStats.Cat, labels},
	{u.Client.DpiStatsRxBytes, counter, c.DpiStats.RxBytes, labels},
	{u.Client.DpiStatsRxPackets, counter, c.DpiStats.RxPackets, labels},
	{u.Client.DpiStatsTxBytes, counter, c.DpiStats.TxBytes, labels},
	{u.Client.DpiStatsTxPackets, counter, c.DpiStats.TxPackets, labels},
	*/
}
