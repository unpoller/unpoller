package unifi

import (
	"encoding/json"
	"time"
)

// USG represents all the data from the Ubiquiti Controller for a Unifi Security Gateway.
type USG struct {
	site                  *Site
	SourceName            string               `json:"-"`
	ID                    string               `json:"_id"`
	Adopted               FlexBool             `json:"adopted"`
	Cfgversion            string               `json:"cfgversion"`
	ConfigNetwork         *ConfigNetwork       `json:"config_network"`
	EthernetTable         []*EthernetTable     `json:"ethernet_table"`
	FwCaps                FlexInt              `json:"fw_caps"`
	InformIP              string               `json:"inform_ip"`
	InformURL             string               `json:"inform_url"`
	IP                    string               `json:"ip"`
	LedOverride           string               `json:"led_override"`
	LicenseState          string               `json:"license_state"`
	Mac                   string               `json:"mac"`
	Model                 string               `json:"model"`
	Name                  string               `json:"name"`
	OutdoorModeOverride   string               `json:"outdoor_mode_override"`
	Serial                string               `json:"serial"`
	SiteID                string               `json:"site_id"`
	SiteName              string               `json:"-"`
	Type                  string               `json:"type"`
	UsgCaps               FlexInt              `json:"usg_caps"`
	Version               string               `json:"version"`
	RequiredVersion       string               `json:"required_version"`
	EthernetOverrides     []*EthernetOverrides `json:"ethernet_overrides"`
	HwCaps                FlexInt              `json:"hw_caps"`
	BoardRev              FlexInt              `json:"board_rev"`
	Unsupported           FlexBool             `json:"unsupported"`
	UnsupportedReason     FlexInt              `json:"unsupported_reason"`
	DeviceID              string               `json:"device_id"`
	State                 FlexInt              `json:"state"`
	LastSeen              FlexInt              `json:"last_seen"`
	Upgradable            FlexBool             `json:"upgradable"`
	AdoptableWhenUpgraded FlexBool             `json:"adoptable_when_upgraded"`
	Rollupgrade           FlexBool             `json:"rollupgrade"`
	KnownCfgversion       string               `json:"known_cfgversion"`
	Uptime                FlexInt              `json:"uptime"`
	Locating              FlexBool             `json:"locating"`
	ConnectRequestIP      string               `json:"connect_request_ip"`
	ConnectRequestPort    string               `json:"connect_request_port"`
	SysStats              SysStats             `json:"sys_stats"`
	Temperatures          []Temperature        `json:"temperatures,omitempty"`
	SystemStats           SystemStats          `json:"system-stats"`
	GuestToken            string               `json:"guest_token"`
	SpeedtestStatus       SpeedtestStatus      `json:"speedtest-status"`
	SpeedtestStatusSaved  FlexBool             `json:"speedtest-status-saved"`
	Wan1                  Wan                  `json:"wan1"`
	Wan2                  Wan                  `json:"wan2"`
	PortTable             []*Port              `json:"port_table"`
	NetworkTable          NetworkTable         `json:"network_table"`
	Uplink                Uplink               `json:"uplink"`
	Stat                  USGStat              `json:"stat"`
	TxBytes               FlexInt              `json:"tx_bytes"`
	RxBytes               FlexInt              `json:"rx_bytes"`
	Bytes                 FlexInt              `json:"bytes"`
	NumSta                FlexInt              `json:"num_sta"`
	UserNumSta            FlexInt              `json:"user-num_sta"`
	GuestNumSta           FlexInt              `json:"guest-num_sta"`
	NumDesktop            FlexInt              `json:"num_desktop"`
	NumMobile             FlexInt              `json:"num_mobile"`
	NumHandheld           FlexInt              `json:"num_handheld"`
}

