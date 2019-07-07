package unifi

import (
	"encoding/json"
	"time"
)

// USW represents all the data from the Ubiquiti Controller for a Unifi Switch.
type USW struct {
	SiteName      string   `json:"-"`
	ID            string   `json:"_id"`
	Adopted       FlexBool `json:"adopted"`
	BoardRev      FlexInt  `json:"board_rev"`
	Cfgversion    string   `json:"cfgversion"`
	ConfigNetwork struct {
		Type string `json:"type"`
		IP   string `json:"ip"`
	} `json:"config_network"`
	Dot1XPortctrlEnabled FlexBool `json:"dot1x_portctrl_enabled"`
	EthernetTable        []struct {
		Mac     string  `json:"mac"`
		NumPort FlexInt `json:"num_port,omitempty"`
		Name    string  `json:"name"`
	} `json:"ethernet_table"`
	FlowctrlEnabled     FlexBool `json:"flowctrl_enabled"`
	FwCaps              FlexInt  `json:"fw_caps"`
	HasFan              FlexBool `json:"has_fan"`
	HasTemperature      FlexBool `json:"has_temperature"`
	InformIP            string   `json:"inform_ip"`
	InformURL           string   `json:"inform_url"`
	IP                  string   `json:"ip"`
	JumboframeEnabled   FlexBool `json:"jumboframe_enabled"`
	LedOverride         string   `json:"led_override"`
	LicenseState        string   `json:"license_state"`
	Mac                 string   `json:"mac"`
	Model               string   `json:"model"`
	Name                string   `json:"name"`
	OutdoorModeOverride string   `json:"outdoor_mode_override"`
	PortOverrides       []struct {
		Name       string  `json:"name,omitempty"`
		PoeMode    string  `json:"poe_mode,omitempty"`
		PortIdx    FlexInt `json:"port_idx"`
		PortconfID string  `json:"portconf_id"`
	} `json:"port_overrides"`
	PortTable []struct {
		PortIdx      FlexInt  `json:"port_idx"`
		Media        string   `json:"media"`
		PortPoe      FlexBool `json:"port_poe"`
		PoeCaps      FlexInt  `json:"poe_caps"`
		SpeedCaps    FlexInt  `json:"speed_caps"`
		OpMode       string   `json:"op_mode"`
		PortconfID   string   `json:"portconf_id"`
		PoeMode      string   `json:"poe_mode,omitempty"`
		Autoneg      FlexBool `json:"autoneg"`
		Dot1XMode    string   `json:"dot1x_mode"`
		Dot1XStatus  string   `json:"dot1x_status"`
		Enable       FlexBool `json:"enable"`
		FlowctrlRx   FlexBool `json:"flowctrl_rx"`
		FlowctrlTx   FlexBool `json:"flowctrl_tx"`
		FullDuplex   FlexBool `json:"full_duplex"`
		IsUplink     FlexBool `json:"is_uplink"`
		Jumbo        FlexBool `json:"jumbo"`
		PoeClass     string   `json:"poe_class,omitempty"`
		PoeCurrent   FlexInt  `json:"poe_current,omitempty"`
		PoeEnable    FlexBool `json:"poe_enable,omitempty"`
		PoeGood      FlexBool `json:"poe_good,omitempty"`
		PoePower     FlexInt  `json:"poe_power,omitempty"`
		PoeVoltage   FlexInt  `json:"poe_voltage,omitempty"`
		RxBroadcast  FlexInt  `json:"rx_broadcast"`
		RxBytes      FlexInt  `json:"rx_bytes"`
		RxDropped    FlexInt  `json:"rx_dropped"`
		RxErrors     FlexInt  `json:"rx_errors"`
		RxMulticast  FlexInt  `json:"rx_multicast"`
		RxPackets    FlexInt  `json:"rx_packets"`
		Satisfaction FlexInt  `json:"satisfaction"`
		Speed        FlexInt  `json:"speed"`
		StpPathcost  FlexInt  `json:"stp_pathcost"`
		StpState     string   `json:"stp_state"`
		TxBroadcast  FlexInt  `json:"tx_broadcast"`
		TxBytes      FlexInt  `json:"tx_bytes"`
		TxDropped    FlexInt  `json:"tx_dropped"`
		TxErrors     FlexInt  `json:"tx_errors"`
		TxMulticast  FlexInt  `json:"tx_multicast"`
		TxPackets    FlexInt  `json:"tx_packets"`
		Up           FlexBool `json:"up"`
		TxBytesR     FlexInt  `json:"tx_bytes-r"`
		RxBytesR     FlexInt  `json:"rx_bytes-r"`
		BytesR       FlexInt  `json:"bytes-r"`
		Name         string   `json:"name"`
		Masked       FlexBool `json:"masked"`
		AggregatedBy FlexBool `json:"aggregated_by"`
		SfpFound     FlexBool `json:"sfp_found,omitempty"`
	} `json:"port_table"`
	Serial          string `json:"serial"`
	SiteID          string `json:"site_id"`
	StpPriority     string `json:"stp_priority"`
	StpVersion      string `json:"stp_version"`
	Type            string `json:"type"`
	Version         string `json:"version"`
	RequiredVersion string `json:"required_version"`
	SwitchCaps      struct {
		FeatureCaps          FlexInt `json:"feature_caps"`
		MaxMirrorSessions    FlexInt `json:"max_mirror_sessions"`
		MaxAggregateSessions FlexInt `json:"max_aggregate_sessions"`
	} `json:"switch_caps"`
	HwCaps                FlexInt  `json:"hw_caps"`
	Unsupported           FlexBool `json:"unsupported"`
	UnsupportedReason     FlexInt  `json:"unsupported_reason"`
	SysErrorCaps          FlexInt  `json:"sys_error_caps"`
	DeviceID              string   `json:"device_id"`
	State                 FlexInt  `json:"state"`
	LastSeen              FlexInt  `json:"last_seen"`
	Upgradable            FlexBool `json:"upgradable"`
	AdoptableWhenUpgraded FlexBool `json:"adoptable_when_upgraded"`
	Rollupgrade           FlexBool `json:"rollupgrade"`
	KnownCfgversion       string   `json:"known_cfgversion"`
	Uptime                FlexInt  `json:"uptime"`
	Locating              FlexBool `json:"locating"`
	ConnectRequestIP      string   `json:"connect_request_ip"`
	ConnectRequestPort    string   `json:"connect_request_port"`
	SysStats              struct {
		Loadavg1  FlexInt `json:"loadavg_1"`
		Loadavg15 FlexInt `json:"loadavg_15"`
		Loadavg5  FlexInt `json:"loadavg_5"`
		MemBuffer FlexInt `json:"mem_buffer"`
		MemTotal  FlexInt `json:"mem_total"`
		MemUsed   FlexInt `json:"mem_used"`
	} `json:"sys_stats"`
	SystemStats struct {
		CPU    FlexInt `json:"cpu"`
		Mem    FlexInt `json:"mem"`
		Uptime FlexInt `json:"uptime"`
	} `json:"system-stats"`
	FanLevel           FlexInt  `json:"fan_level"`
	GeneralTemperature FlexInt  `json:"general_temperature"`
	Overheating        FlexBool `json:"overheating"`
	TotalMaxPower      FlexInt  `json:"total_max_power"`
	DownlinkTable      []struct {
		PortIdx    FlexInt  `json:"port_idx"`
		Speed      FlexInt  `json:"speed"`
		FullDuplex FlexBool `json:"full_duplex"`
		Mac        string   `json:"mac"`
	} `json:"downlink_table"`
	Uplink struct {
		FullDuplex  FlexBool `json:"full_duplex"`
		IP          string   `json:"ip"`
		Mac         string   `json:"mac"`
		Name        string   `json:"name"`
		Netmask     string   `json:"netmask"`
		NumPort     FlexInt  `json:"num_port"`
		RxBytes     FlexInt  `json:"rx_bytes"`
		RxDropped   FlexInt  `json:"rx_dropped"`
		RxErrors    FlexInt  `json:"rx_errors"`
		RxMulticast FlexInt  `json:"rx_multicast"`
		RxPackets   FlexInt  `json:"rx_packets"`
		Speed       FlexInt  `json:"speed"`
		TxBytes     FlexInt  `json:"tx_bytes"`
		TxDropped   FlexInt  `json:"tx_dropped"`
		TxErrors    FlexInt  `json:"tx_errors"`
		TxPackets   FlexInt  `json:"tx_packets"`
		Up          FlexBool `json:"up"`
		PortIdx     FlexInt  `json:"port_idx"`
		Media       string   `json:"media"`
		MaxSpeed    FlexInt  `json:"max_speed"`
		UplinkMac   string   `json:"uplink_mac"`
		Type        string   `json:"type"`
		TxBytesR    FlexInt  `json:"tx_bytes-r"`
		RxBytesR    FlexInt  `json:"rx_bytes-r"`
	} `json:"uplink"`
	LastUplink struct {
		UplinkMac string `json:"uplink_mac"`
	} `json:"last_uplink"`
	UplinkDepth FlexInt  `json:"uplink_depth"`
	Stat        *USWStat `json:"stat"`
	TxBytes     FlexInt  `json:"tx_bytes"`
	RxBytes     FlexInt  `json:"rx_bytes"`
	Bytes       FlexInt  `json:"bytes"`
	NumSta      FlexInt  `json:"num_sta"`
	UserNumSta  FlexInt  `json:"user-num_sta"`
	GuestNumSta FlexInt  `json:"guest-num_sta"`
}

