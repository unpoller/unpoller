package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

type uap struct {
	Uptime       *prometheus.Desc
	TotalTxBytes *prometheus.Desc
	TotalRxBytes *prometheus.Desc
	TotalBytes   *prometheus.Desc
	BytesD       *prometheus.Desc
	TxBytesD     *prometheus.Desc
	RxBytesD     *prometheus.Desc
	BytesR       *prometheus.Desc
	NumSta       *prometheus.Desc
	UserNumSta   *prometheus.Desc
	GuestNumSta  *prometheus.Desc
	// System Stats
	Loadavg1  *prometheus.Desc
	Loadavg5  *prometheus.Desc
	Loadavg15 *prometheus.Desc
	MemBuffer *prometheus.Desc
	MemTotal  *prometheus.Desc
	MemUsed   *prometheus.Desc
	CPU       *prometheus.Desc
	Mem       *prometheus.Desc
	ApBytes   *prometheus.Desc
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
	RadioCuTotal            *prometheus.Desc
	RadioExtchannel         *prometheus.Desc
	RadioGain               *prometheus.Desc
	RadioGuestNumSta        *prometheus.Desc
	RadioNumSta             *prometheus.Desc
	RadioUserNumSta         *prometheus.Desc
	RadioTxPackets          *prometheus.Desc
	RadioTxRetries          *prometheus.Desc
}