// Uplink is the Internet connection (or uplink) on a UniFi device.
type Uplink struct {
	BytesR           FlexInt  `json:"bytes-r"`
	Drops            FlexInt  `json:"drops"`
	Enable           FlexBool `json:"enable,omitempty"`
	FullDuplex       FlexBool `json:"full_duplex"`
	Gateways         []string `json:"gateways,omitempty"`
	IP               string   `json:"ip"`
	Latency          FlexInt  `json:"latency"`
	Mac              string   `json:"mac,omitempty"`
	MaxSpeed         FlexInt  `json:"max_speed"`
	Name             string   `json:"name"`
	Nameservers      []string `json:"nameservers"`
	Netmask          string   `json:"netmask"`
	NumPort          FlexInt  `json:"num_port"`
	Media            string   `json:"media"`
	PortIdx          FlexInt  `json:"port_idx"`
	RxBytes          FlexInt  `json:"rx_bytes"`
	RxBytesR         FlexInt  `json:"rx_bytes-r"`
	RxDropped        FlexInt  `json:"rx_dropped"`
	RxErrors         FlexInt  `json:"rx_errors"`
	RxMulticast      FlexInt  `json:"rx_multicast"`
	RxPackets        FlexInt  `json:"rx_packets"`
	RxRate           FlexInt  `json:"rx_rate"`
	Speed            FlexInt  `json:"speed"`
	SpeedtestLastrun FlexInt  `json:"speedtest_lastrun,omitempty"`
	SpeedtestPing    FlexInt  `json:"speedtest_ping,omitempty"`
	SpeedtestStatus  string   `json:"speedtest_status,omitempty"`
	TxBytes          FlexInt  `json:"tx_bytes"`
	TxBytesR         FlexInt  `json:"tx_bytes-r"`
	TxDropped        FlexInt  `json:"tx_dropped"`
	TxErrors         FlexInt  `json:"tx_errors"`
	TxPackets        FlexInt  `json:"tx_packets"`
	TxRate           FlexInt  `json:"tx_rate"`
	Type             string   `json:"type"`
	Up               FlexBool `json:"up"`
	Uptime           FlexInt  `json:"uptime"`
	XputDown         FlexInt  `json:"xput_down,omitempty"`
	XputUp           FlexInt  `json:"xput_up,omitempty"`
}

// Wan is a Wan interface on a USG or UDM.
type Wan struct {
	Autoneg     FlexBool `json:"autoneg"`
	BytesR      FlexInt  `json:"bytes-r"`
	DNS         []string `json:"dns"` // may be deprecated
	Enable      FlexBool `json:"enable"`
	FlowctrlRx  FlexBool `json:"flowctrl_rx"`
	FlowctrlTx  FlexBool `json:"flowctrl_tx"`
	FullDuplex  FlexBool `json:"full_duplex"`
	Gateway     string   `json:"gateway"` // may be deprecated
	IP          string   `json:"ip"`
	Ifname      string   `json:"ifname"`
	IsUplink    FlexBool `json:"is_uplink"`
	Mac         string   `json:"mac"`
	MaxSpeed    FlexInt  `json:"max_speed"`
	Media       string   `json:"media"`
	Name        string   `json:"name"`
	Netmask     string   `json:"netmask"` // may be deprecated
	NumPort     int      `json:"num_port"`
	PortIdx     int      `json:"port_idx"`
	PortPoe     FlexBool `json:"port_poe"`
	RxBroadcast FlexInt  `json:"rx_broadcast"`
	RxBytes     FlexInt  `json:"rx_bytes"`
	RxBytesR    FlexInt  `json:"rx_bytes-r"`
	RxDropped   FlexInt  `json:"rx_dropped"`
	RxErrors    FlexInt  `json:"rx_errors"`
	RxMulticast FlexInt  `json:"rx_multicast"`
	RxPackets   FlexInt  `json:"rx_packets"`
	RxRate      FlexInt  `json:"rx_rate"`
	Speed       FlexInt  `json:"speed"`
	SpeedCaps   FlexInt  `json:"speed_caps"`
	TxBroadcast FlexInt  `json:"tx_broadcast"`
	TxBytes     FlexInt  `json:"tx_bytes"`
	TxBytesR    FlexInt  `json:"tx_bytes-r"`
	TxDropped   FlexInt  `json:"tx_dropped"`
	TxErrors    FlexInt  `json:"tx_errors"`
	TxMulticast FlexInt  `json:"tx_multicast"`
	TxPackets   FlexInt  `json:"tx_packets"`
	TxRate      FlexInt  `json:"tx_rate"`
	Type        string   `json:"type"`
	Up          FlexBool `json:"up"`
}

