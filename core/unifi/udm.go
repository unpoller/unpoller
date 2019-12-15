package unifi

// UDM represents all the data from the Ubiquiti Controller for a Unifi Dream Machine.
// The UDM shares several structs/type-data with USW and USG.
type UDM struct {
	SourceName            string   `json:"-"`
	SiteID                string   `json:"site_id"`
	SiteName              string   `json:"-"`
	Mac                   string   `json:"mac"`
	Adopted               FlexBool `json:"adopted"`
	Serial                string   `json:"serial"`
	IP                    string   `json:"ip"`
	Uptime                FlexInt  `json:"uptime"`
	Model                 string   `json:"model"`
	Version               string   `json:"version"`
	Name                  string   `json:"name"`
	Default               FlexBool `json:"default"`
	Locating              FlexBool `json:"locating"`
	Type                  string   `json:"type"`
	Unsupported           FlexBool `json:"unsupported"`
	UnsupportedReason     FlexInt  `json:"unsupported_reason"`
	DiscoveredVia         string   `json:"discovered_via"`
	AdoptIP               string   `json:"adopt_ip"`
	AdoptURL              string   `json:"adopt_url"`
	State                 FlexInt  `json:"state"`
	AdoptStatus           FlexInt  `json:"adopt_status"`
	UpgradeState          FlexInt  `json:"upgrade_state"`
	LastSeen              FlexInt  `json:"last_seen"`
	AdoptableWhenUpgraded FlexBool `json:"adoptable_when_upgraded"`
	Cfgversion            string   `json:"cfgversion"`
	ConfigNetwork         struct {
		Type string `json:"type"`
		IP   string `json:"ip"`
	} `json:"config_network"`
	VwireTable             []interface{} `json:"vwire_table"`
	Dot1XPortctrlEnabled   FlexBool      `json:"dot1x_portctrl_enabled"`
	JumboframeEnabled      FlexBool      `json:"jumboframe_enabled"`
	FlowctrlEnabled        FlexBool      `json:"flowctrl_enabled"`
	StpVersion             string        `json:"stp_version"`
	StpPriority            FlexInt       `json:"stp_priority"`
	PowerSourceCtrlEnabled FlexBool      `json:"power_source_ctrl_enabled"`
	LicenseState           string        `json:"license_state"`
	ID                     string        `json:"_id"`
	DeviceID               string        `json:"device_id"`
	AdoptState             FlexInt       `json:"adopt_state"`
	AdoptTries             FlexInt       `json:"adopt_tries"`
	AdoptManual            FlexBool      `json:"adopt_manual"`
	InformURL              string        `json:"inform_url"`
	InformIP               string        `json:"inform_ip"`
	RequiredVersion        string        `json:"required_version"`
	BoardRev               FlexInt       `json:"board_rev"`
	EthernetTable          []struct {
		Mac     string  `json:"mac"`
		NumPort FlexInt `json:"num_port"`
		Name    string  `json:"name"`
	} `json:"ethernet_table"`
	PortTable         []Port `json:"port_table"`
	EthernetOverrides []struct {
		Ifname       string `json:"ifname"`
		Networkgroup string `json:"networkgroup"`
	} `json:"ethernet_overrides"`
	UsgCaps    FlexInt  `json:"usg_caps"`
	HasSpeaker FlexBool `json:"has_speaker"`
	HasEth1    FlexBool `json:"has_eth1"`
	FwCaps     FlexInt  `json:"fw_caps"`
	HwCaps     FlexInt  `json:"hw_caps"`
	WifiCaps   FlexInt  `json:"wifi_caps"`
	SwitchCaps struct {
		MaxMirrorSessions    FlexInt `json:"max_mirror_sessions"`
		MaxAggregateSessions FlexInt `json:"max_aggregate_sessions"`
	} `json:"switch_caps"`
	HasFan            FlexBool    `json:"has_fan"`
	HasTemperature    FlexBool    `json:"has_temperature"`
	RulesetInterfaces interface{} `json:"ruleset_interfaces"`
	/* struct {
		Br0  string `json:"br0"`
		Eth0 string `json:"eth0"`
		Eth1 string `json:"eth1"`
		Eth2 string `json:"eth2"`
		Eth3 string `json:"eth3"`
		Eth4 string `json:"eth4"`
		Eth5 string `json:"eth5"`
		Eth6 string `json:"eth6"`
		Eth7 string `json:"eth7"`
		Eth8 string `json:"eth8"`
	} */
	KnownCfgversion      string           `json:"known_cfgversion"`
	SysStats             SysStats         `json:"sys_stats"`
	SystemStats          SystemStats      `json:"system-stats"`
	GuestToken           string           `json:"guest_token"`
	Overheating          FlexBool         `json:"overheating"`
	SpeedtestStatus      SpeedtestStatus  `json:"speedtest-status"`
	SpeedtestStatusSaved FlexBool         `json:"speedtest-status-saved"`
	Wan1                 Wan              `json:"wan1"`
	Wan2                 Wan              `json:"wan2"`
	Uplink               Uplink           `json:"uplink"`
	ConnectRequestIP     string           `json:"connect_request_ip"`
	ConnectRequestPort   string           `json:"connect_request_port"`
	DownlinkTable        []interface{}    `json:"downlink_table"`
	WlangroupIDNa        string           `json:"wlangroup_id_na"`
	WlangroupIDNg        string           `json:"wlangroup_id_ng"`
	BandsteeringMode     string           `json:"bandsteering_mode"`
	RadioTable           *RadioTable      `json:"radio_table,omitempty"`
	RadioTableStats      *RadioTableStats `json:"radio_table_stats,omitempty"`
	VapTable             *VapTable        `json:"vap_table"`
	XInformAuthkey       string           `json:"x_inform_authkey"`
	NetworkTable         NetworkTable     `json:"network_table"`
	PortOverrides        []struct {
		PortIdx    FlexInt `json:"port_idx"`
		PortconfID string  `json:"portconf_id"`
	} `json:"port_overrides"`
	Stat            UDMStat `json:"stat"`
	TxBytes         FlexInt `json:"tx_bytes"`
	RxBytes         FlexInt `json:"rx_bytes"`
	Bytes           FlexInt `json:"bytes"`
	BytesD          FlexInt `json:"bytes-d"`
	TxBytesD        FlexInt `json:"tx_bytes-d"`
	RxBytesD        FlexInt `json:"rx_bytes-d"`
	BytesR          FlexInt `json:"bytes-r"`
	NumSta          FlexInt `json:"num_sta"`            // USG
	WlanNumSta      FlexInt `json:"wlan-num_sta"`       // UAP
	LanNumSta       FlexInt `json:"lan-num_sta"`        // USW
	UserWlanNumSta  FlexInt `json:"user-wlan-num_sta"`  // UAP
	UserLanNumSta   FlexInt `json:"user-lan-num_sta"`   // USW
	UserNumSta      FlexInt `json:"user-num_sta"`       // USG
	GuestWlanNumSta FlexInt `json:"guest-wlan-num_sta"` // UAP
	GuestLanNumSta  FlexInt `json:"guest-lan-num_sta"`  // USW
	GuestNumSta     FlexInt `json:"guest-num_sta"`      // USG
	NumDesktop      FlexInt `json:"num_desktop"`        // USG
	NumMobile       FlexInt `json:"num_mobile"`         // USG
	NumHandheld     FlexInt `json:"num_handheld"`       // USG
}

