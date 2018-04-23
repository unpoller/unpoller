package unidev

import (
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

// Point generates a client's datapoint for InfluxDB.
func (u UCL) Point() (*influx.Point, error) {
	if u.Name == "" && u.Hostname != "" {
		u.Name = u.Hostname
	} else if u.Hostname == "" && u.Name != "" {
		u.Hostname = u.Name
	} else if u.Hostname == "" && u.Name == "" {
		u.Hostname = "-no-name-"
		u.Name = "-no-name-"
	}
	tags := map[string]string{
		"id":                 u.ID,
		"mac":                u.Mac,
		"user_id":            u.UserID,
		"site_id":            u.SiteID,
		"network_id":         u.NetworkID,
		"usergroup_id":       u.UserGroupID,
		"ap_mac":             u.ApMac,
		"gw_mac":             u.GwMac,
		"sw_mac":             u.SwMac,
		"oui":                u.Oui,
		"radio_name":         u.RadioName,
		"radio":              u.Radio,
		"radio_proto":        u.RadioProto,
		"name":               u.Name,
		"fixed_ip":           u.FixedIP,
		"sw_port":            strconv.Itoa(u.SwPort),
		"os_class":           strconv.Itoa(u.OsClass),
		"os_name":            strconv.Itoa(u.OsName),
		"dev_cat":            strconv.Itoa(u.DevCat),
		"dev_id":             strconv.Itoa(u.DevID),
		"dev_family":         strconv.Itoa(u.DevFamily),
		"authorized":         strconv.FormatBool(u.Authorized),
		"is_11r":             strconv.FormatBool(u.Is11R),
		"is_wired":           strconv.FormatBool(u.IsWired),
		"is_guest":           strconv.FormatBool(u.IsGuest),
		"is_guest_by_uap":    strconv.FormatBool(u.IsGuestByUAP),
		"is_guest_by_ugw":    strconv.FormatBool(u.IsGuestByUGW),
		"is_guest_by_usw":    strconv.FormatBool(u.IsGuestByUSW),
		"noted":              strconv.FormatBool(u.Noted),
		"powersave_enabled":  strconv.FormatBool(u.PowersaveEnabled),
		"qos_policy_applied": strconv.FormatBool(u.QosPolicyApplied),
		"use_fixedip":        strconv.FormatBool(u.UseFixedIP),
		"channel":            strconv.Itoa(u.Channel),
		"vlan":               strconv.Itoa(u.Vlan),
	}
	fields := map[string]interface{}{
		"ip":                     u.IP,
		"essid":                  u.Essid,
		"bssid":                  u.Bssid,
		"hostname":               u.Hostname,
		"dpi_stats_last_updated": u.DpiStatsLastUpdated,
		"last_seen_by_uap":       u.LastSeenByUAP,
		"last_seen_by_ugw":       u.LastSeenByUGW,
		"last_seen_by_usw":       u.LastSeenByUSW,
		"uptime_by_uap":          u.UptimeByUAP,
		"uptime_by_ugw":          u.UptimeByUGW,
		"uptime_by_usw":          u.UptimeByUSW,
		"assoc_time":             u.AssocTime,
		"bytes_r":                u.BytesR,
		"ccq":                    u.Ccq,
		"first_seen":             u.FirstSeen,
		"idle_time":              u.IdleTime,
		"last_seen":              u.LastSeen,
		"latest_assoc_time":      u.LatestAssocTime,
		"network":                u.Network,
		"noise":                  u.Noise,
		"note":                   u.Note,
		"roam_count":             u.RoamCount,
		"rssi":                   u.Rssi,
		"rx_bytes":               u.RxBytes,
		"rx_bytes_r":             u.RxBytesR,
		"rx_packets":             u.RxPackets,
		"rx_rate":                u.RxRate,
		"signal":                 u.Signal,
		"tx_bytes":               u.TxBytes,
		"tx_bytes_r":             u.TxBytesR,
		"tx_packets":             u.TxPackets,
		"tx_power":               u.TxPower,
		"tx_rate":                u.TxRate,
		"uptime":                 u.Uptime,
		"wired-rx_bytes":         u.WiredRxBytes,
		"wired-rx_bytes-r":       u.WiredRxBytesR,
		"wired-rx_packets":       u.WiredRxPackets,
		"wired-tx_bytes":         u.WiredTxBytes,
		"wired-tx_bytes-r":       u.WiredTxBytesR,
		"wired-tx_packets":       u.WiredTxPackets,
	}

	return influx.NewPoint("clients", tags, fields, time.Now())
}
