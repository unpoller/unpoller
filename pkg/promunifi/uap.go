package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
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
	RadioCuTotal            *prometheus.Desc
	RadioExtchannel         *prometheus.Desc
	RadioGain               *prometheus.Desc
	RadioNumSta             *prometheus.Desc
	RadioTxPackets          *prometheus.Desc
	RadioTxRetries          *prometheus.Desc
}

type rogueap struct {
	Age        *prometheus.Desc
	BW         *prometheus.Desc
	CenterFreq *prometheus.Desc
	Channel    *prometheus.Desc
	Freq       *prometheus.Desc
	Noise      *prometheus.Desc
	RSSI       *prometheus.Desc
	RSSIAge    *prometheus.Desc
	Signal     *prometheus.Desc
}

func descRogueAP(ns string) *rogueap {
	label := []string{
		"security", "oui", "band", "mac", "ap_mac", "radio", "radio_name", "site_name", "name", "source",
	}

	return &rogueap{
		Age:        prometheus.NewDesc(ns+"age", "RogueAP Age", label, nil),
		BW:         prometheus.NewDesc(ns+"bw", "RogueAP BW", label, nil),
		CenterFreq: prometheus.NewDesc(ns+"center_freq", "RogueAP Center Frequency", label, nil),
		Channel:    prometheus.NewDesc(ns+"channel", "RogueAP Channel", label, nil),
		Freq:       prometheus.NewDesc(ns+"frequency", "RogueAP Frequency", label, nil),
		Noise:      prometheus.NewDesc(ns+"noise", "RogueAP Noise", label, nil),
		RSSI:       prometheus.NewDesc(ns+"rssi", "RogueAP RSSI", label, nil),
		RSSIAge:    prometheus.NewDesc(ns+"rssi_age", "RogueAP RSSI Age", label, nil),
		Signal:     prometheus.NewDesc(ns+"signal", "RogueAP Signal", label, nil),
	}
}

