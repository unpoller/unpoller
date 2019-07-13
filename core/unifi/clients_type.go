package unifi

// Clients contains a list that contains all of the unifi clients from a controller.
type Clients []*Client

// Client defines all the data a connected-network client contains.
type Client struct {
	Anomalies   int64   `json:"anomalies,omitempty"`
	ApMac       string  `json:"ap_mac"`
	ApName      string  `json:"-"`
	AssocTime   int64   `json:"assoc_time"`
	Blocked     bool    `json:"blocked,omitempty"`
	Bssid       string  `json:"bssid"`
	BytesR      int64   `json:"bytes-r"`
	Ccq         int64   `json:"ccq"`
	Channel     FlexInt `json:"channel"`
	DevCat      FlexInt `json:"dev_cat"`
	DevFamily   FlexInt `json:"dev_family"`
	DevID       FlexInt `json:"dev_id"`
	DevVendor   FlexInt `json:"dev_vendor,omitempty"`
	DhcpendTime int     `json:"dhcpend_time,omitempty"`
	DpiStats    struct {
		App       FlexInt
		Cat       FlexInt
		RxBytes   FlexInt
		RxPackets FlexInt
		TxBytes   FlexInt
		TxPackets FlexInt
	} `json:"dpi_stats"`
	DpiStatsLastUpdated int64    `json:"dpi_stats_last_updated"`
	Essid               string   `json:"essid"`
	FirstSeen           int64    `json:"first_seen"`
	FixedIP             string   `json:"fixed_ip"`
	GwMac               string   `json:"gw_mac"`
	GwName              string   `json:"-"`
	Hostname            string   `json:"hostname"`
	ID                  string   `json:"_id"`
	IP                  string   `json:"ip"`
	IdleTime            int64    `json:"idle_time"`
	Is11R               FlexBool `json:"is_11r"`
	IsGuest             FlexBool `json:"is_guest"`
	IsGuestByUAP        FlexBool `json:"_is_guest_by_uap"`
	IsGuestByUGW        FlexBool `json:"_is_guest_by_ugw"`
	IsGuestByUSW        FlexBool `json:"_is_guest_by_usw"`
	IsWired             FlexBool `json:"is_wired"`
	LastSeen            int64    `json:"last_seen"`
	LastSeenByUAP       int64    `json:"_last_seen_by_uap"`
	LastSeenByUGW       int64    `json:"_last_seen_by_ugw"`
	LastSeenByUSW       int64    `json:"_last_seen_by_usw"`
	LatestAssocTime     int64    `json:"latest_assoc_time"`
	Mac                 string   `json:"mac"`
	Name                string   `json:"name"`
	Network             string   `json:"network"`
	NetworkID           string   `json:"network_id"`
	Noise               int64    `json:"noise"`
	Note                string   `json:"note"`
	Noted               FlexBool `json:"noted"`
	OsClass             FlexInt  `json:"os_class"`
	OsName              FlexInt  `json:"os_name"`
	Oui                 string   `json:"oui"`
	PowersaveEnabled    FlexBool `json:"powersave_enabled"`
	QosPolicyApplied    FlexBool `json:"qos_policy_applied"`
	Radio               string   `json:"radio"`
	RadioName           string   `json:"radio_name"`
	RadioProto          string   `json:"radio_proto"`
	RoamCount           int64    `json:"roam_count"`
	Rssi                int64    `json:"rssi"`
	RxBytes             int64    `json:"rx_bytes"`
	RxBytesR            int64    `json:"rx_bytes-r"`
	RxPackets           int64    `json:"rx_packets"`
	RxRate              int64    `json:"rx_rate"`
	Signal              int64    `json:"signal"`
	SiteID              string   `json:"site_id"`
	SiteName            string   `json:"-"`
	SwDepth             int      `json:"sw_depth"`
	SwMac               string   `json:"sw_mac"`
	SwName              string   `json:"-"`
	SwPort              FlexInt  `json:"sw_port"`
	TxBytes             int64    `json:"tx_bytes"`
	TxBytesR            int64    `json:"tx_bytes-r"`
	TxPackets           int64    `json:"tx_packets"`
	TxPower             int64    `json:"tx_power"`
	TxRate              int64    `json:"tx_rate"`
	Uptime              int64    `json:"uptime"`
	UptimeByUAP         int64    `json:"_uptime_by_uap"`
	UptimeByUGW         int64    `json:"_uptime_by_ugw"`
	UptimeByUSW         int64    `json:"_uptime_by_usw"`
	UseFixedIP          FlexBool `json:"use_fixedip"`
	UserGroupID         string   `json:"usergroup_id"`
	UserID              string   `json:"user_id"`
	Vlan                FlexInt  `json:"vlan"`
	WifiTxAttempts      int64    `json:"wifi_tx_attempts"`
	WiredRxBytes        int64    `json:"wired-rx_bytes"`
	WiredRxBytesR       int64    `json:"wired-rx_bytes-r"`
	WiredRxPackets      int64    `json:"wired-rx_packets"`
	WiredTxBytes        int64    `json:"wired-tx_bytes"`
	WiredTxBytesR       int64    `json:"wired-tx_bytes-r"`
	WiredTxPackets      int64    `json:"wired-tx_packets"`
}