// NetworkTable is the list of networks on a gateway.
type NetworkTable []struct {
	ID                     string        `json:"_id"`
	AttrNoDelete           FlexBool      `json:"attr_no_delete"`
	AttrHiddenID           string        `json:"attr_hidden_id"`
	Name                   string        `json:"name"`
	SiteID                 string        `json:"site_id"`
	VlanEnabled            FlexBool      `json:"vlan_enabled"`
	Purpose                string        `json:"purpose"`
	IPSubnet               string        `json:"ip_subnet"`
	Ipv6InterfaceType      string        `json:"ipv6_interface_type"`
	DomainName             string        `json:"domain_name"`
	IsNat                  FlexBool      `json:"is_nat"`
	DhcpdEnabled           FlexBool      `json:"dhcpd_enabled"`
	DhcpdStart             string        `json:"dhcpd_start"`
	DhcpdStop              string        `json:"dhcpd_stop"`
	Dhcpdv6Enabled         FlexBool      `json:"dhcpdv6_enabled"`
	Ipv6RaEnabled          FlexBool      `json:"ipv6_ra_enabled"`
	LteLanEnabled          FlexBool      `json:"lte_lan_enabled"`
	Networkgroup           string        `json:"networkgroup"`
	DhcpdLeasetime         FlexInt       `json:"dhcpd_leasetime"`
	DhcpdDNSEnabled        FlexBool      `json:"dhcpd_dns_enabled"`
	DhcpdGatewayEnabled    FlexBool      `json:"dhcpd_gateway_enabled"`
	DhcpdTimeOffsetEnabled FlexBool      `json:"dhcpd_time_offset_enabled"`
	Ipv6PdStart            string        `json:"ipv6_pd_start"`
	Ipv6PdStop             string        `json:"ipv6_pd_stop"`
	DhcpdDNS1              string        `json:"dhcpd_dns_1"`
	DhcpdDNS2              string        `json:"dhcpd_dns_2"`
	DhcpdDNS3              string        `json:"dhcpd_dns_3"`
	DhcpdDNS4              string        `json:"dhcpd_dns_4"`
	Enabled                FlexBool      `json:"enabled"`
	DhcpRelayEnabled       FlexBool      `json:"dhcp_relay_enabled"`
	Mac                    string        `json:"mac"`
	IsGuest                FlexBool      `json:"is_guest"`
	IP                     string        `json:"ip"`
	Up                     FlexBool      `json:"up"`
	DPIStatsTable          DPIStatsTable `json:"dpistats_table,omitempty"`
	NumSta                 FlexInt       `json:"num_sta"`
	RxBytes                FlexInt       `json:"rx_bytes"`
	RxPackets              FlexInt       `json:"rx_packets"`
	TxBytes                FlexInt       `json:"tx_bytes"`
	TxPackets              FlexInt       `json:"tx_packets"`
}

