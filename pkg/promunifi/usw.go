package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type usw struct {
	// Switch "total" traffic stats
	SwRxPackets   *prometheus.Desc
	SwRxBytes     *prometheus.Desc
	SwRxErrors    *prometheus.Desc
	SwRxDropped   *prometheus.Desc
	SwRxCrypts    *prometheus.Desc
	SwRxFrags     *prometheus.Desc
	SwTxPackets   *prometheus.Desc
	SwTxBytes     *prometheus.Desc
	SwTxErrors    *prometheus.Desc
	SwTxDropped   *prometheus.Desc
	SwTxRetries   *prometheus.Desc
	SwRxMulticast *prometheus.Desc
	SwRxBroadcast *prometheus.Desc
	SwTxMulticast *prometheus.Desc
	SwTxBroadcast *prometheus.Desc
	SwBytes       *prometheus.Desc
	// Port data.
	PoeCurrent     *prometheus.Desc
	PoePower       *prometheus.Desc
	PoeVoltage     *prometheus.Desc
	RxBroadcast    *prometheus.Desc
	RxBytes        *prometheus.Desc
	RxBytesR       *prometheus.Desc
	RxDropped      *prometheus.Desc
	RxErrors       *prometheus.Desc
	RxMulticast    *prometheus.Desc
	RxPackets      *prometheus.Desc
	Satisfaction   *prometheus.Desc
	Speed          *prometheus.Desc
	TxBroadcast    *prometheus.Desc
	TxBytes        *prometheus.Desc
	TxBytesR       *prometheus.Desc
	TxDropped      *prometheus.Desc
	TxErrors       *prometheus.Desc
	TxMulticast    *prometheus.Desc
	TxPackets      *prometheus.Desc
	SFPCurrent     *prometheus.Desc
	SFPRxPower     *prometheus.Desc
	SFPTemperature *prometheus.Desc
	SFPTxPower     *prometheus.Desc
	SFPVoltage     *prometheus.Desc
	// other
	Upgradeable *prometheus.Desc
}