// SpeedtestStatus is the speed test info on a USG or UDM.
type SpeedtestStatus struct {
	Latency         FlexInt          `json:"latency"`
	Rundate         FlexInt          `json:"rundate"`
	Runtime         FlexInt          `json:"runtime"`
	ServerDesc      string           `json:"server_desc,omitempty"`
	Server          *SpeedtestServer `json:"server"`
	SourceInterface string           `json:"source_interface"`
	StatusDownload  FlexInt          `json:"status_download"`
	StatusPing      FlexInt          `json:"status_ping"`
	StatusSummary   FlexInt          `json:"status_summary"`
	StatusUpload    FlexInt          `json:"status_upload"`
	XputDownload    FlexInt          `json:"xput_download"`
	XputUpload      FlexInt          `json:"xput_upload"`
}

type SpeedtestServer struct {
	Cc          string  `json:"cc"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Lat         FlexInt `json:"lat"`
	Lon         FlexInt `json:"lon"`
	Provider    string  `json:"provider"`
	ProviderURL string  `json:"provider_url"`
}

// ConfigNetwork comes from gateways.
type ConfigNetwork struct {
	Type string `json:"type"`
	IP   string `json:"ip"`
}

// SystemStats is system info for a UDM, USG, USW.
type SystemStats struct {
	CPU    FlexInt `json:"cpu"`
	Mem    FlexInt `json:"mem"`
	Uptime FlexInt `json:"uptime"`
	// This exists on at least USG4, may others, maybe not.
	// {"Board (CPU)":"51 C","Board (PHY)":"51 C","CPU":"72 C","PHY":"77 C"}
	Temps map[string]string `json:"temps,omitempty"`
}

// SysStats is load info for a UDM, USG, USW.
type SysStats struct {
	Loadavg1  FlexInt `json:"loadavg_1"`
	Loadavg15 FlexInt `json:"loadavg_15"`
	Loadavg5  FlexInt `json:"loadavg_5"`
	MemBuffer FlexInt `json:"mem_buffer"`
	MemTotal  FlexInt `json:"mem_total"`
	MemUsed   FlexInt `json:"mem_used"`
}

// USGStat holds the "stat" data for a gateway.
// This is split out because of a JSON data format change from 5.10 to 5.11.
type USGStat struct {
	*Gw
}

// Gw is a subtype of USGStat to make unmarshalling of different controller versions possible.
type Gw struct {
	SiteID       string    `json:"site_id"`
	O            string    `json:"o"`
	Oid          string    `json:"oid"`
	Gw           string    `json:"gw"`
	Time         FlexInt   `json:"time"`
	Datetime     time.Time `json:"datetime"`
	Duration     FlexInt   `json:"duration"`
	WanRxPackets FlexInt   `json:"wan-rx_packets"`
	WanRxBytes   FlexInt   `json:"wan-rx_bytes"`
	WanRxDropped FlexInt   `json:"wan-rx_dropped"`
	WanTxPackets FlexInt   `json:"wan-tx_packets"`
	WanTxBytes   FlexInt   `json:"wan-tx_bytes"`
	LanRxPackets FlexInt   `json:"lan-rx_packets"`
	LanRxBytes   FlexInt   `json:"lan-rx_bytes"`
	LanTxPackets FlexInt   `json:"lan-tx_packets"`
	LanTxBytes   FlexInt   `json:"lan-tx_bytes"`
	LanRxDropped FlexInt   `json:"lan-rx_dropped"`
	WanRxErrors  FlexInt   `json:"wan-rx_errors,omitempty"`
	LanRxErrors  FlexInt   `json:"lan-rx_errors,omitempty"`
}

// UnmarshalJSON unmarshalls 5.10 or 5.11 formatted Gateway Stat data.
func (v *USGStat) UnmarshalJSON(data []byte) error {
	var n struct {
		Gw `json:"gw"`
	}

	v.Gw = &n.Gw

	err := json.Unmarshal(data, v.Gw) // controller version 5.10.
	if err != nil {
		return json.Unmarshal(data, &n) // controller version 5.11.
	}

	return nil
}
