package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type uap struct {
	// Ap Traffic Stats
	ApWifiTxDropped     *prometheus.Desc
	ApRxErrors          *prometheus.Desc
	ApRxDropped         *prometheus.Desc
	ApRxFrags           *prometheus.Desc
	ApRxCrypts          *prometheus.Desc
	ApTxPackets         *prometheus.Desc
	ApTxBytes           *prometheus.Desc
	ApTxErrors          *prometheus.Desc
	ApTxDropped         *prometheus.Desc
	ApTxRetries         *prometheus.Desc
	ApRxPackets         *prometheus.Desc
	ApRxBytes           *prometheus.Desc
	WifiTxAttempts      *prometheus.Desc
	MacFilterRejections *prometheus.Desc
	// VAP Stats
	VAPCcq                   *prometheus.Desc
	VAPMacFilterRejections   *prometheus.Desc
	VAPNumSatisfactionSta    *prometheus.Desc
	VAPAvgClientSignal       *prometheus.Desc
	VAPSatisfaction          *prometheus.Desc
	VAPSatisfactionNow       *prometheus.Desc
	VAPDNSAvgLatency         *prometheus.Desc
	VAPRxBytes               *prometheus.Desc
	VAPRxCrypts              *prometheus.Desc
	VAPRxDropped             *prometheus.Desc
	VAPRxErrors              *prometheus.Desc
	VAPRxFrags               *prometheus.Desc
	VAPRxNwids               *prometheus.Desc
	VAPRxPackets             *prometheus.Desc
	VAPTxBytes               *prometheus.Desc
	VAPTxDropped             *prometheus.Desc
	VAPTxErrors              *prometheus.Desc
	VAPTxPackets             *prometheus.Desc
	VAPTxPower               *prometheus.Desc
	VAPTxRetries             *prometheus.Desc
	VAPTxCombinedRetries     *prometheus.Desc
	VAPTxDataMpduBytes       *prometheus.Desc
	VAPTxRtsRetries          *prometheus.Desc
	VAPTxSuccess             *prometheus.Desc
	VAPTxTotal               *prometheus.Desc
	VAPTxGoodbytes           *prometheus.Desc
	VAPTxLatAvg              *prometheus.Desc
	VAPTxLatMax              *prometheus.Desc
	VAPTxLatMin              *prometheus.Desc
	VAPRxGoodbytes           *prometheus.Desc
	VAPRxLatAvg              *prometheus.Desc
	VAPRxLatMax              *prometheus.Desc
	VAPRxLatMin              *prometheus.Desc
	VAPWifiTxLatencyMovAvg   *prometheus.Desc
	VAPWifiTxLatencyMovMax   *prometheus.Desc
	VAPWifiTxLatencyMovMin   *prometheus.Desc
	VAPWifiTxLatencyMovTotal *prometheus.Desc
	VAPWifiTxLatencyMovCount *prometheus.Desc
	// Radio Stats
	RadioCurrentAntennaGain *prometheus.Desc
	RadioHt                 *prometheus.Desc
	RadioMaxTxpower         *prometheus.Desc
	RadioMinTxpower         *prometheus.Desc
	RadioNss                *prometheus.Desc
	RadioRadioCaps          *prometheus.Desc
	RadioTxPower            *prometheus.Desc
	RadioAstBeXmit          *prometheus.Desc
	RadioChannel            *prometheus.Desc
	RadioCuSelfRx           *prometheus.Desc
	RadioCuSelfTx           *prometheus.Desc
	RadioExtchannel         *prometheus.Desc
	RadioGain               *prometheus.Desc
	RadioGuestNumSta        *prometheus.Desc
	RadioNumSta             *prometheus.Desc
	RadioUserNumSta         *prometheus.Desc
	RadioTxPackets          *prometheus.Desc
	RadioTxRetries          *prometheus.Desc
}

