package promunifi

import (
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v6"
)

// umbb holds the cellular modem descriptors for UMBB (mobile broadband) devices, such as the U5G Max.
type umbb struct {
	Internet          *prometheus.Desc
	Rsrp              *prometheus.Desc
	Rsrq              *prometheus.Desc
	Snr               *prometheus.Desc
	RsrpNr            *prometheus.Desc
	RsrqNr            *prometheus.Desc
	SnrNr             *prometheus.Desc
	SignalBars        *prometheus.Desc
	SignalPercent     *prometheus.Desc
	Coverage          *prometheus.Desc
	Roaming           *prometheus.Desc
	RegistrationState *prometheus.Desc
	Carriers          *prometheus.Desc
	BandwidthMhz      *prometheus.Desc
	MaxBitrate        *prometheus.Desc
	Info              *prometheus.Desc
	SimActive         *prometheus.Desc
	SimCardPresent    *prometheus.Desc
	SimPinBlocked     *prometheus.Desc
	SimDataWarning    *prometheus.Desc
	SimDataLimited    *prometheus.Desc
	SimRxBytes        *prometheus.Desc
	SimTxBytes        *prometheus.Desc
	SimStateSeconds   *prometheus.Desc
	SimInfo           *prometheus.Desc
	ApnInfo           *prometheus.Desc
	InetState         *prometheus.Desc
	CellularIPv4      *prometheus.Desc
	CellularIPv4Info  *prometheus.Desc
	CellularIPv6      *prometheus.Desc
	CellularIPv6Info  *prometheus.Desc
	CellularMTU       *prometheus.Desc
	IPv6PrefixInfo    *prometheus.Desc
	IPv6PrefixLength  *prometheus.Desc
	DataLimitEnabled  *prometheus.Desc
	DataSoftLimit     *prometheus.Desc
	DataWarningPct    *prometheus.Desc
	DataUsedRatio     *prometheus.Desc
}

