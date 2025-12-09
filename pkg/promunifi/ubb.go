package promunifi

import (
	"github.com/unpoller/unifi/v5"
)

// exportUBB is a collection of stats from UBB (UniFi Building Bridge).
// UBB devices are point-to-point wireless bridges with dual radios:
//   - wifi0: 5GHz radio (802.11ac)
//   - terra2/wlan0/ad: 60GHz radio (802.11ad - Terragraph/WiGig)
func (u *promUnifi) exportUBB(r report, d *unifi.UBB) {
	if !d.Adopted.Val || d.Locating.Val {
		return
	}

	labels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	infoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}

	// Export UBB-specific stats if available
	u.exportUBBstats(r, labels, d.Stat)

	// Export VAP table (Virtual Access Point table - wireless interface stats)
	u.exportVAPtable(r, labels, d.VapTable)

	// Export Radio tables (includes 5GHz wifi0 and 60GHz terra2/ad radios)
	u.exportRADtable(r, labels, d.RadioTable, d.RadioTableStats)

	// Shared device stats
	u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)

	if d.SysStats != nil && d.SystemStats != nil {
		u.exportSYSstats(r, labels, *d.SysStats, *d.SystemStats)
	}

	// Device info, uptime, and temperature
	r.send([]*metric{
		{u.Device.Info, gauge, 1.0, append(labels, infoLabels...)},
		{u.Device.Uptime, gauge, d.Uptime, labels},
		{u.Device.Temperature, gauge, d.GeneralTemperature.Val, append(labels, d.Name, "general")},
	})

	// UBB-specific metrics
	if d.P2PStats != nil {
		u.exportP2Pstats(r, labels, d.P2PStats)
	}

	// Link quality metrics for point-to-point links
	r.send([]*metric{
		{u.Device.Counter, gauge, d.LinkQuality.Val, append(labels, "link_quality")},
		{u.Device.Counter, gauge, d.LinkQualityCurrent.Val, append(labels, "link_quality_current")},
		{u.Device.Counter, gauge, d.LinkCapacity.Val, append(labels, "link_capacity")},
	})
}

