package unifi

import (
	"encoding/json"
	"time"
)

// USW represents all the data from the Ubiquiti Controller for a Unifi Switch.
type USW struct {
	site                 *Site
	SourceName           string           `json:"-"`
	SiteName             string           `json:"-"`
	ID                   string           `json:"_id"`
	Adopted              FlexBool         `json:"adopted"`
	BoardRev             FlexInt          `json:"board_rev"`
	Cfgversion           string           `json:"cfgversion"`
	ConfigNetwork        *ConfigNetwork   `json:"config_network"`
	Dot1XPortctrlEnabled FlexBool         `json:"dot1x_portctrl_enabled"`
	EthernetTable        []*EthernetTable `json:"ethernet_table"`
	FlowctrlEnabled      FlexBool         `json:"flowctrl_enabled"`
	FwCaps               FlexInt          `json:"fw_caps"`
	HasFan               FlexBool         `json:"has_fan"`
	HasTemperature       FlexBool         `json:"has_temperature"`
	InformIP             string           `json:"inform_ip"`
	InformURL            string           `json:"inform_url"`
	IP                   string           `json:"ip"`
	JumboframeEnabled    FlexBool         `json:"jumboframe_enabled"`
	LedOverride          string           `json:"led_override"`
	LicenseState         string           `json:"license_state"`
	Mac                  string           `json:"mac"`
	Model                string           `json:"model"`
	Name                 string           `json:"name"`
	OutdoorModeOverride  string           `json:"outdoor_mode_override"`
	PortOverrides        []struct {
		Name       string  `json:"name,omitempty"`
		PoeMode    string  `json:"poe_mode,omitempty"`
		PortIdx    FlexInt `json:"port_idx"`
		PortconfID string  `json:"portconf_id"`
	} `json:"port_overrides"`
	PortTable             []Port           `json:"port_table"`
	Serial                string           `json:"serial"`
	SiteID                string           `json:"site_id"`
	StpPriority           FlexInt          `json:"stp_priority"`
	StpVersion            string           `json:"stp_version"`
	Type                  string           `json:"type"`
	Version               string           `json:"version"`
	RequiredVersion       string           `json:"required_version"`
	SwitchCaps            *SwitchCaps      `json:"switch_caps"`
	HwCaps                FlexInt          `json:"hw_caps"`
	Unsupported           FlexBool         `json:"unsupported"`
	UnsupportedReason     FlexInt          `json:"unsupported_reason"`
	SysErrorCaps          FlexInt          `json:"sys_error_caps"`
	DeviceID              string           `json:"device_id"`
	State                 FlexInt          `json:"state"`
	LastSeen              FlexInt          `json:"last_seen"`
	Upgradable            FlexBool         `json:"upgradable,omitempty"`
	AdoptableWhenUpgraded FlexBool         `json:"adoptable_when_upgraded,omitempty"`
	Rollupgrade           FlexBool         `json:"rollupgrade,omitempty"`
	KnownCfgversion       string           `json:"known_cfgversion"`
	Uptime                FlexInt          `json:"uptime"`
	Locating              FlexBool         `json:"locating"`
	ConnectRequestIP      string           `json:"connect_request_ip"`
	ConnectRequestPort    string           `json:"connect_request_port"`
	SysStats              SysStats         `json:"sys_stats"`
	SystemStats           SystemStats      `json:"system-stats"`
	FanLevel              FlexInt          `json:"fan_level"`
	GeneralTemperature    FlexInt          `json:"general_temperature"`
	Overheating           FlexBool         `json:"overheating"`
	TotalMaxPower         FlexInt          `json:"total_max_power"`
	DownlinkTable         []*DownlinkTable `json:"downlink_table"`
	Uplink                Uplink           `json:"uplink"`
	LastUplink            struct {
		UplinkMac string `json:"uplink_mac"`
	} `json:"last_uplink"`
	UplinkDepth FlexInt `json:"uplink_depth"`
	Stat        USWStat `json:"stat"`
	TxBytes     FlexInt `json:"tx_bytes"`
	RxBytes     FlexInt `json:"rx_bytes"`
	Bytes       FlexInt `json:"bytes"`
	NumSta      FlexInt `json:"num_sta"`
	UserNumSta  FlexInt `json:"user-num_sta"`
	GuestNumSta FlexInt `json:"guest-num_sta"`
}

type SwitchCaps struct {
	FeatureCaps          FlexInt `json:"feature_caps"`
	MaxMirrorSessions    FlexInt `json:"max_mirror_sessions"`
	MaxAggregateSessions FlexInt `json:"max_aggregate_sessions"`
}

