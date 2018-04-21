package main

import (
	"net/http"
	"time"
)

const (
	// LoginPath is Unifi Controller Login API Path
	LoginPath = "/api/login"
	// ClientPath is Unifi Clients API Path
	ClientPath = "/api/s/default/stat/sta"
)

// ClientResponse marshalls the payload from the controller.
type ClientResponse struct {
	Clients []Client `json:"data"`
	Meta    struct {
		Rc string `json:"rc"`
	} `json:"meta"`
}

// Config represents the data needed to poll a controller and report to influxdb.
type Config struct {
	Interval   time.Duration `json:"interval",toml:"interval",yaml:"interval"`
	InfluxAddr string        `json:"influx_addr",toml:"influx_addr",yaml:"influx_addr"`
	InfluxUser string        `json:"influx_user",toml:"influx_user",yaml:"influx_user"`
	InfluxPass string        `json:"influx_pass",toml:"influx_pass",yaml:"influx_pass"`
	InfluxDB   string        `json:"influx_db",toml:"influx_db",yaml:"influx_db"`
	UnifiUser  string        `json:"unifi_user",toml:"unifi_user",yaml:"unifi_user"`
	UnifiPass  string        `json:"unifi_pass"toml:"unifi_pass",yaml:"unifi_pass"`
	UnifiBase  string        `json:"unifi_url",toml:"unifi_url",yaml:"unifi_url"`
	uniClient  *http.Client
}

// DpiStat is for deep packet inspection stats.
// Does not seem to exist in Unifi 5.7.20.
type DpiStat struct {
	App       int64
	Cat       int64
	RxBytes   int64
	RxPackets int64
	TxBytes   int64
	TxPackets int64
}

// Client defines all the data a connected-network client contains.
type Client struct {
	ID                  string    `json:"_id"`
	IsGuestByUAP        bool      `json:"_is_guest_by_uap"`
	IsGuestByUGW        bool      `json:"_is_guest_by_ugw"`
	IsGuestByUSW        bool      `json:"_is_guest_by_usw"`
	LastSeenByUAP       int64     `json:"_last_seen_by_uap"`
	LastSeenByUGW       int64     `json:"_last_seen_by_ugw"`
	LastSeenByUSW       int64     `json:"_last_seen_by_usw"`
	UptimeByUAP         int64     `json:"_uptime_by_uap"`
	UptimeByUGW         int64     `json:"_uptime_by_ugw"`
	UptimeByUSW         int64     `json:"_uptime_by_usw"`
	ApMac               string    `json:"ap_mac"`
	AssocTime           int64     `json:"assoc_time"`
	Authorized          bool      `json:"authorized"`
	Bssid               string    `json:"bssid"`
	BytesR              int64     `json:"bytes-r"`
	Ccq                 int64     `json:"ccq"`
	Channel             int       `json:"channel"`
	DpiStats            []DpiStat `json:"dpi_stats"`
	DpiStatsLastUpdated int64     `json:"dpi_stats_last_updated"`
	Essid               string    `json:"essid"`
	FirstSeen           int64     `json:"first_seen"`
	FixedIP             string    `json:"fixed_ip"`
	Hostname            string    `json:"hostname"`
	GwMac               string    `json:"gw_mac"`
	IdleTime            int64     `json:"idle_time"`
	IP                  string    `json:"ip"`
	Is11R               bool      `json:"is_11r"`
	IsGuest             bool      `json:"is_guest"`
	IsWired             bool      `json:"is_wired"`
	LastSeen            int64     `json:"last_seen"`
	LatestAssocTime     int64     `json:"latest_assoc_time"`
	Mac                 string    `json:"mac"`
	Name                string    `json:"name"`
	Network             string    `json:"network"`
	NetworkID           string    `json:"network_id"`
	Noise               int64     `json:"noise"`
	Note                string    `json:"note"`
	Noted               bool      `json:"noted"`
	Oui                 string    `json:"oui"`
	PowersaveEnabled    bool      `json:"powersave_enabled"`
	QosPolicyApplied    bool      `json:"qos_policy_applied"`
	Radio               string    `json:"radio"`
	RadioName           string    `json:"radio_name"`
	RadioProto          string    `json:"radio_proto"`
	RoamCount           int64     `json:"roam_count"`
	Rssi                int64     `json:"rssi"`
	RxBytes             int64     `json:"rx_bytes"`
	RxBytesR            int64     `json:"rx_bytes-r"`
	RxPackets           int64     `json:"rx_packets"`
	RxRate              int64     `json:"rx_rate"`
	Signal              int64     `json:"signal"`
	SiteID              string    `json:"site_id"`
	SwDepth             int       `json:"sw_depth"`
	SwMac               string    `json:"sw_mac"`
	SwPort              int       `json:"sw_port"`
	TxBytes             int64     `json:"tx_bytes"`
	TxBytesR            int64     `json:"tx_bytes-r"`
	TxPackets           int64     `json:"tx_packets"`
	TxPower             int64     `json:"tx_power"`
	TxRate              int64     `json:"tx_rate"`
	Uptime              int64     `json:"uptime"`
	UserID              string    `json:"user_id"`
	UserGroupID         string    `json:"usergroup_id"`
	UseFixedIP          bool      `json:"use_fixedip"`
	Vlan                int       `json:"vlan"`
}
