package unifi

// UXG represents all the data from the Ubiquiti Controller for a UniFi 10Gb Gateway.
// The UDM shares several structs/type-data with USW and USG.
type UXG struct {
	site                       *Site
	SourceName                 string                  `json:"-"`
	SiteName                   string                  `json:"-"`
	ID                         string                  `json:"_id"`
	IP                         string                  `json:"ip"`
	Mac                        string                  `json:"mac"`
	Model                      string                  `json:"model"`
	ModelInLts                 FlexBool                `json:"model_in_lts"`
	ModelInEol                 FlexBool                `json:"model_in_eol"`
	Type                       string                  `json:"type"`
	Version                    string                  `json:"version"`
	Adopted                    FlexBool                `json:"adopted"`
	SiteID                     string                  `json:"site_id"`
	Cfgversion                 string                  `json:"cfgversion"`
	SyslogKey                  string                  `json:"syslog_key"`
	ConfigNetwork              *ConfigNetwork          `json:"config_network"`
	SetupID                    string                  `json:"setup_id"`
	LicenseState               string                  `json:"license_state"`
	ConfigNetworkLan           *ConfigNetworkLan       `json:"config_network_lan"`
	InformURL                  string                  `json:"inform_url"`
	InformIP                   string                  `json:"inform_ip"`
	RequiredVersion            string                  `json:"required_version"`
	KernelVersion              string                  `json:"kernel_version"`
	Architecture               string                  `json:"architecture"`
	BoardRev                   FlexInt                 `json:"board_rev"`
	ManufacturerID             FlexInt                 `json:"manufacturer_id"`
	Internet                   FlexBool                `json:"internet"`
	ModelIncompatible          FlexBool                `json:"model_incompatible"`
	EthernetTable              []*EthernetTable        `json:"ethernet_table"`
	PortTable                  []Port                  `json:"port_table"`
	EthernetOverrides          []*EthernetOverrides    `json:"ethernet_overrides"`
	UsgCaps                    FlexInt                 `json:"usg_caps"`
	HasSpeaker                 FlexBool                `json:"has_speaker"`
	HasEth1                    FlexBool                `json:"has_eth1"`
	FwCaps                     FlexInt                 `json:"fw_caps"`
	HwCaps                     FlexInt                 `json:"hw_caps"`
	WifiCaps                   FlexInt                 `json:"wifi_caps"`
	SwitchCaps                 *SwitchCaps             `json:"switch_caps"`
	HasFan                     FlexBool                `json:"has_fan"`
	HasTemperature             FlexBool                `json:"has_temperature"`
	Temperatures               []Temperature           `json:"temperatures"`
	Storage                    []*Storage              `json:"storage"`
	RulesetInterfaces          interface{}             `json:"ruleset_interfaces"`
	ConnectedAt                FlexInt                 `json:"connected_at"`
	ProvisionedAt              FlexInt                 `json:"provisioned_at"`
	LedOverride                string                  `json:"led_override"`
	LedOverrideColor           string                  `json:"led_override_color"`
	LedOverrideColorBrightness FlexInt                 `json:"led_override_color_brightness"`
	OutdoorModeOverride        string                  `json:"outdoor_mode_override"`
	LcmBrightnessOverride      FlexBool                `json:"lcm_brightness_override"`
	LcmIdleTimeoutOverride     FlexBool                `json:"lcm_idle_timeout_override"`
	Name                       string                  `json:"name"`
	Unsupported                FlexBool                `json:"unsupported"`
	UnsupportedReason          FlexInt                 `json:"unsupported_reason"`
	Serial                     string                  `json:"serial"`
	HashID                     string                  `json:"hash_id"`
	TwoPhaseAdopt              FlexBool                `json:"two_phase_adopt"`
	DeviceID                   string                  `json:"device_id"`
	State                      FlexInt                 `json:"state"`
	StartDisconnectedMillis    FlexInt                 `json:"start_disconnected_millis"`
	UpgradeState               FlexInt                 `json:"upgrade_state"`
	StartConnectedMillis       FlexInt                 `json:"start_connected_millis"`
	LastSeen                   FlexInt                 `json:"last_seen"`
	Uptime                     FlexInt                 `json:"uptime"`
	UnderscoreUptime           FlexInt                 `json:"_uptime"`
	Locating                   FlexBool                `json:"locating"`
	SysStats                   SysStats                `json:"sys_stats"`
	SystemStats                SystemStats             `json:"system-stats"`
	GuestKicks                 FlexInt                 `json:"guest_kicks"`
	GuestToken                 string                  `json:"guest_token"`
	UptimeStats                map[string]*UptimeStats `json:"uptime_stats"`
	Overheating                FlexBool                `json:"overheating"`
	GeoInfo                    map[string]*GeoInfo     `json:"geo_info"`
	LedState                   *LedState               `json:"led_state"`
	SpeedtestStatus            SpeedtestStatus         `json:"speedtest-status"`
	SpeedtestStatusSaved       FlexBool                `json:"speedtest-status-saved"`
	Wan1                       Wan                     `json:"wan1"`
	Wan2                       Wan                     `json:"wan2"`
	Uplink                     Uplink                  `json:"uplink"`
	DownlinkTable              []*DownlinkTable        `json:"downlink_table"`
	NetworkTable               NetworkTable            `json:"network_table"`
	KnownCfgversion            string                  `json:"known_cfgversion"`
	ConnectRequestIP           string                  `json:"connect_request_ip"`
	ConnectRequestPort         string                  `json:"connect_request_port"`
	NextInterval               FlexInt                 `json:"next_interval"`
	NextHeartbeatAt            FlexInt                 `json:"next_heartbeat_at"`
	ConsideredLostAt           FlexInt                 `json:"considered_lost_at"`
	Stat                       *UXGStat                `json:"stat"`
	TxBytes                    FlexInt                 `json:"tx_bytes"`
	RxBytes                    FlexInt                 `json:"rx_bytes"`
	Bytes                      FlexInt                 `json:"bytes"`
	NumSta                     FlexInt                 `json:"num_sta"`
	WlanNumSta                 FlexInt                 `json:"wlan-num_sta"`
	LanNumSta                  FlexInt                 `json:"lan-num_sta"`
	UserWlanNumSta             FlexInt                 `json:"user-wlan-num_sta"`
	UserLanNumSta              FlexInt                 `json:"user-lan-num_sta"`
	UserNumSta                 FlexInt                 `json:"user-num_sta"`
	GuestWlanNumSta            FlexInt                 `json:"guest-wlan-num_sta"`
	GuestLanNumSta             FlexInt                 `json:"guest-lan-num_sta"`
	GuestNumSta                FlexInt                 `json:"guest-num_sta"`
	NumDesktop                 FlexInt                 `json:"num_desktop"`
	NumMobile                  FlexInt                 `json:"num_mobile"`
	NumHandheld                FlexInt                 `json:"num_handheld"`
}