// MacTable is a newer feature on some switched ports.
type MacTable struct {
	Age           int64    `json:"age"`
	Authorized    FlexBool `json:"authorized"`
	Hostname      string   `json:"hostname"`
	IP            string   `json:"ip"`
	LastReachable int64    `json:"lastReachable"`
	Mac           string   `json:"mac"`
}

// Port is a physical connection on a USW or Gateway.
// Not every port has the same capabilities.
type Port struct {
	AggregatedBy       FlexBool   `json:"aggregated_by"`
	Autoneg            FlexBool   `json:"autoneg,omitempty"`
	BytesR             FlexInt    `json:"bytes-r"`
	DNS                []string   `json:"dns,omitempty"`
	Dot1XMode          string     `json:"dot1x_mode"`
	Dot1XStatus        string     `json:"dot1x_status"`
	Enable             FlexBool   `json:"enable"`
	FlowctrlRx         FlexBool   `json:"flowctrl_rx"`
	FlowctrlTx         FlexBool   `json:"flowctrl_tx"`
	FullDuplex         FlexBool   `json:"full_duplex"`
	IP                 string     `json:"ip,omitempty"`
	Ifname             string     `json:"ifname,omitempty"`
	IsUplink           FlexBool   `json:"is_uplink"`
	Mac                string     `json:"mac,omitempty"`
	MacTable           []MacTable `json:"mac_table,omitempty"`
	Jumbo              FlexBool   `json:"jumbo,omitempty"`
	Masked             FlexBool   `json:"masked"`
	Media              string     `json:"media"`
	Name               string     `json:"name"`
	NetworkName        string     `json:"network_name,omitempty"`
	Netmask            string     `json:"netmask,omitempty"`
	NumPort            int        `json:"num_port,omitempty"`
	OpMode             string     `json:"op_mode"`
	PoeCaps            FlexInt    `json:"poe_caps"`
	PoeClass           string     `json:"poe_class,omitempty"`
	PoeCurrent         FlexInt    `json:"poe_current,omitempty"`
	PoeEnable          FlexBool   `json:"poe_enable,omitempty"`
	PoeGood            FlexBool   `json:"poe_good,omitempty"`
	PoeMode            string     `json:"poe_mode,omitempty"`
	PoePower           FlexInt    `json:"poe_power,omitempty"`
	PoeVoltage         FlexInt    `json:"poe_voltage,omitempty"`
	PortDelta          PortDelta  `json:"port_delta,omitempty"`
	PortIdx            FlexInt    `json:"port_idx"`
	PortPoe            FlexBool   `json:"port_poe"`
	PortconfID         string     `json:"portconf_id"`
	RxBroadcast        FlexInt    `json:"rx_broadcast"`
	RxBytes            FlexInt    `json:"rx_bytes"`
	RxBytesR           FlexInt    `json:"rx_bytes-r"`
	RxDropped          FlexInt    `json:"rx_dropped"`
	RxErrors           FlexInt    `json:"rx_errors"`
	RxMulticast        FlexInt    `json:"rx_multicast"`
	RxPackets          FlexInt    `json:"rx_packets"`
	RxRate             FlexInt    `json:"rx_rate,omitempty"`
	Satisfaction       FlexInt    `json:"satisfaction,omitempty"`
	SatisfactionReason FlexInt    `json:"satisfaction_reason"`
	SFPCompliance      string     `json:"sfp_compliance"`
	SFPCurrent         FlexInt    `json:"sfp_current"`
	SFPFound           FlexBool   `json:"sfp_found"`
	SFPPart            string     `json:"sfp_part"`
	SFPRev             string     `json:"sfp_rev"`
	SFPRxfault         FlexBool   `json:"sfp_rxfault"`
	SFPRxpower         FlexInt    `json:"sfp_rxpower"`
	SFPSerial          string     `json:"sfp_serial"`
	SFPTemperature     FlexInt    `json:"sfp_temperature"`
	SFPTxfault         FlexBool   `json:"sfp_txfault"`
	SFPTxpower         FlexInt    `json:"sfp_txpower"`
	SFPVendor          string     `json:"sfp_vendor"`
	SFPVoltage         FlexInt    `json:"sfp_voltage"`
	Speed              FlexInt    `json:"speed"`
	SpeedCaps          FlexInt    `json:"speed_caps"`
	StpPathcost        FlexInt    `json:"stp_pathcost"`
	StpState           string     `json:"stp_state"`
	TxBroadcast        FlexInt    `json:"tx_broadcast"`
	TxBytes            FlexInt    `json:"tx_bytes"`
	TxBytesR           FlexInt    `json:"tx_bytes-r"`
	TxDropped          FlexInt    `json:"tx_dropped"`
	TxErrors           FlexInt    `json:"tx_errors"`
	TxMulticast        FlexInt    `json:"tx_multicast"`
	TxPackets          FlexInt    `json:"tx_packets"`
	TxRate             FlexInt    `json:"tx_rate,omitempty"`
	Type               string     `json:"type,omitempty"`
	Up                 FlexBool   `json:"up"`
}

