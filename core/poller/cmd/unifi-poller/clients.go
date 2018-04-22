package main

import (
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

// ClientResponse marshalls the payload from the controller.
type ClientResponse struct {
	Clients []Client `json:"data"`
	Meta    struct {
		Rc string `json:"rc"`
	} `json:"meta"`
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
	DevCat              int       `json:"dev_cat"`
	DevFamily           int       `json:"dev_family"`
	DevID               int       `json:"dev_id"`
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
	OsClass             int       `json:"os_class"`
	OsName              int       `json:"os_name"`
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
	WiredRxBytes        int64     `json:"wired-rx_bytes"`
	WiredRxBytesR       int64     `json:"wired-rx_bytes-r"`
	WiredRxPackets      int64     `json:"wired-rx_packets"`
	WiredTxBytes        int64     `json:"wired-tx_bytes"`
	WiredTxBytesR       int64     `json:"wired-tx_bytes-r"`
	WiredTxPackets      int64     `json:"wired-tx_packets"`
}

// Point generates a client's datapoint for InfluxDB.
func (c Client) Point() (*influx.Point, error) {
	if c.Name == "" && c.Hostname != "" {
		c.Name = c.Hostname
	} else if c.Hostname == "" && c.Name != "" {
		c.Hostname = c.Name
	} else if c.Hostname == "" && c.Name == "" {
		c.Hostname = "-no-name-"
		c.Name = "-no-name-"
	}
	tags := map[string]string{
		"id":                 c.ID,
		"mac":                c.Mac,
		"user_id":            c.UserID,
		"site_id":            c.SiteID,
		"network_id":         c.NetworkID,
		"usergroup_id":       c.UserGroupID,
		"ap_mac":             c.ApMac,
		"gw_mac":             c.GwMac,
		"sw_mac":             c.SwMac,
		"oui":                c.Oui,
		"radio_name":         c.RadioName,
		"radio":              c.Radio,
		"radio_proto":        c.RadioProto,
		"name":               c.Name,
		"fixed_ip":           c.FixedIP,
		"sw_port":            strconv.Itoa(c.SwPort),
		"os_class":           strconv.Itoa(c.OsClass),
		"os_name":            strconv.Itoa(c.OsName),
		"dev_cat":            strconv.Itoa(c.DevCat),
		"dev_id":             strconv.Itoa(c.DevID),
		"dev_family":         strconv.Itoa(c.DevFamily),
		"authorized":         strconv.FormatBool(c.Authorized),
		"is_11r":             strconv.FormatBool(c.Is11R),
		"is_wired":           strconv.FormatBool(c.IsWired),
		"is_guest":           strconv.FormatBool(c.IsGuest),
		"is_guest_by_uap":    strconv.FormatBool(c.IsGuestByUAP),
		"is_guest_by_ugw":    strconv.FormatBool(c.IsGuestByUGW),
		"is_guest_by_usw":    strconv.FormatBool(c.IsGuestByUSW),
		"noted":              strconv.FormatBool(c.Noted),
		"powersave_enabled":  strconv.FormatBool(c.PowersaveEnabled),
		"qos_policy_applied": strconv.FormatBool(c.QosPolicyApplied),
		"use_fixedip":        strconv.FormatBool(c.UseFixedIP),
		"channel":            strconv.Itoa(c.Channel),
		"vlan":               strconv.Itoa(c.Vlan),
	}
	fields := map[string]interface{}{
		"ip":                     c.IP,
		"essid":                  c.Essid,
		"bssid":                  c.Bssid,
		"hostname":               c.Hostname,
		"dpi_stats_last_updated": c.DpiStatsLastUpdated,
		"last_seen_by_uap":       c.LastSeenByUAP,
		"last_seen_by_ugw":       c.LastSeenByUGW,
		"last_seen_by_usw":       c.LastSeenByUSW,
		"uptime_by_uap":          c.UptimeByUAP,
		"uptime_by_ugw":          c.UptimeByUGW,
		"uptime_by_usw":          c.UptimeByUSW,
		"assoc_time":             c.AssocTime,
		"bytes_r":                c.BytesR,
		"ccq":                    c.Ccq,
		"first_seen":             c.FirstSeen,
		"idle_time":              c.IdleTime,
		"last_seen":              c.LastSeen,
		"latest_assoc_time":      c.LatestAssocTime,
		"network":                c.Network,
		"noise":                  c.Noise,
		"note":                   c.Note,
		"roam_count":             c.RoamCount,
		"rssi":                   c.Rssi,
		"rx_bytes":               c.RxBytes,
		"rx_bytes_r":             c.RxBytesR,
		"rx_packets":             c.RxPackets,
		"rx_rate":                c.RxRate,
		"signal":                 c.Signal,
		"tx_bytes":               c.TxBytes,
		"tx_bytes_r":             c.TxBytesR,
		"tx_packets":             c.TxPackets,
		"tx_power":               c.TxPower,
		"tx_rate":                c.TxRate,
		"uptime":                 c.Uptime,
		"wired-rx_bytes":         c.WiredRxBytes,
		"wired-rx_bytes-r":       c.WiredRxBytesR,
		"wired-rx_packets":       c.WiredRxPackets,
		"wired-tx_bytes":         c.WiredTxBytes,
		"wired-tx_bytes-r":       c.WiredTxBytesR,
		"wired-tx_packets":       c.WiredTxPackets,
	}

	return influx.NewPoint("clients", tags, fields, time.Now())
}
