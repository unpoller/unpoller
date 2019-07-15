package unifi

import (
	"encoding/json"
	"time"
)

// USG represents all the data from the Ubiquiti Controller for a Unifi Security Gateway.
type USG struct {
	ID            string   `json:"_id"`
	Adopted       FlexBool `json:"adopted"`
	Cfgversion    string   `json:"cfgversion"`
	ConfigNetwork struct {
		Type string `json:"type"`
		IP   string `json:"ip"`
	} `json:"config_network"`
	EthernetTable []struct {
		Mac     string  `json:"mac"`
		NumPort FlexInt `json:"num_port"`
		Name    string  `json:"name"`
	} `json:"ethernet_table"`
	FwCaps              FlexInt `json:"fw_caps"`
	InformIP            string  `json:"inform_ip"`
	InformURL           string  `json:"inform_url"`
	IP                  string  `json:"ip"`
	LedOverride         string  `json:"led_override"`
	LicenseState        string  `json:"license_state"`
	Mac                 string  `json:"mac"`
	Model               string  `json:"model"`
	Name                string  `json:"name"`
	OutdoorModeOverride string  `json:"outdoor_mode_override"`
	Serial              string  `json:"serial"`
	SiteID              string  `json:"site_id"`
	SiteName            string  `json:"-"`
	Type                string  `json:"type"`
	UsgCaps             FlexInt `json:"usg_caps"`
	Version             string  `json:"version"`
	RequiredVersion     string  `json:"required_version"`
	EthernetOverrides   []struct {
		Ifname       string `json:"ifname"`
		Networkgroup string `json:"networkgroup"`
	} `json:"ethernet_overrides"`
	HwCaps                FlexInt  `json:"hw_caps"`
	BoardRev              FlexInt  `json:"board_rev"`
	Unsupported           FlexBool `json:"unsupported"`
	UnsupportedReason     FlexInt  `json:"unsupported_reason"`
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
	GuestToken      string `json:"guest_token"`
	SpeedtestStatus struct {
		Latency        FlexInt `json:"latency"`
		Rundate        FlexInt `json:"rundate"`
		Runtime        FlexInt `json:"runtime"`
		StatusDownload FlexInt `json:"status_download"`
		StatusPing     FlexInt `json:"status_ping"`
		StatusSummary  FlexInt `json:"status_summary"`
		StatusUpload   FlexInt `json:"status_upload"`
		XputDownload   FlexInt `json:"xput_download"`
		XputUpload     FlexInt `json:"xput_upload"`
	} `json:"speedtest-status"`
	SpeedtestStatusSaved FlexBool `json:"speedtest-status-saved"`
	Wan1                 struct {
		TxBytesR    FlexInt  `json:"tx_bytes-r"`
		RxBytesR    FlexInt  `json:"rx_bytes-r"`
		BytesR      FlexInt  `json:"bytes-r"`
		MaxSpeed    FlexInt  `json:"max_speed"`
		Type        string   `json:"type"`
		Name        string   `json:"name"`
		Ifname      string   `json:"ifname"`
		IP          string   `json:"ip"`
		Netmask     string   `json:"netmask"`
		Mac         string   `json:"mac"`
		Up          FlexBool `json:"up"`
		Speed       FlexInt  `json:"speed"`
		FullDuplex  FlexBool `json:"full_duplex"`
		RxBytes     FlexInt  `json:"rx_bytes"`
		RxDropped   FlexInt  `json:"rx_dropped"`
		RxErrors    FlexInt  `json:"rx_errors"`
		RxPackets   FlexInt  `json:"rx_packets"`
		TxBytes     FlexInt  `json:"tx_bytes"`
		TxDropped   FlexInt  `json:"tx_dropped"`
		TxErrors    FlexInt  `json:"tx_errors"`
		TxPackets   FlexInt  `json:"tx_packets"`
		RxMulticast FlexInt  `json:"rx_multicast"`
		Enable      FlexBool `json:"enable"`
		DNS         []string `json:"dns"`
		Gateway     string   `json:"gateway"`
	} `json:"wan1"`
	Wan2 struct {
		TxBytesR    FlexInt  `json:"tx_bytes-r"`
		RxBytesR    FlexInt  `json:"rx_bytes-r"`
		BytesR      FlexInt  `json:"bytes-r"`
		MaxSpeed    FlexInt  `json:"max_speed"`
		Type        string   `json:"type"`
		Name        string   `json:"name"`
		Ifname      string   `json:"ifname"`
		IP          string   `json:"ip"`
		Netmask     string   `json:"netmask"`
		Mac         string   `json:"mac"`
		Up          FlexBool `json:"up"`
		Speed       FlexInt  `json:"speed"`
		FullDuplex  FlexBool `json:"full_duplex"`
		RxBytes     FlexInt  `json:"rx_bytes"`
		RxDropped   FlexInt  `json:"rx_dropped"`
		RxErrors    FlexInt  `json:"rx_errors"`
		RxPackets   FlexInt  `json:"rx_packets"`
		TxBytes     FlexInt  `json:"tx_bytes"`
		TxDropped   FlexInt  `json:"tx_dropped"`
		TxErrors    FlexInt  `json:"tx_errors"`
		TxPackets   FlexInt  `json:"tx_packets"`
		RxMulticast FlexInt  `json:"rx_multicast"`
		Enable      FlexBool `json:"enable"`
		DNS         []string `json:"dns"`
		Gateway     string   `json:"gateway"`
	} `json:"wan2"`
	PortTable []struct {
		Name        string   `json:"name"`
		Ifname      string   `json:"ifname"`
		IP          string   `json:"ip"`
		Netmask     string   `json:"netmask"`
		Mac         string   `json:"mac"`
		Up          FlexBool `json:"up"`
		Speed       FlexInt  `json:"speed"`
		FullDuplex  FlexBool `json:"full_duplex"`
		RxBytes     FlexInt  `json:"rx_bytes"`
		RxDropped   FlexInt  `json:"rx_dropped"`
		RxErrors    FlexInt  `json:"rx_errors"`
		RxPackets   FlexInt  `json:"rx_packets"`
		TxBytes     FlexInt  `json:"tx_bytes"`
		TxDropped   FlexInt  `json:"tx_dropped"`
		TxErrors    FlexInt  `json:"tx_errors"`
		TxPackets   FlexInt  `json:"tx_packets"`
		RxMulticast FlexInt  `json:"rx_multicast"`
		Enable      FlexBool `json:"enable"`
		DNS         []string `json:"dns,omitempty"`
		Gateway     string   `json:"gateway,omitempty"`
	} `json:"port_table"`
	NetworkTable []struct {
		ID                     string   `json:"_id"`
		IsNat                  FlexBool `json:"is_nat"`
		DhcpdDNSEnabled        FlexBool `json:"dhcpd_dns_enabled"`
		Purpose                string   `json:"purpose"`
		DhcpdLeasetime         FlexInt  `json:"dhcpd_leasetime"`
		IgmpSnooping           FlexBool `json:"igmp_snooping"`
		DhcpguardEnabled       FlexBool `json:"dhcpguard_enabled,omitempty"`
		DhcpdStart             string   `json:"dhcpd_start"`
		Enabled                FlexBool `json:"enabled"`
		DhcpdStop              string   `json:"dhcpd_stop"`
		DhcpdWinsEnabled       FlexBool `json:"dhcpd_wins_enabled,omitempty"`
		DomainName             string   `json:"domain_name"`
		DhcpdEnabled           FlexBool `json:"dhcpd_enabled"`
		IPSubnet               string   `json:"ip_subnet"`
		Vlan                   FlexInt  `json:"vlan,omitempty"`
		Networkgroup           string   `json:"networkgroup"`
		Name                   string   `json:"name"`
		SiteID                 string   `json:"site_id"`
		DhcpdIP1               string   `json:"dhcpd_ip_1,omitempty"`
		VlanEnabled            FlexBool `json:"vlan_enabled"`
		DhcpdGatewayEnabled    FlexBool `json:"dhcpd_gateway_enabled"`
		DhcpdTimeOffsetEnabled FlexBool `json:"dhcpd_time_offset_enabled"`
		Ipv6InterfaceType      string   `json:"ipv6_interface_type"`
		DhcpRelayEnabled       FlexBool `json:"dhcp_relay_enabled"`
		Mac                    string   `json:"mac"`
		IsGuest                FlexBool `json:"is_guest"`
		IP                     string   `json:"ip"`
		Up                     FlexBool `json:"up"`
		NumSta                 FlexInt  `json:"num_sta"`
		RxBytes                FlexInt  `json:"rx_bytes"`
		RxPackets              FlexInt  `json:"rx_packets"`
		TxBytes                FlexInt  `json:"tx_bytes"`
		TxPackets              FlexInt  `json:"tx_packets"`
		DhcpdNtp1              string   `json:"dhcpd_ntp_1,omitempty"`
		DhcpdNtpEnabled        FlexBool `json:"dhcpd_ntp_enabled,omitempty"`
		DhcpdUnifiController   string   `json:"dhcpd_unifi_controller,omitempty"`
		UpnpLanEnabled         FlexBool `json:"upnp_lan_enabled,omitempty"`
		AttrNoDelete           FlexBool `json:"attr_no_delete,omitempty"`
		AttrHiddenID           string   `json:"attr_hidden_id,omitempty"`
	} `json:"network_table"`
	Uplink struct {
		Drops            FlexInt  `json:"drops"`
		Enable           FlexBool `json:"enable"`
		FullDuplex       FlexBool `json:"full_duplex"`
		Gateways         []string `json:"gateways"`
		IP               string   `json:"ip"`
		Latency          FlexInt  `json:"latency"`
		Mac              string   `json:"mac"`
		Name             string   `json:"name"`
		Nameservers      []string `json:"nameservers"`
		Netmask          string   `json:"netmask"`
		NumPort          FlexInt  `json:"num_port"`
		RxBytes          FlexInt  `json:"rx_bytes"`
		RxDropped        FlexInt  `json:"rx_dropped"`
		RxErrors         FlexInt  `json:"rx_errors"`
		RxMulticast      FlexInt  `json:"rx_multicast"`
		RxPackets        FlexInt  `json:"rx_packets"`
		Speed            FlexInt  `json:"speed"`
		SpeedtestLastrun FlexInt  `json:"speedtest_lastrun"`
		SpeedtestPing    FlexInt  `json:"speedtest_ping"`
		SpeedtestStatus  string   `json:"speedtest_status"`
		TxBytes          FlexInt  `json:"tx_bytes"`
		TxDropped        FlexInt  `json:"tx_dropped"`
		TxErrors         FlexInt  `json:"tx_errors"`
		TxPackets        FlexInt  `json:"tx_packets"`
		Up               FlexBool `json:"up"`
		Uptime           FlexInt  `json:"uptime"`
		XputDown         FlexInt  `json:"xput_down"`
		XputUp           FlexInt  `json:"xput_up"`
		TxBytesR         FlexInt  `json:"tx_bytes-r"`
		RxBytesR         FlexInt  `json:"rx_bytes-r"`
		BytesR           FlexInt  `json:"bytes-r"`
		MaxSpeed         FlexInt  `json:"max_speed"`
		Type             string   `json:"type"`
	} `json:"uplink"`
	Stat        USGStat `json:"stat"`
	TxBytes     FlexInt `json:"tx_bytes"`
	RxBytes     FlexInt `json:"rx_bytes"`
	Bytes       FlexInt `json:"bytes"`
	NumSta      FlexInt `json:"num_sta"`
	UserNumSta  FlexInt `json:"user-num_sta"`
	GuestNumSta FlexInt `json:"guest-num_sta"`
	NumDesktop  FlexInt `json:"num_desktop"`
	NumMobile   FlexInt `json:"num_mobile"`
	NumHandheld FlexInt `json:"num_handheld"`
}