func descUAP(ns string) *uap { // nolint: funlen
	labelA := []string{"stat", "site_name", "name", "source"} // stat + labels[1:]
	labelV := []string{"vap_name", "bssid", "radio", "radio_name", "essid", "usage", "site_name", "name", "source"}
	labelR := []string{"radio_name", "radio", "site_name", "name", "source"}
	nd := prometheus.NewDesc

	return &uap{
		// 3x each - stat table: total, guest, user
		ApWifiTxDropped:     nd(ns+"stat_wifi_transmt_dropped_total", "Wifi Transmissions Dropped", labelA, nil),
		ApRxErrors:          nd(ns+"stat_receive_errors_total", "Receive Errors", labelA, nil),
		ApRxDropped:         nd(ns+"stat_receive_dropped_total", "Receive Dropped", labelA, nil),
		ApRxFrags:           nd(ns+"stat_receive_frags_total", "Received Frags", labelA, nil),
		ApRxCrypts:          nd(ns+"stat_receive_crypts_total", "Receive Crypts", labelA, nil),
		ApTxPackets:         nd(ns+"stat_transmit_packets_total", "Transmit Packets", labelA, nil),
		ApTxBytes:           nd(ns+"stat_transmit_bytes_total", "Transmit Bytes", labelA, nil),
		ApTxErrors:          nd(ns+"stat_transmit_errors_total", "Transmit Errors", labelA, nil),
		ApTxDropped:         nd(ns+"stat_transmit_dropped_total", "Transmit Dropped", labelA, nil),
		ApTxRetries:         nd(ns+"stat_retries_tx_total", "Transmit Retries", labelA, nil),
		ApRxPackets:         nd(ns+"stat_receive_packets_total", "Receive Packets", labelA, nil),
		ApRxBytes:           nd(ns+"stat_receive_bytes_total", "Receive Bytes", labelA, nil),
		WifiTxAttempts:      nd(ns+"stat_wifi_transmit_attempts_total", "Wifi Transmission Attempts", labelA, nil),
		MacFilterRejections: nd(ns+"stat_mac_filter_rejects_total", "MAC Filter Rejections", labelA, nil),
		// N each - 1 per Virtual AP (VAP)
		VAPCcq:                   nd(ns+"vap_ccq_ratio", "VAP Client Connection Quality", labelV, nil),
		VAPMacFilterRejections:   nd(ns+"vap_mac_filter_rejects_total", "VAP MAC Filter Rejections", labelV, nil),
		VAPNumSatisfactionSta:    nd(ns+"vap_satisfaction_stations", "VAP Number Satisifaction Stations", labelV, nil),
		VAPAvgClientSignal:       nd(ns+"vap_average_client_signal", "VAP Average Client Signal", labelV, nil),
		VAPSatisfaction:          nd(ns+"vap_satisfaction_ratio", "VAP Satisfaction", labelV, nil),
		VAPSatisfactionNow:       nd(ns+"vap_satisfaction_now_ratio", "VAP Satisfaction Now", labelV, nil),
		VAPDNSAvgLatency:         nd(ns+"vap_dns_latency_average_seconds", "VAP DNS Latency Average", labelV, nil),
		VAPRxBytes:               nd(ns+"vap_receive_bytes_total", "VAP Bytes Received", labelV, nil),
		VAPRxCrypts:              nd(ns+"vap_receive_crypts_total", "VAP Crypts Received", labelV, nil),
		VAPRxDropped:             nd(ns+"vap_receive_dropped_total", "VAP Dropped Received", labelV, nil),
		VAPRxErrors:              nd(ns+"vap_receive_errors_total", "VAP Errors Received", labelV, nil),
		VAPRxFrags:               nd(ns+"vap_receive_frags_total", "VAP Frags Received", labelV, nil),
		VAPRxNwids:               nd(ns+"vap_receive_nwids_total", "VAP Nwids Received", labelV, nil),
		VAPRxPackets:             nd(ns+"vap_receive_packets_total", "VAP Packets Received", labelV, nil),
		VAPTxBytes:               nd(ns+"vap_transmit_bytes_total", "VAP Bytes Transmitted", labelV, nil),
		VAPTxDropped:             nd(ns+"vap_transmit_dropped_total", "VAP Dropped Transmitted", labelV, nil),
		VAPTxErrors:              nd(ns+"vap_transmit_errors_total", "VAP Errors Transmitted", labelV, nil),
		VAPTxPackets:             nd(ns+"vap_transmit_packets_total", "VAP Packets Transmitted", labelV, nil),
		VAPTxPower:               nd(ns+"vap_transmit_power", "VAP Transmit Power", labelV, nil),
		VAPTxRetries:             nd(ns+"vap_transmit_retries_total", "VAP Retries Transmitted", labelV, nil),
		VAPTxCombinedRetries:     nd(ns+"vap_transmit_retries_combined_total", "VAP Retries Combined Tx", labelV, nil),
		VAPTxDataMpduBytes:       nd(ns+"vap_data_mpdu_transmit_bytes_total", "VAP Data MPDU Bytes Tx", labelV, nil),
		VAPTxRtsRetries:          nd(ns+"vap_transmit_rts_retries_total", "VAP RTS Retries Transmitted", labelV, nil),
		VAPTxSuccess:             nd(ns+"vap_transmit_success_total", "VAP Success Transmits", labelV, nil),
		VAPTxTotal:               nd(ns+"vap_transmit_total", "VAP Transmit Total", labelV, nil),
		VAPTxGoodbytes:           nd(ns+"vap_transmit_goodbyes", "VAP Goodbyes Transmitted", labelV, nil),
		VAPTxLatAvg:              nd(ns+"vap_transmit_latency_average_seconds", "VAP Latency Average Tx", labelV, nil),
		VAPTxLatMax:              nd(ns+"vap_transmit_latency_maximum_seconds", "VAP Latency Maximum Tx", labelV, nil),
		VAPTxLatMin:              nd(ns+"vap_transmit_latency_minimum_seconds", "VAP Latency Minimum Tx", labelV, nil),
		VAPRxGoodbytes:           nd(ns+"vap_receive_goodbyes", "VAP Goodbyes Received", labelV, nil),
		VAPRxLatAvg:              nd(ns+"vap_receive_latency_average_seconds", "VAP Latency Average Rx", labelV, nil),
		VAPRxLatMax:              nd(ns+"vap_receive_latency_maximum_seconds", "VAP Latency Maximum Rx", labelV, nil),
		VAPRxLatMin:              nd(ns+"vap_receive_latency_minimum_seconds", "VAP Latency Minimum Rx", labelV, nil),
		VAPWifiTxLatencyMovAvg:   nd(ns+"vap_transmit_latency_moving_avg_seconds", "VAP Latency Moving Avg Tx", labelV, nil),
		VAPWifiTxLatencyMovMax:   nd(ns+"vap_transmit_latency_moving_max_seconds", "VAP Latency Moving Min Tx", labelV, nil),
		VAPWifiTxLatencyMovMin:   nd(ns+"vap_transmit_latency_moving_min_seconds", "VAP Latency Moving Max Tx", labelV, nil),
		VAPWifiTxLatencyMovTotal: nd(ns+"vap_transmit_latency_moving_total", "VAP Latency Moving Total Tramsit", labelV, nil),
		VAPWifiTxLatencyMovCount: nd(ns+"vap_transmit_latency_moving_count", "VAP Latency Moving Count Tramsit", labelV, nil),
		// N each - 1 per Radio. 1-4 radios per AP usually
		RadioCurrentAntennaGain: nd(ns+"radio_current_antenna_gain", "Radio Current Antenna Gain", labelR, nil),
		RadioHt:                 nd(ns+"radio_ht", "Radio HT", labelR, nil),
		RadioMaxTxpower:         nd(ns+"radio_max_transmit_power", "Radio Maximum Transmit Power", labelR, nil),
		RadioMinTxpower:         nd(ns+"radio_min_transmit_power", "Radio Minimum Transmit Power", labelR, nil),
		RadioNss:                nd(ns+"radio_nss", "Radio Nss", labelR, nil),
		RadioRadioCaps:          nd(ns+"radio_caps", "Radio Capabilities", labelR, nil),
		RadioTxPower:            nd(ns+"radio_transmit_power", "Radio Transmit Power", labelR, nil),
		RadioAstBeXmit:          nd(ns+"radio_ast_be_xmit", "Radio AstBe Transmit", labelR, nil),
		RadioChannel:            nd(ns+"radio_channel", "Radio Channel", labelR, nil),
		RadioCuSelfRx:           nd(ns+"radio_channel_utilization_receive_ratio", "Channel Utilization Rx", labelR, nil),
		RadioCuSelfTx:           nd(ns+"radio_channel_utilization_transmit_ratio", "Channel Utilization Tx", labelR, nil),
		RadioCuTotal:            nd(ns+"radio_channel_utilization_total_ratio", "Channel Utilization Total", labelR, nil),
		RadioExtchannel:         nd(ns+"radio_ext_channel", "Radio Ext Channel", labelR, nil),
		RadioGain:               nd(ns+"radio_gain", "Radio Gain", labelR, nil),
		RadioNumSta:             nd(ns+"radio_stations", "Radio Total Station Count", append(labelR, "station_type"), nil),
		RadioTxPackets:          nd(ns+"radio_transmit_packets", "Radio Transmitted Packets", labelR, nil),
		RadioTxRetries:          nd(ns+"radio_transmit_retries", "Radio Transmit Retries", labelR, nil),
	}
}