// USWStat holds the "stat" data for a switch.
// This is split out because of a JSON data format change from 5.10 to 5.11.
type USWStat struct {
	*sw
}

type sw struct {
	SiteID            string    `json:"site_id"`
	O                 string    `json:"o"`
	Oid               string    `json:"oid"`
	Sw                string    `json:"sw"`
	Time              FlexInt   `json:"time"`
	Datetime          time.Time `json:"datetime"`
	RxPackets         FlexInt   `json:"rx_packets"`
	RxBytes           FlexInt   `json:"rx_bytes"`
	RxErrors          FlexInt   `json:"rx_errors"`
	RxDropped         FlexInt   `json:"rx_dropped"`
	RxCrypts          FlexInt   `json:"rx_crypts"`
	RxFrags           FlexInt   `json:"rx_frags"`
	TxPackets         FlexInt   `json:"tx_packets"`
	TxBytes           FlexInt   `json:"tx_bytes"`
	TxErrors          FlexInt   `json:"tx_errors"`
	TxDropped         FlexInt   `json:"tx_dropped"`
	TxRetries         FlexInt   `json:"tx_retries"`
	RxMulticast       FlexInt   `json:"rx_multicast"`
	RxBroadcast       FlexInt   `json:"rx_broadcast"`
	TxMulticast       FlexInt   `json:"tx_multicast"`
	TxBroadcast       FlexInt   `json:"tx_broadcast"`
	Bytes             FlexInt   `json:"bytes"`
	Duration          FlexInt   `json:"duration"`
	Port1RxPackets    FlexInt   `json:"port_1-rx_packets"`
	Port1RxBytes      FlexInt   `json:"port_1-rx_bytes"`
	Port1TxPackets    FlexInt   `json:"port_1-tx_packets"`
	Port1TxBytes      FlexInt   `json:"port_1-tx_bytes"`
	Port1TxMulticast  FlexInt   `json:"port_1-tx_multicast"`
	Port1TxBroadcast  FlexInt   `json:"port_1-tx_broadcast"`
	Port3RxPackets    FlexInt   `json:"port_3-rx_packets"`
	Port3RxBytes      FlexInt   `json:"port_3-rx_bytes"`
	Port3TxPackets    FlexInt   `json:"port_3-tx_packets"`
	Port3TxBytes      FlexInt   `json:"port_3-tx_bytes"`
	Port3RxBroadcast  FlexInt   `json:"port_3-rx_broadcast"`
	Port3TxMulticast  FlexInt   `json:"port_3-tx_multicast"`
	Port3TxBroadcast  FlexInt   `json:"port_3-tx_broadcast"`
	Port6RxPackets    FlexInt   `json:"port_6-rx_packets"`
	Port6RxBytes      FlexInt   `json:"port_6-rx_bytes"`
	Port6TxPackets    FlexInt   `json:"port_6-tx_packets"`
	Port6TxBytes      FlexInt   `json:"port_6-tx_bytes"`
	Port6RxMulticast  FlexInt   `json:"port_6-rx_multicast"`
	Port6TxMulticast  FlexInt   `json:"port_6-tx_multicast"`
	Port6TxBroadcast  FlexInt   `json:"port_6-tx_broadcast"`
	Port7RxPackets    FlexInt   `json:"port_7-rx_packets"`
	Port7RxBytes      FlexInt   `json:"port_7-rx_bytes"`
	Port7TxPackets    FlexInt   `json:"port_7-tx_packets"`
	Port7TxBytes      FlexInt   `json:"port_7-tx_bytes"`
	Port7TxMulticast  FlexInt   `json:"port_7-tx_multicast"`
	Port7TxBroadcast  FlexInt   `json:"port_7-tx_broadcast"`
	Port9RxPackets    FlexInt   `json:"port_9-rx_packets"`
	Port9RxBytes      FlexInt   `json:"port_9-rx_bytes"`
	Port9TxPackets    FlexInt   `json:"port_9-tx_packets"`
	Port9TxBytes      FlexInt   `json:"port_9-tx_bytes"`
	Port9TxMulticast  FlexInt   `json:"port_9-tx_multicast"`
	Port9TxBroadcast  FlexInt   `json:"port_9-tx_broadcast"`
	Port10RxPackets   FlexInt   `json:"port_10-rx_packets"`
	Port10RxBytes     FlexInt   `json:"port_10-rx_bytes"`
	Port10TxPackets   FlexInt   `json:"port_10-tx_packets"`
	Port10TxBytes     FlexInt   `json:"port_10-tx_bytes"`
	Port10RxMulticast FlexInt   `json:"port_10-rx_multicast"`
	Port10TxMulticast FlexInt   `json:"port_10-tx_multicast"`
	Port10TxBroadcast FlexInt   `json:"port_10-tx_broadcast"`
	Port11RxPackets   FlexInt   `json:"port_11-rx_packets"`
	Port11RxBytes     FlexInt   `json:"port_11-rx_bytes"`
	Port11TxPackets   FlexInt   `json:"port_11-tx_packets"`
	Port11TxBytes     FlexInt   `json:"port_11-tx_bytes"`
	Port11TxMulticast FlexInt   `json:"port_11-tx_multicast"`
	Port11TxBroadcast FlexInt   `json:"port_11-tx_broadcast"`
	Port12RxPackets   FlexInt   `json:"port_12-rx_packets"`
	Port12RxBytes     FlexInt   `json:"port_12-rx_bytes"`
	Port12TxPackets   FlexInt   `json:"port_12-tx_packets"`
	Port12TxBytes     FlexInt   `json:"port_12-tx_bytes"`
	Port12TxMulticast FlexInt   `json:"port_12-tx_multicast"`
	Port12TxBroadcast FlexInt   `json:"port_12-tx_broadcast"`
	Port13RxPackets   FlexInt   `json:"port_13-rx_packets"`
	Port13RxBytes     FlexInt   `json:"port_13-rx_bytes"`
	Port13TxPackets   FlexInt   `json:"port_13-tx_packets"`
	Port13TxBytes     FlexInt   `json:"port_13-tx_bytes"`
	Port13RxMulticast FlexInt   `json:"port_13-rx_multicast"`
	Port13RxBroadcast FlexInt   `json:"port_13-rx_broadcast"`
	Port13TxMulticast FlexInt   `json:"port_13-tx_multicast"`
	Port13TxBroadcast FlexInt   `json:"port_13-tx_broadcast"`
	Port15RxPackets   FlexInt   `json:"port_15-rx_packets"`
	Port15RxBytes     FlexInt   `json:"port_15-rx_bytes"`
	Port15TxPackets   FlexInt   `json:"port_15-tx_packets"`
	Port15TxBytes     FlexInt   `json:"port_15-tx_bytes"`
	Port15RxBroadcast FlexInt   `json:"port_15-rx_broadcast"`
	Port15TxMulticast FlexInt   `json:"port_15-tx_multicast"`
	Port15TxBroadcast FlexInt   `json:"port_15-tx_broadcast"`
	Port16RxPackets   FlexInt   `json:"port_16-rx_packets"`
	Port16RxBytes     FlexInt   `json:"port_16-rx_bytes"`
	Port16TxPackets   FlexInt   `json:"port_16-tx_packets"`
	Port16TxBytes     FlexInt   `json:"port_16-tx_bytes"`
	Port16TxMulticast FlexInt   `json:"port_16-tx_multicast"`
	Port16TxBroadcast FlexInt   `json:"port_16-tx_broadcast"`
	Port17RxPackets   FlexInt   `json:"port_17-rx_packets"`
	Port17RxBytes     FlexInt   `json:"port_17-rx_bytes"`
	Port17TxPackets   FlexInt   `json:"port_17-tx_packets"`
	Port17TxBytes     FlexInt   `json:"port_17-tx_bytes"`
	Port17TxMulticast FlexInt   `json:"port_17-tx_multicast"`
	Port17TxBroadcast FlexInt   `json:"port_17-tx_broadcast"`
	Port18RxPackets   FlexInt   `json:"port_18-rx_packets"`
	Port18RxBytes     FlexInt   `json:"port_18-rx_bytes"`
	Port18TxPackets   FlexInt   `json:"port_18-tx_packets"`
	Port18TxBytes     FlexInt   `json:"port_18-tx_bytes"`
	Port18RxMulticast FlexInt   `json:"port_18-rx_multicast"`
	Port18TxMulticast FlexInt   `json:"port_18-tx_multicast"`
	Port18TxBroadcast FlexInt   `json:"port_18-tx_broadcast"`
	Port19RxPackets   FlexInt   `json:"port_19-rx_packets"`
	Port19RxBytes     FlexInt   `json:"port_19-rx_bytes"`
	Port19TxPackets   FlexInt   `json:"port_19-tx_packets"`
	Port19TxBytes     FlexInt   `json:"port_19-tx_bytes"`
	Port19TxMulticast FlexInt   `json:"port_19-tx_multicast"`
	Port19TxBroadcast FlexInt   `json:"port_19-tx_broadcast"`
	Port21RxPackets   FlexInt   `json:"port_21-rx_packets"`
	Port21RxBytes     FlexInt   `json:"port_21-rx_bytes"`
	Port21TxPackets   FlexInt   `json:"port_21-tx_packets"`
	Port21TxBytes     FlexInt   `json:"port_21-tx_bytes"`
	Port21RxBroadcast FlexInt   `json:"port_21-rx_broadcast"`
	Port21TxMulticast FlexInt   `json:"port_21-tx_multicast"`
	Port21TxBroadcast FlexInt   `json:"port_21-tx_broadcast"`
	Port22RxPackets   FlexInt   `json:"port_22-rx_packets"`
	Port22RxBytes     FlexInt   `json:"port_22-rx_bytes"`
	Port22TxPackets   FlexInt   `json:"port_22-tx_packets"`
	Port22TxBytes     FlexInt   `json:"port_22-tx_bytes"`
	Port22RxMulticast FlexInt   `json:"port_22-rx_multicast"`
	Port22TxMulticast FlexInt   `json:"port_22-tx_multicast"`
	Port22TxBroadcast FlexInt   `json:"port_22-tx_broadcast"`
	Port23RxPackets   FlexInt   `json:"port_23-rx_packets"`
	Port23RxBytes     FlexInt   `json:"port_23-rx_bytes"`
	Port23RxDropped   FlexInt   `json:"port_23-rx_dropped"`
	Port23TxPackets   FlexInt   `json:"port_23-tx_packets"`
	Port23TxBytes     FlexInt   `json:"port_23-tx_bytes"`
	Port23RxMulticast FlexInt   `json:"port_23-rx_multicast"`
	Port23RxBroadcast FlexInt   `json:"port_23-rx_broadcast"`
	Port23TxMulticast FlexInt   `json:"port_23-tx_multicast"`
	Port23TxBroadcast FlexInt   `json:"port_23-tx_broadcast"`
	Port24RxPackets   FlexInt   `json:"port_24-rx_packets"`
	Port24RxBytes     FlexInt   `json:"port_24-rx_bytes"`
	Port24TxPackets   FlexInt   `json:"port_24-tx_packets"`
	Port24TxBytes     FlexInt   `json:"port_24-tx_bytes"`
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
	Port4RxPackets    FlexInt   `json:"port_4-rx_packets"`
	Port4RxBytes      FlexInt   `json:"port_4-rx_bytes"`
	Port4RxDropped    FlexInt   `json:"port_4-rx_dropped"`
	Port4TxPackets    FlexInt   `json:"port_4-tx_packets"`
	Port4TxBytes      FlexInt   `json:"port_4-tx_bytes"`
	Port4TxDropped    FlexInt   `json:"port_4-tx_dropped"`
	Port4RxMulticast  FlexInt   `json:"port_4-rx_multicast"`
	Port4RxBroadcast  FlexInt   `json:"port_4-rx_broadcast"`
	Port4TxMulticast  FlexInt   `json:"port_4-tx_multicast"`
	Port4TxBroadcast  FlexInt   `json:"port_4-tx_broadcast"`
}

// UnmarshalJSON unmarshalls 5.10 or 5.11 formatted Switch Stat data.
func (v *USWStat) UnmarshalJSON(data []byte) error {
	var n struct {
		sw `json:"sw"`
	}
	v.sw = &n.sw
	err := json.Unmarshal(data, v.sw) // controller version 5.10.
	if err != nil {
		return json.Unmarshal(data, &n) // controller version 5.11.
	}
	return nil
}
