package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

// DeviceResponse represents the payload from Unifi Controller.
type DeviceResponse struct {
	Devices []Device `json:"data"`
	Meta    struct {
		Rc string `json:"rc"`
	} `json:"meta"`
}

// Device represents all the information a device may contain. oh my..
type Device struct {
	// created using https://mholt.github.io/json-to-go
	// with data from https://unifi.ctrlr/api/s/default/stat/device
	ID           string `json:"_id"`
	UUptime      int    `json:"_uptime"`
	AdoptIP      string `json:"adopt_ip,omitempty"`
	AdoptURL     string `json:"adopt_url,omitempty"`
	Adopted      bool   `json:"adopted"`
	AntennaTable []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Wifi0Gain int    `json:"wifi0_gain"`
		Wifi1Gain int    `json:"wifi1_gain"`
	} `json:"antenna_table,omitempty"`
	BandsteeringMode string  `json:"bandsteering_mode,omitempty"`
	BoardRev         int     `json:"board_rev,omitempty"`
	Bytes            float64 `json:"bytes"`
	BytesD           int     `json:"bytes-d,omitempty"`
	BytesR           int     `json:"bytes-r,omitempty"`
	Cfgversion       string  `json:"cfgversion"`
	ConfigNetwork    struct {
		IP   string `json:"ip"`
		Type string `json:"type"`
	} `json:"config_network"`
	ConnectRequestIP   string        `json:"connect_request_ip"`
	ConnectRequestPort string        `json:"connect_request_port"`
	ConsideredLostAt   int           `json:"considered_lost_at"`
	CountrycodeTable   []int         `json:"countrycode_table,omitempty"`
	Default            bool          `json:"default,omitempty"`
	DeviceID           string        `json:"device_id"`
	DiscoveredVia      string        `json:"discovered_via,omitempty"`
	DownlinkTable      []interface{} `json:"downlink_table,omitempty"`
	EthernetTable      []struct {
		Mac     string `json:"mac"`
		Name    string `json:"name"`
		NumPort int    `json:"num_port"`
	} `json:"ethernet_table"`
	FwCaps          int    `json:"fw_caps"`
	GuestNumSta     int    `json:"guest-num_sta"`
	GuestToken      string `json:"guest_token,omitempty"`
	HasEth1         bool   `json:"has_eth1,omitempty"`
	HasSpeaker      bool   `json:"has_speaker,omitempty"`
	InformIP        string `json:"inform_ip"`
	InformURL       string `json:"inform_url"`
	IP              string `json:"ip"`
	Isolated        bool   `json:"isolated,omitempty"`
	KnownCfgversion string `json:"known_cfgversion"`
	LastSeen        int    `json:"last_seen"`
	LastUplink      struct {
		UplinkMac        string `json:"uplink_mac"`
		UplinkRemotePort int    `json:"uplink_remote_port"`
	} `json:"last_uplink,omitempty"`
	LedOverride         string `json:"led_override"`
	Locating            bool   `json:"locating"`
	Mac                 string `json:"mac"`
	Model               string `json:"model"`
	Name                string `json:"name"`
	NextHeartbeatAt     int    `json:"next_heartbeat_at"`
	NumSta              int    `json:"num_sta"`
	OutdoorModeOverride string `json:"outdoor_mode_override"`
	PortTable           []struct {
		AggregatedBy bool `json:"aggregated_by"`
		AttrNoEdit   bool `json:"attr_no_edit,omitempty"`
		Autoneg      bool `json:"autoneg"`
		BytesR       int  `json:"bytes-r"`
		Enable       bool `json:"enable"`
		FlowctrlRx   bool `json:"flowctrl_rx"`
		FlowctrlTx   bool `json:"flowctrl_tx"`
		FullDuplex   bool `json:"full_duplex"`
		IsUplink     bool `json:"is_uplink"`
		Jumbo        bool `json:"jumbo"`
		MacTable     []struct {
			Age    int    `json:"age"`
			Mac    string `json:"mac"`
			Static bool   `json:"static"`
			Uptime int    `json:"uptime"`
			Vlan   int    `json:"vlan"`
		} `json:"mac_table"`
		Masked    bool   `json:"masked"`
		Media     string `json:"media"`
		Name      string `json:"name"`
		OpMode    string `json:"op_mode"`
		PoeCaps   int    `json:"poe_caps"`
		PortDelta struct {
			RxBytes   int `json:"rx_bytes"`
			RxPackets int `json:"rx_packets"`
			TimeDelta int `json:"time_delta"`
			TxBytes   int `json:"tx_bytes"`
			TxPackets int `json:"tx_packets"`
		} `json:"port_delta"`
		PortIdx     int    `json:"port_idx"`
		PortPoe     bool   `json:"port_poe"`
		PortconfID  string `json:"portconf_id"`
		RxBroadcast int    `json:"rx_broadcast"`
		RxBytes     int64  `json:"rx_bytes"`
		RxBytesR    int    `json:"rx_bytes-r"`
		RxDropped   int    `json:"rx_dropped"`
		RxErrors    int    `json:"rx_errors"`
		RxMulticast int    `json:"rx_multicast"`
		RxPackets   int    `json:"rx_packets"`
		Speed       int    `json:"speed"`
		StpPathcost int    `json:"stp_pathcost"`
		StpState    string `json:"stp_state"`
		TxBroadcast int    `json:"tx_broadcast"`
		TxBytes     int64  `json:"tx_bytes"`
		TxBytesR    int    `json:"tx_bytes-r"`
		TxDropped   int    `json:"tx_dropped"`
		TxErrors    int    `json:"tx_errors"`
		TxMulticast int    `json:"tx_multicast"`
		TxPackets   int    `json:"tx_packets"`
		Up          bool   `json:"up"`
	} `json:"port_table"`
	RadioTable []struct {
		BuiltinAntGain     int    `json:"builtin_ant_gain"`
		BuiltinAntenna     bool   `json:"builtin_antenna"`
		Channel            string `json:"channel"`
		CurrentAntennaGain int    `json:"current_antenna_gain"`
		Ht                 string `json:"ht"`
		MaxTxpower         int    `json:"max_txpower"`
		MinRssiEnabled     bool   `json:"min_rssi_enabled"`
		MinTxpower         int    `json:"min_txpower"`
		Name               string `json:"name"`
		Nss                int    `json:"nss"`
		Radio              string `json:"radio"`
		RadioCaps          int    `json:"radio_caps"`
		TxPower            string `json:"tx_power"`
		TxPowerMode        string `json:"tx_power_mode"`
		WlangroupID        string `json:"wlangroup_id"`
		HasDfs             bool   `json:"has_dfs,omitempty"`
		HasFccdfs          bool   `json:"has_fccdfs,omitempty"`
		Is11Ac             bool   `json:"is_11ac,omitempty"`
	} `json:"radio_table,omitempty"`
	RadioTableStats []struct {
		AstBeXmit   interface{} `json:"ast_be_xmit"`
		AstCst      interface{} `json:"ast_cst"`
		AstTxto     interface{} `json:"ast_txto"`
		Channel     int         `json:"channel"`
		CuSelfRx    int         `json:"cu_self_rx"`
		CuSelfTx    int         `json:"cu_self_tx"`
		CuTotal     int         `json:"cu_total"`
		Extchannel  int         `json:"extchannel"`
		Gain        int         `json:"gain"`
		GuestNumSta int         `json:"guest-num_sta"`
		Name        string      `json:"name"`
		NumSta      int         `json:"num_sta"`
		Radio       string      `json:"radio"`
		State       string      `json:"state"`
		TxPackets   float64     `json:"tx_packets"`
		TxPower     int         `json:"tx_power"`
		TxRetries   int         `json:"tx_retries"`
		UserNumSta  int         `json:"user-num_sta"`
	} `json:"radio_table_stats,omitempty"`
	Rollupgrade      bool          `json:"rollupgrade"`
	RxBytes          float64       `json:"rx_bytes"`
	RxBytesD         float64       `json:"rx_bytes-d,omitempty"`
	ScanRadioTable   []interface{} `json:"scan_radio_table,omitempty"`
	Scanning         bool          `json:"scanning,omitempty"`
	Serial           string        `json:"serial"`
	SiteID           string        `json:"site_id"`
	SpectrumScanning bool          `json:"spectrum_scanning,omitempty"`
	SSHSessionTable  []interface{} `json:"ssh_session_table,omitempty"`
	Stat             struct {
		Ap                 string  `json:"ap"`
		Bytes              float64 `json:"bytes"`
		Datetime           string  `json:"datetime"`
		Duration           float64 `json:"duration"`
		GuestRxBytes       float64 `json:"guest-rx_bytes"`
		GuestRxCrypts      float64 `json:"guest-rx_crypts"`
		GuestRxDropped     float64 `json:"guest-rx_dropped"`
		GuestRxErrors      float64 `json:"guest-rx_errors"`
		GuestRxFrags       float64 `json:"guest-rx_frags"`
		GuestRxPackets     float64 `json:"guest-rx_packets"`
		GuestTxBytes       float64 `json:"guest-tx_bytes"`
		GuestTxDropped     float64 `json:"guest-tx_dropped"`
		GuestTxErrors      float64 `json:"guest-tx_errors"`
		GuestTxPackets     float64 `json:"guest-tx_packets"`
		GuestTxRetries     float64 `json:"guest-tx_retries"`
		O                  string  `json:"o"`
		Oid                string  `json:"oid"`
		Port1RxBroadcast   float64 `json:"port_1-rx_broadcast"`
		Port1RxBytes       float64 `json:"port_1-rx_bytes"`
		Port1RxMulticast   float64 `json:"port_1-rx_multicast"`
		Port1RxPackets     float64 `json:"port_1-rx_packets"`
		Port1TxBroadcast   float64 `json:"port_1-tx_broadcast"`
		Port1TxBytes       float64 `json:"port_1-tx_bytes"`
		Port1TxMulticast   float64 `json:"port_1-tx_multicast"`
		Port1TxPackets     float64 `json:"port_1-tx_packets"`
		RxBytes            float64 `json:"rx_bytes"`
		RxCrypts           float64 `json:"rx_crypts"`
		RxDropped          float64 `json:"rx_dropped"`
		RxErrors           float64 `json:"rx_errors"`
		RxFrags            float64 `json:"rx_frags"`
		RxPackets          float64 `json:"rx_packets"`
		SiteID             string  `json:"site_id"`
		Time               int64   `json:"time"`
		TxBytes            float64 `json:"tx_bytes"`
		TxDropped          float64 `json:"tx_dropped"`
		TxErrors           float64 `json:"tx_errors"`
		TxPackets          float64 `json:"tx_packets"`
		TxRetries          float64 `json:"tx_retries"`
		UserRxBytes        float64 `json:"user-rx_bytes"`
		UserRxCrypts       float64 `json:"user-rx_crypts"`
		UserRxDropped      float64 `json:"user-rx_dropped"`
		UserRxErrors       float64 `json:"user-rx_errors"`
		UserRxFrags        float64 `json:"user-rx_frags"`
		UserRxPackets      float64 `json:"user-rx_packets"`
		UserTxBytes        float64 `json:"user-tx_bytes"`
		UserTxDropped      float64 `json:"user-tx_dropped"`
		UserTxErrors       float64 `json:"user-tx_errors"`
		UserTxPackets      float64 `json:"user-tx_packets"`
		UserTxRetries      float64 `json:"user-tx_retries"`
		UserWifi0RxBytes   float64 `json:"user-wifi0-rx_bytes"`
		UserWifi0RxCrypts  float64 `json:"user-wifi0-rx_crypts"`
		UserWifi0RxDropped float64 `json:"user-wifi0-rx_dropped"`
		UserWifi0RxErrors  float64 `json:"user-wifi0-rx_errors"`
		UserWifi0RxFrags   float64 `json:"user-wifi0-rx_frags"`
		UserWifi0RxPackets float64 `json:"user-wifi0-rx_packets"`
		UserWifi0TxBytes   float64 `json:"user-wifi0-tx_bytes"`
		UserWifi0TxDropped float64 `json:"user-wifi0-tx_dropped"`
		UserWifi0TxErrors  float64 `json:"user-wifi0-tx_errors"`
		UserWifi0TxPackets float64 `json:"user-wifi0-tx_packets"`
		UserWifi0TxRetries float64 `json:"user-wifi0-tx_retries"`
		UserWifi1RxBytes   float64 `json:"user-wifi1-rx_bytes"`
		UserWifi1RxCrypts  float64 `json:"user-wifi1-rx_crypts"`
		UserWifi1RxDropped float64 `json:"user-wifi1-rx_dropped"`
		UserWifi1RxErrors  float64 `json:"user-wifi1-rx_errors"`
		UserWifi1RxFrags   float64 `json:"user-wifi1-rx_frags"`
		UserWifi1RxPackets float64 `json:"user-wifi1-rx_packets"`
		UserWifi1TxBytes   float64 `json:"user-wifi1-tx_bytes"`
		UserWifi1TxDropped float64 `json:"user-wifi1-tx_dropped"`
		UserWifi1TxErrors  float64 `json:"user-wifi1-tx_errors"`
		UserWifi1TxPackets float64 `json:"user-wifi1-tx_packets"`
		UserWifi1TxRetries float64 `json:"user-wifi1-tx_retries"`
		Wifi0RxBytes       float64 `json:"wifi0-rx_bytes"`
		Wifi0RxCrypts      float64 `json:"wifi0-rx_crypts"`
		Wifi0RxDropped     float64 `json:"wifi0-rx_dropped"`
		Wifi0RxErrors      float64 `json:"wifi0-rx_errors"`
		Wifi0RxFrags       float64 `json:"wifi0-rx_frags"`
		Wifi0RxPackets     float64 `json:"wifi0-rx_packets"`
		Wifi0TxBytes       float64 `json:"wifi0-tx_bytes"`
		Wifi0TxDropped     float64 `json:"wifi0-tx_dropped"`
		Wifi0TxErrors      float64 `json:"wifi0-tx_errors"`
		Wifi0TxPackets     float64 `json:"wifi0-tx_packets"`
		Wifi0TxRetries     float64 `json:"wifi0-tx_retries"`
		Wifi1RxBytes       float64 `json:"wifi1-rx_bytes"`
		Wifi1RxCrypts      float64 `json:"wifi1-rx_crypts"`
		Wifi1RxDropped     float64 `json:"wifi1-rx_dropped"`
		Wifi1RxErrors      float64 `json:"wifi1-rx_errors"`
		Wifi1RxFrags       float64 `json:"wifi1-rx_frags"`
		Wifi1RxPackets     float64 `json:"wifi1-rx_packets"`
		Wifi1TxBytes       float64 `json:"wifi1-tx_bytes"`
		Wifi1TxDropped     float64 `json:"wifi1-tx_dropped"`
		Wifi1TxErrors      float64 `json:"wifi1-tx_errors"`
		Wifi1TxPackets     float64 `json:"wifi1-tx_packets"`
		Wifi1TxRetries     float64 `json:"wifi1-tx_retries"`
	} `json:"stat"`
	State    int `json:"state"`
	SysStats struct {
		Loadavg1  string `json:"loadavg_1"`
		Loadavg15 string `json:"loadavg_15"`
		Loadavg5  string `json:"loadavg_5"`
		MemBuffer int    `json:"mem_buffer"`
		MemTotal  int    `json:"mem_total"`
		MemUsed   int    `json:"mem_used"`
	} `json:"sys_stats"`
	SystemStats struct {
		CPU    string `json:"cpu"`
		Mem    string `json:"mem"`
		Uptime string `json:"uptime"`
	} `json:"system-stats"`
	TxBytes    int64  `json:"tx_bytes"`
	TxBytesD   int    `json:"tx_bytes-d,omitempty"`
	Type       string `json:"type"`
	Upgradable bool   `json:"upgradable"`
	Uplink     struct {
		FullDuplex       bool   `json:"full_duplex"`
		IP               string `json:"ip"`
		Mac              string `json:"mac"`
		MaxSpeed         int    `json:"max_speed"`
		MaxVlan          int    `json:"max_vlan"`
		Media            string `json:"media"`
		Name             string `json:"name"`
		Netmask          string `json:"netmask"`
		NumPort          int    `json:"num_port"`
		RxBytes          int    `json:"rx_bytes"`
		RxBytesR         int    `json:"rx_bytes-r"`
		RxDropped        int    `json:"rx_dropped"`
		RxErrors         int    `json:"rx_errors"`
		RxMulticast      int    `json:"rx_multicast"`
		RxPackets        int    `json:"rx_packets"`
		Speed            int    `json:"speed"`
		TxBytes          int64  `json:"tx_bytes"`
		TxBytesR         int    `json:"tx_bytes-r"`
		TxDropped        int    `json:"tx_dropped"`
		TxErrors         int    `json:"tx_errors"`
		TxPackets        int    `json:"tx_packets"`
		Type             string `json:"type"`
		Up               bool   `json:"up"`
		UplinkMac        string `json:"uplink_mac"`
		UplinkRemotePort int    `json:"uplink_remote_port"`
	} `json:"uplink"`
	UplinkTable []interface{} `json:"uplink_table,omitempty"`
	Uptime      int           `json:"uptime"`
	UserNumSta  int           `json:"user-num_sta"`
	VapTable    []struct {
		ApMac               string      `json:"ap_mac"`
		Bssid               string      `json:"bssid"`
		Ccq                 int         `json:"ccq"`
		Channel             int         `json:"channel"`
		Essid               string      `json:"essid"`
		Extchannel          int         `json:"extchannel"`
		ID                  string      `json:"id"`
		IsGuest             bool        `json:"is_guest"`
		IsWep               bool        `json:"is_wep"`
		MacFilterRejections int         `json:"mac_filter_rejections"`
		MapID               interface{} `json:"map_id"`
		Name                string      `json:"name"`
		NumSta              int         `json:"num_sta"`
		Radio               string      `json:"radio"`
		RadioName           string      `json:"radio_name"`
		RxBytes             int         `json:"rx_bytes"`
		RxCrypts            int         `json:"rx_crypts"`
		RxDropped           int         `json:"rx_dropped"`
		RxErrors            int         `json:"rx_errors"`
		RxFrags             int         `json:"rx_frags"`
		RxNwids             int         `json:"rx_nwids"`
		RxPackets           int         `json:"rx_packets"`
		SiteID              string      `json:"site_id"`
		State               string      `json:"state"`
		T                   string      `json:"t"`
		TxBytes             int         `json:"tx_bytes"`
		TxDropped           int         `json:"tx_dropped"`
		TxErrors            int         `json:"tx_errors"`
		TxLatencyAvg        int         `json:"tx_latency_avg"`
		TxLatencyMax        int         `json:"tx_latency_max"`
		TxLatencyMin        int         `json:"tx_latency_min"`
		TxPackets           int         `json:"tx_packets"`
		TxPower             int         `json:"tx_power"`
		TxRetries           int         `json:"tx_retries"`
		Up                  bool        `json:"up"`
		Usage               string      `json:"usage"`
		WlanconfID          string      `json:"wlanconf_id"`
	} `json:"vap_table,omitempty"`
	Version             string        `json:"version"`
	VersionIncompatible bool          `json:"version_incompatible"`
	VwireEnabled        bool          `json:"vwireEnabled,omitempty"`
	VwireTable          []interface{} `json:"vwire_table,omitempty"`
	VwireVapTable       []struct {
		Bssid     string `json:"bssid"`
		Radio     string `json:"radio"`
		RadioName string `json:"radio_name"`
		State     string `json:"state"`
	} `json:"vwire_vap_table,omitempty"`
	WifiCaps             int           `json:"wifi_caps,omitempty"`
	DhcpServerTable      []interface{} `json:"dhcp_server_table,omitempty"`
	Dot1XPortctrlEnabled bool          `json:"dot1x_portctrl_enabled,omitempty"`
	FanLevel             int           `json:"fan_level,omitempty"`
	FlowctrlEnabled      bool          `json:"flowctrl_enabled,omitempty"`
	GeneralTemperature   int           `json:"general_temperature,omitempty"`
	HasFan               bool          `json:"has_fan,omitempty"`
	HasTemperature       bool          `json:"has_temperature,omitempty"`
	JumboframeEnabled    bool          `json:"jumboframe_enabled,omitempty"`
	LicenseState         string        `json:"license_state,omitempty"`
	Overheating          bool          `json:"overheating,omitempty"`
	PortOverrides        []struct {
		Name       string `json:"name,omitempty"`
		PoeMode    string `json:"poe_mode,omitempty"`
		PortIdx    int    `json:"port_idx"`
		PortconfID string `json:"portconf_id"`
	} `json:"port_overrides,omitempty"`
	StpPriority      string `json:"stp_priority,omitempty"`
	StpVersion       string `json:"stp_version,omitempty"`
	UplinkDepth      int    `json:"uplink_depth,omitempty"`
	ConfigNetworkWan struct {
		Type string `json:"type"`
	} `json:"config_network_wan,omitempty"`
	NetworkTable []struct {
		ID                     string      `json:"_id"`
		DhcpdDNSEnabled        bool        `json:"dhcpd_dns_enabled"`
		DhcpdEnabled           bool        `json:"dhcpd_enabled"`
		DhcpdIP1               string      `json:"dhcpd_ip_1,omitempty"`
		DhcpdLeasetime         json.Number `json:"dhcpd_leasetime,Number"`
		DhcpdStart             string      `json:"dhcpd_start"`
		DhcpdStop              string      `json:"dhcpd_stop"`
		DhcpdWinsEnabled       bool        `json:"dhcpd_wins_enabled,omitempty"`
		DhcpguardEnabled       bool        `json:"dhcpguard_enabled,omitempty"`
		DomainName             string      `json:"domain_name"`
		Enabled                bool        `json:"enabled"`
		IgmpSnooping           bool        `json:"igmp_snooping,omitempty"`
		IP                     string      `json:"ip"`
		IPSubnet               string      `json:"ip_subnet"`
		IsGuest                bool        `json:"is_guest"`
		IsNat                  bool        `json:"is_nat"`
		Mac                    string      `json:"mac"`
		Name                   string      `json:"name"`
		Networkgroup           string      `json:"networkgroup"`
		NumSta                 int         `json:"num_sta"`
		Purpose                string      `json:"purpose"`
		RxBytes                int         `json:"rx_bytes"`
		RxPackets              int         `json:"rx_packets"`
		SiteID                 string      `json:"site_id"`
		TxBytes                int         `json:"tx_bytes"`
		TxPackets              int         `json:"tx_packets"`
		Up                     string      `json:"up"`
		Vlan                   string      `json:"vlan,omitempty"`
		VlanEnabled            bool        `json:"vlan_enabled"`
		DhcpRelayEnabled       bool        `json:"dhcp_relay_enabled,omitempty"`
		DhcpdGatewayEnabled    bool        `json:"dhcpd_gateway_enabled,omitempty"`
		DhcpdNtp1              string      `json:"dhcpd_ntp_1,omitempty"`
		DhcpdNtpEnabled        bool        `json:"dhcpd_ntp_enabled,omitempty"`
		DhcpdTimeOffsetEnabled bool        `json:"dhcpd_time_offset_enabled,omitempty"`
		DhcpdUnifiController   string      `json:"dhcpd_unifi_controller,omitempty"`
		Ipv6InterfaceType      string      `json:"ipv6_interface_type,omitempty"`
		AttrHiddenID           string      `json:"attr_hidden_id,omitempty"`
		AttrNoDelete           bool        `json:"attr_no_delete,omitempty"`
		UpnpLanEnabled         bool        `json:"upnp_lan_enabled,omitempty"`
	} `json:"network_table,omitempty"`
	NumDesktop      int `json:"num_desktop,omitempty"`
	NumHandheld     int `json:"num_handheld,omitempty"`
	NumMobile       int `json:"num_mobile,omitempty"`
	SpeedtestStatus struct {
		Latency        int     `json:"latency"`
		Rundate        int     `json:"rundate"`
		Runtime        int     `json:"runtime"`
		StatusDownload int     `json:"status_download"`
		StatusPing     int     `json:"status_ping"`
		StatusSummary  int     `json:"status_summary"`
		StatusUpload   int     `json:"status_upload"`
		XputDownload   float64 `json:"xput_download"`
		XputUpload     float64 `json:"xput_upload"`
	} `json:"speedtest-status,omitempty"`
	SpeedtestStatusSaved bool `json:"speedtest-status-saved,omitempty"`
	UsgCaps              int  `json:"usg_caps,omitempty"`
	Wan1                 struct {
		BytesR      int      `json:"bytes-r"`
		DNS         []string `json:"dns"`
		Enable      bool     `json:"enable"`
		FullDuplex  bool     `json:"full_duplex"`
		Gateway     string   `json:"gateway"`
		Ifname      string   `json:"ifname"`
		IP          string   `json:"ip"`
		Mac         string   `json:"mac"`
		MaxSpeed    int      `json:"max_speed"`
		Name        string   `json:"name"`
		Netmask     string   `json:"netmask"`
		RxBytes     int64    `json:"rx_bytes"`
		RxBytesR    int      `json:"rx_bytes-r"`
		RxDropped   int      `json:"rx_dropped"`
		RxErrors    int      `json:"rx_errors"`
		RxMulticast int      `json:"rx_multicast"`
		RxPackets   int      `json:"rx_packets"`
		Speed       int      `json:"speed"`
		TxBytes     int64    `json:"tx_bytes"`
		TxBytesR    int      `json:"tx_bytes-r"`
		TxDropped   int      `json:"tx_dropped"`
		TxErrors    int      `json:"tx_errors"`
		TxPackets   int      `json:"tx_packets"`
		Type        string   `json:"type"`
		Up          bool     `json:"up"`
	} `json:"wan1,omitempty"`
	Wan2 struct {
		BytesR      int      `json:"bytes-r"`
		DNS         []string `json:"dns"`
		Enable      bool     `json:"enable"`
		FullDuplex  bool     `json:"full_duplex"`
		Gateway     string   `json:"gateway"`
		Ifname      string   `json:"ifname"`
		IP          string   `json:"ip"`
		Mac         string   `json:"mac"`
		MaxSpeed    int      `json:"max_speed"`
		Name        string   `json:"name"`
		Netmask     string   `json:"netmask"`
		RxBytes     int64    `json:"rx_bytes"`
		RxBytesR    int      `json:"rx_bytes-r"`
		RxDropped   int      `json:"rx_dropped"`
		RxErrors    int      `json:"rx_errors"`
		RxMulticast int      `json:"rx_multicast"`
		RxPackets   int      `json:"rx_packets"`
		Speed       int      `json:"speed"`
		TxBytes     int64    `json:"tx_bytes"`
		TxBytesR    int      `json:"tx_bytes-r"`
		TxDropped   int      `json:"tx_dropped"`
		TxErrors    int      `json:"tx_errors"`
		TxPackets   int      `json:"tx_packets"`
		Type        string   `json:"type"`
		Up          bool     `json:"up"`
	} `json:"wan2,omitempty"`
}