// exportUBBstats exports UBB-specific stats from the Bb structure.
// This includes separate metrics for wifi0 (5GHz) and terra2 (60GHz) radios.
func (u *promUnifi) exportUBBstats(r report, labels []string, stat *unifi.UBBStat) {
	if stat == nil || stat.Bb == nil {
		return
	}

	bb := stat.Bb

	// Export aggregated stats (total across both radios)
	labelTotal := []string{"total", labels[1], labels[2], labels[3]}
	r.send([]*metric{
		{u.UAP.ApRxPackets, counter, bb.RxPackets, labelTotal},
		{u.UAP.ApRxBytes, counter, bb.RxBytes, labelTotal},
		{u.UAP.ApRxErrors, counter, bb.RxErrors, labelTotal},
		{u.UAP.ApRxDropped, counter, bb.RxDropped, labelTotal},
		{u.UAP.ApRxCrypts, counter, bb.RxCrypts, labelTotal},
		{u.UAP.ApRxFrags, counter, bb.RxFrags, labelTotal},
		{u.UAP.ApTxPackets, counter, bb.TxPackets, labelTotal},
		{u.UAP.ApTxBytes, counter, bb.TxBytes, labelTotal},
		{u.UAP.ApTxErrors, counter, bb.TxErrors, labelTotal},
		{u.UAP.ApTxDropped, counter, bb.TxDropped, labelTotal},
		{u.UAP.ApTxRetries, counter, bb.TxRetries, labelTotal},
		{u.UAP.WifiTxAttempts, counter, bb.WifiTxAttempts, labelTotal},
		{u.UAP.MacFilterRejections, counter, bb.MacFilterRejections, labelTotal},
		{u.UAP.ApWifiTxDropped, counter, bb.WifiTxDropped, labelTotal},
	})

	// Export wifi0 radio stats (5GHz)
	labelWifi0 := []string{"wifi0", labels[1], labels[2], labels[3]}
	r.send([]*metric{
		{u.UAP.ApRxPackets, counter, bb.Wifi0RxPackets, labelWifi0},
		{u.UAP.ApRxBytes, counter, bb.Wifi0RxBytes, labelWifi0},
		{u.UAP.ApRxErrors, counter, bb.Wifi0RxErrors, labelWifi0},
		{u.UAP.ApRxDropped, counter, bb.Wifi0RxDropped, labelWifi0},
		{u.UAP.ApRxCrypts, counter, bb.Wifi0RxCrypts, labelWifi0},
		{u.UAP.ApRxFrags, counter, bb.Wifi0RxFrags, labelWifi0},
		{u.UAP.ApTxPackets, counter, bb.Wifi0TxPackets, labelWifi0},
		{u.UAP.ApTxBytes, counter, bb.Wifi0TxBytes, labelWifi0},
		{u.UAP.ApTxErrors, counter, bb.Wifi0TxErrors, labelWifi0},
		{u.UAP.ApTxDropped, counter, bb.Wifi0TxDropped, labelWifi0},
		{u.UAP.ApTxRetries, counter, bb.Wifi0TxRetries, labelWifi0},
		{u.UAP.WifiTxAttempts, counter, bb.Wifi0WifiTxAttempts, labelWifi0},
		{u.UAP.MacFilterRejections, counter, bb.Wifi0MacFilterRejections, labelWifi0},
		{u.UAP.ApWifiTxDropped, counter, bb.Wifi0WifiTxDropped, labelWifi0},
	})

	// Export terra2 radio stats (60GHz - 802.11ad)
	labelTerra2 := []string{"terra2", labels[1], labels[2], labels[3]}
	r.send([]*metric{
		{u.UAP.ApRxPackets, counter, bb.Terra2RxPackets, labelTerra2},
		{u.UAP.ApRxBytes, counter, bb.Terra2RxBytes, labelTerra2},
		{u.UAP.ApRxErrors, counter, bb.Terra2RxErrors, labelTerra2},
		{u.UAP.ApRxDropped, counter, bb.Terra2RxDropped, labelTerra2},
		{u.UAP.ApRxCrypts, counter, bb.Terra2RxCrypts, labelTerra2},
		{u.UAP.ApRxFrags, counter, bb.Terra2RxFrags, labelTerra2},
		{u.UAP.ApTxPackets, counter, bb.Terra2TxPackets, labelTerra2},
		{u.UAP.ApTxBytes, counter, bb.Terra2TxBytes, labelTerra2},
		{u.UAP.ApTxErrors, counter, bb.Terra2TxErrors, labelTerra2},
		{u.UAP.ApTxDropped, counter, bb.Terra2TxDropped, labelTerra2},
		{u.UAP.ApTxRetries, counter, bb.Terra2TxRetries, labelTerra2},
		{u.UAP.WifiTxAttempts, counter, bb.Terra2WifiTxAttempts, labelTerra2},
		{u.UAP.MacFilterRejections, counter, bb.Terra2MacFilterRejections, labelTerra2},
		{u.UAP.ApWifiTxDropped, counter, bb.Terra2WifiTxDropped, labelTerra2},
	})

	// Export user stats for wifi0
	labelUserWifi0 := []string{"user-wifi0", labels[1], labels[2], labels[3]}
	r.send([]*metric{
		{u.UAP.ApRxPackets, counter, bb.UserWifi0RxPackets, labelUserWifi0},
		{u.UAP.ApRxBytes, counter, bb.UserWifi0RxBytes, labelUserWifi0},
		{u.UAP.ApRxErrors, counter, bb.UserWifi0RxErrors, labelUserWifi0},
		{u.UAP.ApRxDropped, counter, bb.UserWifi0RxDropped, labelUserWifi0},
		{u.UAP.ApRxCrypts, counter, bb.UserWifi0RxCrypts, labelUserWifi0},
		{u.UAP.ApRxFrags, counter, bb.UserWifi0RxFrags, labelUserWifi0},
		{u.UAP.ApTxPackets, counter, bb.UserWifi0TxPackets, labelUserWifi0},
		{u.UAP.ApTxBytes, counter, bb.UserWifi0TxBytes, labelUserWifi0},
		{u.UAP.ApTxErrors, counter, bb.UserWifi0TxErrors, labelUserWifi0},
		{u.UAP.ApTxDropped, counter, bb.UserWifi0TxDropped, labelUserWifi0},
		{u.UAP.ApTxRetries, counter, bb.UserWifi0TxRetries, labelUserWifi0},
		{u.UAP.WifiTxAttempts, counter, bb.UserWifi0WifiTxAttempts, labelUserWifi0},
		{u.UAP.MacFilterRejections, counter, bb.UserWifi0MacFilterRejections, labelUserWifi0},
		{u.UAP.ApWifiTxDropped, counter, bb.UserWifi0WifiTxDropped, labelUserWifi0},
	})

	// Export user stats for terra2 (60GHz)
	labelUserTerra2 := []string{"user-terra2", labels[1], labels[2], labels[3]}
	r.send([]*metric{
		{u.UAP.ApRxPackets, counter, bb.UserTerra2RxPackets, labelUserTerra2},
		{u.UAP.ApRxBytes, counter, bb.UserTerra2RxBytes, labelUserTerra2},
		{u.UAP.ApRxErrors, counter, bb.UserTerra2RxErrors, labelUserTerra2},
		{u.UAP.ApRxDropped, counter, bb.UserTerra2RxDropped, labelUserTerra2},
		{u.UAP.ApRxCrypts, counter, bb.UserTerra2RxCrypts, labelUserTerra2},
		{u.UAP.ApRxFrags, counter, bb.UserTerra2RxFrags, labelUserTerra2},
		{u.UAP.ApTxPackets, counter, bb.UserTerra2TxPackets, labelUserTerra2},
		{u.UAP.ApTxBytes, counter, bb.UserTerra2TxBytes, labelUserTerra2},
		{u.UAP.ApTxErrors, counter, bb.UserTerra2TxErrors, labelUserTerra2},
		{u.UAP.ApTxDropped, counter, bb.UserTerra2TxDropped, labelUserTerra2},
		{u.UAP.ApTxRetries, counter, bb.UserTerra2TxRetries, labelUserTerra2},
		{u.UAP.WifiTxAttempts, counter, bb.UserTerra2WifiTxAttempts, labelUserTerra2},
		{u.UAP.MacFilterRejections, counter, bb.UserTerra2MacFilterRejections, labelUserTerra2},
		{u.UAP.ApWifiTxDropped, counter, bb.UserTerra2WifiTxDropped, labelUserTerra2},
	})
}

// exportP2Pstats exports point-to-point link statistics for UBB devices.
func (u *promUnifi) exportP2Pstats(r report, labels []string, p2p *unifi.P2PStats) {
	r.send([]*metric{
		{u.Device.Counter, gauge, p2p.RXRate.Val, append(labels, "p2p_rx_rate")},
		{u.Device.Counter, gauge, p2p.TXRate.Val, append(labels, "p2p_tx_rate")},
		{u.Device.Counter, gauge, p2p.Throughput.Val, append(labels, "p2p_throughput")},
	})
}