func descUSW(ns string) *usw {
	pns := ns + "port_"
	sfp := pns + "sfp_"
	labelS := []string{"site_name", "name", "source"}
	labelP := []string{"port_id", "port_num", "port_name", "port_mac", "port_ip", "site_name", "name", "source"}
	labelF := []string{
		"sfp_part", "sfp_vendor", "sfp_serial", "sfp_compliance",
		"port_id", "port_num", "port_name", "port_mac", "port_ip", "site_name", "name", "source",
	}
	nd := prometheus.NewDesc

	return &usw{
		// This data may be derivable by sum()ing the port data.
		SwRxPackets:   nd(ns+"switch_receive_packets_total", "Switch Packets Received Total", labelS, nil),
		SwRxBytes:     nd(ns+"switch_receive_bytes_total", "Switch Bytes Received Total", labelS, nil),
		SwRxErrors:    nd(ns+"switch_receive_errors_total", "Switch Errors Received Total", labelS, nil),
		SwRxDropped:   nd(ns+"switch_receive_dropped_total", "Switch Dropped Received Total", labelS, nil),
		SwRxCrypts:    nd(ns+"switch_receive_crypts_total", "Switch Crypts Received Total", labelS, nil),
		SwRxFrags:     nd(ns+"switch_receive_frags_total", "Switch Frags Received Total", labelS, nil),
		SwTxPackets:   nd(ns+"switch_transmit_packets_total", "Switch Packets Transmit Total", labelS, nil),
		SwTxBytes:     nd(ns+"switch_transmit_bytes_total", "Switch Bytes Transmit Total", labelS, nil),
		SwTxErrors:    nd(ns+"switch_transmit_errors_total", "Switch Errors Transmit Total", labelS, nil),
		SwTxDropped:   nd(ns+"switch_transmit_dropped_total", "Switch Dropped Transmit Total", labelS, nil),
		SwTxRetries:   nd(ns+"switch_transmit_retries_total", "Switch Retries Transmit Total", labelS, nil),
		SwRxMulticast: nd(ns+"switch_receive_multicast_total", "Switch Multicast Receive Total", labelS, nil),
		SwRxBroadcast: nd(ns+"switch_receive_broadcast_total", "Switch Broadcast Receive Total", labelS, nil),
		SwTxMulticast: nd(ns+"switch_transmit_multicast_total", "Switch Multicast Transmit Total", labelS, nil),
		SwTxBroadcast: nd(ns+"switch_transmit_broadcast_total", "Switch Broadcast Transmit Total", labelS, nil),
		SwBytes:       nd(ns+"switch_bytes_total", "Switch Bytes Transferred Total", labelS, nil),
		// per-port data
		PoeCurrent:     nd(pns+"poe_amperes", "POE Current", labelP, nil),
		PoePower:       nd(pns+"poe_watts", "POE Power", labelP, nil),
		PoeVoltage:     nd(pns+"poe_volts", "POE Voltage", labelP, nil),
		RxBroadcast:    nd(pns+"receive_broadcast_total", "Receive Broadcast", labelP, nil),
		RxBytes:        nd(pns+"receive_bytes_total", "Total Receive Bytes", labelP, nil),
		RxBytesR:       nd(pns+"receive_rate_bytes", "Receive Bytes Rate", labelP, nil),
		RxDropped:      nd(pns+"receive_dropped_total", "Total Receive Dropped", labelP, nil),
		RxErrors:       nd(pns+"receive_errors_total", "Total Receive Errors", labelP, nil),
		RxMulticast:    nd(pns+"receive_multicast_total", "Total Receive Multicast", labelP, nil),
		RxPackets:      nd(pns+"receive_packets_total", "Total Receive Packets", labelP, nil),
		Satisfaction:   nd(pns+"satisfaction_ratio", "Satisfaction", labelP, nil),
		Speed:          nd(pns+"port_speed_bps", "Speed", labelP, nil),
		TxBroadcast:    nd(pns+"transmit_broadcast_total", "Total Transmit Broadcast", labelP, nil),
		TxBytes:        nd(pns+"transmit_bytes_total", "Total Transmit Bytes", labelP, nil),
		TxBytesR:       nd(pns+"transmit_rate_bytes", "Transmit Bytes Rate", labelP, nil),
		TxDropped:      nd(pns+"transmit_dropped_total", "Total Transmit Dropped", labelP, nil),
		TxErrors:       nd(pns+"transmit_errors_total", "Total Transmit Errors", labelP, nil),
		TxMulticast:    nd(pns+"transmit_multicast_total", "Total Tranmist Multicast", labelP, nil),
		TxPackets:      nd(pns+"transmit_packets_total", "Total Transmit Packets", labelP, nil),
		SFPCurrent:     nd(sfp+"current", "SFP Current", labelF, nil),
		SFPRxPower:     nd(sfp+"rx_power", "SFP Receive Power", labelF, nil),
		SFPTemperature: nd(sfp+"temperature", "SFP Temperature", labelF, nil),
		SFPTxPower:     nd(sfp+"tx_power", "SFP Transmit Power", labelF, nil),
		SFPVoltage:     nd(sfp+"voltage", "SFP Voltage", labelF, nil),
		// other data
		Upgradeable: nd(ns+"upgradeable", "Upgrade-able", labelS, nil),
	}
}

func (u *promUnifi) exportUSW(r report, d *unifi.USW) {
	if !d.Adopted.Val || d.Locating.Val {
		return
	}

	labels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	infoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}

	u.exportUSWstats(r, labels, d.Stat.Sw)
	u.exportPRTtable(r, labels, d.PortTable)
	u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)
	u.exportSYSstats(r, labels, d.SysStats, d.SystemStats)
	u.exportSTAcount(r, labels, d.UserNumSta, d.GuestNumSta)
	r.send([]*metric{
		{u.Device.Info, gauge, 1.0, append(labels, infoLabels...)},
		{u.Device.Uptime, gauge, d.Uptime, labels},
		{u.Device.Upgradeable, gauge, d.Upgradable.Val, labels},
	})

	// Switch System Data.
	if d.HasTemperature.Val {
		r.send([]*metric{{u.Device.Temperature, gauge, d.GeneralTemperature, append(labels, "general", "board")}})
	}

	if d.HasFan.Val {
		r.send([]*metric{{u.Device.FanLevel, gauge, d.FanLevel, labels}})
	}

	if d.TotalMaxPower.Txt != "" {
		r.send([]*metric{{u.Device.TotalMaxPower, gauge, d.TotalMaxPower, labels}})
	}
}