// USGStat holds the "stat" data for a gateway.
// This is split out because of a JSON data format change from 5.10 to 5.11.
type USGStat struct {
	*gw
}

type gw struct {
	SiteID       string    `json:"site_id"`
	O            string    `json:"o"`
	Oid          string    `json:"oid"`
	Gw           string    `json:"gw"`
	Time         FlexInt   `json:"time"`
	Datetime     time.Time `json:"datetime"`
	Duration     FlexInt   `json:"duration"`
	WanRxPackets FlexInt   `json:"wan-rx_packets"`
	WanRxBytes   FlexInt   `json:"wan-rx_bytes"`
	WanTxPackets FlexInt   `json:"wan-tx_packets"`
	WanTxBytes   FlexInt   `json:"wan-tx_bytes"`
	LanRxPackets FlexInt   `json:"lan-rx_packets"`
	LanRxBytes   FlexInt   `json:"lan-rx_bytes"`
	LanTxPackets FlexInt   `json:"lan-tx_packets"`
	LanTxBytes   FlexInt   `json:"lan-tx_bytes"`
	WanRxDropped FlexInt   `json:"wan-rx_dropped"`
	LanRxDropped FlexInt   `json:"lan-rx_dropped"`
}

// UnmarshalJSON unmarshalls 5.10 or 5.11 formatted Gateway Stat data.
func (v *USGStat) UnmarshalJSON(data []byte) error {
	var n struct {
		gw `json:"gw"`
	}
	v.gw = &n.gw
	err := json.Unmarshal(data, v.gw) // controller version 5.10.
	if err != nil {
		return json.Unmarshal(data, &n) // controller version 5.11.
	}
	return nil
}
