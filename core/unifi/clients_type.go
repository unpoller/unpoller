package unifi

// UCL defines all the data a connected-network client contains.
type UCL struct {
	ID            string   `json:"_id"`
	IsGuestByUAP  FlexBool `json:"_is_guest_by_uap"`
	IsGuestByUGW  FlexBool `json:"_is_guest_by_ugw"`
	IsGuestByUSW  FlexBool `json:"_is_guest_by_usw"`
	LastSeenByUAP int64    `json:"_last_seen_by_uap"`
	LastSeenByUGW int64    `json:"_last_seen_by_ugw"`
	LastSeenByUSW int64    `json:"_last_seen_by_usw"`
	UptimeByUAP   int64    `json:"_uptime_by_uap"`
	UptimeByUGW   int64    `json:"_uptime_by_ugw"`
	UptimeByUSW   int64    `json:"_uptime_by_usw"`
	ApMac         string   `json:"ap_mac"`
	AssocTime     int64    `json:"assoc_time"`
	Authorized    FlexBool `json:"authorized"`
	Bssid         string   `json:"bssid"`
	BytesR        int64    `json:"bytes-r"`
	Ccq           int64    `json:"ccq"`
	Channel       int      `json:"channel"`
	DevCat        int      `json:"dev_cat"`
	DevFamily     int      `json:"dev_family"`
	DevID         int      `json:"dev_id"`
	DpiStats      struct {
		App       int64
		Cat       int64
		RxBytes   int64
		RxPackets int64
		TxBytes   int64
		TxPackets int64
	} `json:"dpi_stats"`
	DpiStatsLastUpdated int64    `json:"dpi_stats_last_updated"`
	Essid               string   `json:"essid"`
	FirstSeen           int64    `json:"first_seen"`
	FixedIP             string   `json:"fixed_ip"`
	Hostname            string   `json:"hostname"`
	GwMac               string   `json:"gw_mac"`
	IdleTime            int64    `json:"idle_time"`
	IP                  string   `json:"ip"`
	Is11R               FlexBool `json:"is_11r"`
	IsGuest             FlexBool `json:"is_guest"`
	IsWired             FlexBool `json:"is_wired"`
	LastSeen            int64    `json:"last_seen"`
	LatestAssocTime     int64    `json:"latest_assoc_time"`
	Mac                 string   `json:"mac"`
	Name                string   `json:"name"`
	Network             string   `json:"network"`
	NetworkID           string   `json:"network_id"`
	Noise               int64    `json:"noise"`
	Note                string   `json:"note"`
	Noted               FlexBool `json:"noted"`
	OsClass             int      `json:"os_class"`
	OsName              int      `json:"os_name"`
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
	SwPort              int      `json:"sw_port"`
	TxBytes             int64    `json:"tx_bytes"`
	TxBytesR            int64    `json:"tx_bytes-r"`
	TxPackets           int64    `json:"tx_packets"`
	TxPower             int64    `json:"tx_power"`
	TxRate              int64    `json:"tx_rate"`
	Uptime              int64    `json:"uptime"`
	UserID              string   `json:"user_id"`
	UserGroupID         string   `json:"usergroup_id"`
	UseFixedIP          FlexBool `json:"use_fixedip"`
	Vlan                int      `json:"vlan"`
	WiredRxBytes        int64    `json:"wired-rx_bytes"`
	WiredRxBytesR       int64    `json:"wired-rx_bytes-r"`
	WiredRxPackets      int64    `json:"wired-rx_packets"`
	WiredTxBytes        int64    `json:"wired-tx_bytes"`
	WiredTxBytesR       int64    `json:"wired-tx_bytes-r"`
	WiredTxPackets      int64    `json:"wired-tx_packets"`
}