func descUAP(ns string) *uap {
	//	labels := []string{"ip", "version", "model", "serial", "type", "mac", "site_name", "name"}
	labelA := []string{"stat", "site_name", "name"} // stat + labels[6:]
	labelV := []string{"vap_name", "bssid", "radio", "radio_name", "essid", "usage", "site_name", "name"}
	labelR := []string{"radio_name", "radio", "site_name", "name"}
	return &uap{
		// 3x each - stat table: total, guest, user
		ApWifiTxDropped:     prometheus.NewDesc(ns+"stat_wifi_transmt_dropped_total", "Wifi Transmissions Dropped", labelA, nil),
		ApRxErrors:          prometheus.NewDesc(ns+"stat_receive_errors_total", "Receive Errors", labelA, nil),
		ApRxDropped:         prometheus.NewDesc(ns+"stat_receive_dropped_total", "Receive Dropped", labelA, nil),
		ApRxFrags:           prometheus.NewDesc(ns+"stat_receive_frags_total", "Received Frags", labelA, nil),
		ApRxCrypts:          prometheus.NewDesc(ns+"stat_receive_crypts_total", "Receive Crypts", labelA, nil),
		ApTxPackets:         prometheus.NewDesc(ns+"stat_transmit_packets_total", "Transmit Packets", labelA, nil),
		ApTxBytes:           prometheus.NewDesc(ns+"stat_transmit_bytes_total", "Transmit Bytes", labelA, nil),
		ApTxErrors:          prometheus.NewDesc(ns+"stat_transmit_errors_total", "Transmit Errors", labelA, nil),
		ApTxDropped:         prometheus.NewDesc(ns+"stat_transmit_dropped_total", "Transmit Dropped", labelA, nil),
		ApTxRetries:         prometheus.NewDesc(ns+"stat_retries_tx_total", "Transmit Retries", labelA, nil),
		ApRxPackets:         prometheus.NewDesc(ns+"stat_receive_packets_total", "Receive Packets", labelA, nil),
		ApRxBytes:           prometheus.NewDesc(ns+"stat_receive_bytes_total", "Receive Bytes", labelA, nil),
		WifiTxAttempts:      prometheus.NewDesc(ns+"stat_wifi_transmit_attempts_total", "Wifi Transmission Attempts", labelA, nil),
		MacFilterRejections: prometheus.NewDesc(ns+"stat_mac_filter_rejects_total", "MAC Filter Rejections", labelA, nil),
		// N each - 1 per Virtual AP (VAP)
		VAPCcq:                   prometheus.NewDesc(ns+"vap_ccq_ratio", "VAP Client Connection Quality", labelV, nil),
		VAPMacFilterRejections:   prometheus.NewDesc(ns+"vap_mac_filter_rejects_total", "VAP MAC Filter Rejections", labelV, nil),
		VAPNumSatisfactionSta:    prometheus.NewDesc(ns+"vap_satisfaction_stations", "VAP Number Satisifaction Stations", labelV, nil),
		VAPAvgClientSignal:       prometheus.NewDesc(ns+"vap_average_client_signal", "VAP Average Client Signal", labelV, nil),
		VAPSatisfaction:          prometheus.NewDesc(ns+"vap_satisfaction_ratio", "VAP Satisfaction", labelV, nil),
		VAPSatisfactionNow:       prometheus.NewDesc(ns+"vap_satisfaction_now_ratio", "VAP Satisfaction Now", labelV, nil),
		VAPDNSAvgLatency:         prometheus.NewDesc(ns+"vap_dns_latency_average_seconds", "VAP DNS Latency Average", labelV, nil),
		VAPRxBytes:               prometheus.NewDesc(ns+"vap_receive_bytes_total", "VAP Bytes Received", labelV, nil),
		VAPRxCrypts:              prometheus.NewDesc(ns+"vap_receive_crypts_total", "VAP Crypts Received", labelV, nil),
		VAPRxDropped:             prometheus.NewDesc(ns+"vap_receive_dropped_total", "VAP Dropped Received", labelV, nil),
		VAPRxErrors:              prometheus.NewDesc(ns+"vap_receive_errors_total", "VAP Errors Received", labelV, nil),
		VAPRxFrags:               prometheus.NewDesc(ns+"vap_receive_frags_total", "VAP Frags Received", labelV, nil),
		VAPRxNwids:               prometheus.NewDesc(ns+"vap_receive_nwids_total", "VAP Nwids Received", labelV, nil),
		VAPRxPackets:             prometheus.NewDesc(ns+"vap_receive_packets_total", "VAP Packets Received", labelV, nil),
		VAPTxBytes:               prometheus.NewDesc(ns+"vap_transmit_bytes_total", "VAP Bytes Transmitted", labelV, nil),
		VAPTxDropped:             prometheus.NewDesc(ns+"vap_transmit_dropped_total", "VAP Dropped Transmitted", labelV, nil),
		VAPTxErrors:              prometheus.NewDesc(ns+"vap_transmit_errors_total", "VAP Errors Transmitted", labelV, nil),
		VAPTxPackets:             prometheus.NewDesc(ns+"vap_transmit_packets_total", "VAP Packets Transmitted", labelV, nil),
		VAPTxPower:               prometheus.NewDesc(ns+"vap_transmit_power", "VAP Transmit Power", labelV, nil),
		VAPTxRetries:             prometheus.NewDesc(ns+"vap_transmit_retries_total", "VAP Retries Transmitted", labelV, nil),
		VAPTxCombinedRetries:     prometheus.NewDesc(ns+"vap_transmit_retries_combined_total", "VAP Retries Combined Transmitted", labelV, nil),
		VAPTxDataMpduBytes:       prometheus.NewDesc(ns+"vap_data_mpdu_transmit_bytes_total", "VAP Data MPDU Bytes Transmitted", labelV, nil),
		VAPTxRtsRetries:          prometheus.NewDesc(ns+"vap_transmit_rts_retries_total", "VAP RTS Retries Transmitted", labelV, nil),
		VAPTxSuccess:             prometheus.NewDesc(ns+"vap_transmit_success_total", "VAP Success Transmits", labelV, nil),
		VAPTxTotal:               prometheus.NewDesc(ns+"vap_transmit_total", "VAP Transmit Total", labelV, nil),
		VAPTxGoodbytes:           prometheus.NewDesc(ns+"vap_transmit_goodbyes", "VAP Goodbyes Transmitted", labelV, nil),
		VAPTxLatAvg:              prometheus.NewDesc(ns+"vap_transmit_latency_average_seconds", "VAP Latency Average Transmit", labelV, nil),
		VAPTxLatMax:              prometheus.NewDesc(ns+"vap_transmit_latency_maximum_seconds", "VAP Latency Maximum Transmit", labelV, nil),
		VAPTxLatMin:              prometheus.NewDesc(ns+"vap_transmit_latency_minimum_seconds", "VAP Latency Minimum Transmit", labelV, nil),
		VAPRxGoodbytes:           prometheus.NewDesc(ns+"vap_receive_goodbyes", "VAP Goodbyes Received", labelV, nil),
		VAPRxLatAvg:              prometheus.NewDesc(ns+"vap_receive_latency_average_seconds", "VAP Latency Average Receive", labelV, nil),
		VAPRxLatMax:              prometheus.NewDesc(ns+"vap_receive_latency_maximum_seconds", "VAP Latency Maximum Receive", labelV, nil),
		VAPRxLatMin:              prometheus.NewDesc(ns+"vap_receive_latency_minimum_seconds", "VAP Latency Minimum Receive", labelV, nil),
		VAPWifiTxLatencyMovAvg:   prometheus.NewDesc(ns+"vap_transmit_latency_moving_avg_seconds", "VAP Latency Moving Average Tramsit", labelV, nil),
		VAPWifiTxLatencyMovMax:   prometheus.NewDesc(ns+"vap_transmit_latency_moving_max_seconds", "VAP Latency Moving Maximum Tramsit", labelV, nil),
		VAPWifiTxLatencyMovMin:   prometheus.NewDesc(ns+"vap_transmit_latency_moving_min_seconds", "VAP Latency Moving Minimum Tramsit", labelV, nil),
		VAPWifiTxLatencyMovTotal: prometheus.NewDesc(ns+"vap_transmit_latency_moving_total", "VAP Latency Moving Total Tramsit", labelV, nil),
		VAPWifiTxLatencyMovCount: prometheus.NewDesc(ns+"vap_transmit_latency_moving_count", "VAP Latency Moving Count Tramsit", labelV, nil),
		// N each - 1 per Radio. 1-4 radios per AP usually
		RadioCurrentAntennaGain: prometheus.NewDesc(ns+"radio_current_antenna_gain", "Radio Current Antenna Gain", labelR, nil),
		RadioHt:                 prometheus.NewDesc(ns+"radio_ht", "Radio HT", labelR, nil),
		RadioMaxTxpower:         prometheus.NewDesc(ns+"radio_max_transmit_power", "Radio Maximum Transmit Power", labelR, nil),
		RadioMinTxpower:         prometheus.NewDesc(ns+"radio_min_transmit_power", "Radio Minimum Transmit Power", labelR, nil),
		RadioNss:                prometheus.NewDesc(ns+"radio_nss", "Radio Nss", labelR, nil),
		RadioRadioCaps:          prometheus.NewDesc(ns+"radio_caps", "Radio Capabilities", labelR, nil),
		RadioTxPower:            prometheus.NewDesc(ns+"radio_transmit_power", "Radio Transmit Power", labelR, nil),
		RadioAstBeXmit:          prometheus.NewDesc(ns+"radio_ast_be_xmit", "Radio AstBe Transmit", labelR, nil),
		RadioChannel:            prometheus.NewDesc(ns+"radio_channel", "Radio Channel", labelR, nil),
		RadioCuSelfRx:           prometheus.NewDesc(ns+"radio_channel_utilization_receive_ratio", "Radio Channel Utilization Receive", labelR, nil),
		RadioCuSelfTx:           prometheus.NewDesc(ns+"radio_channel_utilization_transmit_ratio", "Radio Channel Utilization Transmit", labelR, nil),
		RadioExtchannel:         prometheus.NewDesc(ns+"radio_ext_channel", "Radio Ext Channel", labelR, nil),
		RadioGain:               prometheus.NewDesc(ns+"radio_gain", "Radio Gain", labelR, nil),
		RadioGuestNumSta:        prometheus.NewDesc(ns+"radio_guest_stations", "Radio Guest Station Count", labelR, nil),
		RadioNumSta:             prometheus.NewDesc(ns+"radio_stations", "Radio Total Station Count", labelR, nil),
		RadioUserNumSta:         prometheus.NewDesc(ns+"radio_user_stations", "Radio User Station Count", labelR, nil),
		RadioTxPackets:          prometheus.NewDesc(ns+"radio_transmit_packets", "Radio Transmitted Packets", labelR, nil),
		RadioTxRetries:          prometheus.NewDesc(ns+"radio_transmit_retries", "Radio Transmit Retries", labelR, nil),
	}
}