// DPIStatsTable is the Deep Packet Inspection data for each "network"
type DPIStatsTable struct {
	LastUpdated FlexInt `json:"last_updated"`
	ByCat       []struct {
		Cat       FlexInt   `json:"cat"`
		Apps      []FlexInt `json:"apps"`
		RxBytes   FlexInt   `json:"rx_bytes"`
		TxBytes   FlexInt   `json:"tx_bytes"`
		RxPackets FlexInt   `json:"rx_packets"`
		TxPackets FlexInt   `json:"tx_packets"`
	} `json:"by_cat"`
	ByApp []struct {
		App     FlexInt `json:"app"`
		Cat     FlexInt `json:"cat"`
		Clients []struct {
			Mac       string  `json:"mac"`
			RxBytes   FlexInt `json:"rx_bytes"`
			TxBytes   FlexInt `json:"tx_bytes"`
			RxPackets FlexInt `json:"rx_packets"`
			TxPackets FlexInt `json:"tx_packets"`
		} `json:"clients"`
		KnownClients FlexInt `json:"known_clients"`
		RxBytes      FlexInt `json:"rx_bytes"`
		TxBytes      FlexInt `json:"tx_bytes"`
		RxPackets    FlexInt `json:"rx_packets"`
		TxPackets    FlexInt `json:"tx_packets"`
	} `json:"by_app"`
}

// UDMStat holds the "stat" data for a dream machine.
// A dream machine is a USG + USW + Controller
type UDMStat struct {
	*Gw `json:"gw"`
	*Sw `json:"sw"`
	*Ap `json:"ap,omitempty"`
}