func (u *promUnifi) exportRogueAP(r report, d *unifi.RogueAP) {
	if d.Age.Val == 0 {
		return // only keep things that are recent.
	}

	labels := []string{
		d.Security, d.Oui, d.Band, d.Bssid, d.ApMac, d.Radio, d.RadioName, d.SiteName, d.Essid, d.SourceName,
	}

	r.send([]*metric{
		{u.RogueAP.Age, gauge, d.Age.Val, labels},
		{u.RogueAP.BW, gauge, d.Bw.Val, labels},
		{u.RogueAP.CenterFreq, gauge, d.CenterFreq.Val, labels},
		{u.RogueAP.Channel, gauge, d.Channel, labels},
		{u.RogueAP.Freq, gauge, d.Freq.Val, labels},
		{u.RogueAP.Noise, gauge, d.Noise.Val, labels},
		{u.RogueAP.RSSI, gauge, d.Rssi.Val, labels},
		{u.RogueAP.RSSIAge, gauge, d.RssiAge.Val, labels},
		{u.RogueAP.Signal, gauge, d.Signal.Val, labels},
	})
}

func (u *promUnifi) exportUAP(r report, d *unifi.UAP) {
	if !d.Adopted.Val || d.Locating.Val {
		return
	}

	labels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	infoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}
	u.exportUAPstats(r, labels, d.Stat.Ap, d.BytesD, d.TxBytesD, d.RxBytesD, d.BytesR)
	u.exportVAPtable(r, labels, d.VapTable)
	u.exportPRTtable(r, labels, d.PortTable)
	u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)
	u.exportSYSstats(r, labels, d.SysStats, d.SystemStats)
	u.exportSTAcount(r, labels, d.UserNumSta, d.GuestNumSta)
	u.exportRADtable(r, labels, d.RadioTable, d.RadioTableStats)
	r.send([]*metric{
		{u.Device.Info, gauge, 1.0, append(labels, infoLabels...)},
		{u.Device.Uptime, gauge, d.Uptime, labels},
		{u.Device.Upgradeable, gauge, d.Upgradable.Val, labels},
	})
}