func (u *promUnifi) exportUAP(r report, d *unifi.UAP) {
	labels := []string{d.IP, d.Version, d.Model, d.Serial, d.Type, d.Mac, d.SiteName, d.Name}
	// Wireless System Data.
	r.send([]*metric{
		{u.Device.Uptime, prometheus.GaugeValue, d.Uptime, labels},
		{u.Device.TotalTxBytes, prometheus.CounterValue, d.TxBytes, labels},
		{u.Device.TotalRxBytes, prometheus.CounterValue, d.RxBytes, labels},
		{u.Device.TotalBytes, prometheus.CounterValue, d.Bytes, labels},
		{u.Device.BytesD, prometheus.CounterValue, d.BytesD, labels},     // not sure if these 3 Ds are counters or gauges.
		{u.Device.TxBytesD, prometheus.CounterValue, d.TxBytesD, labels}, // not sure if these 3 Ds are counters or gauges.
		{u.Device.RxBytesD, prometheus.CounterValue, d.RxBytesD, labels}, // not sure if these 3 Ds are counters or gauges.
		{u.Device.BytesR, prometheus.GaugeValue, d.BytesR, labels},
		{u.Device.NumSta, prometheus.GaugeValue, d.NumSta, labels},
		{u.Device.UserNumSta, prometheus.GaugeValue, d.UserNumSta, labels},
		{u.Device.GuestNumSta, prometheus.GaugeValue, d.GuestNumSta, labels},
		{u.Device.Loadavg1, prometheus.GaugeValue, d.SysStats.Loadavg1, labels},
		{u.Device.Loadavg5, prometheus.GaugeValue, d.SysStats.Loadavg5, labels},
		{u.Device.Loadavg15, prometheus.GaugeValue, d.SysStats.Loadavg15, labels},
		{u.Device.MemUsed, prometheus.GaugeValue, d.SysStats.MemUsed, labels},
		{u.Device.MemTotal, prometheus.GaugeValue, d.SysStats.MemTotal, labels},
		{u.Device.MemBuffer, prometheus.GaugeValue, d.SysStats.MemBuffer, labels},
		{u.Device.CPU, prometheus.GaugeValue, d.SystemStats.CPU, labels},
		{u.Device.Mem, prometheus.GaugeValue, d.SystemStats.Mem, labels},
	})

	u.exportUAPstats(r, labels, d.Stat.Ap)
	u.exportVAPtable(r, labels, d.VapTable)
	u.exportRadtable(r, labels, d.RadioTable, d.RadioTableStats)
}