func descUAP(ns string) *uap {
	if ns += "_uap_"; ns == "_uap_" {
		ns = "uap_"
	}
	labels := []string{"ip", "site_name", "mac", "model", "name", "serial", "site_id",
		"type", "version", "device_id"}
	labelA := append([]string{"stat"}, labels[2:]...)
	labelV := append([]string{"vap_name", "bssid", "radio_name", "essid"}, labels[2:]...)
	labelR := append([]string{"radio_name", "radio", "wlan_group_id"}, labels[2:]...)

	return &uap{
		Uptime:       prometheus.NewDesc(ns+"uptime", "Uptime", labels, nil),
		TotalTxBytes: prometheus.NewDesc(ns+"transmit_bytes_total", "Total Transmitted Bytes", labels, nil),
		TotalRxBytes: prometheus.NewDesc(ns+"receive_bytes_total", "Total Received Bytes", labels, nil),
		TotalBytes:   prometheus.NewDesc(ns+"bytes_total", "Total Bytes Transferred", labels, nil),
		BytesD:       prometheus.NewDesc(ns+"d_bytes", "Total Bytes D???", labels, nil),
		TxBytesD:     prometheus.NewDesc(ns+"d_tranmsit_bytes", "Transmit Bytes D???", labels, nil),
		RxBytesD:     prometheus.NewDesc(ns+"d_receive_bytes", "Receive Bytes D???", labels, nil),
		BytesR:       prometheus.NewDesc(ns+"rate_bytes", "Transfer Rate", labels, nil),
		NumSta:       prometheus.NewDesc(ns+"num_stations_total", "Number of Stations", labels, nil),
		UserNumSta:   prometheus.NewDesc(ns+"num_stations_user_total", "Number of User Stations", labels, nil),
		GuestNumSta:  prometheus.NewDesc(ns+"num_stations_guest_total", "Number of Guest Stations", labels, nil),
		Loadavg1:     prometheus.NewDesc(ns+"load_average_1", "System Load Average 1 Minute", labels, nil),
		Loadavg5:     prometheus.NewDesc(ns+"load_average_5", "System Load Average 5 Minutes", labels, nil),
		Loadavg15:    prometheus.NewDesc(ns+"load_average_15", "System Load Average 15 Minutes", labels, nil),
		MemUsed:      prometheus.NewDesc(ns+"memory_used_bytes", "System Memory Used", labels, nil),
		MemTotal:     prometheus.NewDesc(ns+"memory_installed_bytes", "System Installed Memory", labels, nil),
		MemBuffer:    prometheus.NewDesc(ns+"memory_buffer_bytes", "System Memory Buffer", labels, nil),
		CPU:          prometheus.NewDesc(ns+"cpu_utilization_percent", "System CPU % Utilized", labels, nil),
		Mem:          prometheus.NewDesc(ns+"memory_utilization_percent", "System Memory % Utilized", labels, nil),
		ApBytes:      prometheus.NewDesc(ns+"transferred_bytes_total", "Total Bytes Moved", labels, nil),

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
		VAPCcq:                   prometheus.NewDesc(ns+"vap_ccq", "VAP Client Connection Quality", labelV, nil),
		VAPMacFilterRejections:   prometheus.NewDesc(ns+"vap_mac_filter_rejects_total", "VAP MAC Filter Rejections", labelV, nil),
		VAPNumSatisfactionSta:    prometheus.NewDesc(ns+"vap_num_satisfaction_stations", "VAP Number Satisifaction Stations", labelV, nil),
		VAPAvgClientSignal:       prometheus.NewDesc(ns+"vap_average_client_signal", "VAP Average Client Signal", labelV, nil),
		VAPSatisfaction:          prometheus.NewDesc(ns+"vap_satisfaction", "VAP Satisfaction", labelV, nil),
		VAPSatisfactionNow:       prometheus.NewDesc(ns+"vap_satisfaction_now", "VAP Satisfaction Now", labelV, nil),
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
		VAPTxRtsRetries:          prometheus.NewDesc(ns+"vap_trnamsit_rts_retries_total", "VAP RTS Retries Transmitted", labelV, nil),
		VAPTxSuccess:             prometheus.NewDesc(ns+"vap_transmit_success_total", "VAP Success Transmits", labelV, nil),
		VAPTxTotal:               prometheus.NewDesc(ns+"vap_transmit_total", "VAP Transmit Total", labelV, nil),
		VAPTxGoodbytes:           prometheus.NewDesc(ns+"vap_transmit_goodbyes", "VAP Goodbyes Transmitted", labelV, nil),
		VAPTxLatAvg:              prometheus.NewDesc(ns+"vap_transmit_latency_average", "VAP Latency Average Transmit", labelV, nil),
		VAPTxLatMax:              prometheus.NewDesc(ns+"vap_transmit_latency_maximum", "VAP Latency Maximum Transmit", labelV, nil),
		VAPTxLatMin:              prometheus.NewDesc(ns+"vap_transmit_latency_minimum", "VAP Latency Minimum Transmit", labelV, nil),
		VAPRxGoodbytes:           prometheus.NewDesc(ns+"vap_receive_goodbyes", "VAP Goodbyes Received", labelV, nil),
		VAPRxLatAvg:              prometheus.NewDesc(ns+"vap_receive_latency_average", "VAP Latency Average Receive", labelV, nil),
		VAPRxLatMax:              prometheus.NewDesc(ns+"vap_receive_latency_maximum", "VAP Latency Maximum Receive", labelV, nil),
		VAPRxLatMin:              prometheus.NewDesc(ns+"vap_receive_latency_minimum", "VAP Latency Minimum Receive", labelV, nil),
		VAPWifiTxLatencyMovAvg:   prometheus.NewDesc(ns+"vap_transmit_latency_moving_avg", "VAP Latency Moving Average Tramsit", labelV, nil),
		VAPWifiTxLatencyMovMax:   prometheus.NewDesc(ns+"vap_transmit_latency_moving_max", "VAP Latency Moving Maximum Tramsit", labelV, nil),
		VAPWifiTxLatencyMovMin:   prometheus.NewDesc(ns+"vap_transmit_latency_moving_min", "VAP Latency Moving Minimum Tramsit", labelV, nil),
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
		RadioCuSelfRx:           prometheus.NewDesc(ns+"radio_channel_utilization_receive", "Radio Channel Utilization Receive", labelR, nil),
		RadioCuSelfTx:           prometheus.NewDesc(ns+"radio_channel_utilization_transmit", "Radio Channel Utilization Transmit", labelR, nil),
		RadioCuTotal:            prometheus.NewDesc(ns+"radio_channel_utilization_total", "Radio Channel Utilization", labelR, nil),
		RadioExtchannel:         prometheus.NewDesc(ns+"radio_ext_channel", "Radio Ext Channel", labelR, nil),
		RadioGain:               prometheus.NewDesc(ns+"radio_gain", "Radio Gain", labelR, nil),
		RadioGuestNumSta:        prometheus.NewDesc(ns+"radio_num_guest_stations_total", "Radio Guest Station Count", labelR, nil),
		RadioNumSta:             prometheus.NewDesc(ns+"radio_num_stations_total", "Radio Total Station Count", labelR, nil),
		RadioUserNumSta:         prometheus.NewDesc(ns+"radio_num_user_stations_total", "Radio User Station Count", labelR, nil),
		RadioTxPackets:          prometheus.NewDesc(ns+"radio_transmit_packets_total", "Radio Transmitted Packets", labelR, nil),
		RadioTxRetries:          prometheus.NewDesc(ns+"radio_transmit_retries_total", "Radio Transmit Retries", labelR, nil),
	}
}