// Switch Stats.
func (u *promUnifi) exportUSWstats(r report, labels []string, sw *unifi.Sw) {
	if sw == nil {
		return
	}

	labelS := labels[1:]

	r.send([]*metric{
		{u.USW.SwRxPackets, counter, sw.RxPackets, labelS},
		{u.USW.SwRxBytes, counter, sw.RxBytes, labelS},
		{u.USW.SwRxErrors, counter, sw.RxErrors, labelS},
		{u.USW.SwRxDropped, counter, sw.RxDropped, labelS},
		{u.USW.SwRxCrypts, counter, sw.RxCrypts, labelS},
		{u.USW.SwRxFrags, counter, sw.RxFrags, labelS},
		{u.USW.SwTxPackets, counter, sw.TxPackets, labelS},
		{u.USW.SwTxBytes, counter, sw.TxBytes, labelS},
		{u.USW.SwTxErrors, counter, sw.TxErrors, labelS},
		{u.USW.SwTxDropped, counter, sw.TxDropped, labelS},
		{u.USW.SwTxRetries, counter, sw.TxRetries, labelS},
		{u.USW.SwRxMulticast, counter, sw.RxMulticast, labelS},
		{u.USW.SwRxBroadcast, counter, sw.RxBroadcast, labelS},
		{u.USW.SwTxMulticast, counter, sw.TxMulticast, labelS},
		{u.USW.SwTxBroadcast, counter, sw.TxBroadcast, labelS},
		{u.USW.SwBytes, counter, sw.Bytes, labelS},
	})
}

// Switch Port Table.
func (u *promUnifi) exportPRTtable(r report, labels []string, pt []unifi.Port) {
	// Per-port data on a switch
	for _, p := range pt {
		if !u.DeadPorts && (!p.Up.Val || !p.Enable.Val) {
			continue
		}

		// Copy labels, and add four new ones.
		labelP := []string{
			labels[2] + " Port " + p.PortIdx.Txt, p.PortIdx.Txt,
			p.Name, p.Mac, p.IP, labels[1], labels[2], labels[3],
		}

		if p.PoeEnable.Val && p.PortPoe.Val {
			r.send([]*metric{
				{u.USW.PoeCurrent, gauge, p.PoeCurrent, labelP},
				{u.USW.PoePower, gauge, p.PoePower, labelP},
				{u.USW.PoeVoltage, gauge, p.PoeVoltage, labelP},
			})
		}

		if p.SFPFound.Val {
			labelF := []string{
				p.SFPPart, p.SFPVendor, p.SFPSerial, p.SFPCompliance,
				labelP[0], labelP[1], labelP[2], labelP[3], labelP[4], labelP[5], labelP[6], labelP[7],
			}

			r.send([]*metric{
				{u.USW.SFPCurrent, gauge, p.SFPCurrent.Val, labelF},
				{u.USW.SFPVoltage, gauge, p.SFPVoltage.Val, labelF},
				{u.USW.SFPTemperature, gauge, p.SFPTemperature.Val, labelF},
				{u.USW.SFPRxPower, gauge, p.SFPRxpower.Val, labelF},
				{u.USW.SFPTxPower, gauge, p.SFPTxpower.Val, labelF},
			})
		}

		r.send([]*metric{
			{u.USW.RxBroadcast, counter, p.RxBroadcast, labelP},
			{u.USW.RxBytes, counter, p.RxBytes, labelP},
			{u.USW.RxBytesR, gauge, p.RxBytesR, labelP},
			{u.USW.RxDropped, counter, p.RxDropped, labelP},
			{u.USW.RxErrors, counter, p.RxErrors, labelP},
			{u.USW.RxMulticast, counter, p.RxMulticast, labelP},
			{u.USW.RxPackets, counter, p.RxPackets, labelP},
			{u.USW.Satisfaction, gauge, p.Satisfaction.Val / 100.0, labelP},
			{u.USW.Speed, gauge, p.Speed.Val * 1000000, labelP},
			{u.USW.TxBroadcast, counter, p.TxBroadcast, labelP},
			{u.USW.TxBytes, counter, p.TxBytes, labelP},
			{u.USW.TxBytesR, gauge, p.TxBytesR, labelP},
			{u.USW.TxDropped, counter, p.TxDropped, labelP},
			{u.USW.TxErrors, counter, p.TxErrors, labelP},
			{u.USW.TxMulticast, counter, p.TxMulticast, labelP},
			{u.USW.TxPackets, counter, p.TxPackets, labelP},
		})
	}
}