func (u *promUnifi) exportUAPstats(r report, labels []string, ap *unifi.Ap) {
	//	labelA := append([]string{"all"}, labels[2:]...)
	labelU := append([]string{"user"}, labels[6:]...)
	labelG := append([]string{"guest"}, labels[6:]...)
	r.send([]*metric{
		/* // all
		{u.UAP.ApWifiTxDropped, prometheus.CounterValue, ap.WifiTxDropped, labelA},
		{u.UAP.ApRxErrors, prometheus.CounterValue, ap.RxErrors, labelA},
		{u.UAP.ApRxDropped, prometheus.CounterValue, ap.RxDropped, labelA},
		{u.UAP.ApRxFrags, prometheus.CounterValue, ap.RxFrags, labelA},
		{u.UAP.ApRxCrypts, prometheus.CounterValue, ap.RxCrypts, labelA},
		{u.UAP.ApTxPackets, prometheus.CounterValue, ap.TxPackets, labelA},
		{u.UAP.ApTxBytes, prometheus.CounterValue, ap.TxBytes, labelA},
		{u.UAP.ApTxErrors, prometheus.CounterValue, ap.TxErrors, labelA},
		{u.UAP.ApTxDropped, prometheus.CounterValue, ap.TxDropped, labelA},
		{u.UAP.ApTxRetries, prometheus.CounterValue, ap.TxRetries, labelA},
		{u.UAP.ApRxPackets, prometheus.CounterValue, ap.RxPackets, labelA},
		{u.UAP.ApRxBytes, prometheus.CounterValue, ap.RxBytes, labelA},
		{u.UAP.WifiTxAttempts, prometheus.CounterValue, ap.WifiTxAttempts, labelA},
		{u.UAP.MacFilterRejections, prometheus.CounterValue, ap.MacFilterRejections, labelA},
		*/
		// user
		{u.UAP.ApWifiTxDropped, prometheus.CounterValue, ap.UserWifiTxDropped, labelU},
		{u.UAP.ApRxErrors, prometheus.CounterValue, ap.UserRxErrors, labelU},
		{u.UAP.ApRxDropped, prometheus.CounterValue, ap.UserRxDropped, labelU},
		{u.UAP.ApRxFrags, prometheus.CounterValue, ap.UserRxFrags, labelU},
		{u.UAP.ApRxCrypts, prometheus.CounterValue, ap.UserRxCrypts, labelU},
		{u.UAP.ApTxPackets, prometheus.CounterValue, ap.UserTxPackets, labelU},
		{u.UAP.ApTxBytes, prometheus.CounterValue, ap.UserTxBytes, labelU},
		{u.UAP.ApTxErrors, prometheus.CounterValue, ap.UserTxErrors, labelU},
		{u.UAP.ApTxDropped, prometheus.CounterValue, ap.UserTxDropped, labelU},
		{u.UAP.ApTxRetries, prometheus.CounterValue, ap.UserTxRetries, labelU},
		{u.UAP.ApRxPackets, prometheus.CounterValue, ap.UserRxPackets, labelU},
		{u.UAP.ApRxBytes, prometheus.CounterValue, ap.UserRxBytes, labelU},
		{u.UAP.WifiTxAttempts, prometheus.CounterValue, ap.UserWifiTxAttempts, labelU},
		{u.UAP.MacFilterRejections, prometheus.CounterValue, ap.UserMacFilterRejections, labelU},
		// guest
		{u.UAP.ApWifiTxDropped, prometheus.CounterValue, ap.GuestWifiTxDropped, labelG},
		{u.UAP.ApRxErrors, prometheus.CounterValue, ap.GuestRxErrors, labelG},
		{u.UAP.ApRxDropped, prometheus.CounterValue, ap.GuestRxDropped, labelG},
		{u.UAP.ApRxFrags, prometheus.CounterValue, ap.GuestRxFrags, labelG},
		{u.UAP.ApRxCrypts, prometheus.CounterValue, ap.GuestRxCrypts, labelG},
		{u.UAP.ApTxPackets, prometheus.CounterValue, ap.GuestTxPackets, labelG},
		{u.UAP.ApTxBytes, prometheus.CounterValue, ap.GuestTxBytes, labelG},
		{u.UAP.ApTxErrors, prometheus.CounterValue, ap.GuestTxErrors, labelG},
		{u.UAP.ApTxDropped, prometheus.CounterValue, ap.GuestTxDropped, labelG},
		{u.UAP.ApTxRetries, prometheus.CounterValue, ap.GuestTxRetries, labelG},
		{u.UAP.ApRxPackets, prometheus.CounterValue, ap.GuestRxPackets, labelG},
		{u.UAP.ApRxBytes, prometheus.CounterValue, ap.GuestRxBytes, labelG},
		{u.UAP.WifiTxAttempts, prometheus.CounterValue, ap.GuestWifiTxAttempts, labelG},
		{u.UAP.MacFilterRejections, prometheus.CounterValue, ap.GuestMacFilterRejections, labelG},
	})
}