// udm doesn't have these stats exposed yet, so pass 2 or 6 metrics.
func (u *promUnifi) exportUAPstats(r report, labels []string, ap *unifi.Ap, bytes ...unifi.FlexInt) {
	if ap == nil {
		return
	}

	labelU := []string{"user", labels[1], labels[2], labels[3]}
	labelG := []string{"guest", labels[1], labels[2], labels[3]}
	r.send([]*metric{
		// ap only stuff.
		{u.Device.BytesD, counter, bytes[0], labels},   // not sure if these 3 Ds are counters or gauges.
		{u.Device.TxBytesD, counter, bytes[1], labels}, // not sure if these 3 Ds are counters or gauges.
		{u.Device.RxBytesD, counter, bytes[2], labels}, // not sure if these 3 Ds are counters or gauges.
		{u.Device.BytesR, gauge, bytes[3], labels},     // only UAP has this one, and those ^ weird.
		// user
		{u.UAP.ApWifiTxDropped, counter, ap.UserWifiTxDropped, labelU},
		{u.UAP.ApRxErrors, counter, ap.UserRxErrors, labelU},
		{u.UAP.ApRxDropped, counter, ap.UserRxDropped, labelU},
		{u.UAP.ApRxFrags, counter, ap.UserRxFrags, labelU},
		{u.UAP.ApRxCrypts, counter, ap.UserRxCrypts, labelU},
		{u.UAP.ApTxPackets, counter, ap.UserTxPackets, labelU},
		{u.UAP.ApTxBytes, counter, ap.UserTxBytes, labelU},
		{u.UAP.ApTxErrors, counter, ap.UserTxErrors, labelU},
		{u.UAP.ApTxDropped, counter, ap.UserTxDropped, labelU},
		{u.UAP.ApTxRetries, counter, ap.UserTxRetries, labelU},
		{u.UAP.ApRxPackets, counter, ap.UserRxPackets, labelU},
		{u.UAP.ApRxBytes, counter, ap.UserRxBytes, labelU},
		{u.UAP.WifiTxAttempts, counter, ap.UserWifiTxAttempts, labelU},
		{u.UAP.MacFilterRejections, counter, ap.UserMacFilterRejections, labelU},
		// guest
		{u.UAP.ApWifiTxDropped, counter, ap.GuestWifiTxDropped, labelG},
		{u.UAP.ApRxErrors, counter, ap.GuestRxErrors, labelG},
		{u.UAP.ApRxDropped, counter, ap.GuestRxDropped, labelG},
		{u.UAP.ApRxFrags, counter, ap.GuestRxFrags, labelG},
		{u.UAP.ApRxCrypts, counter, ap.GuestRxCrypts, labelG},
		{u.UAP.ApTxPackets, counter, ap.GuestTxPackets, labelG},
		{u.UAP.ApTxBytes, counter, ap.GuestTxBytes, labelG},
		{u.UAP.ApTxErrors, counter, ap.GuestTxErrors, labelG},
		{u.UAP.ApTxDropped, counter, ap.GuestTxDropped, labelG},
		{u.UAP.ApTxRetries, counter, ap.GuestTxRetries, labelG},
		{u.UAP.ApRxPackets, counter, ap.GuestRxPackets, labelG},
		{u.UAP.ApRxBytes, counter, ap.GuestRxBytes, labelG},
		{u.UAP.WifiTxAttempts, counter, ap.GuestWifiTxAttempts, labelG},
		{u.UAP.MacFilterRejections, counter, ap.GuestMacFilterRejections, labelG},
	})
}