// GetUnifiDevices returns a response full of devices' data from the Unifi Controller.
func (c *Config) GetUnifiDevices() ([]Device, error) {
	response := &DeviceResponse{}
	if req, err := c.uniRequest(DevicePath, ""); err != nil {
		return nil, err
	} else if resp, err := c.uniClient.Do(req); err != nil {
		return nil, err
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	} else if err = resp.Body.Close(); err != nil {
		log.Println("resp.Body.Close():", err) // Not fatal? Just log it.
	}
	return response.Devices, nil
}

// Point generates a device's datapoint for InfluxDB.
func (d Device) Point() (*influx.Point, error) {
	tags := map[string]string{
		"id":                      d.ID,
		"mac":                     d.Mac,
		"device_type":             d.Stat.O,
		"device_oid":              d.Stat.Oid,
		"device_ap":               d.Stat.Ap,
		"site_id":                 d.SiteID,
		"name":                    d.Name,
		"addopted":                strconv.FormatBool(d.Adopted),
		"adopt_ip":                d.AdoptIP,
		"adopt_url":               d.AdoptURL,
		"bandsteering_mode":       d.BandsteeringMode,
		"board_rev":               strconv.Itoa(d.BoardRev),
		"cfgversion":              d.Cfgversion,
		"config_network_ip":       d.ConfigNetwork.IP,
		"config_network_type":     d.ConfigNetwork.Type,
		"connect_request_ip":      d.ConnectRequestIP,
		"connect_request_port":    d.ConnectRequestPort,
		"default":                 strconv.FormatBool(d.Default),
		"device_id":               d.DeviceID,
		"discovered_via":          d.DiscoveredVia,
		"fw_caps":                 strconv.Itoa(d.FwCaps),
		"guest-num_sta":           strconv.Itoa(d.GuestNumSta),
		"guest_token":             d.GuestToken,
		"has_eth1":                strconv.FormatBool(d.HasEth1),
		"has_speaker":             strconv.FormatBool(d.HasSpeaker),
		"inform_ip":               d.InformIP,
		"isolated":                strconv.FormatBool(d.Isolated),
		"last_seen":               strconv.Itoa(d.LastSeen),
		"last_uplink_mac":         d.LastUplink.UplinkMac,
		"last_uplink_remote_port": strconv.Itoa(d.LastUplink.UplinkRemotePort),
		"known_cfgversion":        d.KnownCfgversion,
		"led_override":            d.LedOverride,
		"locating":                strconv.FormatBool(d.Locating),
		"model":                   d.Model,
		"outdoor_mode_override":   d.OutdoorModeOverride,
		"serial":                  d.Serial,
		"type":                    d.Type,
		"version_incompatible":    strconv.FormatBool(d.VersionIncompatible),
		"vwireEnabled":            strconv.FormatBool(d.VwireEnabled),
		"wifi_caps":               strconv.Itoa(d.WifiCaps),
		"dot1x_portctrl_enabled":  strconv.FormatBool(d.Dot1XPortctrlEnabled),
		"flowctrl_enabled":        strconv.FormatBool(d.FlowctrlEnabled),
		"has_fan":                 strconv.FormatBool(d.HasFan),
		"has_temperature":         strconv.FormatBool(d.HasTemperature),
		"jumboframe_enabled":      strconv.FormatBool(d.JumboframeEnabled),
		"stp_priority":            d.StpPriority,
		"stp_version":             d.StpVersion,
		"uplink_depth":            strconv.Itoa(d.UplinkDepth),
		"config_network_wan_type": d.ConfigNetworkWan.Type,
		"usg_caps":                strconv.Itoa(d.UsgCaps),
		"speedtest-status-saved":  strconv.FormatBool(d.SpeedtestStatusSaved),
	}
	fields := map[string]interface{}{
		"ip":                             d.IP,
		"bytes":                          d.Bytes,
		"bytes_d":                        d.BytesD,
		"bytes_r":                        d.BytesR,
		"fan_level":                      d.FanLevel,
		"general_temperature":            d.GeneralTemperature,
		"last_seen":                      d.LastSeen,
		"license_state":                  d.LicenseState,
		"overheating":                    d.Overheating,
		"rx_bytes":                       d.RxBytes,
		"rx_bytes-d":                     d.RxBytesD,
		"tx_bytes":                       d.TxBytes,
		"tx_bytes-d":                     d.TxBytesD,
		"uptime":                         d.Uptime,
		"considered_lost_at":             d.ConsideredLostAt,
		"next_heartbeat_at":              d.NextHeartbeatAt,
		"scanning":                       d.Scanning,
		"spectrum_scanning":              d.SpectrumScanning,
		"roll_upgrade":                   d.Rollupgrade,
		"state":                          d.State,
		"upgradable":                     d.Upgradable,
		"user-num_sta":                   d.UserNumSta,
		"version":                        d.Version,
		"num_desktop":                    d.NumDesktop,
		"num_handheld":                   d.NumHandheld,
		"num_mobile":                     d.NumMobile,
		"speedtest-status_latency":       d.SpeedtestStatus.Latency,
		"speedtest-status_rundate":       d.SpeedtestStatus.Rundate,
		"speedtest-status_runtime":       d.SpeedtestStatus.Runtime,
		"speedtest-status_download":      d.SpeedtestStatus.StatusDownload,
		"speedtest-status_ping":          d.SpeedtestStatus.StatusPing,
		"speedtest-status_summary":       d.SpeedtestStatus.StatusSummary,
		"speedtest-status_upload":        d.SpeedtestStatus.StatusUpload,
		"speedtest-status_xput_download": d.SpeedtestStatus.XputDownload,
		"speedtest-status_xput_upload":   d.SpeedtestStatus.XputUpload,
		"wan1_bytes-r":                   d.Wan1.BytesR,
		"wan1_enable":                    d.Wan1.Enable,
		"wan1_full_duplex":               d.Wan1.FullDuplex,
		"wan1_gateway":                   d.Wan1.Gateway,
		"wan1_ifname":                    d.Wan1.Ifname,
		"wan1_ip":                        d.Wan1.IP,
		"wan1_mac":                       d.Wan1.Mac,
		"wan1_max_speed":                 d.Wan1.MaxSpeed,
		"wan1_name":                      d.Wan1.Name,
		"wan1_netmask":                   d.Wan1.Netmask,
		"wan1_rx_bytes":                  d.Wan1.RxBytes,
		"wan1_rx_bytes-r":                d.Wan1.RxBytesR,
		"wan1_rx_dropped":                d.Wan1.RxDropped,
		"wan1_rx_errors":                 d.Wan1.RxErrors,
		"wan1_rx_multicast":              d.Wan1.RxMulticast,
		"wan1_rx_packets":                d.Wan1.RxPackets,
		"wan1_type":                      d.Wan1.Type,
		"wan1_speed":                     d.Wan1.Speed,
		"wan1_up":                        d.Wan1.Up,
		"wan1_tx_bytes":                  d.Wan1.TxBytes,
		"wan1_tx_bytes-r":                d.Wan1.TxBytesR,
		"wan1_tx_dropped":                d.Wan1.TxDropped,
		"wan1_tx_errors":                 d.Wan1.TxErrors,
		"wan1_tx_packets":                d.Wan1.TxPackets,
		"loadavg_1":                      d.SysStats.Loadavg1,
		"loadavg_5":                      d.SysStats.Loadavg5,
		"loadavg_15":                     d.SysStats.Loadavg15,
		"mem_buffer":                     d.SysStats.MemBuffer,
		"mem_total":                      d.SysStats.MemTotal,
		"cpu":                            d.SystemStats.CPU,
		"mem":                            d.SystemStats.Mem,
		"system_uptime":                  d.SystemStats.Uptime,
		"stat_bytes":                     d.Stat.Bytes,
		"stat_duration":                  d.Stat.Duration,
		"stat_guest-rx_bytes":            d.Stat.RxBytes,
		"stat_guest-rx_crypts":           d.Stat.RxCrypts,
		"stat_guest-rx_dropped":          d.Stat.RxDropped,
		"stat_guest-rx_errors":           d.Stat.RxErrors,
		"stat_guest-rx_frags":            d.Stat.RxFrags,
		"stat_guest-rx_packets":          d.Stat.RxPackets,
		"stat_guest-tx_bytes":            d.Stat.TxBytes,
		"stat_guest-tx_dropped":          d.Stat.TxDropped,
		"stat_guest-tx_errors":           d.Stat.TxErrors,
		"stat_guest-tx_packets":          d.Stat.TxPackets,
		"stat_guest-tx_retries":          d.Stat.TxRetries,
		"stat_port_1-rx_broadcast":       d.Stat.Port1RxBroadcast,
		"stat_port_1-rx_bytes":           d.Stat.Port1RxBytes,
		"stat_port_1-rx_multicast":       d.Stat.Port1RxMulticast,
		"stat_port_1-rx_packets":         d.Stat.Port1RxPackets,
		"stat_port_1-tx_broadcast":       d.Stat.Port1TxBroadcast,
		"stat_port_1-tx_bytes":           d.Stat.Port1TxBytes,
		"stat_port_1-tx_multicast":       d.Stat.Port1TxMulticast,
		"stat_port_1-tx_packets":         d.Stat.Port1TxPackets,
		"stat_rx_bytes":                  d.Stat.RxBytes,
		"stat_rx_crypts":                 d.Stat.RxCrypts,
		"stat_rx_dropped":                d.Stat.RxDropped,
		"stat_rx_errors":                 d.Stat.RxErrors,
		"stat_rx_frags":                  d.Stat.RxFrags,
		"stat_rx_packets":                d.Stat.TxPackets,
		"stat_tx_bytes":                  d.Stat.TxBytes,
		"stat_tx_dropped":                d.Stat.TxDropped,
		"stat_tx_errors":                 d.Stat.TxErrors,
		"stat_tx_packets":                d.Stat.TxPackets,
		"stat_tx_retries":                d.Stat.TxRetries,
		"stat_user-rx_bytes":             d.Stat.UserRxBytes,
		"stat_user-rx_crypts":            d.Stat.UserRxCrypts,
		"stat_user-rx_dropped":           d.Stat.UserRxDropped,
		"stat_user-rx_errors":            d.Stat.UserRxErrors,
		"stat_user-rx_frags":             d.Stat.UserRxFrags,
		"stat_user-rx_packets":           d.Stat.UserRxPackets,
		"stat_user-tx_bytes":             d.Stat.UserTxBytes,
		"stat_user-tx_dropped":           d.Stat.UserTxDropped,
		"stat_user-tx_errors":            d.Stat.UserTxErrors,
		"stat_user-tx_packets":           d.Stat.UserTxPackets,
		"stat_user-tx_retries":           d.Stat.UserTxRetries,
		"stat_user-wifi0-rx_bytes":       d.Stat.UserWifi0RxBytes,
		"stat_user-wifi0-rx_crypts":      d.Stat.UserWifi0RxCrypts,
		"stat_user-wifi0-rx_dropped":     d.Stat.UserWifi0RxDropped,
		"stat_user-wifi0-rx_errors":      d.Stat.UserWifi0RxErrors,
		"stat_user-wifi0-rx_frags":       d.Stat.UserWifi0RxFrags,
		"stat_user-wifi0-rx_packets":     d.Stat.UserWifi0RxPackets,
		"stat_user-wifi0-tx_bytes":       d.Stat.UserWifi0TxBytes,
		"stat_user-wifi0-tx_dropped":     d.Stat.UserWifi0TxDropped,
		"stat_user-wifi0-tx_errors":      d.Stat.UserWifi0TxErrors,
		"stat_user-wifi0-tx_packets":     d.Stat.UserWifi0TxPackets,
		"stat_user-wifi0-tx_retries":     d.Stat.UserWifi0TxRetries,
		"stat_user-wifi1-rx_bytes":       d.Stat.UserWifi1RxBytes,
		"stat_user-wifi1-rx_crypts":      d.Stat.UserWifi1RxCrypts,
		"stat_user-wifi1-rx_dropped":     d.Stat.UserWifi1RxDropped,
		"stat_user-wifi1-rx_errors":      d.Stat.UserWifi1RxErrors,
		"stat_user-wifi1-rx_frags":       d.Stat.UserWifi1RxFrags,
		"stat_user-wifi1-rx_packets":     d.Stat.UserWifi1RxPackets,
		"stat_user-wifi1-tx_bytes":       d.Stat.UserWifi1TxBytes,
		"stat_user-wifi1-tx_dropped":     d.Stat.UserWifi1TxDropped,
		"stat_user-wifi1-tx_errors":      d.Stat.UserWifi1TxErrors,
		"stat_user-wifi1-tx_packets":     d.Stat.UserWifi1TxPackets,
		"stat_user-wifi1-tx_retries":     d.Stat.UserWifi1TxRetries,
		"stat_wifi0-rx_bytes":            d.Stat.Wifi0RxBytes,
		"stat_wifi0-rx_crypts":           d.Stat.Wifi0RxCrypts,
		"stat_wifi0-rx_dropped":          d.Stat.Wifi0RxDropped,
		"stat_wifi0-rx_errors":           d.Stat.Wifi0RxErrors,
		"stat_wifi0-rx_frags":            d.Stat.Wifi0RxFrags,
		"stat_wifi0-rx_packets":          d.Stat.Wifi0RxPackets,
		"stat_wifi0-tx_bytes":            d.Stat.Wifi0TxBytes,
		"stat_wifi0-tx_dropped":          d.Stat.Wifi0TxDropped,
		"stat_wifi0-tx_errors":           d.Stat.Wifi0TxErrors,
		"stat_wifi0-tx_packets":          d.Stat.Wifi0TxPackets,
		"stat_wifi0-tx_retries":          d.Stat.Wifi0TxRetries,
		"stat_wifi1-rx_bytes":            d.Stat.Wifi1RxBytes,
		"stat_wifi1-rx_crypts":           d.Stat.Wifi1RxCrypts,
		"stat_wifi1-rx_dropped":          d.Stat.Wifi1RxDropped,
		"stat_wifi1-rx_errors":           d.Stat.Wifi1RxErrors,
		"stat_wifi1-rx_frags":            d.Stat.Wifi1RxFrags,
		"stat_wifi1-rx_packets":          d.Stat.Wifi1RxPackets,
		"stat_wifi1-tx_bytes":            d.Stat.Wifi1TxBytes,
		"stat_wifi1-tx_dropped":          d.Stat.Wifi1TxDropped,
		"stat_wifi1-tx_errors":           d.Stat.Wifi1TxErrors,
		"stat_wifi1-tx_packets":          d.Stat.Wifi1TxPackets,
		"stat_wifi1-tx_retries":          d.Stat.Wifi1TxRetries,
	}
	return influx.NewPoint("devices", tags, fields, time.Now())
}
