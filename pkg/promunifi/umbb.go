package promunifi

import (
	"math"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v6"
)

// umbb holds the cellular modem descriptors for UMBB (mobile broadband) devices, such as the U5G Max.
type umbb struct {
	Rsrp              *prometheus.Desc
	Rsrq              *prometheus.Desc
	Snr               *prometheus.Desc
	SignalBars        *prometheus.Desc
	SignalPercent     *prometheus.Desc
	Coverage          *prometheus.Desc
	Roaming           *prometheus.Desc
	RegistrationState *prometheus.Desc
	Carriers          *prometheus.Desc
	BandwidthMhz      *prometheus.Desc
	Info              *prometheus.Desc
	SimActive         *prometheus.Desc
	SimRxBytes        *prometheus.Desc
	SimTxBytes        *prometheus.Desc
	SimStateSeconds   *prometheus.Desc
	SimRejectAge      *prometheus.Desc
	SimInfo           *prometheus.Desc
}

func descUMBB(ns string) *umbb {
	labels := []string{"type", "site_name", "name", "source", "tag"}
	labelS := append(append([]string{}, labels...), "slot", "esim", "spn")
	nd := prometheus.NewDesc

	return &umbb{
		Rsrp:              nd(ns+"mbb_rsrp_dbm", "Cellular RSRP in dBm.", labels, nil),
		Rsrq:              nd(ns+"mbb_rsrq_db", "Cellular RSRQ in dB.", labels, nil),
		Snr:               nd(ns+"mbb_snr_db", "Cellular signal-to-noise ratio in dB.", labels, nil),
		SignalBars:        nd(ns+"mbb_signal_bars", "Cellular signal strength in bars (0-5).", labels, nil),
		SignalPercent:     nd(ns+"mbb_signal_percent", "Cellular signal strength percentage.", labels, nil),
		Coverage:          nd(ns+"mbb_coverage", "Modem has cellular coverage (1) or not (0).", labels, nil),
		Roaming:           nd(ns+"mbb_roaming", "Modem is roaming (1) or on its home network (0).", labels, nil),
		RegistrationState: nd(ns+"mbb_registration_state", "Cellular network registration state code.", labels, nil),
		Carriers:          nd(ns+"mbb_carriers", "Aggregated component carriers in use.", append(labels, "carrier_type"), nil),
		BandwidthMhz:      nd(ns+"mbb_bandwidth_mhz", "Total aggregated carrier bandwidth in MHz.", append(labels, "direction"), nil),
		Info: nd(ns+"mbb_info", "Cellular modem state. Always 1; state and network details are labels.",
			append(labels, "mbb_state", "mbb_mode", "rat", "band", "operator", "sa_mode", "current_slot"), nil),
		SimActive:       nd(ns+"mbb_sim_active", "SIM slot is the active data SIM (1) or not (0).", labelS, nil),
		SimRxBytes:      nd(ns+"mbb_sim_receive_bytes_total", "Bytes received over this SIM.", labelS, nil),
		SimTxBytes:      nd(ns+"mbb_sim_transmit_bytes_total", "Bytes transmitted over this SIM.", labelS, nil),
		SimStateSeconds: nd(ns+"mbb_sim_state_seconds", "Seconds the SIM has been in its current display state.", labelS, nil),
		SimRejectAge: nd(ns+"mbb_sim_network_reject_age",
			"Age of the last carrier network reject, as reported by the controller.", labelS, nil),
		SimInfo: nd(ns+"mbb_sim_info", "SIM state. Always 1; display state and last network reject are labels.",
			append(labelS, "display_state", "network_reject_type", "network_reject_text"), nil),
	}
}