func (u *promUnifi) exportVAPtable(r report, labels []string, vt unifi.VapTable) {
	// vap table stats
	for _, v := range vt {
		if !v.Up.Val {
			continue
		}
		labelV := append([]string{v.Name, v.Bssid, v.Radio, v.RadioName, v.Essid, v.Usage}, labels[6:]...)

		r.send([]*metric{
			{u.UAP.VAPCcq, prometheus.GaugeValue, float64(v.Ccq) / 1000.0, labelV},
			{u.UAP.VAPMacFilterRejections, prometheus.CounterValue, v.MacFilterRejections, labelV},
			{u.UAP.VAPNumSatisfactionSta, prometheus.GaugeValue, v.NumSatisfactionSta, labelV},
			{u.UAP.VAPAvgClientSignal, prometheus.GaugeValue, v.AvgClientSignal.Val, labelV},
			{u.UAP.VAPSatisfaction, prometheus.GaugeValue, v.Satisfaction.Val / 100.0, labelV},
			{u.UAP.VAPSatisfactionNow, prometheus.GaugeValue, v.SatisfactionNow.Val / 100.0, labelV},
			{u.UAP.VAPDNSAvgLatency, prometheus.GaugeValue, v.DNSAvgLatency.Val / 1000, labelV},
			{u.UAP.VAPRxBytes, prometheus.CounterValue, v.RxBytes, labelV},
			{u.UAP.VAPRxCrypts, prometheus.CounterValue, v.RxCrypts, labelV},
			{u.UAP.VAPRxDropped, prometheus.CounterValue, v.RxDropped, labelV},
			{u.UAP.VAPRxErrors, prometheus.CounterValue, v.RxErrors, labelV},
			{u.UAP.VAPRxFrags, prometheus.CounterValue, v.RxFrags, labelV},
			{u.UAP.VAPRxNwids, prometheus.CounterValue, v.RxNwids, labelV},
			{u.UAP.VAPRxPackets, prometheus.CounterValue, v.RxPackets, labelV},
			{u.UAP.VAPTxBytes, prometheus.CounterValue, v.TxBytes, labelV},
			{u.UAP.VAPTxDropped, prometheus.CounterValue, v.TxDropped, labelV},
			{u.UAP.VAPTxErrors, prometheus.CounterValue, v.TxErrors, labelV},
			{u.UAP.VAPTxPackets, prometheus.CounterValue, v.TxPackets, labelV},
			{u.UAP.VAPTxPower, prometheus.GaugeValue, v.TxPower, labelV},
			{u.UAP.VAPTxRetries, prometheus.CounterValue, v.TxRetries, labelV},
			{u.UAP.VAPTxCombinedRetries, prometheus.CounterValue, v.TxCombinedRetries, labelV},
			{u.UAP.VAPTxDataMpduBytes, prometheus.CounterValue, v.TxDataMpduBytes, labelV},
			{u.UAP.VAPTxRtsRetries, prometheus.CounterValue, v.TxRtsRetries, labelV},
			{u.UAP.VAPTxTotal, prometheus.CounterValue, v.TxTotal, labelV},
			{u.UAP.VAPTxGoodbytes, prometheus.CounterValue, v.TxTCPStats.Goodbytes, labelV},
			{u.UAP.VAPTxLatAvg, prometheus.GaugeValue, v.TxTCPStats.LatAvg.Val / 1000, labelV},
			{u.UAP.VAPTxLatMax, prometheus.GaugeValue, v.TxTCPStats.LatMax.Val / 1000, labelV},
			{u.UAP.VAPTxLatMin, prometheus.GaugeValue, v.TxTCPStats.LatMin.Val / 1000, labelV},
			{u.UAP.VAPRxGoodbytes, prometheus.CounterValue, v.RxTCPStats.Goodbytes, labelV},
			{u.UAP.VAPRxLatAvg, prometheus.GaugeValue, v.RxTCPStats.LatAvg.Val / 1000, labelV},
			{u.UAP.VAPRxLatMax, prometheus.GaugeValue, v.RxTCPStats.LatMax.Val / 1000, labelV},
			{u.UAP.VAPRxLatMin, prometheus.GaugeValue, v.RxTCPStats.LatMin.Val / 1000, labelV},
			{u.UAP.VAPWifiTxLatencyMovAvg, prometheus.GaugeValue, v.WifiTxLatencyMov.Avg.Val / 1000, labelV},
			{u.UAP.VAPWifiTxLatencyMovMax, prometheus.GaugeValue, v.WifiTxLatencyMov.Max.Val / 1000, labelV},
			{u.UAP.VAPWifiTxLatencyMovMin, prometheus.GaugeValue, v.WifiTxLatencyMov.Min.Val / 1000, labelV},
			{u.UAP.VAPWifiTxLatencyMovTotal, prometheus.CounterValue, v.WifiTxLatencyMov.Total, labelV},      // not sure if gauge or counter.
			{u.UAP.VAPWifiTxLatencyMovCount, prometheus.CounterValue, v.WifiTxLatencyMov.TotalCount, labelV}, // not sure if gauge or counter.
		})
	}
}