func (u *unifiCollector) exportUAPs(uaps []*unifi.UAP, r *Report) {
	for _, a := range uaps {
		labels := []string{a.IP, a.SiteName, a.Mac, a.Model, a.Name, a.Serial, a.SiteID,
			a.Type, a.Version, a.DeviceID}

		// AP data.
		r.ch <- append(append([]*metricExports{
			{u.UAP.Uptime, prometheus.GaugeValue, a.Uptime, labels},
			{u.UAP.TotalTxBytes, prometheus.CounterValue, a.TxBytes, labels},
			{u.UAP.TotalRxBytes, prometheus.CounterValue, a.RxBytes, labels},
			{u.UAP.TotalBytes, prometheus.CounterValue, a.Bytes, labels},
			{u.UAP.BytesD, prometheus.CounterValue, a.BytesD, labels},     // not sure if these 3 Ds are counters or gauges.
			{u.UAP.TxBytesD, prometheus.CounterValue, a.TxBytesD, labels}, // not sure if these 3 Ds are counters or gauges.
			{u.UAP.RxBytesD, prometheus.CounterValue, a.RxBytesD, labels}, // not sure if these 3 Ds are counters or gauges.
			{u.UAP.BytesR, prometheus.GaugeValue, a.BytesR, labels},
			{u.UAP.NumSta, prometheus.GaugeValue, a.NumSta, labels},
			{u.UAP.UserNumSta, prometheus.GaugeValue, a.UserNumSta, labels},
			{u.UAP.GuestNumSta, prometheus.GaugeValue, a.GuestNumSta, labels},
			{u.UAP.Loadavg1, prometheus.GaugeValue, a.SysStats.Loadavg1, labels},
			{u.UAP.Loadavg5, prometheus.GaugeValue, a.SysStats.Loadavg5, labels},
			{u.UAP.Loadavg15, prometheus.GaugeValue, a.SysStats.Loadavg15, labels},
			{u.UAP.MemUsed, prometheus.GaugeValue, a.SysStats.MemUsed, labels},
			{u.UAP.MemTotal, prometheus.GaugeValue, a.SysStats.MemTotal, labels},
			{u.UAP.MemBuffer, prometheus.GaugeValue, a.SysStats.MemBuffer, labels},
			{u.UAP.CPU, prometheus.GaugeValue, a.SystemStats.CPU, labels},
			{u.UAP.Mem, prometheus.GaugeValue, a.SystemStats.Mem, labels},
		},
			u.exportUAPstat(labels[2:], a.Stat.Ap)...),
			u.exportVAPtable(labels[2:], a.VapTable, a.RadioTable, a.RadioTableStats)...)
	}
}

