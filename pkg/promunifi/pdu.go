package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type pdu struct {
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
	// power
	CycleEnabled      *prometheus.Desc
	RelayState        *prometheus.Desc
	OutletCaps        *prometheus.Desc
	OutletCurrent     *prometheus.Desc
	OutletPower       *prometheus.Desc
	OutletPowerFactor *prometheus.Desc
	OutletVoltage     *prometheus.Desc
}

func descPDU(ns string) *pdu {
	outlet := ns + "outlet_"
	pns := ns + "port_"
	sfp := pns + "sfp_"
	labelS := []string{"site_name", "name", "source"}
	labelP := []string{"port_id", "port_num", "port_name", "port_mac", "port_ip", "site_name", "name", "source"}
	labelF := []string{
		"sfp_part", "sfp_vendor", "sfp_serial", "sfp_compliance",
		"port_id", "port_num", "port_name", "port_mac", "port_ip", "site_name", "name", "source",
	}
	labelO := []string{
		"outlet_description", "outlet_index", "outlet_name", "site_name", "name", "source",
	}
	nd := prometheus.NewDesc

	return &pdu{
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
		// power
		CycleEnabled:      nd(outlet+"cycle_enabled", "Cycle Enabled", labelO, nil),
		RelayState:        nd(outlet+"relay_state", "Relay State", labelO, nil),
		OutletCaps:        nd(outlet+"outlet_caps", "Outlet Caps", labelO, nil),
		OutletCurrent:     nd(outlet+"outlet_current", "Outlet Current", labelO, nil),
		OutletPower:       nd(outlet+"outlet_power", "Outlet Power", labelO, nil),
		OutletPowerFactor: nd(outlet+"outlet_power_factor", "Outlet Power Factor", labelO, nil),
		OutletVoltage:     nd(outlet+"outlet_voltage", "Outlet Voltage", labelO, nil),
	}
}

func (u *promUnifi) exportPDU(r report, d *unifi.PDU) {
	if !d.Adopted.Val || d.Locating.Val {
		return
	}

	labels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	infoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}

	u.exportPDUstats(r, labels, d.Stat.Sw)
	u.exportPDUPrtTable(r, labels, d.PortTable)
	u.exportPDUOutletTable(r, labels, d.OutletTable, d.OutletOverrides)
	u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)
	u.exportSYSstats(r, labels, d.SysStats, d.SystemStats)
	u.exportSTAcount(r, labels, d.UserNumSta, d.GuestNumSta)
	r.send([]*metric{
		{u.Device.Info, gauge, 1.0, append(labels, infoLabels...)},
		{u.Device.Uptime, gauge, d.Uptime, labels},
		{u.Device.Upgradeable, gauge, d.Upgradeable.Val, labels},
	})

	// Switch System Data.
	if d.OutletACPowerConsumption.Txt != "" {
		r.send([]*metric{{u.Device.OutletACPowerConsumption, gauge, d.OutletACPowerConsumption, labels}})
	}

	if d.PowerSource.Txt != "" {
		r.send([]*metric{{u.Device.PowerSource, gauge, d.PowerSource, labels}})
	}

	if d.TotalMaxPower.Txt != "" {
		r.send([]*metric{{u.Device.TotalMaxPower, gauge, d.TotalMaxPower, labels}})
	}
}

// Switch Stats.
func (u *promUnifi) exportPDUstats(r report, labels []string, sw *unifi.Sw) {
	if sw == nil {
		return
	}

	labelS := labels[1:]

	r.send([]*metric{
		{u.PDU.SwRxPackets, counter, sw.RxPackets, labelS},
		{u.PDU.SwRxBytes, counter, sw.RxBytes, labelS},
		{u.PDU.SwRxErrors, counter, sw.RxErrors, labelS},
		{u.PDU.SwRxDropped, counter, sw.RxDropped, labelS},
		{u.PDU.SwRxCrypts, counter, sw.RxCrypts, labelS},
		{u.PDU.SwRxFrags, counter, sw.RxFrags, labelS},
		{u.PDU.SwTxPackets, counter, sw.TxPackets, labelS},
		{u.PDU.SwTxBytes, counter, sw.TxBytes, labelS},
		{u.PDU.SwTxErrors, counter, sw.TxErrors, labelS},
		{u.PDU.SwTxDropped, counter, sw.TxDropped, labelS},
		{u.PDU.SwTxRetries, counter, sw.TxRetries, labelS},
		{u.PDU.SwRxMulticast, counter, sw.RxMulticast, labelS},
		{u.PDU.SwRxBroadcast, counter, sw.RxBroadcast, labelS},
		{u.PDU.SwTxMulticast, counter, sw.TxMulticast, labelS},
		{u.PDU.SwTxBroadcast, counter, sw.TxBroadcast, labelS},
		{u.PDU.SwBytes, counter, sw.Bytes, labelS},
	})
}