// UAP VAP Table.
func (u *promUnifi) exportVAPtable(r report, labels []string, vt unifi.VapTable) {
	// vap table stats
	for _, v := range vt {
		if !v.Up.Val {
			continue
		}

		labelV := []string{v.Name, v.Bssid, v.Radio, v.RadioName, v.Essid, v.Usage, labels[1], labels[2], labels[3]}
		r.send([]*metric{
			{u.UAP.VAPCcq, gauge, float64(v.Ccq) / 1000.0, labelV},
			{u.UAP.VAPMacFilterRejections, counter, v.MacFilterRejections, labelV},
			{u.UAP.VAPNumSatisfactionSta, gauge, v.NumSatisfactionSta, labelV},
			{u.UAP.VAPAvgClientSignal, gauge, v.AvgClientSignal.Val, labelV},
			{u.UAP.VAPSatisfaction, gauge, v.Satisfaction.Val / 100.0, labelV},
			{u.UAP.VAPSatisfactionNow, gauge, v.SatisfactionNow.Val / 100.0, labelV},
			{u.UAP.VAPDNSAvgLatency, gauge, v.DNSAvgLatency.Val / 1000, labelV},
			{u.UAP.VAPRxBytes, counter, v.RxBytes, labelV},
			{u.UAP.VAPRxCrypts, counter, v.RxCrypts, labelV},
			{u.UAP.VAPRxDropped, counter, v.RxDropped, labelV},
			{u.UAP.VAPRxErrors, counter, v.RxErrors, labelV},
			{u.UAP.VAPRxFrags, counter, v.RxFrags, labelV},
			{u.UAP.VAPRxNwids, counter, v.RxNwids, labelV},
			{u.UAP.VAPRxPackets, counter, v.RxPackets, labelV},
			{u.UAP.VAPTxBytes, counter, v.TxBytes, labelV},
			{u.UAP.VAPTxDropped, counter, v.TxDropped, labelV},
			{u.UAP.VAPTxErrors, counter, v.TxErrors, labelV},
			{u.UAP.VAPTxPackets, counter, v.TxPackets, labelV},
			{u.UAP.VAPTxPower, gauge, v.TxPower, labelV},
			{u.UAP.VAPTxRetries, counter, v.TxRetries, labelV},
			{u.UAP.VAPTxCombinedRetries, counter, v.TxCombinedRetries, labelV},
			{u.UAP.VAPTxDataMpduBytes, counter, v.TxDataMpduBytes, labelV},
			{u.UAP.VAPTxRtsRetries, counter, v.TxRtsRetries, labelV},
			{u.UAP.VAPTxTotal, counter, v.TxTotal, labelV},
			{u.UAP.VAPTxGoodbytes, counter, v.TxTCPStats.Goodbytes, labelV},
			{u.UAP.VAPTxLatAvg, gauge, v.TxTCPStats.LatAvg.Val / 1000, labelV},
			{u.UAP.VAPTxLatMax, gauge, v.TxTCPStats.LatMax.Val / 1000, labelV},
			{u.UAP.VAPTxLatMin, gauge, v.TxTCPStats.LatMin.Val / 1000, labelV},
			{u.UAP.VAPRxGoodbytes, counter, v.RxTCPStats.Goodbytes, labelV},
			{u.UAP.VAPRxLatAvg, gauge, v.RxTCPStats.LatAvg.Val / 1000, labelV},
			{u.UAP.VAPRxLatMax, gauge, v.RxTCPStats.LatMax.Val / 1000, labelV},
			{u.UAP.VAPRxLatMin, gauge, v.RxTCPStats.LatMin.Val / 1000, labelV},
			{u.UAP.VAPWifiTxLatencyMovAvg, gauge, v.WifiTxLatencyMov.Avg.Val / 1000, labelV},
			{u.UAP.VAPWifiTxLatencyMovMax, gauge, v.WifiTxLatencyMov.Max.Val / 1000, labelV},
			{u.UAP.VAPWifiTxLatencyMovMin, gauge, v.WifiTxLatencyMov.Min.Val / 1000, labelV},
			{u.UAP.VAPWifiTxLatencyMovTotal, counter, v.WifiTxLatencyMov.Total, labelV},      // not sure if gauge or counter.
			{u.UAP.VAPWifiTxLatencyMovCount, counter, v.WifiTxLatencyMov.TotalCount, labelV}, // not sure if gauge or counter.
		})
	}
}

