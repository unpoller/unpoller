package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type client struct {
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

func descClient(ns string) *client {
	if ns += "_client_"; ns == "_client_" {
		ns = "client_"
	}

	labels := []string{"id", "mac", "user_id", "site_id", "site_name",
		"network_id", "ap_mac", "gw_mac", "sw_mac", "ap_name", "gw_name",
		"sw_name", "radio_name", "radio", "radio_proto", "name", "channel",
		"vlan", "ip", "essid", "bssid", "radio_desc"}

	return &client{
		Anomalies:         prometheus.NewDesc(ns+"anomalies_total", "Client Anomalies", labels, nil),
		BytesR:            prometheus.NewDesc(ns+"bytes_rate", "Client Data Rate", labels, nil),
		CCQ:               prometheus.NewDesc(ns+"ccq", "Client Connection Quality", labels, nil),
		Noise:             prometheus.NewDesc(ns+"noise", "Client AP Noise", labels, nil),
		RoamCount:         prometheus.NewDesc(ns+"roam_total", "Client Roam Counter", labels, nil),
		RSSI:              prometheus.NewDesc(ns+"rssi", "Client RSSI", labels, nil),
		RxBytes:           prometheus.NewDesc(ns+"bytes_rx_total", "Client Receive Bytes", labels, nil),
		RxBytesR:          prometheus.NewDesc(ns+"bytes_rx_rate", "Client Receive Data Rate", labels, nil),
		RxPackets:         prometheus.NewDesc(ns+"packets_rx_total", "Client Receive Packets", labels, nil),
		RxRate:            prometheus.NewDesc(ns+"rate", "Client Receive Rate", labels, nil),
		Signal:            prometheus.NewDesc(ns+"signal", "Client Signal Strength", labels, nil),
		TxBytes:           prometheus.NewDesc(ns+"bytes_tx_total", "Client Transmit Bytes", labels, nil),
		TxBytesR:          prometheus.NewDesc(ns+"bytes_tx_rate", "Client Transmit Data Rate", labels, nil),
		TxPackets:         prometheus.NewDesc(ns+"packets_tx_total", "Client Transmit Packets", labels, nil),
		TxPower:           prometheus.NewDesc(ns+"tx_power", "Client Transmit Power", labels, nil),
		TxRate:            prometheus.NewDesc(ns+"tx_rate", "Client Transmit Rate", labels, nil),
		Uptime:            prometheus.NewDesc(ns+"uptime", "Client Uptime", labels, nil),
		WifiTxAttempts:    prometheus.NewDesc(ns+"wifi_attempts_tx_total", "Client Wifi Transmit Attempts", labels, nil),
		WiredRxBytes:      prometheus.NewDesc(ns+"wired_bytes_rx_total", "Client Wired Receive Bytes", labels, nil),
		WiredRxBytesR:     prometheus.NewDesc(ns+"wired_bytes_rx_rate", "Client Wired Receive Data Rate", labels, nil),
		WiredRxPackets:    prometheus.NewDesc(ns+"wired_packets_rx_total", "Client Wired Receive Packets", labels, nil),
		WiredTxBytes:      prometheus.NewDesc(ns+"wired_bytes_tx_total", "Client Wired Transmit Bytes", labels, nil),
		WiredTxBytesR:     prometheus.NewDesc(ns+"wired_bytes_tx_rate", "Client Wired Data Rate", labels, nil),
		WiredTxPackets:    prometheus.NewDesc(ns+"wired_packets_tx_total", "Client Wired Transmit Packets", labels, nil),
		DpiStatsApp:       prometheus.NewDesc(ns+"dpi_stats_app", "Client DPI Stats App", labels, nil),
		DpiStatsCat:       prometheus.NewDesc(ns+"dpi_stats_cat", "Client DPI Stats Cat", labels, nil),
		DpiStatsRxBytes:   prometheus.NewDesc(ns+"dpi_stats_bytes_rx_total", "Client DPI Stats Receive Bytes", labels, nil),
		DpiStatsRxPackets: prometheus.NewDesc(ns+"dpi_stats_packets_rx_total", "Client DPI Stats Receive Packets", labels, nil),
		DpiStatsTxBytes:   prometheus.NewDesc(ns+"dpi_stats_bytes_tx_total", "Client DPI Stats Transmit Bytes", labels, nil),
		DpiStatsTxPackets: prometheus.NewDesc(ns+"dpi_stats_packets_tx_total", "Client DPI Stats Transmit Packets", labels, nil),
	}
}

// CollectClient exports Clients' Data
func (u *unifiCollector) exportClient(c *unifi.Client) []*metricExports {
	labels := []string{c.ID, c.Mac, c.UserID, c.SiteID, c.SiteName,
		c.NetworkID, c.ApMac, c.GwMac, c.SwMac, c.ApName, c.GwName,
		c.SwName, c.RadioName, c.Radio, c.RadioProto, c.Name, c.Channel.Txt,
		c.Vlan.Txt, c.IP, c.Essid, c.Bssid, c.RadioDescription,
	}

	return []*metricExports{
		{u.Client.Anomalies, prometheus.CounterValue, c.Anomalies, labels},
		{u.Client.BytesR, prometheus.GaugeValue, c.BytesR, labels},
		{u.Client.CCQ, prometheus.GaugeValue, c.Ccq, labels},
		{u.Client.Noise, prometheus.GaugeValue, c.Noise, labels},
		{u.Client.RoamCount, prometheus.CounterValue, c.RoamCount, labels},
		{u.Client.RSSI, prometheus.GaugeValue, c.Rssi, labels},
		{u.Client.RxBytes, prometheus.CounterValue, c.RxBytes, labels},
		{u.Client.RxBytesR, prometheus.GaugeValue, c.RxBytesR, labels},
		{u.Client.RxPackets, prometheus.CounterValue, c.RxPackets, labels},
		{u.Client.RxRate, prometheus.GaugeValue, c.RxRate, labels},
		{u.Client.Signal, prometheus.GaugeValue, c.Signal, labels},
		{u.Client.TxBytes, prometheus.CounterValue, c.TxBytes, labels},
		{u.Client.TxBytesR, prometheus.GaugeValue, c.TxBytesR, labels},
		{u.Client.TxPackets, prometheus.CounterValue, c.TxPackets, labels},
		{u.Client.TxPower, prometheus.GaugeValue, c.TxPower, labels},
		{u.Client.TxRate, prometheus.GaugeValue, c.TxRate, labels},
		{u.Client.Uptime, prometheus.GaugeValue, c.Uptime, labels},
		{u.Client.WifiTxAttempts, prometheus.CounterValue, c.WifiTxAttempts, labels},
		{u.Client.WiredRxBytes, prometheus.CounterValue, c.WiredRxBytes, labels},
		{u.Client.WiredRxBytesR, prometheus.GaugeValue, c.WiredRxBytesR, labels},
		{u.Client.WiredRxPackets, prometheus.CounterValue, c.WiredRxPackets, labels},
		{u.Client.WiredTxBytes, prometheus.CounterValue, c.WiredTxBytes, labels},
		{u.Client.WiredTxBytesR, prometheus.GaugeValue, c.WiredTxBytesR, labels},
		{u.Client.WiredTxPackets, prometheus.CounterValue, c.WiredTxPackets, labels},
		{u.Client.DpiStatsApp, prometheus.GaugeValue, c.DpiStats.App.Val, labels},
		{u.Client.DpiStatsCat, prometheus.GaugeValue, c.DpiStats.Cat.Val, labels},
		{u.Client.DpiStatsRxBytes, prometheus.CounterValue, c.DpiStats.RxBytes.Val, labels},
		{u.Client.DpiStatsRxPackets, prometheus.CounterValue, c.DpiStats.RxPackets.Val, labels},
		{u.Client.DpiStatsTxBytes, prometheus.CounterValue, c.DpiStats.TxBytes.Val, labels},
		{u.Client.DpiStatsTxPackets, prometheus.CounterValue, c.DpiStats.TxPackets.Val, labels},
	}
}