// Switch Port Table.
func (u *promUnifi) exportPDUPrtTable(r report, labels []string, pt []unifi.Port) {
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
				{u.PDU.PoeCurrent, gauge, p.PoeCurrent, labelP},
				{u.PDU.PoePower, gauge, p.PoePower, labelP},
				{u.PDU.PoeVoltage, gauge, p.PoeVoltage, labelP},
			})
		}

		if p.SFPFound.Val {
			labelF := []string{
				p.SFPPart, p.SFPVendor, p.SFPSerial, p.SFPCompliance,
				labelP[0], labelP[1], labelP[2], labelP[3], labelP[4], labelP[5], labelP[6], labelP[7],
			}

			r.send([]*metric{
				{u.PDU.SFPCurrent, gauge, p.SFPCurrent.Val, labelF},
				{u.PDU.SFPVoltage, gauge, p.SFPVoltage.Val, labelF},
				{u.PDU.SFPTemperature, gauge, p.SFPTemperature.Val, labelF},
				{u.PDU.SFPRxPower, gauge, p.SFPRxpower.Val, labelF},
				{u.PDU.SFPTxPower, gauge, p.SFPTxpower.Val, labelF},
			})
		}

		r.send([]*metric{
			{u.PDU.RxBroadcast, counter, p.RxBroadcast, labelP},
			{u.PDU.RxBytes, counter, p.RxBytes, labelP},
			{u.PDU.RxBytesR, gauge, p.RxBytesR, labelP},
			{u.PDU.RxDropped, counter, p.RxDropped, labelP},
			{u.PDU.RxErrors, counter, p.RxErrors, labelP},
			{u.PDU.RxMulticast, counter, p.RxMulticast, labelP},
			{u.PDU.RxPackets, counter, p.RxPackets, labelP},
			{u.PDU.Satisfaction, gauge, p.Satisfaction.Val / 100.0, labelP},
			{u.PDU.Speed, gauge, p.Speed.Val * 1000000, labelP},
			{u.PDU.TxBroadcast, counter, p.TxBroadcast, labelP},
			{u.PDU.TxBytes, counter, p.TxBytes, labelP},
			{u.PDU.TxBytesR, gauge, p.TxBytesR, labelP},
			{u.PDU.TxDropped, counter, p.TxDropped, labelP},
			{u.PDU.TxErrors, counter, p.TxErrors, labelP},
			{u.PDU.TxMulticast, counter, p.TxMulticast, labelP},
			{u.PDU.TxPackets, counter, p.TxPackets, labelP},
		})
	}
}

// Switch Port Table.
func (u *promUnifi) exportPDUOutletTable(r report, labels []string, ot []unifi.OutletTable, oto []unifi.OutletOverride) {
	// Per-outlet data on a switch
	for _, o := range ot {
		// Copy labels, and add four new ones.
		labelOutlet := []string{
			labels[2] + " Outlet " + o.Index.Txt, o.Index.Txt,
			o.Name, labels[1], labels[2], labels[3],
		}

		r.send([]*metric{
			{u.PDU.CycleEnabled, counter, o.CycleEnabled, labelOutlet},
			{u.PDU.RelayState, counter, o.RelayState, labelOutlet},
			{u.PDU.OutletCaps, counter, o.OutletCaps, labelOutlet},
			{u.PDU.OutletCurrent, gauge, o.OutletCurrent, labelOutlet},
			{u.PDU.OutletPower, counter, o.OutletPower, labelOutlet},
			{u.PDU.OutletPowerFactor, counter, o.OutletPowerFactor, labelOutlet},
			{u.PDU.OutletVoltage, counter, o.OutletVoltage, labelOutlet},
		})
	}

	// Per-outlet data on a switch
	for _, o := range oto {
		// Copy labels, and add four new ones.
		labelOutlet := []string{
			labels[2] + " Outlet Override " + o.Index.Txt, o.Index.Txt,
			o.Name, labels[1], labels[2], labels[3],
		}

		r.send([]*metric{
			{u.PDU.CycleEnabled, counter, o.CycleEnabled, labelOutlet},
			{u.PDU.RelayState, counter, o.RelayState, labelOutlet},
		})
	}
}