// exportUMBB exports metrics for UMBB (mobile broadband) devices, such as the U5G Max.
func (u *promUnifi) exportUMBB(r report, d *unifi.UMBB) {
	if !d.Adopted.Val {
		return
	}

	baseLabels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	baseInfoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}

	u.exportWithTags(r, d.Tags, func(tagLabels []string) {
		tag := tagLabels[0]
		labels := append(append([]string{}, baseLabels...), tag)
		infoLabels := append(append([]string{}, baseInfoLabels...), tag)

		u.exportSYSstats(r, labels, d.SysStats, d.SystemStats)
		u.exportDeviceUplink(r, labels, d.Uplink)
		u.exportMBB(r, labels, d.Mbb)

		r.send([]*metric{
			{u.Device.Info, gauge, 1.0, append(append([]string{}, baseLabels...), infoLabels...)},
			{u.Device.Uptime, gauge, d.Uptime, labels},
			{u.Device.Upgradeable, gauge, d.Upgradable.Val, labels},
		})
	})
}

func (u *promUnifi) exportMBB(r report, labels []string, m unifi.MBB) {
	rad := m.Radio
	dl, ul := mbbBandwidth(rad)
	info := append(append([]string{}, labels...),
		m.State, m.Mode, rad.Rat, rad.Band, rad.NetworkOperator,
		strconv.FormatBool(rad.SA5GMode.Val), rad.CurrentSlot.Txt)

	r.send([]*metric{
		{u.UMBB.Rsrp, gauge, rad.Rsrp, labels},
		{u.UMBB.Rsrq, gauge, rad.Rsrq, labels},
		{u.UMBB.Snr, gauge, rad.Snr, labels},
		{u.UMBB.SignalBars, gauge, rad.Signal, labels},
		{u.UMBB.SignalPercent, gauge, rad.SignalPercent, labels},
		{u.UMBB.Coverage, gauge, rad.HasCoverage.Val, labels},
		{u.UMBB.Roaming, gauge, rad.Roaming.Val, labels},
		{u.UMBB.RegistrationState, gauge, rad.RegistrationState, labels},
		{u.UMBB.Carriers, gauge, float64(len(rad.CaNr)), append(append([]string{}, labels...), "nr")},
		{u.UMBB.Carriers, gauge, float64(len(rad.CaLte)), append(append([]string{}, labels...), "lte")},
		{u.UMBB.BandwidthMhz, gauge, dl, append(append([]string{}, labels...), "dl")},
		{u.UMBB.BandwidthMhz, gauge, ul, append(append([]string{}, labels...), "ul")},
		{u.UMBB.Info, gauge, 1.0, info},
	})

	for _, sim := range m.Sim {
		simLabels := append(append([]string{}, labels...),
			sim.Slot.Txt, strconv.FormatBool(sim.Esim.Val), sim.Spn)
		simInfo := append(append([]string{}, simLabels...),
			sim.DisplayState, sim.NetworkRejectType, sim.NetworkRejectText)

		r.send([]*metric{
			{u.UMBB.SimActive, gauge, sim.Active.Val, simLabels},
			{u.UMBB.SimRxBytes, counter, sim.RxBytes, simLabels},
			{u.UMBB.SimTxBytes, counter, sim.TxBytes, simLabels},
			{u.UMBB.SimStateSeconds, gauge, sim.DisplayStateElapsed, simLabels},
			{u.UMBB.SimRejectAge, gauge, sim.NetworkRejectAge, simLabels},
			{u.UMBB.SimInfo, gauge, 1.0, simInfo},
		})
	}
}

// mbbBandwidth sums the downlink and uplink bandwidth of every aggregated carrier, NR and LTE.
// A carrier with no uplink reports a null bandwidth, which decodes as zero.
// A sum that overflows is reported as zero: InfluxDB rejects a whole point with an infinite field.
func mbbBandwidth(rad unifi.MBBRadio) (dl, ul float64) {
	for _, c := range rad.CaNr {
		dl += c.DlBwMhz.Val
		ul += c.UlBwMhz.Val
	}

	for _, c := range rad.CaLte {
		dl += c.DlBwMhz.Val
		ul += c.UlBwMhz.Val
	}

	return finiteOrZero(dl), finiteOrZero(ul)
}

func finiteOrZero(v float64) float64 {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return 0
	}

	return v
}