// PortDelta is part of a Port.
type PortDelta struct {
	TimeDelta         FlexInt `json:"time_delta"`
	TimeDeltaActivity FlexInt `json:"time_delta_activity"`
}

// USWStat holds the "stat" data for a switch.
// This is split out because of a JSON data format change from 5.10 to 5.11.
type USWStat struct {
	*Sw
}

// Sw is a subtype of USWStat to make unmarshalling of different controller versions possible.
type Sw struct {
	SiteID      string    `json:"site_id"`
	O           string    `json:"o"`
	Oid         string    `json:"oid"`
	Sw          string    `json:"sw"`
	Time        FlexInt   `json:"time"`
	Datetime    time.Time `json:"datetime"`
	RxPackets   FlexInt   `json:"rx_packets"`
	RxBytes     FlexInt   `json:"rx_bytes"`
	RxErrors    FlexInt   `json:"rx_errors"`
	RxDropped   FlexInt   `json:"rx_dropped"`
	RxCrypts    FlexInt   `json:"rx_crypts"`
	RxFrags     FlexInt   `json:"rx_frags"`
	TxPackets   FlexInt   `json:"tx_packets"`
	TxBytes     FlexInt   `json:"tx_bytes"`
	TxErrors    FlexInt   `json:"tx_errors"`
	TxDropped   FlexInt   `json:"tx_dropped"`
	TxRetries   FlexInt   `json:"tx_retries"`
	RxMulticast FlexInt   `json:"rx_multicast"`
	RxBroadcast FlexInt   `json:"rx_broadcast"`
	TxMulticast FlexInt   `json:"tx_multicast"`
	TxBroadcast FlexInt   `json:"tx_broadcast"`
	Bytes       FlexInt   `json:"bytes"`
	Duration    FlexInt   `json:"duration"`
	/* These are all in port table */
	/*
		Port1RxPackets    FlexInt   `json:"port_1-rx_packets,omitempty"`
		Port1RxBytes      FlexInt   `json:"port_1-rx_bytes,omitempty"`
		Port1TxPackets    FlexInt   `json:"port_1-tx_packets,omitempty"`
		Port1TxBytes      FlexInt   `json:"port_1-tx_bytes,omitempty"`
		Port1TxMulticast  FlexInt   `json:"port_1-tx_multicast"`
		Port1TxBroadcast  FlexInt   `json:"port_1-tx_broadcast"`
		Port3RxPackets    FlexInt   `json:"port_3-rx_packets,omitempty"`
		Port3RxBytes      FlexInt   `json:"port_3-rx_bytes,omitempty"`
		Port3TxPackets    FlexInt   `json:"port_3-tx_packets,omitempty"`
		Port3TxBytes      FlexInt   `json:"port_3-tx_bytes,omitempty"`
		Port3RxBroadcast  FlexInt   `json:"port_3-rx_broadcast"`
		Port3TxMulticast  FlexInt   `json:"port_3-tx_multicast"`
		Port3TxBroadcast  FlexInt   `json:"port_3-tx_broadcast"`
		Port6RxPackets    FlexInt   `json:"port_6-rx_packets,omitempty"`
		Port6RxBytes      FlexInt   `json:"port_6-rx_bytes,omitempty"`
		Port6TxPackets    FlexInt   `json:"port_6-tx_packets,omitempty"`
		Port6TxBytes      FlexInt   `json:"port_6-tx_bytes,omitempty"`
		Port6RxMulticast  FlexInt   `json:"port_6-rx_multicast"`
		Port6TxMulticast  FlexInt   `json:"port_6-tx_multicast"`
		Port6TxBroadcast  FlexInt   `json:"port_6-tx_broadcast"`
		Port7RxPackets    FlexInt   `json:"port_7-rx_packets,omitempty"`
		Port7RxBytes      FlexInt   `json:"port_7-rx_bytes,omitempty"`
		Port7TxPackets    FlexInt   `json:"port_7-tx_packets,omitempty"`
		Port7TxBytes      FlexInt   `json:"port_7-tx_bytes,omitempty"`
		Port7TxMulticast  FlexInt   `json:"port_7-tx_multicast"`
		Port7TxBroadcast  FlexInt   `json:"port_7-tx_broadcast"`
		Port9RxPackets    FlexInt   `json:"port_9-rx_packets,omitempty"`
		Port9RxBytes      FlexInt   `json:"port_9-rx_bytes,omitempty"`
		Port9TxPackets    FlexInt   `json:"port_9-tx_packets,omitempty"`
		Port9TxBytes      FlexInt   `json:"port_9-tx_bytes,omitempty"`
		Port9TxMulticast  FlexInt   `json:"port_9-tx_multicast"`
		Port9TxBroadcast  FlexInt   `json:"port_9-tx_broadcast"`
		Port10RxPackets   FlexInt   `json:"port_10-rx_packets,omitempty"`
		Port10RxBytes     FlexInt   `json:"port_10-rx_bytes,omitempty"`
		Port10TxPackets   FlexInt   `json:"port_10-tx_packets,omitempty"`
		Port10TxBytes     FlexInt   `json:"port_10-tx_bytes,omitempty"`
		Port10RxMulticast FlexInt   `json:"port_10-rx_multicast"`
		Port10TxMulticast FlexInt   `json:"port_10-tx_multicast"`
		Port10TxBroadcast FlexInt   `json:"port_10-tx_broadcast"`
		Port11RxPackets   FlexInt   `json:"port_11-rx_packets,omitempty"`
		Port11RxBytes     FlexInt   `json:"port_11-rx_bytes,omitempty"`
		Port11TxPackets   FlexInt   `json:"port_11-tx_packets,omitempty"`
		Port11TxBytes     FlexInt   `json:"port_11-tx_bytes,omitempty"`
		Port11TxMulticast FlexInt   `json:"port_11-tx_multicast"`
		Port11TxBroadcast FlexInt   `json:"port_11-tx_broadcast"`
		Port12RxPackets   FlexInt   `json:"port_12-rx_packets,omitempty"`
		Port12RxBytes     FlexInt   `json:"port_12-rx_bytes,omitempty"`
		Port12TxPackets   FlexInt   `json:"port_12-tx_packets,omitempty"`
		Port12TxBytes     FlexInt   `json:"port_12-tx_bytes,omitempty"`
		Port12TxMulticast FlexInt   `json:"port_12-tx_multicast"`
		Port12TxBroadcast FlexInt   `json:"port_12-tx_broadcast"`
		Port13RxPackets   FlexInt   `json:"port_13-rx_packets,omitempty"`
		Port13RxBytes     FlexInt   `json:"port_13-rx_bytes,omitempty"`
		Port13TxPackets   FlexInt   `json:"port_13-tx_packets,omitempty"`
		Port13TxBytes     FlexInt   `json:"port_13-tx_bytes,omitempty"`
		Port13RxMulticast FlexInt   `json:"port_13-rx_multicast"`
		Port13RxBroadcast FlexInt   `json:"port_13-rx_broadcast"`
		Port13TxMulticast FlexInt   `json:"port_13-tx_multicast"`
		Port13TxBroadcast FlexInt   `json:"port_13-tx_broadcast"`
		Port15RxPackets   FlexInt   `json:"port_15-rx_packets,omitempty"`
		Port15RxBytes     FlexInt   `json:"port_15-rx_bytes,omitempty"`
		Port15TxPackets   FlexInt   `json:"port_15-tx_packets,omitempty"`
		Port15TxBytes     FlexInt   `json:"port_15-tx_bytes,omitempty"`
		Port15RxBroadcast FlexInt   `json:"port_15-rx_broadcast"`
		Port15TxMulticast FlexInt   `json:"port_15-tx_multicast"`
		Port15TxBroadcast FlexInt   `json:"port_15-tx_broadcast"`
		Port16RxPackets   FlexInt   `json:"port_16-rx_packets,omitempty"`
		Port16RxBytes     FlexInt   `json:"port_16-rx_bytes,omitempty"`
		Port16TxPackets   FlexInt   `json:"port_16-tx_packets,omitempty"`
		Port16TxBytes     FlexInt   `json:"port_16-tx_bytes,omitempty"`
		Port16TxMulticast FlexInt   `json:"port_16-tx_multicast"`
		Port16TxBroadcast FlexInt   `json:"port_16-tx_broadcast"`
		Port17RxPackets   FlexInt   `json:"port_17-rx_packets,omitempty"`
		Port17RxBytes     FlexInt   `json:"port_17-rx_bytes,omitempty"`
		Port17TxPackets   FlexInt   `json:"port_17-tx_packets,omitempty"`
		Port17TxBytes     FlexInt   `json:"port_17-tx_bytes,omitempty"`
		Port17TxMulticast FlexInt   `json:"port_17-tx_multicast"`
		Port17TxBroadcast FlexInt   `json:"port_17-tx_broadcast"`
		Port18RxPackets   FlexInt   `json:"port_18-rx_packets,omitempty"`
		Port18RxBytes     FlexInt   `json:"port_18-rx_bytes,omitempty"`
		Port18TxPackets   FlexInt   `json:"port_18-tx_packets,omitempty"`
		Port18TxBytes     FlexInt   `json:"port_18-tx_bytes,omitempty"`
		Port18RxMulticast FlexInt   `json:"port_18-rx_multicast"`
		Port18TxMulticast FlexInt   `json:"port_18-tx_multicast"`
		Port18TxBroadcast FlexInt   `json:"port_18-tx_broadcast"`
		Port19RxPackets   FlexInt   `json:"port_19-rx_packets,omitempty"`
		Port19RxBytes     FlexInt   `json:"port_19-rx_bytes,omitempty"`
		Port19TxPackets   FlexInt   `json:"port_19-tx_packets,omitempty"`
		Port19TxBytes     FlexInt   `json:"port_19-tx_bytes,omitempty"`
		Port19TxMulticast FlexInt   `json:"port_19-tx_multicast"`
		Port19TxBroadcast FlexInt   `json:"port_19-tx_broadcast"`
		Port21RxPackets   FlexInt   `json:"port_21-rx_packets,omitempty"`
		Port21RxBytes     FlexInt   `json:"port_21-rx_bytes,omitempty"`
		Port21TxPackets   FlexInt   `json:"port_21-tx_packets,omitempty"`
		Port21TxBytes     FlexInt   `json:"port_21-tx_bytes,omitempty"`
		Port21RxBroadcast FlexInt   `json:"port_21-rx_broadcast"`
		Port21TxMulticast FlexInt   `json:"port_21-tx_multicast"`
		Port21TxBroadcast FlexInt   `json:"port_21-tx_broadcast"`
		Port22RxPackets   FlexInt   `json:"port_22-rx_packets,omitempty"`
		Port22RxBytes     FlexInt   `json:"port_22-rx_bytes,omitempty"`
		Port22TxPackets   FlexInt   `json:"port_22-tx_packets,omitempty"`
		Port22TxBytes     FlexInt   `json:"port_22-tx_bytes,omitempty"`
		Port22RxMulticast FlexInt   `json:"port_22-rx_multicast"`
		Port22TxMulticast FlexInt   `json:"port_22-tx_multicast"`
		Port22TxBroadcast FlexInt   `json:"port_22-tx_broadcast"`
		Port23RxPackets   FlexInt   `json:"port_23-rx_packets,omitempty"`
		Port23RxBytes     FlexInt   `json:"port_23-rx_bytes,omitempty"`
		Port23RxDropped   FlexInt   `json:"port_23-rx_dropped"`
		Port23TxPackets   FlexInt   `json:"port_23-tx_packets,omitempty"`
		Port23TxBytes     FlexInt   `json:"port_23-tx_bytes,omitempty"`
		Port23RxMulticast FlexInt   `json:"port_23-rx_multicast"`
		Port23RxBroadcast FlexInt   `json:"port_23-rx_broadcast"`
		Port23TxMulticast FlexInt   `json:"port_23-tx_multicast"`
		Port23TxBroadcast FlexInt   `json:"port_23-tx_broadcast"`
		Port24RxPackets   FlexInt   `json:"port_24-rx_packets,omitempty"`
		Port24RxBytes     FlexInt   `json:"port_24-rx_bytes,omitempty"`
		Port24TxPackets   FlexInt   `json:"port_24-tx_packets,omitempty"`
		Port24TxBytes     FlexInt   `json:"port_24-tx_bytes,omitempty"`
		Port24RxMulticast FlexInt   `json:"port_24-rx_multicast"`
		Port24TxMulticast FlexInt   `json:"port_24-tx_multicast"`
		Port24TxBroadcast FlexInt   `json:"port_24-tx_broadcast"`
		Port1RxMulticast  FlexInt   `json:"port_1-rx_multicast"`
		Port3RxDropped    FlexInt   `json:"port_3-rx_dropped"`
		Port3RxMulticast  FlexInt   `json:"port_3-rx_multicast"`
		Port6RxDropped    FlexInt   `json:"port_6-rx_dropped"`
		Port7RxDropped    FlexInt   `json:"port_7-rx_dropped"`
		Port7RxMulticast  FlexInt   `json:"port_7-rx_multicast"`
		Port9RxDropped    FlexInt   `json:"port_9-rx_dropped"`
		Port9RxMulticast  FlexInt   `json:"port_9-rx_multicast"`
		Port9RxBroadcast  FlexInt   `json:"port_9-rx_broadcast"`
		Port10RxBroadcast FlexInt   `json:"port_10-rx_broadcast"`
		Port12RxDropped   FlexInt   `json:"port_12-rx_dropped"`
		Port12RxMulticast FlexInt   `json:"port_12-rx_multicast"`
		Port13RxDropped   FlexInt   `json:"port_13-rx_dropped"`
		Port17RxDropped   FlexInt   `json:"port_17-rx_dropped"`
		Port17RxMulticast FlexInt   `json:"port_17-rx_multicast"`
		Port17RxBroadcast FlexInt   `json:"port_17-rx_broadcast"`
		Port19RxDropped   FlexInt   `json:"port_19-rx_dropped"`
		Port19RxMulticast FlexInt   `json:"port_19-rx_multicast"`
		Port19RxBroadcast FlexInt   `json:"port_19-rx_broadcast"`
		Port21RxDropped   FlexInt   `json:"port_21-rx_dropped"`
		Port21RxMulticast FlexInt   `json:"port_21-rx_multicast"`
		Port7RxBroadcast  FlexInt   `json:"port_7-rx_broadcast"`
		Port18RxBroadcast FlexInt   `json:"port_18-rx_broadcast"`
		Port16RxMulticast FlexInt   `json:"port_16-rx_multicast"`
		Port15RxDropped   FlexInt   `json:"port_15-rx_dropped"`
		Port15RxMulticast FlexInt   `json:"port_15-rx_multicast"`
		Port16RxBroadcast FlexInt   `json:"port_16-rx_broadcast"`
		Port11RxBroadcast FlexInt   `json:"port_11-rx_broadcast"`
		Port12RxBroadcast FlexInt   `json:"port_12-rx_broadcast"`
		Port6RxBroadcast  FlexInt   `json:"port_6-rx_broadcast"`
		Port24RxBroadcast FlexInt   `json:"port_24-rx_broadcast"`
		Port22RxBroadcast FlexInt   `json:"port_22-rx_broadcast"`
		Port10TxDropped   FlexInt   `json:"port_10-tx_dropped"`
		Port16TxDropped   FlexInt   `json:"port_16-tx_dropped"`
		Port1RxBroadcast  FlexInt   `json:"port_1-rx_broadcast"`
		Port4RxPackets    FlexInt   `json:"port_4-rx_packets,omitempty"`
		Port4RxBytes      FlexInt   `json:"port_4-rx_bytes,omitempty"`
		Port4RxDropped    FlexInt   `json:"port_4-rx_dropped"`
		Port4TxPackets    FlexInt   `json:"port_4-tx_packets,omitempty"`
		Port4TxBytes      FlexInt   `json:"port_4-tx_bytes,omitempty"`
		Port4TxDropped    FlexInt   `json:"port_4-tx_dropped"`
		Port4RxMulticast  FlexInt   `json:"port_4-rx_multicast"`
		Port4RxBroadcast  FlexInt   `json:"port_4-rx_broadcast"`
		Port4TxMulticast  FlexInt   `json:"port_4-tx_multicast"`
		Port4TxBroadcast  FlexInt   `json:"port_4-tx_broadcast"`
	*/
}

// UnmarshalJSON unmarshalls 5.10 or 5.11 formatted Switch Stat data.
func (v *USWStat) UnmarshalJSON(data []byte) error {
	var n struct {
		Sw `json:"sw"`
	}

	v.Sw = &n.Sw

	err := json.Unmarshal(data, v.Sw) // controller version 5.10.
	if err != nil {
		return json.Unmarshal(data, &n) // controller version 5.11.
	}

	return nil
}