// UAP Radio Table.
func (u *promUnifi) exportRADtable(r report, labels []string, rt unifi.RadioTable, rts unifi.RadioTableStats) {
	// radio table
	for _, p := range rt {
		labelR := []string{p.Name, p.Radio, labels[1], labels[2], labels[3]}
		labelRUser := append(labelR, "user")
		labelRGuest := append(labelR, "guest")

		r.send([]*metric{
			{u.UAP.RadioCurrentAntennaGain, gauge, p.CurrentAntennaGain, labelR},
			{u.UAP.RadioHt, gauge, p.Ht, labelR},
			{u.UAP.RadioMaxTxpower, gauge, p.MaxTxpower, labelR},
			{u.UAP.RadioMinTxpower, gauge, p.MinTxpower, labelR},
			{u.UAP.RadioNss, gauge, p.Nss, labelR},
			{u.UAP.RadioRadioCaps, gauge, p.RadioCaps, labelR},
		})

		// combine radio table with radio stats table.
		for _, t := range rts {
			if t.Name != p.Name {
				continue
			}

			r.send([]*metric{
				{u.UAP.RadioTxPower, gauge, t.TxPower, labelR},
				{u.UAP.RadioAstBeXmit, gauge, t.AstBeXmit, labelR},
				{u.UAP.RadioChannel, gauge, t.Channel, labelR},
				{u.UAP.RadioCuSelfRx, gauge, t.CuSelfRx.Val / 100.0, labelR},
				{u.UAP.RadioCuSelfTx, gauge, t.CuSelfTx.Val / 100.0, labelR},
				{u.UAP.RadioCuTotal, gauge, t.CuTotal.Val / 100.0, labelR},
				{u.UAP.RadioExtchannel, gauge, t.Extchannel, labelR},
				{u.UAP.RadioGain, gauge, t.Gain, labelR},
				{u.UAP.RadioNumSta, gauge, t.GuestNumSta, labelRGuest},
				{u.UAP.RadioNumSta, gauge, t.UserNumSta, labelRUser},
				{u.UAP.RadioTxPackets, gauge, t.TxPackets, labelR},
				{u.UAP.RadioTxRetries, gauge, t.TxRetries, labelR},
			})

			break
		}
	}
}