// ConfigNetworkLan is part of a UXG, maybe others.
type ConfigNetworkLan struct {
	DhcpEnabled FlexBool `json:"dhcp_enabled"`
	Vlan        int      `json:"vlan"`
}

// DownlinkTable is part of a UXG and UDM output.
type DownlinkTable struct {
	PortIdx    FlexInt  `json:"port_idx"`
	Speed      FlexInt  `json:"speed"`
	FullDuplex FlexBool `json:"full_duplex"`
	Mac        string   `json:"mac"`
}

// LedState is incuded with newer devices.
type LedState struct {
	Pattern string  `json:"pattern"`
	Tempo   FlexInt `json:"tempo"`
}

// GeoInfo is incuded with certain devices.
type GeoInfo struct {
	Accuracy        FlexInt `json:"accuracy"`
	Address         string  `json:"address"`
	Asn             FlexInt `json:"asn"`
	City            string  `json:"city"`
	ContinentCode   string  `json:"continent_code"`
	CountryCode     string  `json:"country_code"`
	CountryName     string  `json:"country_name"`
	IspName         string  `json:"isp_name"`
	IspOrganization string  `json:"isp_organization"`
	Latitude        FlexInt `json:"latitude"`
	Longitude       FlexInt `json:"longitude"`
	Timezone        string  `json:"timezone"`
}

// UptimeStats is incuded with certain devices.
type UptimeStats struct {
	Availability   FlexInt `json:"availability"`
	LatencyAverage FlexInt `json:"latency_average"`
	TimePeriod     FlexInt `json:"time_period"`
}

// UXGStat holds the "stat" data for a 10Gb gateway.
type UXGStat struct {
	*Gw `json:"gw"`
	*Sw `json:"sw"`
}