func (u *unifiCollector) exportUAPstat(labels []string, a *unifi.Ap) []*metricExports {
	labelA := append([]string{"all"}, labels...)
	labelU := append([]string{"user"}, labels...)
	labelG := append([]string{"guest"}, labels...)

	return []*metricExports{
		// all
		{u.UAP.ApWifiTxDropped, prometheus.CounterValue, a.WifiTxDropped, labelA},
		{u.UAP.ApRxErrors, prometheus.CounterValue, a.RxErrors, labelA},
		{u.UAP.ApRxDropped, prometheus.CounterValue, a.RxDropped, labelA},
		{u.UAP.ApRxFrags, prometheus.CounterValue, a.RxFrags, labelA},
		{u.UAP.ApRxCrypts, prometheus.CounterValue, a.RxCrypts, labelA},
		{u.UAP.ApTxPackets, prometheus.CounterValue, a.TxPackets, labelA},
		{u.UAP.ApTxBytes, prometheus.CounterValue, a.TxBytes, labelA},
		{u.UAP.ApTxErrors, prometheus.CounterValue, a.TxErrors, labelA},
		{u.UAP.ApTxDropped, prometheus.CounterValue, a.TxDropped, labelA},
		{u.UAP.ApTxRetries, prometheus.CounterValue, a.TxRetries, labelA},
		{u.UAP.ApRxPackets, prometheus.CounterValue, a.RxPackets, labelA},
		{u.UAP.ApRxBytes, prometheus.CounterValue, a.RxBytes, labelA},
		{u.UAP.WifiTxAttempts, prometheus.CounterValue, a.WifiTxAttempts, labelA},
		{u.UAP.MacFilterRejections, prometheus.CounterValue, a.MacFilterRejections, labelA},
		// user
		{u.UAP.ApWifiTxDropped, prometheus.CounterValue, a.UserWifiTxDropped, labelU},
		{u.UAP.ApRxErrors, prometheus.CounterValue, a.UserRxErrors, labelU},
		{u.UAP.ApRxDropped, prometheus.CounterValue, a.UserRxDropped, labelU},
		{u.UAP.ApRxFrags, prometheus.CounterValue, a.UserRxFrags, labelU},
		{u.UAP.ApRxCrypts, prometheus.CounterValue, a.UserRxCrypts, labelU},
		{u.UAP.ApTxPackets, prometheus.CounterValue, a.UserTxPackets, labelU},
		{u.UAP.ApTxBytes, prometheus.CounterValue, a.UserTxBytes, labelU},
		{u.UAP.ApTxErrors, prometheus.CounterValue, a.UserTxErrors, labelU},
		{u.UAP.ApTxDropped, prometheus.CounterValue, a.UserTxDropped, labelU},
		{u.UAP.ApTxRetries, prometheus.CounterValue, a.UserTxRetries, labelU},
		{u.UAP.ApRxPackets, prometheus.CounterValue, a.UserRxPackets, labelU},
		{u.UAP.ApRxBytes, prometheus.CounterValue, a.UserRxBytes, labelU},
		{u.UAP.WifiTxAttempts, prometheus.CounterValue, a.UserWifiTxAttempts, labelU},
		{u.UAP.MacFilterRejections, prometheus.CounterValue, a.UserMacFilterRejections, labelU},
		// guest
		{u.UAP.ApWifiTxDropped, prometheus.CounterValue, a.GuestWifiTxDropped, labelG},
		{u.UAP.ApRxErrors, prometheus.CounterValue, a.GuestRxErrors, labelG},
		{u.UAP.ApRxDropped, prometheus.CounterValue, a.GuestRxDropped, labelG},
		{u.UAP.ApRxFrags, prometheus.CounterValue, a.GuestRxFrags, labelG},
		{u.UAP.ApRxCrypts, prometheus.CounterValue, a.GuestRxCrypts, labelG},
		{u.UAP.ApTxPackets, prometheus.CounterValue, a.GuestTxPackets, labelG},
		{u.UAP.ApTxBytes, prometheus.CounterValue, a.GuestTxBytes, labelG},
		{u.UAP.ApTxErrors, prometheus.CounterValue, a.GuestTxErrors, labelG},
		{u.UAP.ApTxDropped, prometheus.CounterValue, a.GuestTxDropped, labelG},
		{u.UAP.ApTxRetries, prometheus.CounterValue, a.GuestTxRetries, labelG},
		{u.UAP.ApRxPackets, prometheus.CounterValue, a.GuestRxPackets, labelG},
		{u.UAP.ApRxBytes, prometheus.CounterValue, a.GuestRxBytes, labelG},
		{u.UAP.WifiTxAttempts, prometheus.CounterValue, a.GuestWifiTxAttempts, labelG},
		{u.UAP.MacFilterRejections, prometheus.CounterValue, a.GuestMacFilterRejections, labelG},
	}
}