func descUMBB(ns string) *umbb {
	labels := []string{"type", "site_name", "name", "source", "tag"}
	labelS := append(append([]string{}, labels...), "slot", "esim", "spn")
	nd := prometheus.NewDesc

	return &umbb{
		Internet: nd(ns+"mbb_internet", "Modem reports internet connectivity (1) or not (0).", labels, nil),
		Rsrp:     nd(ns+"mbb_rsrp_dbm", "Cellular RSRP of the serving cell in dBm. LTE anchor in 5G NSA mode.", labels, nil),
		Rsrq:     nd(ns+"mbb_rsrq_db", "Cellular RSRQ of the serving cell in dB. LTE anchor in 5G NSA mode.", labels, nil),
		Snr:      nd(ns+"mbb_snr_db", "Cellular signal-to-noise ratio of the serving cell in dB. LTE anchor in 5G NSA mode.", labels, nil),
		RsrpNr:   nd(ns+"mbb_nr_rsrp_dbm", "5G NR RSRP in dBm. Same as mbb_rsrp_dbm in 5G SA mode.", labels, nil),
		RsrqNr:   nd(ns+"mbb_nr_rsrq_db", "5G NR RSRQ in dB. Same as mbb_rsrq_db in 5G SA mode.", labels, nil),
		SnrNr:    nd(ns+"mbb_nr_snr_db", "5G NR signal-to-noise ratio in dB. Same as mbb_snr_db in 5G SA mode.", labels, nil),
		SignalBars: nd(ns+"mbb_signal_bars",
			"Cellular signal strength in bars, as shown by the controller.", labels, nil),
		SignalPercent: nd(ns+"mbb_signal_percent",
			"Cellular signal strength percentage (0-100), as computed by the controller.", labels, nil),
		Coverage: nd(ns+"mbb_coverage", "Modem has cellular coverage (1) or not (0).", labels, nil),
		Roaming:  nd(ns+"mbb_roaming", "Modem is roaming (1) or on its home network (0).", labels, nil),
		RegistrationState: nd(ns+"mbb_registration_state",
			"Cellular network registration state code. Undocumented; likely 0 not registered, 1 registered home, "+
				"2 searching, 3 denied, 4 unknown, 5 registered roaming.", labels, nil),
		Carriers: nd(ns+"mbb_carriers", "Aggregated component carriers in use.", append(labels, "carrier_type"), nil),
		BandwidthMhz: nd(ns+"mbb_bandwidth_mhz", "Total aggregated carrier bandwidth in MHz.",
			append(labels, "direction"), nil),
		MaxBitrate: nd(ns+"mbb_max_bitrate_bps", "Maximum bitrate negotiated with the network, in bits per second.",
			append(labels, "carrier_type", "direction"), nil),
		Info: nd(ns+"mbb_info", "Cellular modem state. Always 1; state and network details are labels.",
			append(labels, "mbb_state", "failover_mode", "rat", "band", "operator", "sa_mode", "current_slot",
				"primary_slot", "mcc", "mnc", "cell_id", "pci", "channel"), nil),
		SimActive:      nd(ns+"mbb_sim_active", "SIM slot is the active data SIM (1) or not (0).", labelS, nil),
		SimCardPresent: nd(ns+"mbb_sim_card_present", "SIM card or eSIM profile is present (1) or not (0).", labelS, nil),
		SimPinBlocked:  nd(ns+"mbb_sim_pin_blocked", "SIM is PIN blocked (1) or not (0).", labelS, nil),
		SimDataWarning: nd(ns+"mbb_sim_data_warning", "SIM has passed its data usage warning threshold (1) or not (0).",
			labelS, nil),
		SimDataLimited: nd(ns+"mbb_sim_data_limited", "SIM has reached its data limit (1) or not (0).", labelS, nil),
		SimRxBytes: nd(ns+"mbb_sim_receive_bytes_total",
			"Bytes received over this SIM. Cumulative; not reset when the modem reboots.", labelS, nil),
		SimTxBytes: nd(ns+"mbb_sim_transmit_bytes_total",
			"Bytes transmitted over this SIM. Cumulative; not reset when the modem reboots.", labelS, nil),
		SimStateSeconds: nd(ns+"mbb_sim_state_seconds", "Seconds the SIM has been in its current display state.",
			labelS, nil),
		SimInfo: nd(ns+"mbb_sim_info", "SIM state. Always 1; display state and last network reject are labels.",
			append(labelS, "display_state", "has_carrier", "metered", "network_reject_type", "network_reject_text"), nil),
		ApnInfo: nd(ns+"mbb_apn_info", "APN profile in use. Always 1; APN and PDP details are labels.",
			append(labelS, "apn", "pdp_type", "auth", "clat_enabled", "roaming"), nil),
		InetState: nd(ns+"mbb_inet_state", "Cellular data session state as reported by the controller.", labelS, nil),
		CellularIPv4: nd(ns+"mbb_cellular_ipv4_present",
			"Modem has a cellular IPv4 address (1) or not (0).", labels, nil),
		CellularIPv4Info: nd(ns+"mbb_cellular_ipv4_info",
			"Cellular IPv4 address. Always 1; address, netmask, gateway and DNS are labels.",
			append(labels, "address", "netmask", "gateway", "dns"), nil),
		CellularIPv6: nd(ns+"mbb_cellular_ipv6_present",
			"Modem has a cellular IPv6 address (1) or not (0).", labels, nil),
		CellularIPv6Info: nd(ns+"mbb_cellular_ipv6_info",
			"Cellular IPv6 address. Always 1; address, gateway and DNS are labels.",
			append(labels, "address", "gateway", "dns"), nil),
		CellularMTU: nd(ns+"mbb_cellular_mtu", "Cellular IPv4 MTU in bytes.", labels, nil),
		IPv6PrefixInfo: nd(ns+"mbb_ipv6_prefix_info",
			"Cellular IPv6 prefix. Always 1; prefix is a label.", append(labels, "prefix"), nil),
		IPv6PrefixLength: nd(ns+"mbb_ipv6_prefix_length",
			"Cellular IPv6 prefix length in bits.", labels, nil),
		DataLimitEnabled: nd(ns+"mbb_sim_data_limit_enabled",
			"SIM has a configured data limit (1) or not (0).", labelS, nil),
		DataSoftLimit: nd(ns+"mbb_sim_data_soft_limit_bytes",
			"Configured data soft limit in bytes; 0 means no limit.", labelS, nil),
		DataWarningPct: nd(ns+"mbb_sim_data_warning_threshold_percent",
			"Configured data warning threshold in percent.", labelS, nil),
		DataUsedRatio: nd(ns+"mbb_sim_data_used_ratio",
			"Data used divided by the configured soft limit, if a limit is set.", labelS, nil),
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
		u.exportMBB(r, labels, d)

		r.send([]*metric{
			{u.Device.Info, gauge, 1.0, append(append([]string{}, baseLabels...), infoLabels...)},
			{u.Device.Uptime, gauge, d.Uptime, labels},
			{u.Device.Upgradeable, gauge, d.Upgradable.Val, labels},
		})
	})
}