func (u *promUnifi) exportRadtable(r report, labels []string, rt unifi.RadioTable, rts unifi.RadioTableStats) {
	// radio table
	for _, p := range rt {
		labelR := append([]string{p.Name, p.Radio}, labels[6:]...)
		r.send([]*metric{
			{u.UAP.RadioCurrentAntennaGain, prometheus.GaugeValue, p.CurrentAntennaGain, labelR},
			{u.UAP.RadioHt, prometheus.GaugeValue, p.Ht, labelR},
			{u.UAP.RadioMaxTxpower, prometheus.GaugeValue, p.MaxTxpower, labelR},
			{u.UAP.RadioMinTxpower, prometheus.GaugeValue, p.MinTxpower, labelR},
			{u.UAP.RadioNss, prometheus.GaugeValue, p.Nss, labelR},
			{u.UAP.RadioRadioCaps, prometheus.GaugeValue, p.RadioCaps, labelR},
		})

		// combine radio table with radio stats table.
		for _, t := range rts {
			if t.Name != p.Name {
				continue
			}
			r.send([]*metric{
				{u.UAP.RadioTxPower, prometheus.GaugeValue, t.TxPower, labelR},
				{u.UAP.RadioAstBeXmit, prometheus.GaugeValue, t.AstBeXmit, labelR},
				{u.UAP.RadioChannel, prometheus.GaugeValue, t.Channel, labelR},
				{u.UAP.RadioCuSelfRx, prometheus.GaugeValue, t.CuSelfRx.Val / 100.0, labelR},
				{u.UAP.RadioCuSelfTx, prometheus.GaugeValue, t.CuSelfTx.Val / 100.0, labelR},
				{u.UAP.RadioExtchannel, prometheus.GaugeValue, t.Extchannel, labelR},
				{u.UAP.RadioGain, prometheus.GaugeValue, t.Gain, labelR},
				{u.UAP.RadioGuestNumSta, prometheus.GaugeValue, t.GuestNumSta, labelR},
				{u.UAP.RadioNumSta, prometheus.GaugeValue, t.NumSta, labelR},
				{u.UAP.RadioUserNumSta, prometheus.GaugeValue, t.UserNumSta, labelR},
				{u.UAP.RadioTxPackets, prometheus.GaugeValue, t.TxPackets, labelR},
				{u.UAP.RadioTxRetries, prometheus.GaugeValue, t.TxRetries, labelR},
			})
		}
	}
}