func (u *unifiCollector) exportVAPtable(labels []string, vt unifi.VapTable, rt unifi.RadioTable, rts unifi.RadioTableStats) []*metricExports {
	metrics := []*metricExports{}

	// vap table stats
	for _, v := range vt {
		labelV := append([]string{v.Name, v.Bssid, v.RadioName, v.Essid}, labels...)
		metrics = append(metrics, []*metricExports{
			{u.UAP.VAPCcq, prometheus.GaugeValue, v.Ccq, labelV},
			{u.UAP.VAPMacFilterRejections, prometheus.CounterValue, v.MacFilterRejections, labelV},
			{u.UAP.VAPNumSatisfactionSta, prometheus.GaugeValue, v.NumSatisfactionSta, labelV},
			{u.UAP.VAPAvgClientSignal, prometheus.GaugeValue, v.AvgClientSignal, labelV},
			{u.UAP.VAPSatisfaction, prometheus.GaugeValue, v.Satisfaction, labelV},
			{u.UAP.VAPSatisfactionNow, prometheus.GaugeValue, v.SatisfactionNow, labelV},
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
			{u.UAP.VAPTxLatAvg, prometheus.GaugeValue, v.TxTCPStats.LatAvg, labelV},
			{u.UAP.VAPTxLatMax, prometheus.GaugeValue, v.TxTCPStats.LatMax, labelV},
			{u.UAP.VAPTxLatMin, prometheus.GaugeValue, v.TxTCPStats.LatMin, labelV},
			{u.UAP.VAPRxGoodbytes, prometheus.CounterValue, v.RxTCPStats.Goodbytes, labelV},
			{u.UAP.VAPRxLatAvg, prometheus.GaugeValue, v.RxTCPStats.LatAvg, labelV},
			{u.UAP.VAPRxLatMax, prometheus.GaugeValue, v.RxTCPStats.LatMax, labelV},
			{u.UAP.VAPRxLatMin, prometheus.GaugeValue, v.RxTCPStats.LatMin, labelV},
			{u.UAP.VAPWifiTxLatencyMovAvg, prometheus.GaugeValue, v.WifiTxLatencyMov.Avg, labelV},
			{u.UAP.VAPWifiTxLatencyMovMax, prometheus.GaugeValue, v.WifiTxLatencyMov.Max, labelV},
			{u.UAP.VAPWifiTxLatencyMovMin, prometheus.GaugeValue, v.WifiTxLatencyMov.Min, labelV},
			{u.UAP.VAPWifiTxLatencyMovTotal, prometheus.CounterValue, v.WifiTxLatencyMov.Total, labelV},      // not sure if gauge or counter.
			{u.UAP.VAPWifiTxLatencyMovCount, prometheus.CounterValue, v.WifiTxLatencyMov.TotalCount, labelV}, // not sure if gauge or counter.
		}...)
	}

	// radio table
	for _, p := range rt {
		labelR := append([]string{p.Name, p.Radio, p.WlangroupID}, labels...)
		metrics = append(metrics, []*metricExports{
			{u.UAP.RadioCurrentAntennaGain, prometheus.GaugeValue, p.CurrentAntennaGain, labelR},
			{u.UAP.RadioHt, prometheus.GaugeValue, p.Ht, labelR},
			{u.UAP.RadioMaxTxpower, prometheus.GaugeValue, p.MaxTxpower, labelR},
			{u.UAP.RadioMinTxpower, prometheus.GaugeValue, p.MinTxpower, labelR},
			{u.UAP.RadioNss, prometheus.GaugeValue, p.Nss, labelR},
			{u.UAP.RadioRadioCaps, prometheus.GaugeValue, p.RadioCaps, labelR},
		}...)

		// combine radio table with radio stats table.
		for _, t := range rts {
			if t.Name != p.Name {
				continue
			}
			metrics = append(metrics, []*metricExports{
				{u.UAP.RadioTxPower, prometheus.GaugeValue, t.TxPower, labelR},
				{u.UAP.RadioAstBeXmit, prometheus.GaugeValue, t.AstBeXmit, labelR},
				{u.UAP.RadioChannel, prometheus.GaugeValue, t.Channel, labelR},
				{u.UAP.RadioCuSelfRx, prometheus.GaugeValue, t.CuSelfRx, labelR},
				{u.UAP.RadioCuSelfTx, prometheus.GaugeValue, t.CuSelfTx, labelR},
				{u.UAP.RadioCuTotal, prometheus.GaugeValue, t.CuTotal, labelR},
				{u.UAP.RadioExtchannel, prometheus.GaugeValue, t.Extchannel, labelR},
				{u.UAP.RadioGain, prometheus.GaugeValue, t.Gain, labelR},
				{u.UAP.RadioGuestNumSta, prometheus.GaugeValue, t.GuestNumSta, labelR},
				{u.UAP.RadioNumSta, prometheus.GaugeValue, t.NumSta, labelR},
				{u.UAP.RadioUserNumSta, prometheus.GaugeValue, t.UserNumSta, labelR},
				{u.UAP.RadioTxPackets, prometheus.CounterValue, t.TxPackets, labelR},
				{u.UAP.RadioTxRetries, prometheus.CounterValue, t.TxRetries, labelR},
			}...)
		}
	}
	return metrics
}