func (u *promUnifi) exportMBB(r report, labels []string, d *unifi.UMBB) {
	m := d.Mbb
	rad := m.Radio
	dl, ul := mbbBandwidth(rad)
	info := append(append([]string{}, labels...),
		m.State, m.Mode, rad.Rat, rad.Band, rad.NetworkOperator,
		strconv.FormatBool(rad.SA5GMode.Val), rad.CurrentSlot.Txt, d.MbbOverrides.PrimarySlot.Txt,
		rad.Mcc.Txt, rad.Mnc.Txt, rad.CellID.Txt, rad.Pci.Txt, rad.Channel.Txt)

	r.send([]*metric{
		{u.UMBB.Internet, gauge, d.Internet.Val, labels},
		{u.UMBB.Coverage, gauge, rad.HasCoverage.Val, labels},
		{u.UMBB.Roaming, gauge, rad.Roaming.Val, labels},
		{u.UMBB.RegistrationState, gauge, rad.RegistrationState, labels},
		{u.UMBB.Carriers, gauge, float64(len(rad.CaNr)), append(append([]string{}, labels...), "nr")},
		{u.UMBB.Carriers, gauge, float64(len(rad.CaLte)), append(append([]string{}, labels...), "lte")},
		{u.UMBB.BandwidthMhz, gauge, dl, append(append([]string{}, labels...), "dl")},
		{u.UMBB.BandwidthMhz, gauge, ul, append(append([]string{}, labels...), "ul")},
		{u.UMBB.Info, gauge, 1.0, info},
	})

	ip4 := m.IPSettings
	ip6 := m.IPv6Settings
	r.send([]*metric{
		{u.UMBB.CellularIPv4, gauge, boolToFloat(ip4.IPv4Address != ""), labels},
		{u.UMBB.CellularIPv6, gauge, boolToFloat(ip6.IPv6Address != ""), labels},
	})

	if ip4.IPv4Address != "" {
		r.send([]*metric{{u.UMBB.CellularIPv4Info, gauge, 1.0, append(append([]string{}, labels...),
			ip4.IPv4Address, ip4.IPv4Netmask, ip4.IPv4Gateway, ip4.IPv4DNS)}})
	}

	if ip4.MTU.Val > 0 {
		r.send([]*metric{{u.UMBB.CellularMTU, gauge, ip4.MTU.Val, labels}})
	}

	if ip6.IPv6Address != "" {
		r.send([]*metric{{u.UMBB.CellularIPv6Info, gauge, 1.0, append(append([]string{}, labels...),
			ip6.IPv6Address, ip6.IPv6Gateway, ip6.IPv6DNS)}})

		if prefix, length := mbbPrefix(ip6.IPv6Address); prefix != "" {
			r.send([]*metric{
				{u.UMBB.IPv6PrefixInfo, gauge, 1.0, append(append([]string{}, labels...), prefix)},
				{u.UMBB.IPv6PrefixLength, gauge, length, labels},
			})
		}
	}

	// Signal readings and bitrates are left out when the modem does not report them,
	// so a missing reading is not mistaken for 0 dBm.
	for _, s := range []struct {
		desc *prometheus.Desc
		txt  string
		val  float64
		dims []string
	}{
		{u.UMBB.Rsrp, rad.Rsrp.Txt, rad.Rsrp.Val, nil},
		{u.UMBB.Rsrq, rad.Rsrq.Txt, rad.Rsrq.Val, nil},
		{u.UMBB.Snr, rad.Snr.Txt, rad.Snr.Val, nil},
		{u.UMBB.RsrpNr, rad.RsrpNr.Txt, rad.RsrpNr.Val, nil},
		{u.UMBB.RsrqNr, rad.RsrqNr.Txt, rad.RsrqNr.Val, nil},
		{u.UMBB.SnrNr, rad.SnrNr.Txt, rad.SnrNr.Val, nil},
		{u.UMBB.SignalBars, rad.Signal.Txt, rad.Signal.Val, nil},
		{u.UMBB.SignalPercent, rad.SignalPercent.Txt, rad.SignalPercent.Val, nil},
		{u.UMBB.MaxBitrate, rad.MaxBitrateDl.Txt, rad.MaxBitrateDl.Val, []string{"lte", "dl"}},
		{u.UMBB.MaxBitrate, rad.MaxBitrateUl.Txt, rad.MaxBitrateUl.Val, []string{"lte", "ul"}},
		{u.UMBB.MaxBitrate, rad.MaxBitrateDlNr.Txt, rad.MaxBitrateDlNr.Val, []string{"nr", "dl"}},
		{u.UMBB.MaxBitrate, rad.MaxBitrateUlNr.Txt, rad.MaxBitrateUlNr.Val, []string{"nr", "ul"}},
	} {
		if s.txt != "" {
			r.send([]*metric{{s.desc, gauge, s.val, append(append([]string{}, labels...), s.dims...)}})
		}
	}

	for _, sim := range m.Sim {
		simLabels := append(append([]string{}, labels...),
			sim.Slot.Txt, strconv.FormatBool(sim.Esim.Val), sim.Spn)
		simInfo := append(append([]string{}, simLabels...),
			sim.DisplayState, strconv.FormatBool(sim.HasCarrier.Val), strconv.FormatBool(sim.Metered.Val),
			sim.NetworkRejectType, sim.NetworkRejectText)

		r.send([]*metric{
			{u.UMBB.SimActive, gauge, sim.Active.Val, simLabels},
			{u.UMBB.SimCardPresent, gauge, sim.CardPresent.Val, simLabels},
			{u.UMBB.SimPinBlocked, gauge, sim.PinBlocked.Val, simLabels},
			{u.UMBB.SimDataWarning, gauge, sim.DataWarning.Val, simLabels},
			{u.UMBB.SimDataLimited, gauge, sim.DataLimited.Val, simLabels},
			{u.UMBB.SimRxBytes, counter, sim.RxBytes, simLabels},
			{u.UMBB.SimTxBytes, counter, sim.TxBytes, simLabels},
			{u.UMBB.SimStateSeconds, gauge, sim.DisplayStateElapsed, simLabels},
			{u.UMBB.SimInfo, gauge, 1.0, simInfo},
		})

		apn := sim.CurrentApn
		if apn.Apn != "" {
			r.send([]*metric{{u.UMBB.ApnInfo, gauge, 1.0, append(append([]string{}, simLabels...),
				apn.Apn, apn.PdpType, apn.AuthType, strconv.FormatBool(apn.ClatEnabled.Val),
				strconv.FormatBool(apn.Roaming.Val))}})
		}

		if sim.ConnectionInfo.InetState.Txt != "" {
			r.send([]*metric{{u.UMBB.InetState, gauge, sim.ConnectionInfo.InetState.Val, simLabels}})
		}

		if override, ok := findMBBOverride(d.MbbOverrides.Sim, sim.Slot.Val); ok {
			r.send([]*metric{
				{u.UMBB.DataLimitEnabled, gauge, override.DataLimitEnabled.Val, simLabels},
				{u.UMBB.DataSoftLimit, gauge, override.DataSoftLimitBytes.Val, simLabels},
				{u.UMBB.DataWarningPct, gauge, override.DataWarningThreshold.Val, simLabels},
			})
			if override.DataSoftLimitBytes.Val > 0 {
				used := sim.RxBytes.Val + sim.TxBytes.Val
				r.send([]*metric{{u.UMBB.DataUsedRatio, gauge, used / override.DataSoftLimitBytes.Val, simLabels}})
			}
		}
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

func boolToFloat(v bool) float64 {
	if v {
		return 1
	}

	return 0
}

// mbbPrefix returns the network and prefix length from an address such as
// "2a01:598:ba0a:7710::1/64". The controller reports the cellular address with
// the prefix length, so the prefix is available even on library versions that
// do not yet decode ipv6_ui_prefix.
func mbbPrefix(address string) (string, float64) {
	_, network, err := net.ParseCIDR(address)
	if err != nil {
		return "", 0
	}

	ones, _ := network.Mask.Size()
	if ones == 0 && !strings.Contains(address, "/0") {
		return "", 0
	}

	return network.String(), float64(ones)
}

// findMBBOverride returns the data-plan override for a SIM slot. Some firmware
// omits the slot when only one SIM is configured, so a single entry is used as
// the fallback.
func findMBBOverride(sims []unifi.MBBSimOverride, slot float64) (unifi.MBBSimOverride, bool) {
	for _, sim := range sims {
		if sim.Slot.Val == slot {
			return sim, true
		}
	}

	if len(sims) == 1 {
		return sims[0], true
	}

	return unifi.MBBSimOverride{}, false
}
